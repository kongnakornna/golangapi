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

 Date: 06/10/2026 11:58:43
*/


-- ----------------------------
-- Table structure for sd_user_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_role";
CREATE TABLE "public"."sd_user_role" (
  "id" int4 NOT NULL DEFAULT nextval('sd_user_role_id_seq'::regclass),
  "role_id" int8 NOT NULL,
  "title" varchar(50) COLLATE "pg_catalog"."default",
  "createddate" timestamp(6),
  "updateddate" timestamp(6),
  "create_by" int8 NOT NULL,
  "lastupdate_by" int8 NOT NULL,
  "status" int2 NOT NULL,
  "type_id" int8 NOT NULL,
  "lang" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Records of sd_user_role
-- ----------------------------
INSERT INTO "public"."sd_user_role" VALUES (1, 1, 'Dev', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (2, 2, 'Administrator', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (3, 3, 'Company', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (4, 4, 'Staff', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (5, 5, 'Helpdask', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (6, 6, 'Customer', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (7, 7, 'Donate', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (8, 8, 'Edittor', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (9, 9, 'User', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');
INSERT INTO "public"."sd_user_role" VALUES (10, 10, 'gaust', '2025-04-06 20:56:00.301605', '2025-04-06 20:56:00.301605', 1, 1, 1, 1, 'en');

-- ----------------------------
-- Primary Key structure for table sd_user_role
-- ----------------------------
ALTER TABLE "public"."sd_user_role" ADD CONSTRAINT "PK_ce286bbce9874c345c85ba7c6e4" PRIMARY KEY ("id");
