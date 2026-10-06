DROP TABLE IF EXISTS "public"."sd_user_roles_access" CASCADE;
CREATE TABLE "public"."sd_user_roles_access" (
  "create" timestamp(6) NOT NULL,
  "update" timestamp(6) NOT NULL,
  "role_id" int4 NOT NULL,
  "role_type_id" int8 NOT NULL
);

ALTER TABLE "public"."sd_user_roles_access"
  ADD CONSTRAINT "sd_user_roles_access_pkey" PRIMARY KEY ("role_id", "role_type_id");