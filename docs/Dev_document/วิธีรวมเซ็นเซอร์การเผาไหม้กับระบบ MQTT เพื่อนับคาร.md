<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# วิธีรวมเซ็นเซอร์การเผาไหม้กับระบบ MQTT เพื่อนับคาร์บอน

การรวมเซนเซอร์การเผาไหม้ (Combustion Sensors) เข้ากับระบบ MQTT เพื่อคำนวณคาร์บอน (Scope 1) เป็นงานที่ต้องแม่นยำและปลอดภัยสูง (เพราะเกี่ยวกับเชื้อเพลิงและความร้อน) นี่คือแนวทางปฏิบัติทีละขั้นตอนครับ

***

### 1. เลือกและเชื่อมต่อ Hardware

เซนเซอร์กลุ่มนี้มักเป็นเกรดอุตสาหกรรม (Industrial Grade) ไม่ค่อยส่ง MQTT ได้โดยตรง ต้องมี **Gateway** คั่นกลาง

#### A. Fuel Flow Meter (วัดปริมาณเชื้อเพลิง) - *วิธีที่แม่นยำที่สุดสำหรับ Carbon Footprint*

* **Output:** มักเป็น Pulse (Digital Pulse) หรือ 4-20mA หรือ RS-485 (Modbus RTU).
* **การเชื่อมต่อ:**
    * ใช้ **Industrial IoT Gateway** (เช่น Teltonika, Advantech หรือ ESP32 เกรดอุตสาหกรรม) อ่านค่า Pulse/Modbus.
    * แปลงค่าสะสม (Totalizer) เป็นหน่วย Liters หรือ $m^3$.


#### B. Flue Gas Analyzer / CEMS (วัดก๊าซไอเสีย) - *วิธีตรวจสอบประสิทธิภาพ*

* **Output:** RS-485 (Modbus RTU) หรือ 4-20mA.
* **ค่าที่ต้องดึง:** % $CO_2$, % $O_2$, Temp ($^\circ C$), ppm $CO$.
* **การเชื่อมต่อ:** ใช้ Gateway อ่าน Modbus Register ตามคู่มือของผู้ผลิต (เช่น Testo, SICK, Horiba).

***

### 2. การออกแบบ MQTT Payload (Edge Calculation)

เพื่อให้ Server คำนวณคาร์บอนได้ง่าย ควร "Pre-process" ข้อมูลบางส่วนที่ Edge Gateway ก่อนส่ง

**Topic:** `factory/{site_id}/boiler/{device_id}/combustion`

**Payload Example (JSON):**

```json
{
  "ts": "2025-11-25T15:30:00Z",
  "fuel_type": "diesel",  // ระบุชนิดเชื้อเพลิงเพื่อเลือก Factor ถูก
  "flow_rate": {
    "value": 15.5,
    "unit": "liters/hour"
  },
  "total_consumption": {
    "value": 15040.0,
    "unit": "liters",
    "reset_ts": "2025-01-01T00:00:00Z" // จุดอ้างอิงค่าสะสม
  },
  "exhaust_gas": {
    "co2_percent": 12.5,
    "o2_percent": 4.2,
    "temp_c": 180.5
  },
  "status": {
    "burner_on": true,
    "efficiency": 88.5 // คำนวณที่ Edge ได้ถ้า Gateway ฉลาดพอ
  }
}
```


***

### 3. สูตรคำนวณคาร์บอน (Server-Side Logic)

เมื่อ MQTT Broker ได้รับข้อมูล ให้ Backend (Node-RED, Python, หรือ Server) นำไปคำนวณ

#### สูตร 1: คำนวณจากปริมาณเชื้อเพลิง (แม่นยำที่สุดสำหรับ Scope 1 Report)

ใช้ค่า `total_consumption` (Liters) × Emission Factor

\$ Carbon (kgCO_2e) = Fuel Used (Liters) \times Emission Factor \$

* *Diesel Factor:* $\approx 2.6-2.7 \text{ kg}CO_2/\text{L}$
* *LPG Factor:* $\approx 1.5-1.6 \text{ kg}CO_2/\text{L}$ (หรือตาม kg)
* *Natural Gas:* คำนวณตาม $m^3$ หรือ MMBTU


#### สูตร 2: ตรวจสอบจากความเข้มข้นไอเสีย ($Mass Balance Check$)

ใช้ค่า `co2_percent` จาก CEMS × อัตราการไหลของลมไอเสีย (Exhaust Flow Rate)
*วิธีนี้ยากกว่าเพราะต้องรู้ Air Flow ที่แน่นอน แต่นิยมใช้ในโรงไฟฟ้าขนาดใหญ่เพื่อ Cross-check*

***

### 4. แนวทางการ Integration จริง (Architecture Blueprint)

1. **Edge Layer (หน้างาน):**
    * ติดตั้ง **Flow Meter** ตัดท่อน้ำมันเข้า Burner.
    * ติดตั้ง **Modbus Gateway** (เช่น ESP32 + MAX485 หรือ RUT240) อ่านค่า Flow (Pulse) และ CEMS (Modbus).
    * Gateway แปลงข้อมูลเป็น JSON แล้ว Publish ขึ้น MQTT Broker.
2. **Communication Layer:**
    * ใช้ MQTT (TLS Secured) ส่งผ่าน 4G/LAN.
    * Topic: `factory/bangkadi/boiler-01/data`
    * QoS: 1 (ห้ามข้อมูลหาย เพราะเกี่ยวกับยอดสะสมเชื้อเพลิง).
3. **AoT/Data Layer:**
    * **InfluxDB/TimescaleDB:** เก็บ Time-series ของ Flow Rate และ % $CO_2$ เพื่อดูประสิทธิภาพการเผาไหม้ (Efficiency Trend).
    * **SQL Database:** เก็บสรุปยอดรายวัน (Daily Consumption) เพื่อคูณ Factor ออก Report รายเดือน.
4. **Dashboard \& Alert:**
    * แสดง **Real-time Carbon Emission Rate** (kg $CO_2$/hour).
    * แจ้งเตือนเมื่อ **Efficiency ตก** (เช่น $O_2$ สูงเกินไป แปลว่าเผาไหม้ไม่สมบูรณ์/ลมเกิน → เปลืองน้ำมัน → คาร์บอนพุ่ง).

### ข้อควรระวัง (Safety \& Reliability)

* **Intrinsically Safe:** หากติดตั้งเซนเซอร์ในโซนแก๊ส/น้ำมันไวไฟ อุปกรณ์ต้องกันระเบิด (Ex-proof) หรือติดตั้งนอกโซนอันตราย
* **Bypass Counter:** กรณี Gateway ล่ม ต้องจดมิเตอร์ Flow Meter ด้วยตาเปล่าได้ (Mechanical Register) เพื่อไม่ให้ข้อมูลคาร์บอนหายไปจากระบบบัญชีครับ
<span style="display:none">[^1][^2][^3]</span>

<div align="center">⁂</div>

[^1]: https://www.sciencedirect.com/science/article/pii/S2405844025006887

[^2]: https://www.automation.co.th/products/factory-automation/environmental-monitoring/continuous-emission-monitoring-system-cems/

[^3]: https://envea.global/design/pdf/envea_catalogue_cems_emission-monitoring_en.pdf

