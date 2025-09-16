-- 1. Insert Users first (no dependencies)
-- INSERT INTO users (email, password_hash) VALUES
--     ('01901901901@mailinator.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1tJXKAv1j6vK8OfKNl1YHbdO8YhLwPq'); -- password: 321dsa

-- 2. Insert Cities (no dependencies)
INSERT INTO cities (name, base_delivery_fee) VALUES
                                                 ('Dhaka', 60.00),
                                                 ('Chittagong', 70.00),
                                                 ('Sylhet', 80.00),
                                                 ('Rajshahi', 75.00),
                                                 ('Khulna', 75.00),
                                                 ('Barishal', 80.00),
                                                 ('Rangpur', 85.00),
                                                 ('Mymensingh', 70.00);

-- 3. Insert Zones (depends on cities)
INSERT INTO zones (city_id, name) VALUES
-- Dhaka zones
(1, 'Dhanmondi'),
(1, 'Gulshan'),
(1, 'Uttara'),
(1, 'Mirpur'),
(1, 'Wari'),
(1, 'Motijheel'),
(1, 'Tejgaon'),
(1, 'Banani'),
-- Chittagong zones
(2, 'Agrabad'),
(2, 'Nasirabad'),
(2, 'Panchlaish'),
(2, 'Khulshi'),
-- Sylhet zones
(3, 'Zindabazar'),
(3, 'Amberkhana'),
(3, 'Subhanighat');

-- 4. Insert Item Types (no dependencies)
INSERT INTO item_types (name) VALUES
                                  ('Electronics'),
                                  ('Clothing'),
                                  ('Food & Beverage'),
                                  ('Documents'),
                                  ('Cosmetics'),
                                  ('Books'),
                                  ('Home Appliances'),
                                  ('Jewelry'),
                                  ('Toys'),
                                  ('Medicine');

-- 5. Insert Delivery Types (no dependencies)
INSERT INTO delivery_types (name) VALUES
                                      ('Standard Delivery'),
                                      ('Express Delivery'),
                                      ('Same Day Delivery'),
                                      ('Next Day Delivery'),
                                      ('Urgent Delivery');

-- 6. Insert Stores (no dependencies)
INSERT INTO stores (name, contact_phone, address) VALUES
                                                      ('Tech World Dhaka', '01712345678', '123 Dhanmondi R/A, Dhaka-1205'),
                                                      ('Fashion Plaza', '01812345679', '456 Gulshan Avenue, Dhaka-1212'),
                                                      ('Grocery Mart', '01912345680', '789 Uttara Sector 7, Dhaka-1230'),
                                                      ('Book Corner', '01612345681', '321 New Market, Dhaka-1205'),
                                                      ('Electronics Hub', '01512345682', '654 Mirpur 10, Dhaka-1216'),
                                                      ('Clothing Store', '01412345683', '987 Wari, Dhaka-1203'),
                                                      ('Food Court', '01312345684', '147 Motijheel C/A, Dhaka-1000'),
                                                      ('Medical Store', '01212345685', '258 Tejgaon I/A, Dhaka-1208');

-- 7. Sample Orders (depends on users, stores, cities, zones, item_types, delivery_types)
INSERT INTO orders (
    consignment_id, user_id, store_id, merchant_order_id,
    recipient_name, recipient_phone, recipient_address,
    recipient_city, recipient_zone, recipient_area,
    order_type, delivery_type_id, item_type,
    item_quantity, item_weight, item_description, special_instruction,
    order_amount, amount_to_collect, delivery_fee, cod_fee,
    promo_discount, discount, total_fee, order_status
) VALUES
-- Order 1
('CON202501001', 1, 1, 'MERCHANT001',
 'John Doe', '01711111111', 'House 10, Road 15, Block C',
 1, 1, 'Near Dhanmondi Lake',
 'delivery', 1, 1,
 1, 0.50, 'iPhone 15 Pro Max', 'Handle with care',
 120000.00, 120060.00, 60.00, 0.00,
 0.00, 0.00, 60.00, 'pending'),

-- Order 2
('CON202501002', 1, 2, 'MERCHANT002',
 'Jane Smith', '01722222222', 'Apt 5B, Road 11, Gulshan 2',
 1, 2, 'Opposite Gulshan Park',
 'delivery', 2, 2,
 3, 1.20, 'Cotton T-Shirts', 'Size: Medium, Color: Blue',
 1500.00, 1580.00, 80.00, 0.00,
 0.00, 0.00, 80.00, 'confirmed'),

-- Order 3
('CON202501003', 1, 3, 'MERCHANT003',
 'Mike Johnson', '01733333333', 'House 25, Sector 11, Uttara',
 1, 3, 'Near Uttara University',
 'delivery', 3, 3,
 2, 2.00, 'Fresh Vegetables and Fruits', 'Deliver before 6 PM',
 800.00, 900.00, 100.00, 0.00,
 0.00, 0.00, 100.00, 'picked_up'),

-- Order 4
('CON202501004', 1, 4, 'MERCHANT004',
 'Sarah Wilson', '01744444444', 'Flat 3A, Building 7, Mirpur 12',
 1, 4, 'DOHS Area',
 'delivery', 1, 6,
 5, 1.50, 'Programming Books Collection', 'Educational books',
 2500.00, 2560.00, 60.00, 0.00,
 0.00, 0.00, 60.00, 'in_transit'),

-- Order 5
('CON202501005', 1, 5, 'MERCHANT005',
 'David Brown', '01755555555', 'House 12, Lane 8, Wari',
 1, 5, 'Old Dhaka Area',
 'delivery', 4, 7,
 1, 5.00, 'Microwave Oven', 'Fragile item',
 15000.00, 15070.00, 70.00, 0.00,
 0.00, 0.00, 70.00, 'delivered'),

-- Order 6 (Return order example)
('CON202501006', 1, 1, 'RETURN001',
 'John Doe', '01711111111', 'House 10, Road 15, Block C',
 1, 1, 'Near Dhanmondi Lake',
 'return', 2, 1,
 1, 0.50, 'Defective Phone Case', 'Product return - defective',
 0.00, 0.00, 80.00, 0.00,
 0.00, 0.00, 80.00, 'pending');