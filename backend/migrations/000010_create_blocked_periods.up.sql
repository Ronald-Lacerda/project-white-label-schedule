CREATE TABLE blocked_periods (
    id              CHAR(26)    NOT NULL,
    professional_id CHAR(26)    NOT NULL,
    starts_at       DATETIME(6) NOT NULL,
    ends_at         DATETIME(6) NOT NULL,
    reason          VARCHAR(250),
    google_event_id VARCHAR(250),

    PRIMARY KEY (id),
    INDEX idx_blocked_professional_range (professional_id, starts_at, ends_at),
    CONSTRAINT fk_blocked_professional FOREIGN KEY (professional_id)
        REFERENCES professionals (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
