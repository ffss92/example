-- +goose Up
PRAGMA foreign_keys=off;
-- +goose StatementBegin
ALTER TABLE "users" RENAME TO "old_users";
CREATE TABLE "users" (
    "id" INTEGER PRIMARY KEY,
    "username" TEXT NOT NULL,
    "password_hash" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX "idx_users_username" ON "users"("username");
INSERT INTO "users" ("username", "password_hash") SELECT "email", "password_hash" FROM "old_users";
DROP TABLE "old_users";
-- +goose StatementEnd
PRAGMA foreign_keys=on;

-- +goose Down
PRAGMA foreign_keys=off;
-- +goose StatementBegin
ALTER TABLE "users" RENAME TO "old_users";
CREATE TABLE "users" (
    "id" INTEGER PRIMARY KEY,
    "email" TEXT NOT NULL,
    "password_hash" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX "idx_users_email" ON "users"("email");
INSERT INTO "users" ("email", "password_hash") SELECT "username", "password_hash" FROM "old_users";
DROP TABLE "old_users";
-- +goose StatementEnd
PRAGMA foreign_keys=on;
