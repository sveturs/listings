package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/domain/models"
)

// SetWorkingHours устанавливает часы работы витрины
func (r *storefrontRepo) SetWorkingHours(ctx context.Context, hours []*models.StorefrontHours) error {
	if len(hours) == 0 {
		return nil
	}

	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибку если транзакция уже была завершена
		}
	}()

	// Удаляем старые часы работы
	storefrontID := hours[0].StorefrontID
	_, err = tx.Exec(ctx, "DELETE FROM storefront_hours WHERE storefront_id = $1", storefrontID)
	if err != nil {
		return fmt.Errorf("failed to delete old hours: %w", err)
	}

	// Вставляем новые
	for _, h := range hours {
		_, err = tx.Exec(ctx, `
			INSERT INTO storefront_hours (
				storefront_id, day_of_week, open_time, close_time, is_closed, special_date, special_note
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, h.StorefrontID, h.DayOfWeek, h.OpenTime, h.CloseTime, h.IsClosed, h.SpecialDate, h.SpecialNote)
		if err != nil {
			return fmt.Errorf("failed to insert hours: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetWorkingHours получает часы работы витрины
func (r *storefrontRepo) GetWorkingHours(ctx context.Context, storefrontID int) ([]*models.StorefrontHours, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, storefront_id, day_of_week, open_time, close_time, is_closed, special_date, special_note
		FROM storefront_hours
		WHERE storefront_id = $1
		ORDER BY day_of_week, special_date
	`, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get working hours: %w", err)
	}
	defer rows.Close()

	var hours []*models.StorefrontHours
	for rows.Next() {
		h := &models.StorefrontHours{}
		err := rows.Scan(&h.ID, &h.StorefrontID, &h.DayOfWeek, &h.OpenTime, &h.CloseTime,
			&h.IsClosed, &h.SpecialDate, &h.SpecialNote)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hours: %w", err)
		}
		hours = append(hours, h)
	}

	return hours, nil
}

// IsOpenNow проверяет открыта ли витрина сейчас
func (r *storefrontRepo) IsOpenNow(ctx context.Context, storefrontID int) (bool, error) {
	var isOpen bool
	err := r.db.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM storefront_hours
			WHERE storefront_id = $1
			AND (
				-- Обычные часы работы
				(special_date IS NULL 
				 AND day_of_week = EXTRACT(DOW FROM CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Belgrade')
				 AND NOT is_closed
				 AND CURRENT_TIME AT TIME ZONE 'Europe/Belgrade' BETWEEN open_time AND close_time)
				OR
				-- Специальные часы на сегодня
				(special_date = CURRENT_DATE
				 AND NOT is_closed
				 AND CURRENT_TIME AT TIME ZONE 'Europe/Belgrade' BETWEEN open_time AND close_time)
			)
		)
	`, storefrontID).Scan(&isOpen)

	return isOpen, err
}

// SetPaymentMethods устанавливает методы оплаты
func (r *storefrontRepo) SetPaymentMethods(ctx context.Context, methods []*models.StorefrontPaymentMethod) error {
	if len(methods) == 0 {
		return nil
	}

	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибку если транзакция уже была завершена
		}
	}()

	storefrontID := methods[0].StorefrontID

	// Обновляем существующие или вставляем новые
	for _, m := range methods {
		settingsJSON, _ := json.Marshal(m.Settings)

		_, err = tx.Exec(ctx, `
			INSERT INTO storefront_payment_methods (
				storefront_id, method_type, is_enabled, provider, settings,
				transaction_fee, min_amount, max_amount
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (storefront_id, method_type)
			DO UPDATE SET
				is_enabled = EXCLUDED.is_enabled,
				provider = EXCLUDED.provider,
				settings = EXCLUDED.settings,
				transaction_fee = EXCLUDED.transaction_fee,
				min_amount = EXCLUDED.min_amount,
				max_amount = EXCLUDED.max_amount
		`, storefrontID, m.MethodType, m.IsEnabled, m.Provider, settingsJSON,
			m.TransactionFee, m.MinAmount, m.MaxAmount)
		if err != nil {
			return fmt.Errorf("failed to upsert payment method: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetPaymentMethods получает методы оплаты витрины
func (r *storefrontRepo) GetPaymentMethods(ctx context.Context, storefrontID int) ([]*models.StorefrontPaymentMethod, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, storefront_id, method_type, is_enabled, provider, settings,
			   transaction_fee, min_amount, max_amount, created_at
		FROM storefront_payment_methods
		WHERE storefront_id = $1
		ORDER BY method_type
	`, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	defer rows.Close()

	var methods []*models.StorefrontPaymentMethod
	for rows.Next() {
		m := &models.StorefrontPaymentMethod{}
		var settings []byte

		err := rows.Scan(&m.ID, &m.StorefrontID, &m.MethodType, &m.IsEnabled,
			&m.Provider, &settings, &m.TransactionFee, &m.MinAmount, &m.MaxAmount, &m.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}

		if settings != nil {
			if err := json.Unmarshal(settings, &m.Settings); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}

		methods = append(methods, m)
	}

	return methods, nil
}

// SetDeliveryOptions устанавливает опции доставки
func (r *storefrontRepo) SetDeliveryOptions(ctx context.Context, options []*models.StorefrontDeliveryOption) error {
	if len(options) == 0 {
		return nil
	}

	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибку если транзакция уже была завершена
		}
	}()

	for _, opt := range options {
		zonesJSON, _ := json.Marshal(opt.Zones)
		availableDaysJSON, _ := json.Marshal(opt.AvailableDays)
		providerConfigJSON, _ := json.Marshal(opt.ProviderConfig)
		supportedPaymentsJSON, _ := json.Marshal(opt.SupportedPaymentMethods)

		if opt.ID > 0 {
			// Обновляем существующую
			_, err = tx.Exec(ctx, `
				UPDATE storefront_delivery_options SET
					name = $2, description = $3, provider = $4,
					base_price = $5, price_per_km = $6, price_per_kg = $7, free_above_amount = $8,
					cod_fee = $9, insurance_fee = $10, fragile_handling = $11,
					min_order_amount = $12, max_weight_kg = $13, max_distance_km = $14,
					estimated_days_min = $15, estimated_days_max = $16,
					zones = $17, available_days = $18, cutoff_time = $19,
					supported_payment_methods = $20, provider_config = $21,
					is_enabled = $22, updated_at = CURRENT_TIMESTAMP
				WHERE id = $1
			`, opt.ID, opt.Name, opt.Description, opt.Provider,
				opt.BasePrice, opt.PricePerKm, opt.PricePerKg, opt.FreeAboveAmount,
				opt.CODFee, opt.InsuranceFee, opt.FragileHandling,
				opt.MinOrderAmount, opt.MaxWeightKg, opt.MaxDistanceKm,
				opt.EstimatedDaysMin, opt.EstimatedDaysMax,
				zonesJSON, availableDaysJSON, opt.CutoffTime,
				supportedPaymentsJSON, providerConfigJSON, opt.IsEnabled)
		} else {
			// Создаем новую
			_, err = tx.Exec(ctx, `
				INSERT INTO storefront_delivery_options (
					storefront_id, name, description, provider,
					base_price, price_per_km, price_per_kg, free_above_amount,
					cod_fee, insurance_fee, fragile_handling,
					min_order_amount, max_weight_kg, max_distance_km,
					estimated_days_min, estimated_days_max,
					zones, available_days, cutoff_time,
					supported_payment_methods, provider_config, is_enabled
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
			`, opt.StorefrontID, opt.Name, opt.Description, opt.Provider,
				opt.BasePrice, opt.PricePerKm, opt.PricePerKg, opt.FreeAboveAmount,
				opt.CODFee, opt.InsuranceFee, opt.FragileHandling,
				opt.MinOrderAmount, opt.MaxWeightKg, opt.MaxDistanceKm,
				opt.EstimatedDaysMin, opt.EstimatedDaysMax,
				zonesJSON, availableDaysJSON, opt.CutoffTime,
				supportedPaymentsJSON, providerConfigJSON, opt.IsEnabled)
		}

		if err != nil {
			return fmt.Errorf("failed to save delivery option: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetDeliveryOptions получает опции доставки
func (r *storefrontRepo) GetDeliveryOptions(ctx context.Context, storefrontID int) ([]*models.StorefrontDeliveryOption, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, storefront_id, name, description, provider,
			   base_price, price_per_km, price_per_kg, free_above_amount,
			   cod_fee, insurance_fee, fragile_handling,
			   min_order_amount, max_weight_kg, max_distance_km,
			   estimated_days_min, estimated_days_max,
			   zones, available_days, cutoff_time,
			   supported_payment_methods, provider_config,
			   is_enabled, created_at, updated_at
		FROM storefront_delivery_options
		WHERE storefront_id = $1
		ORDER BY name
	`, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery options: %w", err)
	}
	defer rows.Close()

	var options []*models.StorefrontDeliveryOption
	for rows.Next() {
		opt := &models.StorefrontDeliveryOption{}
		var zonesJSON, availableDaysJSON, providerConfigJSON, supportedPaymentsJSON []byte

		err := rows.Scan(
			&opt.ID, &opt.StorefrontID, &opt.Name, &opt.Description, &opt.Provider,
			&opt.BasePrice, &opt.PricePerKm, &opt.PricePerKg, &opt.FreeAboveAmount,
			&opt.CODFee, &opt.InsuranceFee, &opt.FragileHandling,
			&opt.MinOrderAmount, &opt.MaxWeightKg, &opt.MaxDistanceKm,
			&opt.EstimatedDaysMin, &opt.EstimatedDaysMax,
			&zonesJSON, &availableDaysJSON, &opt.CutoffTime,
			&supportedPaymentsJSON, &providerConfigJSON,
			&opt.IsEnabled, &opt.CreatedAt, &opt.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan delivery option: %w", err)
		}

		// Парсим JSON поля
		if zonesJSON != nil {
			if err := json.Unmarshal(zonesJSON, &opt.Zones); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}
		if availableDaysJSON != nil {
			if err := json.Unmarshal(availableDaysJSON, &opt.AvailableDays); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}
		if providerConfigJSON != nil {
			if err := json.Unmarshal(providerConfigJSON, &opt.ProviderConfig); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}
		if supportedPaymentsJSON != nil {
			if err := json.Unmarshal(supportedPaymentsJSON, &opt.SupportedPaymentMethods); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}

		options = append(options, opt)
	}

	return options, nil
}

// AddStaff добавляет сотрудника
func (r *storefrontRepo) AddStaff(ctx context.Context, staff *models.StorefrontStaff) error {
	permissionsJSON, _ := json.Marshal(staff.Permissions)

	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO storefront_staff (storefront_id, user_id, role, permissions)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (storefront_id, user_id) 
		DO UPDATE SET role = EXCLUDED.role, permissions = EXCLUDED.permissions, updated_at = CURRENT_TIMESTAMP
	`, staff.StorefrontID, staff.UserID, staff.Role, permissionsJSON)

	return err
}

// UpdateStaff обновляет права сотрудника
func (r *storefrontRepo) UpdateStaff(ctx context.Context, id int, permissions models.JSONB) error {
	permissionsJSON, _ := json.Marshal(permissions)

	result, err := r.db.pool.Exec(ctx, `
		UPDATE storefront_staff 
		SET permissions = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id, permissionsJSON)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// RemoveStaff удаляет сотрудника
func (r *storefrontRepo) RemoveStaff(ctx context.Context, storefrontID, userID int) error {
	result, err := r.db.pool.Exec(ctx, `
		DELETE FROM storefront_staff 
		WHERE storefront_id = $1 AND user_id = $2 AND role != 'owner'
	`, storefrontID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// GetStaff получает список сотрудников
func (r *storefrontRepo) GetStaff(ctx context.Context, storefrontID int) ([]*models.StorefrontStaff, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, storefront_id, user_id, role, permissions, last_active_at, actions_count, created_at, updated_at
		FROM storefront_staff
		WHERE storefront_id = $1
		ORDER BY role, created_at
	`, storefrontID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staff []*models.StorefrontStaff
	for rows.Next() {
		s := &models.StorefrontStaff{}
		var permissionsJSON []byte

		err := rows.Scan(&s.ID, &s.StorefrontID, &s.UserID, &s.Role, &permissionsJSON,
			&s.LastActiveAt, &s.ActionsCount, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if permissionsJSON != nil {
			if err := json.Unmarshal(permissionsJSON, &s.Permissions); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}

		staff = append(staff, s)
	}

	return staff, nil
}

// RecordView записывает просмотр витрины
func (r *storefrontRepo) RecordView(ctx context.Context, storefrontID int) error {
	_, err := r.db.pool.Exec(ctx, `
		UPDATE storefronts 
		SET views_count = views_count + 1
		WHERE id = $1
	`, storefrontID)

	return err
}

// RecordAnalytics записывает аналитику
func (r *storefrontRepo) RecordAnalytics(ctx context.Context, analytics *models.StorefrontAnalytics) error {
	trafficJSON, _ := json.Marshal(analytics.TrafficSources)
	topProductsJSON, _ := json.Marshal(analytics.TopProducts)
	topCategoriesJSON, _ := json.Marshal(analytics.TopCategories)
	paymentMethodsJSON, _ := json.Marshal(analytics.PaymentMethodsUsage)
	ordersByCityJSON, _ := json.Marshal(analytics.OrdersByCity)

	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO storefront_analytics (
			storefront_id, date,
			page_views, unique_visitors, bounce_rate, avg_session_time,
			orders_count, revenue, avg_order_value, conversion_rate,
			payment_methods_usage,
			product_views, add_to_cart_count, checkout_count,
			traffic_sources, top_products, top_categories, orders_by_city
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		ON CONFLICT (storefront_id, date)
		DO UPDATE SET
			page_views = storefront_analytics.page_views + EXCLUDED.page_views,
			unique_visitors = GREATEST(storefront_analytics.unique_visitors, EXCLUDED.unique_visitors),
			bounce_rate = (storefront_analytics.bounce_rate + EXCLUDED.bounce_rate) / 2,
			avg_session_time = (storefront_analytics.avg_session_time + EXCLUDED.avg_session_time) / 2,
			orders_count = storefront_analytics.orders_count + EXCLUDED.orders_count,
			revenue = storefront_analytics.revenue + EXCLUDED.revenue,
			avg_order_value = CASE 
				WHEN (storefront_analytics.orders_count + EXCLUDED.orders_count) > 0 
				THEN (storefront_analytics.revenue + EXCLUDED.revenue) / (storefront_analytics.orders_count + EXCLUDED.orders_count)
				ELSE 0 
			END,
			conversion_rate = CASE 
				WHEN (storefront_analytics.unique_visitors + EXCLUDED.unique_visitors) > 0 
				THEN (storefront_analytics.orders_count + EXCLUDED.orders_count)::float / (storefront_analytics.unique_visitors + EXCLUDED.unique_visitors) * 100
				ELSE 0 
			END,
			payment_methods_usage = EXCLUDED.payment_methods_usage,
			product_views = storefront_analytics.product_views + EXCLUDED.product_views,
			add_to_cart_count = storefront_analytics.add_to_cart_count + EXCLUDED.add_to_cart_count,
			checkout_count = storefront_analytics.checkout_count + EXCLUDED.checkout_count,
			traffic_sources = EXCLUDED.traffic_sources,
			top_products = EXCLUDED.top_products,
			top_categories = EXCLUDED.top_categories,
			orders_by_city = EXCLUDED.orders_by_city
	`,
		analytics.StorefrontID, analytics.Date,
		analytics.PageViews, analytics.UniqueVisitors, analytics.BounceRate, analytics.AvgSessionTime,
		analytics.OrdersCount, analytics.Revenue, analytics.AvgOrderValue, analytics.ConversionRate,
		paymentMethodsJSON,
		analytics.ProductViews, analytics.AddToCartCount, analytics.CheckoutCount,
		trafficJSON, topProductsJSON, topCategoriesJSON, ordersByCityJSON,
	)

	return err
}

// GetAnalytics получает аналитику за период
func (r *storefrontRepo) GetAnalytics(ctx context.Context, storefrontID int, from, to time.Time) ([]*models.StorefrontAnalytics, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT 
			id, storefront_id, date,
			page_views, unique_visitors, bounce_rate, avg_session_time,
			orders_count, revenue, avg_order_value, conversion_rate,
			payment_methods_usage,
			product_views, add_to_cart_count, checkout_count,
			traffic_sources, top_products, top_categories, orders_by_city,
			created_at
		FROM storefront_analytics
		WHERE storefront_id = $1 AND date BETWEEN $2 AND $3
		ORDER BY date DESC
	`, storefrontID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []*models.StorefrontAnalytics
	for rows.Next() {
		a := &models.StorefrontAnalytics{}
		var trafficJSON, topProductsJSON, topCategoriesJSON, paymentMethodsJSON, ordersByCityJSON []byte

		err := rows.Scan(
			&a.ID, &a.StorefrontID, &a.Date,
			&a.PageViews, &a.UniqueVisitors, &a.BounceRate, &a.AvgSessionTime,
			&a.OrdersCount, &a.Revenue, &a.AvgOrderValue, &a.ConversionRate,
			&paymentMethodsJSON,
			&a.ProductViews, &a.AddToCartCount, &a.CheckoutCount,
			&trafficJSON, &topProductsJSON, &topCategoriesJSON, &ordersByCityJSON,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Парсим JSON поля
		if trafficJSON != nil {
			if err := json.Unmarshal(trafficJSON, &a.TrafficSources); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}
		if topProductsJSON != nil {
			if err := json.Unmarshal(topProductsJSON, &a.TopProducts); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}
		if topCategoriesJSON != nil {
			if err := json.Unmarshal(topCategoriesJSON, &a.TopCategories); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}
		if paymentMethodsJSON != nil {
			if err := json.Unmarshal(paymentMethodsJSON, &a.PaymentMethodsUsage); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}
		if ordersByCityJSON != nil {
			if err := json.Unmarshal(ordersByCityJSON, &a.OrdersByCity); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}

		analytics = append(analytics, a)
	}

	return analytics, nil
}

// GetClusters получает кластеры для карты
func (r *storefrontRepo) GetClusters(ctx context.Context, bounds GeoBounds, zoomLevel int) ([]*models.MapCluster, error) {
	// Определяем размер кластера в зависимости от зума
	var clusterSize float64
	switch {
	case zoomLevel <= 10:
		clusterSize = 1.0 // 1 градус
	case zoomLevel <= 14:
		clusterSize = 0.1 // 0.1 градуса
	case zoomLevel <= 16:
		clusterSize = 0.01 // 0.01 градуса
	default:
		clusterSize = 0.001 // 0.001 градуса
	}

	rows, err := r.db.pool.Query(ctx, `
		SELECT 
			ROUND(latitude/$1)*$1 as cluster_lat,
			ROUND(longitude/$1)*$1 as cluster_lng,
			COUNT(*) as count
		FROM storefronts
		WHERE is_active = true
		AND latitude BETWEEN $2 AND $3
		AND longitude BETWEEN $4 AND $5
		GROUP BY cluster_lat, cluster_lng
		HAVING COUNT(*) > 1
	`, clusterSize, bounds.MinLat, bounds.MaxLat, bounds.MinLng, bounds.MaxLng)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []*models.MapCluster
	for rows.Next() {
		c := &models.MapCluster{}
		err := rows.Scan(&c.Latitude, &c.Longitude, &c.Count)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, c)
	}

	return clusters, nil
}
