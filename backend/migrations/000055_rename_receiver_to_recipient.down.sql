-- Rename back recipient columns to receiver
ALTER TABLE post_express_shipments RENAME COLUMN recipient_name TO receiver_name;
ALTER TABLE post_express_shipments RENAME COLUMN recipient_address TO receiver_address;
ALTER TABLE post_express_shipments RENAME COLUMN recipient_city TO receiver_city;
ALTER TABLE post_express_shipments RENAME COLUMN recipient_postal_code TO receiver_postal_code;
ALTER TABLE post_express_shipments RENAME COLUMN recipient_phone TO receiver_phone;
ALTER TABLE post_express_shipments RENAME COLUMN recipient_email TO receiver_email;