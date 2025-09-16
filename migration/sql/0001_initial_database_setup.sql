-- Create ENUM types for order types and statuses
CREATE TYPE order_type_enum AS ENUM ('delivery', 'return', 'exchange', 'pickup');
CREATE TYPE order_status_enum AS ENUM ('pending', 'confirmed', 'picked_up', 'in_transit', 'out_for_delivery', 'delivered', 'failed_delivery', 'returned', 'cancelled');

-- Users table
CREATE TABLE IF NOT EXISTS users (
                       id BIGSERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Stores table (referenced in orders)
CREATE TABLE IF NOT EXISTS stores (
                        id BIGSERIAL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        contact_phone VARCHAR(20),
                        address TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Cities table for delivery locations
CREATE TABLE IF NOT EXISTS cities (
                        id BIGSERIAL PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
                        base_delivery_fee DECIMAL(10,2) DEFAULT 100.00,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Zones table (districts within cities)
CREATE TABLE IF NOT EXISTS zones (
                       id BIGSERIAL PRIMARY KEY,
                       city_id BIGINT REFERENCES cities(id),
                       name VARCHAR(100) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Item types lookup (keeping as table since these might vary more)
CREATE TABLE IF NOT EXISTS item_types (
                            id BIGSERIAL PRIMARY KEY,
                            name VARCHAR(50) NOT NULL,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Delivery types lookup (keeping as table for flexibility)
CREATE TABLE IF NOT EXISTS delivery_types (
                                id BIGSERIAL PRIMARY KEY,
                                name VARCHAR(50) NOT NULL,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Main orders table
CREATE TABLE IF NOT EXISTS orders (
                        id BIGSERIAL PRIMARY KEY,
                        consignment_id VARCHAR(50) UNIQUE NOT NULL,
                        user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
                        store_id BIGINT REFERENCES stores(id),
                        merchant_order_id VARCHAR(100),
                        recipient_name VARCHAR(255) NOT NULL,
                        recipient_phone VARCHAR(20) NOT NULL,
                        recipient_address TEXT NOT NULL,
                        recipient_city_id BIGINT REFERENCES cities(id),
                        recipient_zone_id BIGINT REFERENCES zones(id),
                        recipient_area_address TEXT,
                        order_type order_type_enum NOT NULL DEFAULT 'delivery',
                        delivery_type_id BIGINT REFERENCES delivery_types(id),
                        item_type_id BIGINT REFERENCES item_types(id),
                        item_quantity INTEGER NOT NULL DEFAULT 1,
                        item_weight DECIMAL(8,2) NOT NULL,
                        item_description TEXT,
                        special_instruction TEXT,
                        amount_to_collect DECIMAL(10,2) NOT NULL,
                        delivery_fee DECIMAL(10,2) NOT NULL,
                        cod_fee DECIMAL(10,2) NOT NULL DEFAULT 0,
                        promo_discount DECIMAL(10,2) DEFAULT 0,
                        discount DECIMAL(10,2) DEFAULT 0,
                        total_fee DECIMAL(10,2) NOT NULL,
                        order_status order_status_enum NOT NULL DEFAULT 'pending',
                        archived BOOLEAN DEFAULT FALSE,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- User sessions for JWT token management
CREATE TABLE IF NOT EXISTS user_sessions (
                               id BIGSERIAL PRIMARY KEY,
                               user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
                               token_hash VARCHAR(255) NOT NULL,
                               expires_at TIMESTAMP NOT NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_consignment_id ON orders(consignment_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(order_status);
CREATE INDEX IF NOT EXISTS idx_orders_order_type ON orders(order_type);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_archived ON orders(archived);
CREATE INDEX IF NOT EXISTS idx_orders_recipient_city_zone ON orders(recipient_city_id, recipient_zone_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);