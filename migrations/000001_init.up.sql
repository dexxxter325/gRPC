create table users(
    id serial primary key ,
    email varchar(255) unique not null ,
    password varchar(255) not null ,

);
create table investments(
    id serial primary key ,
    amount integer not null,
    currency varchar not null ,
);