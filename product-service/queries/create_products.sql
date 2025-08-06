CREATE SCHEMA categories;
CREATE TABLE categories.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    image_url VARCHAR(512),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE SCHEMA products;
CREATE TABLE products.products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12, 2) NOT NULL,
    category_id UUID NOT NULL REFERENCES categories.categories(id),
    images JSONB,
    quantity DECIMAL(10, 3) NOT NULL,
    unit VARCHAR(20) NOT NULL CHECK (unit IN ('kg', 'g', 'pcs', 'pack', 'l', 'ml', 'm', 'cm')),
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    weight DECIMAL(10, 3)
);


CREATE INDEX idx_products_category ON products.products(category_id);
CREATE INDEX idx_products_name ON products.products(name);
CREATE INDEX idx_products_price ON products.products(price);