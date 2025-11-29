CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL, 
    email TEXT NOT NULL UNIQUE, 
    password TEXT NOT NULL
); 

CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(name, version)
);

CREATE TABLE IF NOT EXISTS orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id),
    product_id INT NOT NULL REFERENCES products(product_id),
    quantity INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, product_id)
);
