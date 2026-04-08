CREATE TABLE professional_hours (
    id              CHAR(26)  NOT NULL,
    professional_id CHAR(26)  NOT NULL,
    day_of_week     TINYINT   NOT NULL COMMENT '0=Dom, 1=Seg, 2=Ter, 3=Qua, 4=Qui, 5=Sex, 6=Sab',
    start_time      TIME      NOT NULL,
    end_time        TIME      NOT NULL,
    is_unavailable  BOOLEAN   NOT NULL DEFAULT FALSE,

    PRIMARY KEY (id),
    UNIQUE INDEX idx_prof_hours_day (professional_id, day_of_week),
    CONSTRAINT fk_prof_hours_professional FOREIGN KEY (professional_id)
        REFERENCES professionals (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
