-- +goose Up
CREATE TABLE IF NOT EXISTS "posts" (
    "id" INTEGER PRIMARY KEY,
    "title" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "user_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX "idx_user_post" ON "posts"("title", "user_id");

CREATE TABLE IF NOT EXISTS "post_likes" (
    "user_id" INTEGER NOT NULL,
    "post_id" INTEGER NOT NULL,
    "liked_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("user_id", "post_id")
);

CREATE TABLE IF NOT EXISTS "post_comments" (
    "id" INTEGER PRIMARY KEY,
    "post_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "comment" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE,  
    FOREIGN KEY ("post_id") REFERENCES "posts"("id") ON DELETE CASCADE  
);

-- +goose Down
DROP INDEX "idx_user_post";
DROP TABLE IF EXISTS "post_comments";
DROP TABLE IF EXISTS "post_likes";
DROP TABLE IF EXISTS "posts";
