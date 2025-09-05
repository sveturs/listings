-- Add missing columns to post_express_shipments
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS invoice_url TEXT;
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS invoice_number VARCHAR(255);
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS invoice_date TIMESTAMPTZ;