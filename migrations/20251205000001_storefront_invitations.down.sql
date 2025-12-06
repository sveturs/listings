-- Rollback: Remove Storefront Invitations System
-- Date: 2025-12-05

-- Drop trigger
DROP TRIGGER IF EXISTS trigger_storefront_invitations_updated_at ON storefront_invitations;
DROP FUNCTION IF EXISTS update_storefront_invitations_updated_at();

-- Remove invitation_id column from storefront_staff
ALTER TABLE storefront_staff DROP COLUMN IF EXISTS invitation_id;

-- Drop main table
DROP TABLE IF EXISTS storefront_invitations;

-- Drop enum types
DROP TYPE IF EXISTS storefront_invitation_status;
DROP TYPE IF EXISTS storefront_invitation_type;
