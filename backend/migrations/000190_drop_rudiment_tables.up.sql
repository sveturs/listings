-- Migration: Drop rudiment tables (unused/legacy)
-- Date: 2025-10-13
-- Phase 2, Task 2.11

-- Drop address change log (not used)
DROP TABLE IF EXISTS address_change_log CASCADE;

-- Drop GIS cache tables (not used, data from external APIs)
DROP TABLE IF EXISTS gis_poi_cache CASCADE;
DROP TABLE IF EXISTS gis_filter_analytics CASCADE;
DROP TABLE IF EXISTS gis_isochrone_cache CASCADE;

-- Drop old geo table (replaced by unified_geo)
DROP TABLE IF EXISTS listings_geo CASCADE;

-- Drop import sources (not used)
DROP TABLE IF EXISTS import_sources CASCADE;

-- Drop custom UI tables (not implemented)
DROP TABLE IF EXISTS custom_ui_component_usage CASCADE;
DROP TABLE IF EXISTS custom_ui_components CASCADE;
DROP TABLE IF EXISTS custom_ui_templates CASCADE;
