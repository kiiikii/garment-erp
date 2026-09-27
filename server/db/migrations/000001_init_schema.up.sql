CREATE TABLE customers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  phone VARCHAR(50) NOT NULL,
  address TEXT NOT NULL
);

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  customer_id INT NOT NULL,
  total_quantity INT NOT NULL,
  production_type VARCHAR(50) NOT NULL,
  status VARCHAR(50) DEFAULT 'WAITING_FOR_SAMPLE',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_customer FOREIGN KEY(customer_id) REFERENCES customers(id)
);

CREATE TABLE order_sizes (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL,
  size_table VARCHAR(10) NOT NULL,
  quantity INT NOT NULL,
  CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);