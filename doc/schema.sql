-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2024-03-26T04:54:25.373Z

CREATE TABLE "account" (
  "id" bigserial PRIMARY KEY,
  "role" varchar NOT NULL DEFAULT 'user',
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "user_name" varchar UNIQUE NOT NULL,
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

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "account_id" bigserial NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "account" ("user_name");

CREATE UNIQUE INDEX ON "account" ("user_name");

CREATE INDEX ON "todo" ("account_id");

CREATE INDEX ON "sessions" ("account_id");

ALTER TABLE "todo" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");
