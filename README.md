# SAC_Server
Shop Around The Corner - Server

Pull the code
Create data directory : mkdir -p docker/pq-server/data
Run "docker-compose up"

Setting up postgres DB:
Connect to the postgres DB : `docker exec -it sacserver_pq-server_1 psql -U postgres`
Run the following commands in psql console :
`CREATE database sac;`
`\c sac`
`CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`
`CREATE EXTENSION IF NOT EXISTS postgis;`
 
`CREATE TABLE shop_categories (
  id SERIAL,
  name VARCHAR(255)
);`

`INSERT INTO shop_categories(name) VALUES('Electronics');
INSERT INTO shop_categories(name) VALUES('Restaurant');
INSERT INTO shop_categories(name) VALUES('Grocery');
INSERT INTO shop_categories(name) VALUES('Stationary');
INSERT INTO shop_categories(name) VALUES('Electricals');
INSERT INTO shop_categories(name) VALUES('Fancy Store');
INSERT INTO shop_categories(name) VALUES('Medicals');
INSERT INTO shop_categories(name) VALUES('Fruit & Veg');
INSERT INTO shop_categories(name) VALUES('Others');`

`CREATE TABLE tags (
  id SERIAL,
  name VARCHAR(255),
  category_id BIGINT
);`

`INSERT INTO tags(name, category_id) VALUES('Tea', 2);
INSERT INTO tags(name, category_id) VALUES('Lemon Tea', 2);
INSERT INTO tags(name, category_id) VALUES('Coffee', 2);
INSERT INTO tags(name, category_id) VALUES('Dosa', 2);
INSERT INTO tags(name, category_id) VALUES('Mutton Biriyani', 2);
INSERT INTO tags(name, category_id) VALUES('Chicken Biriyani', 2);
INSERT INTO tags(name, category_id) VALUES('Prawn Biriyani', 2);`

`select * from tags where lower(name) like '%tea%';`

`CREATE TABLE shops (
  id SERIAL,
  name VARCHAR(255),
  description VARCHAR(1024),
  phone VARCHAR(255),
  owner VARCHAR(255),
  address VARCHAR(1024),
  category_id BIGINT,
  latitude NUMERIC,
  longitude NUMERIC,
  location_geom geometry(POINT,2163)
);`

`INSERT INTO shops(name, description, phone, owner, address, category_id, latitude, longitude, location_geom) 
  VALUES('Junior Kuppanna', 'Traditional South Indian Restaurant', '0987654321', 'Kuppanna', '26, Srinivasa Nagar, Kandanchavadi, Perungudi, Chennai', 2, 12.9693, 80.2486, ST_Transform(ST_SetSRID(ST_MakePoint(12.9693, 80.2486),4326),2163));`

`CREATE TABLE shop_tags (
  shop_id BIGINT, 
  tag_id BIGINT
);`

`INSERT INTO shop_tags VALUES(1, 3);
INSERT INTO shop_tags VALUES(1, 5);
INSERT INTO shop_tags VALUES(1, 6);
INSERT INTO shop_tags VALUES(1, 7);`

`CREATE TABLE search_requests (
  request_time TIMESTAMPTZ NOT NULL,
  tag_id BIGINT,
  category_id BIGINT,
  latitude NUMERIC,
  longitude NUMERIC,
  request_geom geometry(POINT,2163)
);`

`SELECT create_hypertable('search_requests', 'request_time', 'tag_id', 128);`

`INSERT INTO search_requests VALUES('2017-08-19 13:00:00', 3, 2, 12.9697583, 80.2436885, ST_Transform(ST_SetSRID(ST_MakePoint(12.9697583, 80.2436885),4326),2163));`

`select time_bucket('1 minute', request_time) AS minute_bucket, tag_id, count(*) from search_requests where tag_id = 3 group by minute_bucket, tag_id order by count(*) desc limit 25;`
