-- +goose Up
CREATE TABLE IF NOT EXISTS "users" (
    "id" INTEGER PRIMARY KEY,
    "username" TEXT NOT NULL,
    "password_hash" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX "idx_users_username" ON "users"("username");

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
