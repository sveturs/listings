-- Remove added columns from post_express_shipments table
ALTER TABLE post_express_shipments DROP CONSTRAINT IF EXISTS fk_sender_location;
ALTER TABLE post_express_shipments DROP CONSTRAINT IF EXISTS fk_recipient_location;

ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS sender_location_id;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS recipient_location_id;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS cod_reference;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS base_price;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS insurance_fee;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS cod_fee;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS total_price;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS delivery_status;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS delivery_instructions;
ALTER TABLE post_express_shipments DROP COLUMN IF EXISTS notes;