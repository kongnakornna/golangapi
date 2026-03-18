Vibe coding and Prompt Engineering: โครงสร้างและการเขียน
1.สร้างบทนำ
2.สร้างบทนิยาม
3.สร้างบทหัวข้อ
4.ออกแบบ workflow
5.case study
6.โคงสร้าง
6.1 Prompt Engineering
6.2 Vibe coding
6.3 AI Tool
6.2 แบบจำลองภาษาขนาดใหญ่ (LLM) เป็นโปรแกรม ปัญญาประดิษฐ์ (AI)
7.Template 
8.case study  IoT  monitoring fram
8.0 ออก แบบ REST API ด้วย nust js  แบบ Microservices CRUD  + ORM  +JWT + MVC 
8.1.node red
8.2.mqtt
8.3.gafana
8.4.n8n
8.5.python ai llm
8.6.node ja nustjs api
8.7.php codeingator3  hmvc  fontend  ,jqury ,css ,html
8.8.database postgreaasl
8.9.Redis
8.10.influcdb
8.11.kafka
8.12.web socket io
8.13 python fast api

 ระบบ notification  system โดยใช้
NustJS framwork By Node JS

1.type ORM 
2.databse postgresssql
3.redis cache

IoT Notification System Architecture Documentation
1. บทนำ
ระบบ IoT Notification System เป็นระบบแจ้งเตือนอัจฉริยะที่พัฒนาขึ้นด้วย NestJS Framework โดยมีวัตถุประสงค์หลักในการจัดการและแจ้งเตือนข้อมูลจากเซ็นเซอร์ IoT ต่างๆ ระบบนี้รองรับการแจ้งเตือนผ่านช่องทางหลากหลายรูปแบบ พร้อมทั้งมีกลไกการจัดการสถานะ การหน่วงเวลา และการบันทึกประวัติอย่างครบวงจร

วัตถุประสงค์
จัดการข้อมูลจากเซ็นเซอร์ IoT ผ่าน MQTT
ตรวจสอบเงื่อนไขและแจ้งเตือนตามระดับความสำคัญ
รองรับช่องทางการแจ้งเตือนที่หลากหลาย
บันทึกประวัติและสร้างรายงาน
รองรับการทำงานแบบ Real-time

2. บทนิยาม
2.1 คำศัพท์สำคัญ
IoT Device: อุปกรณ์ IoT ที่เชื่อมต่อกับระบบ

Sensor: เซ็นเซอร์ที่วัดค่าต่างๆ เช่น อุณหภูมิ ความชื้น
Notification Channel: ช่องทางการแจ้งเตือน
Alarm Status: สถานะการแจ้งเตือน
Cooldown: เวลาหน่วงระหว่างการแจ้งเตือนซ้ำ
Recovery: สถานะเมื่อค่ากลับสู่ภาวะปกติ

2.2 สถานะการแจ้งเตือน

Warning (1): เตือนเบื้องต้น
Alarm (2): แจ้งเตือนฉุกเฉิน
Recovery Warning (3): ฟื้นตัวจาก Warning
Recovery Alarm (4): ฟื้นตัวจาก Alarm
Normal (5): สภาวะปกติ

3  ส่งการแจ้งเตือน
    1. ส่งการแจ้งเตือน  email   
    2. ส่งการแจ้งเตือน line notification
    3. ส่งการแจ้งเตือน discord notification
	4. ส่งการแจ้งเตือน  telegram notification
	5. ส่งการแจ้งเตือน  sms notification
    6. ส่งการแจ้งเตือน web dashboard notification
    7. ส่งการแจ้งเตือน AI CHART bot  fast api pythone



หลักการทำงานของระบบ
1.รับข้อมูลจากเซ็นเซอร์ ผ่าน MQTT
2.ตรวจสอบเงื่อนไข กับค่าที่ตั้งไว้
3.อัปเดตสถานะ และบันทึกประวัติ
4.สร้างการแจ้งเตือน ตามระดับความสำคัญ
5.ตรวจสอบ cooldown ด้วย Redis หรือ notification  alarm configuration 
6.ส่งผ่านช่องทางต่างๆ ตามการตั้งค่า
    6.1.email notification 
	6.2.line notification
	6.3.discord notification
	6.4.telegram notification
	6.5.sms notification
	6.6.web dashboard notification
	6.7.ส่งค่าไป function สั่ง เปิด ปิด อุปกรณ์  หรือ สั่ง ให้ หุ่นยนต์ ทำงาน หรือ หยุดทำงาน 
7.บันทึกผลการส่ง และอัปเดตสถานะ แยก ประเภท 
    7.1.email notification 
	7.2.line notification
	7.3.discord notification
	7.4.telegram notification
	7.5.sms notification
	7.6.web dashboard notification
	7.7.ส่งค่าไป function สั่ง เปิด ปิด อุปกรณ์  หรือ สั่ง ให้ หุ่นยนต์ ทำงาน หรือ หยุดทำงาน 
8.แจ้งเตือนผ่าน WebSocket สำหรับ Real-time
	8.1.web dashboard notification
	8.2.ส่งค่าไป function สั่ง เปิด ปิด อุปกรณ์  หรือ สั่ง ให้ หุ่นยนต์ ทำงาน หรือ หยุดทำงาน 
9.จัดการการแจ้งซ้ำ ตามสถานะและเวลาที่กำหนด
  9.1.นำเวลาแจ่งเตียน ล่าสุด มา เทียเวลาปจุบัน เกิน ค่าที่กำหนดไหม เช่น 10 นาที 
     table prefix sd_
     table sd_notification_type
	   1.Normal  หากยัง มีสถานะ  
	   2.Worming 3.Recovery Worming  
	   3.Alarm    
	   4.Recovery Alarm  
   ให้เจ้งเตียนซ้ำทุก  10 นาที  หาก  Normal  ให้แเจ้งเตียน หยุด  สถานะ Normal แล้วหยุด แเจ้งเตียน
        table sd_notification_type
	   1.Normal  หากยัง มีสถานะ  
	   2.Worming 3.Recovery Worming  
	   3.Alarm    
	   4.Recovery Alarm  
   Notification condition 
   9.2.หาก สถานะ 1.Normal ไม่มีการ  แเจ้งเตียน
   9.3.เก็นประวัติการ แเจ้งเตียน
   9.4.กำหนด ค่าการแจ้งเตียชนแต่ละช่วง  1.Normal  2.Worming 3.Recovery Worming  3.Alarm    4.Recovery Alarm
   
10.เก็บข้อมูลรายงาน สำหรับ Dashboard
  10.1 แยกประเภท เก็บข้อมูลรายงาน
  10.2 ข้อมูลรายงาน สำหรับ สร้างกราฟ
  10.3 ข้อมูลรายงาน สำหรับ ข้อมูล ดิบ
  
 11.notification  group 
  11.1 sensor device เช่น  
  11.2 IO device (Input  optput) เช่น สถานะ การ 1.เปิด  2.ปิด 3.เชื่อมต่ออุปกรณ์ไม่ได้

12.หลัก ตั้งค่า แจ้งเตียน
12.1.สร้าง หัวข้อแจ้างเตียน
12.2 แจ้งเตียนพร้อม  สั่ง ON  OFF Device  ตาม  ฟังก์ชัน AlarmDetailValidate

การ หน่วงการแจ้งเตียน   เช่น  10 นาที
การทำงาน 
1.นำค่าที่ตั้งไว้่ เช่น 10 นาที
2.สร้าง case key ใน Redis ไว้ตรวจสอบ เช่น ในเวลา 10 นาที ที่มีการแจ้งเตียนจะไม่แจ้งเตียนซ้ำ 
3.หาก ค่า ปกติ ลบ case key ใน Redis ออก
4.บนทึกเหตุการลง log


## 1. ขั้นตอนการทำงานของฟังก์ชัน AlarmDetailValidate

### ขั้นตอนที่ 1: รับและเตรียมข้อมูลพื้นฐาน
- รับข้อมูลจาก DTO (Data Transfer Object)
- ดึงค่า `unit` หรือกำหนดเป็น string ว่างถ้าไม่มี
- ดึงค่า `type_id` และแปลงเป็นตัวเลข
- ถ้ามี `alarmTypeId` ให้ใช้ค่าแทน `type_id`

### ขั้นตอนที่ 2: ประมวลผล sensorValues
- ดึงค่า `value_data` จาก DTO
- ถ้าค่าไม่ใช่ null, undefined หรือ string ว่าง:
  - พยายามแปลงค่าเป็นตัวเลข
  - ถ้าแปลงสำเร็จให้ใช้ค่าตัวเลข

### ขั้นตอนที่ 3: กำหนดค่าพื้นฐาน
- แปลงค่าต่างๆ เป็นตัวเลขหรือ string:
  - `status_alert` → `statusAlert` (ค่าเริ่มต้น 0)
  - `status_warning` → `statusWarning` (ค่าเริ่มต้น 0)
  - `recovery_warning` → `recoveryWarning` (ค่าเริ่มต้น 0)
  - `recovery_alert` → `recoveryAlert` (ค่าเริ่มต้น 0)
  - `mqtt_name` → `mqttName` (string ว่าง)
  - `device_name` → `deviceName`
  - `action_name` → `alarmActionName`
  - `mqtt_control_on` → `mqttControlOn`
  - `mqtt_control_off` → `mqttControlOff`
  - `count_alarm` → `count_alarm` (ค่าเริ่มต้น 0)
  - `event` → `event` (ค่าเริ่มต้น 0)
- กำหนดค่าเริ่มต้น:
  - `dataAlarm` = 999
  - `eventControl` = ค่า event
  - `messageMqttControl` = ขึ้นกับค่า event (1 → mqttControlOn, อื่น → mqttControlOff)

### ขั้นตอนที่ 4: กำหนดตัวแปรสำหรับเงื่อนไข
- กำหนดค่าเริ่มต้นสำหรับตัวแปรที่จะใช้ในเงื่อนไข:
  - `alarmStatusSet`: 999
  - `subject`: string ว่าง
  - `content`: string ว่าง
  - `status`: 5
  - `data_alarm`: 0
  - `value_data`: จาก DTO
  - `value_alarm`: จาก DTO หรือ string ว่าง
  - `value_relay`: จาก DTO หรือ string ว่าง
  - `value_control_relay`: จาก DTO หรือ string ว่าง
  - `sensor_data`: null
  - `title`: 'Normal'

### ขั้นตอนที่ 5: ตรวจสอบสถานะ alarm
- แปลง `value_data` เป็นตัวเลข → `sensor_data`
- ใช้ `sensor_data` เป็น `sensorValue`
- **เงื่อนไขการตรวจสอบ:**
  1. **Warning** (สถานะ 1):
     - เมื่อ: `sensorValue ≥ statusWarning` และ `statusWarning < statusAlert`
     - กำหนด: `alarmStatusSet=1`, `title='Warning'`
  
  2. **Alarm** (สถานะ 2):
     - เมื่อ: `sensorValue ≥ statusAlert` และ `statusAlert > statusWarning`
     - กำหนด: `alarmStatusSet=2`, `title='Alarm'`
  
  3. **Recovery Warning** (สถานะ 3):
     - เมื่อ: `count_alarm ≥ 1` และ `sensorValue ≤ recoveryWarning` และ `recoveryWarning ≤ recoveryAlert`
     - กำหนด: `alarmStatusSet=3`, `title='Recovery Warning'`
     - สลับสถานะ eventControl และ messageMqttControl
  
  4. **Recovery Alarm** (สถานะ 4):
     - เมื่อ: `count_alarm ≥ 1` และ `sensorValue ≤ recoveryAlert` และ `recoveryAlert ≥ recoveryWarning`
     - กำหนด: `alarmStatusSet=4`, `title='Recovery Alarm'`
     - สลับสถานะ eventControl และ messageMqttControl
  
  5. **Normal** (สถานะ 5):
     - เมื่อ: ไม่ตรงกับเงื่อนไขใดๆ ข้างต้น
     - กำหนด: `alarmStatusSet=999`, `title='Normal'`

### ขั้นตอนที่ 6: สร้างผลลัพธ์
- รวมค่าทั้งหมดเป็น object ผลลัพธ์
- ประกอบด้วยข้อมูลทั้งหมดที่ประมวลผลได้

### ขั้นตอนที่ 7: จัดการข้อผิดพลาด
- จับข้อผิดพลาดด้วย try-catch
- log ข้อผิดพลาด
- throw ข้อผิดพลาดต่อไป

---

## 2. Workflow Design

```mermaid
graph TD
    A[เริ่มต้น: รับ DTO] --> B[ขั้นตอนที่ 1: เตรียมข้อมูลพื้นฐาน]
    B --> C[ขั้นตอนที่ 2: ประมวลผล sensor values]
    C --> D[ขั้นตอนที่ 3: กำหนดค่าพื้นฐาน]
    D --> E[ขั้นตอนที่ 4: กำหนดตัวแปรเงื่อนไข]
    E --> F[ขั้นตอนที่ 5: ตรวจสอบสถานะ alarm]
    
    F --> G{ตรวจสอบเงื่อนไข}
    
    G -->|sensorValue ≥ statusWarning<br>และ statusWarning < statusAlert| H[สถานะ: Warning]
    G -->|sensorValue ≥ statusAlert<br>และ statusAlert > statusWarning| I[สถานะ: Alarm]
    G -->|count_alarm ≥ 1<br>และ sensorValue ≤ recoveryWarning<br>และ recoveryWarning ≤ recoveryAlert| J[สถานะ: Recovery Warning]
    G -->|count_alarm ≥ 1<br>และ sensorValue ≤ recoveryAlert<br>และ recoveryAlert ≥ recoveryWarning| K[สถานะ: Recovery Alarm]
    G -->|ไม่ตรงเงื่อนไขใดๆ| L[สถานะ: Normal]
    
    H --> M[ขั้นตอนที่ 6: สร้างผลลัพธ์]
    I --> M
    J --> M
    K --> M
    L --> M
    
    M --> N{มีข้อผิดพลาด?}
    N -->|ใช่| O[บันทึกและ throw error]
    N -->|ไม่ใช่| P[ส่งคืนผลลัพธ์]
    
    O --> Q[สิ้นสุดกระบวนการ]
    P --> Q
```

---

## 3. Flowchart อย่างง่าย

```
┌─────────────────────────────────────────────────┐
│            เริ่มต้น: AlarmDetailValidate        │
├─────────────────────────────────────────────────┤
│  1. รับและเตรียมข้อมูลพื้นฐานจาก DTO           │
│     - unit, type_id, alarmTypeId                │
├─────────────────────────────────────────────────┤
│  2. ประมวลผล sensor values                     │
│     - แปลง value_data เป็นตัวเลข (ถ้าเป็นไปได้) │
├─────────────────────────────────────────────────┤
│  3. กำหนดค่าพื้นฐาน                             │
│     - statusAlert, statusWarning, etc.          │
│     - กำหนดค่าเริ่มต้นให้ตัวแปรต่างๆ             │
├─────────────────────────────────────────────────┤
│  4. กำหนดตัวแปรสำหรับเงื่อนไข                   │
│     - alarmStatusSet, subject, content, etc.    │
├─────────────────────────────────────────────────┤
│  5. ตรวจสอบสถานะ alarm                         │
│     ├── IF sensorValue ≥ statusWarning          │
│     │      AND statusWarning < statusAlert      │
│     │   → Warning (สถานะ 1)                     │
│     ├── ELSE IF sensorValue ≥ statusAlert       │
│     │      AND statusAlert > statusWarning      │
│     │   → Alarm (สถานะ 2)                       │
│     ├── ELSE IF count_alarm ≥ 1                 │
│     │      AND sensorValue ≤ recoveryWarning    │
│     │      AND recoveryWarning ≤ recoveryAlert  │
│     │   → Recovery Warning (สถานะ 3)            │
│     ├── ELSE IF count_alarm ≥ 1                 │
│     │      AND sensorValue ≤ recoveryAlert      │
│     │      AND recoveryAlert ≥ recoveryWarning  │
│     │   → Recovery Alarm (สถานะ 4)              │
│     └── ELSE                                    │
│         → Normal (สถานะ 5)                      │
├─────────────────────────────────────────────────┤
│  6. สร้าง object ผลลัพธ์                        │
│     - รวมค่าทั้งหมด                             │
├─────────────────────────────────────────────────┤
│  7. จัดการข้อผิดพลาด (try-catch)                │
└─────────────────────────────────────────────────┘
```

---

## 4. เงื่อนไขสำคัญที่ควรทราบ

### ลำดับความสำคัญของสถานะ:
1. **Alarm** (สูงสุด) - เมื่อค่าเกินระดับแจ้งเตือนสูง
2. **Warning** - เมื่อค่าเกินระดับแจ้งเตือนต่ำ
3. **Recovery Alarm** - เมื่อค่าลดลงต่ำกว่าระดับ recovery สูง
4. **Recovery Warning** - เมื่อค่าลดลงต่ำกว่าระดับ recovery ต่ำ
5. **Normal** - สภาวะปกติ

### เงื่อนไขพิเศษ:
- การ Recovery (สถานะ 3 และ 4) เกิดขึ้นได้เมื่อ `count_alarm ≥ 1` เท่านั้น
- การ Recovery จะสลับสถานะของ `eventControl` และ `messageMqttControl`
- ค่าเริ่มต้นของ `dataAlarm` คือ 999 แสดงถึงสถานะ "ไม่ทราบ"

### ค่าคืนที่สำคัญในผลลัพธ์:
- `status`: สถานะปัจจุบัน (1-5)
- `alarmStatusSet`: รหัสสถานะ alarm
- `title`: หัวข้อสถานะ
- `subject` และ `content`: ข้อความสำหรับการแจ้งเตือน
- `eventControl` และ `messageMqttControl`: สำหรับควบคุมอุปกรณ์
- `dataAlarm` และ `data_alarm`: ค่าขอบเขตของ alarm



step 1
	-สร้าง entities  type orm
	-สร้าง modules rest api  
	-สร้าง level notification
	1.type orm
	2.database postgress sql
	3.socket.io
	4.mqtt
	5.redis
step 2
	1.ออกแบบ ระบบ notification data flow 
	2.ออกแบบ table database postgress sql แยก ทุก ค่าที่จำเป็น
	3.ออกแบบ Quit ใช่้ socket.io / mqtt / redis
	4.ออกแบบ  โคร้งสร้างข้อมูล สำหรับ หน้ารายงาน
	5.ออกแบบ ระบบตั้งค้า notification

step 3 notification  alarm configuration 
   1.ชื่อ Device 
   2.ข้อมูล Device เช่น  humidity  56 %
   3.วันเวลา
   4.สถานะ เช่น  1.Normal  2.Worming 3.Recovery Worming  3.Alarm    4.Recovery Alarm
   5.notification log  เก็บเวลา 
   6.นำเวลาแจ่งเตียน ล่าสุด มา เทียเวลาปจุบัน เกิน ค่าที่กำหนดไหม เช่น 10 นาที 
     table notification_type
	   1.Normal  หากยัง มีสถานะ  
	   2.Worming 3.Recovery Worming  
	   3.Alarm    
	   4.Recovery Alarm  
    7.icon notification
   ให้เจ้งเตียนซ้ำทุก  10 นาที  หาก  Normal  ให้แเจ้งเตียน หยุด  สถานะ Normal แล้วหยุด แเจ้งเตียน
        table notification_type
	   1.Normal  หากยัง มีสถานะ  
	   2.Worming 3.Recovery Worming  
	   3.Alarm    
	   4.Recovery Alarm  
   Notification condition 
   7.หาก สถานะ 1.Normal ไม่มีการ  แเจ้งเตียน
   8.เก็นประวัติการ แเจ้งเตียน
   9.กำหนด ค่าการแจ้งเตียชนแต่ละช่วง  1.Normal  2.Worming 3.Recovery Worming  3.Alarm    4.Recovery Alarm

step 4 ช่องทาง การ ให้เลืกได้มากว่า 1 ช่องทาง โดยมี รูปแบบดังนี้
	1.line notification
	2.discord notification
	3.telegram notification
	4.sms notification
	5.web dashboard notification
	7.ส่งค่าไป function สั่ง เปิด ปิด อุปกรณ์  หรือ สั่ง ให้ หุ่นยนต์ ทำงาน หรือ หยุดทำงาน 
	8.AI CHERT BOT  notification  python  fast api
 

 

ออกแบบรูป monitorring  htnl 5 CSS  ICON 
URL : https://preview.tabler.io
ICON : https://preview.tabler.io/icons.html

การแสดง ข้อมูล 
1.ชื่อ Device sensor
2.ข้อมูล sensor เช่น  humidity  56 %
3.ข้อมูล สถานะ  เช่น  1.Normal  2.Worming 3.Recovery Worming  3.Alarm    4.Recovery Alarm
4.แสดงสัญลักษณ์ icon
5.แสดงสี  ไฟกระพริบ 1.Normal สีเขียว 2.Worming สีส้ม 3.Alarm  สีแดง
6.แสดง history log สถานะ

1.ชื่อ Device IO เช่น  Air Conditioner
2.ข้อมูล Device IO เช่น   Air  ON / OFF
3.ข้อมูล สถานะการทำงาน  Device IO เช่น  ไฟกระพริบ 1.ON สีเขียว 2.OFF  สีแดง
4.แสดงสัญลักษณ์ icon
5.แสดง history log สถานะ


ทำเป็น Rest api เรียกให้ทำงานด้วย  node-red  หรือ  nood js consjob  มี .env configuration
 

