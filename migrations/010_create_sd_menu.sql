DROP TABLE IF EXISTS "public"."sd_menu" CASCADE;
CREATE SEQUENCE IF NOT EXISTS sd_menu_menu_id_seq;

CREATE TABLE "public"."sd_menu" (
  "menu_id"         int4 NOT NULL DEFAULT nextval('sd_menu_menu_id_seq'::regclass),
  "menu_parent_id"  int4,
  "modules_id"      int4 NOT NULL,
  "name"            text NOT NULL,
  "path"            text,
  "icon"            text,
  "sort"            int4 NOT NULL DEFAULT 0,
  "create_date"     timestamptz(6) NOT NULL DEFAULT now(),
  "update_date"     timestamptz(6) NOT NULL DEFAULT now(),
  "status"          int2 NOT NULL DEFAULT 1
);

ALTER TABLE "public"."sd_menu" ADD CONSTRAINT "sd_menu_pkey" PRIMARY KEY ("menu_id");

CREATE INDEX "idx_sd_menu_parent"  ON "public"."sd_menu" ("menu_parent_id");
CREATE INDEX "idx_sd_menu_modules" ON "public"."sd_menu" ("modules_id");
CREATE INDEX "idx_sd_menu_status"  ON "public"."sd_menu" ("status");
CREATE INDEX "idx_sd_menu_sort"    ON "public"."sd_menu" ("sort" ASC);

INSERT INTO "public"."sd_menu" (menu_parent_id, modules_id, name, path, icon, sort, status)
VALUES
  (NULL, 2, 'User List', '/users/list',  'list',    1, 1),
  (NULL, 2, 'Roles',     '/users/roles', 'shield',  2, 1),
  (NULL, 3, 'Devices',   '/iot/devices', 'devices', 1, 1),
  (NULL, 3, 'Charts',    '/iot/charts',  'chart',   2, 1)
ON CONFLICT DO NOTHING;