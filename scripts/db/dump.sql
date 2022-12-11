CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "username" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "password" varchar NOT NULL,
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS categories (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS types (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS monsters (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "category_id" uuid NOT NULL,
  "description" text NOT NULL,
  "length" float8 NOT NULL,
  "weight" int NOT NULL,
  "hp" int NOT NULL,
  "attack" int NOT NULL,
  "defends" int NOT NULL,
  "speed" int NOT NULL,
  "catched" boolean NOT NULL DEFAULT false,
  "image" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS  monster_type (
  "monster_id" uuid NOT NULL,
  "type_id" uuid NOT NULL
);

ALTER TABLE "monsters" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "monster_type" ADD FOREIGN KEY ("monster_id") REFERENCES "monsters" ("id");

ALTER TABLE "monster_type" ADD FOREIGN KEY ("type_id") REFERENCES "types" ("id");

-- Seed data
INSERT INTO users (username, fullname, password, role) VALUES('admin', 'ADMIN', '5f4dcc3b5aa765d61d8327deb882cf99', 'admin'), ('user', 'USER', '5f4dcc3b5aa765d61d8327deb882cf99', 'user');

INSERT INTO categories (name) VALUES('Leaf Monster'), ('Diving Monster'), ('Lizard Monster');

INSERT INTO types (name) VALUES('GRASS'), ('PSYCHIC'), ('FLYING'), ('FIRE'), ('WATER'), ('ELECTRIC'), ('BUG');
