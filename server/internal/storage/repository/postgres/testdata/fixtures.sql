INSERT INTO drops (id, expiration_date) VALUES ('test1', '2050-01-01');
INSERT INTO drops (id, expiration_date) VALUES ('test2', '2050-01-01');
INSERT INTO drops (id, expiration_date) VALUES ('expd1', '2000-01-01');
INSERT INTO drops (id, expiration_date) VALUES ('expd2', '2005-01-01');

INSERT INTO resources (id, type, drop_id, name) VALUES ('11111111-1111-1111-1111-111111111111', 'file', 'test1', 'file1.txt');
INSERT INTO resources (id, type, drop_id, name) VALUES ('22222222-2222-2222-2222-222222222222', 'file', 'test1', 'file2.txt');
INSERT INTO resources (id, type, drop_id, name) VALUES ('33333333-3333-3333-3333-333333333333', 'text', 'test2', NULL);
INSERT INTO resources (id, type, drop_id, name) VALUES ('44444444-4444-4444-4444-444444444444', 'text', 'expd1', NULL);
INSERT INTO resources (id, type, drop_id, name) VALUES ('55555555-5555-5555-5555-555555555555', 'text', 'expd1', NULL);
INSERT INTO resources (id, type, drop_id, name) VALUES ('66666666-6666-6666-6666-666666666666', 'file', 'expd2', 'file3.txt');
