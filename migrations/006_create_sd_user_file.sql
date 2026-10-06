DROP TABLE IF EXISTS "public"."sd_user_file" CASCADE;
CREATE SEQUENCE IF NOT EXISTS sd_user_file_id_seq;

CREATE TABLE "public"."sd_user_file" (
  "id"           int4 NOT NULL DEFAULT nextval('sd_user_file_id_seq'::regclass),
  "file_name"    text NOT NULL,
  "file_type"    text NOT NULL,
  "file_path"    text NOT NULL,
  "file_type_id" int8 NOT NULL,
  "uid"          text,
  "file_date"    timestamp(6) NOT NULL DEFAULT now(),
  "status"       int2 NOT NULL
);

ALTER TABLE "public"."sd_user_file" ADD CONSTRAINT "sd_user_file_pkey" PRIMARY KEY ("id");
CREATE INDEX "idx_sd_user_file_uid" ON "public"."sd_user_file" ("uid");
CREATE INDEX "idx_sd_user_file_file_type_id" ON "public"."sd_user_file" ("file_type_id");
CREATE INDEX "idx_sd_user_file_status" ON "public"."sd_user_file" ("status");