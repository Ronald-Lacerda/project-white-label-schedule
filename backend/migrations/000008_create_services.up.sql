CREATE TABLE services (
    id               CHAR(26)     NOT NULL,
    establishment_id CHAR(26)     NOT NULL,
    name             VARCHAR(120) NOT NULL,
    description      VARCHAR(500),
    duration_minutes SMALLINT     NOT NULL,
    price_cents      INT          COMMENT 'NULL = preço não exibido',
    active           BOOLEAN      NOT NULL DEFAULT TRUE,
    display_order    SMALLINT     NOT NULL DEFAULT 0,
    created_at       DATETIME(6)  NOT NULL,

    PRIMARY KEY (id),
    INDEX idx_services_establishment (establishment_id),
    CONSTRAINT fk_services_establishment FOREIGN KEY (establishment_id)
        REFERENCES establishments (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
