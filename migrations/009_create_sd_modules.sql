DROP TABLE IF EXISTS "public"."sd_modules" CASCADE;
CREATE SEQUENCE IF NOT EXISTS sd_modules_module_id_seq;

CREATE TABLE "public"."sd_modules" (
  "module_id"         int4 NOT NULL DEFAULT nextval('sd_modules_module_id_seq'::regclass),
  "module_parent_id"  int4,
  "name"              text NOT NULL,
  "path"              text,
  "icon"              text,
  "sort"              int4 NOT NULL DEFAULT 0,
  "create_date"       timestamptz(6) NOT NULL DEFAULT now(),
  "update_date"       timestamptz(6) NOT NULL DEFAULT now(),
  "status"            int2 NOT NULL DEFAULT 1
);

ALTER TABLE "public"."sd_modules" ADD CONSTRAINT "sd_modules_pkey" PRIMARY KEY ("module_id");

CREATE INDEX "idx_sd_modules_parent" ON "public"."sd_modules" ("module_parent_id");
CREATE INDEX "idx_sd_modules_status" ON "public"."sd_modules" ("status");
CREATE INDEX "idx_sd_modules_sort"   ON "public"."sd_modules" ("sort" ASC);

INSERT INTO "public"."sd_modules" (module_parent_id, name, path, icon, sort, status)
VALUES
  (NULL, 'Dashboard', '/dashboard', 'dashboard', 1, 1),
  (NULL, 'Users',     '/users',     'people',    2, 1),
  (NULL, 'IOT',       '/iot',       'sensors',   3, 1),
  (NULL, 'Reports',   '/reports',   'bar-chart', 4, 1)
ON CONFLICT DO NOTHING;