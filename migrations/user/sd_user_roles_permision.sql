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

 Date: 06/10/2026 11:59:04
*/


-- ----------------------------
-- Table structure for sd_user_roles_permision
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_roles_permision";
CREATE TABLE "public"."sd_user_roles_permision" (
  "role_type_id" int4 NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "detail" text COLLATE "pg_catalog"."default",
  "created" timestamp(6) NOT NULL,
  "updated" timestamp(6),
  "insert" int8,
  "update" int8,
  "delete" int8,
  "select" int8,
  "log" int8,
  "config" int8,
  "truncate" int8
)
;

-- ----------------------------
-- Records of sd_user_roles_permision
-- ----------------------------
INSERT INTO "public"."sd_user_roles_permision" VALUES (1, 'System Admin', 'All', '2021-05-05 09:23:32', '2021-05-05 09:23:33', 1, 1, 1, 1, 1, 1, 1);
INSERT INTO "public"."sd_user_roles_permision" VALUES (2, 'Admin Level', 'Admin access', '2021-05-09 19:52:04', '2021-05-09 19:52:09', 1, 1, 1, 1, 1, 0, 0);
INSERT INTO "public"."sd_user_roles_permision" VALUES (3, 'Customer Level', 'Customer access', '2021-05-09 19:52:30', '2021-05-09 19:52:32', 1, 1, 0, 1, 1, 0, 0);

-- ----------------------------
-- Primary Key structure for table sd_user_roles_permision
-- ----------------------------
ALTER TABLE "public"."sd_user_roles_permision" ADD CONSTRAINT "PK_4df7386cc58a6712f2bef59c507" PRIMARY KEY ("role_type_id");
