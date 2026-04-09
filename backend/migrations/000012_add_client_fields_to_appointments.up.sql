ALTER TABLE appointments
    ADD COLUMN client_email VARCHAR(120) NULL AFTER client_name,
    ADD COLUMN client_birth_date DATE NULL AFTER client_phone;
