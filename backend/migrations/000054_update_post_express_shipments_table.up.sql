-- Add missing columns to post_express_shipments table
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS sender_location_id INTEGER;
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS recipient_location_id INTEGER;
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS cod_reference VARCHAR(255);
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS base_price DECIMAL(10,2);
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS insurance_fee DECIMAL(10,2);
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS cod_fee DECIMAL(10,2);
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS total_price DECIMAL(10,2);
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS delivery_status VARCHAR(100);
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS delivery_instructions TEXT;
ALTER TABLE post_express_shipments ADD COLUMN IF NOT EXISTS notes TEXT;

-- Add foreign key constraints
ALTER TABLE post_express_shipments 
ADD CONSTRAINT fk_sender_location 
FOREIGN KEY (sender_location_id) 
REFERENCES post_express_locations(id)
ON DELETE SET NULL;

ALTER TABLE post_express_shipments 
ADD CONSTRAINT fk_recipient_location 
FOREIGN KEY (recipient_location_id) 
REFERENCES post_express_locations(id)
ON DELETE SET NULL;