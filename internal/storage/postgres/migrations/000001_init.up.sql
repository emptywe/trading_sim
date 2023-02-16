CREATE TABLE IF NOT EXISTS users
(
    id serial not null unique,
    firstname varchar(255) ,
    lastname varchar(255) ,
    username varchar(255) not null unique,
    email varchar(255) not null,
    password_hash varchar(255) not null,
    status varchar(255) ,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    last_signed_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS currencies
(
    id serial not null unique,
    name varchar(255) not null unique,
    value float4 not null
);

CREATE TABLE IF NOT EXISTS basket
(
    id serial not null unique,
    user_id int references users(id)  not null,
    currency_id int references currencies(id)  not null,
    Currency varchar(255) references currencies(name)  not null,
    amount float4
);

INSERT INTO currencies (name, value)  values ('usdt', 1);

