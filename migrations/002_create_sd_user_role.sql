DROP TABLE IF EXISTS "public"."sd_user_role" CASCADE;
CREATE TABLE "public"."sd_user_role" (
  "id" int4 NOT NULL DEFAULT nextval('sd_user_role_id_seq'::regclass),
  "role_id" int8 NOT NULL,
  "title" varchar(50),
  "createddate" timestamp(6),
  "updateddate" timestamp(6),
  "create_by" int8 NOT NULL,
  "lastupdate_by" int8 NOT NULL,
  "status" int2 NOT NULL,
  "type_id" int8 NOT NULL,
  "lang" varchar(255) NOT NULL
);

ALTER TABLE "public"."sd_user_role" ADD CONSTRAINT "sd_user_role_pkey" PRIMARY KEY ("id");