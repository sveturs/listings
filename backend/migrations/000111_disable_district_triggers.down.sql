-- Re-enable automatic district/municipality assignment trigger

-- Enable the trigger
ALTER TABLE listings_geo ENABLE TRIGGER trigger_assign_district_municipality;

-- Remove comment
COMMENT ON TRIGGER trigger_assign_district_municipality ON listings_geo IS NULL;