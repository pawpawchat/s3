CREATE TABLE media (
    media_id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL,
    content_type VARCHAR(10),
    file_ext VARCHAR(8) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);