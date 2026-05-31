INSERT INTO drops (id, expiration_date) VALUES ('test1', '2050-01-01');
INSERT INTO drops (id, expiration_date) VALUES ('test2', '2050-01-01');

INSERT INTO resources (id, type, drop_id, name) VALUES ('11111111-1111-1111-1111-111111111111', 'file', 'test1', 'file1.txt');
INSERT INTO resources (id, type, drop_id, name) VALUES ('22222222-2222-2222-2222-222222222222', 'file', 'test1', 'file2.txt');
INSERT INTO resources (id, type, drop_id, name) VALUES ('33333333-3333-3333-3333-333333333333', 'text', 'test2', NULL);
