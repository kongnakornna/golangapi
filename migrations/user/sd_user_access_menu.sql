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

 Date: 06/10/2026 11:48:10
*/


-- ----------------------------
-- Table structure for sd_user_access_menu
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_access_menu";
CREATE TABLE "public"."sd_user_access_menu" (
  "user_access_id" int4 NOT NULL DEFAULT nextval('sd_user_access_menu_user_access_id_seq'::regclass),
  "user_type_id" int8,
  "menu_id" int8,
  "parent_id" int8
)
;

-- ----------------------------
-- Records of sd_user_access_menu
-- ----------------------------
INSERT INTO "public"."sd_user_access_menu" VALUES (1, 1, 1, 1);

-- ----------------------------
-- Primary Key structure for table sd_user_access_menu
-- ----------------------------
ALTER TABLE "public"."sd_user_access_menu" ADD CONSTRAINT "PK_b08610dd9113be8c7df7774dbc8" PRIMARY KEY ("user_access_id");
