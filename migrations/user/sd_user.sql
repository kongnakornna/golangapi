/*
 Navicat Premium Dump SQL

 Source Server         : postgres-Docker-5435
 Source Server Type    : PostgreSQL
 Source Server Version : 150019 (150019)
 Source Host           : localhost:5435
 Source Catalog        : icmongolang
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 150019 (150019)
 File Encoding         : 65001

 Date: 06/10/2026 12:00:00
*/


-- ----------------------------
-- Table structure for sd_user
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user";
CREATE TABLE "public"."sd_user" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "deletedate" date,
  "role_id" int8 NOT NULL,
  "email" text COLLATE "pg_catalog"."default" NOT NULL,
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "password_temp" text COLLATE "pg_catalog"."default",
  "firstname" text COLLATE "pg_catalog"."default",
  "lastname" text COLLATE "pg_catalog"."default",
  "fullname" text COLLATE "pg_catalog"."default",
  "nickname" text COLLATE "pg_catalog"."default",
  "idcard" text COLLATE "pg_catalog"."default",
  "lastsignindate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int2 NOT NULL,
  "active_status" int2,
  "network_id" int8 DEFAULT 1,
  "remark" text COLLATE "pg_catalog"."default",
  "infomation_agree_status" int2 DEFAULT 0,
  "gender" text COLLATE "pg_catalog"."default",
  "birthday" date,
  "online_status" text COLLATE "pg_catalog"."default" DEFAULT '0'::text,
  "message" text COLLATE "pg_catalog"."default",
  "network_type_id" int8 DEFAULT 0,
  "public_status" int2 DEFAULT 0,
  "type_id" int8 DEFAULT 0,
  "avatarpath" text COLLATE "pg_catalog"."default",
  "avatar" text COLLATE "pg_catalog"."default",
  "refresh_token" text COLLATE "pg_catalog"."default",
  "loginfailed" int2,
  "public_notification" int2 DEFAULT 0,
  "sms_notification" int2 DEFAULT 0,
  "email_notification" int2 DEFAULT 0,
  "line_notification" int2 DEFAULT 0,
  "mobile_number" text COLLATE "pg_catalog"."default" DEFAULT '0'::text,
  "phone_number" text COLLATE "pg_catalog"."default" DEFAULT '0'::text,
  "lineid" text COLLATE "pg_catalog"."default" DEFAULT '0'::text,
  "system_id" text COLLATE "pg_catalog"."default" DEFAULT '1'::text,
  "location_id" text COLLATE "pg_catalog"."default" DEFAULT '1'::text,
  "verified" bool DEFAULT false,
  "verification_code" varchar(64) COLLATE "pg_catalog"."default",
  "password_reset_token" varchar(64) COLLATE "pg_catalog"."default",
  "password_reset_at" timestamptz(6),
  "is_superuser" bool DEFAULT false
)
;

-- ----------------------------
-- Records of sd_user
-- ----------------------------
INSERT INTO "public"."sd_user" VALUES ('96731149-2dcd-4419-8d45-c85d91cb3724', '2026-08-22 19:03:51.223072+07', '2026-08-22 19:03:51.223072+07', NULL, 2, 'newuser@example.com', 'new_user', '$2a$10$HhftAkwsAHxykmCLGNvy.emzVKDcuvb8n0vB2rLEazrPGSGhcQH4.', 'password123', 'New', 'User', 'New User', NULL, NULL, '2026-08-22 19:03:51.223072+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, NULL, NULL, NULL, NULL, 0, 0, 0, 0, '0', '0', '0', '1', '1', 'f', '8f1cf28d0fc08cb72aec49d52e2e544b', 'daed31c8f8777f12a9a0971810f4ad6f', '2026-09-30 16:39:44.662766+07', 'f');
INSERT INTO "public"."sd_user" VALUES ('6057d5cb-fcc9-4226-9ea8-490a93670838', '2026-08-15 14:58:32.976431+07', '2026-08-15 14:58:32.976431+07', NULL, 2, 'uitest01@test.local', 'uitest01', '$2a$10$l59ZB/MR2z9OehcnrCDiCOroTbKKHWfUicPXJg9mnwZY/2/JRKeg2', 'password123', 'UI2', 'Test', 'UI Test', NULL, NULL, '2026-08-15 14:58:32.976431+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, NULL, NULL, NULL, NULL, 0, 0, 0, 0, '0899999999', '0', '0', '1', '1', 'f', '458b8ef9c25348f720c1b8b938bfc357', '', NULL, 'f');
INSERT INTO "public"."sd_user" VALUES ('6d762f44-deba-4a08-bd13-d8bfe385c5e8', '2026-10-03 09:01:41.84842+07', '2026-10-03 09:01:41.84842+07', NULL, 2, 'demo1@gmail.com', 'demo1', '$2a$10$ptDJwlYnM/BvUyWYKLjN6uBowIPm2nFy2F8dOUWSZIzoJPnRinjWG', NULL, 'demo1', 'demo1', 'demo1', NULL, NULL, '2026-10-03 09:01:41.84842+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, NULL, NULL, NULL, NULL, 0, 0, 0, 0, '0812345678', '021234567', 'demo1_line', '1', 'loc_001', 'f', '9a6d96653814f3bdcce94064224fa7ee', '', NULL, 'f');
INSERT INTO "public"."sd_user" VALUES ('934dcedc-0221-44a7-a0b4-fdaab3f774fd', '2026-08-25 20:11:33.972872+07', '2026-08-25 20:11:33.972872+07', NULL, 2, 'john@example.com', 'john_doe', '$2a$10$KpnSxNKs5tTDaOJXuoUvPO/CxIWcEqiEMjhp0PLqHcyLjpXwOq5nG', NULL, 'John', 'Doe', 'John Doe', NULL, NULL, '2026-08-25 20:11:33.972872+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, NULL, NULL, NULL, NULL, 0, 0, 0, 0, '0812345678', '0', '0', '1', '1', 'f', '1af7663d8a8e36be364c7c8f55846965', '', NULL, 'f');
INSERT INTO "public"."sd_user" VALUES ('8032b3ce-32a0-4b45-ac72-069201322639', '2026-08-27 10:46:26.259654+07', '2026-08-27 10:46:26.259654+07', NULL, 1, 'root@gmail.com', 'root@gmail.com', '$2a$10$khy6X.Mr2Vvb8b.NuiwInO.3phuK3OfWVKOJbr1yg0fI01P1r.yd2', NULL, NULL, NULL, 'root', NULL, NULL, '2026-08-27 10:46:26.259654+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, NULL, NULL, NULL, NULL, 0, 0, 0, 0, '0', '0', '0', '1', '1', 't', '', '', NULL, 't');
INSERT INTO "public"."sd_user" VALUES ('e9c08475-03bf-4b00-a7cd-1c58b54fba5e', '2026-10-03 09:02:20.018662+07', '2026-10-03 09:02:20.018662+07', NULL, 2, 'demoiot@gmail.com', 'demoiot', '$2a$10$jmuXrwR9iKCusddtNCtPwul27x3CVnFAgb78JxZNHbFGjWbn/dbYO', NULL, 'demoiot', 'demoiot', 'demoiot', NULL, NULL, '2026-10-03 09:02:20.018662+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, NULL, NULL, NULL, NULL, 0, 0, 0, 0, '0812345678', '021234567', 'demo1_line', '1', 'loc_001', 'f', '899b17fe4e271eefcacfc61944dda0a6', '', NULL, 'f');
INSERT INTO "public"."sd_user" VALUES ('3c8f7d9f-8298-4f01-bf6f-ecba1a5d423d', '2026-08-30 20:38:55.225655+07', '2026-10-04 21:06:58.494148+07', NULL, 1, 'kongnakornna@gmail.com', 'kongnakornna', '$2a$10$McG4Pmo55vwxvsi6NTSH2eH.VlmzNus45IG1CpAs8gbNwsGOUv9Fm', 'password', 'Kongnakorn', 'Jantakun', 'Kongnakorn Jantakun', NULL, NULL, '2026-08-30 20:38:55.225655+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, '/uploads/avatar/3c8f7d9f-8298-4f01-bf6f-ecba1a5d423d.jpg', '3c8f7d9f-8298-4f01-bf6f-ecba1a5d423d.jpg', NULL, NULL, 0, 0, 0, 0, '0955088091', '021234888', 'kongnakorn_line', '1', 'loc_001', 'f', '426f6c204fbbb6a238d042bc8a40eaf9', 'f248bd0c75778b598abcc452674d9772', '2026-09-27 19:19:56.667558+07', 't');
INSERT INTO "public"."sd_user" VALUES ('9e8365bd-0d08-4b62-8887-f42e628ad93e', '2026-08-15 11:08:49.000755+07', '2026-08-15 11:08:49.000755+07', NULL, 1, 'admin@icmongolang.local', 'admin@icmongolang.local', '$2a$10$NuNKBGcXMERsync3fgenvuWJ4gjCcRt00VYKvXdoVuUSHDghc3fBi', NULL, NULL, NULL, 'Admin', NULL, NULL, '2026-08-15 11:08:49.000755+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, NULL, NULL, NULL, NULL, 0, 0, 0, 0, '0', '0', '0', '1', '1', 't', '', '1402cddbf5cf3e41542e70663438eaf0', '2026-09-30 16:39:12.841566+07', 't');

-- ----------------------------
-- Indexes structure for table sd_user
-- ----------------------------
CREATE INDEX "idx_sd_user_createddate" ON "public"."sd_user" USING btree (
  "createddate" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE UNIQUE INDEX "idx_sd_user_email" ON "public"."sd_user" USING btree (
  "email" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sd_user_password_reset_token" ON "public"."sd_user" USING btree (
  "password_reset_token" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sd_user_status" ON "public"."sd_user" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_sd_user_username" ON "public"."sd_user" USING btree (
  "username" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sd_user_verification_code" ON "public"."sd_user" USING btree (
  "verification_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sd_user
-- ----------------------------
ALTER TABLE "public"."sd_user" ADD CONSTRAINT "sd_user_pkey" PRIMARY KEY ("id");
