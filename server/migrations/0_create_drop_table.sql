CREATE TABLE drops (
  id varchar(5) PRIMARY KEY,
  expiration_date date
);

CREATE TYPE resource_type AS ENUM ('file', 'text');

CREATE TABLE resources (
  link text PRIMARY KEY,
  type resource_type
);

CREATE TABLE drops_resources (
  drop_id varchar(5) REFERENCES drops(id),
  resource_link text REFERENCES resources(link)
);
