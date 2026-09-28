CREATE TABLE order_images (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(id) ON DELETE CASCADE
);