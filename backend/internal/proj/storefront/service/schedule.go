// backend/internal/proj/storefront/service/schedule.go
package service

import (
	"backend/internal/storage"
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ScheduleService сервис для работы с расписанием импорта
type ScheduleService struct {
	storage    storage.Storage
	storefront StorefrontServiceInterface
	ticker     *time.Ticker
	stopCh     chan struct{}
	wg         sync.WaitGroup
	running    bool
	mutex      sync.Mutex
}

// NewScheduleService создает новый сервис планирования
func NewScheduleService(storage storage.Storage, storefront StorefrontServiceInterface) *ScheduleService {
	return &ScheduleService{
		storage:    storage,
		storefront: storefront,
		stopCh:     make(chan struct{}),
	}
}

// Start запускает фоновую задачу проверки расписания
func (s *ScheduleService) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		return
	}

	// Запускаем тикер каждые 15 минут
	s.ticker = time.NewTicker(15 * time.Minute)
	s.running = true

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// Запускаем проверку сразу после старта
		s.checkSchedules()

		for {
			select {
			case <-s.ticker.C:
				s.checkSchedules()
			case <-s.stopCh:
				s.ticker.Stop()
				return
			}
		}
	}()

	log.Println("Schedule service started")
}

// Stop останавливает фоновую задачу
func (s *ScheduleService) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return
	}

	close(s.stopCh)
	s.wg.Wait()
	s.running = false

	log.Println("Schedule service stopped")
}

// checkSchedules проверяет все источники импорта и запускает импорт по расписанию
func (s *ScheduleService) checkSchedules() {
	log.Println("Checking import schedules...")

	// Получаем все источники импорта с расписанием
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Исполняем SQL-запрос для получения источников с расписанием
	query := `
		SELECT 
			s.id, s.storefront_id, s.type, s.url, s.schedule,
			COALESCE(s.last_import_at, '1970-01-01'::timestamp) as last_import,
			sf.user_id
		FROM 
			import_sources s
			JOIN user_storefronts sf ON s.storefront_id = sf.id
		WHERE 
			s.schedule IS NOT NULL 
			AND s.schedule != ''
			AND s.url IS NOT NULL
			AND s.url != ''
	`

	rows, err := s.storage.Query(ctx, query)
	if err != nil {
		log.Printf("Failed to fetch scheduled imports: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id           int
			storefrontID int
			importType   string
			url          string
			schedule     string
			lastImport   time.Time
			userID       int
		)

		if err := rows.Scan(&id, &storefrontID, &importType, &url, &schedule, &lastImport, &userID); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Проверяем, нужно ли запустить импорт по расписанию
		if shouldRunImport(schedule, lastImport) {
			log.Printf("Running scheduled import for source ID: %d", id)

			// Запускаем импорт в отдельной горутине
			go s.runScheduledImport(id, url, userID)
		}
	}

	if err := rows.Scan(); err != nil && err.Error() != "EOF" {
		log.Printf("Error after scanning rows: %v", err)
	}
}

// shouldRunImport проверяет, нужно ли запустить импорт по расписанию
func shouldRunImport(schedule string, lastImport time.Time) bool {
	// Поддерживаем простые форматы расписания:
	// - hourly: каждый час
	// - daily: каждый день
	// - weekly: каждую неделю
	// - monthly: каждый месяц

	now := time.Now()
	var interval time.Duration

	switch strings.ToLower(schedule) {
	case "hourly":
		interval = time.Hour
	case "daily":
		interval = 24 * time.Hour
	case "weekly":
		interval = 7 * 24 * time.Hour
	case "monthly":
		// Приблизительно 30 дней
		interval = 30 * 24 * time.Hour
	default:
		// Если формат расписания неизвестен, возвращаем false
		return false
	}

	// Проверяем, прошло ли указанное время с последнего импорта
	return now.Sub(lastImport) >= interval
}

// runScheduledImport запускает импорт по расписанию
func (s *ScheduleService) runScheduledImport(sourceID int, url string, userID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Проверяем, если URL оканчивается на .zip, предполагаем, что это XML в ZIP
	if strings.HasSuffix(strings.ToLower(url), ".zip") {
		// Загружаем ZIP-архив
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error downloading ZIP from URL for scheduled import: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Bad status when downloading ZIP for scheduled import: %s", resp.Status)
			return
		}

		_, err = s.storefront.ImportXMLFromZip(ctx, sourceID, resp.Body, userID)
		if err != nil {
			log.Printf("Error importing XML from ZIP URL for scheduled import: %v", err)
			return
		}
	} else {
		_, err := s.storefront.RunImport(ctx, sourceID, userID)
		if err != nil {
			log.Printf("Error running import from URL for scheduled import: %v", err)
			return
		}
	}

	log.Printf("Scheduled import for source ID %d completed successfully", sourceID)
}