package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Translation struct {
	EntityType string
	EntityID   string
	FieldName  string
	Language   string
	Text       string
}

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	ctx := context.Background()

	// Create translations for categories
	log.Println("Creating translations for categories...")
	if err := createCategoryTranslations(ctx, db); err != nil {
		log.Printf("Error creating category translations: %v", err)
	}

	// Create translations for attributes
	log.Println("Creating translations for attributes...")
	if err := createAttributeTranslations(ctx, db); err != nil {
		log.Printf("Error creating attribute translations: %v", err)
	}

	log.Println("Translation creation completed!")
}

func createCategoryTranslations(ctx context.Context, db *sql.DB) error {
	// Get all categories
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, description, seo_title, seo_description
		FROM marketplace_categories
		ORDER BY id
	`)
	if err != nil {
		return fmt.Errorf("failed to query categories: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	translations := []Translation{}

	for rows.Next() {
		var id int
		var name, description, seoTitle, seoDescription sql.NullString

		if err := rows.Scan(&id, &name, &description, &seoTitle, &seoDescription); err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}

		idStr := fmt.Sprintf("%d", id)

		// Create translations for each field and language
		languages := []string{"en", "ru", "sr"}

		for _, lang := range languages {
			if name.Valid {
				translations = append(translations, Translation{
					EntityType: "category",
					EntityID:   idStr,
					FieldName:  "name",
					Language:   lang,
					Text:       translateText(name.String, lang),
				})
			}

			if description.Valid && description.String != "" {
				translations = append(translations, Translation{
					EntityType: "category",
					EntityID:   idStr,
					FieldName:  "description",
					Language:   lang,
					Text:       translateText(description.String, lang),
				})
			}

			if seoTitle.Valid && seoTitle.String != "" {
				translations = append(translations, Translation{
					EntityType: "category",
					EntityID:   idStr,
					FieldName:  "seo_title",
					Language:   lang,
					Text:       translateText(seoTitle.String, lang),
				})
			}

			if seoDescription.Valid && seoDescription.String != "" {
				translations = append(translations, Translation{
					EntityType: "category",
					EntityID:   idStr,
					FieldName:  "seo_description",
					Language:   lang,
					Text:       translateText(seoDescription.String, lang),
				})
			}
		}
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over categories: %w", err)
	}

	// Insert translations
	for _, t := range translations {
		_, err := db.ExecContext(ctx, `
			INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, is_machine_translated, is_verified, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, true, false, $6, $6)
			ON CONFLICT (entity_type, entity_id, language, field_name) 
			DO UPDATE SET translated_text = $5, updated_at = $6
		`, t.EntityType, t.EntityID, t.FieldName, t.Language, t.Text, time.Now())

		if err != nil {
			log.Printf("Error inserting translation for category %s, field %s, lang %s: %v", t.EntityID, t.FieldName, t.Language, err)
		} else {
			log.Printf("Created translation for category %s, field %s, lang %s", t.EntityID, t.FieldName, t.Language)
		}
	}

	return nil
}

func createAttributeTranslations(ctx context.Context, db *sql.DB) error {
	// Get all attributes
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, display_name
		FROM category_attributes
		ORDER BY id
	`)
	if err != nil {
		return fmt.Errorf("failed to query attributes: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	translations := []Translation{}

	for rows.Next() {
		var id int
		var name, displayName string

		if err := rows.Scan(&id, &name, &displayName); err != nil {
			log.Printf("Error scanning attribute: %v", err)
			continue
		}

		idStr := fmt.Sprintf("%d", id)

		// Create translations for each language
		languages := []string{"en", "ru", "sr"}

		for _, lang := range languages {
			translations = append(translations, Translation{
				EntityType: "attribute",
				EntityID:   idStr,
				FieldName:  "display_name",
				Language:   lang,
				Text:       translateText(displayName, lang),
			})
		}
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over attributes: %w", err)
	}

	// Insert translations
	for _, t := range translations {
		_, err := db.ExecContext(ctx, `
			INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, is_machine_translated, is_verified, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, true, false, $6, $6)
			ON CONFLICT (entity_type, entity_id, language, field_name) 
			DO UPDATE SET translated_text = $5, updated_at = $6
		`, t.EntityType, t.EntityID, t.FieldName, t.Language, t.Text, time.Now())

		if err != nil {
			log.Printf("Error inserting translation for attribute %s, lang %s: %v", t.EntityID, t.Language, err)
		} else {
			log.Printf("Created translation for attribute %s, lang %s", t.EntityID, t.Language)
		}
	}

	return nil
}

// Simple translation function - in real implementation should use Google Translate API
func translateText(text string, targetLang string) string {
	// For demonstration, we'll just return the original text
	// In production, this should call the actual translation service

	// Map of simple translations for common terms
	translations := map[string]map[string]string{
		// Main Categories (50 categories)
		"Elektronika": {
			"en": "Electronics",
			"ru": "Электроника",
			"sr": "Електроника",
		},
		"Moda": {
			"en": "Fashion",
			"ru": "Мода",
			"sr": "Мода",
		},
		"Automobili": {
			"en": "Automotive",
			"ru": "Автомобили",
			"sr": "Аутомобили",
		},
		"Nekretnine": {
			"en": "Real Estate",
			"ru": "Недвижимость",
			"sr": "Некретнине",
		},
		"Dom i bašta": {
			"en": "Home & Garden",
			"ru": "Дом и сад",
			"sr": "Дом и башта",
		},
		"Poljoprivreda": {
			"en": "Agriculture",
			"ru": "Сельское хозяйство",
			"sr": "Пољопривреда",
		},
		"Industrija": {
			"en": "Industry",
			"ru": "Промышленность",
			"sr": "Индустрија",
		},
		"Hrana i piće": {
			"en": "Food & Beverages",
			"ru": "Еда и напитки",
			"sr": "Храна и пиће",
		},
		"Usluge": {
			"en": "Services",
			"ru": "Услуги",
			"sr": "Услуге",
		},
		"Sport i rekreacija": {
			"en": "Sports & Recreation",
			"ru": "Спорт и отдых",
			"sr": "Спорт и рекреација",
		},
		"Kućni ljubimci": {
			"en": "Pets",
			"ru": "Домашние животные",
			"sr": "Кућни љубимци",
		},
		"Deca": {
			"en": "Children",
			"ru": "Дети",
			"sr": "Деца",
		},
		"Knjige": {
			"en": "Books",
			"ru": "Книги",
			"sr": "Књиге",
		},
		"Zdravlje i lepota": {
			"en": "Health & Beauty",
			"ru": "Здоровье и красота",
			"sr": "Здравље и лепота",
		},
		"Umetnička dela": {
			"en": "Art",
			"ru": "Искусство",
			"sr": "Уметничка дела",
		},
		"Muzika": {
			"en": "Music",
			"ru": "Музыка",
			"sr": "Музика",
		},
		"Poslovi": {
			"en": "Jobs",
			"ru": "Работа",
			"sr": "Послови",
		},
		"Obrazovanje": {
			"en": "Education",
			"ru": "Образование",
			"sr": "Образовање",
		},
		"Turizam": {
			"en": "Tourism",
			"ru": "Туризм",
			"sr": "Туризам",
		},
		"Antikviteti": {
			"en": "Antiques",
			"ru": "Антиквариат",
			"sr": "Антиквитети",
		},
		// Subcategories
		"Pametni telefoni": {
			"en": "Smartphones",
			"ru": "Смартфоны",
			"sr": "Паметни телефони",
		},
		"Računari": {
			"en": "Computers",
			"ru": "Компьютеры",
			"sr": "Рачунари",
		},
		"TV i audio": {
			"en": "TV & Audio",
			"ru": "ТВ и аудио",
			"sr": "ТВ и аудио",
		},
		"Fotoaparati": {
			"en": "Cameras",
			"ru": "Фотоаппараты",
			"sr": "Фотоапарати",
		},
		"Igre i konzole": {
			"en": "Games & Consoles",
			"ru": "Игры и консоли",
			"sr": "Игре и конзоле",
		},
		"Tableti": {
			"en": "Tablets",
			"ru": "Планшеты",
			"sr": "Таблети",
		},
		"Muška odeća": {
			"en": "Men's Clothing",
			"ru": "Мужская одежда",
			"sr": "Мушка одећа",
		},
		"Ženska odeća": {
			"en": "Women's Clothing",
			"ru": "Женская одежда",
			"sr": "Женска одећа",
		},
		"Obuća": {
			"en": "Footwear",
			"ru": "Обувь",
			"sr": "Обућа",
		},
		"Stanovi": {
			"en": "Apartments",
			"ru": "Квартиры",
			"sr": "Станови",
		},
		"Kuće": {
			"en": "Houses",
			"ru": "Дома",
			"sr": "Куће",
		},
		"Zemljišta": {
			"en": "Land",
			"ru": "Земельные участки",
			"sr": "Земљишта",
		},
		"Nameštaj": {
			"en": "Furniture",
			"ru": "Мебель",
			"sr": "Намештај",
		},
		"Alati": {
			"en": "Tools",
			"ru": "Инструменты",
			"sr": "Алати",
		},
		"Bela tehnika": {
			"en": "Appliances",
			"ru": "Бытовая техника",
			"sr": "Бела техника",
		},
		"Građevinski materijal": {
			"en": "Building Materials",
			"ru": "Строительные материалы",
			"sr": "Грађевински материјал",
		},
		"Poljoprivredne mašine": {
			"en": "Agricultural Machinery",
			"ru": "Сельскохозяйственная техника",
			"sr": "Пољопривредне машине",
		},
		"Stoka": {
			"en": "Livestock",
			"ru": "Скот",
			"sr": "Стока",
		},
		"Semena i sadnice": {
			"en": "Seeds & Seedlings",
			"ru": "Семена и саженцы",
			"sr": "Семена и саднице",
		},
		"Popravke i održavanje": {
			"en": "Repairs & Maintenance",
			"ru": "Ремонт и обслуживание",
			"sr": "Поправке и одржавање",
		},
		"Transport": {
			"en": "Transportation",
			"ru": "Транспорт",
			"sr": "Транспорт",
		},
		"Čišćenje": {
			"en": "Cleaning",
			"ru": "Уборка",
			"sr": "Чишћење",
		},
		"Sportska oprema": {
			"en": "Sports Equipment",
			"ru": "Спортивное снаряжение",
			"sr": "Спортска опрема",
		},
		"Fitness": {
			"en": "Fitness",
			"ru": "Фитнес",
			"sr": "Фитнес",
		},
		"Bicikli": {
			"en": "Bicycles",
			"ru": "Велосипеды",
			"sr": "Бицикли",
		},
		"Psi": {
			"en": "Dogs",
			"ru": "Собаки",
			"sr": "Пси",
		},
		"Mačke": {
			"en": "Cats",
			"ru": "Кошки",
			"sr": "Мачке",
		},
		"Igračke": {
			"en": "Toys",
			"ru": "Игрушки",
			"sr": "Играчке",
		},
		"Oprema za bebe": {
			"en": "Baby Equipment",
			"ru": "Товары для малышей",
			"sr": "Опрема за бебе",
		},
		// All 32 Attributes
		"Cena": {
			"en": "Price",
			"ru": "Цена",
			"sr": "Цена",
		},
		"Stanje": {
			"en": "Condition",
			"ru": "Состояние",
			"sr": "Стање",
		},
		"Brend": {
			"en": "Brand",
			"ru": "Бренд",
			"sr": "Бренд",
		},
		"Boja": {
			"en": "Color",
			"ru": "Цвет",
			"sr": "Боја",
		},
		"Veličina": {
			"en": "Size",
			"ru": "Размер",
			"sr": "Величина",
		},
		"Materijal": {
			"en": "Material",
			"ru": "Материал",
			"sr": "Материјал",
		},
		"Lokacija": {
			"en": "Location",
			"ru": "Местоположение",
			"sr": "Локација",
		},
		"Godina proizvodnje": {
			"en": "Year of Manufacture",
			"ru": "Год производства",
			"sr": "Година производње",
		},
		"Kilometraža": {
			"en": "Mileage",
			"ru": "Пробег",
			"sr": "Километража",
		},
		"Gorivo": {
			"en": "Fuel",
			"ru": "Топливо",
			"sr": "Гориво",
		},
		"Snaga motora": {
			"en": "Engine Power",
			"ru": "Мощность двигателя",
			"sr": "Снага мотора",
		},
		"Tip": {
			"en": "Type",
			"ru": "Тип",
			"sr": "Тип",
		},
		"Model": {
			"en": "Model",
			"ru": "Модель",
			"sr": "Модел",
		},
		"Površina": {
			"en": "Area",
			"ru": "Площадь",
			"sr": "Површина",
		},
		"Broj soba": {
			"en": "Number of Rooms",
			"ru": "Количество комнат",
			"sr": "Број соба",
		},
		"Sprat": {
			"en": "Floor",
			"ru": "Этаж",
			"sr": "Спрат",
		},
		"Parking": {
			"en": "Parking",
			"ru": "Парковка",
			"sr": "Паркинг",
		},
		"Grejanje": {
			"en": "Heating",
			"ru": "Отопление",
			"sr": "Грејање",
		},
		"Nameštenost": {
			"en": "Furnished",
			"ru": "Меблированность",
			"sr": "Намештеност",
		},
		"Starost": {
			"en": "Age",
			"ru": "Возраст",
			"sr": "Старост",
		},
		"Pol": {
			"en": "Gender",
			"ru": "Пол",
			"sr": "Пол",
		},
		"Rasa": {
			"en": "Breed",
			"ru": "Порода",
			"sr": "Раса",
		},
		"Garancija": {
			"en": "Warranty",
			"ru": "Гарантия",
			"sr": "Гаранција",
		},
		"Memorija": {
			"en": "Memory",
			"ru": "Память",
			"sr": "Меморија",
		},
		"Procesor": {
			"en": "Processor",
			"ru": "Процессор",
			"sr": "Процесор",
		},
		"Ekran": {
			"en": "Screen",
			"ru": "Экран",
			"sr": "Екран",
		},
		"Kamera": {
			"en": "Camera",
			"ru": "Камера",
			"sr": "Камера",
		},
		"Baterija": {
			"en": "Battery",
			"ru": "Батарея",
			"sr": "Батерија",
		},
		"Težina": {
			"en": "Weight",
			"ru": "Вес",
			"sr": "Тежина",
		},
		"Dimenzije": {
			"en": "Dimensions",
			"ru": "Размеры",
			"sr": "Димензије",
		},
		"Dostava": {
			"en": "Delivery",
			"ru": "Доставка",
			"sr": "Достава",
		},
		"Zamena": {
			"en": "Exchange",
			"ru": "Обмен",
			"sr": "Замена",
		},
	}

	if trans, ok := translations[text]; ok {
		if langTrans, ok := trans[targetLang]; ok {
			return langTrans
		}
	}

	// If no translation found, return original text
	return text
}
