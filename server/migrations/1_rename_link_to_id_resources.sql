ALTER TABLE resources RENAME COLUMN link TO id;
ALTER TABLE drops_resources RENAME COLUMN resource_link TO resource_id;
