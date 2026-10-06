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

 Date: 06/10/2026 11:58:52
*/


-- ----------------------------
-- Table structure for sd_user_roles_access
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_roles_access";
CREATE TABLE "public"."sd_user_roles_access" (
  "create" timestamp(6) NOT NULL,
  "update" timestamp(6) NOT NULL,
  "role_id" int4 NOT NULL,
  "role_type_id" int8 NOT NULL
)
;

-- ----------------------------
-- Records of sd_user_roles_access
-- ----------------------------
INSERT INTO "public"."sd_user_roles_access" VALUES ('2021-05-05 09:23:46', '2021-05-05 09:23:48', 1, 1);
INSERT INTO "public"."sd_user_roles_access" VALUES ('2023-05-01 17:50:20', '2023-05-01 17:50:25', 2, 2);
INSERT INTO "public"."sd_user_roles_access" VALUES ('2023-05-01 17:50:59', '2023-05-01 17:51:03', 3, 3);
INSERT INTO "public"."sd_user_roles_access" VALUES ('2023-05-01 17:51:13', '2023-05-01 17:51:15', 4, 4);
INSERT INTO "public"."sd_user_roles_access" VALUES ('2023-05-01 17:51:23', '2023-05-01 17:51:25', 5, 5);
INSERT INTO "public"."sd_user_roles_access" VALUES ('2023-05-01 17:51:31', '2023-05-01 17:51:34', 5, 6);

-- ----------------------------
-- Primary Key structure for table sd_user_roles_access
-- ----------------------------
ALTER TABLE "public"."sd_user_roles_access" ADD CONSTRAINT "PK_ea1374b87e00872215780b096f7" PRIMARY KEY ("role_id", "role_type_id");
