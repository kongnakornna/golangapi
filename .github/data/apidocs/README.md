# ICMON API คู่มือ — ดัชนีโมดูล

> Base URL: `/api` · Framework: Go + Chi · Auth: JWT (RS256) Bearer

นี้คือคู่มือ API แบบแยกไฟล์ตาม module ครอบคลุมทุก endpoint ของโปรเจกต์ icmongolang

## ดัชนีไฟล์

| ไฟล์ | ครอบคลุม module |
|------|-----------------|
| [`auth-core.md`](./auth-core.md) | Authentication `/auth`, Users `/user`, Items `/item`, Orders (Kafka), WebSocket `/ws`, System & Health |
| [`iot-mqtt.md`](./iot-mqtt.md) | IoT `/iot` (รวม `/devicemqtt` ใหม่), MQTT Gateway `/mqtt`, InfluxDB `/influx`, Alarm `/alarm` |
| [`settings.md`](./settings.md) | Settings `/settings` (Master + Device + Schedule + Integration + Alarm + Logs) |
| [`business.md`](./business.md) | Quotation `/quotation`, Purchase Order `/purchase-orders`, Payment `/payments`, Receipt `/receipts`, Job `/job`, Customer `/customer`, Car `/car` |
| [`infra.md`](./infra.md) | Dashboard `/dashboard`, Document `/documents`, Email `/email`, Batch `/batch`, I18n `/i18n`, WOS `/wos`, Report `/reports`, VectorData `/vectordata`, Elasticsearch `/elasticsearch` |

## ภาพรวม

- **Public endpoint:** `/api/auth/*` (login/signin/publickey/verifyemail/forgotpassword/resetpassword), `/api/register`, `/api/login`, `/api/signin`, `/api/iot/*`, `/api/mqtt/{subscriptions,status,gettopicdata,devicecontrol}`, `/api/elasticsearch/health`, `/api/vectordata/health`, `/`, `/health`, `/api/ping`, `/api/health`, `/metrics`, `/apimetric`, `/uploads/avatar/*`, `/swagger/*`, `/api/ws`
- **ทุก module business ใช้:** `Verifier(true)` → `Authenticator()` → `CurrentUser()` → `ActiveUser()` (Bearer token)

## Enum ที่ใช้บ่อย

- **hardware_type_name**: 1=Sensor, 2=IO Sensor, 3=IO Control, 4=Critical Sensor
- **layoutapp**: 1=Right Menu, 2=Card, 3=Left Menu, 4=Footer Menu
- **PO status**: DRAFT / SENT / CONFIRMED / SHIPPED / RECEIVED / CANCELLED
- **WOS order status**: pending / confirmed / shipped / delivered / cancelled
- **Report type**: daily_sales / inventory_summary / customer_list / invoice / credit_note / debit_note

## ฉบับรวม

ดู [`../apidocs_full.md`](../apidocs_full.md) สำหรับคู่มือฉบับรวมทุก module ในไฟล์เดียว
