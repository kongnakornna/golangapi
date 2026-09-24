**โมดูลที่ 3: 👥 Customer Management (การจัดการข้อมูลลูกค้าและรถยนต์)**

โมดูลนี้เป็นข้อมูลหลัก (Master Data) ที่สำคัญที่สุดของระบบ เพราะทุกใบงานต้องเชื่อมโยงกับลูกค้าและรถยนต์ ประกอบด้วย 2 ส่วนหลัก:
1. **ข้อมูลลูกค้า (Customer)** – บุคคลหรือนิติบุคคล
2. **ข้อมูลรถยนต์ (Vehicle/Car)** – ยานพาหนะที่เข้ารับบริการ

ผมจะนำเสนอโครงสร้างที่ครบถ้วนตามข้อกำหนดของคุณทุกประการ

---

## 📁 โครงสร้างโมดูล Customer Management (`modules/customer`)

```
modules/customer/
├── application/
│   ├── interfaces/
│   │   ├── CustomerService.java
│   │   └── CarService.java
│   ├── impl/
│   │   ├── CustomerServiceImpl.java
│   │   └── CarServiceImpl.java
│   └── usecase/
│       ├── CreateCustomerUseCase.java
│       ├── UpdateCustomerUseCase.java
│       ├── GetCustomerUseCase.java
│       ├── DeleteCustomerUseCase.java
│       ├── SearchCustomerUseCase.java
│       ├── GetCustomerHistoryUseCase.java
│       ├── CreateCarUseCase.java
│       ├── UpdateCarUseCase.java
│       ├── GetCarUseCase.java
│       ├── DeleteCarUseCase.java
│       ├── GetCarByLicensePlateUseCase.java
│       └── GetCustomerCarsUseCase.java
├── domain/
│   ├── MCustomer.java
│   ├── MCar.java
│   ├── enums/
│   │   ├── CustomerType.java        // INDIVIDUAL, CORPORATE
│   │   ├── CustomerStatus.java      // ACTIVE, INACTIVE, BLACKLISTED
│   │   ├── CarBrand.java
│   │   ├── CarFuelType.java         // GASOLINE, DIESEL, EV, HYBRID, LPG, CNG
│   │   └── CarTransmissionType.java // MANUAL, AUTOMATIC, CVT, DCT
│   └── valueobjects/
│       ├── TaxId.java               // เลขประจำตัวผู้เสียภาษี
│       ├── PhoneNumber.java
│       └── LicensePlate.java
├── infrastructure/
│   ├── repository/
│   │   ├── CustomerRepository.java
│   │   ├── CarRepository.java
│   │   └── impl/
│   │       ├── CustomerRepositoryImpl.java
│   │       └── CarRepositoryImpl.java
│   ├── cache/                              // ⬅️ ระบบ Cache สำหรับ Customer/Car
│   │   ├── CustomerCacheService.java
│   │   └── CarCacheService.java
│   ├── entity/
│   │   ├── CustomerEntity.java
│   │   └── CarEntity.java
│   └── mapper/
│       ├── CustomerMapper.java
│       └── CarMapper.java
└── presentation/
    ├── controller/
    │   ├── CustomerController.java
    │   └── CarController.java
    ├── dto/
    │   ├── request/
    │   │   ├── CustomerCreateRequestDTO.java
    │   │   ├── CustomerUpdateRequestDTO.java
    │   │   ├── CustomerSearchRequestDTO.java
    │   │   ├── CarCreateRequestDTO.java
    │   │   ├── CarUpdateRequestDTO.java
    │   │   └── CarSearchRequestDTO.java
    │   └── response/
    │       ├── CustomerResponseDTO.java
    │       ├── CustomerDetailResponseDTO.java
    │       ├── CarResponseDTO.java
    │       └── CustomerHistoryResponseDTO.java
    └── validator/
        ├── CustomerValidator.java
        └── CarValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Customer

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V3__customer_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_customer (ข้อมูลลูกค้า)
-- Master table for customer information.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_customer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_code VARCHAR(20) UNIQUE NOT NULL,         -- รหัสลูกค้า (เช่น CUST-2026-001)
    full_name VARCHAR(200) NOT NULL,                   -- ชื่อ-นามสกุล หรือชื่อบริษัท
    display_name VARCHAR(200),                         -- ชื่อที่ใช้แสดง
    customer_type VARCHAR(20) NOT NULL DEFAULT 'INDIVIDUAL', -- INDIVIDUAL, CORPORATE
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',      -- ACTIVE, INACTIVE, BLACKLISTED
    tax_id VARCHAR(20),                                -- เลขประจำตัวผู้เสียภาษี (TIN)
    email VARCHAR(100),
    phone_number VARCHAR(20) NOT NULL,
    secondary_phone VARCHAR(20),
    address TEXT,
    province VARCHAR(100),
    city VARCHAR(100),
    district VARCHAR(100),
    postal_code VARCHAR(10),
    country VARCHAR(50) DEFAULT 'Thailand',
    contact_person VARCHAR(100),                       -- ผู้ติดต่อ (กรณีเป็นนิติบุคคล)
    contact_phone VARCHAR(20),
    notes TEXT,
    last_visit_date TIMESTAMP,                         -- วันที่มาใช้บริการล่าสุด
    total_visit_count INTEGER DEFAULT 0,               -- จำนวนครั้งที่มาใช้บริการทั้งหมด
    total_spent DECIMAL(15,2) DEFAULT 0,               -- ยอดเงินที่ใช้บริการทั้งหมด
    -- Audit Fields (from GenericBusinessEntity)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_customer_phone ON m_customer(phone_number);
CREATE INDEX idx_m_customer_email ON m_customer(email);
CREATE INDEX idx_m_customer_tax_id ON m_customer(tax_id);
CREATE INDEX idx_m_customer_code ON m_customer(customer_code);
CREATE INDEX idx_m_customer_status ON m_customer(status);
CREATE INDEX idx_m_customer_whitelabel ON m_customer(whitelabel_id);
CREATE INDEX idx_m_customer_deleted ON m_customer(deleted);

-- ==============================================
-- ฟังก์ชันสร้างรหัสลูกค้าอัตโนมัติ (Auto-generate Customer Code)
-- Function to generate unique customer code with year and sequence.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_customer_code()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SPLIT_PART(customer_code, '-', 3) AS INTEGER)), 0) + 1
        FROM m_customer
        WHERE customer_code LIKE 'CUST-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.customer_code := 'CUST-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_customer_code ON m_customer;
CREATE TRIGGER trg_generate_customer_code
BEFORE INSERT ON m_customer
FOR EACH ROW
EXECUTE FUNCTION generate_customer_code();

-- ==============================================
-- ตาราง: m_car (ข้อมูลรถยนต์)
-- Master table for vehicle information.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_car (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES m_customer(id) ON DELETE CASCADE,
    license_plate VARCHAR(20) UNIQUE NOT NULL,        -- ทะเบียนรถ
    province VARCHAR(50),                              -- จังหวัดที่จดทะเบียน
    brand VARCHAR(50) NOT NULL,                       -- ยี่ห้อ (Toyota, Honda, etc.)
    model VARCHAR(100) NOT NULL,                      -- รุ่น (Camry, Accord, etc.)
    sub_model VARCHAR(100),                           -- รุ่นย่อย (2.0 Hybrid, etc.)
    year INTEGER,                                     -- ปีที่ผลิต
    color VARCHAR(30),
    engine_number VARCHAR(50),
    chassis_number VARCHAR(50),                       -- เลขตัวถัง (VIN)
    fuel_type VARCHAR(20),                            -- GASOLINE, DIESEL, EV, HYBRID, LPG, CNG
    transmission_type VARCHAR(20),                    -- MANUAL, AUTOMATIC, CVT, DCT
    engine_cc INTEGER,                                -- ขนาดเครื่องยนต์ (cc)
    seating_capacity INTEGER,                         -- จำนวนที่นั่ง
    mileage INTEGER DEFAULT 0,                        -- ระยะทางสะสม (กม.)
    last_service_date TIMESTAMP,                      -- วันที่เข้ารับบริการล่าสุด
    next_service_mileage INTEGER,                     -- ระยะทางที่ต้องเข้ารับบริการครั้งถัดไป
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    -- Audit Fields
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_car_customer ON m_car(customer_id);
CREATE INDEX idx_m_car_license_plate ON m_car(license_plate);
CREATE INDEX idx_m_car_brand ON m_car(brand);
CREATE INDEX idx_m_car_whitelabel ON m_car(whitelabel_id);
CREATE INDEX idx_m_car_deleted ON m_car(deleted);

-- ==============================================
-- ตาราง: m_car_service_history (ประวัติการซ่อมบำรุงของรถ)
-- Vehicle service history summary (cached/aggregated data for quick display).
-- ==============================================
CREATE TABLE IF NOT EXISTS m_car_service_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    car_id UUID NOT NULL REFERENCES m_car(id) ON DELETE CASCADE,
    job_id UUID NOT NULL REFERENCES t_job(id) ON DELETE CASCADE,
    service_date TIMESTAMP NOT NULL,
    service_type VARCHAR(50),                         -- ชนิดงาน (OIL_CHANGE, BRAKE, etc.)
    description TEXT,
    total_cost DECIMAL(15,2),
    mileage_at_service INTEGER,
    mechanic_name VARCHAR(100),
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_car_service_history_car ON m_car_service_history(car_id);
CREATE INDEX idx_m_car_service_history_date ON m_car_service_history(service_date);
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Customer Management

### `infrastructure/cache/CustomerCacheService.java`

```java
package com.template.app.modules.customer.infrastructure.cache;

import com.template.app.modules.customer.domain.MCustomer;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class CustomerCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลลูกค้าจาก Cache หากมีอยู่จะไม่เรียกฐานข้อมูล (เพิ่มประสิทธิภาพการค้นหาบ่อย)
        This function retrieves customer data from cache. If present, it avoids DB hits (improves frequent lookups).
        Redis Key: customer:{id}
    */
    @Cacheable(value = "customers", key = "#customerId")
    public MCustomer getCustomer(UUID customerId) {
        return null; // Spring จะใช้ Cache หากมี / Spring uses cache if available.
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลลูกค้าจากหมายเลขโทรศัพท์ (Cache เฉพาะฟิลด์ที่ค้นหาบ่อย)
        This function retrieves customer data by phone number (caches frequently searched fields).
        Redis Key: customer_phone:{phone}
    */
    @Cacheable(value = "customer_phone", key = "#phone")
    public MCustomer getCustomerByPhone(String phone) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกข้อมูลลูกค้า
        This function updates the cache when customer data is saved.
    */
    @CachePut(value = "customers", key = "#customer.id")
    public MCustomer saveCustomer(MCustomer customer) {
        return customer;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลลูกค้าออกจาก Cache ทั้งสองรูปแบบ (ID และ Phone)
        This function evicts customer data from both cache forms (ID and Phone).
    */
    @CacheEvict(value = {"customers", "customer_phone"}, key = "#customerId")
    public void evictCustomer(UUID customerId) {
        // ลบ Cache ทั้งสองรูปแบบ / Evict both cache entries.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของลูกค้า (ใช้เมื่อระบบอัปเดตข้อมูลจำนวนมาก)
        This function clears all customer caches (used during bulk updates).
    */
    @CacheEvict(value = {"customers", "customer_phone"}, allEntries = true)
    public void evictAllCustomers() {
        // ลบทุก key ใน customers และ customer_phone / Evict all keys in both caches.
    }
}
```

### `infrastructure/cache/CarCacheService.java`

```java
package com.template.app.modules.customer.infrastructure.cache;

import com.template.app.modules.customer.domain.MCar;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class CarCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลรถยนต์จาก Cache (ลดภาระฐานข้อมูลเมื่อค้นหาทะเบียนบ่อยๆ)
        This function retrieves vehicle data from cache (reduces DB load on frequent plate lookups).
    */
    @Cacheable(value = "cars", key = "#carId")
    public MCar getCar(UUID carId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลรถยนต์จากทะเบียน (Cache สำหรับการค้นหาทะเบียน)
        This function retrieves vehicle data by license plate (cached for quick plate searches).
    */
    @Cacheable(value = "car_plate", key = "#licensePlate")
    public MCar getCarByPlate(String licensePlate) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อบันทึกรถยนต์
        This function updates the cache when a car is saved.
    */
    @CachePut(value = "cars", key = "#car.id")
    public MCar saveCar(MCar car) {
        return car;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลรถยนต์ออกจาก Cache ทั้งสองรูปแบบ (ID และ ทะเบียน)
        This function evicts car data from both cache forms (ID and license plate).
    */
    @CacheEvict(value = {"cars", "car_plate"}, key = "#carId")
    public void evictCar(UUID carId) {
        // ลบ Cache / Evict cache.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Customer Controller

```java
package com.template.app.modules.customer.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.customer.application.interfaces.CustomerService;
import com.template.app.modules.customer.presentation.dto.request.CustomerCreateRequestDTO;
import com.template.app.modules.customer.presentation.dto.request.CustomerSearchRequestDTO;
import com.template.app.modules.customer.presentation.dto.request.CustomerUpdateRequestDTO;
import com.template.app.modules.customer.presentation.dto.response.CustomerDetailResponseDTO;
import com.template.app.modules.customer.presentation.dto.response.CustomerResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/customers")
@Tag(name = "Customer", description = "Customer and Vehicle Management APIs")
@RequiredArgsConstructor
public class CustomerController {

    private final CustomerService customerService;

    // ========================================================================
    // 1. CREATE CUSTOMER
    // ========================================================================

    /*
        API: POST /api/v1/customers
        ฟังก์ชันนี้สร้างข้อมูลลูกค้าใหม่ ทั้งบุคคลธรรมดาและนิติบุคคล
        This function creates a new customer (individual or corporate).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที (ป้องกันการสร้างข้อมูลซ้ำ)
        Rate Limit: Allows 20 requests per minute (prevent duplicate entries).
    */
    @PostMapping
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Create a new customer")
    public ResponseEntity<CustomerResponseDTO> createCustomer(@Valid @RequestBody CustomerCreateRequestDTO request)
            throws SystemGlobalException {
        CustomerResponseDTO response = customerService.createCustomer(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET CUSTOMER BY ID
    // ========================================================================

    /*
        API: GET /api/v1/customers/{id}
        ฟังก์ชันนี้ดึงข้อมูลลูกค้าตาม ID (ใช้ Cache ช่วยลดภาระ DB)
        This function retrieves customer by ID (uses caching to reduce DB load).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get customer by ID")
    public ResponseEntity<CustomerDetailResponseDTO> getCustomer(@PathVariable UUID id)
            throws SystemGlobalException {
        CustomerDetailResponseDTO response = customerService.getCustomer(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. SEARCH CUSTOMERS (Pagination + Filters)
    // ========================================================================

    /*
        API: POST /api/v1/customers/search
        ฟังก์ชันนี้ค้นหาลูกค้าด้วยตัวกรองหลายรูปแบบ เช่น ชื่อ, เบอร์โทร, ประเภท, สถานะ
        This function searches customers with multiple filters: name, phone, type, status.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @PostMapping("/search")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Search customers with filters")
    public ResponseEntity<Page<CustomerResponseDTO>> searchCustomers(
            @Valid @RequestBody CustomerSearchRequestDTO request,
            Pageable pageable) throws SystemGlobalException {
        Page<CustomerResponseDTO> page = customerService.searchCustomers(request, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 4. UPDATE CUSTOMER
    // ========================================================================

    /*
        API: PUT /api/v1/customers/{id}
        ฟังก์ชันนี้แก้ไขข้อมูลลูกค้า (ชื่อ, ที่อยู่, เบอร์โทร, สถานะ)
        This function updates customer details (name, address, phone, status).
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @PutMapping("/{id}")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Update customer details")
    public ResponseEntity<CustomerResponseDTO> updateCustomer(
            @PathVariable UUID id,
            @Valid @RequestBody CustomerUpdateRequestDTO request) throws SystemGlobalException {
        CustomerResponseDTO response = customerService.updateCustomer(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. DELETE CUSTOMER (SOFT DELETE)
    // ========================================================================

    /*
        API: DELETE /api/v1/customers/{id}
        ฟังก์ชันนี้ลบลูกค้าแบบ Soft Delete (จะไม่หายจากระบบ แต่ถูกซ่อน)
        This function soft-deletes a customer (hidden from views but retained in DB).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete customer (soft delete)")
    public ResponseEntity<Void> deleteCustomer(@PathVariable UUID id) throws SystemGlobalException {
        customerService.deleteCustomer(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 6. GET CUSTOMER SERVICE HISTORY
    // ========================================================================

    /*
        API: GET /api/v1/customers/{id}/history
        ฟังก์ชันนี้ดึงประวัติการซ่อมบำรุงของลูกค้า (แสดงรายการ Job Card ทั้งหมดของลูกค้า)
        This function retrieves the service history of a customer (all Job Cards).
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/{id}/history")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get customer service history")
    public ResponseEntity<List<CustomerHistoryResponseDTO>> getCustomerHistory(@PathVariable UUID id)
            throws SystemGlobalException {
        List<CustomerHistoryResponseDTO> history = customerService.getCustomerHistory(id);
        return ResponseEntity.ok(history);
    }

    // ========================================================================
    // 7. GET CUSTOMER BY PHONE (Quick Lookup)
    // ========================================================================

    /*
        API: GET /api/v1/customers/phone/{phone}
        ฟังก์ชันนี้ค้นหาลูกค้าด้วยเบอร์โทรศัพท์ (ใช้ Cache ช่วยค้นหาเร็วมาก)
        This function finds a customer by phone number (uses cache for ultra-fast lookup).
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที (ใช้บ่อยเวลารับรถ)
        Rate Limit: Allows 60 requests per minute (frequently used during vehicle intake).
    */
    @GetMapping("/phone/{phone}")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Find customer by phone number")
    public ResponseEntity<CustomerResponseDTO> getCustomerByPhone(@PathVariable String phone)
            throws SystemGlobalException {
        CustomerResponseDTO response = customerService.getCustomerByPhone(phone);
        return ResponseEntity.ok(response);
    }
}
```

### `CarController.java`

```java
package com.template.app.modules.customer.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.customer.application.interfaces.CarService;
import com.template.app.modules.customer.presentation.dto.request.CarCreateRequestDTO;
import com.template.app.modules.customer.presentation.dto.request.CarSearchRequestDTO;
import com.template.app.modules.customer.presentation.dto.request.CarUpdateRequestDTO;
import com.template.app.modules.customer.presentation.dto.response.CarResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/cars")
@Tag(name = "Vehicle", description = "Vehicle Management APIs")
@RequiredArgsConstructor
public class CarController {

    private final CarService carService;

    // ========================================================================
    // 1. CREATE CAR
    // ========================================================================

    /*
        API: POST /api/v1/cars
        ฟังก์ชันนี้เพิ่มรถยนต์ใหม่ให้กับลูกค้า โดยต้องระบุ Customer ID
        This function adds a new vehicle for a customer (requires Customer ID).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @PostMapping
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Add a new vehicle")
    public ResponseEntity<CarResponseDTO> createCar(@Valid @RequestBody CarCreateRequestDTO request)
            throws SystemGlobalException {
        CarResponseDTO response = carService.createCar(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET CAR BY ID
    // ========================================================================

    /*
        API: GET /api/v1/cars/{id}
        ฟังก์ชันนี้ดึงข้อมูลรถยนต์ตาม ID
        This function retrieves vehicle details by ID.
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get vehicle by ID")
    public ResponseEntity<CarResponseDTO> getCar(@PathVariable UUID id)
            throws SystemGlobalException {
        CarResponseDTO response = carService.getCar(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET CARS BY CUSTOMER
    // ========================================================================

    /*
        API: GET /api/v1/cars/customer/{customerId}
        ฟังก์ชันนี้ดึงรายการรถยนต์ทั้งหมดของลูกค้า (ใช้เวลาเปิด Job Card)
        This function lists all vehicles owned by a customer (used when opening a Job Card).
        Rate Limit: อนุญาต 80 ครั้งต่อ 1 นาที
        Rate Limit: Allows 80 requests per minute.
    */
    @GetMapping("/customer/{customerId}")
    @RateLimit(limit = 80, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get all vehicles by customer")
    public ResponseEntity<List<CarResponseDTO>> getCarsByCustomer(@PathVariable UUID customerId)
            throws SystemGlobalException {
        List<CarResponseDTO> cars = carService.getCarsByCustomer(customerId);
        return ResponseEntity.ok(cars);
    }

    // ========================================================================
    // 4. GET CAR BY LICENSE PLATE
    // ========================================================================

    /*
        API: GET /api/v1/cars/plate/{licensePlate}
        ฟังก์ชันนี้ค้นหารถยนต์ด้วยทะเบียนรถ (ใช้บ่อยในการรับรถเข้าอู่)
        This function finds a vehicle by license plate (frequently used during vehicle intake).
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/plate/{licensePlate}")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Find vehicle by license plate")
    public ResponseEntity<CarResponseDTO> getCarByPlate(@PathVariable String licensePlate)
            throws SystemGlobalException {
        CarResponseDTO response = carService.getCarByPlate(licensePlate);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. SEARCH CARS
    // ========================================================================

    /*
        API: POST /api/v1/cars/search
        ฟังก์ชันนี้ค้นหารถยนต์ด้วยตัวกรอง เช่น ยี่ห้อ, รุ่น, ปี, ประเภทเกียร์
        This function searches vehicles with filters: brand, model, year, transmission.
        Rate Limit: อนุญาต 40 ครั้งต่อ 1 นาที
        Rate Limit: Allows 40 requests per minute.
    */
    @PostMapping("/search")
    @RateLimit(limit = 40, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Search vehicles with filters")
    public ResponseEntity<Page<CarResponseDTO>> searchCars(
            @Valid @RequestBody CarSearchRequestDTO request,
            Pageable pageable) throws SystemGlobalException {
        Page<CarResponseDTO> page = carService.searchCars(request, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 6. UPDATE CAR
    // ========================================================================

    /*
        API: PUT /api/v1/cars/{id}
        ฟังก์ชันนี้แก้ไขข้อมูลรถยนต์ (ทะเบียน, ยี่ห้อ, รุ่น, ปี)
        This function updates vehicle details (plate, brand, model, year).
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @PutMapping("/{id}")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Update vehicle details")
    public ResponseEntity<CarResponseDTO> updateCar(
            @PathVariable UUID id,
            @Valid @RequestBody CarUpdateRequestDTO request) throws SystemGlobalException {
        CarResponseDTO response = carService.updateCar(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 7. DELETE CAR (SOFT DELETE)
    // ========================================================================

    /*
        API: DELETE /api/v1/cars/{id}
        ฟังก์ชันนี้ลบรถยนต์แบบ Soft Delete
        This function soft-deletes a vehicle.
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete vehicle (soft delete)")
    public ResponseEntity<Void> deleteCar(@PathVariable UUID id) throws SystemGlobalException {
        carService.deleteCar(id);
        return ResponseEntity.noContent().build();
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/CustomerType.java`

```java
package com.template.app.modules.customer.domain.enums;

/*
    ประเภทของลูกค้า / Customer type.
*/
public enum CustomerType {
    INDIVIDUAL,   // บุคคลธรรมดา / Individual person.
    CORPORATE     // นิติบุคคล / Corporate entity.
}
```

### `domain/enums/CustomerStatus.java`

```java
package com.template.app.modules.customer.domain.enums;

/*
    สถานะของลูกค้า / Customer status.
*/
public enum CustomerStatus {
    ACTIVE,       // ใช้งานได้ / Active.
    INACTIVE,     // ใช้งานไม่ได้ (พัก) / Inactive.
    BLACKLISTED   // ถูกแบน / Blacklisted (prohibited from service).
}
```

### `domain/MCustomer.java`

```java
package com.template.app.modules.customer.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.customer.domain.enums.CustomerStatus;
import com.template.app.modules.customer.domain.enums.CustomerType;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@EqualsAndHashCode(callSuper = true)
public class MCustomer extends GenericBusinessClass {

    private String customerCode;          // รหัสลูกค้า / Customer code.
    private String fullName;              // ชื่อ-นามสกุล / Full name.
    private String displayName;           // ชื่อที่ใช้แสดง / Display name.
    private CustomerType customerType;    // ประเภท / Type.
    private CustomerStatus status;        // สถานะ / Status.
    private String taxId;                 // เลขประจำตัวผู้เสียภาษี / Tax ID.
    private String email;                 // อีเมล / Email.
    private String phoneNumber;           // เบอร์โทรหลัก / Primary phone.
    private String secondaryPhone;        // เบอร์โทรสำรอง / Secondary phone.
    private String address;               // ที่อยู่ / Address.
    private String province;              // จังหวัด / Province.
    private String city;                  // อำเภอ/เขต / City/District.
    private String district;              // ตำบล/แขวง / Sub-district.
    private String postalCode;            // รหัสไปรษณีย์ / Postal code.
    private String country;               // ประเทศ / Country.
    private String contactPerson;         // ผู้ติดต่อ / Contact person.
    private String contactPhone;          // เบอร์ติดต่อ / Contact phone.
    private String notes;                 // หมายเหตุ / Notes.
    private LocalDateTime lastVisitDate;  // วันที่มาใช้บริการล่าสุด / Last visit date.
    private Integer totalVisitCount;      // จำนวนครั้งที่มาใช้บริการ / Total visits.
    private BigDecimal totalSpent;        // ยอดเงินที่ใช้บริการทั้งหมด / Total spent.

    /*
        ฟังก์ชันนี้ใช้บันทึกการมาใช้บริการครั้งใหม่ (อัปเดต lastVisitDate และเพิ่ม visit count)
        This function records a new visit (updates lastVisitDate and increments visit count).
    */
    public void recordVisit(BigDecimal spentAmount) {
        this.lastVisitDate = LocalDateTime.now();
        this.totalVisitCount = (this.totalVisitCount != null ? this.totalVisitCount : 0) + 1;
        this.totalSpent = (this.totalSpent != null ? this.totalSpent : BigDecimal.ZERO).add(spentAmount);
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าลูกค้าสามารถใช้บริการได้หรือไม่
        This function checks if the customer is eligible for service.
    */
    public boolean canReceiveService() {
        return this.status == CustomerStatus.ACTIVE;
    }

    /*
        ฟังก์ชันนี้เปลี่ยนสถานะลูกค้าเป็น BLACKLISTED (ห้ามใช้บริการ)
        This function changes the customer status to BLACKLISTED (service ban).
    */
    public void blacklist(String reason) {
        this.status = CustomerStatus.BLACKLISTED;
        this.notes = (this.notes != null ? this.notes + "\n" : "") + "Blacklisted: " + reason;
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/CustomerServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.customer.application.impl;

import com.template.app.modules.customer.application.interfaces.CustomerService;
import com.template.app.modules.customer.domain.MCustomer;
import com.template.app.modules.customer.infrastructure.cache.CustomerCacheService;
import com.template.app.modules.customer.infrastructure.repository.CustomerRepository;
import com.template.app.modules.customer.presentation.dto.request.CustomerCreateRequestDTO;
import com.template.app.modules.customer.presentation.dto.response.CustomerResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.stereotype.Service;

@Service
public class CustomerServiceImpl extends GenericServiceImpl<MCustomer, CustomerRepository> implements CustomerService {

    private final CustomerCacheService cacheService;

    public CustomerServiceImpl(CustomerRepository repository, CustomerCacheService cacheService) {
        super(repository);
        this.cacheService = cacheService;
    }

    /*
        ฟังก์ชันนี้สร้างลูกค้าใหม่ ตรวจสอบความซ้ำของเบอร์โทรและอีเมล จากนั้นบันทึกและอัปเดต Cache
        This function creates a new customer, checks for duplicate phone/email, saves, and updates cache.
    */
    @Override
    public CustomerResponseDTO createCustomer(CustomerCreateRequestDTO request) throws SystemGlobalException {
        // 1. ตรวจสอบว่ามีลูกค้าที่มีเบอร์โทรนี้แล้วหรือไม่ / Check if phone already exists.
        if (repository.existsByPhone(request.getPhoneNumber())) {
            throw new SystemGlobalException("Phone number already registered.", null);
        }

        // 2. สร้าง Domain Entity และกำหนดค่า / Create domain entity and set values.
        MCustomer customer = new MCustomer();
        customer.setFullName(request.getFullName());
        customer.setPhoneNumber(request.getPhoneNumber());
        customer.setEmail(request.getEmail());
        customer.setCustomerType(request.getCustomerType());
        customer.setStatus(CustomerStatus.ACTIVE);
        customer.setAddress(request.getAddress());
        customer.setProvince(request.getProvince());
        customer.setCity(request.getCity());

        // 3. บันทึกผ่าน Generic Repository (จะสร้าง customerCode อัตโนมัติผ่าน Trigger) / Save via Generic Repository.
        MCustomer savedCustomer = this.create(customer);

        // 4. บันทึกลง Cache เพื่อให้ค้นหาเร็ว / Store in cache for fast lookups.
        cacheService.saveCustomer(savedCustomer);

        // 5. แปลงเป็น DTO และคืนค่า / Convert to DTO and return.
        return CustomerResponseDTO.fromEntity(savedCustomer);
    }

    /*
        ฟังก์ชันนี้ค้นหาลูกค้าด้วยเบอร์โทร โดยใช้ Cache เป็นหลัก (ถ้ามี)
        This function finds a customer by phone, primarily using Cache (if available).
    */
    @Override
    public CustomerResponseDTO getCustomerByPhone(String phone) throws SystemGlobalException {
        // 1. ลองค้นหาจาก Cache ก่อน / Try to get from cache first.
        MCustomer customer = cacheService.getCustomerByPhone(phone);
        if (customer == null) {
            // 2. ถ้าไม่พบใน Cache ให้ค้นหาจากฐานข้อมูล / If not in cache, fetch from DB.
            customer = repository.findByPhone(phone, getRepositoryAuth())
                    .orElseThrow(() -> new SystemGlobalException("Customer not found with phone: " + phone, null));
            // 3. บันทึกลง Cache เพื่อใช้ครั้งต่อไป / Store in cache for next time.
            cacheService.saveCustomer(customer);
        }
        return CustomerResponseDTO.fromEntity(customer);
    }
}
```

---

## ✅ สรุปส่วนที่เพิ่มเติมในโมดูล Customer Management

| ฟีเจอร์ | สถานะ | คำอธิบาย |
|---------|--------|----------|
| **Database Design** | ✅ ครบถ้วน | `m_customer`, `m_car`, `m_car_service_history` พร้อม Trigger และ Index |
| **Cache (Redis)** | ✅ เพิ่มแล้ว | `CustomerCacheService` (ID + Phone), `CarCacheService` (ID + Plate) |
| **Rate Limit** | ✅ เพิ่มแล้ว | Customer Controller และ Car Controller ทุก Endpoint |
| **API Routing** | ✅ ชัดเจน | CRUD + Search + History + Phone Lookup + Plate Lookup |
| **Bilingual Comments** | ✅ ครบถ้วน | ทุก Function มีคำอธิบายภาษาไทยและภาษาอังกฤษแยกบรรทัด |

---

 