CREATE TABLE appointments (
    id               CHAR(26)    NOT NULL,
    establishment_id CHAR(26)    NOT NULL,
    professional_id  CHAR(26)    NOT NULL,
    service_id       CHAR(26)    NOT NULL,
    client_name      VARCHAR(120) NOT NULL,
    client_phone     VARCHAR(20)  NOT NULL,
    starts_at        DATETIME(6)  NOT NULL,
    ends_at          DATETIME(6)  NOT NULL,
    status           ENUM('pending','confirmed','cancelled','completed','no_show')
                     NOT NULL DEFAULT 'confirmed',
    source           ENUM('booking_link','manager','api')
                     NOT NULL DEFAULT 'booking_link',
    google_event_id  VARCHAR(250),
    notes            VARCHAR(500),
    idempotency_key  CHAR(36),
    created_at       DATETIME(6)  NOT NULL,
    updated_at       DATETIME(6)  NOT NULL,

    PRIMARY KEY (id),
    UNIQUE INDEX idx_appointments_idempotency (idempotency_key),
    INDEX idx_appointments_establishment_time (establishment_id, starts_at, ends_at),
    INDEX idx_appointments_professional_time (professional_id, starts_at, ends_at),
    INDEX idx_appointments_status (status),
    CONSTRAINT fk_appt_establishment FOREIGN KEY (establishment_id)
        REFERENCES establishments (id),
    CONSTRAINT fk_appt_professional FOREIGN KEY (professional_id)
        REFERENCES professionals (id),
    CONSTRAINT fk_appt_service FOREIGN KEY (service_id)
        REFERENCES services (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
