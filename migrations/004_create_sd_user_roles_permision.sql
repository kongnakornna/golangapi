DROP TABLE IF EXISTS "public"."sd_user_roles_permision" CASCADE;
CREATE TABLE "public"."sd_user_roles_permision" (
  "role_type_id" int4 NOT NULL,
  "name" text NOT NULL,
  "detail" text,
  "created" timestamp(6) NOT NULL,
  "updated" timestamp(6),
  "insert" int8,
  "update" int8,
  "delete" int8,
  "select" int8,
  "log" int8,
  "config" int8,
  "truncate" int8
);

ALTER TABLE "public"."sd_user_roles_permision"
  ADD CONSTRAINT "sd_user_roles_permision_pkey" PRIMARY KEY ("role_type_id");