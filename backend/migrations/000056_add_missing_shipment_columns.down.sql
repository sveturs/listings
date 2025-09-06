-- Remove added columns
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS invoice_url;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS invoice_number;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS invoice_date;