-- Rename receiver columns to recipient for consistency
ALTER TABLE post_express_shipments RENAME COLUMN receiver_name TO recipient_name;
ALTER TABLE post_express_shipments RENAME COLUMN receiver_address TO recipient_address;
ALTER TABLE post_express_shipments RENAME COLUMN receiver_city TO recipient_city;
ALTER TABLE post_express_shipments RENAME COLUMN receiver_postal_code TO recipient_postal_code;
ALTER TABLE post_express_shipments RENAME COLUMN receiver_phone TO recipient_phone;
ALTER TABLE post_express_shipments RENAME COLUMN receiver_email TO recipient_email;