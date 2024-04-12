CREATE TABLE orders
(
    order_uid TEXT UNIQUE PRIMARY KEY,
    track_number TEXT NOT NULL ,
    entry TEXT NOT NULL ,
    locale TEXT NOT NULL ,
    internal_signature TEXT NOT NULL  DEFAULT '',
    customer_id TEXT NOT NULL ,
    delivery_service TEXT NOT NULL,
    shardkey TEXT NOT NULL,
    sm_id INT NOT NULL ,
    date_created TEXT NOT NULL ,
    oof_shard TEXT NOT NULL
);

CREATE TABLE delivery
(
    id SERIAL PRIMARY KEY,
    name VARCHAR (70) NOT NULL ,
    phone VARCHAR(11) NOT NULL ,
    zip VARCHAR(10) NOT NULL ,
    city VARCHAR(100) NOT NULL ,
    address TEXT NOT NULL ,
    region varchar(100) NOT NULL ,
    email varchar(150) NOT NULL,
    order_uid TEXT NOT NULL ,
    CONSTRAINT fk_delivery_order FOREIGN KEY (order_uid) REFERENCES orders (order_uid)
);

CREATE TABLE  payment
(
    id SERIAL PRIMARY KEY,
    transaction TEXT NOT NULL UNIQUE ,
    request_id text NOT NULL  DEFAULT '',
    currency VARCHAR(30) NOT NULL,
    provider TEXT NOT NULL,
    amount  INT NOT NULL,
    payment_dt INT NOT NULL ,
    bank TEXT NOT NULL ,
    delivery_cost INT NOT NULL ,
    goods_total INT NOT NULL ,
    custom_fee INT NOT NULL,
    order_uid TEXT NOT NULL,
    CONSTRAINT fk_payment_order FOREIGN KEY (order_uid) REFERENCES orders (order_uid)
);

CREATE TABLE item
(
    id SERIAL PRIMARY KEY,
    chrt_id INT NOT NULL,
    track_number TEXT NOT NULL,
    price  INT NOT NULL ,
    rid TEXT NOT NULL ,
    name TEXT NOT NULL,
    sale INT NOT NULL ,
    size TEXT NOT NULL ,
    total_price INT NOT NULL ,
    nm_id INT NOT NULL ,
    brand TEXT NOT NULL ,
    status INT NOT NULL
);

CREATE TABLE orders_item
(
    item_id  INT REFERENCES  item (id),
    order_uid TEXT REFERENCES  orders(order_uid),
    CONSTRAINT orders_item_pkey PRIMARY KEY (item_id, order_uid)
);