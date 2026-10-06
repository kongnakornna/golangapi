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

 Date: 06/10/2026 11:48:36
*/


-- ----------------------------
-- Table structure for sd_user_file
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_file";
CREATE TABLE "public"."sd_user_file" (
  "id" int4 NOT NULL DEFAULT nextval('sd_user_file_id_seq'::regclass),
  "file_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "file_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "file_path" text COLLATE "pg_catalog"."default" NOT NULL,
  "file_type_id" int8 NOT NULL,
  "uid" text COLLATE "pg_catalog"."default",
  "file_date" timestamp(6) NOT NULL DEFAULT now(),
  "status" int2 NOT NULL
)
;

-- ----------------------------
-- Records of sd_user_file
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table sd_user_file
-- ----------------------------
ALTER TABLE "public"."sd_user_file" ADD CONSTRAINT "PK_bee867c384da15706056a6d4d79" PRIMARY KEY ("id");
