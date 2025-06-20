package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryTreeNode struct {
	ID             int                        `json:"id"`
	Name           string                     `json:"name"`
	Slug           string                     `json:"slug"`
	Icon           string                     `json:"icon"`
	ParentID       *int                       `json:"parent_id"`
	CreatedAt      string                     `json:"created_at"`
	Level          int                        `json:"level"`
	Path           string                     `json:"path"`
	ListingCount   int                        `json:"listing_count"`
	ChildrenCount  int                        `json:"children_count"`
	Translations   map[string]string          `json:"translations"`
	Children       []CategoryTreeNode         `json:"children"`
}

func testGetCategoryTree(pool *pgxpool.Pool) {
	log.Printf("GetCategoryTree in storage called")

	query := `
WITH RECURSIVE category_tree AS (
    SELECT 
        c.id,
        c.name,
        c.slug,
        c.icon,
        c.parent_id,
        to_char(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as created_at,
        ARRAY[c.id] as category_path,
        1 as level,
        COALESCE(clc.listing_count, 0) as listing_count,
        (SELECT COUNT(*) FROM marketplace_categories sc WHERE sc.parent_id = c.id) as children_count
    FROM marketplace_categories c
    LEFT JOIN category_listing_counts clc ON clc.category_id = c.id
    WHERE c.parent_id IS NULL

    UNION ALL

    SELECT 
        c.id,
        c.name,
        c.slug,
        c.icon,
        c.parent_id,
        to_char(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as created_at,
        ct.category_path || c.id,
        ct.level + 1,
        COALESCE(clc.listing_count, 0),
        (SELECT COUNT(*) FROM marketplace_categories sc WHERE sc.parent_id = c.id)
    FROM marketplace_categories c
    LEFT JOIN category_listing_counts clc ON clc.category_id = c.id
    INNER JOIN category_tree ct ON ct.id = c.parent_id
    WHERE ct.level < 10
),
categories_with_translations AS (
    SELECT 
        ct.*,
        COALESCE(
            jsonb_object_agg(
                t.language, 
                t.translated_text
            ) FILTER (WHERE t.language IS NOT NULL),
            '{}'::jsonb
        ) as translations
    FROM category_tree ct
    LEFT JOIN translations t ON 
        t.entity_type = 'category' 
        AND t.entity_id = ct.id 
        AND t.field_name = 'name'
    GROUP BY 
        ct.id, ct.name, ct.slug, ct.icon, ct.parent_id, 
        ct.created_at, ct.category_path, ct.level, ct.listing_count, 
        ct.children_count
)
SELECT 
    c1.id,
    c1.name,
    c1.slug,
    c1.icon,
    c1.parent_id,
    c1.created_at,
    c1.level,
    array_to_string(c1.category_path, ',') as path,
    c1.listing_count,
    c1.children_count,
    c1.translations,
    COALESCE(
        json_agg(
            json_build_object(
                'id', c2.id,
                'name', c2.name,
                'slug', c2.slug,
                'icon', c2.icon,
                'parent_id', c2.parent_id,
                'created_at', c2.created_at,
                'level', c2.level,
                'path', array_to_string(c2.category_path, ','),
                'listing_count', c2.listing_count,
                'children_count', c2.children_count,
                'translations', c2.translations
            ) ORDER BY c2.name ASC
        ) FILTER (WHERE c2.id IS NOT NULL),
        '[]'::json
    ) as children
FROM categories_with_translations c1
LEFT JOIN categories_with_translations c2 ON c2.parent_id = c1.id
GROUP BY 
    c1.id, c1.name, c1.slug, c1.icon, c1.parent_id, 
    c1.created_at, c1.level, c1.category_path, c1.listing_count,
    c1.children_count, c1.translations
ORDER BY c1.name ASC;
`

	ctx := context.Background()
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	var rootCategories []CategoryTreeNode
	
	for rows.Next() {
		var node CategoryTreeNode
		var translationsJson, childrenJson []byte
		var pathStr string
		var icon sql.NullString

		err := rows.Scan(
			&node.ID,
			&node.Name,
			&node.Slug,
			&icon,
			&node.ParentID,
			&node.CreatedAt,
			&node.Level,
			&pathStr,
			&node.ListingCount,
			&node.ChildrenCount,
			&translationsJson,
			&childrenJson,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return
		}

		// Обработка NULL icon
		if icon.Valid {
			node.Icon = icon.String
		}

		if err := json.Unmarshal(translationsJson, &node.Translations); err != nil {
			log.Printf("Error unmarshaling translations for category %d: %v", node.ID, err)
			node.Translations = make(map[string]string)
		}

		var children []CategoryTreeNode
		if err := json.Unmarshal(childrenJson, &children); err != nil {
			log.Printf("Error unmarshaling children for category %d: %v", node.ID, err)
			node.Children = make([]CategoryTreeNode, 0)
		} else {
			node.Children = children
		}

		rootCategories = append(rootCategories, node)
	}

	log.Printf("Returning %d root categories with tree", len(rootCategories))
	for _, cat := range rootCategories {
		log.Printf("Category: ID=%d, Name=%s, Children=%d", cat.ID, cat.Name, len(cat.Children))
	}
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}
	defer pool.Close()

	testGetCategoryTree(pool)
}