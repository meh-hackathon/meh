INSERT INTO users (username, email)
VALUES (
  'admin',
  'admin@example.com'
);

INSERT INTO user_credentials (user_id, password_hash)
SELECT id, '$2a$12$QCF/IQohEtSK5SbRDeiCNOTduZWUDVUR8HvGVbjo5IH1Jhyw9LLXW' -- admin123
FROM users
WHERE username = 'admin';
