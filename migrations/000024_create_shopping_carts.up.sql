-- Migration: Create Shopping Carts Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14
--
-- Purpose: Create shopping_carts table to support both authenticated and anonymous users
--          Each cart is tied to ONE storefront (multi-storefront = multiple carts)
--
-- Features:
-- - Authenticated users: identified by user_id
-- - Anonymous users: identified by session_id (cookie/localStorage)
-- - One cart per user per storefront
-- - Auto-updated timestamps

-- =====================================================
-- SHOPPING_CARTS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS shopping_carts (
    id BIGSERIAL PRIMARY KEY,

    -- User identification (mutually exclusive: user_id OR session_id)
    user_id BIGINT NULL,                    -- Authenticated user (FK to auth service)
    session_id VARCHAR(255) NULL,           -- Anonymous user session

    -- Storefront association
    storefront_id BIGINT NOT NULL,          -- FK to storefronts

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_shopping_carts_storefront FOREIGN KEY (storefront_id)
        REFERENCES storefronts(id) ON DELETE CASCADE,

    -- Business Logic Constraints
    -- Ensure EXACTLY ONE of user_id or session_id is set (not both, not neither)
    CONSTRAINT chk_shopping_carts_user_or_session CHECK (
        (user_id IS NOT NULL AND session_id IS NULL) OR
        (user_id IS NULL AND session_id IS NOT NULL)
    )
);

-- =====================================================
-- INDEXES
-- =====================================================

-- Primary lookup: find cart by authenticated user + storefront
CREATE INDEX idx_shopping_carts_user_storefront
    ON shopping_carts(user_id, storefront_id)
    WHERE user_id IS NOT NULL;

-- Primary lookup: find cart by anonymous session + storefront
CREATE INDEX idx_shopping_carts_session_storefront
    ON shopping_carts(session_id, storefront_id)
    WHERE session_id IS NOT NULL;

-- General lookup by storefront (for admin/analytics)
CREATE INDEX idx_shopping_carts_storefront_id
    ON shopping_carts(storefront_id);

-- Cleanup: find old anonymous carts by updated_at
CREATE INDEX idx_shopping_carts_updated_at
    ON shopping_carts(updated_at);

-- =====================================================
-- UNIQUE CONSTRAINTS
-- =====================================================

-- One cart per authenticated user per storefront
CREATE UNIQUE INDEX idx_shopping_carts_unique_user_per_storefront
    ON shopping_carts(user_id, storefront_id)
    WHERE user_id IS NOT NULL;

-- One cart per anonymous session per storefront
CREATE UNIQUE INDEX idx_shopping_carts_unique_session_per_storefront
    ON shopping_carts(session_id, storefront_id)
    WHERE session_id IS NOT NULL;

-- =====================================================
-- AUTO-UPDATE TRIGGER
-- =====================================================

CREATE OR REPLACE FUNCTION update_shopping_carts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_shopping_carts_updated_at
BEFORE UPDATE ON shopping_carts
FOR EACH ROW
EXECUTE FUNCTION update_shopping_carts_updated_at();

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE shopping_carts IS
    'Shopping carts for authenticated and anonymous users. One cart per user/session per storefront.';

COMMENT ON COLUMN shopping_carts.user_id IS
    'Authenticated user ID (FK to auth service). Mutually exclusive with session_id.';

COMMENT ON COLUMN shopping_carts.session_id IS
    'Anonymous user session ID (UUID from cookie/localStorage). Mutually exclusive with user_id.';

COMMENT ON COLUMN shopping_carts.storefront_id IS
    'Storefront this cart belongs to. Multi-storefront shopping requires multiple carts.';

COMMENT ON COLUMN shopping_carts.updated_at IS
    'Auto-updated on any cart modification. Used for cleanup of old anonymous carts (e.g., >30 days).';
