# Field Mapping Matrix

This document provides a comprehensive mapping of field names across all layers of the Sve Tu Platform.

## User Entity

| Field | PostgreSQL | Go Struct | Go JSON | TypeScript | OpenSearch | Notes |
|-------|------------|-----------|----------|------------|------------|-------|
| Primary Key | `id` | `ID` | `id` | `id` | `id` | integer/int64/number |
| Email | `email` | `Email` | `email` | `email` | `email` | string |
| Name | `name` | `Name` | `name` | `name` | `name` | string |
| Phone | `phone` | `Phone` | `phone` | `phone` | `phone` | string |
| Picture URL | `picture_url` | `PictureURL` | `picture_url` | `picture_url` | `picture_url` | string |
| Google ID | `google_id` | `GoogleID` | `google_id` | `google_id` | - | string |
| Provider | `provider` | `Provider` | `provider` | `provider` | - | string |
| Created At | `created_at` | `CreatedAt` | `created_at` | `created_at` | `created_at` | timestamp |

## MarketplaceListing Entity

| Field | PostgreSQL | Go Struct | Go JSON | TypeScript | OpenSearch | Notes |
|-------|------------|-----------|----------|------------|------------|-------|
| Primary Key | `id` | `ID` | `id` | `id` | `id` | integer/int64/number |
| User ID | `user_id` | `UserID` | `user_id` | `user_id` | `user_id` | integer |
| Category ID | `category_id` | `CategoryID` | `category_id` | `category_id` | `category_id` | integer |
| Storefront ID | `storefront_id` | `StorefrontID` | `storefront_id` | `storefront_id` | `storefront_id` | integer |
| Title | `title` | `Title` | `title` | `title` | `title` | string |
| Description | `description` | `Description` | `description` | `description` | `description` | text |
| Price | `price` | `Price` | `price` | `price` | `price` | numeric/float64/number |
| Condition | `condition` | `Condition` | `condition` | `condition` | `condition` | string |
| Status | `status` | `Status` | `status` | `status` | `status` | string |
| Location | `location` | `Location` | `location` | `location` | `location` | string |
| Latitude | `latitude` | `Latitude` | `latitude` | `latitude` | `latitude` | numeric/float64/number |
| Longitude | `longitude` | `Longitude` | `longitude` | `longitude` | `longitude` | numeric/float64/number |
| City | `city` | `City` | `city` | `city` | `city` | string (was address_city) |
| Country | `country` | `Country` | `country` | `country` | `country` | string (was address_country) |
| Views Count | `views_count` | `ViewsCount` | `views_count` | `views_count` | `views_count` | integer |
| Show on Map | `show_on_map` | `ShowOnMap` | `show_on_map` | `show_on_map` | `show_on_map` | boolean |
| Created At | `created_at` | `CreatedAt` | `created_at` | `created_at` | `created_at` | timestamp |
| Updated At | `updated_at` | `UpdatedAt` | `updated_at` | `updated_at` | `updated_at` | timestamp |
| External ID | `external_id` | `ExternalID` | `external_id` | `external_id` | - | string |
| Metadata | `metadata` | `Metadata` | `metadata` | `metadata` | - | jsonb/map/object |

## Common Patterns

1. **ID Fields**: Always use `_id` suffix in database and JSON (e.g., `user_id`, not `userId`)
2. **Timestamps**: Always use `_at` suffix (e.g., `created_at`, `updated_at`)
3. **Booleans**: Use `is_` or `has_` prefix (e.g., `is_active`, `has_discount`)
4. **URLs**: Use `_url` suffix (e.g., `picture_url`, `logo_url`)
5. **Counts**: Use `_count` suffix (e.g., `views_count`, `review_count`)

## Migration History

- **Migration 000093**: Renamed `address_city` → `city` and `address_country` → `country` in `marketplace_listings` table
EOF < /dev/null