CREATE SCHEMA IF NOT EXISTS sellers_schema;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE sellers_schema.organization_type AS ENUM ('LLC', 'JSC', 'OOO','IP');

CREATE TABLE IF NOT EXISTS sellers_schema.sellers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_name VARCHAR(255) NOT NULL,
    store_description TEXT,
    logo_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    tin VARCHAR(32) UNIQUE NOT NULL,
    kpp VARCHAR(32) UNIQUE NOT NULL,
    organization_form sellers_schema.organization_type NOT NULL,
    psrn VARCHAR(32) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(100) NOT NULL,
    building VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS idx_sellers_tin ON sellers_schema.sellers(tin);

ALTER DATABASE "Hailow" SET TIMEZONE TO 'Europe/Moscow';