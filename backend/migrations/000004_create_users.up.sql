CREATE TABLE users (
    id               CHAR(26)     NOT NULL,
    establishment_id CHAR(26)     NOT NULL,
    name             VARCHAR(120) NOT NULL,
    email            VARCHAR(120) NOT NULL,
    password_hash    VARCHAR(255) NOT NULL,
    role             ENUM('owner','manager') NOT NULL DEFAULT 'owner',
    active           BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at       DATETIME(6)  NOT NULL,

    PRIMARY KEY (id),
    UNIQUE INDEX idx_users_email (email),
    INDEX idx_users_establishment (establishment_id),
    CONSTRAINT fk_users_establishment FOREIGN KEY (establishment_id)
        REFERENCES establishments (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
