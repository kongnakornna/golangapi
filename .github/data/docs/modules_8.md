**โมดูลที่ 8: 📊 Dashboard & Reports (ระบบรายงานและแดชบอร์ด)**

โมดูล Dashboard & Reports เป็นศูนย์กลางสำหรับการแสดงผลข้อมูลเชิงวิเคราะห์และรายงานสรุปสำหรับผู้บริหารและผู้จัดการ ครอบคลุมการทำงานดังนี้:

1. **Dashboard ภาพรวม (Overview Dashboard)** – แสดงสถิติสำคัญแบบ Real-time
2. **รายงานยอดขาย (Sales Reports)** – รายงานสรุปยอดขายรายวัน/รายเดือน/รายปี
3. **รายงานสินค้าคงคลัง (Inventory Reports)** – สรุปสถานะสินค้าคงคลัง
4. **รายงานการซ่อม (Service Reports)** – สถิติการซ่อมแซมและประสิทธิภาพ
5. **รายงานทางการเงิน (Financial Reports)** – สรุปรายได้-รายจ่าย, ลูกหนี้
6. **การส่งออกข้อมูล (Export)** – Excel, PDF, CSV
7. **Widget แบบกำหนดเอง (Customizable Widgets)** – ผู้ใช้สามารถปรับแต่ง Dashboard ได้

---

## 📁 โครงสร้างโมดูล Dashboard & Reports (`modules/dashboard`)

```
modules/dashboard/
├── application/
│   ├── interfaces/
│   │   ├── DashboardService.java
│   │   ├── ReportService.java
│   │   ├── ExportService.java
│   │   └── WidgetService.java
│   ├── impl/
│   │   ├── DashboardServiceImpl.java
│   │   ├── ReportServiceImpl.java
│   │   ├── ExportServiceImpl.java
│   │   └── WidgetServiceImpl.java
│   └── usecase/
│       ├── GetSalesOverviewUseCase.java
│       ├── GetInventoryOverviewUseCase.java
│       ├── GetJobStatusSummaryUseCase.java
│       ├── GetRevenueByPeriodUseCase.java
│       ├── GetTopSellingPartsUseCase.java
│       ├── GetServiceCategoryUseCase.java
│       ├── GenerateDailyReportUseCase.java
│       ├── GenerateMonthlyReportUseCase.java
│       ├── GenerateYearlyReportUseCase.java
│       ├── ExportExcelUseCase.java
│       └── ExportPDFUseCase.java
├── domain/
│   ├── DSalesOverview.java
│   ├── DInventoryOverview.java
│   ├── DJobStatusSummary.java
│   ├── DRevenueByPeriod.java
│   ├── DTopSellingParts.java
│   ├── DServiceCategory.java
│   ├── DFinancialSummary.java
│   ├── DWidgetConfig.java
│   ├── enums/
│   │   ├── ReportType.java           // DAILY, MONTHLY, YEARLY, CUSTOM
│   │   ├── ExportFormat.java         // EXCEL, PDF, CSV
│   │   └── WidgetType.java           // CHART, TABLE, METRIC, MAP
│   └── valueobjects/
│       ├── DateRange.java
│       └── Period.java
├── infrastructure/
│   ├── repository/
│   │   ├── DashboardRepository.java
│   │   ├── WidgetConfigRepository.java
│   │   └── impl/
│   │       └── DashboardRepositoryImpl.java
│   ├── cache/                                           // ⬅️ ระบบ Cache สำหรับ Dashboard
│   │   ├── DashboardCacheService.java
│   │   └── ReportCacheService.java
│   ├── report/                                          // ⬅️ ระบบสร้างรายงาน
│   │   ├── ReportGenerator.java
│   │   ├── ReportDataProvider.java
│   │   ├── excel/
│   │   │   └── ExcelReportGenerator.java
│   │   └── pdf/
│   │       └── PDFReportGenerator.java
│   ├── query/                                           // ⬅️ Query คอมเพล็กซ์
│   │   ├── DashboardQueryBuilder.java
│   │   └── DashboardNativeQuery.java
│   └── entity/
│       ├── DashboardWidgetEntity.java
│       └── DashboardView.java
└── presentation/
    ├── controller/
    │   ├── DashboardController.java      // Dashboard APIs
    │   ├── ReportController.java         // Report Generation APIs
    │   ├── ExportController.java         // Export APIs
    │   └── WidgetController.java         // Widget Management APIs
    ├── dto/
    │   ├── request/
    │   │   ├── DashboardFilterRequestDTO.java
    │   │   ├── ReportRequestDTO.java
    │   │   ├── ExportRequestDTO.java
    │   │   └── WidgetConfigRequestDTO.java
    │   └── response/
    │       ├── DashboardOverviewResponseDTO.java
    │       ├── SalesOverviewResponseDTO.java
    │       ├── InventoryOverviewResponseDTO.java
    │       ├── JobStatusSummaryResponseDTO.java
    │       ├── RevenueResponseDTO.java
    │       ├── TopPartsResponseDTO.java
    │       ├── ServiceCategoryResponseDTO.java
    │       └── WidgetConfigResponseDTO.java
    └── validator/
        └── ReportValidator.java
```

---

## 🗄️ ออกแบบฐานข้อมูล (Database Design) สำหรับโมดูล Dashboard

### 📄 SQL DDL (ไฟล์: `src/main/resources/db/migration/V8__dashboard_schema.sql`)

```sql
-- ==============================================
-- หมายเหตุ: ตาราง Dashboard ส่วนใหญ่เป็น VIEW หรือ Materialized View
-- ที่ดึงข้อมูลจากตารางหลัก (t_job, t_invoice, t_payment, t_inventory)
-- NOTE: Most Dashboard tables are VIEWs or Materialized Views
-- that pull data from main tables.
-- ==============================================

-- ==============================================
-- VIEW: v_dashboard_sales_overview (ภาพรวมยอดขาย)
-- Sales overview dashboard view.
-- ==============================================
CREATE OR REPLACE VIEW v_dashboard_sales_overview AS
WITH sales_data AS (
    SELECT 
        i.id AS invoice_id,
        i.customer_id,
        i.job_id,
        i.invoice_date,
        i.total AS invoice_total,
        COALESCE((
            SELECT SUM(p.amount) 
            FROM t_payment p 
            WHERE p.invoice_id = i.id 
            AND p.status = 'COMPLETED'
        ), 0) AS amount_paid,
        i.whitelabel_id
    FROM t_invoice_adjustment i
    WHERE i.deleted = false 
    AND i.type = 'INVOICE'
)
SELECT 
    COUNT(DISTINCT invoice_id) AS total_invoices,
    COUNT(DISTINCT customer_id) AS total_customers,
    COALESCE(SUM(invoice_total), 0) AS total_revenue,
    COALESCE(SUM(invoice_total - amount_paid), 0) AS total_outstanding,
    COALESCE(AVG(invoice_total), 0) AS average_invoice,
    DATE_TRUNC('month', invoice_date) AS period,
    whitelabel_id
FROM sales_data
GROUP BY DATE_TRUNC('month', invoice_date), whitelabel_id;

-- ==============================================
-- VIEW: v_dashboard_job_status (สรุปสถานะใบงาน)
-- Job status summary dashboard view.
-- ==============================================
CREATE OR REPLACE VIEW v_dashboard_job_status AS
SELECT 
    status,
    COUNT(*) AS count,
    whitelabel_id
FROM t_job
WHERE deleted = false
GROUP BY status, whitelabel_id;

-- ==============================================
-- VIEW: v_dashboard_inventory_overview (ภาพรวมสินค้าคงคลัง)
-- Inventory overview dashboard view.
-- ==============================================
CREATE OR REPLACE VIEW v_dashboard_inventory_overview AS
SELECT 
    COUNT(*) AS total_parts,
    SUM(stock_quantity) AS total_quantity,
    COALESCE(SUM(stock_quantity * unit_cost), 0) AS total_value,
    COUNT(CASE WHEN stock_quantity <= reorder_level THEN 1 END) AS low_stock_count,
    COUNT(CASE WHEN status = 'ACTIVE' THEN 1 END) AS active_parts,
    whitelabel_id
FROM m_part_master
WHERE deleted = false;

-- ==============================================
-- VIEW: v_dashboard_top_parts (อะไหล่ขายดี)
-- Top selling parts dashboard view.
-- ==============================================
CREATE OR REPLACE VIEW v_dashboard_top_parts AS
SELECT 
    p.id AS part_id,
    p.part_code,
    p.part_name,
    COALESCE(SUM(jps.quantity), 0) AS total_sold,
    COALESCE(SUM(jps.net_price), 0) AS total_revenue,
    p.whitelabel_id,
    ROW_NUMBER() OVER (PARTITION BY p.whitelabel_id ORDER BY SUM(jps.quantity) DESC) AS rank
FROM m_part_master p
LEFT JOIN t_job_part_sales jps ON jps.part_id = p.id
WHERE p.deleted = false
GROUP BY p.id, p.part_code, p.part_name, p.whitelabel_id;

-- ==============================================
-- VIEW: v_dashboard_service_category (บริการแยกประเภท)
-- Service by category dashboard view.
-- ==============================================
CREATE OR REPLACE VIEW v_dashboard_service_category AS
SELECT 
    c.name AS category_name,
    COUNT(js.id) AS service_count,
    COALESCE(SUM(js.net_price), 0) AS total_revenue,
    js.whitelabel_id
FROM t_job_service js
JOIN m_service s ON s.id = js.service_id
JOIN m_category c ON c.id = s.category_id
WHERE js.deleted = false
GROUP BY c.name, js.whitelabel_id;

-- ==============================================
-- VIEW: v_dashboard_financial_summary (สรุปการเงิน)
-- Financial summary dashboard view.
-- ==============================================
CREATE OR REPLACE VIEW v_dashboard_financial_summary AS
WITH invoice_data AS (
    SELECT 
        DATE_TRUNC('month', invoice_date) AS month,
        COALESCE(SUM(total), 0) AS total_invoice,
        whitelabel_id
    FROM t_invoice_adjustment
    WHERE deleted = false AND type = 'INVOICE'
    GROUP BY DATE_TRUNC('month', invoice_date), whitelabel_id
),
payment_data AS (
    SELECT 
        DATE_TRUNC('month', payment_date) AS month,
        COALESCE(SUM(amount), 0) AS total_payment,
        whitelabel_id
    FROM t_payment
    WHERE deleted = false AND status = 'COMPLETED'
    GROUP BY DATE_TRUNC('month', payment_date), whitelabel_id
)
SELECT 
    COALESCE(i.month, p.month) AS month,
    COALESCE(i.total_invoice, 0) AS total_invoice,
    COALESCE(p.total_payment, 0) AS total_payment,
    COALESCE(i.total_invoice, 0) - COALESCE(p.total_payment, 0) AS net_income,
    COALESCE(i.whitelabel_id, p.whitelabel_id) AS whitelabel_id
FROM invoice_data i
FULL OUTER JOIN payment_data p ON i.month = p.month AND i.whitelabel_id = p.whitelabel_id;

-- ==============================================
-- VIEW: v_dashboard_revenue_by_period (รายได้แยกช่วงเวลา)
-- Revenue by period dashboard view.
-- ==============================================
CREATE OR REPLACE VIEW v_dashboard_revenue_by_period AS
SELECT 
    DATE_TRUNC('day', invoice_date) AS period,
    COUNT(*) AS invoice_count,
    COALESCE(SUM(total), 0) AS revenue,
    COALESCE(AVG(total), 0) AS average_revenue,
    whitelabel_id
FROM t_invoice_adjustment
WHERE deleted = false AND type = 'INVOICE'
GROUP BY DATE_TRUNC('day', invoice_date), whitelabel_id;

-- ==============================================
-- ตาราง: d_widget_config (กำหนดค่า Widget บน Dashboard)
-- Dashboard widget configuration.
-- ==============================================
CREATE TABLE IF NOT EXISTS d_widget_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES m_user(id) ON DELETE CASCADE,
    widget_id VARCHAR(50) NOT NULL,              -- เช่น 'sales_overview', 'job_status'
    widget_type VARCHAR(20) NOT NULL,            -- CHART, TABLE, METRIC, MAP
    widget_title VARCHAR(100),
    position INTEGER DEFAULT 0,                  -- ลำดับการแสดงผล
    width INTEGER DEFAULT 4,                     -- ความกว้าง (หน่วย: grid)
    height INTEGER DEFAULT 2,                    -- ความสูง (หน่วย: grid)
    config JSONB,                                -- ตัวเลือกเพิ่มเติม
    is_active BOOLEAN DEFAULT TRUE,
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_d_widget_config_user ON d_widget_config(user_id);
CREATE INDEX idx_d_widget_config_whitelabel ON d_widget_config(whitelabel_id);

-- ==============================================
-- ฟังก์ชัน Refresh Dashboard Views (ใช้ใน Batch Job)
-- Function to refresh dashboard views (used in Batch Jobs).
-- ==============================================
CREATE OR REPLACE FUNCTION refresh_dashboard_views()
RETURNS VOID AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY v_dashboard_sales_overview;
    REFRESH MATERIALIZED VIEW CONCURRENTLY v_dashboard_job_status;
    REFRESH MATERIALIZED VIEW CONCURRENTLY v_dashboard_inventory_overview;
    REFRESH MATERIALIZED VIEW CONCURRENTLY v_dashboard_top_parts;
    REFRESH MATERIALIZED VIEW CONCURRENTLY v_dashboard_service_category;
    REFRESH MATERIALIZED VIEW CONCURRENTLY v_dashboard_financial_summary;
    REFRESH MATERIALIZED VIEW CONCURRENTLY v_dashboard_revenue_by_period;
END;
$$ LANGUAGE plpgsql;
```

---

## 🧠 ระบบ Cache (Redis) สำหรับ Dashboard

### `infrastructure/cache/DashboardCacheService.java`

```java
package com.template.app.modules.dashboard.infrastructure.cache;

import com.template.app.modules.dashboard.presentation.dto.response.*;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
public class DashboardCacheService {

    /*
        ฟังก์ชันนี้ดึงข้อมูลภาพรวม Dashboard จาก Cache (ลดภาระฐานข้อมูลอย่างมาก)
        This function retrieves dashboard overview data from cache (significantly reduces DB load).
        Redis Key: dashboard_overview:{whitelabelId}
    */
    @Cacheable(value = "dashboard_overview", key = "#whitelabelId")
    public DashboardOverviewResponseDTO getDashboardOverview(UUID whitelabelId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลยอดขายแยกช่วงเวลาจาก Cache
        This function retrieves sales by period from cache.
        Redis Key: sales_by_period:{whitelabelId}:{period}
    */
    @Cacheable(value = "sales_by_period", key = "#whitelabelId + ':' + #period")
    public List<RevenueResponseDTO> getSalesByPeriod(UUID whitelabelId, String period) {
        return null;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลอะไหล่ขายดีจาก Cache
        This function retrieves top selling parts from cache.
        Redis Key: top_parts:{whitelabelId}
    */
    @Cacheable(value = "top_parts", key = "#whitelabelId")
    public List<TopPartsResponseDTO> getTopParts(UUID whitelabelId) {
        return null;
    }

    /*
        ฟังก์ชันนี้ลบ Cache Dashboard เมื่อมีข้อมูลเปลี่ยนแปลง (ใช้ใน Batch Job)
        This function evicts dashboard cache when data changes (used in Batch Jobs).
    */
    @CacheEvict(value = {
        "dashboard_overview", 
        "sales_by_period", 
        "top_parts", 
        "job_status_summary",
        "inventory_overview"
    }, allEntries = true)
    public void evictAllDashboardCache() {
        // ลบทุก key ใน Dashboard caches ทั้งหมด / Evict all keys in all Dashboard caches.
    }

    /*
        ฟังก์ชันนี้ลบ Cache Dashboard เฉพาะ Whitelabel ID
        This function evicts dashboard cache for a specific Whitelabel ID.
    */
    @CacheEvict(value = "dashboard_overview", key = "#whitelabelId")
    public void evictDashboardOverview(UUID whitelabelId) {
        // ลบ Cache ของ whitelabel นี้ / Evict cache for this whitelabel.
    }
}
```

---

## ⏱️ ระบบ Rate Limit สำหรับ Dashboard Controller

```java
package com.template.app.modules.dashboard.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.dashboard.application.interfaces.DashboardService;
import com.template.app.modules.dashboard.application.interfaces.ReportService;
import com.template.app.modules.dashboard.presentation.dto.request.DashboardFilterRequestDTO;
import com.template.app.modules.dashboard.presentation.dto.request.ReportRequestDTO;
import com.template.app.modules.dashboard.presentation.dto.response.*;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/dashboard")
@Tag(name = "Dashboard", description = "Dashboard and Analytics APIs")
@RequiredArgsConstructor
public class DashboardController {

    private final DashboardService dashboardService;
    private final ReportService reportService;

    // ========================================================================
    // 1. GET DASHBOARD OVERVIEW (ภาพรวมแดชบอร์ด)
    // ========================================================================

    /*
        API: GET /api/v1/dashboard/overview
        ฟังก์ชันนี้ดึงข้อมูลภาพรวมของระบบ ทั้งยอดขาย, จำนวนงาน, สินค้าคงคลัง
        This function retrieves overall system overview including sales, job count, inventory.
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที (ใช้บ่อยในการแสดง Dashboard)
        Rate Limit: Allows 30 requests per minute (frequently used for Dashboard display).
    */
    @GetMapping("/overview")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get dashboard overview")
    public ResponseEntity<DashboardOverviewResponseDTO> getDashboardOverview() throws SystemGlobalException {
        DashboardOverviewResponseDTO response = dashboardService.getDashboardOverview();
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GET SALES OVERVIEW (ภาพรวมยอดขาย)
    // ========================================================================

    /*
        API: GET /api/v1/dashboard/sales
        ฟังก์ชันนี้แสดงข้อมูลยอดขายแยกตามช่วงเวลา (วัน/เดือน/ปี)
        This function displays sales data grouped by period (day/month/year).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/sales")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get sales overview")
    public ResponseEntity<SalesOverviewResponseDTO> getSalesOverview(
            @RequestParam(required = false) String period,
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate) throws SystemGlobalException {
        SalesOverviewResponseDTO response = dashboardService.getSalesOverview(period, startDate, endDate);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GET JOB STATUS SUMMARY (สรุปสถานะใบงาน)
    // ========================================================================

    /*
        API: GET /api/v1/dashboard/job-status
        ฟังก์ชันนี้แสดงสรุปสถานะของใบงานทั้งหมด (OPEN, IN_PROGRESS, CLOSED, etc.)
        This function shows summary of all job statuses (OPEN, IN_PROGRESS, CLOSED, etc.).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/job-status")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get job status summary")
    public ResponseEntity<List<JobStatusSummaryResponseDTO>> getJobStatusSummary() throws SystemGlobalException {
        List<JobStatusSummaryResponseDTO> response = dashboardService.getJobStatusSummary();
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. GET INVENTORY OVERVIEW (ภาพรวมสินค้าคงคลัง)
    // ========================================================================

    /*
        API: GET /api/v1/dashboard/inventory
        ฟังก์ชันนี้แสดงภาพรวมสินค้าคงคลัง (จำนวนทั้งหมด, มูลค่า, สินค้าใกล้หมด)
        This function shows inventory overview (total items, value, low stock).
        Rate Limit: อนุญาต 30 ครั้งต่อ 1 นาที
        Rate Limit: Allows 30 requests per minute.
    */
    @GetMapping("/inventory")
    @RateLimit(limit = 30, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get inventory overview")
    public ResponseEntity<InventoryOverviewResponseDTO> getInventoryOverview() throws SystemGlobalException {
        InventoryOverviewResponseDTO response = dashboardService.getInventoryOverview();
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 5. GET TOP SELLING PARTS (อะไหล่ขายดี)
    // ========================================================================

    /*
        API: GET /api/v1/dashboard/top-parts
        ฟังก์ชันนี้แสดงรายการอะไหล่ที่ขายดีที่สุด (ใช้สำหรับการตัดสินใจสั่งซื้อ)
        This function shows the best-selling parts (used for purchasing decisions).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @GetMapping("/top-parts")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get top selling parts")
    public ResponseEntity<List<TopPartsResponseDTO>> getTopParts(
            @RequestParam(defaultValue = "10") int limit) throws SystemGlobalException {
        List<TopPartsResponseDTO> response = dashboardService.getTopParts(limit);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 6. GET SERVICE CATEGORY (บริการแยกประเภท)
    // ========================================================================

    /*
        API: GET /api/v1/dashboard/service-category
        ฟังก์ชันนี้แสดงรายได้แยกตามประเภทบริการ (ช่วยให้เห็นบริการที่ทำรายได้มากที่สุด)
        This function shows revenue by service category (helps identify most profitable services).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @GetMapping("/service-category")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get revenue by service category")
    public ResponseEntity<List<ServiceCategoryResponseDTO>> getServiceCategory() throws SystemGlobalException {
        List<ServiceCategoryResponseDTO> response = dashboardService.getServiceCategory();
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 7. GET REVENUE BY PERIOD (รายได้แยกช่วงเวลา)
    // ========================================================================

    /*
        API: GET /api/v1/dashboard/revenue
        ฟังก์ชันนี้แสดงรายได้แยกตามช่วงเวลา (รายวัน/รายเดือน/รายปี) แบบกราฟ
        This function shows revenue by period (day/month/year) in chart format.
        Rate Limit: อนุญาต 25 ครั้งต่อ 1 นาที
        Rate Limit: Allows 25 requests per minute.
    */
    @GetMapping("/revenue")
    @RateLimit(limit = 25, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get revenue by period")
    public ResponseEntity<List<RevenueResponseDTO>> getRevenueByPeriod(
            @RequestParam(defaultValue = "MONTH") String period,
            @RequestParam(required = false) Integer months) throws SystemGlobalException {
        List<RevenueResponseDTO> response = dashboardService.getRevenueByPeriod(period, months);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 8. GET FINANCIAL SUMMARY (สรุปการเงิน)
    // ========================================================================

    /*
        API: GET /api/v1/dashboard/financial
        ฟังก์ชันนี้แสดงสรุปทางการเงิน (รายได้, ค่าใช้จ่าย, กำไรสุทธิ)
        This function shows financial summary (income, expenses, net profit).
        Rate Limit: อนุญาต 15 ครั้งต่อ 1 นาที
        Rate Limit: Allows 15 requests per minute.
    */
    @GetMapping("/financial")
    @RateLimit(limit = 15, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get financial summary")
    public ResponseEntity<FinancialSummaryResponseDTO> getFinancialSummary(
            @RequestParam(required = false) String startDate,
            @RequestParam(required = false) String endDate) throws SystemGlobalException {
        FinancialSummaryResponseDTO response = dashboardService.getFinancialSummary(startDate, endDate);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 9. GET CUSTOM FILTERED DASHBOARD (Dashboard แบบกำหนดเอง)
    // ========================================================================

    /*
        API: POST /api/v1/dashboard/filtered
        ฟังก์ชันนี้ดึงข้อมูล Dashboard ตามตัวกรองที่กำหนดเอง (ลูกค้า, ช่วงเวลา, สาขา)
        This function retrieves dashboard data with custom filters (customer, date range, branch).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @PostMapping("/filtered")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get filtered dashboard data")
    public ResponseEntity<DashboardOverviewResponseDTO> getFilteredDashboard(
            @Valid @RequestBody DashboardFilterRequestDTO request) throws SystemGlobalException {
        DashboardOverviewResponseDTO response = dashboardService.getFilteredDashboard(request);
        return ResponseEntity.ok(response);
    }
}
```

### `ReportController.java` (สร้างรายงาน)

```java
package com.template.app.modules.dashboard.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.dashboard.application.interfaces.ReportService;
import com.template.app.modules.dashboard.presentation.dto.request.ReportRequestDTO;
import com.template.app.modules.dashboard.presentation.dto.response.ReportResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/reports")
@Tag(name = "Reports", description = "Report Generation APIs")
@RequiredArgsConstructor
public class ReportController {

    private final ReportService reportService;

    // ========================================================================
    // 1. GENERATE DAILY REPORT (รายงานรายวัน)
    // ========================================================================

    /*
        API: POST /api/v1/reports/daily
        ฟังก์ชันนี้สร้างรายงานสรุปประจำวัน (ยอดขาย, งานซ่อม, การชำระเงิน)
        This function generates a daily summary report (sales, repairs, payments).
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping("/daily")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Generate daily report")
    public ResponseEntity<ReportResponseDTO> generateDailyReport(@Valid @RequestBody ReportRequestDTO request)
            throws SystemGlobalException {
        ReportResponseDTO response = reportService.generateDailyReport(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 2. GENERATE MONTHLY REPORT (รายงานรายเดือน)
    // ========================================================================

    /*
        API: POST /api/v1/reports/monthly
        ฟังก์ชันนี้สร้างรายงานสรุปรายเดือน (ยอดขายรวม, สินค้าคงคลัง, ลูกหนี้)
        This function generates a monthly summary report (total sales, inventory, receivables).
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping("/monthly")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Generate monthly report")
    public ResponseEntity<ReportResponseDTO> generateMonthlyReport(@Valid @RequestBody ReportRequestDTO request)
            throws SystemGlobalException {
        ReportResponseDTO response = reportService.generateMonthlyReport(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 3. GENERATE YEARLY REPORT (รายงานรายปี)
    // ========================================================================

    /*
        API: POST /api/v1/reports/yearly
        ฟังก์ชันนี้สร้างรายงานสรุปรายปี (สำหรับผู้บริหาร)
        This function generates a yearly summary report (for executives).
        Rate Limit: อนุญาต 5 ครั้งต่อ 1 ชั่วโมง
        Rate Limit: Allows 5 requests per 1 hour.
    */
    @PostMapping("/yearly")
    @RateLimit(limit = 5, duration = 3600, keyType = "USER_ID")
    @Operation(summary = "Generate yearly report")
    public ResponseEntity<ReportResponseDTO> generateYearlyReport(@Valid @RequestBody ReportRequestDTO request)
            throws SystemGlobalException {
        ReportResponseDTO response = reportService.generateYearlyReport(request);
        return ResponseEntity.ok(response);
    }

    // ========================================================================
    // 4. GET REPORT STATUS (สถานะการสร้างรายงาน)
    // ========================================================================

    /*
        API: GET /api/v1/reports/status/{reportId}
        ฟังก์ชันนี้ตรวจสอบสถานะการสร้างรายงาน (ใช้สำหรับรายงานที่ใช้เวลานาน)
        This function checks the status of report generation (for long-running reports).
        Rate Limit: อนุญาต 20 ครั้งต่อ 1 นาที
        Rate Limit: Allows 20 requests per minute.
    */
    @GetMapping("/status/{reportId}")
    @RateLimit(limit = 20, duration = 60, keyType = "USER_ID")
    @Operation(summary = "Get report generation status")
    public ResponseEntity<ReportResponseDTO> getReportStatus(@PathVariable String reportId)
            throws SystemGlobalException {
        ReportResponseDTO response = reportService.getReportStatus(reportId);
        return ResponseEntity.ok(response);
    }
}
```

### `ExportController.java` (ส่งออกข้อมูล)

```java
package com.template.app.modules.dashboard.presentation.controller;

import com.template.app.modules.auth.infrastructure.ratelimit.RateLimit;
import com.template.app.modules.dashboard.application.interfaces.ExportService;
import com.template.app.modules.dashboard.presentation.dto.request.ExportRequestDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/export")
@Tag(name = "Export", description = "Data Export APIs")
@RequiredArgsConstructor
public class ExportController {

    private final ExportService exportService;

    // ========================================================================
    // 1. EXPORT TO EXCEL (.xlsx)
    // ========================================================================

    /*
        API: POST /api/v1/export/excel
        ฟังก์ชันนี้ส่งออกข้อมูลเป็นไฟล์ Excel (.xlsx) สำหรับวิเคราะห์ข้อมูลเพิ่มเติม
        This function exports data as Excel file (.xlsx) for further analysis.
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping("/excel")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Export data to Excel")
    public ResponseEntity<byte[]> exportToExcel(@Valid @RequestBody ExportRequestDTO request)
            throws SystemGlobalException {
        byte[] excelData = exportService.exportToExcel(request);
        return ResponseEntity.ok()
                .header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
                .header("Content-Disposition", "attachment; filename=report_" + request.getReportType() + ".xlsx")
                .body(excelData);
    }

    // ========================================================================
    // 2. EXPORT TO PDF (.pdf)
    // ========================================================================

    /*
        API: POST /api/v1/export/pdf
        ฟังก์ชันนี้ส่งออกข้อมูลเป็นไฟล์ PDF สำหรับพิมพ์หรือส่งให้ผู้บริหาร
        This function exports data as PDF file for printing or executive delivery.
        Rate Limit: อนุญาต 10 ครั้งต่อ 5 นาที
        Rate Limit: Allows 10 requests per 5 minutes.
    */
    @PostMapping("/pdf")
    @RateLimit(limit = 10, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Export data to PDF")
    public ResponseEntity<byte[]> exportToPDF(@Valid @RequestBody ExportRequestDTO request)
            throws SystemGlobalException {
        byte[] pdfData = exportService.exportToPDF(request);
        return ResponseEntity.ok()
                .header("Content-Type", "application/pdf")
                .header("Content-Disposition", "attachment; filename=report_" + request.getReportType() + ".pdf")
                .body(pdfData);
    }

    // ========================================================================
    // 3. EXPORT TO CSV (.csv)
    // ========================================================================

    /*
        API: POST /api/v1/export/csv
        ฟังก์ชันนี้ส่งออกข้อมูลเป็นไฟล์ CSV สำหรับนำไปใช้กับโปรแกรมอื่น
        This function exports data as CSV file for use with other applications.
        Rate Limit: อนุญาต 15 ครั้งต่อ 5 นาที
        Rate Limit: Allows 15 requests per 5 minutes.
    */
    @PostMapping("/csv")
    @RateLimit(limit = 15, duration = 300, keyType = "USER_ID")
    @Operation(summary = "Export data to CSV")
    public ResponseEntity<byte[]> exportToCSV(@Valid @RequestBody ExportRequestDTO request)
            throws SystemGlobalException {
        byte[] csvData = exportService.exportToCSV(request);
        return ResponseEntity.ok()
                .header("Content-Type", "text/csv")
                .header("Content-Disposition", "attachment; filename=report_" + request.getReportType() + ".csv")
                .body(csvData);
    }
}
```

---

## 🧩 ตัวอย่าง Domain Entity และ Enum

### `domain/enums/ReportType.java`

```java
package com.template.app.modules.dashboard.domain.enums;

/*
    ประเภทของรายงาน / Report type.
*/
public enum ReportType {
    DAILY,          // รายงานประจำวัน / Daily report.
    MONTHLY,        // รายงานประจำเดือน / Monthly report.
    YEARLY,         // รายงานประจำปี / Yearly report.
    CUSTOM,         // รายงานตามกำหนดเอง / Custom report.
    INVENTORY,      // รายงานสินค้าคงคลัง / Inventory report.
    FINANCIAL,      // รายงานการเงิน / Financial report.
    PERFORMANCE     // รายงานประสิทธิภาพ / Performance report.
}
```

### `domain/DSalesOverview.java`

```java
package com.template.app.modules.dashboard.domain;

import lombok.Data;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.util.UUID;

@Data
public class DSalesOverview {

    private UUID whitelabelId;          // ID สาขา / Branch ID.
    private LocalDate period;           // ช่วงเวลา / Period.
    private Integer totalInvoices;      // จำนวนใบแจ้งหนี้ทั้งหมด / Total invoices.
    private Integer totalCustomers;     // จำนวนลูกค้าทั้งหมด / Total customers.
    private BigDecimal totalRevenue;    // รายได้รวม / Total revenue.
    private BigDecimal totalOutstanding;// ยอดคงค้างรวม / Total outstanding.
    private BigDecimal averageInvoice;  // ค่าเฉลี่ยต่อใบแจ้งหนี้ / Average per invoice.

    /*
        ฟังก์ชันนี้คำนวณอัตราการเติบโตของรายได้เมื่อเทียบกับงวดก่อนหน้า
        This function calculates revenue growth rate compared to previous period.
    */
    public BigDecimal getGrowthRate(DSalesOverview previous) {
        if (previous == null || previous.getTotalRevenue().compareTo(BigDecimal.ZERO) == 0) {
            return BigDecimal.ZERO;
        }
        return this.totalRevenue.subtract(previous.getTotalRevenue())
                .divide(previous.getTotalRevenue(), 4, java.math.RoundingMode.HALF_UP)
                .multiply(new BigDecimal(100));
    }
}
```

### `domain/DJobStatusSummary.java`

```java
package com.template.app.modules.dashboard.domain;

import lombok.Data;

import java.util.UUID;

@Data
public class DJobStatusSummary {

    private String status;              // สถานะใบงาน / Job status.
    private Long count;                 // จำนวน / Count.
    private UUID whitelabelId;          // ID สาขา / Branch ID.

    /*
        ฟังก์ชันนี้คำนวณสัดส่วนของสถานะนี้เทียบกับทั้งหมด
        This function calculates the percentage of this status compared to total.
    */
    public Double getPercentage(Long total) {
        if (total == null || total == 0) {
            return 0.0;
        }
        return (this.count.doubleValue() / total.doubleValue()) * 100;
    }
}
```

---

## 🔧 ฟังก์ชัน Service ตัวอย่าง (พร้อมคำอธิบาย)

### `application/impl/DashboardServiceImpl.java` (บางส่วน)

```java
package com.template.app.modules.dashboard.application.impl;

import com.template.app.modules.dashboard.application.interfaces.DashboardService;
import com.template.app.modules.dashboard.infrastructure.cache.DashboardCacheService;
import com.template.app.modules.dashboard.infrastructure.repository.DashboardRepository;
import com.template.app.modules.dashboard.presentation.dto.response.DashboardOverviewResponseDTO;
import com.template.app.modules.dashboard.presentation.dto.response.SalesOverviewResponseDTO;
import com.template.app.exception.SystemGlobalException;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.UUID;

@Service
public class DashboardServiceImpl implements DashboardService {

    private final DashboardRepository dashboardRepository;
    private final DashboardCacheService cacheService;

    public DashboardServiceImpl(DashboardRepository dashboardRepository,
                                DashboardCacheService cacheService) {
        this.dashboardRepository = dashboardRepository;
        this.cacheService = cacheService;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลภาพรวม Dashboard โดยใช้ Cache เป็นหลัก (อัปเดตทุก 5 นาที)
        This function retrieves dashboard overview data primarily from cache (updated every 5 minutes).
    */
    @Override
    public DashboardOverviewResponseDTO getDashboardOverview() throws SystemGlobalException {
        // 1. ดึง Whitelabel ID จาก Context / Get Whitelabel ID from context.
        UUID whitelabelId = getWhitelabelId();

        // 2. ลองดึงจาก Cache ก่อน / Try to get from cache first.
        DashboardOverviewResponseDTO cached = cacheService.getDashboardOverview(whitelabelId);
        if (cached != null) {
            // ถ้ามี Cache และไม่เก่าเกินไป ให้คืนค่า Cache / If cache exists and is recent, return it.
            if (isCacheValid(cached.getCacheTimestamp())) {
                return cached;
            }
        }

        // 3. ถ้าไม่มี Cache ให้ดึงจากฐานข้อมูล / If no cache, fetch from database.
        DashboardOverviewResponseDTO response = dashboardRepository.getDashboardOverview(whitelabelId);

        // 4. บันทึก Cache เพื่อใช้ครั้งต่อไป / Store in cache for next time.
        response.setCacheTimestamp(LocalDateTime.now());

        return response;
    }

    /*
        ฟังก์ชันนี้ดึงข้อมูลยอดขายแยกตามช่วงเวลา (DAY, MONTH, YEAR) สำหรับแสดงกราฟ
        This function retrieves sales data by period (DAY, MONTH, YEAR) for chart display.
    */
    @Override
    public SalesOverviewResponseDTO getSalesOverview(String period, String startDate, String endDate)
            throws SystemGlobalException {
        // กำหนดค่าเริ่มต้นหากไม่มีการระบุ / Set defaults if not specified.
        if (period == null || period.isEmpty()) {
            period = "MONTH";
        }

        // สร้าง DateRange จากพารามิเตอร์ / Create DateRange from parameters.
        DateRange dateRange = DateRange.of(startDate, endDate);

        // ดึงข้อมูลจากฐานข้อมูล / Fetch data from database.
        SalesOverviewResponseDTO response = dashboardRepository.getSalesOverview(period, dateRange);

        return response;
    }

    /*
        ฟังก์ชันนี้ตรวจสอบว่า Cache ยังใช้งานได้หรือไม่ (หมดอายุ 5 นาที)
        This function checks if cache is still valid (expires in 5 minutes).
    */
    private boolean isCacheValid(LocalDateTime cacheTimestamp) {
        if (cacheTimestamp == null) {
            return false;
        }
        return cacheTimestamp.plusMinutes(5).isAfter(LocalDateTime.now());
    }
}
```

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
| 7 | 💰 Payment Management | ✅ ครบถ้วน | Payment Record + Receipt + Outstanding Balance + Refund + Cache |
| 8 | 📊 Dashboard & Reports | ✅ ครบถ้วน | Overview + Sales + Inventory + Job Status + Top Parts + Financial + Export |

--- 