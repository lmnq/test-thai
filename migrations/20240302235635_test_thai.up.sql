CREATE TABLE IF NOT EXISTS "tbl_items" (
    "id" SERIAL PRIMARY KEY,
    "item_name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_tbl_items_item_name" ON "tbl_items" ("item_name");


CREATE TABLE IF NOT EXISTS "tbl_groups" (
    "id" SERIAL PRIMARY KEY,
    "group_name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_tbl_groups_group_name" ON "tbl_groups" ("group_name");


CREATE TABLE IF NOT EXISTS "tbl_categories" (
    "id" SERIAL PRIMARY KEY,
    "category_name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_tbl_categories_category_name" ON "tbl_categories" ("category_name");


CREATE TABLE IF NOT EXISTS "tbl_item_details" (
    "id" SERIAL PRIMARY KEY,
    "item_id" INTEGER NOT NULL,
    FOREIGN KEY ("item_id") REFERENCES "tbl_items" ("id") ON DELETE CASCADE,
    "category_id" INTEGER NOT NULL,
    FOREIGN KEY ("category_id") REFERENCES "tbl_categories" ("id") ON DELETE CASCADE,
    "group_id" INTEGER NOT NULL,
    FOREIGN KEY ("group_id") REFERENCES "tbl_groups" ("id") ON DELETE CASCADE,
    "cost" DECIMAL(10,2) NOT NULL,
    "price" DECIMAL(10,2) NOT NULL,
    "sort" INTEGER NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);
