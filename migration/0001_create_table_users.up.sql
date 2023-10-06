begin;
create table if not exists users (
    id serial primary key,
    email varchar(255) unique not null,
    name varchar(255) not null,
    password varchar(255) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

commit;