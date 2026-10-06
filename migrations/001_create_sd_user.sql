DROP TABLE IF EXISTS "public"."sd_user" CASCADE;
CREATE TABLE "public"."sd_user" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "deletedate" date,
  "role_id" int8 NOT NULL,
  "email" text NOT NULL,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "password_temp" text,
  "firstname" text,
  "lastname" text,
  "fullname" text,
  "nickname" text,
  "idcard" text,
  "lastsignindate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int2 NOT NULL,
  "active_status" int2,
  "network_id" int8 DEFAULT 1,
  "remark" text,
  "infomation_agree_status" int2 DEFAULT 0,
  "gender" text,
  "birthday" date,
  "online_status" text DEFAULT '0',
  "message" text,
  "network_type_id" int8 DEFAULT 0,
  "public_status" int2 DEFAULT 0,
  "type_id" int8 DEFAULT 0,
  "avatarpath" text,
  "avatar" text,
  "refresh_token" text,
  "loginfailed" int2,
  "public_notification" int2 DEFAULT 0,
  "sms_notification" int2 DEFAULT 0,
  "email_notification" int2 DEFAULT 0,
  "line_notification" int2 DEFAULT 0,
  "mobile_number" text DEFAULT '0',
  "phone_number" text DEFAULT '0',
  "lineid" text DEFAULT '0',
  "system_id" text DEFAULT '1',
  "location_id" text DEFAULT '1',
  "verified" bool DEFAULT false,
  "verification_code" varchar(64),
  "password_reset_token" varchar(64),
  "password_reset_at" timestamptz(6),
  "is_superuser" bool DEFAULT false
);

ALTER TABLE "public"."sd_user" ADD CONSTRAINT "sd_user_pkey" PRIMARY KEY ("id");
CREATE UNIQUE INDEX "idx_sd_user_email" ON "public"."sd_user" ("email");
CREATE UNIQUE INDEX "idx_sd_user_username" ON "public"."sd_user" ("username");
CREATE INDEX "idx_sd_user_status" ON "public"."sd_user" ("status");
CREATE INDEX "idx_sd_user_createddate" ON "public"."sd_user" ("createddate" DESC);