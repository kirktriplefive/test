CREATE TABLE users
(
    id serial primary key,
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE delivery
(
    d_id serial primary key,
    name varchar(255) not null,
    phone varchar(255) not null,
    zip varchar(255) not null,
    city varchar(255) not null,
    address varchar(255) not null,
    region varchar(255),
    email varchar(255)
);

CREATE TABLE payment
(
    id serial primary key,
    request_id varchar(255) unique,
    currency varchar(255) not null,
    provider varchar(255) not null,
    amount int not null,
    payment_dt int not null,
    bank varchar(255) not null,
    delivery_cost int not null,
    goods_total varchar(255) not null,
    custom_fee int not null default 0
);

CREATE TABLE item
(
    chrt_id int not null,
    track_number varchar(255) not null,
    price int not null,
    rid varchar(255) primary key,
    name varchar(255) not null,
    sale int not null default 0,
    size varchar(255) not null,
    total_price int not null,
    nm_id int not null,
    brand varchar(255),
    status int not null
);

CREATE TABLE orders
(
    order_uid serial primary key,
    track_number varchar(255) not null,
    entry varchar(255) not null,
    locale varchar(255) not null,
    internal_signature varchar(255),
    customer_id varchar(255) not null,
    delivery_service varchar(255) not null,
    shardkey varchar(255) not null,
    sm_id int not null,
    date_created timestamp not null,
    oof_shard varchar(255) not null,
    payment_id int references payment(id) on delete cascade not null,
    delivery_id int references delivery(d_id) on delete cascade not null
);

CREATE TABLE itemsinorder
(
    id serial primary key,
    order_id int references orders(order_uid) on delete cascade not null,
    item_id varchar(255) references item(rid) on delete cascade not null
);