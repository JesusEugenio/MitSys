-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    full_name VARCHAR(60) NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_cashier BOOLEAN NOT NULL DEFAULT FALSE,
    is_delivery BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(60) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    cut_name VARCHAR(50) NOT NULL,
    price_per_kilo NUMERIC(10,2) NOT NULL,
    available_stock NUMERIC(10,3) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE sale_transactions (
    id SERIAL PRIMARY KEY,
    total_amount NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL,
    sale_type VARCHAR(20) NOT NULL,
    user_id INT,
    customer_id INT,

    CONSTRAINT chk_sale_transactions_status
        CHECK (status IN('pending','paid','cancelled')),

    CONSTRAINT chk_sale_transactions_sale_type
        CHECK (sale_type IN('delivery','in_store')),

    CONSTRAINT fk_sale_transactions_users
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT,

    CONSTRAINT fk_sale_transactions_customers
        FOREIGN KEY (customer_id)
        REFERENCES customers(id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
);

CREATE TABLE delivery_orders (
    id SERIAL PRIMARY KEY,
    delivery_status BOOLEAN NOT NULL DEFAULT FALSE,
    sale_id INT,
    delivery_id INT,

    CONSTRAINT fk_delivery_orders_sale_transactions
        FOREIGN KEY (sale_id)
        REFERENCES sale_transactions(id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT,

    CONSTRAINT fk_delivery_orders_users
        FOREIGN KEY (delivery_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
);

CREATE TABLE sale_transactions_products (
    sale_transaction_id INT NOT NULL,
    product_id INT NOT NULL,
    weight_sold NUMERIC(10,3) NOT NULL,
    unit_price NUMERIC(10,2) NOT NULL,

    CONSTRAINT fk_sale_transactions
        FOREIGN KEY (sale_transaction_id)
        REFERENCES sale_transactions(id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT,

    CONSTRAINT fk_products
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT,

    PRIMARY KEY (sale_transaction_id, product_id)
);

-- +goose Down

DROP TABLE IF EXISTS sale_transactions_products;
DROP TABLE IF EXISTS delivery_orders;
DROP TABLE IF EXISTS sale_transactions;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS users;
