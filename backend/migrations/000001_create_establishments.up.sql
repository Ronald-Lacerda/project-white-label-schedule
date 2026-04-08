CREATE TABLE establishments (
    id                       CHAR(26)     NOT NULL,
    name                     VARCHAR(120) NOT NULL,
    slug                     VARCHAR(60)  NOT NULL,
    timezone                 VARCHAR(50)  NOT NULL DEFAULT 'America/Sao_Paulo',
    contact_email            VARCHAR(120),
    contact_phone            VARCHAR(20),
    min_advance_cancel_hours SMALLINT     NOT NULL DEFAULT 0,
    active                   BOOLEAN      NOT NULL DEFAULT TRUE,
    google_calendar_connected BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at               DATETIME(6)  NOT NULL,
    updated_at               DATETIME(6)  NOT NULL,

    PRIMARY KEY (id),
    UNIQUE INDEX idx_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
