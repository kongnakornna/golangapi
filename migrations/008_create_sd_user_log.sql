DROP TABLE IF EXISTS "public"."sd_user_log" CASCADE;
CREATE SEQUENCE IF NOT EXISTS sd_user_log_id_seq;

CREATE TABLE "public"."sd_user_log" (
  "id"             int4 NOT NULL DEFAULT nextval('sd_user_log_id_seq'::regclass),
  "log_type_id"    int8 NOT NULL,
  "uid"            uuid NOT NULL,
  "name"           text NOT NULL,
  "detail"         text NOT NULL,
  "ipaddress"      text,
  "device_access"  text,
  "browser_agent"  text,
  "select_status"  int8,
  "insert_status"  int8,
  "update_status"  int8,
  "delete_status"  int8,
  "status"         int8,
  "create"         timestamp(6) NOT NULL DEFAULT now(),
  "update"         timestamp(6) NOT NULL DEFAULT now(),
  "lang"           text
);

ALTER TABLE "public"."sd_user_log"
  ADD CONSTRAINT "sd_user_log_pkey" PRIMARY KEY ("id", "log_type_id");

CREATE INDEX "idx_sd_user_log_uid"          ON "public"."sd_user_log" ("uid");
CREATE INDEX "idx_sd_user_log_log_type_id"  ON "public"."sd_user_log" ("log_type_id");
CREATE INDEX "idx_sd_user_log_create"       ON "public"."sd_user_log" ("create" DESC);
CREATE INDEX "idx_sd_user_log_ipaddress"    ON "public"."sd_user_log" ("ipaddress");