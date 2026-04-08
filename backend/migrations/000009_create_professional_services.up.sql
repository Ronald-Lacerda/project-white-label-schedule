CREATE TABLE professional_services (
    professional_id CHAR(26) NOT NULL,
    service_id      CHAR(26) NOT NULL,

    PRIMARY KEY (professional_id, service_id),
    CONSTRAINT fk_ps_professional FOREIGN KEY (professional_id)
        REFERENCES professionals (id),
    CONSTRAINT fk_ps_service FOREIGN KEY (service_id)
        REFERENCES services (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
