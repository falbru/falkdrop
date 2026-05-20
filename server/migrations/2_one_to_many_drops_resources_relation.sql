ALTER TABLE resources ADD COLUMN drop_id varchar(5) REFERENCES drops (id);

UPDATE resources
SET drop_id = dr.drop_id
FROM drops_resources dr
WHERE resources.id = dr.resource_id;

DROP TABLE drops_resources;
