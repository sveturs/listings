-- Migration: Add Storefront Invitations System
-- Date: 2025-12-05
-- Purpose: Add invitation system for storefront staff
--          Supports both email invitations and shareable links

-- =====================================================
-- 1. CREATE ENUM TYPES
-- =====================================================

CREATE TYPE storefront_invitation_type AS ENUM ('email', 'link');
CREATE TYPE storefront_invitation_status AS ENUM ('pending', 'accepted', 'declined', 'expired', 'revoked');

-- =====================================================
-- 2. STOREFRONT_INVITATIONS TABLE
-- =====================================================

CREATE TABLE storefront_invitations (
    id              BIGSERIAL PRIMARY KEY,
    storefront_id   BIGINT NOT NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'staff',
    type            storefront_invitation_type NOT NULL,

    -- For type='email' invitations
    invited_email   VARCHAR(255),
    invited_user_id BIGINT,  -- cross-DB reference to users in monolith

    -- For type='link' invitations
    invite_code     VARCHAR(32) UNIQUE,
    expires_at      TIMESTAMPTZ,
    max_uses        INT,
    current_uses    INT DEFAULT 0,

    -- Invitation metadata
    invited_by_id   BIGINT NOT NULL,  -- cross-DB reference to users in monolith

    -- Status tracking
    status          storefront_invitation_status NOT NULL DEFAULT 'pending',
    comment         TEXT,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    accepted_at     TIMESTAMPTZ,
    declined_at     TIMESTAMPTZ,

    -- Foreign Keys
    CONSTRAINT fk_storefront_invitations_storefront FOREIGN KEY (storefront_id)
        REFERENCES storefronts(id) ON DELETE CASCADE,

    -- Validation Constraints
    CONSTRAINT valid_sf_email_invitation CHECK (
        type != 'email' OR invited_email IS NOT NULL
    ),
    CONSTRAINT valid_sf_link_invitation CHECK (
        type != 'link' OR invite_code IS NOT NULL
    ),
    CONSTRAINT valid_sf_role CHECK (
        role IN ('owner', 'manager', 'staff', 'cashier')
    )
);

-- =====================================================
-- 3. ADD INVITATION_ID TO STOREFRONT_STAFF
-- =====================================================

ALTER TABLE storefront_staff
ADD COLUMN IF NOT EXISTS invitation_id BIGINT
    REFERENCES storefront_invitations(id) ON DELETE SET NULL;

-- =====================================================
-- 4. INDEXES
-- =====================================================

CREATE INDEX idx_storefront_invitations_storefront
    ON storefront_invitations(storefront_id);

CREATE INDEX idx_storefront_invitations_email
    ON storefront_invitations(invited_email)
    WHERE invited_email IS NOT NULL;

CREATE INDEX idx_storefront_invitations_user
    ON storefront_invitations(invited_user_id)
    WHERE invited_user_id IS NOT NULL;

CREATE INDEX idx_storefront_invitations_code
    ON storefront_invitations(invite_code)
    WHERE invite_code IS NOT NULL;

CREATE INDEX idx_storefront_invitations_status
    ON storefront_invitations(status);

CREATE INDEX idx_storefront_invitations_invited_by
    ON storefront_invitations(invited_by_id);

CREATE INDEX idx_storefront_invitations_expires
    ON storefront_invitations(expires_at)
    WHERE expires_at IS NOT NULL AND status = 'pending';

CREATE INDEX idx_storefront_staff_invitation
    ON storefront_staff(invitation_id)
    WHERE invitation_id IS NOT NULL;

-- =====================================================
-- 5. UPDATE TRIGGER
-- =====================================================

CREATE OR REPLACE FUNCTION update_storefront_invitations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_storefront_invitations_updated_at
    BEFORE UPDATE ON storefront_invitations
    FOR EACH ROW
    EXECUTE FUNCTION update_storefront_invitations_updated_at();

-- =====================================================
-- 6. COMMENTS
-- =====================================================

COMMENT ON TABLE storefront_invitations IS 'Staff invitation system for storefronts';
COMMENT ON COLUMN storefront_invitations.type IS 'Invitation type: email (one-time) or link (shareable)';
COMMENT ON COLUMN storefront_invitations.role IS 'Role to assign after accepting: owner, manager, staff, cashier';
COMMENT ON COLUMN storefront_invitations.invite_code IS 'Unique code for shareable link invitations (sf-XXXXXXXX)';
COMMENT ON COLUMN storefront_invitations.max_uses IS 'Max number of uses for link invitations (NULL = unlimited)';
COMMENT ON COLUMN storefront_invitations.invited_by_id IS 'User ID who created the invitation';
COMMENT ON COLUMN storefront_invitations.status IS 'Current status: pending, accepted, declined, expired, revoked';
