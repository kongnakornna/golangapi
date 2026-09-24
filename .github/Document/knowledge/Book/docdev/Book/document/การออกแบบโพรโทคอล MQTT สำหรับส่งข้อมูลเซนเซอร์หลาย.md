
# การออกแบบโพรโทคอล MQTT สำหรับส่งข้อมูลเซนเซอร์หลายชนิด

การออกแบบโพรโทคอล MQTT สำหรับเซนเซอร์หลายชนิดควรคิดเป็น “มาตรฐานกลาง” ของ topic + payload + metadata เพื่อให้ทุก device พูดภาษาเดียวกันและระบบ AoT ต่อยอดง่าย ไม่ควรปล่อยให้แต่ละ node ออกแบบเองแบบ ad‑hoc.

ด้านล่างเป็นแนวทางออกแบบเป็นขั้นๆ โดยอ้างอิงแนวปฏิบัติจาก MQTT community + convention อย่าง Homie และ Sparkplug B ที่ใช้จริงใน IIoT ครับ.[^1][^2][^3][^4][^5][^6]

***

## 1) หลักคิดเรื่อง Topic Tree

หัวใจคือออกแบบ topic hierarchy ให้:

- อ่านแล้วเดาได้ทันทีว่าเป็น site ไหน / device ไหน / sensor อะไร / data หรือ command.
- ใช้ wildcard ได้ดี เช่น subscribe ทั้งอาคาร, ทั้งชั้น, หรือทั้งเซนเซอร์ชนิดเดียวกันง่ายๆ ตามแนวแนะนำของผู้ให้บริการ MQTT หลายเจ้า.[^7][^8][^9][^1]

โครงสร้างที่เหมาะกับหลาย site + หลาย sensor type (ตัวอย่าง):

`iot/{site}/{building}/{floor}/{deviceId}/{channel}/{property}`

เช่น:

- `iot/siteA/bld1/f2/node-01/env/temperature`
- `iot/siteA/bld1/f2/node-01/env/humidity`
- `iot/siteA/bld1/f2/node-01/pm/pm2_5`
- `iot/siteA/bld1/f1/lidar-01/occupancy/zone1`
- `iot/siteA/bld1/f1/cctv-01/event/person_detected`

แนวทางจากหลายบทความ MQTT คือ:

- ใส่ “context” ไว้ใน topic (เช่น สถานที่, device) แต่ใส่ค่าตัวเลข/สถานะลงใน payload แทน.[^7]
- แยก data กับ command คนละ topic เช่น
    - Data: `.../{deviceId}/data/...`
    - Command: `.../{deviceId}/cmd/...`.

***

## 2) การจัด payload (JSON) สำหรับ Multi-sensor

มีสอง pattern หลัก ๆ ที่ใช้กันจริง:

1. **หนึ่ง topic ต่อหนึ่งค่า** – payload เป็นเลขเดียว (เรียบง่าย, scale ได้ดี).
2. **หนึ่ง topic ต่อหนึ่ง device** – payload เป็น JSON ที่รวมหลาย sensor field เหมาะเมื่อ device ตัวเดียวมีเซนเซอร์หลายค่า.[^10][^7]

สำหรับ MQTT+AoT ที่ต้องการความยืดหยุ่น แนะนำ pattern 2 สำหรับ “node” และยังเปิดทางใช้ pattern 1 กับ topic ย่อยเฉพาะ sensor ได้ถ้าต้องการ:

ตัวอย่าง payload JSON แบบ multi-sensor ต่อ 1 device:

```json
{
  "ts": "2025-11-25T07:55:00Z",
  "deviceId": "node-01",
  "sensors": {
    "temperature": { "value": 27.4, "unit": "C" },
    "humidity":    { "value": 56.2, "unit": "%" },
    "co2":         { "value": 650,  "unit": "ppm" },
    "pm2_5":       { "value": 18,   "unit": "ug/m3" }
  }
}
```

แนวคิดนี้สอดคล้องกับคำแนะนำของหลายบทความที่เสนอให้ใช้ JSON payload แทนการพยายามยัดหลายฟิลด์ใน string เดียวอ่านยาก.[^10][^7]

***

## 3) ใช้ Convention สำเร็จรูป (Homie / Sparkplug B)

ถ้าต้องการความ “มาตรฐาน + auto-discovery” จริงจัง แนะนำดู 2 ตัวนี้:

- **Homie MQTT Convention**
    - กำหนดรูปแบบ topic + payload + metadata สำหรับ IoT device อย่างชัดเจน เช่น level: device → node → property.
    - รองรับ auto-discovery, การประกาศ metadata (ชื่อ, unit, data type) ผ่าน MQTT เอง ทำให้ controller / dashboard ต่อได้ง่ายโดยไม่ต้อง config ทีละตัว.[^2][^11][^6]
    - เน้น IoT ทั่วไป (smart home, building).
- **Sparkplug B**
    - เป็น spec บน MQTT ที่ออกแบบมาสำหรับ Industrial IoT, มี topic namespace, payload format, และ state management (birth/death) มาตรฐานเดียวกัน.[^3][^4][^12][^13]
    - ดีมากเมื่อมี PLC/SCADA/IIoT gateway จำนวนมากและต้องการ interoperability สูง.

สำหรับโปรเจกต์ที่ออกแบบเองทั้ง stack:

- ถ้า scope ใกล้ smart building / smart home / dashboard เอง → ใช้ idea จาก Homie (แต่ปรับง่าย ๆ) ก็พอ.
- ถ้าจะ integrate กับ SCADA/PLC หรือ vendor อุตสาหกรรม → พิจารณาใช้ Sparkplug B ตรงๆ เพื่อไม่ต้องเขียน parser แปลกๆ เอง.[^4][^5][^3]

***

## 4) การใส่ Metadata ของ Footprint และเซนเซอร์

เพื่อให้ AoT layer ทำ spatial analytics ได้ ต้องมี metadata มาตรฐานของแต่ละ sensor:

- แนะนำแยกเป็น **config topic** กับ **data topic** ตามแนวที่หลาย convention ทำ (Homie ก็ใช้แนวคล้ายกัน).[^6][^2]
    - Config (retained):
        - `iot/{site}/{building}/{floor}/{deviceId}/$config`
    - Data (stream):
        - `iot/{site}/{building}/{floor}/{deviceId}/data`

ตัวอย่าง payload config (publish เป็น retained message):

```json
{
  "deviceId": "node-01",
  "type": "env_node",
  "location": { "x": 12.3, "y": 4.5, "z": 2.7, "floor": 2 },
  "footprint": {
    "shape": "circle",
    "radius_m": 5.0
  },
  "sensors": {
    "temperature": { "unit": "C", "min": -10, "max": 60 },
    "humidity":    { "unit": "%", "min": 0,  "max": 100 },
    "co2":         { "unit": "ppm", "min": 400, "max": 5000 },
    "pm2_5":       { "unit": "ug/m3", "min": 0, "max": 500 }
  }
}
```

การใช้ retained + QoS ที่เหมาะสมสำหรับ config เป็นแนวปฏิบัติที่ถูกแนะนำใน Homie เพื่อให้ client ใหม่อ่าน config ได้ทันทีหลัง subscribe โดยไม่ต้องรออุปกรณ์ online.[^11][^2]

***

## 5) QoS, Retain, และกลยุทธ์ส่งข้อมูล

สำหรับเซนเซอร์หลายชนิดให้คิดแบบนี้:

- **Data ที่สำคัญแต่อัพเดตเรื่อย ๆ** เช่น Temp/RH/CO₂/PM
    - QoS 0 หรือ 1 ก็พอ ขึ้นกับความคับคั่งของ network.
    - ไม่จำเป็นต้อง retain ทุกค่า (จะเปลือง broker) แต่อาจ retain “last known value” แยก topic เช่น `.../last` ถ้าจำเป็น.
- **State / Metadata / Config / Birth-Death**
    - ใช้ QoS 1 + retained ตามแนว Homie และ Sparkplug B เพื่อให้มีความน่าเชื่อถือของสถานะอุปกรณ์.[^2][^3][^4]
- **Multi-sensor ใน device เดียว**
    - เลือกว่าจะส่ง:
        - เป็น snapshot รวมทุกค่าใน message เดียวทุก X วินาที (ง่ายต่อ time-series DB, ใช้ JSON) หรือ
        - ใช้ report-by-exception (Sparkplug B concept) คือส่งเมื่อค่ามีการเปลี่ยนแปลงเกิน threshold เพื่อลด traffic.[^12][^5]

***

## 6) ตัวอย่างโครงสร้างรวม (สำหรับโปรเจกต์คุณ)

สมมติคุณมี node ที่เก็บ LiDAR+CCTV+Env+PM ในพื้นที่เดียวกัน:

- Config:
    - `iot/solarFarmA/bld1/f1/node-01/$config`
- Data stream หลัก:
    - `iot/solarFarmA/bld1/f1/node-01/data`  → JSON รวมทุกค่า + timestamp.
- Event เฉพาะ:
    - LiDAR occupancy: `iot/solarFarmA/bld1/f1/lidar-01/occupancy/zone1`
    - CCTV detection: `iot/solarFarmA/bld1/f1/cctv-01/event/person`

แบบนี้ AoT side (InfluxDB, PostgreSQL+Timescale, หรือ data lake) จะ map ได้ง่ายว่า:

- node‑01 = โหนดเซนเซอร์สิ่งแวดล้อมใน footprint โซน X.
- lidar‑01 / cctv‑01 = เซนเซอร์เชิงภาพ/ระยะ สำหรับ occupancy ใน footprint เดียวกันแต่มี event-based topic แยกต่างหาก.

***

ถ้าบอกเพิ่มได้ว่าตอนนี้ใช้ broker / stack อะไร (เช่น EMQX, Mosquitto, AWS IoT, ThingsBoard, Node‑RED) กับต้องการ integrate กับ InfluxDB / PostgreSQL ยังไง จะช่วยออกแบบ schema + topic tree + retention policy ให้ match ระบบของคุณแบบละเอียดได้อีกชั้นครับ.
<span style="display:none">[^14][^15][^16][^17][^18][^19][^20]</span>

<div align="center">⁂</div>

[^1]: https://www.hivemq.com/blog/mqtt-essentials-part-5-mqtt-topics-best-practices/

[^2]: https://homieiot.github.io/specification/spec-core-v3_0_1/

[^3]: https://ubidots.com/blog/mqtt-sparkplug-vs-plain-mqtt/

[^4]: https://www.racoman.com/blog/industrial-protocols/fundamentals-of-mqtt/sparkplug-b

[^5]: https://corsosystems.com/posts/mqtt-and-sparkplug-b-simplified

[^6]: https://homieiot.github.io/specification/

[^7]: https://pi3g.com/mqtt-topic-tree-design-best-practices-tips-examples/

[^8]: https://www.emqx.com/en/blog/advanced-features-of-mqtt-topics

[^9]: https://thingsboard.io/docs/mqtt-broker/user-guide/topics/

[^10]: https://community.home-assistant.io/t/multi-sensor-using-same-mqtt-topic/172552

[^11]: https://homieiot.github.io

[^12]: https://blog.isa.org/iot-architecture-with-mqtt-sparkplugb

[^13]: https://sparkplug.eclipse.org/specification/

[^14]: https://erlangforums.com/t/what-is-the-better-practice-when-designing-mqtt-topics-for-thousands-of-devices-of-multiple-customers/1547

[^15]: https://docs.aws.amazon.com/whitepapers/latest/designing-mqtt-topics-aws-iot-core/mqtt-design-best-practices.html

[^16]: https://github.com/homieiot/convention

[^17]: https://stackoverflow.com/questions/72640570/what-is-the-better-practice-when-designing-mqtt-topics-for-thousands-of-devices

[^18]: https://documentation.meraki.com/IoT/MT_-_Sensors/Design_and_Configure/MT_MQTT_Setup_Guide

[^19]: https://flows.nodered.org/node/node-red-contrib-homie-convention

[^20]: https://2smart.com/docs-resources/articles/mqtt-convention-2smart

