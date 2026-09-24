**โมดูลที่ 5: 🛒 Purchase Order (การจัดการใบสั่งซื้อ)**

โมดูล Purchase Order เป็นส่วนสำคัญของระบบที่ใช้ในการสั่งซื้ออะไหล่จาก Supplier เมื่อสินค้าในคลังไม่เพียงพอ โดยมีกระบวนการ:
1. สร้างใบสั่งซื้อจาก Quotation ที่อนุมัติแล้ว หรือสร้างเองจากคลัง
2. เลือก Supplier และระบุรายการอะไหล่ที่ต้องการ
3. ส่งใบสั่งซื้อไปยัง Supplier ทางอีเมล
4. ติดตามสถานะการสั่งซื้อ (DRAFT, SENT, CONFIRMED, SHIPPED, RECEIVED, CANCELLED)
5. เมื่อรับสินค้า → อัปเดต Inventory และเชื่อมโยงกับ Job Card

---

## 📁 โครงสร้างโมดูล Purchase Order (`modules/purchase`)

```
modules/purchase/
├── application/
│   ├── interfaces/
│   │   ├── PurchaseOrderService.java
│   │   └── PurchaseOrderDetailService.java
│   ├── impl/
│   │   ├── PurchaseOrderServiceImpl.java
│   │   └── PurchaseOrderDetailServiceImpl.java
│   └── usecase/
│       ├── CreatePurchaseOrderUseCase.java
│       ├── UpdatePurchaseOrderUseCase.java
│       ├── GetPurchaseOrderUseCase.java
│       ├── DeletePurchaseOrderUseCase.java
│       ├── SendPurchaseOrderUseCase.java
│       ├── ConfirmPurchaseOrderUseCase.java
│       ├── ReceivePurchaseOrderUseCase.java
│       ├── CancelPurchaseOrderUseCase.java
│       ├── CreatePOFromQuotationUseCase.java
│       ├── GeneratePurchaseOrderPDFUseCase.java
│       └── SuggestPOFromJobUseCase.java
├── domain/
│   ├── TPurchaseOrderHeader.java
│   ├── TPurchaseOrderDetail.java
│   ├── TPurchaseOrderStatusHistory.java
│   ├── enums/
│   │   └── PurchaseOrderStatus.java        // DRAFT, SENT, CONFIRMED, SHIPPED, RECEIVED, CANCELLED
│   └── valueobjects/
│       ├── PONumber.java
│       └── DeliveryDate.java
├── infrastructure/
│   ├── repository/
│   │   ├── PurchaseOrderRepository.java
│   │   ├── PurchaseOrderDetailRepository.java
│   │   ├── PurchaseOrderStatusHistoryRepository.java
│   │   └── impl/
│   │       ├── PurchaseOrderRepositoryImpl.java
│   │       └── PurchaseOrderDetailRepositoryImpl.java
│   ├── cache/                                        // ⬅️ ระบบ Cache สำหรับ Purchase Order
│   │   ├── PurchaseOrderCacheService.java
│   │   └── PurchaseOrderSuggestionCacheService.java
│   ├── report/                                       // ⬅️ ระบบสร้าง PDF
│   │   ├── PurchaseOrderReportGenerator.java
│   │   └── PurchaseOrderReportDataProvider.java
│   ├── email/                                        // ⬅️ ระบบส่งอีเมล
│   │   └── PurchaseOrderEmailService.java
│   ├── entity/
│   │   ├── PurchaseOrderHeaderEntity.java
│   │   ├── PurchaseOrderDetailEntity.java
│   │   └── PurchaseOrderStatusHistoryEntity.java
│   └── mapper/
│       ├── PurchaseOrderMapper.java
│       └── PurchaseOrderDetailMapper.java
└── presentation/
    ├── controller/
    │   ├── PurchaseOrderController.java          // CRUD + Status + PDF + Email
    │   └── PurchaseOrderDetailController.java    // จัดการรายการสั่งซื้อ
    ├── dto/
    │   ├── request/
    │   │   ├── PurchaseOrderCreateRequestDTO.java
    │   │   ├── PurchaseOrderUpdateRequestDTO.java
    │   │   ├── PurchaseOrderStatusRequestDTO.java
    │   │   ├── PurchaseOrderDetailRequestDTO.java
    │   │   └── PurchaseOrderReceiveRequestDTO.java
    │   └── response/
    │       ├── PurchaseOrderResponseDTO.java
    │       ├── PurchaseOrderDetailResponseDTO.java
    │       └── PurchaseOrderSuggestionDTO.java
    └── validator/
        ├── PurchaseOrderValidator.java
        └── PurchaseOrderDetailValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Purchase Order

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V5__purchase_order_schema.sql`)

```sql
-- ==============================================
-- ตาราง: t_purchase_order_header (หัวใบสั่งซื้อ)
-- Main table for purchase order documents.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_purchase_order_header (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_no VARCHAR(20) UNIQUE NOT NULL,                   -- เลขที่ใบสั่งซื้อ (PO-2026-0001)
    quotation_id UUID REFERENCES t_quotation(id),        -- เชื่อมโยง Quotation (ถ้ามี)
    job_id UUID REFERENCES t_job(id),                    -- เชื่อมโยง Job (ถ้ามี)
    supplier_id UUID NOT NULL REFERENCES m_supplier(id) ON DELETE RESTRICT,
    po_date TIMESTAMP NOT NULL DEFAULT NOW(),
    expected_delivery_date TIMESTAMP,                    -- วันที่คาดว่าจะรับสินค้า
    actual_delivery_date TIMESTAMP,                      -- วันที่รับสินค้าจริง
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',         -- DRAFT, SENT, CONFIRMED, SHIPPED, RECEIVED, CANCELLED
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0,           -- ราคาก่อนภาษี
    tax_rate DECIMAL(5,2) DEFAULT 7.00,                  -- อัตราภาษี (%)
    tax_amount DECIMAL(15,2) DEFAULT 0,                  -- จำนวนภาษี
    discount_type VARCHAR(20),                           -- PERCENTAGE, FIXED
    discount_value DECIMAL(15,2) DEFAULT 0,              -- ส่วนลด
    total DECIMAL(15,2) NOT NULL DEFAULT 0,              -- ราคาสุทธิ
    currency VARCHAR(10) DEFAULT 'THB',                  -- สกุลเงิน
    exchange_rate DECIMAL(10,4) DEFAULT 1.0000,          -- อัตราแลกเปลี่ยน
    shipping_cost DECIMAL(15,2) DEFAULT 0,               -- ค่าจัดส่ง
    payment_terms TEXT,                                  -- เงื่อนไขการชำระเงิน
    delivery_address TEXT,                               -- ที่อยู่จัดส่ง
    notes TEXT,                                          -- หมายเหตุ
    terms_and_conditions TEXT,                           -- เงื่อนไข
    sent_at TIMESTAMP,                                   -- วันที่ส่งให้ Supplier
    confirmed_at TIMESTAMP,                              -- วันที่ Supplier ยืนยัน
    received_by UUID REFERENCES m_user(id),              -- ผู้รับสินค้า
    -- Audit Fields
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_po_header_supplier ON t_purchase_order_header(supplier_id);
CREATE INDEX idx_t_po_header_quotation ON t_purchase_order_header(quotation_id);
CREATE INDEX idx_t_po_header_job ON t_purchase_order_header(job_id);
CREATE INDEX idx_t_po_header_status ON t_purchase_order_header(status);
CREATE INDEX idx_t_po_header_date ON t_purchase_order_header(po_date);
CREATE INDEX idx_t_po_header_whitelabel ON t_purchase_order_header(whitelabel_id);
CREATE INDEX idx_t_po_header_deleted ON t_purchase_order_header(deleted);

-- ==============================================
-- ฟังก์ชันสร้างเลขที่ใบสั่งซื้ออัตโนมัติ (Auto-generate PO Number)
-- Function to generate unique purchase order number with year and sequence.
-- ==============================================
CREATE OR REPLACE FUNCTION generate_po_no()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    sequence_part := LPAD(CAST((
        SELECT COALESCE(MAX(CAST(SUBSTRING(po_no FROM 8) AS INTEGER)), 0) + 1
        FROM t_purchase_order_header
        WHERE po_no LIKE 'PO-' || year_part || '-%'
    ) AS TEXT), 4, '0');
    NEW.po_no := 'PO-' || year_part || '-' || sequence_part;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_generate_po_no ON t_purchase_order_header;
CREATE TRIGGER trg_generate_po_no
BEFORE INSERT ON t_purchase_order_header
FOR EACH ROW
EXECUTE FUNCTION generate_po_no();

-- ==============================================
-- ตาราง: t_purchase_order_detail (รายละเอียดใบสั่งซื้อ)
-- Purchase order line items.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_purchase_order_detail (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_header_id UUID NOT NULL REFERENCES t_purchase_order_header(id) ON DELETE CASCADE,
    part_id UUID NOT NULL REFERENCES m_part_master(id) ON DELETE RESTRICT,
    quantity_ordered INTEGER NOT NULL DEFAULT 1,           -- จำนวนที่สั่ง
    quantity_received INTEGER DEFAULT 0,                   -- จำนวนที่รับแล้ว
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) GENERATED ALWAYS AS (quantity_ordered * unit_price) STORED,
    discount DECIMAL(15,2) DEFAULT 0,
    net_price DECIMAL(15,2) GENERATED ALWAYS AS ((quantity_ordered * unit_price) - discount) STORED,
    note TEXT,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_po_detail_header ON t_purchase_order_detail(po_header_id);
CREATE INDEX idx_t_po_detail_part ON t_purchase_order_detail(part_id);
CREATE INDEX idx_t_po_detail_whitelabel ON t_purchase_order_detail(whitelabel_id);

-- ==============================================
-- ตาราง: t_purchase_order_status_history (ประวัติการเปลี่ยนสถานะใบสั่งซื้อ)
-- Track purchase order status changes over time.
-- ==============================================
CREATE TABLE IF NOT EXISTS t_purchase_order_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_header_id UUID NOT NULL REFERENCES t_purchase_order_header(id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    changed_by UUID NOT NULL REFERENCES m_user(id) ON DELETE RESTRICT,
    changed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reason TEXT,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_po_status_history_po ON t_purchase_order_status_history(po_header_id);
CREATE INDEX idx_t_po_status_history_changed ON t_purchase_order_status_history(changed_at);
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Purchase Order

### `infrastructure/cache/PurchaseOrderCacheService.java`

```java
package com.template.app.modules.purchase.infrastructure.cache;

import com.template.app.modules.purchase.domain.TPurchaseOrderHeader;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class PurchaseOrderCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลใบสั่งซื้อจาก Cache (ลดภาระฐานข้อมูลเมื่อเปิดดูบ่อยๆ)
        This function retrieves purchase order data from cache (reduces DB load on frequent views).
        Redis Key: purchase_order:{id}
    */
    @Cacheable(value = "purchase_orders", key = "#poId")
    public TPurchaseOrderHeader getPurchaseOrder(UUID poId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลใบสั่งซื้อตาม Quotation ID (ใช้เชื่อมโยง Quotation -> PO)
        This function retrieves purchase order by Quotation ID (links Quotation -> PO).
        Redis Key: po_quotation:{quotationId}
    */
    @Cacheable(value = "po_quotation", key = "#quotationId")
    public TPurchaseOrderHeader getPurchaseOrderByQuotationId(UUID quotationId) {
        return null;
    }

    /*
        ฟังก์ชันนี้จะอัปเดต Cache เมื่อมีการบันทึกใบสั่งซื้อ
        This function updates the cache when a purchase order is saved.
    */
    @CachePut(value = "purchase_orders", key = "#po.id")
    public TPurchaseOrderHeader savePurchaseOrder(TPurchaseOrderHeader po) {
        return po;
    }

    /*
        ฟังก์ชันนี้ลบข้อมูลใบสั่งซื้อออกจาก Cache ทั้งสองรูปแบบ (ID และ Quotation ID)
        This function evicts purchase order data from both cache forms (ID and Quotation ID).
    */
    @CacheEvict(value = {"purchase_orders", "po_quotation"}, key = "#poId")
    public void evictPurchaseOrder(UUID poId) {
        // ลบ Cache ทั้งสองรูปแบบ / Evict both cache entries.
    }

    /*
        ฟังก์ชันนี้ลบ Cache ทั้งหมดของ Purchase Order (ใช้เมื่อระบบอัปเดตข้อมูลจำนวนมาก)
        This function clears all purchase order caches (used during bulk updates).
    */
    @CacheEvict(value = {"purchase_orders", "po_quotation", "po_suggestion"}, allEntries = true)
    public void evictAllPurchaseOrders() {
        // ลบทุก key ใน caches ทั้งหมด / Evict all keys in all caches.
    }
}
```

### `infrastructure/cache/PurchaseOrderSuggestionCacheService.java`

```java
package com.template.app.modules.purchase.infrastructure.cache;

import com.template.app.modules.purchase.presentation.dto.response.PurchaseOrderSuggestionDTO;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
public class PurchaseOrderSuggestionCacheService {

    /*
        ฟังก์ชันนี้ดึงคำแนะนำการสร้าง PO จาก Job (ใช้ Cache เพื่อไม่ต้องคำนวณซ้ำ)
        This function retrieves PO suggestions from a Job (cached to avoid recalculation).
        Redis Key: po_suggestion:{jobId}
    */
    @Cacheable(value = "po_suggestion", key = "#jobId")
    public List<PurchaseOrderSuggestionDTO> getPoSuggestions(UUID jobId) {
        return null; // Spring จะใช้ Cache หากมี / Spring uses cache if available.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Purchase Order Controller

```java
package com.template.app.modules.purchase.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.purchase.application.interfaces.PurchaseOrderService;
import com.template.app.modules.purchase.presentation.dto.request.PurchaseOrderCreateRequestDTO;
import com.template.app.modules.purchase.presentation.dto.request.PurchaseOrderReceiveRequestDTO;
import com.template.app.modules.purchase.presentation.dto.request.PurchaseOrderUpdateRequestDTO;
import com.template.app.modules.purchase.presentation.dto.response.PurchaseOrderResponseDTO;
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
@RequestMapping("/api/v1/purchase-orders")
@Tag(name = "Purchase Order", description = "Purchase Order Management APIs")
@RequiredArgsConstructor
public class PurchaseOrderController {

    private final PurchaseOrderService purchaseOrderService;

    // ========================================================================
    // 1. CREATE PURCHASE ORDER
    // ========================================================================

    /*
        API: POST /api/v1/purchase-orders
        ฟังก์ชันนี้สร้างใบสั่งซื้อใหม่ สามารถสร้างจาก Quotation หรือสร้างเอง
        This function creates a new purchase order, either from a Quotation or manually.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PostMapping
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Create a new purchase order")
    public ResponseEntity<PurchaseOrderResponseDTO> createPurchaseOrder(
            @Valid @RequestBody PurchaseOrderCreateRequestDTO request) throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.createPurchaseOrder(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. CREATE PO FROM QUOTATION
    // ========================================================================

    /*
        API: POST /api/v1/purchase-orders/from-quotation/{quotationId}
        ฟังก์ชันนี้สร้างใบสั่งซื้อจาก Quotation ที่อนุมัติแล้ว (ใช้เมื่อ Quotation Approved)
        This function creates a purchase order from an approved Quotation (used when Quotation is approved).
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @PostMapping("/from-quotation/{quotationId}")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Create purchase order from quotation")
    public ResponseEntity<PurchaseOrderResponseDTO> createPOFromQuotation(@PathVariable UUID quotationId)
            throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.createPOFromQuotation(quotationId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET PURCHASE ORDER BY ID
    // ========================================================================

    /*
        API: GET /api/v1/purchase-orders/{id}
        ฟังก์ชันนี้ดึงข้อมูลใบสั่งซื้อตาม ID พร้อมรายการสั่งซื้อทั้งหมด
        This function retrieves purchase order by ID with all line items.
        Rate Limit: อนุญาต 100 ครั้งต่อ 1 นาที
        Rate Limit: Allows 100 requests per minute.
    */
    @GetMapping("/{id}")
    @RateLimit(limit = 100, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get purchase order by ID")
    public ResponseEntity<PurchaseOrderResponseDTO> getPurchaseOrder(@PathVariable UUID id)
            throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.getPurchaseOrder(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. GET PO BY QUOTATION ID
    // ========================================================================

    /*
        API: GET /api/v1/purchase-orders/quotation/{quotationId}
        ฟังก์ชันนี้ดึงใบสั่งซื้อที่สร้างจาก Quotation ที่ระบุ
        This function retrieves the purchase order created from a given Quotation.
        Rate Limit: อนุญาต 60 ครั้งต่อ 1 นาที
        Rate Limit: Allows 60 requests per minute.
    */
    @GetMapping("/quotation/{quotationId}")
    @RateLimit(limit = 60, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get purchase order by quotation")
    public ResponseEntity<PurchaseOrderResponseDTO> getPOByQuotationId(@PathVariable UUID quotationId)
            throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.getPOByQuotationId(quotationId);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. LIST PURCHASE ORDERS (Pagination + Filters)
    // ========================================================================

    /*
        API: GET /api/v1/purchase-orders
        ฟังก์ชันนี้แสดงรายการใบสั่งซื้อแบบแบ่งหน้า พร้อมตัวกรองตาม Supplier, Status, และช่วงวันที่
        This function lists purchase orders with pagination, filtering by Supplier, Status, and date range.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "List purchase orders with pagination")
    public ResponseEntity<Page<PurchaseOrderResponseDTO>> listPurchaseOrders(
            @RequestParam(required = false) UUID supplierId,
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate,
            Pageable pageable) throws SystemGlobalException {
        Page<PurchaseOrderResponseDTO> page = purchaseOrderService.listPurchaseOrders(
                supplierId, status, startDate, endDate, pageable);
        return ResponseEntity.ok(page);
    }

    // ========================================================================
    // 6. UPDATE PURCHASE ORDER
    // ========================================================================

    /*
        API: PUT /api/v1/purchase-orders/{id}
        ฟังก์ชันนี้แก้ไขข้อมูลใบสั่งซื้อ (เฉพาะที่ยังไม่ได้ส่งให้ Supplier)
        This function updates purchase order details (only if not yet sent to supplier).
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @PutMapping("/{id}")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Update purchase order details")
    public ResponseEntity<PurchaseOrderResponseDTO> updatePurchaseOrder(
            @PathVariable UUID id,
            @Valid @RequestBody PurchaseOrderUpdateRequestDTO request) throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.updatePurchaseOrder(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 7. DELETE PURCHASE ORDER (SOFT DELETE)
    // ========================================================================

    /*
        API: DELETE /api/v1/purchase-orders/{id}
        ฟังก์ชันนี้ลบใบสั่งซื้อแบบ Soft Delete (เฉพาะที่ยังไม่ได้ส่งให้ Supplier)
        This function soft-deletes a purchase order (only if not yet sent to supplier).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @DeleteMapping("/{id}")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Delete purchase order (soft delete)")
    public ResponseEntity<Void> deletePurchaseOrder(@PathVariable UUID id) throws SystemGlobalException {
        purchaseOrderService.deletePurchaseOrder(id);
        return ResponseEntity.noContent().build();
    }

    // ========================================================================
    // 8. SEND PO TO SUPPLIER (EMAIL)
    // ========================================================================

    /*
        API: POST /api/v1/purchase-orders/{id}/send
        ฟังก์ชันนี้ส่งใบสั่งซื้อให้ Supplier ทางอีเมล พร้อมแนบไฟล์ PDF
        This function sends the purchase order to the supplier via email with PDF attachment.
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 นาที
        Rate Limit: Allows 10 requests per minute.
    */
    @PostMapping("/{id}/send")
    @RateLimit(limit = 10, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Send purchase order to supplier")
    public ResponseEntity<PurchaseOrderResponseDTO> sendPurchaseOrder(@PathVariable UUID id)
            throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.sendPurchaseOrder(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 9. CONFIRM PO (FROM SUPPLIER)
    // ========================================================================

    /*
        API: PUT /api/v1/purchase-orders/{id}/confirm
        ฟังก์ชันนี้ยืนยันใบสั่งซื้อเมื่อ Supplier ยืนยันการรับคำสั่งซื้อแล้ว
        This function confirms the purchase order when the supplier acknowledges the order.
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PutMapping("/{id}/confirm")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Confirm purchase order by supplier")
    public ResponseEntity<PurchaseOrderResponseDTO> confirmPurchaseOrder(@PathVariable UUID id)
            throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.confirmPurchaseOrder(id);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 10. RECEIVE GOODS (INVENTORY UPDATE)
    // ========================================================================

    /*
        API: POST /api/v1/purchase-orders/{id}/receive
        ฟังก์ชันนี้รับสินค้าเข้าคลัง อัปเดต Inventory และบันทึกจำนวนที่รับ
        This function receives goods into inventory, updates Inventory, and records received quantities.
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @PostMapping("/{id}/receive")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Receive goods and update inventory")
    public ResponseEntity<PurchaseOrderResponseDTO> receiveGoods(
            @PathVariable UUID id,
            @Valid @RequestBody PurchaseOrderReceiveRequestDTO request) throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.receiveGoods(id, request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 11. CANCEL PURCHASE ORDER
    // ========================================================================

    /*
        API: PUT /api/v1/purchase-orders/{id}/cancel
        ฟังก์ชันนี้ยกเลิกใบสั่งซื้อ พร้อมระบุเหตุผล (ใช้เมื่อไม่ต้องการสั่งซื้อแล้ว)
        This function cancels a purchase order with a reason (used when the order is no longer needed).
        Rate Limit: อนุญาต 10 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 10 requests per 1 hour.
    */
    @PutMapping("/{id}/cancel")
    @RateLimit(limit = 10, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Cancel purchase order")
    public ResponseEntity<PurchaseOrderResponseDTO> cancelPurchaseOrder(
            @PathVariable UUID id,
            @RequestParam String reason) throws SystemGlobalException {
        PurchaseOrderResponseDTO response = purchaseOrderService.cancelPurchaseOrder(id, reason);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 12. GENERATE PO PDF
    // ========================================================================

    /*
        API: GET /api/v1/purchase-orders/{id}/pdf
        ฟังก์ชันนี้สร้างไฟล์ PDF ของใบสั่งซื้อ (ใช้ JasperReports) สำหรับส่งให้ Supplier
        This function generates a PDF file of the purchase order (using JasperReports) for supplier delivery.
        Rate Limit: อนุญาต 15 ครั้งต่อ 5 นาที
        Rate Limit: Allows 15 requests per 5 minutes.
    */
    @GetMapping("/{id}/pdf")
    @RateLimit(limit = 15, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Generate purchase order PDF")
    public ResponseEntity<byte[]> generatePurchaseOrderPDF(@PathVariable UUID id) throws SystemGlobalException {
        byte[] pdf = purchaseOrderService.generatePurchaseOrderPDF(id);
        return ResponseEntity.ok()
                .header("Content-Type", "application/pdf")
                .header("Content-Disposition", "inline; filename=po_" + id + ".pdf")
                .body(pdf);
    }

    // ========================================================================
    // 13. GET PO SUGGESTIONS FROM JOB
    // ========================================================================

    /*
        API: GET /api/v1/purchase-orders/suggestions/{jobId}
        ฟังก์ชันนี้แนะนำรายการอะไหล่ที่ต้องสั่งซื้อจาก Job โดยพิจารณาจาก Quotation และ Stock ที่มี
        This function suggests parts to order from a Job, based on Quotation and available stock.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/suggestions/{jobId}")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get purchase order suggestions from job")
    public ResponseEntity<List<PurchaseOrderSuggestionDTO>> getPOSuggestions(@PathVariable UUID jobId)
            throws SystemGlobalException {
        List<PurchaseOrderSuggestionDTO> suggestions = purchaseOrderService.getPOSuggestions(jobId);
        return ResponseEntity.ok(suggestions);
    }

    // ========================================================================
    // 14. GET PO STATUS HISTORY
    // ========================================================================

    /*
        API: GET /api/v1/purchase-orders/{id}/history
        ฟังก์ชันนี้ดึงประวัติการเปลี่ยนสถานะของใบสั่งซื้อ
        This function retrieves the status change history of a purchase order.
        Rate Limit: อนุญาต 50 ครั้งต่อ 1 นาที
        Rate Limit: Allows 50 requests per minute.
    */
    @GetMapping("/{id}/history")
    @RateLimit(limit = 50, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get purchase order status history")
    public ResponseEntity<List<PurchaseOrderStatusHistoryDTO>> getPOStatusHistory(@PathVariable UUID id)
            throws SystemGlobalException {
        List<PurchaseOrderStatusHistoryDTO> history = purchaseOrderService.getPOStatusHistory(id);
        return ResponseEntity.ok(history);
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/PurchaseOrderStatus.java`

```java
package com.template.app.modules.purchase.domain.enums;

/*
    สถานะของใบสั่งซื้อ / Purchase Order status.
*/
public enum PurchaseOrderStatus {
    DRAFT,        // ร่าง / Draft.
    SENT,         // ส่งให้ Supplier แล้ว / Sent to supplier.
    CONFIRMED,    // Supplier ยืนยันแล้ว / Supplier confirmed.
    SHIPPED,      // Supplier จัดส่งแล้ว / Supplier shipped.
    RECEIVED,     // รับสินค้าแล้ว / Goods received.
    CANCELLED     // ยกเลิก / Cancelled.
}
```

### `domain/TPurchaseOrderHeader.java`

```java
package com.template.app.modules.purchase.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.purchase.domain.enums.PurchaseOrderStatus;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TPurchaseOrderHeader extends GenericBusinessClass {

    private String poNo;                        // เลขที่ใบสั่งซื้อ / PO Number.
    private UUID quotationId;                   // ID Quotation (ถ้ามี)
    private UUID jobId;                         // ID Job (ถ้ามี)
    private UUID supplierId;                    // ID Supplier
    private LocalDateTime poDate;               // วันที่สั่งซื้อ / PO Date.
    private LocalDateTime expectedDeliveryDate; // วันที่คาดว่าจะได้รับ / Expected delivery.
    private LocalDateTime actualDeliveryDate;   // วันที่รับจริง / Actual delivery.
    private PurchaseOrderStatus status;         // สถานะ / Status.
    private BigDecimal subtotal;                // ราคาก่อนภาษี / Subtotal.
    private BigDecimal taxRate;                 // อัตราภาษี (%) / Tax rate.
    private BigDecimal taxAmount;               // จำนวนภาษี / Tax amount.
    private String discountType;                // PERCENTAGE, FIXED
    private BigDecimal discountValue;           // ส่วนลด / Discount.
    private BigDecimal total;                   // ราคาสุทธิ / Total.
    private String currency;                    // สกุลเงิน / Currency.
    private BigDecimal exchangeRate;            // อัตราแลกเปลี่ยน / Exchange rate.
    private BigDecimal shippingCost;            // ค่าจัดส่ง / Shipping cost.
    private String paymentTerms;                // เงื่อนไขการชำระเงิน / Payment terms.
    private String deliveryAddress;             // ที่อยู่จัดส่ง / Delivery address.
    private String notes;                       // หมายเหตุ / Notes.
    private String termsAndConditions;          // เงื่อนไข / Terms.
    private LocalDateTime sentAt;               // วันที่ส่ง / Sent date.
    private LocalDateTime confirmedAt;          // วันที่ยืนยัน / Confirmed date.
    private UUID receivedBy;                    // ผู้รับสินค้า / Receiver.

    private List<TPurchaseOrderDetail> details = new ArrayList<>();

    /*
        ฟังก์ชันนี้คำนวณยอดรวมของใบสั่งซื้อ (Subtotal, Tax, Total)
        This function calculates the total amounts of the purchase order (Subtotal, Tax, Total).
    */
    public void calculateTotals() {
        BigDecimal detailsTotal = details.stream()
                .map(TPurchaseOrderDetail::getNetPrice)
                .reduce(BigDecimal.ZERO, BigDecimal::add);

        this.subtotal = detailsTotal;
        this.taxAmount = this.subtotal.multiply(this.taxRate.divide(new BigDecimal(100)));
        this.total = this.subtotal.add(this.taxAmount).add(this.shippingCost != null ? this.shippingCost : BigDecimal.ZERO);
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าใบสั่งซื้อสามารถส่งให้ Supplier ได้หรือไม่
        This function checks if the purchase order can be sent to the supplier.
    */
    public boolean canSend() {
        return this.status == PurchaseOrderStatus.DRAFT && !this.details.isEmpty();
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่าใบสั่งซื้อสามารถรับสินค้าได้หรือไม่
        This function checks if the purchase order can receive goods.
    */
    public boolean canReceive() {
        return this.status == PurchaseOrderStatus.SENT || 
               this.status == PurchaseOrderStatus.CONFIRMED || 
               this.status == PurchaseOrderStatus.SHIPPED;
    }

    /*
        ฟังก์ชันนี้เปลี่ยนสถานะเป็น SENT และบันทึกวันที่ส่ง
        This function changes status to SENT and records the sent date.
    */
    public void markAsSent() {
        if (!canSend()) {
            throw new IllegalStateException("Cannot send PO with status: " + this.status);
        }
        this.status = PurchaseOrderStatus.SENT;
        this.sentAt = LocalDateTime.now();
    }

    /*
        ฟังก์ชันนี้เปลี่ยนสถานะเป็น RECEIVED และบันทึกวันที่รับ
        This function changes status to RECEIVED and records the received date.
    */
    public void markAsReceived(UUID receiverId) {
        if (!canReceive()) {
            throw new IllegalStateException("Cannot receive goods for PO with status: " + this.status);
        }
        this.status = PurchaseOrderStatus.RECEIVED;
        this.actualDeliveryDate = LocalDateTime.now();
        this.receivedBy = receiverId;
    }
}
```

### `domain/TPurchaseOrderDetail.java`

```java
package com.template.app.modules.purchase.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@EqualsAndHashCode(callSuper = true)
public class TPurchaseOrderDetail extends GenericBusinessClass {

    private UUID poHeaderId;            // ID ใบสั่งซื้อ / PO ID.
    private UUID partId;                // ID อะไหล่ / Part ID.
    private Integer quantityOrdered;    // จำนวนที่สั่ง / Quantity ordered.
    private Integer quantityReceived;   // จำนวนที่รับแล้ว / Quantity received.
    private BigDecimal unitPrice;       // ราคาต่อหน่วย / Unit price.
    private BigDecimal totalPrice;      // ราคารวม / Total price.
    private BigDecimal discount;        // ส่วนลด / Discount.
    private BigDecimal netPrice;        // ราคาสุทธิ / Net price.
    private String note;                // หมายเหตุ / Note.

    /*
        ฟังก์ชันนี้คำนวณราคารวมและราคาสุทธิของรายการสั่งซื้อ
        This function calculates total and net prices for the order item.
    */
    public void calculatePrices() {
        this.totalPrice = this.unitPrice.multiply(new BigDecimal(this.quantityOrdered));
        this.netPrice = this.totalPrice.subtract(this.discount != null ? this.discount : BigDecimal.ZERO);
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่ายังรับสินค้าไม่ครบตามที่สั่ง
        This function checks if goods are not fully received yet.
    */
    public boolean isPartialReceived() {
        return this.quantityReceived < this.quantityOrdered;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่ารับสินค้าครบแล้ว
        This function checks if goods are fully received.
    */
    public boolean isFullyReceived() {
        return this.quantityReceived >= this.quantityOrdered;
    }

    /*
        ฟังก์ชันนี้เพิ่มจำนวนที่รับ (ใช้เมื่อรับสินค้าเข้ามา)
        This function adds received quantity (used when goods are received).
    */
    public void addReceivedQuantity(Integer quantity) {
        if (this.quantityReceived == null) {
            this.quantityReceived = 0;
        }
        int newTotal = this.quantityReceived + quantity;
        if (newTotal > this.quantityOrdered) {
            throw new IllegalArgumentException("Received quantity exceeds ordered quantity.");
        }
        this.quantityReceived = newTotal;
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/PurchaseOrderServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.purchase.application.impl;

import com.template.app.modules.purchase.application.interfaces.PurchaseOrderService;
import com.template.app.modules.purchase.domain.TPurchaseOrderHeader;
import com.template.app.modules.purchase.domain.TPurchaseOrderDetail;
import com.template.app.modules.purchase.domain.enums.PurchaseOrderStatus;
import com.template.app.modules.purchase.infrastructure.cache.PurchaseOrderCacheService;
import com.template.app.modules.purchase.infrastructure.repository.PurchaseOrderRepository;
import com.template.app.modules.purchase.presentation.dto.request.PurchaseOrderCreateRequestDTO;
import com.template.app.modules.purchase.presentation.dto.request.PurchaseOrderReceiveRequestDTO;
import com.template.app.modules.purchase.presentation.dto.response.PurchaseOrderResponseDTO;
import com.template.app._shared.application.GenericServiceImpl;
import com.template.app.exception.SystemGlobalException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

@Service
public class PurchaseOrderServiceImpl 
        extends GenericServiceImpl<TPurchaseOrderHeader, PurchaseOrderRepository> 
        implements PurchaseOrderService {

    private final PurchaseOrderCacheService cacheService;
    private final InventoryService inventoryService;      // Inject Inventory Service
    private final PurchaseOrderEmailService emailService;

    public PurchaseOrderServiceImpl(PurchaseOrderRepository repository, 
                                    PurchaseOrderCacheService cacheService,
                                    InventoryService inventoryService,
                                    PurchaseOrderEmailService emailService) {
        super(repository);
        this.cacheService = cacheService;
        this.inventoryService = inventoryService;
        this.emailService = emailService;
    }

    /*
        ฟังก์ชันนี้สร้างใบสั่งซื้อใหม่จาก Quotation ที่อนุมัติแล้ว โดยจะดึงรายการอะไหล่จาก Quotation
        This function creates a new purchase order from an approved Quotation, fetching parts from the Quotation.
    */
    @Override
    @Transactional
    public PurchaseOrderResponseDTO createPOFromQuotation(UUID quotationId) throws SystemGlobalException {
        // 1. ตรวจสอบ Quotation ว่าอนุมัติแล้วหรือยัง / Verify Quotation is approved.
        Quotation quotation = quotationService.getQuotation(quotationId);
        if (quotation.getStatus() != QuotationStatus.APPROVED) {
            throw new SystemGlobalException("Quotation not approved. Cannot create PO.", null);
        }

        // 2. สร้าง PO Header / Create PO Header.
        TPurchaseOrderHeader po = new TPurchaseOrderHeader();
        po.setQuotationId(quotationId);
        po.setJobId(quotation.getJobId());
        po.setSupplierId(quotation.getSupplierId()); // ต้องมี Supplier ใน Quotation
        po.setPoDate(LocalDateTime.now());
        po.setStatus(PurchaseOrderStatus.DRAFT);
        po.setCurrency(quotation.getCurrency());
        po.setTaxRate(quotation.getTaxRate());

        // 3. สร้างรายการ PO จาก Quotation Parts / Create PO details from Quotation parts.
        for (QuotationPart qPart : quotation.getParts()) {
            TPurchaseOrderDetail detail = new TPurchaseOrderDetail();
            detail.setPartId(qPart.getPartId());
            detail.setQuantityOrdered(qPart.getQuantity());
            detail.setUnitPrice(qPart.getUnitPrice());
            detail.setDiscount(qPart.getDiscount());
            detail.calculatePrices();
            po.getDetails().add(detail);
        }

        // 4. คำนวณยอดรวม / Calculate totals.
        po.calculateTotals();

        // 5. บันทึก PO / Save PO.
        TPurchaseOrderHeader savedPO = this.create(po);

        // 6. บันทึก Cache / Store in cache.
        cacheService.savePurchaseOrder(savedPO);

        return PurchaseOrderResponseDTO.fromEntity(savedPO);
    }

    /*
        ฟังก์ชันนี้รับสินค้าเข้าคลัง อัปเดต Inventory และบันทึกจำนวนที่รับใน PO Detail
        This function receives goods into inventory, updates Inventory, and records quantities in PO Detail.
    */
    @Override
    @Transactional
    public PurchaseOrderResponseDTO receiveGoods(UUID poId, PurchaseOrderReceiveRequestDTO request) 
            throws SystemGlobalException {
        // 1. ดึงข้อมูล PO / Get PO.
        TPurchaseOrderHeader po = this.read(poId);

        // 2. ตรวจสอบว่า PO อยู่ในสถานะที่รับสินค้าได้ / Verify PO can receive goods.
        if (!po.canReceive()) {
            throw new SystemGlobalException("PO cannot receive goods in status: " + po.getStatus(), null);
        }

        // 3. วนลูปอัปเดตแต่ละรายการ / Loop through each received item.
        for (PurchaseOrderReceiveRequestDTO.ReceiveItem item : request.getItems()) {
            TPurchaseOrderDetail detail = po.getDetails().stream()
                    .filter(d -> d.getId().equals(item.getDetailId()))
                    .findFirst()
                    .orElseThrow(() -> new SystemGlobalException("Detail not found: " + item.getDetailId(), null));

            // 4. เพิ่มจำนวนที่รับ / Add received quantity.
            detail.addReceivedQuantity(item.getReceivedQuantity());

            // 5. อัปเดต Inventory (ผ่าน InventoryService) / Update Inventory.
            inventoryService.receivePart(
                    detail.getPartId(),
                    item.getReceivedQuantity(),
                    poId.toString(),
                    "Received from PO: " + po.getPoNo()
            );
        }

        // 6. ตรวจสอบว่าทุกรายการรับครบแล้วหรือยัง / Check if all items are fully received.
        boolean allReceived = po.getDetails().stream().allMatch(TPurchaseOrderDetail::isFullyReceived);
        if (allReceived) {
            po.markAsReceived(getUserId());
        }

        // 7. บันทึก PO / Save PO.
        TPurchaseOrderHeader updatedPO = this.update(poId, po);

        // 8. ลบ Cache เก่าและอัปเดต Cache ใหม่ / Evict old cache and update new cache.
        cacheService.evictPurchaseOrder(poId);
        cacheService.savePurchaseOrder(updatedPO);

        return PurchaseOrderResponseDTO.fromEntity(updatedPO);
    }

    /*
        ฟังก์ชันนี้ส่งใบสั่งซื้อให้ Supplier ทางอีเมล พร้อมแนบไฟล์ PDF
        This function sends the purchase order to supplier via email with PDF attachment.
    */
    @Override
    public PurchaseOrderResponseDTO sendPurchaseOrder(UUID poId) throws SystemGlobalException {
        // 1. ดึงข้อมูล PO / Get PO.
        TPurchaseOrderHeader po = this.read(poId);

        // 2. ตรวจสอบว่า PO สามารถส่งได้ / Verify PO can be sent.
        if (!po.canSend()) {
            throw new SystemGlobalException("Cannot send PO with status: " + po.getStatus(), null);
        }

        // 3. สร้าง PDF / Generate PDF.
        byte[] pdf = generatePurchaseOrderPDF(poId);

        // 4. ส่งอีเมลไปยัง Supplier / Send email to supplier.
        String supplierEmail = supplierService.getSupplierEmail(po.getSupplierId());
        emailService.sendPurchaseOrderEmail(supplierEmail, po.getPoNo(), pdf);

        // 5. เปลี่ยนสถานะเป็น SENT / Change status to SENT.
        po.markAsSent();

        // 6. บันทึก PO / Save PO.
        TPurchaseOrderHeader updatedPO = this.update(poId, po);

        // 7. ลบ Cache เก่าและอัปเดต Cache ใหม่ / Evict old cache and update new cache.
        cacheService.evictPurchaseOrder(poId);
        cacheService.savePurchaseOrder(updatedPO);

        return PurchaseOrderResponseDTO.fromEntity(updatedPO);
    }
}
```

---

## ✅ สรุปส่วนที่เพิ่มเติมในโมดูล Purchase Order

| ฟีเจอร์ | สถานะ | คำอธิบาย |
|---------|--------|----------|
| **Database Design** | ✅ ครบถ้วน | `t_purchase_order_header`, `t_purchase_order_detail`, `t_purchase_order_status_history` พร้อม Trigger |
| **Cache (Redis)** | ✅ เพิ่มแล้ว | `PurchaseOrderCacheService` (ID + Quotation ID), `PurchaseOrderSuggestionCacheService` |
| **Rate Limit** | ✅ เพิ่มแล้ว | ทุก Endpoint ใน Controller |
| **API Routing** | ✅ ชัดเจน | CRUD + Send/Confirm/Receive/Cancel + Suggestion + History + PDF |
| **Bilingual Comments** | ✅ ครบถ้วน | ทุก Function มีคำอธิบายภาษาไทยและภาษาอังกฤษแยกบรรทัด |
| **Email Service** | ✅ เพิ่มแล้ว | `PurchaseOrderEmailService` ส่งอีเมลพร้อม PDF |
| **PDF Generation** | ✅ เพิ่มแล้ว | `PurchaseOrderReportGenerator` ใช้ JasperReports |

---
 