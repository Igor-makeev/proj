package postgress

const usersTableSchema = `
CREATE TABLE if not exists users_table
(
    id serial not null unique,
    login varchar(255) not null unique,
    password_hash varchar(255) not null
);`

const ordersTableSchema = `
CREATE TABLE if not exists orders_table
(
    user_id int references users_table(id) on delete cascade not null,
    number int not null unique,
    status varchar(10) not null,
	  accrual float,
    uploaded_at timestamp
);`

var LoginIndex = `
CREATE UNIQUE INDEX if not exists login_index_unique
  ON users_table
  USING btree(login);
`
var OrderIndex = `
CREATE UNIQUE INDEX if not exists order_index_unique
  ON order_table
  USING btree(number);
`
