

CREATE  DOMAIN u_datetime    timestamptz; -- date and time with timezone
create function ep(timestamptz) returns bigint as 'select cast(extract(epoch from $1)*1000 as bigint);' language sql immutable;
create function ts(bigint) returns timestamptz as 'select to_timestamp($1/1000.0);' language sql immutable;

CREATE TABLE task(
 id          serial PRIMARY KEY,
 title       char(50),
 description char(50),
 create_time u_datetime default now(),
 update_time u_datetime default now(),
 deadline    u_datetime
);

