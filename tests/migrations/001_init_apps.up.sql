INSERT INTO apps (id, name, secret)
VALUES (1,'test','watermelon')
ON CONFLICT DO NOTHING;