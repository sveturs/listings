-- Temporarily disable automatic district/municipality assignment trigger
-- This is done to reduce database load while district functionality is disabled

-- Disable the trigger (but keep the function for future use)
ALTER TABLE listings_geo DISABLE TRIGGER trigger_assign_district_municipality;

-- Comment explaining the temporary nature
COMMENT ON TRIGGER trigger_assign_district_municipality ON listings_geo IS 
'TEMPORARILY DISABLED: Automatic district/municipality assignment disabled to reduce DB load while district search is not in use';