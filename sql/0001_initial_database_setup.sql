-- Database Schema for Order Management System

-- Users table for authentication
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Stores table (referenced in orders)
CREATE TABLE stores (
                        id BIGINT PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        contact_phone VARCHAR(20),
                        address TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cities table for delivery locations
CREATE TABLE cities (
                        id BIGINT PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
                        base_delivery_fee DECIMAL(10,2) DEFAULT 100.00,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Zones table (districts within cities)
CREATE TABLE zones (
                       id BIGINT PRIMARY KEY,
                       city_id BIGINT REFERENCES cities(id),
                       name VARCHAR(100) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Areas table (specific areas within zones)
CREATE TABLE areas (
                       id BIGINT PRIMARY KEY,
                       zone_id BIGINT REFERENCES zones(id),
                       name VARCHAR(100) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order types lookup
CREATE TABLE order_types (
                             id BIGINT PRIMARY KEY,
                             name VARCHAR(50) NOT NULL, -- 'Delivery', 'Return', etc.
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Item types lookup
CREATE TABLE item_types (
                            id BIGINT PRIMARY KEY,
                            name VARCHAR(50) NOT NULL, -- 'Parcel', 'Document', 'Fragile', etc.
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Delivery types lookup
CREATE TABLE delivery_types (
                                id BIGINT PRIMARY KEY,
                                name VARCHAR(50) NOT NULL, -- 'Regular', 'Express', 'Same Day', etc.
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order statuses lookup
CREATE TABLE order_statuses (
                                id BIGINT PRIMARY KEY,
                                name VARCHAR(50) NOT NULL, -- 'Pending', 'Confirmed', 'Picked Up', 'In Transit', 'Delivered', 'Cancelled', etc.
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Main orders table
CREATE TABLE orders (
                        id BIGSERIAL PRIMARY KEY,
                        consignment_id VARCHAR(50) UNIQUE NOT NULL,
                        user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
                        store_id BIGINT REFERENCES stores(id),
                        merchant_order_id VARCHAR(100),

    -- Recipient information
                        recipient_name VARCHAR(255) NOT NULL,
                        recipient_phone VARCHAR(20) NOT NULL,
                        recipient_address TEXT NOT NULL,
                        recipient_city_id BIGINT REFERENCES cities(id),
                        recipient_zone_id BIGINT REFERENCES zones(id),
                        recipient_area_id BIGINT REFERENCES areas(id),

    -- Order details
                        order_type_id BIGINT REFERENCES order_types(id),
                        delivery_type_id BIGINT REFERENCES delivery_types(id),
                        item_type_id BIGINT REFERENCES item_types(id),
                        item_quantity INTEGER NOT NULL DEFAULT 1,
                        item_weight DECIMAL(8,2) NOT NULL,
                        item_description TEXT,
                        special_instruction TEXT,

    -- Financial information
                        amount_to_collect DECIMAL(10,2) NOT NULL,
                        delivery_fee DECIMAL(10,2) NOT NULL,
                        cod_fee DECIMAL(10,2) NOT NULL DEFAULT 0,
                        promo_discount DECIMAL(10,2) DEFAULT 0,
                        discount DECIMAL(10,2) DEFAULT 0,
                        total_fee DECIMAL(10,2) NOT NULL,

    -- Status and tracking
                        order_status_id BIGINT REFERENCES order_statuses(id),
                        transfer_status INTEGER DEFAULT 1, -- 1: active, 0: inactive
                        archived BOOLEAN DEFAULT FALSE,

    -- Timestamps
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User sessions for JWT token management (optional but recommended)
CREATE TABLE user_sessions (
                               id BIGSERIAL PRIMARY KEY,
                               user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
                               token_hash VARCHAR(255) NOT NULL,
                               expires_at TIMESTAMP NOT NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_consignment_id ON orders(consignment_id);
CREATE INDEX idx_orders_status ON orders(order_status_id);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_transfer_status ON orders(transfer_status);
CREATE INDEX idx_orders_archived ON orders(archived);
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Insert initial data
INSERT INTO cities (id, name, base_delivery_fee) VALUES
                                                     (1, 'Dhaka', 60.00),
                                                     (2, 'Chittagong', 100.00),
                                                     (3, 'Sylhet', 100.00);

INSERT INTO zones (id, city_id, name) VALUES
                                          (1, 1, 'Gulshan'),
                                          (2, 1, 'Dhanmondi'),
                                          (3, 1, 'Banani');

INSERT INTO areas (id, zone_id, name) VALUES
                                          (1, 1, 'Gulshan 1'),
                                          (2, 1, 'Gulshan 2'),
                                          (3, 3, 'Banani');

INSERT INTO order_types (id, name) VALUES
                                       (1, 'Delivery'),
                                       (2, 'Return'),
                                       (3, 'Exchange');

INSERT INTO item_types (id, name) VALUES
                                      (1, 'Document'),
                                      (2, 'Parcel'),
                                      (3, 'Fragile'),
                                      (4, 'Liquid');

INSERT INTO delivery_types (id, name) VALUES
                                          (24, 'Regular'),
                                          (48, 'Express'),
                                          (72, 'Same Day');

INSERT INTO order_statuses (id, name) VALUES
                                          (1, 'Pending'),
                                          (2, 'Confirmed'),
                                          (3, 'Picked Up'),
                                          (4, 'In Transit'),
                                          (5, 'Delivered'),
                                          (6, 'Cancelled'),
                                          (7, 'Returned');

INSERT INTO stores (id, name, contact_phone) VALUES
    (131172, 'Default Store', '+8801234567890');

-- Sample user (password should be hashed in application)
INSERT INTO users (email, password_hash) VALUES
    ('01901901901@mailinator.com', '$2a$10$hash_will_be_generated_by_application');