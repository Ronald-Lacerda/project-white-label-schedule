CREATE TABLE google_oauth_tokens (
    establishment_id CHAR(26)    NOT NULL,
    access_token     TEXT        NOT NULL,
    refresh_token    TEXT        NOT NULL,
    expiry           DATETIME(6) NOT NULL,
    scope            TEXT        NOT NULL,
    updated_at       DATETIME(6) NOT NULL,

    PRIMARY KEY (establishment_id),
    CONSTRAINT fk_goauth_establishment FOREIGN KEY (establishment_id)
        REFERENCES establishments (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
