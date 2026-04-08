CREATE TABLE whitelabel_configs (
    establishment_id CHAR(26)     NOT NULL,
    logo_url         VARCHAR(500),
    primary_color    VARCHAR(7)   NOT NULL DEFAULT '#000000',
    secondary_color  VARCHAR(7),
    custom_domain    VARCHAR(120),
    custom_css       TEXT,

    PRIMARY KEY (establishment_id),
    CONSTRAINT fk_wl_establishment FOREIGN KEY (establishment_id)
        REFERENCES establishments (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
