-- ============================================================
-- sd_user_access_menu — รองรับ role_id
-- ============================================================
DROP TABLE IF EXISTS "public"."sd_user_access_menu" CASCADE;
CREATE SEQUENCE IF NOT EXISTS sd_user_access_menu_user_access_id_seq;

CREATE TABLE "public"."sd_user_access_menu" (
  "user_access_id" int4 NOT NULL DEFAULT nextval('sd_user_access_menu_user_access_id_seq'::regclass),
  "user_type_id"   int8,
  "role_id"        int8,
  "menu_id"        int8,
  "parent_id"      int8
);

ALTER TABLE "public"."sd_user_access_menu"
  ADD CONSTRAINT "sd_user_access_menu_pkey" PRIMARY KEY ("user_access_id");

CREATE INDEX "idx_sd_user_access_menu_user_type" ON "public"."sd_user_access_menu" ("user_type_id");
CREATE INDEX "idx_sd_user_access_menu_role"      ON "public"."sd_user_access_menu" ("role_id");
CREATE INDEX "idx_sd_user_access_menu_menu"      ON "public"."sd_user_access_menu" ("menu_id");
CREATE INDEX "idx_sd_user_access_menu_parent"    ON "public"."sd_user_access_menu" ("parent_id");
CREATE INDEX "idx_sd_user_access_menu_role_menu" ON "public"."sd_user_access_menu" ("role_id", "menu_id");