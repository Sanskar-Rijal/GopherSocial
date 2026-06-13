ALTER TABLE IF EXISTS users 
ADD COLUMN role_id INT REFERENCES roles(id) DEFAULT 1; 

UPDATE users SET role_id = (
    SELECT id FROM roles WHERE name='user'
);

-- DEFAULT 1 was needed for EXISTING users
-- but NEW users being created from now on
-- should NOT get automatic default
ALTER TABLE users ALTER COLUMN role_id DROP DEFAULT;

-- NOT NULL means — this column CANNOT be empty
-- every user MUST have a role_id
-- no exceptions
ALTER TABLE
  users
ALTER COLUMN
  role_id
SET
  NOT NULL;