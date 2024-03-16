CREATE TABLE "account" (
"id" bigserial PRIMARY KEY,
"first_name" varchar NOT NULL,
"last_name" varchar NOT NULL,
"user_name" varchar NOT NULL,
"password" varchar NOT NULL
);

CREATE TABLE "todo" (
"id" bigserial PRIMARY KEY,
"account_id" bigserial NOT NULL,
"title" varchar NOT NULL,
"time" varchar NOT NULL,
"date" varchar NOT NULL,
"complete" varchar NOT NULL
);

ALTER TABLE "todo" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");
