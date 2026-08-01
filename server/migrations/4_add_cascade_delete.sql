ALTER TABLE resources DROP CONSTRAINT IF EXISTS resources_drop_id_fkey;
ALTER TABLE resources ADD CONSTRAINT resources_drop_id_fkey FOREIGN KEY (drop_id) REFERENCES drops(id) ON DELETE CASCADE;
