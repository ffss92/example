-- +goose Up
ALTER TABLE "users" ADD COLUMN "username" TEXT NOT NULL DEFAULT '';
CREATE UNIQUE INDEX "idx_users_username" ON "users"("username");
-- Set the default username to the email
UPDATE "users" SET "username" = "email" WHERE "username" = ''; 

-- +goose Down
ALTER TABLE "users" DROP COLUMN "username";
