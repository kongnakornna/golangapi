 **โมดูลที่ 6: 📦 Inventory Management (การจัดการสินค้าคงคลัง)**

โมดูล Inventory Management เป็นหัวใจสำคัญของการจัดการอะไหล่และสินค้าคงคลังทั้งหมดในระบบ ครอบคลุมการทำงานดังนี้:

1. **การจัดการสินค้า (Part Master)** – เพิ่ม/แก้ไข/ลบ/ค้นหา อะไหล่
2. **การเคลื่อนไหวสินค้า (Inventory Transaction)** – รับเข้า, เบิกจ่าย, ปรับปรุงสต็อก
3. **การเบิกอะไหล่ (Part Picking)** – เบิกอะไหล่จากคลังเพื่อใช้ใน Job Card
4. **การตรวจนับสต็อก (Stock Take)** – นับสินค้าจริงเพื่อเปรียบเทียบกับระบบ
5. **การปรับปรุงสต็อก (Stock Adjustment)** – ปรับจำนวนสินค้าให้ตรงกับความเป็นจริง
6. **การจัดการตำแหน่งจัดเก็บ (Stock Location)** – จัดเก็บสินค้าในตำแหน่งต่างๆ
7. **การแจ้งเตือนสินค้าต่ำกว่าเกณฑ์ (Low Stock Alert)** – เมื่อสินค้าใกล้หมด

---

## 📁 โครงสร้างโมดูล Inventory Management (`modules/inventory`)

```
modules/inventory/
├── application/
│   ├── interfaces/
│   │   ├── InventoryService.java
│   │   ├── PartMasterService.java
│   │   ├── PartPickingService.java
│   │   ├── StockAdjustmentService.java
│   │   ├── StockTakeService.java
│   │   ├── StockLocationService.java
│   │   └── StockSummaryService.java
│   ├── impl/
│   │   ├── InventoryServiceImpl.java
│   │   ├── PartMasterServiceImpl.java
│   │   ├── PartPickingServiceImpl.java
│   │   ├── StockAdjustmentServiceImpl.java
│   │   ├── StockTakeServiceImpl.java
│   │   ├── StockLocationServiceImpl.java
│   │   └── StockSummaryServiceImpl.java
│   └── usecase/
│       ├── CreatePartUseCase.java
│       ├── UpdatePartUseCase.java
│       ├── GetPartUseCase.java
│       ├── DeletePartUseCase.java
│       ├── SearchPartUseCase.java
│       ├── ReceiveInventoryUseCase.java
│       ├── IssueInventoryUseCase.java
│       ├── CreatePartPickingUseCase.java
│       ├── ConfirmPartPickingUseCase.java
│       ├── CreateStockAdjustmentUseCase.java
│       ├── CreateStockTakeUseCase.java
│       ├── CompleteStockTakeUseCase.java
│       ├── CheckLowStockUseCase.java
│       └── GetStockSummaryUseCase.java
├── domain/
│   ├── MPartMaster.java
│   ├── MStockLocation.java
│   ├── TInventory.java
│   ├── TInventoryAdjustmentHeader.java
│   ├── TInventoryAdjustmentDetail.java
│   ├── TPartPickingRequest.java
│   ├── TPartPickingDetail.java
│   ├── TStockTakeHeader.java
│   ├── TStockTakeDetail.java
│   ├── enums/
│   │   ├── InventoryTransactionType.java    // RECEIVE, ISSUE, ADJUSTMENT, RETURN
│   │   ├── StockAdjustmentType.java         // INCREASE, DECREASE
│   │   ├── PickingStatus.java               // DRAFT, PENDING, PICKED, CONFIRMED, CANCELLED
│   │   └── StockTakeStatus.java             // DRAFT, IN_PROGRESS, COMPLETED, CANCELLED
│   └── valueobjects/
│       ├── PartNumber.java
│       ├── StockQuantity.java
│       └── ReorderLevel.java
├── infrastructure/
│   ├── repository/
│   │   ├── PartMasterRepository.java
│   │   ├── StockLocationRepository.java
│   │   ├── InventoryRepository.java
│   │   ├── StockAdjustmentRepository.java
│   │   ├── PartPickingRequestRepository.java
│   │   ├── PartPickingDetailRepository.java
│   │   ├── StockTakeRepository.java
│   │   └── impl/
│   │       ├── PartMasterRepositoryImpl.java
│   │       ├── InventoryRepositoryImpl.java
│   │       ├── PartPickingRequestRepositoryImpl.java
│   │       └── StockTakeRepositoryImpl.java
│   ├── cache/                                           // ⬅️ ระบบ Cache สำหรับ Inventory
│   │   ├── PartMasterCacheService.java
│   │   ├── InventoryCacheService.java
│   │   └── StockSummaryCacheService.java
│   ├── report/                                          // ⬅️ ระบบสร้างรายงาน
│   │   ├── PartPickingReportGenerator.java
│   │   └── InventoryReportGenerator.java
│   ├── entity/
│   │   ├── PartMasterEntity.java
│   │   ├── StockLocationEntity.java
│   │   ├── InventoryEntity.java
│   │   ├── InventoryAdjustmentHeaderEntity.java
│   │   ├── InventoryAdjustmentDetailEntity.java
│   │   ├── PartPickingRequestEntity.java
│   │   ├── PartPickingDetailEntity.java
│   │   ├── StockTakeHeaderEntity.java
│   │   └── StockTakeDetailEntity.java
│   └── mapper/
│       ├── PartMasterMapper.java
│       ├── StockLocationMapper.java
│       ├── InventoryMapper.java
│       ├── PartPickingMapper.java
│       └── StockTakeMapper.java
└── presentation/
    ├── controller/
    │   ├── PartMasterController.java        // CRUD Part
    │   ├── InventoryController.java         // Receive, Issue, Stock Summary
    │   ├── PartPickingController.java       // Create, Confirm Picking
    │   ├── StockAdjustmentController.java   // Adjust Stock
    │   ├── StockTakeController.java         // Stock Take
    │   └── StockLocationController.java     // Location Management
    ├── dto/
    │   ├── request/
    │   │   ├── PartMasterCreateRequestDTO.java
    │   │   ├── PartMasterUpdateRequestDTO.java
    │   │   ├── PartMasterSearchRequestDTO.java
    │   │   ├── InventoryReceiveRequestDTO.java
    │   │   ├── InventoryIssueRequestDTO.java
    │   │   ├── PartPickingCreateRequestDTO.java
    │   │   ├── PartPickingConfirmRequestDTO.java
    │   │   ├── StockAdjustmentRequestDTO.java
    │   │   └── StockTakeRequestDTO.java
    │   └── response/
    │       ├── PartMasterResponseDTO.java
    │       ├── InventoryResponseDTO.java
    │       ├── InventorySummaryDTO.java
    │       ├── PartPickingResponseDTO.java
    │       ├── StockAdjustmentResponseDTO.java
    │       └── StockTakeResponseDTO.java
    └── validator/
        ├── PartMasterValidator.java
        ├── InventoryValidator.java
        └── PartPickingValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Inventory

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V6__inventory_schema.sql`)

```sql
-- ==============================================
-- ตาราง: m_part_master (ข้อมูลอะไหล่หลัก)
-- Master table for parts and spare parts.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_part_master (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    part_code VARCHAR(50) UNIQUE NOT NULL,          -- รหัสอะไหล่ (Part Code)
    part_name VARCHAR(200) NOT NULL,                -- ชื่ออะไหล่
    part_name_en VARCHAR(200),                      -- ชื่ออะไหล่ (ภาษาอังกฤษ)
    category_id UUID REFERENCES m_category(id),     -- หมวดหมู่
    brand VARCHAR(50),                              -- ยี่ห้อ / Brand
    model VARCHAR(100),                             -- รุ่น / Model
    oem_number VARCHAR(50),                         -- เลข OEM (Original Equipment Manufacturer)
    description TEXT,                               -- รายละเอียด
    unit VARCHAR(20) DEFAULT 'PIECE',               -- หน่วย: PIECE, SET, BOX, LITER, KG
    reorder_level INTEGER DEFAULT 0,                -- ระดับที่ต้องสั่งซื้อ (Reorder Level)
    reorder_quantity INTEGER DEFAULT 0,             -- จำนวนที่ต้องสั่งซื้อเมื่อถึง Reorder Level
    stock_quantity INTEGER DEFAULT 0,               -- จำนวนในสต็อกปัจจุบัน
    min_stock INTEGER DEFAULT 0,                    -- จำนวนขั้นต่ำที่ต้องมี
    max_stock INTEGER DEFAULT 0,                    -- จำนวนสูงสุดที่ควรมี
    unit_cost DECIMAL(15,2),                        -- ต้นทุนต่อหน่วย
    selling_price DECIMAL(15,2),                    -- ราคาขาย
    location_id UUID REFERENCES m_stock_location(id), -- ตำแหน่งจัดเก็บ
    status VARCHAR(20) DEFAULT 'ACTIVE',            -- ACTIVE, INACTIVE, DISCONTINUED
    image_url TEXT,                                 -- URL รูปภาพ
    notes TEXT,                                     -- หมายเหตุ
    last_updated_stock TIMESTAMP,                   -- วันที่อัปเดตสต็อกครั้งล่าสุด
    -- Audit Fields
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_part_master_code ON m_part_master(part_code);
CREATE INDEX idx_m_part_master_name ON m_part_master(part_name);
CREATE INDEX idx_m_part_master_category ON m_part_master(category_id);
CREATE INDEX idx_m_part_master_brand ON m_part_master(brand);
CREATE INDEX idx_m_part_master_oem ON m_part_master(oem_number);
CREATE INDEX idx_m_part_master_location ON m_part_master(location_id);
CREATE INDEX idx_m_part_master_status ON m_part_master(status);
CREATE INDEX idx_m_part_master_whitelabel ON m_part_master(whitelabel_id);
CREATE INDEX idx_m_part_master_deleted ON m_part_master(deleted);

-- ==============================================
-- ตาราง: m_stock_location (ตำแหน่งจัดเก็บสินค้า)
-- Stock location management.
-- ==============================================
CREATE TABLE IF NOT EXISTS m_stock_location (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_code VARCHAR(20) UNIQUE NOT NULL,      -- รหัสตำแหน่ง (เช่น A-01, B-03)
    location_name VARCHAR(100) NOT NULL,            -- ชื่อตำแหน่ง
    location_type VARCHAR(20) DEFAULT 'SHELF',      -- SHELF, RACK, WAREHOUSE, BAY
    zone VARCHAR(50),                               -- โซน (Zone)
    capacity INTEGER,                               -- ความจุสูงสุด
    current_usage INTEGER DEFAULT 0,                -- จำนวนที่ใช้อยู่
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_m_stock_location_code ON m_stock_location(location_code);
CREATE INDEX idx_m_stock_location_zone ON m_stock_location(zone);

-- ==============================================
-- ตาราง: t_inventory (รายการเคลื่อนไหวสินค้าคงคลัง)
-- Inventory transaction (movement) log.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    part_id UUID NOT NULL REFERENCES m_part_master(id) ON DELETE RESTRICT,
    transaction_type VARCHAR(20) NOT NULL,          -- RECEIVE, ISSUE, ADJUSTMENT, RETURN
    reference_type VARCHAR(30),                     -- PO, JOB, ADJUSTMENT, STOCK_TAKE
    reference_id UUID,                              -- Reference ID (PO ID, Job ID, etc.)
    quantity INTEGER NOT NULL,                      -- จำนวนที่เคลื่อนไหว (บวก = เพิ่ม, ลบ = ลด)
    previous_quantity INTEGER NOT NULL,             -- จำนวนก่อนการเคลื่อนไหว
    new_quantity INTEGER NOT NULL,                  -- จำนวนหลังการเคลื่อนไหว
    unit_cost DECIMAL(15,2),                        -- ต้นทุนต่อหน่วย ณ เวลาที่เคลื่อนไหว
    total_cost DECIMAL(15,2),                       -- ต้นทุนรวม
    transaction_date TIMESTAMP NOT NULL DEFAULT NOW(),
    note TEXT,                                      -- หมายเหตุ
    performed_by UUID NOT NULL REFERENCES m_user(id), -- ผู้ทำรายการ
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_inventory_part ON t_inventory(part_id);
CREATE INDEX idx_t_inventory_type ON t_inventory(transaction_type);
CREATE INDEX idx_t_inventory_reference ON t_inventory(reference_type, reference_id);
CREATE INDEX idx_t_inventory_date ON t_inventory(transaction_date);
CREATE INDEX idx_t_inventory_whitelabel ON t_inventory(whitelabel_id);

-- ==============================================
-- ตาราง: t_inventory_adjustment_header (หัวการปรับปรุงสต็อก)
-- Stock adjustment header.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_inventory_adjustment_header (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    adjustment_no VARCHAR(20) UNIQUE NOT NULL,      -- เลขที่ปรับปรุง (ADJ-2026-0001)
    adjustment_date TIMESTAMP NOT NULL DEFAULT NOW(),
    adjustment_type VARCHAR(20) NOT NULL,           -- INCREASE, DECREASE
    reason VARCHAR(50) NOT NULL,                    -- DAMAGE, LOST, RETURN, CORRECTION, OTHER
    status VARCHAR(20) DEFAULT 'DRAFT',             -- DRAFT, APPROVED, CANCELLED
    description TEXT,                               -- รายละเอียด
    approved_by UUID REFERENCES m_user(id),
    approved_at TIMESTAMP,
    total_adjustment_value DECIMAL(15,2),
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_inv_adj_header_no ON t_inventory_adjustment_header(adjustment_no);
CREATE INDEX idx_t_inv_adj_header_status ON t_inventory_adjustment_header(status);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่ปรับปรุงสต็อกอัตโนมัติ
-- Function to generate adjustment number.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_adjustment_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(adjustment_no FROM 9) AS INTEGER)), 0) + 1
        FROM t_inventory_adjustment_header
        WHERE adjustment_no LIKE 'ADJ-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.adjustment_no := 'ADJ-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_adjustment_no ON t_inventory_adjustment_header;
CREATE TRIGGER trg_generate_adjustment_no
BEFORE INSERT ON t_inventory_adjustment_header
FOR EACH ROW
EXECUTE FUNCTION generate_adjustment_no();

-- ==============================================
-- ตาราง: t_inventory_adjustment_detail (รายละเอียดการปรับปรุงสต็อก)
-- Stock adjustment details.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_inventory_adjustment_detail (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    adjustment_header_id UUID NOT NULL REFERENCES t_inventory_adjustment_header(id) ON DELETE CASCADE,
    part_id UUID NOT NULL REFERENCES m_part_master(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL,
    unit_cost DECIMAL(15,2),
    total_cost DECIMAL(15,2),
    note TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_inv_adj_detail_header ON t_inventory_adjustment_detail(adjustment_header_id);
CREATE INDEX idx_t_inv_adj_detail_part ON t_inventory_adjustment_detail(part_id);

-- ==============================================
-- ตาราง: t_part_picking_request (คำขอเบิกอะไหล่)
-- Part picking request.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_part_picking_request (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    picking_no VARCHAR(20) UNIQUE NOT NULL,         -- เลขที่คำขอเบิก (PK-2026-0001)
    job_id UUID REFERENCES t_job(id) ON DELETE RESTRICT,
    quotation_id UUID REFERENCES t_quotation(id) ON DELETE RESTRICT,
    requested_date TIMESTAMP NOT NULL DEFAULT NOW(),
    requested_by UUID NOT NULL REFERENCES m_user(id),
    status VARCHAR(20) DEFAULT 'DRAFT',             -- DRAFT, PENDING, PICKED, CONFIRMED, CANCELLED
    priority VARCHAR(20) DEFAULT 'NORMAL',          -- NORMAL, URGENT
    notes TEXT,
    picked_by UUID REFERENCES m_user(id),
    picked_date TIMESTAMP,
    confirmed_by UUID REFERENCES m_user(id),
    confirmed_date TIMESTAMP,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_picking_request_job ON t_part_picking_request(job_id);
CREATE INDEX idx_t_picking_request_quotation ON t_part_picking_request(quotation_id);
CREATE INDEX idx_t_picking_request_status ON t_part_picking_request(status);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่คำขอเบิกอัตโนมัติ
-- Function to generate picking number.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_picking_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(picking_no FROM 8) AS INTEGER)), 0) + 1
        FROM t_part_picking_request
        WHERE picking_no LIKE 'PK-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.picking_no := 'PK-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_picking_no ON t_part_picking_request;
CREATE TRIGGER trg_generate_picking_no
BEFORE INSERT ON t_part_picking_request
FOR EACH ROW
EXECUTE FUNCTION generate_picking_no();

-- ==============================================
-- ตาราง: t_part_picking_detail (รายละเอียดคำขอเบิก)
-- Part picking details.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_part_picking_detail (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    picking_request_id UUID NOT NULL REFERENCES t_part_picking_request(id) ON DELETE CASCADE,
    part_id UUID NOT NULL REFERENCES m_part_master(id) ON DELETE RESTRICT,
    requested_quantity INTEGER NOT NULL,
    picked_quantity INTEGER DEFAULT 0,
    unit_price DECIMAL(15,2),
    total_price DECIMAL(15,2),
    note TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_picking_detail_request ON t_part_picking_detail(picking_request_id);
CREATE INDEX idx_t_picking_detail_part ON t_part_picking_detail(part_id);

-- ==============================================
-- ตาราง: t_stocktake_header (หัวการตรวจนับสต็อก)
-- Stock take header.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_stocktake_header (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stocktake_no VARCHAR(20) UNIQUE NOT NULL,       -- เลขที่ตรวจนับ (ST-2026-0001)
    stocktake_date TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'DRAFT',             -- DRAFT, IN_PROGRESS, COMPLETED, CANCELLED
    started_by UUID REFERENCES m_user(id),
    started_at TIMESTAMP,
    completed_by UUID REFERENCES m_user(id),
    completed_at TIMESTAMP,
    total_discrepancy INTEGER DEFAULT 0,            -- ผลต่างรวม
    notes TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_stocktake_header_status ON t_stocktake_header(status);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่ตรวจนับอัตโนมัติ
-- Function to generate stocktake number.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_stocktake_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(stocktake_no FROM 8) AS INTEGER)), 0) + 1
        FROM t_stocktake_header
        WHERE stocktake_no LIKE 'ST-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.stocktake_no := 'ST-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_stocktake_no ON t_stocktake_header;
CREATE TRIGGER trg_generate_stocktake_no
BEFORE INSERT ON t_stocktake_header
FOR EACH ROW
EXECUTE FUNCTION generate_stocktake_no();

-- ==============================================
-- ตาราง: t_stocktake_detail (รายละเอียดการตรวจนับ)
-- Stock take details.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_stocktake_detail (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stocktake_header_id UUID NOT NULL REFERENCES t_stocktake_header(id) ON DELETE CASCADE,
    part_id UUID NOT NULL REFERENCES m_part_master(id) ON DELETE RESTRICT,
    system_quantity INTEGER NOT NULL,               -- จำนวนในระบบ
    counted_quantity INTEGER NOT NULL,              -- จำนวนที่นับได้จริง
    discrepancy INTEGER GENERATED ALWAYS AS (counted_quantity - system_quantity) STORED,
    note TEXT,
    counted_by UUID REFERENCES m_user(id),
    counted_at TIMESTAMP,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_stocktake_detail_header ON t_stocktake_detail(stocktake_header_id);
CREATE INDEX idx_t_stocktake_detail_part ON t_stocktake_detail(part_id);
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Inventory

### `infrastructure/cache/PartMasterCacheService.java`

```java
package com.template.app.modules.inventory.infrastructure.cache;

import com.template.app.modules.inventory.domain.MPartMaster;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class PartMasterCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลอะไหล่จาก Cache (ลดภาระฐานข้อมูลเมื่อค้นหาบ่อยๆ)
        This function retrieves part data from cache (reduces DB load on frequent searches).
        Redis Key: part:{id}
    */
    @Cacheable(value = "parts", key = "#partId")
    public MPartMaster getPart(UUID partId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลอะไหล่จากรหัสสินค้า (Part Code)
        This function retrieves part data by part code.
        Redis Key: part_code:{partCode}
    */
    @Cacheable(value = "part_code", key = "#partCode")
    public MPartMaster getPartByCode(String partCode) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกอะไหล่
        This function updates the cache when a part is saved.
    */
    @CachePut(value = "parts", key = "#part.id")
    public MPartMaster savePart(MPartMaster part) {
        return part;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลอะไหล่ออกจาก Cache ทั้งสองรูปแบบ (ID และ Part Code)
        This function evicts part data from both cache forms (ID and Part Code).
    */
    @CacheEvict(value = {"parts", "part_code"}, key = "#partId")
    public void evictPart(UUID partId) {
        // ลบ Cache ทั้งสองรูปแบบ / Evict both cache entries.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของอะไหล่ (ใช้เมื่อระบบอัปเดตข้อมูลจำนวนมาก)
        This function clears all part caches (used during bulk updates).
    */
    @CacheEvict(value = {"parts", "part_code"}, allEntries = true)
    public void evictAllParts() {
        // ลบทุก key ใน caches / Evict all keys.
    }
}
```

### `infrastructure/cache/InventoryCacheService.java`

```java
package com.template.app.modules.inventory.infrastructure.cache;

import com.template.app.modules.inventory.presentation.dto.response.InventorySummaryDTO;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class InventoryCacheService {

    /*
        ฟังก์ชันนี้ดึงสรุปสต็อกของอะไหล่จาก Cache (ใช้แสดงใน Dashboard)
        This function retrieves part stock summary from cache (used in Dashboard display).
        Redis Key: stock_summary:{partId}
    */
    @Cacheable(value = "stock_summary", key = "#partId")
    public InventorySummaryDTO getStockSummary(UUID partId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ลบสรุปสต็อกของอะไหล่ออกจาก Cache เมื่อมีการเคลื่อนไหวสต็อก
        This function evicts part stock summary when stock changes.
    */
    @CacheEvict(value = "stock_summary", key = "#partId")
    public void evictStockSummary(UUID partId) {
        // ลบ Cache / Evict cache.
    }

    /*
        ฟังก์ชันนี้ลบสรุปสต็อกทั้งหมด (ใช้เมื่อมีการปรับปรุงสต็อกจำนวนมาก)
        This function clears all stock summaries (used during bulk stock updates).
    */
    @CacheEvict(value = "stock_summary", allEntries = true)
    public void evictAllStockSummaries() {
        // ลบทุก key ใน stock_summary / Evict all keys.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Inventory Controllers

### `PartMasterController.java`

```java
package com.template.app.modules.inventory.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.inventory.application.interfaces.PartMasterService;
import com.template.app.modules.inventory.presentation.dto.request.PartMasterCreateRequestDTO;
import com.template.app.modules.inventory.presentation.dto.request.PartMasterSearchRequestDTO;
import com.template.app.modules.inventory.presentation.dto.request.PartMasterUpdateRequestDTO;
import com.template.app.modules.inventory.presentation.dto.response.PartMasterResponseDTO;
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
@RequestMapping("/api/v1/parts")
@Tag(name = "Part Master", description = "Part Master Data Management APIs")
@RequiredArgsConstructor
public class PartMasterController {

    private final PartMasterService partMasterService;

    // ========================================================================
    // 1. CREATE PART
    // ========================================================================

    /*
        API: POST /api/v1/parts
        ฟังก์ชันนี้เพิ่มอะไหล่ใหม่ในระบบ พร้อมกำหนดราคาและตำแหน่งจัดเก็บ
        This function adds a new part to the system, including price and location.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PostMapping
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Create a new part")
    public ResponseEntity<PartMasterResponseDTO> createPart(@Valid @RequestBody PartMasterCreateRequestDTO request)
            throws SystemGlobalException {
        PartMasterResponseDTO response = partMasterService.createPart(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET PART BY ID
    // ========================================================================

    /*
        API: GET /api/v1/parts/{id}
        ฟังก์ชันนี้ดึงข้อมูลอะไหล่ตาม ID (ใช้ Cache ช่วยลดภาระ DB)
        This function retrieves part by ID (uses caching to reduce DB load).
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get part by ID")
    public ResponseEntity<PartMasterResponseDTO> getPart(@PathVariable UUID id)
            throws SystemGlobalException {
        PartMasterResponseDTO response = partMasterService.getPart(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET PART BY CODE
    // ========================================================================

    /*
        API: GET /api/v1/parts/code/{partCode}
        ฟังก์ชันนี้ค้นหาอะไหล่ด้วยรหัสอะไหล่ (ใช้บ่อยในการค้นหาอะไหล่)
        This function finds a part by part code (frequently used for part lookups).
        Rate Limit: อนุญาต 80 ครั้งต่อ 1 นาที
        Rate Limit: Allows 80 requests per minute.
    */
    @GetMapping("/code/{partCode}")
    @RateLimit(limit = 80, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Find part by part code")
    public ResponseEntity<PartMasterResponseDTO> getPartByCode(@PathVariable String partCode)
            throws SystemGlobalException {
        PartMasterResponseDTO response = partMasterService.getPartByCode(partCode);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. SEARCH PARTS (Pagination + Filters)
    // ========================================================================

    /*
        API: POST /api/v1/parts/search
        ฟังก์ชันนี้ค้นหาอะไหล่ด้วยตัวกรอง เช่น ชื่อ, หมวดหมู่, ยี่ห้อ, สถานะ
        This function searches parts with filters: name, category, brand, status.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @PostMapping("/search")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Search parts with filters")
    public ResponseEntity<Page<PartMasterResponseDTO>> searchParts(
            @Valid @RequestBody PartMasterSearchRequestDTO request,
            Pageable pageable) throws SystemGlobalException {
        Page<PartMasterResponseDTO> page = partMasterService.searchParts(request, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 5. UPDATE PART
    // ========================================================================

    /*
        API: PUT /api/v1/parts/{id}
        ฟังก์ชันนี้แก้ไขข้อมูลอะไหล่ (ราคา, ชื่อ, ตำแหน่งจัดเก็บ)
        This function updates part details (price, name, location).
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @PutMapping("/{id}")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Update part details")
    public ResponseEntity<PartMasterResponseDTO> updatePart(
            @PathVariable UUID id,
            @Valid @RequestBody PartMasterUpdateRequestDTO request) throws SystemGlobalException {
        PartMasterResponseDTO response = partMasterService.updatePart(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 6. DELETE PART (SOFT DELETE)
    // ========================================================================

    /*
        API: DELETE /api/v1/parts/{id}
        ฟังก์ชันนี้ลบอะไหล่แบบ Soft Delete (จะไม่หายจากระบบ)
        This function soft-deletes a part (retained in system).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete part (soft delete)")
    public ResponseEntity<Void> deletePart(@PathVariable UUID id) throws SystemGlobalException {
        partMasterService.deletePart(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 7. GET LOW STOCK PARTS
    // ========================================================================

    /*
        API: GET /api/v1/parts/low-stock
        ฟังก์ชันนี้แสดงรายการอะไหล่ที่ต่ำกว่า Reorder Level (เพื่อแจ้งเตือนการสั่งซื้อ)
        This function lists parts below reorder level (for purchase alerts).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/low-stock")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get low stock parts")
    public ResponseEntity<List<PartMasterResponseDTO>> getLowStockParts() throws SystemGlobalException {
        List<PartMasterResponseDTO> lowStockParts = partMasterService.getLowStockParts();
        return ResponseEntity.ok(lowStockParts);
    }
}
```

### `PartPickingController.java`

```java
package com.template.app.modules.inventory.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.inventory.application.interfaces.PartPickingService;
import com.template.app.modules.inventory.presentation.dto.request.PartPickingConfirmRequestDTO;
import com.template.app.modules.inventory.presentation.dto.request.PartPickingCreateRequestDTO;
import com.template.app.modules.inventory.presentation.dto.response.PartPickingResponseDTO;
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
@RequestMapping("/api/v1/part-picking")
@Tag(name = "Part Picking", description = "Part Picking Request Management APIs")
@RequiredArgsConstructor
public class PartPickingController {

    private final PartPickingService partPickingService;

    // ========================================================================
    // 1. CREATE PICKING REQUEST
    // ========================================================================

    /*
        API: POST /api/v1/part-picking
        ฟังก์ชันนี้สร้างคำขอเบิกอะไหล่จาก Job Card หรือ Quotation
        This function creates a part picking request from a Job Card or Quotation.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @PostMapping
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Create part picking request")
    public ResponseEntity<PartPickingResponseDTO> createPickingRequest(
            @Valid @RequestBody PartPickingCreateRequestDTO request) throws SystemGlobalException {
        PartPickingResponseDTO response = partPickingService.createPickingRequest(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET PICKING REQUEST BY ID
    // ========================================================================

    /*
        API: GET /api/v1/part-picking/{id}
        ฟังก์ชันนี้ดึงข้อมูลคำขอเบิกตาม ID พร้อมรายการอะไหล่ที่ต้องการเบิก
        This function retrieves picking request by ID with all parts.
        Rate Limit: อนุญาต 80 ครั้งต่อ 1 นาที
        Rate Limit: Allows 80 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 80, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get picking request by ID")
    public ResponseEntity<PartPickingResponseDTO> getPickingRequest(@PathVariable UUID id)
            throws SystemGlobalException {
        PartPickingResponseDTO response = partPickingService.getPickingRequest(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET PICKING BY JOB ID
    // ========================================================================

    /*
        API: GET /api/v1/part-picking/job/{jobId}
        ฟังก์ชันนี้ดึงคำขอเบิกตาม Job ID (ใช้เชื่อมโยง Job -> Picking)
        This function retrieves picking request by Job ID (links Job -> Picking).
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/job/{jobId}")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get picking request by Job ID")
    public ResponseEntity<PartPickingResponseDTO> getPickingByJobId(@PathVariable UUID jobId)
            throws SystemGlobalException {
        PartPickingResponseDTO response = partPickingService.getPickingByJobId(jobId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. LIST PICKING REQUESTS (Pagination)
    // ========================================================================

    /*
        API: GET /api/v1/part-picking
        ฟังก์ชันนี้แสดงรายการคำขอเบิกแบบแบ่งหน้า
        This function lists picking requests with pagination.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List picking requests")
    public ResponseEntity<Page<PartPickingResponseDTO>> listPickingRequests(
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate,
            Pageable pageable) throws SystemGlobalException {
        Page<PartPickingResponseDTO> page = partPickingService.listPickingRequests(status, startDate, endDate, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 5. CONFIRM PICKING
    // ========================================================================

    /*
        API: PUT /api/v1/part-picking/{id}/confirm
        ฟังก์ชันนี้ยืนยันการเบิกอะไหล่ เมื่อพนักงานคลังเบิกสินค้าเรียบร้อยแล้ว
        This function confirms part picking when warehouse staff has picked the goods.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PutMapping("/{id}/confirm")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Confirm part picking")
    public ResponseEntity<PartPickingResponseDTO> confirmPicking(
            @PathVariable UUID id,
            @Valid @RequestBody PartPickingConfirmRequestDTO request) throws SystemGlobalException {
        PartPickingResponseDTO response = partPickingService.confirmPicking(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 6. CANCEL PICKING
    // ========================================================================

    /*
        API: PUT /api/v1/part-picking/{id}/cancel
        ฟังก์ชันนี้ยกเลิกคำขอเบิก พร้อมระบุเหตุผล
        This function cancels a picking request with a reason.
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @PutMapping("/{id}/cancel")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Cancel picking request")
    public ResponseEntity<PartPickingResponseDTO> cancelPicking(
            @PathVariable UUID id,
            @RequestParam String reason) throws SystemGlobalException {
        PartPickingResponseDTO response = partPickingService.cancelPicking(id, reason);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 7. GENERATE PICKING PDF
    // ========================================================================

    /*
        API: GET /api/v1/part-picking/{id}/pdf
        ฟังก์ชันนี้สร้าง PDF เอกสารเบิกอะไหล่สำหรับพนักงานคลัง
        This function generates a PDF picking document for warehouse staff.
        Rate Limit: อนุญาต 15 ครั้งต่อ 5 นาที
        Rate Limit: Allows 15 requests per 5 minutes.
    */
    @GetMapping("/{id}/pdf")
    @RateLimit(limit = 15, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Generate picking PDF")
    public ResponseEntity<byte[]> generatePickingPDF(@PathVariable UUID id) throws SystemGlobalException {
        byte[] pdf = partPickingService.generatePickingPDF(id);
        return ResponseEntity.ok()
                .header("Content-Type", "application/pdf")
                .header("Content-Disposition", "inline; filename=picking_" + id + ".pdf")
                .body(pdf);
    }
}
```

### `InventoryController.java` (รับเข้า / เบิกจ่าย)

```java
package com.template.app.modules.inventory.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.inventory.application.interfaces.InventoryService;
import com.template.app.modules.inventory.presentation.dto.request.InventoryIssueRequestDTO;
import com.template.app.modules.inventory.presentation.dto.request.InventoryReceiveRequestDTO;
import com.template.app.modules.inventory.presentation.dto.response.InventoryResponseDTO;
import com.template.app.modules.inventory.presentation.dto.response.InventorySummaryDTO;
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
@RequestMapping("/api/v1/inventory")
@Tag(name = "Inventory", description = "Inventory Transaction Management APIs")
@RequiredArgsConstructor
public class InventoryController {

    private final InventoryService inventoryService;

    // ========================================================================
    // 1. RECEIVE GOODS (รับสินค้าเข้า)
    // ========================================================================

    /*
        API: POST /api/v1/inventory/receive
        ฟังก์ชันนี้รับสินค้าเข้า Inventory (จาก Purchase Order หรือจาก Supplier โดยตรง)
        This function receives goods into inventory (from Purchase Order or directly from supplier).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PostMapping("/receive")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Receive goods into inventory")
    public ResponseEntity<InventoryResponseDTO> receiveGoods(
            @Valid @RequestBody InventoryReceiveRequestDTO request) throws SystemGlobalException {
        InventoryResponseDTO response = inventoryService.receiveGoods(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. ISSUE GOODS (เบิกจ่ายสินค้า)
    // ========================================================================

    /*
        API: POST /api/v1/inventory/issue
        ฟังก์ชันนี้เบิกจ่ายสินค้าออกจาก Inventory (ใช้ใน Job Card หรือ Part Picking)
        This function issues goods from inventory (used in Job Card or Part Picking).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @PostMapping("/issue")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Issue goods from inventory")
    public ResponseEntity<InventoryResponseDTO> issueGoods(
            @Valid @RequestBody InventoryIssueRequestDTO request) throws SystemGlobalException {
        InventoryResponseDTO response = inventoryService.issueGoods(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET INVENTORY BY PART ID
    // ========================================================================

    /*
        API: GET /api/v1/inventory/part/{partId}
        ฟังก์ชันนี้ดึงข้อมูล Inventory ของอะไหล่ (จำนวน, ต้นทุน, ประวัติการเคลื่อนไหว)
        This function retrieves inventory data for a part (quantity, cost, movement history).
        Rate Limit: อนุญาต 80 ครั้งต่อ 1 นาที
        Rate Limit: Allows 80 requests per minute.
    */
    @GetMapping("/part/{partId}")
    @RateLimit(limit = 80, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get inventory by part ID")
    public ResponseEntity<InventorySummaryDTO> getInventoryByPart(@PathVariable UUID partId)
            throws SystemGlobalException {
        InventorySummaryDTO response = inventoryService.getInventoryByPart(partId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. LIST INVENTORY TRANSACTIONS
    // ========================================================================

    /*
        API: GET /api/v1/inventory/transactions
        ฟังก์ชันนี้แสดงประวัติการเคลื่อนไหวสินค้าทั้งหมดแบบแบ่งหน้า
        This function lists all inventory transaction history with pagination.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/transactions")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List inventory transactions")
    public ResponseEntity<Page<InventoryResponseDTO>> listTransactions(
            @RequestParam(required = false) UUID partId,
            @RequestParam(required = false) String transactionType,
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate,
            Pageable pageable) throws SystemGlobalException {
        Page<InventoryResponseDTO> page = inventoryService.listTransactions(
                partId, transactionType, startDate, endDate, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 5. GET INVENTORY SUMMARY (Dashboard)
    // ========================================================================

    /*
        API: GET /api/v1/inventory/summary
        ฟังก์ชันนี้แสดงสรุปสินค้าคงคลังทั้งหมด (จำนวนสินค้า, มูลค่ารวม, สินค้าใกล้หมด)
        This function shows total inventory summary (total items, total value, low stock).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @GetMapping("/summary")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get inventory summary")
    public ResponseEntity<InventorySummaryDTO> getInventorySummary() throws SystemGlobalException {
        InventorySummaryDTO response = inventoryService.getInventorySummary();
        return ResponseEntity.ok(response);
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/InventoryTransactionType.java`

```java
package com.template.app.modules.inventory.domain.enums;

/*
    ประเภทของการเคลื่อนไหวสินค้า / Inventory transaction type.
*/
public enum InventoryTransactionType {
    RECEIVE,        // รับสินค้าเข้า / Goods received.
    ISSUE,          // เบิกจ่ายสินค้า / Goods issued.
    ADJUSTMENT,     // ปรับปรุงสต็อก / Stock adjustment.
    RETURN,         // คืนสินค้า / Goods returned.
    STOCK_TAKE      // ตรวจนับสต็อก / Stock take.
}
```

### `domain/enums/PickingStatus.java`

```java
package com.template.app.modules.inventory.domain.enums;

/*
    สถานะของคำขอเบิกอะไหล่ / Part picking status.
*/
public enum PickingStatus {
    DRAFT,        // ร่าง / Draft.
    PENDING,      // รอการดำเนินการ / Pending.
    PICKED,       // เบิกสินค้าแล้ว / Picked.
    CONFIRMED,    // ยืนยันแล้ว / Confirmed.
    CANCELLED     // ยกเลิก / Cancelled.
}
```

### `domain/MPartMaster.java`

```java
package com.template.app.modules.inventory.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class MPartMaster extends GenericBusinessClass {

    private String partCode;                // รหัสอะไหล่ / Part code.
    private String partName;                // ชื่ออะไหล่ / Part name.
    private String partNameEn;              // ชื่ออะไหล่ (อังกฤษ) / Part name (English).
    private UUID categoryId;                // ID หมวดหมู่ / Category ID.
    private String brand;                   // ยี่ห้อ / Brand.
    private String model;                   // รุ่น / Model.
    private String oemNumber;               // เลข OEM / OEM number.
    private String description;             // รายละเอียด / Description.
    private String unit;                    // หน่วย / Unit.
    private Integer reorderLevel;           // ระดับสั่งซื้อ / Reorder level.
    private Integer reorderQuantity;        // จำนวนสั่งซื้อ / Reorder quantity.
    private Integer stockQuantity;          // จำนวนในสต็อก / Stock quantity.
    private Integer minStock;               // จำนวนขั้นต่ำ / Min stock.
    private Integer maxStock;               // จำนวนสูงสุด / Max stock.
    private BigDecimal unitCost;            // ต้นทุนต่อหน่วย / Unit cost.
    private BigDecimal sellingPrice;        // ราคาขาย / Selling price.
    private UUID locationId;                // ID ตำแหน่งจัดเก็บ / Location ID.
    private String status;                  // สถานะ / Status.
    private String imageUrl;                // URL รูปภาพ / Image URL.
    private String notes;                   // หมายเหตุ / Notes.
    private LocalDateTime lastUpdatedStock; // อัปเดตสต็อกล่าสุด / Last stock update.

    /*
        ฟังก์ชันนี้ตรวจสอบว่าสินค้าต่ำกว่า Reorder Level หรือไม่
        This function checks if stock is below reorder level.
    */
    public boolean isLowStock() {
        return this.stockQuantity != null && this.reorderLevel != null
                && this.stockQuantity <= this.reorderLevel;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าสินค้าในสต็อกเพียงพอตามจำนวนที่ต้องการหรือไม่
        This function checks if there is sufficient stock for a required quantity.
    */
    public boolean hasSufficientStock(Integer requiredQuantity) {
        return this.stockQuantity != null && this.stockQuantity >= requiredQuantity;
    }

    /*
        ฟังก์ชันนี้เพิ่มจำนวนสต็อก (ใช้เมื่อรับสินค้าเข้า)
        This function increases stock quantity (used when receiving goods).
    */
    public void increaseStock(Integer quantity) {
        if (quantity == null || quantity <= 0) {
            throw new IllegalArgumentException("Quantity must be positive.");
        }
        this.stockQuantity = (this.stockQuantity != null ? this.stockQuantity : 0) + quantity;
        this.lastUpdatedStock = LocalDateTime.now();
    }

    /*
        ฟังก์ชันนี้ลดจำนวนสต็อก (ใช้เมื่อเบิกจ่ายสินค้า)
        This function decreases stock quantity (used when issuing goods).
    */
    public void decreaseStock(Integer quantity) {
        if (quantity == null || quantity <= 0) {
            throw new IllegalArgumentException("Quantity must be positive.");
        }
        if (!hasSufficientStock(quantity)) {
            throw new IllegalArgumentException("Insufficient stock. Available: " + this.stockQuantity);
        }
        this.stockQuantity -= quantity;
        this.lastUpdatedStock = LocalDateTime.now();
    }

    /*
        ฟังก์ชันนี้คำนวณมูลค่าสต็อกทั้งหมด (stockQuantity * unitCost)
        This function calculates total stock value (stockQuantity * unitCost).
    */
    public BigDecimal getTotalStockValue() {
        if (this.stockQuantity == null || this.unitCost == null) {
            return BigDecimal.ZERO;
        }
        return this.unitCost.multiply(new BigDecimal(this.stockQuantity));
    }
}
```

### `domain/TInventory.java`

```java
package com.template.app.modules.inventory.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.inventory.domain.enums.InventoryTransactionType;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TInventory extends GenericBusinessClass {

    private UUID partId;                        // ID อะไหล่ / Part ID.
    private InventoryTransactionType transactionType; // ประเภทการเคลื่อนไหว
    private String referenceType;               // ประเภทอ้างอิง (PO, JOB, ADJUSTMENT)
    private UUID referenceId;                   // ID อ้างอิง
    private Integer quantity;                   // จำนวนที่เคลื่อนไหว
    private Integer previousQuantity;           // จำนวนก่อนเคลื่อนไหว
    private Integer newQuantity;                // จำนวนหลังเคลื่อนไหว
    private BigDecimal unitCost;                // ต้นทุนต่อหน่วย
    private BigDecimal totalCost;               // ต้นทุนรวม
    private LocalDateTime transactionDate;      // วันที่ทำรายการ
    private String note;                        // หมายเหตุ
    private UUID performedBy;                   // ผู้ทำรายการ

    /*
        ฟังก์ชันนี้สร้างบันทึกการเคลื่อนไหวสินค้าจากการรับสินค้าเข้า
        This function creates an inventory transaction record for receiving goods.
    */
    public static TInventory createReceiveTransaction(
            UUID partId,
            Integer quantity,
            BigDecimal unitCost,
            String referenceType,
            UUID referenceId,
            Integer previousQuantity,
            Integer newQuantity,
            UUID performedBy,
            String note) {
        TInventory transaction = new TInventory();
        transaction.setPartId(partId);
        transaction.setTransactionType(InventoryTransactionType.RECEIVE);
        transaction.setReferenceType(referenceType);
        transaction.setReferenceId(referenceId);
        transaction.setQuantity(quantity);
        transaction.setPreviousQuantity(previousQuantity);
        transaction.setNewQuantity(newQuantity);
        transaction.setUnitCost(unitCost);
        transaction.setTotalCost(unitCost != null ? unitCost.multiply(new BigDecimal(quantity)) : BigDecimal.ZERO);
        transaction.setTransactionDate(LocalDateTime.now());
        transaction.setPerformedBy(performedBy);
        transaction.setNote(note);
        return transaction;
    }

    /*
        ฟังก์ชันนี้สร้างบันทึกการเคลื่อนไหวสินค้าจากการเบิกจ่าย
        This function creates an inventory transaction record for issuing goods.
    */
    public static TInventory createIssueTransaction(
            UUID partId,
            Integer quantity,
            BigDecimal unitCost,
            String referenceType,
            UUID referenceId,
            Integer previousQuantity,
            Integer newQuantity,
            UUID performedBy,
            String note) {
        TInventory transaction = new TInventory();
        transaction.setPartId(partId);
        transaction.setTransactionType(InventoryTransactionType.ISSUE);
        transaction.setReferenceType(referenceType);
        transaction.setReferenceId(referenceId);
        transaction.setQuantity(quantity);
        transaction.setPreviousQuantity(previousQuantity);
        transaction.setNewQuantity(newQuantity);
        transaction.setUnitCost(unitCost);
        transaction.setTotalCost(unitCost != null ? unitCost.multiply(new BigDecimal(quantity)) : BigDecimal.ZERO);
        transaction.setTransactionDate(LocalDateTime.now());
        transaction.setPerformedBy(performedBy);
        transaction.setNote(note);
        return transaction;
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/InventoryServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.inventory.application.impl;

import com.template.app.modules.inventory.application.interfaces.InventoryService;
import com.template.app.modules.inventory.domain.MPartMaster;
import com.template.app.modules.inventory.domain.TInventory;
import com.template.app.modules.inventory.infrastructure.cache.InventoryCacheService;
import com.template.app.modules.inventory.infrastructure.cache.PartMasterCacheService;
import com.template.app.modules.inventory.infrastructure.repository.InventoryRepository;
import com.template.app.modules.inventory.infrastructure.repository.PartMasterRepository;
import com.template.app.modules.inventory.presentation.dto.request.InventoryReceiveRequestDTO;
import com.template.app.modules.inventory.presentation.dto.response.InventoryResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
public class InventoryServiceImpl extends GenericServiceImpl<TInventory, InventoryRepository>
        implements InventoryService {

    private final PartMasterRepository partMasterRepository;
    private final PartMasterCacheService partCacheService;
    private final InventoryCacheService inventoryCacheService;

    public InventoryServiceImpl(InventoryRepository repository,
                                PartMasterRepository partMasterRepository,
                                PartMasterCacheService partCacheService,
                                InventoryCacheService inventoryCacheService) {
        super(repository);
        this.partMasterRepository = partMasterRepository;
        this.partCacheService = partCacheService;
        this.inventoryCacheService = inventoryCacheService;
    }

    /*
        ฟังก์ชันนี้รับสินค้าเข้า Inventory ตรวจสอบว่ามี Part อยู่หรือไม่ อัปเดตสต็อก และบันทึกประวัติ
        This function receives goods into inventory, validates part exists, updates stock, and logs history.
    */
    @Override
    @Transactional
    public InventoryResponseDTO receiveGoods(InventoryReceiveRequestDTO request) throws SystemGlobalException {
        // 1. ดึงข้อมูล Part (จาก Cache หรือ DB) / Get Part (from cache or DB).
        MPartMaster part = partCacheService.getPart(request.getPartId());
        if (part == null) {
            part = partMasterRepository.findById(request.getPartId(), getRepositoryAuth())
                    .orElseThrow(() -> new SystemGlobalException("Part not found: " + request.getPartId(), null));
            partCacheService.savePart(part);
        }

        // 2. บันทึกจำนวนสต็อกเดิม / Record previous stock quantity.
        Integer previousQuantity = part.getStockQuantity() != null ? part.getStockQuantity() : 0;

        // 3. อัปเดตสต็อก / Update stock.
        part.increaseStock(request.getQuantity());

        // 4. อัปเดต Part ในฐานข้อมูลและ Cache / Update Part in DB and Cache.
        partMasterRepository.update(part, getRepositoryAuth());
        partCacheService.savePart(part);

        // 5. สร้างบันทึกการเคลื่อนไหวสินค้า (TInventory) / Create inventory transaction record.
        TInventory transaction = TInventory.createReceiveTransaction(
                part.getId(),
                request.getQuantity(),
                request.getUnitCost(),
                request.getReferenceType(),
                request.getReferenceId(),
                previousQuantity,
                part.getStockQuantity(),
                getUserId(),
                request.getNote()
        );

        // 6. บันทึกประวัติการเคลื่อนไหว / Save transaction history.
        TInventory savedTransaction = this.create(transaction);

        // 7. ลบ Cache สรุปสต็อกของ Part นี้ / Evict stock summary cache for this part.
        inventoryCacheService.evictStockSummary(part.getId());

        return InventoryResponseDTO.fromEntity(savedTransaction, part);
    }

    /*
        ฟังก์ชันนี้เบิกจ่ายสินค้าออกจาก Inventory ตรวจสอบว่ามีสต็อกเพียงพอหรือไม่
        This function issues goods from inventory, checking if there is sufficient stock.
    */
    @Override
    @Transactional
    public InventoryResponseDTO issueGoods(InventoryIssueRequestDTO request) throws SystemGlobalException {
        // 1. ดึงข้อมูล Part / Get Part.
        MPartMaster part = partCacheService.getPart(request.getPartId());
        if (part == null) {
            part = partMasterRepository.findById(request.getPartId(), getRepositoryAuth())
                    .orElseThrow(() -> new SystemGlobalException("Part not found: " + request.getPartId(), null));
            partCacheService.savePart(part);
        }

        // 2. ตรวจสอบสต็อกเพียงพอหรือไม่ / Check if stock is sufficient.
        if (!part.hasSufficientStock(request.getQuantity())) {
            throw new SystemGlobalException(
                    "Insufficient stock. Available: " + part.getStockQuantity() + ", Requested: " + request.getQuantity(),
                    null);
        }

        // 3. บันทึกจำนวนสต็อกเดิม / Record previous stock quantity.
        Integer previousQuantity = part.getStockQuantity();

        // 4. ลดสต็อก / Decrease stock.
        part.decreaseStock(request.getQuantity());

        // 5. อัปเดต Part ในฐานข้อมูลและ Cache / Update Part in DB and Cache.
        partMasterRepository.update(part, getRepositoryAuth());
        partCacheService.savePart(part);

        // 6. สร้างบันทึกการเคลื่อนไหว / Create inventory transaction record.
        TInventory transaction = TInventory.createIssueTransaction(
                part.getId(),
                request.getQuantity(),
                part.getUnitCost(),
                request.getReferenceType(),
                request.getReferenceId(),
                previousQuantity,
                part.getStockQuantity(),
                getUserId(),
                request.getNote()
        );

        // 7. บันทึกประวัติ / Save transaction.
        TInventory savedTransaction = this.create(transaction);

        // 8. ลบ Cache สรุปสต็อก / Evict stock summary cache.
        inventoryCacheService.evictStockSummary(part.getId());

        return InventoryResponseDTO.fromEntity(savedTransaction, part);
    }
}
```

---

## ✅ สรุปส่วนที่เพิ่มเติมในโมดูล Inventory Management

| ฟีเจอร์ | สถานะ | คำอธิบาย |
|---------|--------|----------|
| **Database Design** | ✅ ครบถ้วน | `m_part_master`, `m_stock_location`, `t_inventory`, `t_inventory_adjustment_header/detail`, `t_part_picking_request/detail`, `t_stocktake_header/detail` พร้อม Trigger |
| **Cache (Redis)** | ✅ เพิ่มแล้ว | `PartMasterCacheService` (ID + Code), `InventoryCacheService` (Stock Summary) |
| **Rate Limit** | ✅ เพิ่มแล้ว | ทุก Controller และ Endpoint |
| **API Routing** | ✅ ชัดเจน | Part Master CRUD + Inventory Receive/Issue + Part Picking + Stock Adjustment + Stock Take + Stock Location |
| **Bilingual Comments** | ✅ ครบถ้วน | ทุก Function มีคำอธิบายภาษาไทยและภาษาอังกฤษแยกบรรทัด |
| **PDF Generation** | ✅ เพิ่มแล้ว | Part Picking Report, Inventory Report |

---

## 📊 สรุปโมดูลที่ดำเนินการแล้ว (Modules Completed)

| # | โมดูล | สถานะ | รายละเอียด |
|---|-------|--------|-----------|
| 1 | 🔑 Authentication & Permission | ✅ ครบถ้วน | JWT + Role/Permission + Rate Limit + Redis Cache |
| 2 | 🚗 Job Card Management | ✅ ครบถ้วน | 14 Statuses + Service/Part + History + Cache |
| 3 | 👥 Customer Management | ✅ ครบถ้วน | Customer + Car + History + Cache (ID/Phone/Plate) |
| 4 | 📋 Quotation | ✅ ครบถ้วน | Quotation + Part/Service + Approve/Reject + PDF + Cache |
| 5 | 🛒 Purchase Order | ✅ ครบถ้วน | PO + Status + Send/Receive + PDF + Email + Cache |
| 6 | 📦 Inventory Management | ✅ ครบถ้วน | Part Master + Receive/Issue + Picking + Stock Take + Adjustment + Cache |

---
 