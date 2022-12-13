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

CREATE TABLE IF NOT EXISTS  monster_types (
  "monster_id" uuid NOT NULL,
  "type_id" uuid NOT NULL
);

ALTER TABLE "monsters" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "monster_types" ADD FOREIGN KEY ("monster_id") REFERENCES "monsters" ("id");

ALTER TABLE "monster_types" ADD FOREIGN KEY ("type_id") REFERENCES "types" ("id");

-- Delete data
DELETE FROM users;
DELETE FROM categories;
DELETE FROM types;

-- Seed data
INSERT INTO users (username, fullname, password, role) VALUES('admin', 'ADMIN', '$2a$04$euYwgSigV4MDtKR0pvnBXumov0IsFsfumR0fsjgwGcEqXNOpmp0Ju', 'admin'), ('user', 'USER', '$2a$04$yYhf5Y3wsZoYmlGWc.uX8OCfgA2oJgGl5GX73n5rvRlUpZQtOuOFG', 'user');

INSERT INTO categories (name) VALUES('Leaf Monster'), ('Diving Monster'), ('Lizard Monster');

INSERT INTO types (name) VALUES('GRASS'), ('PSYCHIC'), ('FLYING'), ('FIRE'), ('WATER'), ('ELECTRIC'), ('BUG');

-- Select
SELECT * FROM users;
SELECT * FROM categories;
SELECT * FROM types;