package postgress

const usersTableSchema = `
CREATE TABLE if not exists users_table
(
    id serial not null unique,
    login varchar(255) not null unique,
    password_hash varchar(255) not null
);`

var Index = `
CREATE UNIQUE INDEX if not exists login_index_unique
  ON users_table
  USING btree(login);
`
