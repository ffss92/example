-- +goose Up
CREATE TABLE IF NOT EXISTS "users" (
    "id" INTEGER PRIMARY KEY,
    "email" TEXT NOT NULL UNIQUE, -- Adding a unique idx like this is a big mistake. This column cannot be altered.
    "password_hash" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "tokens" (
    "hash" BLOB PRIMARY KEY,
    "scope" TEXT NOT NULL,
    "expiry" TIMESTAMP NOT NULL,
    "user_id" INTEGER NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS "tokens";
DROP TABLE IF EXISTS "users";
