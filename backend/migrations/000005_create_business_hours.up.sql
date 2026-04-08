CREATE TABLE business_hours (
    id               CHAR(26)   NOT NULL,
    establishment_id CHAR(26)   NOT NULL,
    day_of_week      TINYINT    NOT NULL COMMENT '0=Dom, 1=Seg, 2=Ter, 3=Qua, 4=Qui, 5=Sex, 6=Sab',
    open_time        TIME       NOT NULL,
    close_time       TIME       NOT NULL,
    is_closed        BOOLEAN    NOT NULL DEFAULT FALSE,

    PRIMARY KEY (id),
    UNIQUE INDEX idx_business_hours_day (establishment_id, day_of_week),
    CONSTRAINT fk_bh_establishment FOREIGN KEY (establishment_id)
        REFERENCES establishments (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
