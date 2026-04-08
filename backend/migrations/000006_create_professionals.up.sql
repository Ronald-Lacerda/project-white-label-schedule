CREATE TABLE professionals (
    id                 CHAR(26)     NOT NULL,
    establishment_id   CHAR(26)     NOT NULL,
    name               VARCHAR(120) NOT NULL,
    avatar_url         VARCHAR(500),
    phone              VARCHAR(20),
    google_calendar_id VARCHAR(250),
    display_order      SMALLINT     NOT NULL DEFAULT 0,
    active             BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at         DATETIME(6)  NOT NULL,
    updated_at         DATETIME(6)  NOT NULL,

    PRIMARY KEY (id),
    INDEX idx_professionals_establishment (establishment_id),
    CONSTRAINT fk_professionals_establishment FOREIGN KEY (establishment_id)
        REFERENCES establishments (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
