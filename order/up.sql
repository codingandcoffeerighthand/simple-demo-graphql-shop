CREATE TABLE IF NOT EXISTS "orders" (
    "id" CHAR(27) PRIMARY KEY, 
    "account_id" CHAR(27) NOT NULL,
    "created_at" TIMESTAMP WiTH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "total_price" MONEY NOT NULL
);

CREATE TABLE IF NOT EXISTS "order_products" ( 
    "product_id" CHAR(27) NOT NULL,
    "order_id" CHAR(27) REFERENCES "orders"(id) ON DELETE CASCADE,
    "quantity" INT NOT NULL,
    PRIMARY KEY ("product_id", "order_id")
);
