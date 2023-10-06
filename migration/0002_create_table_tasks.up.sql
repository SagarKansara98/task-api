begin;
create table if not exists tasks (
    id serial primary key ,
    user_id integer not null REFERENCES users(id),
    title varchar(256) not null,
    description text not null,
    status varchar(256) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
commit;