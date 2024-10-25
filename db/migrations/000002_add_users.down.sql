ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "unique_owner_currency";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
DROP TABLE IF EXISTS "users";