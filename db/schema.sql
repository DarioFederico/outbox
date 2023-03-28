CREATE TABLE IF NOT EXISTS categories
(
    id          int          auto_increment primary key,
    name        varchar(50)  not null,
    description varchar(100) not null,
    created_at  datetime     not null,
    updated_at  datetime     null,
    constraint categories_unique unique (name)
);

CREATE TABLE IF NOT EXISTS outbox
(
    id         int         auto_increment primary key,
    type       varchar(50) null,
    message    text        not null,
    status     varchar(50) not null,
    created_at datetime    not null,
    updated_at datetime    null
);

