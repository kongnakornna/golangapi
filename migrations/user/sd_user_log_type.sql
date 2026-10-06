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

 Date: 06/10/2026 11:58:33
*/


-- ----------------------------
-- Table structure for sd_user_log_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_log_type";
CREATE TABLE "public"."sd_user_log_type" (
  "log_type_id" int4 NOT NULL DEFAULT nextval('sd_user_log_type_log_type_id_seq'::regclass),
  "type_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "type_detail" text COLLATE "pg_catalog"."default" NOT NULL,
  "status" int8,
  "create" timestamp(6) NOT NULL DEFAULT now(),
  "update" timestamp(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_user_log_type
-- ----------------------------
INSERT INTO "public"."sd_user_log_type" VALUES (1, '1', 'User SignIn', 1, '2026-10-06 11:49:35', '2026-10-06 11:57:22');
INSERT INTO "public"."sd_user_log_type" VALUES (2, '2', 'User Signout', 1, '2026-10-06 11:58:06', '2026-10-06 11:58:08');

-- ----------------------------
-- Primary Key structure for table sd_user_log_type
-- ----------------------------
ALTER TABLE "public"."sd_user_log_type" ADD CONSTRAINT "PK_3f8b97a85e0528d6c18c4fd20b3" PRIMARY KEY ("log_type_id");
