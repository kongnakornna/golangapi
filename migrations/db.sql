/*
 Navicat Premium Dump SQL

 Source Server         : 127.0.0.1
 Source Server Type    : PostgreSQL
 Source Server Version : 170006 (170006)
 Source Host           : 127.0.0.1:5432
 Source Catalog        : db
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 170006 (170006)
 File Encoding         : 65001

 Date: 06/04/2026 16:34:11
*/


-- ----------------------------
-- Type structure for seg
-- ----------------------------
DROP TYPE IF EXISTS "public"."seg";
CREATE TYPE "public"."seg" (
  INPUT = "public"."seg_in",
  OUTPUT = "public"."seg_out",
  INTERNALLENGTH = 12,
  CATEGORY = U,
  DELIMITER = ','
);
ALTER TYPE "public"."seg" OWNER TO "postgres";
COMMENT ON TYPE "public"."seg" IS 'floating point interval ''FLOAT .. FLOAT'', ''.. FLOAT'', ''FLOAT ..'' or ''FLOAT''';

-- ----------------------------
-- Type structure for user_role_enum
-- ----------------------------
DROP TYPE IF EXISTS "public"."user_role_enum";
CREATE TYPE "public"."user_role_enum" AS ENUM (
  'SUPERADMIN',
  'ADMIN',
  'EDITOR',
  'MONITOR',
  'USER'
);
ALTER TYPE "public"."user_role_enum" OWNER TO "postgres";

-- ----------------------------
-- Type structure for user_usertype_enum
-- ----------------------------
DROP TYPE IF EXISTS "public"."user_usertype_enum";
CREATE TYPE "public"."user_usertype_enum" AS ENUM (
  'therapist',
  'supervisor',
  'superadmin',
  'system',
  'admin',
  'support',
  'enduser'
);
ALTER TYPE "public"."user_usertype_enum" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for activity_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."activity_log_id_seq";
CREATE SEQUENCE "public"."activity_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for activity_log_id_seq1
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."activity_log_id_seq1";
CREATE SEQUENCE "public"."activity_log_id_seq1" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for command_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."command_log_id_seq";
CREATE SEQUENCE "public"."command_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for device_alert_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."device_alert_id_seq";
CREATE SEQUENCE "public"."device_alert_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for device_config_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."device_config_id_seq";
CREATE SEQUENCE "public"."device_config_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for device_status_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."device_status_id_seq";
CREATE SEQUENCE "public"."device_status_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for iot_data_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."iot_data_id_seq";
CREATE SEQUENCE "public"."iot_data_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for noti_notification_logs_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."noti_notification_logs_log_id_seq";
CREATE SEQUENCE "public"."noti_notification_logs_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for noti_notification_rules_rule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."noti_notification_rules_rule_id_seq";
CREATE SEQUENCE "public"."noti_notification_rules_rule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for noti_notification_types_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."noti_notification_types_type_id_seq";
CREATE SEQUENCE "public"."noti_notification_types_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for notification_devices_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."notification_devices_id_seq";
CREATE SEQUENCE "public"."notification_devices_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for notification_groups_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."notification_groups_id_seq";
CREATE SEQUENCE "public"."notification_groups_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for notification_logs_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."notification_logs_id_seq";
CREATE SEQUENCE "public"."notification_logs_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for notification_types_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."notification_types_id_seq";
CREATE SEQUENCE "public"."notification_types_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_activity_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_activity_log_id_seq";
CREATE SEQUENCE "public"."sd_activity_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_activity_type_log_typeId_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_activity_type_log_typeId_seq";
CREATE SEQUENCE "public"."sd_activity_type_log_typeId_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_admin_access_menu_admin_access_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_admin_access_menu_admin_access_id_seq";
CREATE SEQUENCE "public"."sd_admin_access_menu_admin_access_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_air_control_air_control_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_air_control_air_control_id_seq";
CREATE SEQUENCE "public"."sd_air_control_air_control_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_air_mod_air_mod_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_air_mod_air_mod_id_seq";
CREATE SEQUENCE "public"."sd_air_mod_air_mod_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_air_period_air_period_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_air_period_air_period_id_seq";
CREATE SEQUENCE "public"."sd_air_period_air_period_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_air_setting_warning_air_setting_warning_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_air_setting_warning_air_setting_warning_id_seq";
CREATE SEQUENCE "public"."sd_air_setting_warning_air_setting_warning_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_air_warning_air_warning_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_air_warning_air_warning_id_seq";
CREATE SEQUENCE "public"."sd_air_warning_air_warning_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_api_key_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_api_key_id_seq";
CREATE SEQUENCE "public"."sd_api_key_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_audit_log_audit_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_audit_log_audit_id_seq";
CREATE SEQUENCE "public"."sd_audit_log_audit_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_channel_template_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_channel_template_id_seq";
CREATE SEQUENCE "public"."sd_channel_template_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_dashboard_config_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_dashboard_config_id_seq";
CREATE SEQUENCE "public"."sd_dashboard_config_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_device_category_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_device_category_id_seq";
CREATE SEQUENCE "public"."sd_device_category_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_device_group_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_device_group_id_seq";
CREATE SEQUENCE "public"."sd_device_group_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_device_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_device_log_id_seq";
CREATE SEQUENCE "public"."sd_device_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_device_member_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_device_member_id_seq";
CREATE SEQUENCE "public"."sd_device_member_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_device_notification_config_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_device_notification_config_id_seq";
CREATE SEQUENCE "public"."sd_device_notification_config_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_device_schedule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_device_schedule_id_seq";
CREATE SEQUENCE "public"."sd_device_schedule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_device_status_history_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_device_status_history_id_seq";
CREATE SEQUENCE "public"."sd_device_status_history_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_group_notification_config_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_group_notification_config_id_seq";
CREATE SEQUENCE "public"."sd_group_notification_config_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_api_api_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_api_api_id_seq";
CREATE SEQUENCE "public"."sd_iot_api_api_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_device_action_device_action_user_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_device_action_device_action_user_id_seq";
CREATE SEQUENCE "public"."sd_iot_device_action_device_action_user_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_device_action_log_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_device_action_log_log_id_seq";
CREATE SEQUENCE "public"."sd_iot_device_action_log_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_device_action_user_device_action_user_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_device_action_user_device_action_user_id_seq";
CREATE SEQUENCE "public"."sd_iot_device_action_user_device_action_user_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_device_alarm_action_alarm_action_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_device_alarm_action_alarm_action_id_seq";
CREATE SEQUENCE "public"."sd_iot_device_alarm_action_alarm_action_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_device_device_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_device_device_id_seq";
CREATE SEQUENCE "public"."sd_iot_device_device_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_device_type_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_device_type_type_id_seq";
CREATE SEQUENCE "public"."sd_iot_device_type_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_group_group_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_group_group_id_seq";
CREATE SEQUENCE "public"."sd_iot_group_group_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_location_location_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_location_location_id_seq";
CREATE SEQUENCE "public"."sd_iot_location_location_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_mqtt_mqtt_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_mqtt_mqtt_id_seq";
CREATE SEQUENCE "public"."sd_iot_mqtt_mqtt_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_schedule_schedule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_schedule_schedule_id_seq";
CREATE SEQUENCE "public"."sd_iot_schedule_schedule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_sensor_sensor_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_sensor_sensor_id_seq";
CREATE SEQUENCE "public"."sd_iot_sensor_sensor_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_setting_setting_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_setting_setting_id_seq";
CREATE SEQUENCE "public"."sd_iot_setting_setting_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_token_token_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_token_token_id_seq";
CREATE SEQUENCE "public"."sd_iot_token_token_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_iot_type_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_iot_type_type_id_seq";
CREATE SEQUENCE "public"."sd_iot_type_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_module_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_module_log_id_seq";
CREATE SEQUENCE "public"."sd_module_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_notification_channel_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_notification_channel_id_seq";
CREATE SEQUENCE "public"."sd_notification_channel_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_notification_condition_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_notification_condition_id_seq";
CREATE SEQUENCE "public"."sd_notification_condition_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_notification_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_notification_log_id_seq";
CREATE SEQUENCE "public"."sd_notification_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_notification_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_notification_type_id_seq";
CREATE SEQUENCE "public"."sd_notification_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_report_data_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_report_data_id_seq";
CREATE SEQUENCE "public"."sd_report_data_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_sensor_data_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_sensor_data_id_seq";
CREATE SEQUENCE "public"."sd_sensor_data_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_system_setting_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_system_setting_id_seq";
CREATE SEQUENCE "public"."sd_system_setting_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_user_access_menu_user_access_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_user_access_menu_user_access_id_seq";
CREATE SEQUENCE "public"."sd_user_access_menu_user_access_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_user_file_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_user_file_id_seq";
CREATE SEQUENCE "public"."sd_user_file_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_user_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_user_log_id_seq";
CREATE SEQUENCE "public"."sd_user_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_user_log_type_log_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_user_log_type_log_type_id_seq";
CREATE SEQUENCE "public"."sd_user_log_type_log_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_user_role_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_user_role_id_seq";
CREATE SEQUENCE "public"."sd_user_role_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sd_user_roles_permision_role_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sd_user_roles_permision_role_type_id_seq";
CREATE SEQUENCE "public"."sd_user_roles_permision_role_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for tnb_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."tnb_id_seq";
CREATE SEQUENCE "public"."tnb_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for activity_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."activity_log";
CREATE TABLE "public"."activity_log" (
  "id" int8 NOT NULL DEFAULT nextval('activity_log_id_seq1'::regclass),
  "type" text COLLATE "pg_catalog"."default" NOT NULL,
  "deviceId" text COLLATE "pg_catalog"."default",
  "userId" text COLLATE "pg_catalog"."default",
  "details" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" jsonb,
  "severity" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'info'::text,
  "ipAddress" text COLLATE "pg_catalog"."default",
  "userAgent" text COLLATE "pg_catalog"."default",
  "sessionId" text COLLATE "pg_catalog"."default",
  "correlationId" text COLLATE "pg_catalog"."default",
  "timestamp" timestamptz(6) NOT NULL,
  "createdAt" timestamptz(6) NOT NULL DEFAULT now(),
  "stackTrace" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of activity_log
-- ----------------------------

-- ----------------------------
-- Table structure for command_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."command_log";
CREATE TABLE "public"."command_log" (
  "id" int8 NOT NULL DEFAULT nextval('command_log_id_seq'::regclass),
  "deviceId" text COLLATE "pg_catalog"."default" NOT NULL,
  "action" text COLLATE "pg_catalog"."default" NOT NULL,
  "parameters" jsonb,
  "metadata" jsonb,
  "status" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'pending'::text,
  "issuedBy" text COLLATE "pg_catalog"."default",
  "clientIp" text COLLATE "pg_catalog"."default",
  "response" jsonb,
  "error" text COLLATE "pg_catalog"."default",
  "issuedAt" timestamptz(6) NOT NULL,
  "sentAt" timestamptz(6),
  "executedAt" timestamptz(6),
  "failedAt" timestamptz(6),
  "createdAt" timestamptz(6) NOT NULL DEFAULT now(),
  "updatedAt" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of command_log
-- ----------------------------

-- ----------------------------
-- Table structure for device_alert
-- ----------------------------
DROP TABLE IF EXISTS "public"."device_alert";
CREATE TABLE "public"."device_alert" (
  "id" int8 NOT NULL DEFAULT nextval('device_alert_id_seq'::regclass),
  "deviceId" text COLLATE "pg_catalog"."default" NOT NULL,
  "type" text COLLATE "pg_catalog"."default" NOT NULL,
  "metric" text COLLATE "pg_catalog"."default",
  "value" float8,
  "threshold" jsonb,
  "severity" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'low'::text,
  "message" text COLLATE "pg_catalog"."default" NOT NULL,
  "details" jsonb,
  "resolved" bool NOT NULL DEFAULT false,
  "resolutionNotes" text COLLATE "pg_catalog"."default",
  "resolvedBy" text COLLATE "pg_catalog"."default",
  "resolvedAt" timestamptz(6),
  "acknowledged" bool NOT NULL DEFAULT false,
  "acknowledgedBy" text COLLATE "pg_catalog"."default",
  "acknowledgedAt" timestamptz(6),
  "escalation" jsonb,
  "dataId" int8,
  "createdAt" timestamptz(6) NOT NULL DEFAULT now(),
  "updatedAt" timestamptz(6) NOT NULL DEFAULT now(),
  "expiresAt" timestamptz(6),
  "notificationCount" int8 NOT NULL DEFAULT 0
)
;

-- ----------------------------
-- Records of device_alert
-- ----------------------------

-- ----------------------------
-- Table structure for device_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."device_config";
CREATE TABLE "public"."device_config" (
  "id" int8 NOT NULL DEFAULT nextval('device_config_id_seq'::regclass),
  "deviceId" text COLLATE "pg_catalog"."default" NOT NULL,
  "config" jsonb,
  "status" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'active'::text,
  "notes" text COLLATE "pg_catalog"."default",
  "updatedBy" text COLLATE "pg_catalog"."default",
  "createdAt" timestamptz(6) NOT NULL DEFAULT now(),
  "updatedAt" timestamptz(6) NOT NULL DEFAULT now(),
  "lastAppliedAt" timestamptz(6)
)
;

-- ----------------------------
-- Records of device_config
-- ----------------------------

-- ----------------------------
-- Table structure for device_status
-- ----------------------------
DROP TABLE IF EXISTS "public"."device_status";
CREATE TABLE "public"."device_status" (
  "id" int8 NOT NULL DEFAULT nextval('device_status_id_seq'::regclass),
  "deviceId" text COLLATE "pg_catalog"."default" NOT NULL,
  "isOnline" bool NOT NULL DEFAULT true,
  "isActive" bool NOT NULL DEFAULT true,
  "lastSeen" timestamptz(6) NOT NULL,
  "lastData" jsonb,
  "batteryLevel" int8,
  "signalStrength" int8,
  "temperature" float8,
  "humidity" float8,
  "firmwareVersion" text COLLATE "pg_catalog"."default",
  "uptime" int8,
  "location" jsonb,
  "networkInfo" jsonb,
  "hardwareInfo" jsonb,
  "metrics" jsonb,
  "statusMessage" text COLLATE "pg_catalog"."default",
  "customFields" jsonb,
  "createdAt" timestamptz(6) NOT NULL DEFAULT now(),
  "updatedAt" timestamptz(6) NOT NULL DEFAULT now(),
  "firstSeen" timestamptz(6),
  "lastMaintenance" timestamptz(6),
  "connectionCount" int8 NOT NULL DEFAULT 0
)
;

-- ----------------------------
-- Records of device_status
-- ----------------------------

-- ----------------------------
-- Table structure for iot_data
-- ----------------------------
DROP TABLE IF EXISTS "public"."iot_data";
CREATE TABLE "public"."iot_data" (
  "id" int8 NOT NULL DEFAULT nextval('iot_data_id_seq'::regclass),
  "data" jsonb NOT NULL,
  "createdAt" timestamptz(6) NOT NULL DEFAULT now(),
  "location" jsonb,
  "metadata" jsonb,
  "dataType" text COLLATE "pg_catalog"."default",
  "dataQuality" float8,
  "deviceId" text COLLATE "pg_catalog"."default" NOT NULL,
  "timestamp" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of iot_data
-- ----------------------------

-- ----------------------------
-- Table structure for item
-- ----------------------------
DROP TABLE IF EXISTS "public"."item";
CREATE TABLE "public"."item" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "title" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "description" varchar(200) COLLATE "pg_catalog"."default" NOT NULL,
  "owner_id" uuid
)
;

-- ----------------------------
-- Records of item
-- ----------------------------
INSERT INTO "public"."item" VALUES ('90394643-f5d3-43c4-8610-58882d60a13e', '2026-04-06 16:33:09.791098+07', '2026-04-06 16:33:09.791098+07', NULL, 'item 1', 'item1 description', 'e32af9b1-a773-4757-9a03-09c3c01dc0ce');

-- ----------------------------
-- Table structure for noti_notification_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."noti_notification_logs";
CREATE TABLE "public"."noti_notification_logs" (
  "log_id" int8 NOT NULL DEFAULT nextval('noti_notification_logs_log_id_seq'::regclass),
  "notification_id" uuid NOT NULL,
  "channel" text COLLATE "pg_catalog"."default" NOT NULL,
  "payload" jsonb NOT NULL,
  "response" jsonb,
  "status" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'pending'::text,
  "retry_count" int8 DEFAULT 0,
  "error_message" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) DEFAULT now(),
  "sent_at" timestamptz(6),
  "delivered_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of noti_notification_logs
-- ----------------------------

-- ----------------------------
-- Table structure for noti_notification_rules
-- ----------------------------
DROP TABLE IF EXISTS "public"."noti_notification_rules";
CREATE TABLE "public"."noti_notification_rules" (
  "rule_id" int8 NOT NULL DEFAULT nextval('noti_notification_rules_rule_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default" NOT NULL,
  "event_trigger" text COLLATE "pg_catalog"."default" NOT NULL,
  "conditions" jsonb NOT NULL,
  "actions" jsonb NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "priority" int8 NOT NULL DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of noti_notification_rules
-- ----------------------------

-- ----------------------------
-- Table structure for noti_notification_types
-- ----------------------------
DROP TABLE IF EXISTS "public"."noti_notification_types";
CREATE TABLE "public"."noti_notification_types" (
  "type_id" int8 NOT NULL DEFAULT nextval('noti_notification_types_type_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default" NOT NULL,
  "default_template" jsonb,
  "allowed_channels" jsonb,
  "status" int8 NOT NULL DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of noti_notification_types
-- ----------------------------

-- ----------------------------
-- Table structure for noti_notifications
-- ----------------------------
DROP TABLE IF EXISTS "public"."noti_notifications";
CREATE TABLE "public"."noti_notifications" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "title" text COLLATE "pg_catalog"."default" NOT NULL,
  "message" text COLLATE "pg_catalog"."default" NOT NULL,
  "type" text COLLATE "pg_catalog"."default" NOT NULL,
  "priority" text COLLATE "pg_catalog"."default" NOT NULL,
  "category" text COLLATE "pg_catalog"."default",
  "user_id" int8,
  "user_uuid" uuid,
  "metadata" jsonb,
  "is_read" bool DEFAULT false,
  "read_at" timestamptz(6),
  "is_sent" bool DEFAULT false,
  "channels_sent" jsonb,
  "scheduled_at" timestamptz(6),
  "expires_at" timestamptz(6),
  "status" int8 NOT NULL DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of noti_notifications
-- ----------------------------

-- ----------------------------
-- Table structure for notification_devices
-- ----------------------------
DROP TABLE IF EXISTS "public"."notification_devices";
CREATE TABLE "public"."notification_devices" (
  "id" int8 NOT NULL DEFAULT nextval('notification_devices_id_seq'::regclass),
  "device_id" int8,
  "device_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "device_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "mqtt_topic" text COLLATE "pg_catalog"."default",
  "mqtt_on" text COLLATE "pg_catalog"."default",
  "mqtt_off" text COLLATE "pg_catalog"."default",
  "location" text COLLATE "pg_catalog"."default",
  "unit" text COLLATE "pg_catalog"."default",
  "last_value" text COLLATE "pg_catalog"."default",
  "last_status" int8,
  "last_updated" timestamptz(6),
  "is_online" bool NOT NULL DEFAULT false,
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of notification_devices
-- ----------------------------

-- ----------------------------
-- Table structure for notification_groups
-- ----------------------------
DROP TABLE IF EXISTS "public"."notification_groups";
CREATE TABLE "public"."notification_groups" (
  "id" int8 NOT NULL DEFAULT nextval('notification_groups_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "device_ids" json,
  "isActive" bool NOT NULL DEFAULT true,
  "createdAt" timestamptz(6) NOT NULL DEFAULT now(),
  "updatedAt" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of notification_groups
-- ----------------------------

-- ----------------------------
-- Table structure for notification_groups_devices_notification_devices
-- ----------------------------
DROP TABLE IF EXISTS "public"."notification_groups_devices_notification_devices";
CREATE TABLE "public"."notification_groups_devices_notification_devices" (
  "notificationGroupsId" int8 NOT NULL,
  "notificationDevicesId" int8 NOT NULL
)
;

-- ----------------------------
-- Records of notification_groups_devices_notification_devices
-- ----------------------------

-- ----------------------------
-- Table structure for notification_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."notification_logs";
CREATE TABLE "public"."notification_logs" (
  "id" int8 NOT NULL DEFAULT nextval('notification_logs_id_seq'::regclass),
  "device_id" int8 NOT NULL,
  "device_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "device_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "value_data" text COLLATE "pg_catalog"."default" NOT NULL,
  "numeric_value" float8,
  "notification_type_id" int8 NOT NULL,
  "status" int8 NOT NULL,
  "title" text COLLATE "pg_catalog"."default" NOT NULL,
  "message" text COLLATE "pg_catalog"."default" NOT NULL,
  "channels_sent" jsonb NOT NULL,
  "control_action" jsonb,
  "redis_key" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "config_id" int8
)
;

-- ----------------------------
-- Records of notification_logs
-- ----------------------------

-- ----------------------------
-- Table structure for notification_types
-- ----------------------------
DROP TABLE IF EXISTS "public"."notification_types";
CREATE TABLE "public"."notification_types" (
  "id" int8 NOT NULL DEFAULT nextval('notification_types_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "code" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default" NOT NULL,
  "icon" text COLLATE "pg_catalog"."default",
  "color" text COLLATE "pg_catalog"."default",
  "repeat_cooldown" int8 NOT NULL DEFAULT 10,
  "should_notify" bool NOT NULL DEFAULT true
)
;

-- ----------------------------
-- Records of notification_types
-- ----------------------------

-- ----------------------------
-- Table structure for sd_activity_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_activity_log";
CREATE TABLE "public"."sd_activity_log" (
  "id" int8 NOT NULL DEFAULT nextval('sd_activity_log_id_seq'::regclass),
  "user_id" text COLLATE "pg_catalog"."default",
  "type_id" int8,
  "modules_id" int8,
  "name" text COLLATE "pg_catalog"."default",
  "event" text COLLATE "pg_catalog"."default",
  "detail" text COLLATE "pg_catalog"."default",
  "location" text COLLATE "pg_catalog"."default",
  "date" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_activity_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_activity_type_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_activity_type_log";
CREATE TABLE "public"."sd_activity_type_log" (
  "typeId" int8 NOT NULL DEFAULT nextval('"sd_activity_type_log_typeId_seq"'::regclass),
  "type_name" text COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Records of sd_activity_type_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_admin_access_menu
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_admin_access_menu";
CREATE TABLE "public"."sd_admin_access_menu" (
  "admin_access_id" int8 NOT NULL DEFAULT nextval('sd_admin_access_menu_admin_access_id_seq'::regclass),
  "admin_type_id" int8,
  "admin_menu_id" int8
)
;

-- ----------------------------
-- Records of sd_admin_access_menu
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_control
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_control";
CREATE TABLE "public"."sd_air_control" (
  "air_control_id" int8 NOT NULL DEFAULT nextval('sd_air_control_air_control_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default",
  "data" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default",
  "active" int8,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_air_control
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_control_device_map
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_control_device_map";
CREATE TABLE "public"."sd_air_control_device_map" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "air_control_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_air_control_device_map
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_control_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_control_log";
CREATE TABLE "public"."sd_air_control_log" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "air_control_id" int8,
  "device_id" int8,
  "type_id" int8,
  "temperature" text COLLATE "pg_catalog"."default",
  "warning" text COLLATE "pg_catalog"."default",
  "recovery" text COLLATE "pg_catalog"."default",
  "period" text COLLATE "pg_catalog"."default",
  "percent" text COLLATE "pg_catalog"."default",
  "firealarm" text COLLATE "pg_catalog"."default",
  "humidityalarm" text COLLATE "pg_catalog"."default",
  "air2_alarm" text COLLATE "pg_catalog"."default",
  "air1_alarm" text COLLATE "pg_catalog"."default",
  "temperaturealarm" text COLLATE "pg_catalog"."default",
  "mode" text COLLATE "pg_catalog"."default",
  "state_air1" text COLLATE "pg_catalog"."default",
  "state_air2" text COLLATE "pg_catalog"."default",
  "temperaturealarmoff" text COLLATE "pg_catalog"."default",
  "ups_alarm" text COLLATE "pg_catalog"."default",
  "ups2_alarm" text COLLATE "pg_catalog"."default",
  "hssdalarm" text COLLATE "pg_catalog"."default",
  "waterleakalarm" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_air_control_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_mod
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_mod";
CREATE TABLE "public"."sd_air_mod" (
  "air_mod_id" int8 NOT NULL DEFAULT nextval('sd_air_mod_air_mod_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default",
  "data" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default",
  "active" int8,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_air_mod
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_mod_device_map
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_mod_device_map";
CREATE TABLE "public"."sd_air_mod_device_map" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "air_mod_id" int8,
  "air_control_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_air_mod_device_map
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_period
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_period";
CREATE TABLE "public"."sd_air_period" (
  "air_period_id" int8 NOT NULL DEFAULT nextval('sd_air_period_air_period_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default",
  "data" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default",
  "active" int8,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_air_period
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_period_device_map
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_period_device_map";
CREATE TABLE "public"."sd_air_period_device_map" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "air_period_id" int8,
  "air_control_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_air_period_device_map
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_setting_warning
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_setting_warning";
CREATE TABLE "public"."sd_air_setting_warning" (
  "air_setting_warning_id" int8 NOT NULL DEFAULT nextval('sd_air_setting_warning_air_setting_warning_id_seq'::regclass),
  "type_id" int8,
  "device_id" int8,
  "period_id" int8,
  "event_name" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default",
  "active" int8,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_air_setting_warning
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_setting_warning_device_map
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_setting_warning_device_map";
CREATE TABLE "public"."sd_air_setting_warning_device_map" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "air_setting_warning_id" int8,
  "air_control_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_air_setting_warning_device_map
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_warning
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_warning";
CREATE TABLE "public"."sd_air_warning" (
  "air_warning_id" int8 NOT NULL DEFAULT nextval('sd_air_warning_air_warning_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default",
  "data" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default",
  "active" int8,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_air_warning
-- ----------------------------

-- ----------------------------
-- Table structure for sd_air_warning_device_map
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_air_warning_device_map";
CREATE TABLE "public"."sd_air_warning_device_map" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "air_warning_id" int8,
  "air_control_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_air_warning_device_map
-- ----------------------------

-- ----------------------------
-- Table structure for sd_alarm_process_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_alarm_process_log";
CREATE TABLE "public"."sd_alarm_process_log" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8,
  "type_id" int8,
  "event" text COLLATE "pg_catalog"."default",
  "alarm_type" text COLLATE "pg_catalog"."default",
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "email_alarm" int8,
  "line_alarm" int8,
  "telegram_alarm" int8,
  "sms_alarm" int8,
  "nonc_alarm" int8,
  "status" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "data_alarm" text COLLATE "pg_catalog"."default",
  "alarm_status" text COLLATE "pg_catalog"."default",
  "subject" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_alarm_process_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_alarm_process_log_email
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_alarm_process_log_email";
CREATE TABLE "public"."sd_alarm_process_log_email" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8,
  "type_id" int8,
  "event" text COLLATE "pg_catalog"."default",
  "alarm_type" text COLLATE "pg_catalog"."default",
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "email_alarm" int8,
  "line_alarm" int8,
  "telegram_alarm" int8,
  "sms_alarm" int8,
  "nonc_alarm" int8,
  "status" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "data_alarm" text COLLATE "pg_catalog"."default",
  "alarm_status" text COLLATE "pg_catalog"."default",
  "subject" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_alarm_process_log_email
-- ----------------------------

-- ----------------------------
-- Table structure for sd_alarm_process_log_line
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_alarm_process_log_line";
CREATE TABLE "public"."sd_alarm_process_log_line" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8,
  "type_id" int8,
  "event" text COLLATE "pg_catalog"."default",
  "alarm_type" text COLLATE "pg_catalog"."default",
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "email_alarm" int8,
  "line_alarm" int8,
  "telegram_alarm" int8,
  "sms_alarm" int8,
  "nonc_alarm" int8,
  "status" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "data_alarm" text COLLATE "pg_catalog"."default",
  "alarm_status" text COLLATE "pg_catalog"."default",
  "subject" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_alarm_process_log_line
-- ----------------------------

-- ----------------------------
-- Table structure for sd_alarm_process_log_mqtt
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_alarm_process_log_mqtt";
CREATE TABLE "public"."sd_alarm_process_log_mqtt" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8,
  "type_id" int8,
  "event" text COLLATE "pg_catalog"."default",
  "alarm_type" text COLLATE "pg_catalog"."default",
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "email_alarm" int8,
  "line_alarm" int8,
  "telegram_alarm" int8,
  "sms_alarm" int8,
  "nonc_alarm" int8,
  "status" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "data_alarm" text COLLATE "pg_catalog"."default",
  "alarm_status" text COLLATE "pg_catalog"."default",
  "subject" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_alarm_process_log_mqtt
-- ----------------------------

-- ----------------------------
-- Table structure for sd_alarm_process_log_sms
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_alarm_process_log_sms";
CREATE TABLE "public"."sd_alarm_process_log_sms" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8,
  "type_id" int8,
  "event" text COLLATE "pg_catalog"."default",
  "alarm_type" text COLLATE "pg_catalog"."default",
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "email_alarm" int8,
  "line_alarm" int8,
  "telegram_alarm" int8,
  "sms_alarm" int8,
  "nonc_alarm" int8,
  "status" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "data_alarm" text COLLATE "pg_catalog"."default",
  "alarm_status" text COLLATE "pg_catalog"."default",
  "subject" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_alarm_process_log_sms
-- ----------------------------

-- ----------------------------
-- Table structure for sd_alarm_process_log_telegram
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_alarm_process_log_telegram";
CREATE TABLE "public"."sd_alarm_process_log_telegram" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8,
  "type_id" int8,
  "event" text COLLATE "pg_catalog"."default",
  "alarm_type" text COLLATE "pg_catalog"."default",
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "email_alarm" int8,
  "line_alarm" int8,
  "telegram_alarm" int8,
  "sms_alarm" int8,
  "nonc_alarm" int8,
  "status" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "data_alarm" text COLLATE "pg_catalog"."default",
  "alarm_status" text COLLATE "pg_catalog"."default",
  "subject" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_alarm_process_log_telegram
-- ----------------------------

-- ----------------------------
-- Table structure for sd_alarm_process_log_temp
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_alarm_process_log_temp";
CREATE TABLE "public"."sd_alarm_process_log_temp" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8,
  "type_id" int8,
  "event" text COLLATE "pg_catalog"."default",
  "alarm_type" text COLLATE "pg_catalog"."default",
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "email_alarm" int8,
  "line_alarm" int8,
  "telegram_alarm" int8,
  "sms_alarm" int8,
  "nonc_alarm" int8,
  "status" text COLLATE "pg_catalog"."default",
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "data_alarm" text COLLATE "pg_catalog"."default",
  "alarm_status" text COLLATE "pg_catalog"."default",
  "subject" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_alarm_process_log_temp
-- ----------------------------

-- ----------------------------
-- Table structure for sd_api_key
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_api_key";
CREATE TABLE "public"."sd_api_key" (
  "id" int8 NOT NULL DEFAULT nextval('sd_api_key_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "api_key" text COLLATE "pg_catalog"."default" NOT NULL,
  "api_secret" text COLLATE "pg_catalog"."default" NOT NULL,
  "user_id" text COLLATE "pg_catalog"."default",
  "permissions" jsonb,
  "expires_at" timestamptz(6),
  "last_used_at" timestamptz(6),
  "usage_count" int8 NOT NULL DEFAULT 0,
  "is_active" bool NOT NULL DEFAULT true,
  "ip_whitelist" jsonb,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_api_key
-- ----------------------------

-- ----------------------------
-- Table structure for sd_audit_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_audit_log";
CREATE TABLE "public"."sd_audit_log" (
  "audit_id" int8 NOT NULL DEFAULT nextval('sd_audit_log_audit_id_seq'::regclass),
  "user_id" text COLLATE "pg_catalog"."default",
  "user_name" text COLLATE "pg_catalog"."default",
  "action" text COLLATE "pg_catalog"."default" NOT NULL,
  "entity_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "entity_id" int8 NOT NULL,
  "before" jsonb,
  "after" jsonb,
  "changes" jsonb,
  "ip_address" text COLLATE "pg_catalog"."default",
  "user_agent" text COLLATE "pg_catalog"."default",
  "action_time" timestamptz(6) NOT NULL DEFAULT now(),
  "description" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_audit_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_channel_template
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_channel_template";
CREATE TABLE "public"."sd_channel_template" (
  "id" int8 NOT NULL DEFAULT nextval('sd_channel_template_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "channel_id" int8 NOT NULL,
  "notification_type_id" int8 NOT NULL,
  "template" text COLLATE "pg_catalog"."default" NOT NULL,
  "variables" jsonb,
  "is_active" bool NOT NULL DEFAULT true,
  "is_default" bool NOT NULL DEFAULT false,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_channel_template
-- ----------------------------

-- ----------------------------
-- Table structure for sd_dashboard_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_dashboard_config";
CREATE TABLE "public"."sd_dashboard_config" (
  "id" int8 NOT NULL DEFAULT nextval('sd_dashboard_config_id_seq'::regclass),
  "location_id" int8 NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "config_data" json NOT NULL,
  "status" int8 NOT NULL DEFAULT 1,
  "created_date" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_date" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_dashboard_config
-- ----------------------------

-- ----------------------------
-- Table structure for sd_device_category
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_device_category";
CREATE TABLE "public"."sd_device_category" (
  "id" int8 NOT NULL DEFAULT nextval('sd_device_category_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "icon" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_device_category
-- ----------------------------

-- ----------------------------
-- Table structure for sd_device_group
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_device_group";
CREATE TABLE "public"."sd_device_group" (
  "id" int8 NOT NULL DEFAULT nextval('sd_device_group_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "group_type" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'custom'::text,
  "is_active" bool NOT NULL DEFAULT true,
  "config" jsonb,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_device_group
-- ----------------------------

-- ----------------------------
-- Table structure for sd_device_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_device_log";
CREATE TABLE "public"."sd_device_log" (
  "id" int8 NOT NULL DEFAULT nextval('sd_device_log_id_seq'::regclass),
  "type_id" int8 NOT NULL,
  "sensor_id" int8 NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default" NOT NULL,
  "status" int8,
  "lang" text COLLATE "pg_catalog"."default",
  "create" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_device_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_device_member
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_device_member";
CREATE TABLE "public"."sd_device_member" (
  "id" int8 NOT NULL DEFAULT nextval('sd_device_member_id_seq'::regclass),
  "Device_id" int8 NOT NULL,
  "group_id" int8 NOT NULL,
  "role" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'member'::text,
  "priority" int8 NOT NULL DEFAULT 1,
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_device_member
-- ----------------------------

-- ----------------------------
-- Table structure for sd_device_notification_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_device_notification_config";
CREATE TABLE "public"."sd_device_notification_config" (
  "id" int8 NOT NULL DEFAULT nextval('sd_device_notification_config_id_seq'::regclass),
  "device_id" int8 NOT NULL,
  "notification_channel_id" int8 NOT NULL,
  "notification_type_id" int8 NOT NULL,
  "config" jsonb,
  "is_active" bool NOT NULL DEFAULT true,
  "retry_count" int8 NOT NULL DEFAULT 3,
  "retry_delay_minutes" int8 NOT NULL DEFAULT 5,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_device_notification_config
-- ----------------------------

-- ----------------------------
-- Table structure for sd_device_schedule
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_device_schedule";
CREATE TABLE "public"."sd_device_schedule" (
  "id" int8 NOT NULL DEFAULT nextval('sd_device_schedule_id_seq'::regclass),
  "device_id" int8 NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "schedule_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "schedule_config" jsonb NOT NULL,
  "action" jsonb NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "last_run_at" timestamptz(6),
  "next_run_at" timestamptz(6),
  "run_count" int8 NOT NULL DEFAULT 0,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_device_schedule
-- ----------------------------

-- ----------------------------
-- Table structure for sd_device_status_history
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_device_status_history";
CREATE TABLE "public"."sd_device_status_history" (
  "id" int8 NOT NULL DEFAULT nextval('sd_device_status_history_id_seq'::regclass),
  "device_id" int8 NOT NULL,
  "status" text COLLATE "pg_catalog"."default",
  "value" numeric(10,2),
  "notification_type_id" int8,
  "duration_minutes" int8,
  "previous_status" text COLLATE "pg_catalog"."default",
  "previous_value" numeric(10,2),
  "change_reason" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_device_status_history
-- ----------------------------

-- ----------------------------
-- Table structure for sd_group_notification_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_group_notification_config";
CREATE TABLE "public"."sd_group_notification_config" (
  "id" int8 NOT NULL DEFAULT nextval('sd_group_notification_config_id_seq'::regclass),
  "group_id" int8 NOT NULL,
  "notification_channel_id" int8 NOT NULL,
  "notification_type_id" int8 NOT NULL,
  "config" jsonb,
  "is_active" bool NOT NULL DEFAULT true,
  "escalation_level" int8 NOT NULL DEFAULT 1,
  "escalation_delay_minutes" int8 NOT NULL DEFAULT 30,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_group_notification_config
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_alarm_device
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_alarm_device";
CREATE TABLE "public"."sd_iot_alarm_device" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_iot_alarm_device
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_alarm_device_event
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_alarm_device_event";
CREATE TABLE "public"."sd_iot_alarm_device_event" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "alarm_action_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_iot_alarm_device_event
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_api
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_api";
CREATE TABLE "public"."sd_iot_api" (
  "api_id" int8 NOT NULL DEFAULT nextval('sd_iot_api_api_id_seq'::regclass),
  "api_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "host" int8,
  "port" text COLLATE "pg_catalog"."default" NOT NULL,
  "token_value" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_api
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_device
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_device";
CREATE TABLE "public"."sd_iot_device" (
  "device_id" int8 NOT NULL DEFAULT nextval('sd_iot_device_device_id_seq'::regclass),
  "setting_id" int8,
  "type_id" int8,
  "location_id" int8,
  "device_name" text COLLATE "pg_catalog"."default",
  "sn" text COLLATE "pg_catalog"."default",
  "hardware_id" int8,
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "time_life" int8 DEFAULT 1,
  "period" text COLLATE "pg_catalog"."default",
  "work_status" int8 DEFAULT 1,
  "model" text COLLATE "pg_catalog"."default",
  "vendor" text COLLATE "pg_catalog"."default",
  "comparevalue" text COLLATE "pg_catalog"."default",
  "unit" text COLLATE "pg_catalog"."default",
  "mqtt_id" int8,
  "oid" text COLLATE "pg_catalog"."default",
  "action_id" int8,
  "status_alert_id" int8,
  "mqtt_data_value" text COLLATE "pg_catalog"."default",
  "mqtt_data_control" text COLLATE "pg_catalog"."default",
  "measurement" text COLLATE "pg_catalog"."default",
  "mqtt_control_on" text COLLATE "pg_catalog"."default" DEFAULT '1'::text,
  "mqtt_control_off" text COLLATE "pg_catalog"."default" DEFAULT '0'::text,
  "org" text COLLATE "pg_catalog"."default" NOT NULL,
  "bucket" text COLLATE "pg_catalog"."default" NOT NULL,
  "status" int8,
  "mqtt_device_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "mqtt_status_over_name" text COLLATE "pg_catalog"."default",
  "mqtt_status_data_name" text COLLATE "pg_catalog"."default",
  "mqtt_act_relay_name" text COLLATE "pg_catalog"."default",
  "mqtt_control_relay_name" text COLLATE "pg_catalog"."default",
  "mqtt_config" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "max" text COLLATE "pg_catalog"."default",
  "min" text COLLATE "pg_catalog"."default",
  "layout" int8 DEFAULT 1,
  "alert_set" int8 DEFAULT 1,
  "icon_normal" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '<svg...>'::text,
  "icon_warning" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '<svg...>'::text,
  "icon_alert" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '<svg...>'::text,
  "icon" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '<svg...>'::text,
  "color_normal" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '#22C55E'::text,
  "color_warning" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '#F59E0B'::text,
  "color_alarm" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '#EF4444'::text,
  "code" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'normal'::text,
  "menu" int8 DEFAULT 1,
  "icon_on" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '<svg...>'::text,
  "icon_off" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT '<svg...>'::text,
  "calibration_add" text COLLATE "pg_catalog"."default" DEFAULT '0'::text,
  "calibration_subtract" text COLLATE "pg_catalog"."default" DEFAULT '0'::text,
  "calibration_type" int8 DEFAULT 3
)
;

-- ----------------------------
-- Records of sd_iot_device
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_device_action
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_device_action";
CREATE TABLE "public"."sd_iot_device_action" (
  "device_action_user_id" int8 NOT NULL DEFAULT nextval('sd_iot_device_action_device_action_user_id_seq'::regclass),
  "alarm_action_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_iot_device_action
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_device_action_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_device_action_log";
CREATE TABLE "public"."sd_iot_device_action_log" (
  "log_id" int8 NOT NULL DEFAULT nextval('sd_iot_device_action_log_log_id_seq'::regclass),
  "alarm_action_id" int8,
  "device_id" int8,
  "uid" text COLLATE "pg_catalog"."default",
  "status" int8,
  "createddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_iot_device_action_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_device_action_user
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_device_action_user";
CREATE TABLE "public"."sd_iot_device_action_user" (
  "device_action_user_id" int8 NOT NULL DEFAULT nextval('sd_iot_device_action_user_device_action_user_id_seq'::regclass),
  "alarm_action_id" int8,
  "uid" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of sd_iot_device_action_user
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_device_alarm_action
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_device_alarm_action";
CREATE TABLE "public"."sd_iot_device_alarm_action" (
  "alarm_action_id" int8 NOT NULL DEFAULT nextval('sd_iot_device_alarm_action_alarm_action_id_seq'::regclass),
  "action_name" text COLLATE "pg_catalog"."default",
  "status_warning" text COLLATE "pg_catalog"."default",
  "recovery_warning" text COLLATE "pg_catalog"."default",
  "status_alert" text COLLATE "pg_catalog"."default",
  "recovery_alert" text COLLATE "pg_catalog"."default",
  "email_alarm" int8,
  "line_alarm" int8,
  "telegram_alarm" int8,
  "sms_alarm" int8,
  "nonc_alarm" int8,
  "time_life" int8,
  "event" int8,
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_device_alarm_action
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_device_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_device_type";
CREATE TABLE "public"."sd_iot_device_type" (
  "type_id" int8 NOT NULL DEFAULT nextval('sd_iot_device_type_type_id_seq'::regclass),
  "type_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_device_type
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_email
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_email";
CREATE TABLE "public"."sd_iot_email" (
  "email_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "email_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "host" text COLLATE "pg_catalog"."default" NOT NULL,
  "port" int8,
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_email
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_group
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_group";
CREATE TABLE "public"."sd_iot_group" (
  "group_id" int8 NOT NULL DEFAULT nextval('sd_iot_group_group_id_seq'::regclass),
  "group_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_group
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_host
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_host";
CREATE TABLE "public"."sd_iot_host" (
  "host_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "host_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "port" text COLLATE "pg_catalog"."default" NOT NULL,
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8,
  "idhost" int8
)
;

-- ----------------------------
-- Records of sd_iot_host
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_influxdb
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_influxdb";
CREATE TABLE "public"."sd_iot_influxdb" (
  "influxdb_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "influxdb_name" text COLLATE "pg_catalog"."default",
  "host" text COLLATE "pg_catalog"."default",
  "port" text COLLATE "pg_catalog"."default" NOT NULL,
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "token_value" text COLLATE "pg_catalog"."default",
  "buckets" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_influxdb
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_line
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_line";
CREATE TABLE "public"."sd_iot_line" (
  "line_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "line_name" text COLLATE "pg_catalog"."default",
  "client_id" text COLLATE "pg_catalog"."default",
  "client_secret" text COLLATE "pg_catalog"."default",
  "secret_key" text COLLATE "pg_catalog"."default",
  "redirect_uri" text COLLATE "pg_catalog"."default",
  "grant_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "code" text COLLATE "pg_catalog"."default" NOT NULL,
  "accesstoken" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_line
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_location
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_location";
CREATE TABLE "public"."sd_iot_location" (
  "location_id" int8 NOT NULL DEFAULT nextval('sd_iot_location_location_id_seq'::regclass),
  "location_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "ipaddress" text COLLATE "pg_catalog"."default" NOT NULL,
  "location_detail" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8,
  "configdata" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of sd_iot_location
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_mqtt
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_mqtt";
CREATE TABLE "public"."sd_iot_mqtt" (
  "mqtt_id" int8 NOT NULL DEFAULT nextval('sd_iot_mqtt_mqtt_id_seq'::regclass),
  "mqtt_type_id" int8,
  "sort" int8 NOT NULL DEFAULT 1,
  "mqtt_name" text COLLATE "pg_catalog"."default",
  "host" text COLLATE "pg_catalog"."default",
  "port" int8,
  "username" text COLLATE "pg_catalog"."default",
  "password" text COLLATE "pg_catalog"."default",
  "secret" text COLLATE "pg_catalog"."default",
  "expire_in" text COLLATE "pg_catalog"."default",
  "token_value" text COLLATE "pg_catalog"."default",
  "org" text COLLATE "pg_catalog"."default",
  "bucket" text COLLATE "pg_catalog"."default",
  "envavorment" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8 NOT NULL DEFAULT 1,
  "location_id" int8 DEFAULT 1,
  "latitude" text COLLATE "pg_catalog"."default",
  "longitude" text COLLATE "pg_catalog"."default",
  "mqtt_main_id" int8 NOT NULL DEFAULT 1,
  "configuration" text COLLATE "pg_catalog"."default" DEFAULT '{"0":"temperature1","1":"humidity1"}'::text,
  "zoom" int8 NOT NULL DEFAULT 6
)
;

-- ----------------------------
-- Records of sd_iot_mqtt
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_nodered
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_nodered";
CREATE TABLE "public"."sd_iot_nodered" (
  "nodered_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "nodered_name" text COLLATE "pg_catalog"."default",
  "host" text COLLATE "pg_catalog"."default" NOT NULL,
  "port" text COLLATE "pg_catalog"."default" NOT NULL,
  "routing" text COLLATE "pg_catalog"."default",
  "client_id" text COLLATE "pg_catalog"."default",
  "grant_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "scope" text COLLATE "pg_catalog"."default" NOT NULL,
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_nodered
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_schedule
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_schedule";
CREATE TABLE "public"."sd_iot_schedule" (
  "schedule_id" int8 NOT NULL DEFAULT nextval('sd_iot_schedule_schedule_id_seq'::regclass),
  "schedule_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "device_id" int8,
  "start" text COLLATE "pg_catalog"."default" NOT NULL,
  "event" int8,
  "sunday" int8,
  "monday" int8,
  "tuesday" int8,
  "wednesday" int8,
  "thursday" int8,
  "friday" int8,
  "saturday" int8,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_schedule
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_schedule_device
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_schedule_device";
CREATE TABLE "public"."sd_iot_schedule_device" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "schedule_id" int8,
  "device_id" int8
)
;

-- ----------------------------
-- Records of sd_iot_schedule_device
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_sensor
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_sensor";
CREATE TABLE "public"."sd_iot_sensor" (
  "sensor_id" int8 NOT NULL DEFAULT nextval('sd_iot_sensor_sensor_id_seq'::regclass),
  "setting_id" int8,
  "setting_type_id" int8,
  "sensor_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "sn" text COLLATE "pg_catalog"."default" NOT NULL,
  "max" text COLLATE "pg_catalog"."default" NOT NULL,
  "min" text COLLATE "pg_catalog"."default" NOT NULL,
  "hardware_id" int8,
  "status_high" text COLLATE "pg_catalog"."default" NOT NULL,
  "status_warning" text COLLATE "pg_catalog"."default" NOT NULL,
  "status_alert" text COLLATE "pg_catalog"."default" NOT NULL,
  "model" text COLLATE "pg_catalog"."default" NOT NULL,
  "vendor" text COLLATE "pg_catalog"."default" NOT NULL,
  "comparevalue" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8,
  "unit" text COLLATE "pg_catalog"."default" NOT NULL,
  "mqtt_id" int8,
  "oid" text COLLATE "pg_catalog"."default" NOT NULL,
  "action_id" int8,
  "status_alert_id" int8,
  "mqtt_data_value" text COLLATE "pg_catalog"."default" NOT NULL,
  "mqtt_data_control" text COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Records of sd_iot_sensor
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_setting
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_setting";
CREATE TABLE "public"."sd_iot_setting" (
  "setting_id" int8 NOT NULL DEFAULT nextval('sd_iot_setting_setting_id_seq'::regclass),
  "location_id" int8,
  "setting_type_id" int8,
  "setting_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "sn" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_setting
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_sms
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_sms";
CREATE TABLE "public"."sd_iot_sms" (
  "sms_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "sms_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "host" text COLLATE "pg_catalog"."default" NOT NULL,
  "port" int8,
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "apikey" text COLLATE "pg_catalog"."default" NOT NULL,
  "originator" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_sms
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_telegram
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_telegram";
CREATE TABLE "public"."sd_iot_telegram" (
  "telegram_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "telegram_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "port" text COLLATE "pg_catalog"."default" NOT NULL,
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_telegram
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_token
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_token";
CREATE TABLE "public"."sd_iot_token" (
  "token_id" int8 NOT NULL DEFAULT nextval('sd_iot_token_token_id_seq'::regclass),
  "token_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "host" int8,
  "port" text COLLATE "pg_catalog"."default" NOT NULL,
  "token_value" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_token
-- ----------------------------

-- ----------------------------
-- Table structure for sd_iot_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_iot_type";
CREATE TABLE "public"."sd_iot_type" (
  "type_id" int8 NOT NULL DEFAULT nextval('sd_iot_type_type_id_seq'::regclass),
  "type_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "group_id" int8,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8
)
;

-- ----------------------------
-- Records of sd_iot_type
-- ----------------------------

-- ----------------------------
-- Table structure for sd_module_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_module_log";
CREATE TABLE "public"."sd_module_log" (
  "id" int8 NOT NULL DEFAULT nextval('sd_module_log_id_seq'::regclass),
  "module_name" text COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Records of sd_module_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_mqtt_host
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_mqtt_host";
CREATE TABLE "public"."sd_mqtt_host" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "hostname" text COLLATE "pg_catalog"."default" NOT NULL,
  "host" text COLLATE "pg_catalog"."default" NOT NULL,
  "port" text COLLATE "pg_catalog"."default" NOT NULL,
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int8,
  "idhost" int8
)
;

-- ----------------------------
-- Records of sd_mqtt_host
-- ----------------------------

-- ----------------------------
-- Table structure for sd_mqtt_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_mqtt_log";
CREATE TABLE "public"."sd_mqtt_log" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "name" text COLLATE "pg_catalog"."default",
  "statusmqtt" text COLLATE "pg_catalog"."default",
  "msg" text COLLATE "pg_catalog"."default",
  "type_id" int8,
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now(),
  "Device_id" int8,
  "Device_name" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of sd_mqtt_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_notification_channel
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_notification_channel";
CREATE TABLE "public"."sd_notification_channel" (
  "id" int8 NOT NULL DEFAULT nextval('sd_notification_channel_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "icon" text COLLATE "pg_catalog"."default",
  "handler_class" text COLLATE "pg_catalog"."default",
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_notification_channel
-- ----------------------------

-- ----------------------------
-- Table structure for sd_notification_condition
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_notification_condition";
CREATE TABLE "public"."sd_notification_condition" (
  "id" int8 NOT NULL DEFAULT nextval('sd_notification_condition_id_seq'::regclass),
  "device_id" int8 NOT NULL,
  "notification_type_id" int8 NOT NULL,
  "condition_operator" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'between'::text,
  "priority" int8 NOT NULL DEFAULT 1,
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "minValue" numeric(10,2),
  "maxValue" numeric(10,2)
)
;

-- ----------------------------
-- Records of sd_notification_condition
-- ----------------------------

-- ----------------------------
-- Table structure for sd_notification_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_notification_log";
CREATE TABLE "public"."sd_notification_log" (
  "id" int8 NOT NULL DEFAULT nextval('sd_notification_log_id_seq'::regclass),
  "device_id" int8,
  "notification_type_id" int8,
  "notification_channel_id" int8,
  "message" text COLLATE "pg_catalog"."default" NOT NULL,
  "response_data" jsonb,
  "sent_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "template_id" int8,
  "delivered_at" timestamptz(6),
  "read_at" timestamptz(6),
  "retry_count" int8 NOT NULL DEFAULT 0,
  "error_message" text COLLATE "pg_catalog"."default",
  "message_id" text COLLATE "pg_catalog"."default",
  "recipient" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'pending'::text
)
;

-- ----------------------------
-- Records of sd_notification_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_notification_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_notification_type";
CREATE TABLE "public"."sd_notification_type" (
  "id" int8 NOT NULL DEFAULT nextval('sd_notification_type_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "cooldown_minutes" int8 NOT NULL DEFAULT 10,
  "is_active" bool NOT NULL DEFAULT true,
  "icon" text COLLATE "pg_catalog"."default",
  "color" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_notification_type
-- ----------------------------

-- ----------------------------
-- Table structure for sd_report_data
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_report_data";
CREATE TABLE "public"."sd_report_data" (
  "id" int8 NOT NULL DEFAULT nextval('sd_report_data_id_seq'::regclass),
  "device_id" int8 NOT NULL,
  "report_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "data" jsonb NOT NULL,
  "period_start" timestamptz(6) NOT NULL,
  "period_end" timestamptz(6) NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "template_id" int8,
  "generated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "file_path" text COLLATE "pg_catalog"."default",
  "file_format" text COLLATE "pg_catalog"."default",
  "is_exported" bool NOT NULL DEFAULT false,
  "exported_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of sd_report_data
-- ----------------------------

-- ----------------------------
-- Table structure for sd_schedule_process_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_schedule_process_log";
CREATE TABLE "public"."sd_schedule_process_log" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "schedule_id" int8,
  "device_id" int8,
  "schedule_event_start" text COLLATE "pg_catalog"."default" NOT NULL,
  "day" text COLLATE "pg_catalog"."default" NOT NULL,
  "doday" text COLLATE "pg_catalog"."default" NOT NULL,
  "dotime" text COLLATE "pg_catalog"."default" NOT NULL,
  "schedule_event" text COLLATE "pg_catalog"."default" NOT NULL,
  "device_status" text COLLATE "pg_catalog"."default" NOT NULL,
  "status" text COLLATE "pg_catalog"."default" NOT NULL,
  "date" text COLLATE "pg_catalog"."default" NOT NULL,
  "time" text COLLATE "pg_catalog"."default" NOT NULL,
  "createddate" timestamptz(6) NOT NULL DEFAULT now(),
  "updateddate" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_schedule_process_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_sensor_data
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_sensor_data";
CREATE TABLE "public"."sd_sensor_data" (
  "id" int8 NOT NULL DEFAULT nextval('sd_sensor_data_id_seq'::regclass),
  "device_id" int8 NOT NULL,
  "value" numeric(10,2) NOT NULL,
  "raw_data" jsonb,
  "notification_type_id" int8,
  "timestamp" timestamptz(6) NOT NULL DEFAULT now(),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "battery_level" numeric(5,2),
  "signal_strength" int8
)
;

-- ----------------------------
-- Records of sd_sensor_data
-- ----------------------------

-- ----------------------------
-- Table structure for sd_system_setting
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_system_setting";
CREATE TABLE "public"."sd_system_setting" (
  "id" int8 NOT NULL DEFAULT nextval('sd_system_setting_id_seq'::regclass),
  "key" text COLLATE "pg_catalog"."default" NOT NULL,
  "value" jsonb NOT NULL,
  "category" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "is_public" bool NOT NULL DEFAULT false,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_system_setting
-- ----------------------------

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
INSERT INTO "public"."sd_user" VALUES ('e32af9b1-a773-4757-9a03-09c3c01dc0ce', '2026-04-06 16:23:23.047937+07', '2026-04-06 16:23:23.047937+07', NULL, 2, 'kongnakornna@gmail.com', 'kongnakornna@gmail.com', '$2a$10$mmW7K300LuUPLIcgf0H4puiZn/1GDz3ZMnVegzmM5S24ecVIUSzca', NULL, 'Kongnakorn', 'Jantakun', 'Kongnakorn Jantakun', NULL, NULL, '2026-04-06 16:23:23.047937+07', 1, NULL, 1, NULL, 0, NULL, NULL, '0', NULL, 0, 0, 0, NULL, NULL, NULL, NULL, 0, 0, 0, 0, '0812345678', '021234567', 'kongnakorn_line', '1', 'loc_001', 'f', '5564655fd82272e5be5a259b029e9825', '', NULL, 'f');

-- ----------------------------
-- Table structure for sd_user_access_menu
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_access_menu";
CREATE TABLE "public"."sd_user_access_menu" (
  "user_access_id" int8 NOT NULL DEFAULT nextval('sd_user_access_menu_user_access_id_seq'::regclass),
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
-- Table structure for sd_user_file
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_file";
CREATE TABLE "public"."sd_user_file" (
  "id" int8 NOT NULL DEFAULT nextval('sd_user_file_id_seq'::regclass),
  "file_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "file_type" text COLLATE "pg_catalog"."default" NOT NULL,
  "file_path" text COLLATE "pg_catalog"."default" NOT NULL,
  "file_type_id" int8 NOT NULL,
  "uid" text COLLATE "pg_catalog"."default",
  "file_date" timestamptz(6) NOT NULL DEFAULT now(),
  "status" int2 NOT NULL
)
;

-- ----------------------------
-- Records of sd_user_file
-- ----------------------------

-- ----------------------------
-- Table structure for sd_user_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_log";
CREATE TABLE "public"."sd_user_log" (
  "id" int8 NOT NULL DEFAULT nextval('sd_user_log_id_seq'::regclass),
  "log_type_id" int8 NOT NULL,
  "uid" uuid NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "detail" text COLLATE "pg_catalog"."default" NOT NULL,
  "select_status" int8,
  "insert_status" int8,
  "update_status" int8,
  "delete_status" int8,
  "status" int8,
  "create" timestamptz(6) NOT NULL DEFAULT now(),
  "update" timestamptz(6) NOT NULL DEFAULT now(),
  "lang" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of sd_user_log
-- ----------------------------

-- ----------------------------
-- Table structure for sd_user_log_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_log_type";
CREATE TABLE "public"."sd_user_log_type" (
  "log_type_id" int8 NOT NULL DEFAULT nextval('sd_user_log_type_log_type_id_seq'::regclass),
  "type_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "type_detail" text COLLATE "pg_catalog"."default" NOT NULL,
  "status" int8,
  "create" timestamptz(6) NOT NULL DEFAULT now(),
  "update" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of sd_user_log_type
-- ----------------------------

-- ----------------------------
-- Table structure for sd_user_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_role";
CREATE TABLE "public"."sd_user_role" (
  "id" int8 NOT NULL DEFAULT nextval('sd_user_role_id_seq'::regclass),
  "role_id" int8 NOT NULL,
  "title" varchar(50) COLLATE "pg_catalog"."default",
  "createddate" timestamptz(6),
  "updateddate" timestamptz(6),
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

-- ----------------------------
-- Table structure for sd_user_roles_access
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_roles_access";
CREATE TABLE "public"."sd_user_roles_access" (
  "create" timestamptz(6) NOT NULL,
  "update" timestamptz(6) NOT NULL,
  "role_id" int8 NOT NULL,
  "role_type_id" int8 NOT NULL
)
;

-- ----------------------------
-- Records of sd_user_roles_access
-- ----------------------------

-- ----------------------------
-- Table structure for sd_user_roles_permision
-- ----------------------------
DROP TABLE IF EXISTS "public"."sd_user_roles_permision";
CREATE TABLE "public"."sd_user_roles_permision" (
  "role_type_id" int8 NOT NULL DEFAULT nextval('sd_user_roles_permision_role_type_id_seq'::regclass),
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "detail" text COLLATE "pg_catalog"."default",
  "created" timestamptz(6) NOT NULL,
  "updated" timestamptz(6),
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

-- ----------------------------
-- Table structure for tnb
-- ----------------------------
DROP TABLE IF EXISTS "public"."tnb";
CREATE TABLE "public"."tnb" (
  "id" int8 NOT NULL DEFAULT nextval('tnb_id_seq'::regclass),
  "location_id" int8 NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "config_data" json NOT NULL,
  "status" int8 NOT NULL DEFAULT 1,
  "created_date" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_date" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of tnb
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS "public"."user";
CREATE TABLE "public"."user" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "email" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "password" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "is_super_user" bool NOT NULL DEFAULT false,
  "verified" bool NOT NULL DEFAULT false,
  "verification_code" varchar(32) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "password_reset_token" varchar(32) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "password_reset_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO "public"."user" VALUES ('02404e80-af5f-43cd-aa05-597b2b658b1c', '2026-03-31 20:09:10.535027+07', '2026-03-31 20:09:10.545345+07', NULL, 'kongnakornna', 'kongnakornna@gmail.com', '$2a$10$O.GzDjRPc7EX9j/xuqM6IeK3tEboG1Cu986D6vKYpW9ffSa4CatO6', 't', 'f', 'f', '79758af58a5ed221836325f94e60aebf', NULL, NULL);
INSERT INTO "public"."user" VALUES ('6b049b5e-88c8-4b1f-9a3b-8f8d30e68fda', '2026-03-31 19:07:14.339724+07', '2026-03-31 20:11:19.925615+07', NULL, 'root', 'root@gmail.com', '$2a$10$g33en5/JjLFJxNdlILZE0uZ/fJe.FVH8J4qxhihJSRmh0oH/ZxYyu', 't', 't', 't', NULL, NULL, NULL);

-- ----------------------------
-- Function structure for armor
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."armor"(bytea);
CREATE FUNCTION "public"."armor"(bytea)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_armor'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for armor
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."armor"(bytea, _text, _text);
CREATE FUNCTION "public"."armor"(bytea, _text, _text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_armor'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for crypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."crypt"(text, text);
CREATE FUNCTION "public"."crypt"(text, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_crypt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for dearmor
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."dearmor"(text);
CREATE FUNCTION "public"."dearmor"(text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_dearmor'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."decrypt"(bytea, bytea, text);
CREATE FUNCTION "public"."decrypt"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_decrypt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for decrypt_iv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."decrypt_iv"(bytea, bytea, bytea, text);
CREATE FUNCTION "public"."decrypt_iv"(bytea, bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_decrypt_iv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for digest
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."digest"(bytea, text);
CREATE FUNCTION "public"."digest"(bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_digest'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for digest
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."digest"(text, text);
CREATE FUNCTION "public"."digest"(text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_digest'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."encrypt"(bytea, bytea, text);
CREATE FUNCTION "public"."encrypt"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_encrypt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for encrypt_iv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."encrypt_iv"(bytea, bytea, bytea, text);
CREATE FUNCTION "public"."encrypt_iv"(bytea, bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_encrypt_iv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gen_random_bytes
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gen_random_bytes"(int4);
CREATE FUNCTION "public"."gen_random_bytes"(int4)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_random_bytes'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gen_random_uuid
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gen_random_uuid"();
CREATE FUNCTION "public"."gen_random_uuid"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/pgcrypto', 'pg_random_uuid'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for gen_salt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gen_salt"(text, int4);
CREATE FUNCTION "public"."gen_salt"(text, int4)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_gen_salt_rounds'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gen_salt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gen_salt"(text);
CREATE FUNCTION "public"."gen_salt"(text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_gen_salt'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gseg_consistent
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gseg_consistent"(internal, "public"."seg", int2, oid, internal);
CREATE FUNCTION "public"."gseg_consistent"(internal, "public"."seg", int2, oid, internal)
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'gseg_consistent'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gseg_penalty
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gseg_penalty"(internal, internal, internal);
CREATE FUNCTION "public"."gseg_penalty"(internal, internal, internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/seg', 'gseg_penalty'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gseg_picksplit
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gseg_picksplit"(internal, internal);
CREATE FUNCTION "public"."gseg_picksplit"(internal, internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/seg', 'gseg_picksplit'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gseg_same
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gseg_same"("public"."seg", "public"."seg", internal);
CREATE FUNCTION "public"."gseg_same"("public"."seg", "public"."seg", internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/seg', 'gseg_same'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gseg_union
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gseg_union"(internal, internal);
CREATE FUNCTION "public"."gseg_union"(internal, internal)
  RETURNS "public"."seg" AS '$libdir/seg', 'gseg_union'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for hmac
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hmac"(bytea, bytea, text);
CREATE FUNCTION "public"."hmac"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_hmac'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for hmac
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hmac"(text, text, text);
CREATE FUNCTION "public"."hmac"(text, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_hmac'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_armor_headers
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_armor_headers"(text, OUT "key" text, OUT "value" text);
CREATE FUNCTION "public"."pgp_armor_headers"(IN text, OUT "key" text, OUT "value" text)
  RETURNS SETOF "pg_catalog"."record" AS '$libdir/pgcrypto', 'pgp_armor_headers'
  LANGUAGE c IMMUTABLE STRICT
  COST 1
  ROWS 1000;

-- ----------------------------
-- Function structure for pgp_key_id
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_key_id"(bytea);
CREATE FUNCTION "public"."pgp_key_id"(bytea)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_key_id_w'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt"(bytea, bytea);
CREATE FUNCTION "public"."pgp_pub_decrypt"(bytea, bytea)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt"(bytea, bytea, text, text);
CREATE FUNCTION "public"."pgp_pub_decrypt"(bytea, bytea, text, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt"(bytea, bytea, text);
CREATE FUNCTION "public"."pgp_pub_decrypt"(bytea, bytea, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt_bytea"(bytea, bytea);
CREATE FUNCTION "public"."pgp_pub_decrypt_bytea"(bytea, bytea)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt_bytea"(bytea, bytea, text);
CREATE FUNCTION "public"."pgp_pub_decrypt_bytea"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt_bytea"(bytea, bytea, text, text);
CREATE FUNCTION "public"."pgp_pub_decrypt_bytea"(bytea, bytea, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_encrypt"(text, bytea);
CREATE FUNCTION "public"."pgp_pub_encrypt"(text, bytea)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_encrypt_text'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_encrypt"(text, bytea, text);
CREATE FUNCTION "public"."pgp_pub_encrypt"(text, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_encrypt_text'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_encrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_encrypt_bytea"(bytea, bytea, text);
CREATE FUNCTION "public"."pgp_pub_encrypt_bytea"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_encrypt_bytea'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_encrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_encrypt_bytea"(bytea, bytea);
CREATE FUNCTION "public"."pgp_pub_encrypt_bytea"(bytea, bytea)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_encrypt_bytea'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_decrypt"(bytea, text, text);
CREATE FUNCTION "public"."pgp_sym_decrypt"(bytea, text, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_sym_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_decrypt"(bytea, text);
CREATE FUNCTION "public"."pgp_sym_decrypt"(bytea, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_sym_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_decrypt_bytea"(bytea, text);
CREATE FUNCTION "public"."pgp_sym_decrypt_bytea"(bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_decrypt_bytea"(bytea, text, text);
CREATE FUNCTION "public"."pgp_sym_decrypt_bytea"(bytea, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_encrypt"(text, text, text);
CREATE FUNCTION "public"."pgp_sym_encrypt"(text, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_encrypt_text'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_encrypt"(text, text);
CREATE FUNCTION "public"."pgp_sym_encrypt"(text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_encrypt_text'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_encrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_encrypt_bytea"(bytea, text);
CREATE FUNCTION "public"."pgp_sym_encrypt_bytea"(bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_encrypt_bytea'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_encrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_encrypt_bytea"(bytea, text, text);
CREATE FUNCTION "public"."pgp_sym_encrypt_bytea"(bytea, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_encrypt_bytea'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for seg_center
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_center"("public"."seg");
CREATE FUNCTION "public"."seg_center"("public"."seg")
  RETURNS "pg_catalog"."float4" AS '$libdir/seg', 'seg_center'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for seg_cmp
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_cmp"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_cmp"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."int4" AS '$libdir/seg', 'seg_cmp'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_cmp"("public"."seg", "public"."seg") IS 'btree comparison function';

-- ----------------------------
-- Function structure for seg_contained
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_contained"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_contained"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_contained'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_contained"("public"."seg", "public"."seg") IS 'contained in';

-- ----------------------------
-- Function structure for seg_contains
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_contains"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_contains"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_contains'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_contains"("public"."seg", "public"."seg") IS 'contains';

-- ----------------------------
-- Function structure for seg_different
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_different"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_different"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_different'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_different"("public"."seg", "public"."seg") IS 'different';

-- ----------------------------
-- Function structure for seg_ge
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_ge"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_ge"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_ge'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_ge"("public"."seg", "public"."seg") IS 'greater than or equal';

-- ----------------------------
-- Function structure for seg_gt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_gt"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_gt"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_gt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_gt"("public"."seg", "public"."seg") IS 'greater than';

-- ----------------------------
-- Function structure for seg_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_in"(cstring);
CREATE FUNCTION "public"."seg_in"(cstring)
  RETURNS "public"."seg" AS '$libdir/seg', 'seg_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for seg_inter
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_inter"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_inter"("public"."seg", "public"."seg")
  RETURNS "public"."seg" AS '$libdir/seg', 'seg_inter'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for seg_le
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_le"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_le"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_le'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_le"("public"."seg", "public"."seg") IS 'less than or equal';

-- ----------------------------
-- Function structure for seg_left
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_left"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_left"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_left'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_left"("public"."seg", "public"."seg") IS 'is left of';

-- ----------------------------
-- Function structure for seg_lower
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_lower"("public"."seg");
CREATE FUNCTION "public"."seg_lower"("public"."seg")
  RETURNS "pg_catalog"."float4" AS '$libdir/seg', 'seg_lower'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for seg_lt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_lt"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_lt"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_lt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_lt"("public"."seg", "public"."seg") IS 'less than';

-- ----------------------------
-- Function structure for seg_out
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_out"("public"."seg");
CREATE FUNCTION "public"."seg_out"("public"."seg")
  RETURNS "pg_catalog"."cstring" AS '$libdir/seg', 'seg_out'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for seg_over_left
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_over_left"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_over_left"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_over_left'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_over_left"("public"."seg", "public"."seg") IS 'overlaps or is left of';

-- ----------------------------
-- Function structure for seg_over_right
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_over_right"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_over_right"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_over_right'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_over_right"("public"."seg", "public"."seg") IS 'overlaps or is right of';

-- ----------------------------
-- Function structure for seg_overlap
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_overlap"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_overlap"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_overlap'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_overlap"("public"."seg", "public"."seg") IS 'overlaps';

-- ----------------------------
-- Function structure for seg_right
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_right"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_right"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_right'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_right"("public"."seg", "public"."seg") IS 'is right of';

-- ----------------------------
-- Function structure for seg_same
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_same"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_same"("public"."seg", "public"."seg")
  RETURNS "pg_catalog"."bool" AS '$libdir/seg', 'seg_same'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
COMMENT ON FUNCTION "public"."seg_same"("public"."seg", "public"."seg") IS 'same as';

-- ----------------------------
-- Function structure for seg_size
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_size"("public"."seg");
CREATE FUNCTION "public"."seg_size"("public"."seg")
  RETURNS "pg_catalog"."float4" AS '$libdir/seg', 'seg_size'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for seg_union
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_union"("public"."seg", "public"."seg");
CREATE FUNCTION "public"."seg_union"("public"."seg", "public"."seg")
  RETURNS "public"."seg" AS '$libdir/seg', 'seg_union'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for seg_upper
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."seg_upper"("public"."seg");
CREATE FUNCTION "public"."seg_upper"("public"."seg")
  RETURNS "pg_catalog"."float4" AS '$libdir/seg', 'seg_upper'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v1
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v1"();
CREATE FUNCTION "public"."uuid_generate_v1"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v1'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v1mc
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v1mc"();
CREATE FUNCTION "public"."uuid_generate_v1mc"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v1mc'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v3
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v3"("namespace" uuid, "name" text);
CREATE FUNCTION "public"."uuid_generate_v3"("namespace" uuid, "name" text)
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v3'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v4
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v4"();
CREATE FUNCTION "public"."uuid_generate_v4"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v4'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v5
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v5"("namespace" uuid, "name" text);
CREATE FUNCTION "public"."uuid_generate_v5"("namespace" uuid, "name" text)
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v5'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_nil
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_nil"();
CREATE FUNCTION "public"."uuid_nil"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_nil'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_dns
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_dns"();
CREATE FUNCTION "public"."uuid_ns_dns"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_dns'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_oid
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_oid"();
CREATE FUNCTION "public"."uuid_ns_oid"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_oid'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_url
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_url"();
CREATE FUNCTION "public"."uuid_ns_url"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_url'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_x500
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_x500"();
CREATE FUNCTION "public"."uuid_ns_x500"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_x500'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."activity_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."activity_log_id_seq1"
OWNED BY "public"."activity_log"."id";
SELECT setval('"public"."activity_log_id_seq1"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."command_log_id_seq"
OWNED BY "public"."command_log"."id";
SELECT setval('"public"."command_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."device_alert_id_seq"
OWNED BY "public"."device_alert"."id";
SELECT setval('"public"."device_alert_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."device_config_id_seq"
OWNED BY "public"."device_config"."id";
SELECT setval('"public"."device_config_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."device_status_id_seq"
OWNED BY "public"."device_status"."id";
SELECT setval('"public"."device_status_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."iot_data_id_seq"
OWNED BY "public"."iot_data"."id";
SELECT setval('"public"."iot_data_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."noti_notification_logs_log_id_seq"
OWNED BY "public"."noti_notification_logs"."log_id";
SELECT setval('"public"."noti_notification_logs_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."noti_notification_rules_rule_id_seq"
OWNED BY "public"."noti_notification_rules"."rule_id";
SELECT setval('"public"."noti_notification_rules_rule_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."noti_notification_types_type_id_seq"
OWNED BY "public"."noti_notification_types"."type_id";
SELECT setval('"public"."noti_notification_types_type_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."notification_devices_id_seq"
OWNED BY "public"."notification_devices"."id";
SELECT setval('"public"."notification_devices_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."notification_groups_id_seq"
OWNED BY "public"."notification_groups"."id";
SELECT setval('"public"."notification_groups_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."notification_logs_id_seq"
OWNED BY "public"."notification_logs"."id";
SELECT setval('"public"."notification_logs_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."notification_types_id_seq"
OWNED BY "public"."notification_types"."id";
SELECT setval('"public"."notification_types_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_activity_log_id_seq"
OWNED BY "public"."sd_activity_log"."id";
SELECT setval('"public"."sd_activity_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_activity_type_log_typeId_seq"
OWNED BY "public"."sd_activity_type_log"."typeId";
SELECT setval('"public"."sd_activity_type_log_typeId_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_admin_access_menu_admin_access_id_seq"
OWNED BY "public"."sd_admin_access_menu"."admin_access_id";
SELECT setval('"public"."sd_admin_access_menu_admin_access_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_air_control_air_control_id_seq"
OWNED BY "public"."sd_air_control"."air_control_id";
SELECT setval('"public"."sd_air_control_air_control_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_air_mod_air_mod_id_seq"
OWNED BY "public"."sd_air_mod"."air_mod_id";
SELECT setval('"public"."sd_air_mod_air_mod_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_air_period_air_period_id_seq"
OWNED BY "public"."sd_air_period"."air_period_id";
SELECT setval('"public"."sd_air_period_air_period_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_air_setting_warning_air_setting_warning_id_seq"
OWNED BY "public"."sd_air_setting_warning"."air_setting_warning_id";
SELECT setval('"public"."sd_air_setting_warning_air_setting_warning_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_air_warning_air_warning_id_seq"
OWNED BY "public"."sd_air_warning"."air_warning_id";
SELECT setval('"public"."sd_air_warning_air_warning_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_api_key_id_seq"
OWNED BY "public"."sd_api_key"."id";
SELECT setval('"public"."sd_api_key_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_audit_log_audit_id_seq"
OWNED BY "public"."sd_audit_log"."audit_id";
SELECT setval('"public"."sd_audit_log_audit_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_channel_template_id_seq"
OWNED BY "public"."sd_channel_template"."id";
SELECT setval('"public"."sd_channel_template_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_dashboard_config_id_seq"
OWNED BY "public"."sd_dashboard_config"."id";
SELECT setval('"public"."sd_dashboard_config_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_device_category_id_seq"
OWNED BY "public"."sd_device_category"."id";
SELECT setval('"public"."sd_device_category_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_device_group_id_seq"
OWNED BY "public"."sd_device_group"."id";
SELECT setval('"public"."sd_device_group_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_device_log_id_seq"
OWNED BY "public"."sd_device_log"."id";
SELECT setval('"public"."sd_device_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_device_member_id_seq"
OWNED BY "public"."sd_device_member"."id";
SELECT setval('"public"."sd_device_member_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_device_notification_config_id_seq"
OWNED BY "public"."sd_device_notification_config"."id";
SELECT setval('"public"."sd_device_notification_config_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_device_schedule_id_seq"
OWNED BY "public"."sd_device_schedule"."id";
SELECT setval('"public"."sd_device_schedule_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_device_status_history_id_seq"
OWNED BY "public"."sd_device_status_history"."id";
SELECT setval('"public"."sd_device_status_history_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_group_notification_config_id_seq"
OWNED BY "public"."sd_group_notification_config"."id";
SELECT setval('"public"."sd_group_notification_config_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_api_api_id_seq"
OWNED BY "public"."sd_iot_api"."api_id";
SELECT setval('"public"."sd_iot_api_api_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_device_action_device_action_user_id_seq"
OWNED BY "public"."sd_iot_device_action"."device_action_user_id";
SELECT setval('"public"."sd_iot_device_action_device_action_user_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_device_action_log_log_id_seq"
OWNED BY "public"."sd_iot_device_action_log"."log_id";
SELECT setval('"public"."sd_iot_device_action_log_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_device_action_user_device_action_user_id_seq"
OWNED BY "public"."sd_iot_device_action_user"."device_action_user_id";
SELECT setval('"public"."sd_iot_device_action_user_device_action_user_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_device_alarm_action_alarm_action_id_seq"
OWNED BY "public"."sd_iot_device_alarm_action"."alarm_action_id";
SELECT setval('"public"."sd_iot_device_alarm_action_alarm_action_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_device_device_id_seq"
OWNED BY "public"."sd_iot_device"."device_id";
SELECT setval('"public"."sd_iot_device_device_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_device_type_type_id_seq"
OWNED BY "public"."sd_iot_device_type"."type_id";
SELECT setval('"public"."sd_iot_device_type_type_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_group_group_id_seq"
OWNED BY "public"."sd_iot_group"."group_id";
SELECT setval('"public"."sd_iot_group_group_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_location_location_id_seq"
OWNED BY "public"."sd_iot_location"."location_id";
SELECT setval('"public"."sd_iot_location_location_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_mqtt_mqtt_id_seq"
OWNED BY "public"."sd_iot_mqtt"."mqtt_id";
SELECT setval('"public"."sd_iot_mqtt_mqtt_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_schedule_schedule_id_seq"
OWNED BY "public"."sd_iot_schedule"."schedule_id";
SELECT setval('"public"."sd_iot_schedule_schedule_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_sensor_sensor_id_seq"
OWNED BY "public"."sd_iot_sensor"."sensor_id";
SELECT setval('"public"."sd_iot_sensor_sensor_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_setting_setting_id_seq"
OWNED BY "public"."sd_iot_setting"."setting_id";
SELECT setval('"public"."sd_iot_setting_setting_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_token_token_id_seq"
OWNED BY "public"."sd_iot_token"."token_id";
SELECT setval('"public"."sd_iot_token_token_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_iot_type_type_id_seq"
OWNED BY "public"."sd_iot_type"."type_id";
SELECT setval('"public"."sd_iot_type_type_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_module_log_id_seq"
OWNED BY "public"."sd_module_log"."id";
SELECT setval('"public"."sd_module_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_notification_channel_id_seq"
OWNED BY "public"."sd_notification_channel"."id";
SELECT setval('"public"."sd_notification_channel_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_notification_condition_id_seq"
OWNED BY "public"."sd_notification_condition"."id";
SELECT setval('"public"."sd_notification_condition_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_notification_log_id_seq"
OWNED BY "public"."sd_notification_log"."id";
SELECT setval('"public"."sd_notification_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_notification_type_id_seq"
OWNED BY "public"."sd_notification_type"."id";
SELECT setval('"public"."sd_notification_type_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_report_data_id_seq"
OWNED BY "public"."sd_report_data"."id";
SELECT setval('"public"."sd_report_data_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_sensor_data_id_seq"
OWNED BY "public"."sd_sensor_data"."id";
SELECT setval('"public"."sd_sensor_data_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_system_setting_id_seq"
OWNED BY "public"."sd_system_setting"."id";
SELECT setval('"public"."sd_system_setting_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_user_access_menu_user_access_id_seq"
OWNED BY "public"."sd_user_access_menu"."user_access_id";
SELECT setval('"public"."sd_user_access_menu_user_access_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_user_file_id_seq"
OWNED BY "public"."sd_user_file"."id";
SELECT setval('"public"."sd_user_file_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_user_log_id_seq"
OWNED BY "public"."sd_user_log"."id";
SELECT setval('"public"."sd_user_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_user_log_type_log_type_id_seq"
OWNED BY "public"."sd_user_log_type"."log_type_id";
SELECT setval('"public"."sd_user_log_type_log_type_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_user_role_id_seq"
OWNED BY "public"."sd_user_role"."id";
SELECT setval('"public"."sd_user_role_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sd_user_roles_permision_role_type_id_seq"
OWNED BY "public"."sd_user_roles_permision"."role_type_id";
SELECT setval('"public"."sd_user_roles_permision_role_type_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."tnb_id_seq"
OWNED BY "public"."tnb"."id";
SELECT setval('"public"."tnb_id_seq"', 1, false);

-- ----------------------------
-- Primary Key structure for table activity_log
-- ----------------------------
ALTER TABLE "public"."activity_log" ADD CONSTRAINT "activity_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table command_log
-- ----------------------------
ALTER TABLE "public"."command_log" ADD CONSTRAINT "command_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table device_alert
-- ----------------------------
ALTER TABLE "public"."device_alert" ADD CONSTRAINT "device_alert_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table device_config
-- ----------------------------
ALTER TABLE "public"."device_config" ADD CONSTRAINT "uni_device_config_device_id" UNIQUE ("deviceId");

-- ----------------------------
-- Primary Key structure for table device_config
-- ----------------------------
ALTER TABLE "public"."device_config" ADD CONSTRAINT "device_config_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table device_status
-- ----------------------------
ALTER TABLE "public"."device_status" ADD CONSTRAINT "uni_device_status_device_id" UNIQUE ("deviceId");

-- ----------------------------
-- Primary Key structure for table device_status
-- ----------------------------
ALTER TABLE "public"."device_status" ADD CONSTRAINT "device_status_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table iot_data
-- ----------------------------
ALTER TABLE "public"."iot_data" ADD CONSTRAINT "iot_data_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table item
-- ----------------------------
CREATE INDEX "idx_item_deleted_at" ON "public"."item" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table item
-- ----------------------------
ALTER TABLE "public"."item" ADD CONSTRAINT "item_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table noti_notification_logs
-- ----------------------------
ALTER TABLE "public"."noti_notification_logs" ADD CONSTRAINT "noti_notification_logs_pkey" PRIMARY KEY ("log_id");

-- ----------------------------
-- Primary Key structure for table noti_notification_rules
-- ----------------------------
ALTER TABLE "public"."noti_notification_rules" ADD CONSTRAINT "noti_notification_rules_pkey" PRIMARY KEY ("rule_id");

-- ----------------------------
-- Primary Key structure for table noti_notification_types
-- ----------------------------
ALTER TABLE "public"."noti_notification_types" ADD CONSTRAINT "noti_notification_types_pkey" PRIMARY KEY ("type_id");

-- ----------------------------
-- Primary Key structure for table noti_notifications
-- ----------------------------
ALTER TABLE "public"."noti_notifications" ADD CONSTRAINT "noti_notifications_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table notification_devices
-- ----------------------------
ALTER TABLE "public"."notification_devices" ADD CONSTRAINT "uni_notification_devices_device_id" UNIQUE ("device_id");

-- ----------------------------
-- Primary Key structure for table notification_devices
-- ----------------------------
ALTER TABLE "public"."notification_devices" ADD CONSTRAINT "notification_devices_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table notification_groups
-- ----------------------------
ALTER TABLE "public"."notification_groups" ADD CONSTRAINT "notification_groups_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table notification_groups_devices_notification_devices
-- ----------------------------
ALTER TABLE "public"."notification_groups_devices_notification_devices" ADD CONSTRAINT "notification_groups_devices_notification_devices_pkey" PRIMARY KEY ("notificationGroupsId", "notificationDevicesId");

-- ----------------------------
-- Primary Key structure for table notification_logs
-- ----------------------------
ALTER TABLE "public"."notification_logs" ADD CONSTRAINT "notification_logs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table notification_types
-- ----------------------------
ALTER TABLE "public"."notification_types" ADD CONSTRAINT "uni_notification_types_name" UNIQUE ("name");
ALTER TABLE "public"."notification_types" ADD CONSTRAINT "uni_notification_types_code" UNIQUE ("code");

-- ----------------------------
-- Primary Key structure for table notification_types
-- ----------------------------
ALTER TABLE "public"."notification_types" ADD CONSTRAINT "notification_types_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_activity_log
-- ----------------------------
ALTER TABLE "public"."sd_activity_log" ADD CONSTRAINT "sd_activity_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_activity_type_log
-- ----------------------------
ALTER TABLE "public"."sd_activity_type_log" ADD CONSTRAINT "sd_activity_type_log_pkey" PRIMARY KEY ("typeId");

-- ----------------------------
-- Primary Key structure for table sd_admin_access_menu
-- ----------------------------
ALTER TABLE "public"."sd_admin_access_menu" ADD CONSTRAINT "sd_admin_access_menu_pkey" PRIMARY KEY ("admin_access_id");

-- ----------------------------
-- Primary Key structure for table sd_air_control
-- ----------------------------
ALTER TABLE "public"."sd_air_control" ADD CONSTRAINT "sd_air_control_pkey" PRIMARY KEY ("air_control_id");

-- ----------------------------
-- Primary Key structure for table sd_air_control_device_map
-- ----------------------------
ALTER TABLE "public"."sd_air_control_device_map" ADD CONSTRAINT "sd_air_control_device_map_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_air_control_log
-- ----------------------------
ALTER TABLE "public"."sd_air_control_log" ADD CONSTRAINT "sd_air_control_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_air_mod
-- ----------------------------
ALTER TABLE "public"."sd_air_mod" ADD CONSTRAINT "sd_air_mod_pkey" PRIMARY KEY ("air_mod_id");

-- ----------------------------
-- Primary Key structure for table sd_air_mod_device_map
-- ----------------------------
ALTER TABLE "public"."sd_air_mod_device_map" ADD CONSTRAINT "sd_air_mod_device_map_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_air_period
-- ----------------------------
ALTER TABLE "public"."sd_air_period" ADD CONSTRAINT "sd_air_period_pkey" PRIMARY KEY ("air_period_id");

-- ----------------------------
-- Primary Key structure for table sd_air_period_device_map
-- ----------------------------
ALTER TABLE "public"."sd_air_period_device_map" ADD CONSTRAINT "sd_air_period_device_map_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_air_setting_warning
-- ----------------------------
ALTER TABLE "public"."sd_air_setting_warning" ADD CONSTRAINT "sd_air_setting_warning_pkey" PRIMARY KEY ("air_setting_warning_id");

-- ----------------------------
-- Primary Key structure for table sd_air_setting_warning_device_map
-- ----------------------------
ALTER TABLE "public"."sd_air_setting_warning_device_map" ADD CONSTRAINT "sd_air_setting_warning_device_map_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_air_warning
-- ----------------------------
ALTER TABLE "public"."sd_air_warning" ADD CONSTRAINT "sd_air_warning_pkey" PRIMARY KEY ("air_warning_id");

-- ----------------------------
-- Primary Key structure for table sd_air_warning_device_map
-- ----------------------------
ALTER TABLE "public"."sd_air_warning_device_map" ADD CONSTRAINT "sd_air_warning_device_map_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_alarm_process_log
-- ----------------------------
ALTER TABLE "public"."sd_alarm_process_log" ADD CONSTRAINT "sd_alarm_process_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_alarm_process_log_email
-- ----------------------------
ALTER TABLE "public"."sd_alarm_process_log_email" ADD CONSTRAINT "sd_alarm_process_log_email_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_alarm_process_log_line
-- ----------------------------
ALTER TABLE "public"."sd_alarm_process_log_line" ADD CONSTRAINT "sd_alarm_process_log_line_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_alarm_process_log_mqtt
-- ----------------------------
ALTER TABLE "public"."sd_alarm_process_log_mqtt" ADD CONSTRAINT "sd_alarm_process_log_mqtt_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_alarm_process_log_sms
-- ----------------------------
ALTER TABLE "public"."sd_alarm_process_log_sms" ADD CONSTRAINT "sd_alarm_process_log_sms_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_alarm_process_log_telegram
-- ----------------------------
ALTER TABLE "public"."sd_alarm_process_log_telegram" ADD CONSTRAINT "sd_alarm_process_log_telegram_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_alarm_process_log_temp
-- ----------------------------
ALTER TABLE "public"."sd_alarm_process_log_temp" ADD CONSTRAINT "sd_alarm_process_log_temp_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table sd_api_key
-- ----------------------------
ALTER TABLE "public"."sd_api_key" ADD CONSTRAINT "uni_sd_api_key_api_key" UNIQUE ("api_key");

-- ----------------------------
-- Primary Key structure for table sd_api_key
-- ----------------------------
ALTER TABLE "public"."sd_api_key" ADD CONSTRAINT "sd_api_key_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_audit_log
-- ----------------------------
ALTER TABLE "public"."sd_audit_log" ADD CONSTRAINT "sd_audit_log_pkey" PRIMARY KEY ("audit_id");

-- ----------------------------
-- Primary Key structure for table sd_channel_template
-- ----------------------------
ALTER TABLE "public"."sd_channel_template" ADD CONSTRAINT "sd_channel_template_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_dashboard_config
-- ----------------------------
ALTER TABLE "public"."sd_dashboard_config" ADD CONSTRAINT "sd_dashboard_config_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_device_category
-- ----------------------------
ALTER TABLE "public"."sd_device_category" ADD CONSTRAINT "sd_device_category_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_device_group
-- ----------------------------
ALTER TABLE "public"."sd_device_group" ADD CONSTRAINT "sd_device_group_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_device_log
-- ----------------------------
ALTER TABLE "public"."sd_device_log" ADD CONSTRAINT "sd_device_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_device_member
-- ----------------------------
ALTER TABLE "public"."sd_device_member" ADD CONSTRAINT "sd_device_member_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_device_notification_config
-- ----------------------------
ALTER TABLE "public"."sd_device_notification_config" ADD CONSTRAINT "sd_device_notification_config_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_device_schedule
-- ----------------------------
ALTER TABLE "public"."sd_device_schedule" ADD CONSTRAINT "sd_device_schedule_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_device_status_history
-- ----------------------------
ALTER TABLE "public"."sd_device_status_history" ADD CONSTRAINT "sd_device_status_history_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_group_notification_config
-- ----------------------------
ALTER TABLE "public"."sd_group_notification_config" ADD CONSTRAINT "sd_group_notification_config_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_iot_alarm_device
-- ----------------------------
ALTER TABLE "public"."sd_iot_alarm_device" ADD CONSTRAINT "sd_iot_alarm_device_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_iot_alarm_device_event
-- ----------------------------
ALTER TABLE "public"."sd_iot_alarm_device_event" ADD CONSTRAINT "sd_iot_alarm_device_event_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_iot_api
-- ----------------------------
ALTER TABLE "public"."sd_iot_api" ADD CONSTRAINT "sd_iot_api_pkey" PRIMARY KEY ("api_id");

-- ----------------------------
-- Uniques structure for table sd_iot_device
-- ----------------------------
ALTER TABLE "public"."sd_iot_device" ADD CONSTRAINT "uni_sd_iot_device_sn" UNIQUE ("sn");

-- ----------------------------
-- Primary Key structure for table sd_iot_device
-- ----------------------------
ALTER TABLE "public"."sd_iot_device" ADD CONSTRAINT "sd_iot_device_pkey" PRIMARY KEY ("device_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_device_action
-- ----------------------------
ALTER TABLE "public"."sd_iot_device_action" ADD CONSTRAINT "sd_iot_device_action_pkey" PRIMARY KEY ("device_action_user_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_device_action_log
-- ----------------------------
ALTER TABLE "public"."sd_iot_device_action_log" ADD CONSTRAINT "sd_iot_device_action_log_pkey" PRIMARY KEY ("log_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_device_action_user
-- ----------------------------
ALTER TABLE "public"."sd_iot_device_action_user" ADD CONSTRAINT "sd_iot_device_action_user_pkey" PRIMARY KEY ("device_action_user_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_device_alarm_action
-- ----------------------------
ALTER TABLE "public"."sd_iot_device_alarm_action" ADD CONSTRAINT "sd_iot_device_alarm_action_pkey" PRIMARY KEY ("alarm_action_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_device_type
-- ----------------------------
ALTER TABLE "public"."sd_iot_device_type" ADD CONSTRAINT "sd_iot_device_type_pkey" PRIMARY KEY ("type_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_email
-- ----------------------------
ALTER TABLE "public"."sd_iot_email" ADD CONSTRAINT "sd_iot_email_pkey" PRIMARY KEY ("email_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_group
-- ----------------------------
ALTER TABLE "public"."sd_iot_group" ADD CONSTRAINT "sd_iot_group_pkey" PRIMARY KEY ("group_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_host
-- ----------------------------
ALTER TABLE "public"."sd_iot_host" ADD CONSTRAINT "sd_iot_host_pkey" PRIMARY KEY ("host_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_influxdb
-- ----------------------------
ALTER TABLE "public"."sd_iot_influxdb" ADD CONSTRAINT "sd_iot_influxdb_pkey" PRIMARY KEY ("influxdb_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_line
-- ----------------------------
ALTER TABLE "public"."sd_iot_line" ADD CONSTRAINT "sd_iot_line_pkey" PRIMARY KEY ("line_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_location
-- ----------------------------
ALTER TABLE "public"."sd_iot_location" ADD CONSTRAINT "sd_iot_location_pkey" PRIMARY KEY ("location_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_mqtt
-- ----------------------------
ALTER TABLE "public"."sd_iot_mqtt" ADD CONSTRAINT "sd_iot_mqtt_pkey" PRIMARY KEY ("mqtt_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_nodered
-- ----------------------------
ALTER TABLE "public"."sd_iot_nodered" ADD CONSTRAINT "sd_iot_nodered_pkey" PRIMARY KEY ("nodered_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_schedule
-- ----------------------------
ALTER TABLE "public"."sd_iot_schedule" ADD CONSTRAINT "sd_iot_schedule_pkey" PRIMARY KEY ("schedule_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_schedule_device
-- ----------------------------
ALTER TABLE "public"."sd_iot_schedule_device" ADD CONSTRAINT "sd_iot_schedule_device_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_iot_sensor
-- ----------------------------
ALTER TABLE "public"."sd_iot_sensor" ADD CONSTRAINT "sd_iot_sensor_pkey" PRIMARY KEY ("sensor_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_setting
-- ----------------------------
ALTER TABLE "public"."sd_iot_setting" ADD CONSTRAINT "sd_iot_setting_pkey" PRIMARY KEY ("setting_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_sms
-- ----------------------------
ALTER TABLE "public"."sd_iot_sms" ADD CONSTRAINT "sd_iot_sms_pkey" PRIMARY KEY ("sms_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_telegram
-- ----------------------------
ALTER TABLE "public"."sd_iot_telegram" ADD CONSTRAINT "sd_iot_telegram_pkey" PRIMARY KEY ("telegram_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_token
-- ----------------------------
ALTER TABLE "public"."sd_iot_token" ADD CONSTRAINT "sd_iot_token_pkey" PRIMARY KEY ("token_id");

-- ----------------------------
-- Primary Key structure for table sd_iot_type
-- ----------------------------
ALTER TABLE "public"."sd_iot_type" ADD CONSTRAINT "sd_iot_type_pkey" PRIMARY KEY ("type_id");

-- ----------------------------
-- Primary Key structure for table sd_module_log
-- ----------------------------
ALTER TABLE "public"."sd_module_log" ADD CONSTRAINT "sd_module_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_mqtt_host
-- ----------------------------
ALTER TABLE "public"."sd_mqtt_host" ADD CONSTRAINT "sd_mqtt_host_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_mqtt_log
-- ----------------------------
ALTER TABLE "public"."sd_mqtt_log" ADD CONSTRAINT "sd_mqtt_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_notification_channel
-- ----------------------------
ALTER TABLE "public"."sd_notification_channel" ADD CONSTRAINT "sd_notification_channel_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_notification_condition
-- ----------------------------
ALTER TABLE "public"."sd_notification_condition" ADD CONSTRAINT "sd_notification_condition_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_notification_log
-- ----------------------------
ALTER TABLE "public"."sd_notification_log" ADD CONSTRAINT "sd_notification_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_notification_type
-- ----------------------------
ALTER TABLE "public"."sd_notification_type" ADD CONSTRAINT "sd_notification_type_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_report_data
-- ----------------------------
ALTER TABLE "public"."sd_report_data" ADD CONSTRAINT "sd_report_data_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_schedule_process_log
-- ----------------------------
ALTER TABLE "public"."sd_schedule_process_log" ADD CONSTRAINT "sd_schedule_process_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_sensor_data
-- ----------------------------
ALTER TABLE "public"."sd_sensor_data" ADD CONSTRAINT "sd_sensor_data_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table sd_system_setting
-- ----------------------------
ALTER TABLE "public"."sd_system_setting" ADD CONSTRAINT "uni_sd_system_setting_key" UNIQUE ("key");

-- ----------------------------
-- Primary Key structure for table sd_system_setting
-- ----------------------------
ALTER TABLE "public"."sd_system_setting" ADD CONSTRAINT "sd_system_setting_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sd_user
-- ----------------------------
CREATE UNIQUE INDEX "idx_sd_user_email" ON "public"."sd_user" USING btree (
  "email" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sd_user_password_reset_token" ON "public"."sd_user" USING btree (
  "password_reset_token" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
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

-- ----------------------------
-- Primary Key structure for table sd_user_access_menu
-- ----------------------------
ALTER TABLE "public"."sd_user_access_menu" ADD CONSTRAINT "sd_user_access_menu_pkey" PRIMARY KEY ("user_access_id");

-- ----------------------------
-- Primary Key structure for table sd_user_file
-- ----------------------------
ALTER TABLE "public"."sd_user_file" ADD CONSTRAINT "sd_user_file_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_user_log
-- ----------------------------
ALTER TABLE "public"."sd_user_log" ADD CONSTRAINT "sd_user_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_user_log_type
-- ----------------------------
ALTER TABLE "public"."sd_user_log_type" ADD CONSTRAINT "sd_user_log_type_pkey" PRIMARY KEY ("log_type_id");

-- ----------------------------
-- Primary Key structure for table sd_user_role
-- ----------------------------
ALTER TABLE "public"."sd_user_role" ADD CONSTRAINT "sd_user_role_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sd_user_roles_access
-- ----------------------------
ALTER TABLE "public"."sd_user_roles_access" ADD CONSTRAINT "sd_user_roles_access_pkey" PRIMARY KEY ("role_id", "role_type_id");

-- ----------------------------
-- Primary Key structure for table sd_user_roles_permision
-- ----------------------------
ALTER TABLE "public"."sd_user_roles_permision" ADD CONSTRAINT "sd_user_roles_permision_pkey" PRIMARY KEY ("role_type_id");

-- ----------------------------
-- Primary Key structure for table tnb
-- ----------------------------
ALTER TABLE "public"."tnb" ADD CONSTRAINT "tnb_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table user
-- ----------------------------
CREATE INDEX "idx_user_deleted_at" ON "public"."user" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_user_email" ON "public"."user" USING btree (
  "email" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table user
-- ----------------------------
ALTER TABLE "public"."user" ADD CONSTRAINT "user_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table noti_notification_logs
-- ----------------------------
ALTER TABLE "public"."noti_notification_logs" ADD CONSTRAINT "fk_noti_notifications_notification_logs" FOREIGN KEY ("notification_id") REFERENCES "public"."noti_notifications" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table sd_activity_log
-- ----------------------------
ALTER TABLE "public"."sd_activity_log" ADD CONSTRAINT "fk_sd_activity_log_module" FOREIGN KEY ("modules_id") REFERENCES "public"."sd_module_log" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
