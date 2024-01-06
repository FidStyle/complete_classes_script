-- auto-generated definition
create table orders
(
    id                   bigserial
        primary key,
    pw                   text,
    account              text,
    a                    numeric,
    b                    numeric,
    c                    numeric,
    d                    numeric,
    e                    numeric,
    f                    numeric,
    a0                   numeric,
    public_random        bigint,
    professional_random  bigint,
    specify_public       text,
    specify_professional text,
    b_n                  numeric,
    f_n                  numeric,
    a0_n                 numeric,
    created_at           timestamp with time zone,
    success_at           timestamp with time zone,
    condition            text
);

alter table orders
    owner to test;

-- auto-generated definition
create table users
(
    account    text,
    pw         text,
    created_at timestamp with time zone
);

alter table users
    owner to test;

-- auto-generated definition
create sequence classes_id_seq;

alter sequence classes_id_seq owner to test;

alter sequence classes_id_seq owned by classes.id;

-- auto-generated definition
create sequence orders_id_seq;

alter sequence orders_id_seq owner to test;

alter sequence orders_id_seq owned by orders.id;



