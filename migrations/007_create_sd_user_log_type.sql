DROP TABLE IF EXISTS "public"."sd_user_log_type" CASCADE;
CREATE SEQUENCE IF NOT EXISTS sd_user_log_type_log_type_id_seq;

CREATE TABLE "public"."sd_user_log_type" (
  "log_type_id"  int4 NOT NULL DEFAULT nextval('sd_user_log_type_log_type_id_seq'::regclass),
  "type_name"    text NOT NULL,
  "type_detail"  text NOT NULL,
  "status"       int8,
  "create"       timestamp(6) NOT NULL DEFAULT now(),
  "update"       timestamp(6) NOT NULL DEFAULT now()
);

ALTER TABLE "public"."sd_user_log_type"
  ADD CONSTRAINT "sd_user_log_type_pkey" PRIMARY KEY ("log_type_id");

INSERT INTO "public"."sd_user_log_type" (type_name, type_detail, status)
VALUES
  ('1',  'User SignIn',     1),
  ('2',  'User Signout',    1),
  ('3',  'User Create',     1),
  ('4',  'User Update',     1),
  ('5',  'User Delete',     1),
  ('6',  'Menu Create',     1),
  ('7',  'Menu Update',     1),
  ('8',  'Menu Delete',     1),
  ('9',  'Module Create',   1),
  ('10', 'Module Update',   1),
  ('11', 'Module Delete',   1),
  ('12', 'Menu Sort Update',1),
  ('13', 'Module Sort Update', 1)
ON CONFLICT DO NOTHING;