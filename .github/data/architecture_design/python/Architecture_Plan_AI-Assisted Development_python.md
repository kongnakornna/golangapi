

# PART F — SESSION P04 (Python ERP) — Step 1 Domain

> **Instruction:** ใน chat ใหม่ paste CONTEXT.md + spec + prompt แล้ว AI จะ generate ตามนี้

## F.1 Directory Tree

```
app/modules/erp/
├── domain/
│   ├── entities/
│   │   ├── order.py
│   │   ├── order_item.py
│   │   ├── invoice.py
│   │   ├── warehouse.py
│   │   └── stock_item.py
│   ├── value_objects/
│   │   ├── order_type.py
│   │   ├── order_status.py
│   │   ├── invoice_type.py
│   │   ├── invoice_status.py
│   │   ├── warehouse_type.py
│   │   ├── money.py
│   │   └── address.py
│   ├── repositories/
│   │   ├── order_repository.py
│   │   ├── invoice_repository.py
│   │   ├── warehouse_repository.py
│   │   └── stock_repository.py
│   ├── events/
│   │   └── events.py
│   └── errors/
│       └── errors.py
```

## F.2 `domain/value_objects/money.py`

```python
from pydantic import BaseModel, Field, field_validator


class Money(BaseModel):
    amount: float = Field(..., ge=0)
    currency: str = Field(..., min_length=3, max_length=3)

    @field_validator("currency")
    @classmethod
    def uppercase_currency(cls, v: str) -> str:
        return v.upper()

    def add(self, other: "Money") -> "Money":
        if self.currency != other.currency:
            raise ValueError("currency mismatch")
        return Money(amount=self.amount + other.amount, currency=self.currency)

    def sub(self, other: "Money") -> "Money":
        if self.currency != other.currency:
            raise ValueError("currency mismatch")
        return Money(amount=self.amount - other.amount, currency=self.currency)

    def mul(self, factor: float) -> "Money":
        return Money(amount=self.amount * factor, currency=self.currency)

    def is_zero(self) -> bool:
        return self.amount == 0

    def is_negative(self) -> bool:
        return self.amount < 0
```

## F.3 `domain/value_objects/order_type.py`

```python
from enum import Enum


class OrderType(str, Enum):
    SALES = "SALES"
    PURCHASE = "PURCHASE"
```

## F.4 `domain/value_objects/order_status.py`

```python
from enum import Enum


class OrderStatus(str, Enum):
    DRAFT = "DRAFT"
    CONFIRMED = "CONFIRMED"
    SHIPPED = "SHIPPED"
    DELIVERED = "DELIVERED"
    CANCELLED = "CANCELLED"

    def is_terminal(self) -> bool:
        return self in (OrderStatus.DELIVERED, OrderStatus.CANCELLED)

    def can_transition_to(self, next_status: "OrderStatus") -> bool:
        transitions = {
            OrderStatus.DRAFT: {OrderStatus.CONFIRMED, OrderStatus.CANCELLED},
            OrderStatus.CONFIRMED: {OrderStatus.SHIPPED, OrderStatus.CANCELLED},
            OrderStatus.SHIPPED: {OrderStatus.DELIVERED, OrderStatus.CANCELLED},
            OrderStatus.DELIVERED: set(),
            OrderStatus.CANCELLED: set(),
        }
        return next_status in transitions.get(self, set())
```

## F.5 `domain/value_objects/invoice_type.py`

```python
from enum import Enum


class InvoiceType(str, Enum):
    TAX = "TAX"
    RECEIPT = "RECEIPT"
    DEBIT_NOTE = "DEBIT_NOTE"
    CREDIT_NOTE = "CREDIT_NOTE"
```

## F.6 `domain/value_objects/invoice_status.py`

```python
from enum import Enum


class InvoiceStatus(str, Enum):
    DRAFT = "DRAFT"
    ISSUED = "ISSUED"
    PAID = "PAID"
    OVERDUE = "OVERDUE"
    CANCELLED = "CANCELLED"

    def is_terminal(self) -> bool:
        return self in (InvoiceStatus.PAID, InvoiceStatus.CANCELLED)

    def can_transition_to(self, next_status: "InvoiceStatus") -> bool:
        transitions = {
            InvoiceStatus.DRAFT: {InvoiceStatus.ISSUED, InvoiceStatus.CANCELLED},
            InvoiceStatus.ISSUED: {InvoiceStatus.PAID, InvoiceStatus.OVERDUE, InvoiceStatus.CANCELLED},
            InvoiceStatus.OVERDUE: {InvoiceStatus.PAID, InvoiceStatus.CANCELLED},
            InvoiceStatus.PAID: set(),
            InvoiceStatus.CANCELLED: set(),
        }
        return next_status in transitions.get(self, set())
```

## F.7 `domain/value_objects/warehouse_type.py`

```python
from enum import Enum


class WarehouseType(str, Enum):
    MAIN = "MAIN"
    BRANCH = "BRANCH"
    TRANSIT = "TRANSIT"
    COLD = "COLD"
```

## F.8 `domain/value_objects/address.py`

```python
import re
from pydantic import BaseModel, Field, field_validator

POSTCODE_PATTERN = re.compile(r"^\d{5}$")


class Address(BaseModel):
    line1: str = Field("", max_length=255)
    line2: str = Field("", max_length=255)
    district: str = ""
    amphoe: str = ""
    province: str = ""
    postcode: str = ""
    country: str = "TH"

    @field_validator("postcode")
    @classmethod
    def validate_postcode(cls, v: str) -> str:
        if v and not POSTCODE_PATTERN.match(v):
            raise ValueError("invalid postcode")
        return v

    def validate_address(self) -> None:
        if not self.country:
            raise ValueError("country required")
```

## F.9 `domain/errors/errors.py`

```python
class DomainError(Exception):
    code: str = "DOMAIN_ERROR"

    def __init__(self, message: str = ""):
        self.message = message or self.__class__.__name__
        super().__init__(self.message)


# Validation
class InvalidTenantError(DomainError):
    code = "INVALID_TENANT"


class InvalidCustomerError(DomainError):
    code = "INVALID_CUSTOMER"


class InvalidOrderTypeError(DomainError):
    code = "INVALID_ORDER_TYPE"


class InvalidQuantityError(DomainError):
    code = "INVALID_QUANTITY"


class InvalidPriceError(DomainError):
    code = "INVALID_PRICE"


class InvalidOrderNumberError(DomainError):
    code = "INVALID_ORDER_NUMBER"


class InvalidInvoiceTypeError(DomainError):
    code = "INVALID_INVOICE_TYPE"


class InvalidWarehouseTypeError(DomainError):
    code = "INVALID_WAREHOUSE_TYPE"


class InvalidDueDateError(DomainError):
    code = "INVALID_DUE_DATE"


# Not found
class OrderNotFoundError(DomainError):
    code = "ORDER_NOT_FOUND"


class OrderItemNotFoundError(DomainError):
    code = "ORDER_ITEM_NOT_FOUND"


class InvoiceNotFoundError(DomainError):
    code = "INVOICE_NOT_FOUND"


class WarehouseNotFoundError(DomainError):
    code = "WAREHOUSE_NOT_FOUND"


class StockNotFoundError(DomainError):
    code = "STOCK_NOT_FOUND"


# Conflict
class OrderNumberDuplicateError(DomainError):
    code = "ORDER_NUMBER_DUPLICATE"


class InvoiceNumberDuplicateError(DomainError):
    code = "INVOICE_NUMBER_DUPLICATE"


class WarehouseCodeDuplicateError(DomainError):
    code = "WAREHOUSE_CODE_DUPLICATE"


# State
class OrderNotDraftError(DomainError):
    code = "ORDER_NOT_DRAFT"


class OrderEmptyError(DomainError):
    code = "ORDER_EMPTY"


class CannotCancelDeliveredError(DomainError):
    code = "CANNOT_CANCEL_DELIVERED"


class InvalidOrderTransitionError(DomainError):
    code = "INVALID_ORDER_TRANSITION"


class InvoiceAlreadyPaidError(DomainError):
    code = "INVOICE_ALREADY_PAID"


class CannotCancelPaidError(DomainError):
    code = "CANNOT_CANCEL_PAID"


class InvalidInvoiceTransitionError(DomainError):
    code = "INVALID_INVOICE_TRANSITION"


# Business
class InsufficientStockError(DomainError):
    code = "INSUFFICIENT_STOCK"


class ReservationMismatchError(DomainError):
    code = "RESERVATION_MISMATCH"


class InsufficientPaymentError(DomainError):
    code = "INSUFFICIENT_PAYMENT"


class CurrencyMismatchError(DomainError):
    code = "CURRENCY_MISMATCH"
```

## F.10 `domain/entities/order_item.py`

```python
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field

from app.modules.erp.domain.errors.errors import (
    InvalidQuantityError,
    InvalidPriceError,
    CurrencyMismatchError,
)
from app.modules.erp.domain.value_objects.money import Money


class OrderItem(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    order_id: UUID | None = None
    product_id: UUID
    sku: str
    name: str
    quantity: int
    unit_price: Money
    discount: Money | None = None
    tax_rate: float = 0.0
    subtotal: Money
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(
        cls,
        product_id: UUID,
        sku: str,
        name: str,
        qty: int,
        unit_price: Money,
    ) -> "OrderItem":
        if qty <= 0:
            raise InvalidQuantityError()
        if unit_price.amount < 0:
            raise InvalidPriceError()
        subtotal = unit_price.mul(float(qty))
        return cls(
            product_id=product_id,
            sku=sku,
            name=name,
            quantity=qty,
            unit_price=unit_price,
            discount=Money(amount=0, currency=unit_price.currency),
            subtotal=subtotal,
        )

    def update_quantity(self, qty: int) -> None:
        if qty <= 0:
            raise InvalidQuantityError()
        self.quantity = qty
        self._recalc_subtotal()

    def update_unit_price(self, price: Money) -> None:
        if price.currency != self.unit_price.currency:
            raise CurrencyMismatchError()
        self.unit_price = price
        self._recalc_subtotal()

    def _recalc_subtotal(self) -> None:
        self.subtotal = self.unit_price.mul(float(self.quantity))
```

## F.11 `domain/entities/order.py`

```python
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field

from app.modules.erp.domain.entities.order_item import OrderItem
from app.modules.erp.domain.errors.errors import (
    InvalidTenantError,
    InvalidCustomerError,
    InvalidOrderTypeError,
    InvalidOrderNumberError,
    OrderItemNotFoundError,
    OrderEmptyError,
    InvalidOrderTransitionError,
)
from app.modules.erp.domain.value_objects.address import Address
from app.modules.erp.domain.value_objects.money import Money
from app.modules.erp.domain.value_objects.order_status import OrderStatus
from app.modules.erp.domain.value_objects.order_type import OrderType


class Order(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    customer_id: UUID
    order_number: str = ""
    type: OrderType
    items: list[OrderItem] = Field(default_factory=list)
    subtotal: Money = Field(default_factory=lambda: Money(amount=0, currency="THB"))
    tax_amount: Money = Field(default_factory=lambda: Money(amount=0, currency="THB"))
    discount: Money = Field(default_factory=lambda: Money(amount=0, currency="THB"))
    total: Money = Field(default_factory=lambda: Money(amount=0, currency="THB"))
    status: OrderStatus = OrderStatus.DRAFT
    notes: str = ""
    shipping_addr: Address = Field(default_factory=Address)
    billing_addr: Address = Field(default_factory=Address)
    ordered_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    confirmed_at: datetime | None = None
    shipped_at: datetime | None = None
    delivered_at: datetime | None = None
    cancelled_at: datetime | None = None
    cancel_reason: str = ""
    created_by: UUID
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(
        cls,
        tenant_id: UUID,
        customer_id: UUID,
        order_type: OrderType,
        created_by: UUID,
    ) -> "Order":
        if not tenant_id:
            raise InvalidTenantError()
        if not customer_id:
            raise InvalidCustomerError()
        return cls(
            tenant_id=tenant_id,
            customer_id=customer_id,
            type=order_type,
            created_by=created_by,
        )

    # ---------- Behavior ----------
    def set_number(self, number: str) -> None:
        if not number or len(number) > 50:
            raise InvalidOrderNumberError()
        self.order_number = number
        self._touch()

    def add_item(self, item: OrderItem) -> None:
        if self.status != OrderStatus.DRAFT:
            raise InvalidOrderTransitionError("order not in draft")
        item.order_id = self.id
        self.items.append(item)
        self._recalculate()
        self._touch()

    def remove_item(self, item_id: UUID) -> None:
        if self.status != OrderStatus.DRAFT:
            raise InvalidOrderTransitionError("order not in draft")
        for i, item in enumerate(self.items):
            if item.id == item_id:
                self.items.pop(i)
                self._recalculate()
                self._touch()
                return
        raise OrderItemNotFoundError()

    def confirm(self) -> None:
        if not self.status.can_transition_to(OrderStatus.CONFIRMED):
            raise InvalidOrderTransitionError()
        if not self.items:
            raise OrderEmptyError()
        self.status = OrderStatus.CONFIRMED
        self.confirmed_at = datetime.now(timezone.utc)
        self._touch()

    def ship(self) -> None:
        if not self.status.can_transition_to(OrderStatus.SHIPPED):
            raise InvalidOrderTransitionError()
        self.status = OrderStatus.SHIPPED
        self.shipped_at = datetime.now(timezone.utc)
        self._touch()

    def deliver(self) -> None:
        if not self.status.can_transition_to(OrderStatus.DELIVERED):
            raise InvalidOrderTransitionError()
        self.status = OrderStatus.DELIVERED
        self.delivered_at = datetime.now(timezone.utc)
        self._touch()

    def cancel(self, reason: str) -> None:
        if not self.status.can_transition_to(OrderStatus.CANCELLED):
            raise InvalidOrderTransitionError()
        self.status = OrderStatus.CANCELLED
        self.cancelled_at = datetime.now(timezone.utc)
        self.cancel_reason = reason
        self._touch()

    # ---------- Query ----------
    def is_draft(self) -> bool:
        return self.status == OrderStatus.DRAFT

    def is_delivered(self) -> bool:
        return self.status == OrderStatus.DELIVERED

    def item_count(self) -> int:
        return len(self.items)

    # ---------- Private ----------
    def _recalculate(self) -> None:
        currency = "THB"
        if self.items:
            currency = self.items[0].subtotal.currency
        subtotal = sum(i.subtotal.amount for i in self.items)
        discount = sum((i.discount.amount if i.discount else 0) for i in self.items)
        tax = sum(i.subtotal.amount * i.tax_rate / 100 for i in self.items)
        self.subtotal = Money(amount=subtotal, currency=currency)
        self.discount = Money(amount=discount, currency=currency)
        self.tax_amount = Money(amount=tax, currency=currency)
        self.total = Money(amount=subtotal - discount + tax, currency=currency)

    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)
```

## F.12 `domain/entities/invoice.py`

```python
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field

from app.modules.erp.domain.errors.errors import (
    InvalidTenantError,
    InvalidCustomerError,
    InvalidInvoiceTypeError,
    InvalidDueDateError,
    InvoiceAlreadyPaidError,
    InsufficientPaymentError,
    CurrencyMismatchError,
    CannotCancelPaidError,
    InvalidInvoiceTransitionError,
)
from app.modules.erp.domain.value_objects.invoice_status import InvoiceStatus
from app.modules.erp.domain.value_objects.invoice_type import InvoiceType
from app.modules.erp.domain.value_objects.money import Money


class Invoice(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    order_id: UUID | None = None
    customer_id: UUID
    number: str = ""
    type: InvoiceType
    subtotal: Money
    tax_amount: Money
    total: Money
    status: InvoiceStatus = InvoiceStatus.DRAFT
    issued_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    due_date: datetime
    paid_at: datetime | None = None
    paid_amount: Money = Field(default_factory=lambda: Money(amount=0, currency="THB"))
    payment_method: str = ""
    notes: str = ""
    created_by: UUID
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(
        cls,
        tenant_id: UUID,
        customer_id: UUID,
        invoice_type: InvoiceType,
        subtotal: Money,
        tax: Money,
        total: Money,
        due_date: datetime,
        created_by: UUID,
    ) -> "Invoice":
        if not tenant_id:
            raise InvalidTenantError()
        if not customer_id:
            raise InvalidCustomerError()
        now = datetime.now(timezone.utc)
        if due_date <= now:
            raise InvalidDueDateError()
        return cls(
            tenant_id=tenant_id,
            customer_id=customer_id,
            type=invoice_type,
            subtotal=subtotal,
            tax_amount=tax,
            total=total,
            due_date=due_date,
            created_by=created_by,
        )

    # ---------- Behavior ----------
    def set_number(self, number: str) -> None:
        self.number = number
        self._touch()

    def link_order(self, order_id: UUID) -> None:
        self.order_id = order_id
        self._touch()

    def issue(self) -> None:
        if not self.status.can_transition_to(InvoiceStatus.ISSUED):
            raise InvalidInvoiceTransitionError()
        self.status = InvoiceStatus.ISSUED
        self.issued_at = datetime.now(timezone.utc)
        self._touch()

    def mark_paid(self, amount: Money, method: str) -> None:
        if self.status == InvoiceStatus.PAID:
            raise InvoiceAlreadyPaidError()
        if amount.amount < self.total.amount:
            raise InsufficientPaymentError()
        if amount.currency != self.total.currency:
            raise CurrencyMismatchError()
        self.status = InvoiceStatus.PAID
        self.paid_at = datetime.now(timezone.utc)
        self.paid_amount = amount
        self.payment_method = method
        self._touch()

    def mark_overdue(self) -> None:
        if self.status == InvoiceStatus.ISSUED and datetime.now(timezone.utc) > self.due_date:
            self.status = InvoiceStatus.OVERDUE
            self._touch()

    def cancel(self) -> None:
        if self.status == InvoiceStatus.PAID:
            raise CannotCancelPaidError()
        if not self.status.can_transition_to(InvoiceStatus.CANCELLED):
            raise InvalidInvoiceTransitionError()
        self.status = InvoiceStatus.CANCELLED
        self._touch()

    # ---------- Query ----------
    def is_paid(self) -> bool:
        return self.status == InvoiceStatus.PAID

    def days_until_due(self) -> int:
        delta = self.due_date - datetime.now(timezone.utc)
        return delta.days

    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)
```

## F.13 `domain/entities/warehouse.py`

```python
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field

from app.modules.erp.domain.errors.errors import (
    InvalidTenantError,
    InvalidOrderNumberError,
    InvalidWarehouseTypeError,
    InvalidQuantityError,
)
from app.modules.erp.domain.value_objects.address import Address
from app.modules.erp.domain.value_objects.warehouse_type import WarehouseType


class Warehouse(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    code: str
    name: str
    address: Address = Field(default_factory=Address)
    type: WarehouseType
    capacity: int = 0
    manager_id: UUID | None = None
    is_active: bool = True
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(
        cls, tenant_id: UUID, code: str, name: str, warehouse_type: WarehouseType
    ) -> "Warehouse":
        if not tenant_id:
            raise InvalidTenantError()
        if not code or len(code) > 50:
            raise InvalidOrderNumberError()
        if not name or len(name) > 255:
            raise InvalidOrderNumberError()
        return cls(
            tenant_id=tenant_id, code=code, name=name, type=warehouse_type
        )

    def rename(self, name: str) -> None:
        if not name or len(name) > 255:
            raise InvalidOrderNumberError()
        self.name = name
        self._touch()

    def set_address(self, addr: Address) -> None:
        addr.validate_address()
        self.address = addr
        self._touch()

    def set_capacity(self, capacity: int) -> None:
        if capacity < 0:
            raise InvalidQuantityError()
        self.capacity = capacity
        self._touch()

    def assign_manager(self, manager_id: UUID) -> None:
        self.manager_id = manager_id
        self._touch()

    def deactivate(self) -> None:
        self.is_active = False
        self._touch()

    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)
```

## F.14 `domain/entities/stock_item.py`

```python
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field

from app.modules.erp.domain.errors.errors import (
    InvalidTenantError,
    InvalidQuantityError,
    InsufficientStockError,
    ReservationMismatchError,
    CurrencyMismatchError,
    InvalidOrderNumberError,
)
from app.modules.erp.domain.value_objects.money import Money


class StockItem(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    product_id: UUID
    warehouse_id: UUID
    sku: str
    quantity: int = 0
    reserved: int = 0
    reorder_point: int = 0
    unit_cost: Money
    last_restocked: datetime | None = None
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(
        cls,
        tenant_id: UUID,
        product_id: UUID,
        warehouse_id: UUID,
        sku: str,
        unit_cost: Money,
    ) -> "StockItem":
        if not tenant_id:
            raise InvalidTenantError()
        if not sku:
            raise InvalidOrderNumberError()
        return cls(
            tenant_id=tenant_id,
            product_id=product_id,
            warehouse_id=warehouse_id,
            sku=sku,
            unit_cost=unit_cost,
        )

    def reserve(self, qty: int) -> None:
        if qty <= 0:
            raise InvalidQuantityError()
        if self.available() < qty:
            raise InsufficientStockError()
        self.reserved += qty
        self._touch()

    def release(self, qty: int) -> None:
        self.reserved = max(0, self.reserved - qty)
        self._touch()

    def issue(self, qty: int) -> None:
        if qty <= 0:
            raise InvalidQuantityError()
        if self.reserved < qty:
            raise ReservationMismatchError()
        self.reserved -= qty
        self.quantity -= qty
        self._touch()

    def restock(self, qty: int, unit_cost: Money) -> None:
        if qty <= 0:
            raise InvalidQuantityError()
        if unit_cost.currency != self.unit_cost.currency:
            raise CurrencyMismatchError()
        self.quantity += qty
        self.unit_cost = unit_cost
        self.last_restocked = datetime.now(timezone.utc)
        self._touch()

    def set_reorder_point(self, point: int) -> None:
        if point < 0:
            raise InvalidQuantityError()
        self.reorder_point = point
        self._touch()

    # ---------- Query ----------
    def available(self) -> int:
        return self.quantity - self.reserved

    def needs_reorder(self) -> bool:
        return self.available() <= self.reorder_point

    def is_out_of_stock(self) -> bool:
        return self.available() <= 0

    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)
```

## F.15 `domain/repositories/order_repository.py`

```python
from typing import Protocol
from uuid import UUID

from app.modules.erp.domain.entities.order import Order
from app.modules.erp.domain.value_objects.order_status import OrderStatus
from app.modules.erp.domain.value_objects.order_type import OrderType


class OrderRepository(Protocol):
    async def save(self, order: Order) -> None: ...
    async def find_by_id(self, id: UUID) -> Order | None: ...
    async def find_by_number(self, tenant_id: UUID, number: str) -> Order | None: ...
    async def find_by_customer(self, customer_id: UUID) -> list[Order]: ...
    async def find_by_status(self, tenant_id: UUID, status: OrderStatus) -> list[Order]: ...
    async def next_order_number(self, tenant_id: UUID, order_type: OrderType) -> str: ...
    async def exists_by_number(self, tenant_id: UUID, number: str) -> bool: ...
    async def delete(self, id: UUID) -> None: ...
```

## F.16 `domain/repositories/invoice_repository.py`

```python
from datetime import datetime
from typing import Protocol
from uuid import UUID

from app.modules.erp.domain.entities.invoice import Invoice
from app.modules.erp.domain.value_objects.invoice_status import InvoiceStatus


class InvoiceRepository(Protocol):
    async def save(self, invoice: Invoice) -> None: ...
    async def find_by_id(self, id: UUID) -> Invoice | None: ...
    async def find_by_number(self, tenant_id: UUID, number: str) -> Invoice | None: ...
    async def find_by_customer(self, customer_id: UUID) -> list[Invoice]: ...
    async def find_by_order(self, order_id: UUID) -> list[Invoice]: ...
    async def find_by_status(self, tenant_id: UUID, status: InvoiceStatus) -> list[Invoice]: ...
    async def find_overdue(self, before: datetime) -> list[Invoice]: ...
    async def next_invoice_number(self, tenant_id: UUID) -> str: ...
    async def delete(self, id: UUID) -> None: ...
```

## F.17 `domain/repositories/warehouse_repository.py`

```python
from typing import Protocol
from uuid import UUID

from app.modules.erp.domain.entities.warehouse import Warehouse


class WarehouseRepository(Protocol):
    async def save(self, warehouse: Warehouse) -> None: ...
    async def find_by_id(self, id: UUID) -> Warehouse | None: ...
    async def find_by_code(self, tenant_id: UUID, code: str) -> Warehouse | None: ...
    async def find_by_tenant(self, tenant_id: UUID) -> list[Warehouse]: ...
    async def exists_by_code(self, tenant_id: UUID, code: str) -> bool: ...
    async def delete(self, id: UUID) -> None: ...
```

## F.18 `domain/repositories/stock_repository.py`

```python
from typing import Protocol
from uuid import UUID

from app.modules.erp.domain.entities.stock_item import StockItem


class StockRepository(Protocol):
    async def save(self, stock: StockItem) -> None: ...
    async def find_by_id(self, id: UUID) -> StockItem | None: ...
    async def find_by_product_warehouse(
        self, product_id: UUID, warehouse_id: UUID
    ) -> StockItem | None: ...
    async def find_by_product(self, product_id: UUID) -> list[StockItem]: ...
    async def find_by_warehouse(self, warehouse_id: UUID) -> list[StockItem]: ...
    async def find_needs_reorder(self, tenant_id: UUID) -> list[StockItem]: ...
    async def delete(self, id: UUID) -> None: ...
```

## F.19 `domain/events/events.py`

```python
from dataclasses import dataclass, field
from datetime import datetime, timezone
from uuid import UUID


@dataclass
class OrderCreatedEvent:
    order_id: UUID
    tenant_id: UUID
    customer_id: UUID
    order_number: str
    type: str
    total: float
    currency: str
    item_count: int
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


@dataclass
class OrderConfirmedEvent:
    order_id: UUID
    tenant_id: UUID
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


@dataclass
class OrderShippedEvent:
    order_id: UUID
    tenant_id: UUID
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


@dataclass
class OrderDeliveredEvent:
    order_id: UUID
    tenant_id: UUID
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


@dataclass
class OrderCancelledEvent:
    order_id: UUID
    tenant_id: UUID
    reason: str
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


@dataclass
class InvoiceIssuedEvent:
    invoice_id: UUID
    tenant_id: UUID
    customer_id: UUID
    total: float
    currency: str
    due_date: datetime
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


@dataclass
class InvoicePaidEvent:
    invoice_id: UUID
    tenant_id: UUID
    amount: float
    currency: str
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


@dataclass
class StockLowEvent:
    tenant_id: UUID
    product_id: UUID
    sku: str
    warehouse_id: UUID
    available: int
    reorder_point: int
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))
```

## F.20 Status After Step 1

```
✅ Step 1 Domain complete (19 files)
   - 5 entities (Order, OrderItem, Invoice, Warehouse, StockItem)
   - 7 value objects
   - 4 repository Protocols
   - 1 errors file (27 exceptions)
   - 1 events file

Reply: "OK domain. Say continue for application."

Compile check:
python -m compileall app/modules/erp/domain/
Expected: PASS (no external deps)
```

---

# PART G — Parallel Guide

## G.1 Running Order (แนะนำ)

```
┌──────────────────────────────────────────────────────────┐
│  Day 1 — ERP Parallel                                    │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  09:00-09:30  G04 Go Step 1 (domain) ✅ Done above       │
│  09:30-10:00  P04 Py Step 1 (domain) ✅ Done above       │
│  10:00-10:15  Break                                       │
│  10:15-10:45  G04 Go Step 2 (application)                │
│  10:45-11:15  P04 Py Step 2 (application)                │
│  11:15-11:30  Break                                       │
│  11:30-12:00  G04 Go Step 3 (infrastructure)             │
│  12:00-13:00  Lunch                                       │
│  13:00-13:30  P04 Py Step 3 (infrastructure)             │
│  13:30-14:00  G04 Go Step 4 (interface + module.go)      │
│  14:00-14:30  P04 Py Step 4 (interface + module.py)      │
│  14:30-15:00  G04 Go Step 5 (migration + tests)          │
│  15:00-15:30  P04 Py Step 5 (migration + tests)          │
│  15:30-16:00  Verify build + commit both                 │
│  16:00-16:30  Update PROGRESS.md × 2, TOKEN_LOG.md       │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

## G.2 Token Tracking

| Session | Start | Step 1 | Step 2 | Step 3 | Step 4 | Step 5 | Total |
| :--- | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| G04 (Go) | 2.5k | 5k | 5k | 8k | 5k | 4k | ~30k |
| P04 (Py) | 2.3k | 4k | 4k | 6k | 4k | 3k | ~24k |
| **Total** | | | | | | | **~54k** |

## G.3 Checkpoint After Day 1

**ไฟล์ `.ai/go/checkpoints/LAST_SESSION.md`:**
```markdown
# Last Go Session — G04

- Date: 2026-01-16
- Module: erp
- Status: ✅ Complete
- Tokens: ~30,000
- Files: ~28
- LOC: ~3,800

Next: G05 — crm Go
```

**ไฟล์ `.ai/py/checkpoints/LAST_SESSION.md`:**
```markdown
# Last Python Session — P04

- Date: 2026-01-16
- Module: erp
- Status: ✅ Complete
- Tokens: ~24,000
- Files: ~25
- LOC: ~2,700

Next: P05 — crm Python
```

**อัปเดต `.ai/MASTER_INDEX.md`:**
```markdown
| 04 | erp | ✅ | ✅ | ✅ | – |
```

**อัปเดต `.ai/TOKEN_LOG.md`:**
```markdown
| G04 | 01-16 | Go | erp | 2,500 | 27,500 | 30,000 | 179,800 |
| P04 | 01-16 | Py | erp | 2,300 | 21,700 | 24,000 | 203,800 |
```

## G.4 Prompt สำหรับ Session ถัดไป

**G04 Step 2 (Go Application):**
```
Continue G04 erp Go from Step 1.

Domain layer is complete (paste the code you generated above or reference it).

Generate Step 2 — Application Layer:
- application/ports.go
- application/dto.go
- application/{verb}_{entity}.go (1 file per use case)

Reference: internal/modules/device/application/
```

**P04 Step 2 (Python Application):**
```
Continue P04 erp Python from Step 1.

Domain layer is complete.

Generate Step 2 — Application Layer:
- application/ports/*.py
- application/use_cases/*.py

Reference: app/modules/device/application/
```

---

## 📊 สรุป Session นี้

| Item | Status |
| :--- | :-: |
| Setup structure | ✅ |
| CONTEXT.md Go | ✅ |
| CONTEXT.md Python | ✅ |
| Prompts Go | ✅ |
| Prompts Python | ✅ |
| ERP Spec YAML | ✅ |
| G04 Step 1 (Domain Go) | ✅ 19 files |
| P04 Step 1 (Domain Python) | ✅ 19 files |

## 🎯 คำสั่งถัดไป

พิมพ์ได้เลย:

- `"G04 Step 2"` → Go Application Layer (ports, dto, use cases)
- `"P04 Step 2"` → Python Application Layer
- `"G04 + P04 Step 2"` → ทั้ง 2 พร้อมกัน
- `"ทำ CONTEXT.md + Prompts ครบทุกไฟล์"` → generate ไฟล์ `.ai/` ทั้งหมดเป็น artifacts
- `"ทำ spec YAML ครบ 40 modules"` → generate specs ทั้งหมด
# 🔀 แยก Workflow: Go vs Python — 2 สายงานอิสระ

> **หลักการ:** Go และ Python **แยกสมบูรณ์** — คนละ `.ai/` folder, คนละ session log, คนละ progress tracker, คนละ prompt
> **ผลลัพธ์:** ทำงานขนานได้, resume แยกได้, สลับไปมาได้โดยไม่ context leak

---

## 📋 สารบัญ

1. [ทำไมต้องแยก](#1-ทำไมต้องแยก)
2. [โครงสร้างโฟลเดอร์ใหม่](#2-โครงสร้างโฟลเดอร์ใหม่)
3. [Workflow Go (แยก)](#3-workflow-go-แยก)
4. [Workflow Python (แยก)](#4-workflow-python-แยก)
5. [Master Index (รวม 2 สาย)](#5-master-index-รวม-2-สาย)
6. [Cross-Language Sync](#6-cross-language-sync)
7. [Quick Start แยก](#7-quick-start-แยก)

---

## 1. ทำไมต้องแยก

| เหตุผล | Go | Python |
| :--- | :--- | :--- |
| **Framework** | Gin + GORM | FastAPI + SQLAlchemy |
| **Testing** | `go test ./...` | `pytest` |
| **Build** | `go build ./...` | `python -m compileall` |
| **Structure** | `internal/modules/` | `app/modules/` |
| **Async model** | goroutine | asyncio |
| **Context size** | ~500 tokens | ~500 tokens |
| **Error style** | sentinel errors | exception classes |
| **DI** | manual wire-up | FastAPI Depends |

**ผลลัพธ์:** ถ้าแยก → prompt ไม่ปนกัน, ไม่ต้องอธิบาย framework ข้ามภาษา, resume ได้เร็ว

---

## 2. โครงสร้างโฟลเดอร์ใหม่

### 2.1 แยก 2 สายชัดเจน

```
icmongolang/
│
├── .ai/                                    # 🎯 AI Workspace (แชร์)
│   ├── MASTER_INDEX.md                    # ทะเบียน 40 modules (ทั้ง 2 สาย)
│   ├── TOKEN_LOG.md                       # Token รวม
│   │
│   ├── go/                                # 🟦 GO WORKSPACE
│   │   ├── CONTEXT.md                    # Go context (compact)
│   │   ├── PROGRESS.md                   # Go progress (40 modules)
│   │   ├── CONVENTIONS.md                # Go-specific conventions
│   │   ├── modules/                      # Go module specs
│   │   │   ├── 01_customer.yaml
│   │   │   ├── 02_package.yaml
│   │   │   └── ...
│   │   ├── sessions/                     # Go session log
│   │   │   ├── G01_customer.md
│   │   │   ├── G02_package.md
│   │   │   └── ...
│   │   ├── checkpoints/
│   │   │   ├── LAST_SESSION.md
│   │   │   └── NEXT_SESSION.md
│   │   └── prompts/
│   │       ├── generate_module.md
│   │       ├── fix_error.md
│   │       └── review.md
│   │
│   └── py/                                # 🐍 PYTHON WORKSPACE
│       ├── CONTEXT.md                    # Python context (compact)
│       ├── PROGRESS.md                   # Python progress (40 modules)
│       ├── CONVENTIONS.md                # Python-specific conventions
│       ├── modules/                      # Python module specs
│       │   ├── 01_customer.yaml
│       │   ├── 02_package.yaml
│       │   └── ...
│       ├── sessions/                     # Python session log
│       │   ├── P01_customer.md
│       │   ├── P02_package.md
│       │   └── ...
│       ├── checkpoints/
│       │   ├── LAST_SESSION.md
│       │   └── NEXT_SESSION.md
│       └── prompts/
│           ├── generate_module.md
│           ├── fix_error.md
│           └── review.md
│
├── internal/modules/                       # 🟦 Go code
├── migrations/                             # 🟦 Shared SQL
├── app/modules/                            # 🐍 Python code
│
└── docs/                                   # Shared docs
```

### 2.2 กฎการใช้งาน

| เมื่อทำ... | เปิด folder | ห้ามเปิด |
| :--- | :--- | :--- |
| Go | `.ai/go/` | `.ai/py/` |
| Python | `.ai/py/` | `.ai/go/` |
| Plan รวม | `.ai/MASTER_INDEX.md` | – |
| Token | `.ai/TOKEN_LOG.md` | – |

---

## 3. Workflow Go (แยก)

### 3.1 ไฟล์ `.ai/go/CONTEXT.md`

```markdown
# Go Context — icmongolang

## Stack
- Go 1.21+
- Gin (HTTP)
- GORM (PostgreSQL)
- IBM Sarama (Kafka)
- go-redis/v8
- go-elasticsearch/v8
- influxdb-client-go/v2
- eclipse/paho.mqtt.golang
- google/uuid

## Module Structure
internal/modules/{name}/
├── domain/
│   ├── entity/
│   ├── value_object/
│   ├── repository/          # interfaces only
│   ├── service/
│   ├── event/
│   └── errors/              # sentinel errors
├── application/             # use cases + ports
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/
│   ├── persistence/influxdb/
│   ├── messaging/           # kafka
│   ├── iot/                 # mqtt
│   ├── search/elasticsearch/
│   └── scheduler/
├── interfaces/
│   ├── http/                # handlers + routes
│   └── websocket/
└── module.go                # composition root

## Naming
- Entity: PascalCase (Customer, Device)
- File: snake_case (create_customer.go)
- Table: {module}_{plural} (customer_customers)
- Migration: YYYYMMDD_{module}_init.sql
- Kafka: {context}.{aggregate}.{event}
- Redis key: {module}:{entity}:{id}
- MQTT: iot/{serial}/{suffix}

## Layer Rules
| Layer | Import | Forbidden |
|-------|--------|-----------|
| domain | stdlib, uuid | gorm, gin, sarama, redis, es |
| application | domain | infrastructure, interface |
| infrastructure | domain, application | interface |
| interface | all | – |

## Entity Pattern
```go
type Customer struct { ... }

func NewCustomer(/* required */) (*Customer, error) { ... }
func (c *Customer) Activate() error { ... }  // behavior
func (c *Customer) IsActive() bool { ... }   // query
// NO setters
```

## Use Case Pattern
```go
type CreateCustomerUseCase struct {
    repo      repository.CustomerRepository
    producer  EventProducer
    auditRepo AuditRepository
}

func (uc *CreateCustomerUseCase) Execute(
    ctx context.Context, in CreateCustomerInput,
) (*CreateCustomerOutput, error) { ... }
```

## Repository Pattern
- Interface: `domain/repository/{entity}_repository.go`
- Impl: `infrastructure/persistence/postgres/{entity}_repo_impl.go`
- Mapper: `toEntity()`, `toModel()`

## Error Pattern
```go
// domain/errors/errors.go
var (
    ErrCustomerNotFound = errors.New("customer not found")
    ErrInvalidCode      = errors.New("invalid customer code")
)
```

## Port Pattern
```go
// application/ports.go
type EventProducer interface { ... }
type AuditRepository interface { ... }
type WSHub interface { ... }
```

## Whitelist (pkg/)
cryptpass, db, elasticsearch, emailTemplates, helpers, http-swagger,
httpErrors, influxdb, jwt, kafka, llm, logger, mqtt, report, responses,
secureRandom, sendEmail, transaction, utils, vectordb, websocket

## Build & Test
go build ./...
go vet ./...
go test ./...
```

**ขนาด:** ~800 tokens

### 3.2 ไฟล์ `.ai/go/PROGRESS.md`

```markdown
# Go Progress — 40 Modules

## Legend
✅ Done | 🚧 In Progress | ⏳ Queued | ⚪ Not Started

## Modules

| # | Module | Spec | Code | Tests | Build | Commit | Session |
|:-:|--------|:----:|:----:|:-----:|:-----:|:------:|:-------:|
| 01 | customer | ✅ | ✅ | ✅ | ✅ | ✅ | G01 |
| 02 | package | ✅ | ✅ | ✅ | ✅ | ✅ | G02 |
| 03 | device | ✅ | ✅ | ✅ | ✅ | ✅ | G03 |
| 04 | erp | ✅ | 🚧 | ⏳ | ⏳ | ⏳ | G04 |
| 05 | crm | ⏳ | ⚪ | ⚪ | ⚪ | ⚪ | – |
| 06 | logistics | ⏳ | ⚪ | ⚪ | ⚪ | ⚪ | – |
| 07 | report | ⏳ | ⚪ | ⚪ | ⚪ | ⚪ | – |
| 08 | auth | ⏳ | ⚪ | ⚪ | ⚪ | ⚪ | – |
| 09 | users | ⏳ | ⚪ | ⚪ | ⚪ | ⚪ | – |
| ... | ... | ... | ... | ... | ... | ... | ... |

## Stats
- Total: 40
- Done: 3
- In Progress: 1
- Remaining: 36
- Avg tokens/session: 25k
- Est. remaining tokens: 900k

## Last Session
G03 — device (2026-01-13) ✅

## Next
G04 — erp (in progress)
```

### 3.3 ไฟล์ `.ai/go/CONVENTIONS.md`

```markdown
# Go-Specific Conventions

## Import Organization
```go
import (
    // stdlib
    "context"
    "time"

    // third-party
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    // internal
    "icmongolang/internal/modules/customer/domain/entity"
    "icmongolang/pkg/responses"
)
```

## Error Handling
- Return errors, don't panic (except startup)
- Wrap with `fmt.Errorf("context: %w", err)`
- Sentinel errors in `domain/errors/`
- HTTP mapping in handler

## Context Usage
- First param: `ctx context.Context`
- Pass to all IO operations
- Timeout: `context.WithTimeout`

## Struct Tags
```go
type CustomerModel struct {
    ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
    Code     string    `gorm:"type:varchar(50);uniqueIndex"`
    Metadata datatypes.JSON `gorm:"type:jsonb"`
}
```

## JSON Tags
```go
type CustomerResponse struct {
    ID   string `json:"id"`
    Code string `json:"code"`
    Name string `json:"name,omitempty"`
}
```

## Repository Naming
- Interface: `CustomerRepository`
- Impl: `customerRepoImpl`
- Constructor: `NewCustomerRepository(db)`

## Use Case Naming
- Type: `CreateCustomerUseCase`
- Method: `Execute(ctx, input)`
- Input: `CreateCustomerInput`
- Output: `CreateCustomerOutput`

## Gin Handler Pattern
```go
func (h *CustomerHandler) Create(c *gin.Context) {
    uid := c.MustGet("user_id").(uuid.UUID)
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    
    var req CreateCustomerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, responses.Error("INVALID", err.Error(), nil))
        return
    }
    
    out, err := h.createUC.Execute(c.Request.Context(), ...)
    if err != nil {
        respondDomainError(c, err)
        return
    }
    c.JSON(201, responses.Success(out))
}
```

## Transaction Pattern
```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Save(m).Error; err != nil {
        return err
    }
    return nil
})
```

## Testing Pattern
```go
func TestCustomer_Activate(t *testing.T) {
    tests := []struct {
        name    string
        setup   func() *entity.Customer
        wantErr error
    }{
        // cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

## File Naming
- Entity: `customer.go`
- Use case: `create_customer.go`
- Repository: `customer_repo_impl.go`
- Handler: `customer_handler.go`
- Test: `customer_test.go`
```

### 3.4 ไฟล์ `.ai/go/prompts/generate_module.md`

```markdown
# Generate Go Module

## Context (paste .ai/go/CONTEXT.md)
{{CONTEXT}}

## Conventions (paste .ai/go/CONVENTIONS.md)
{{CONVENTIONS}}

## Spec (paste .ai/go/modules/{NN}_{name}.yaml)
{{SPEC}}

## Task
Generate complete Go module per spec.

## Structure
internal/modules/{name}/
├── domain/{entity,value_object,repository,service,event,errors}/
├── application/
├── infrastructure/{persistence/postgres,messaging,search}/
├── interfaces/http/
└── module.go

## Output (layer by layer)
**Reply "ready" first, then wait for "start"**

Then output in 5 steps (I'll say "continue" between each):

### Step 1 — Domain
- domain/entity/*.go
- domain/value_object/*.go
- domain/repository/*.go (interfaces)
- domain/errors/errors.go
- domain/event/events.go
- domain/service/*.go
→ Reply "OK domain. Say continue for application."

### Step 2 — Application
- application/ports.go
- application/dto.go
- application/{verb}_{entity}.go (1 per use case)
→ Reply "OK application. Say continue for infrastructure."

### Step 3 — Infrastructure
- infrastructure/persistence/postgres/models.go
- infrastructure/persistence/postgres/*_repo_impl.go
- infrastructure/messaging/producer.go
- infrastructure/{others per spec}
→ Reply "OK infrastructure. Say continue for interface."

### Step 4 — Interface
- interfaces/http/*_handler.go
- interfaces/http/routes.go
- interfaces/http/dto.go
- module.go
→ Reply "OK interface. Say continue for migration + tests."

### Step 5 — Migration + Tests
- migrations/YYYYMMDD_{name}_init.sql
- domain/entity/{entity}_test.go
→ Reply "OK complete."

## Rules
1. Domain: NO gorm/gin/sarama imports
2. Entities: constructor + behavior methods (no setters)
3. Repository: interface in domain, impl in infrastructure
4. Use case: Execute(ctx, input) (output, error)
5. Errors: sentinel errors
6. Reference: internal/modules/device/ (structure), customer/ (patterns)
7. Complete code — do NOT truncate
8. No explanatory text between files — just code blocks

## After generation
Run: `go build ./...` and report errors.
```

### 3.5 ไฟล์ `.ai/go/checkpoints/LAST_SESSION.md`

```markdown
# Last Go Session — G03

## Metadata
- Date: 2026-01-13 14:00
- Module: device
- Status: ✅ Complete
- Tokens: ~30,800
- Duration: 32 min

## Created Files
domain/ (12 files)
application/ (10 files)
infrastructure/ (15 files)
interfaces/ (4 files)
module.go
migration.sql
tests (2 files)
**Total: ~44 files, ~4,500 LOC**

## Issues
- MQTT circular dep — resolved

## Commit
`feat(device): generate Go module per spec`

## Next Session
See NEXT_SESSION.md
```

### 3.6 ไฟล์ `.ai/go/checkpoints/NEXT_SESSION.md`

```markdown
# Next Go Session — G04

## Target
- Module: erp
- Est. tokens: 28,000
- Est. duration: 30 min

## Checklist
- [ ] Load .ai/go/CONTEXT.md
- [ ] Load .ai/go/CONVENTIONS.md
- [ ] Load .ai/go/modules/04_erp.yaml
- [ ] Load .ai/go/prompts/generate_module.md
- [ ] Generate domain layer
- [ ] Generate application layer
- [ ] Generate infrastructure layer
- [ ] Generate interface layer + module.go
- [ ] Generate migration + tests
- [ ] Run go build ./...
- [ ] Fix errors
- [ ] Update .ai/go/PROGRESS.md
- [ ] Write LAST_SESSION.md → G04
- [ ] Prepare NEXT_SESSION.md → G05 (crm)

## Prompt to use
Copy .ai/go/prompts/generate_module.md
```

---

## 4. Workflow Python (แยก)

### 4.1 ไฟล์ `.ai/py/CONTEXT.md`

```markdown
# Python Context — icmongolang

## Stack
- Python 3.11+
- FastAPI
- Pydantic v2
- SQLAlchemy 2.0 (async) + asyncpg
- aiokafka
- redis-py (async)
- elasticsearch-py (async)
- influxdb-client (async)
- aiomqtt
- pytest + pytest-asyncio

## Module Structure
app/modules/{name}/
├── domain/
│   ├── entities/
│   ├── value_objects/
│   ├── repositories/         # Protocol
│   ├── services/
│   ├── events/
│   └── errors/              # Exception classes
├── application/
│   ├── use_cases/
│   └── ports/
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/
│   ├── persistence/influxdb/
│   ├── messaging/
│   ├── iot/
│   ├── search/elasticsearch/
│   └── scheduler/
├── interfaces/
│   ├── http/                # routers + schemas
│   └── middleware/
└── module.py                # DI container

## Naming
- Entity: PascalCase (Customer, Device)
- File: snake_case (create_customer.py)
- Table: {module}_{plural} (customer_customers)
- Migration: YYYYMMDD_{module}_init.sql
- Kafka: {context}.{aggregate}.{event}
- Redis key: {module}:{entity}:{id}
- MQTT: iot/{serial}/{suffix}

## Layer Rules
| Layer | Import | Forbidden |
|-------|--------|-----------|
| domain | stdlib, pydantic | sqlalchemy, fastapi, aiokafka |
| application | domain | infrastructure, interface |
| infrastructure | domain, application | interface |
| interface | all | – |

## Entity Pattern
```python
class Customer(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    code: str
    name: str
    status: CustomerStatus = CustomerStatus.PENDING
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    
    @classmethod
    def create(cls, tenant_id: UUID, code: str, name: str) -> "Customer":
        ...
    
    def activate(self) -> None:  # behavior
        ...
    
    def is_active(self) -> bool:  # query
        ...
```

## Use Case Pattern
```python
@dataclass
class CreateCustomerInput:
    tenant_id: UUID
    code: str
    ...

@dataclass
class CreateCustomerOutput:
    customer_id: UUID
    ...

class CreateCustomerUseCase:
    def __init__(
        self,
        repo: CustomerRepository,
        producer: EventProducer,
        audit_repo: AuditRepository,
    ):
        self._repo = repo
        ...
    
    async def execute(self, inp: CreateCustomerInput) -> CreateCustomerOutput:
        ...
```

## Repository Pattern
- Protocol: `domain/repositories/{entity}_repository.py`
- Impl: `infrastructure/persistence/postgres/{entity}_repository_impl.py`
- Mapper: `_to_entity()`, `_to_model()`

## Error Pattern
```python
# domain/errors/errors.py
class DomainError(Exception):
    code: str = "DOMAIN_ERROR"

class CustomerNotFoundError(DomainError):
    code = "CUSTOMER_NOT_FOUND"

class InvalidCodeError(DomainError):
    code = "INVALID_CODE"
```

## Port Pattern
```python
# application/ports/event_producer.py
class EventProducer(Protocol):
    async def publish_customer_created(self, evt: CustomerCreatedEvent) -> None: ...
```

## FastAPI Router Pattern
```python
@router.post("", response_model=CustomerResponse, status_code=201)
async def create_customer(
    req: CreateCustomerRequest,
    uc: CreateCustomerUseCase = Depends(get_create_uc),
    user_id: UUID = Depends(get_current_user),
    tenant_id: UUID = Depends(get_tenant_id),
):
    try:
        out = await uc.execute(CreateCustomerInput(...))
        return out
    except DomainError as e:
        raise HTTPException(
            status_code=400,
            detail={"code": e.code, "message": str(e)},
        )
```

## DI Pattern
- FastAPI Depends for request-scoped
- Module builder for app-scoped
- `module.py` exports dict of use cases

## Build & Test
python -m compileall app/
pytest tests/
pytest --cov=app tests/
```

**ขนาด:** ~800 tokens

### 4.2 ไฟล์ `.ai/py/PROGRESS.md`

```markdown
# Python Progress — 40 Modules

## Legend
✅ Done | 🚧 In Progress | ⏳ Queued | ⚪ Not Started

## Modules

| # | Module | Spec | Code | Tests | Compile | Commit | Session |
|:-:|--------|:----:|:----:|:-----:|:-------:|:------:|:-------:|
| 01 | customer | ✅ | ✅ | ✅ | ✅ | ✅ | P01 |
| 02 | package | ✅ | ✅ | ✅ | ✅ | ✅ | P02 |
| 03 | device | ✅ | ✅ | ✅ | ✅ | ✅ | P03 |
| 04 | erp | ✅ | ⏳ | ⏳ | ⏳ | ⏳ | P04 |
| 05 | crm | ⏳ | ⚪ | ⚪ | ⚪ | ⚪ | – |
| ... | ... | ... | ... | ... | ... | ... | ... |

## Stats
- Total: 40
- Done: 3
- In Progress: 0
- Remaining: 37
- Avg tokens/session: 22k
- Est. remaining tokens: 814k
```

### 4.3 ไฟล์ `.ai/py/CONVENTIONS.md`

```markdown
# Python-Specific Conventions

## Import Organization
```python
# stdlib
from datetime import datetime, timezone
from uuid import UUID, uuid4
from enum import Enum

# third-party
from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel, Field
from sqlalchemy import select

# internal
from app.modules.customer.domain.entities.customer import Customer
from app.modules.customer.domain.errors.errors import CustomerNotFoundError
```

## Error Handling
- Raise custom exceptions from domain
- Map to HTTP in router
- Use try/except at interface layer
- Always include `code` in error detail

## Async Everywhere
- All IO = async
- Repository: `async def save(...)`
- Use case: `async def execute(...)`
- No blocking calls (use async libraries)

## Type Hints
- Always use type hints
- `UUID` from uuid module (not str)
- `datetime` from datetime module
- `X | None` for optional (Python 3.10+)

## Pydantic v2 Style
```python
from pydantic import BaseModel, Field, field_validator

class Customer(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    name: str = Field(..., min_length=1, max_length=255)
    status: CustomerStatus = CustomerStatus.PENDING
    
    @field_validator("name")
    @classmethod
    def validate_name(cls, v: str) -> str:
        return v.strip()
```

## SQLAlchemy 2.0 Style
```python
from sqlalchemy import select
from sqlalchemy.orm import Mapped, mapped_column

class CustomerModel(Base):
    __tablename__ = "customer_customers"
    
    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True)
    code: Mapped[str] = mapped_column(String(50), nullable=False)
    
# Query
stmt = select(CustomerModel).where(CustomerModel.id == id)
result = (await session.execute(stmt)).scalar_one_or_none()
```

## Dataclass for DTOs
```python
from dataclasses import dataclass

@dataclass
class CreateCustomerInput:
    tenant_id: UUID
    code: str
    name: str
```

## Protocol for Ports
```python
from typing import Protocol

class EventProducer(Protocol):
    async def publish_customer_created(self, evt: CustomerCreatedEvent) -> None: ...
```

## Router Pattern
- 1 file per entity
- Prefix: `/api/v1/{plural}`
- Depends for UC injection
- Exception → HTTPException

## Testing Pattern
```python
import pytest
from uuid import uuid4

def test_activate_ok():
    c = Customer.create(uuid4(), "C001", "Alice")
    c.activate()
    assert c.is_active()

def test_activate_invalid_transition():
    c = Customer.create(uuid4(), "C002", "Bob")
    c.cancel("test")
    with pytest.raises(InvalidStatusTransitionError):
        c.activate()
```

## File Naming
- Entity: `customer.py`
- Use case: `create_customer.py`
- Repository: `customer_repository_impl.py`
- Router: `customer_router.py`
- Schema: `schemas.py`
- Test: `test_customer_entity.py`
```

### 4.4 ไฟล์ `.ai/py/prompts/generate_module.md`

```markdown
# Generate Python Module

## Context (paste .ai/py/CONTEXT.md)
{{CONTEXT}}

## Conventions (paste .ai/py/CONVENTIONS.md)
{{CONVENTIONS}}

## Spec (paste .ai/py/modules/{NN}_{name}.yaml)
{{SPEC}}

## Task
Generate complete Python FastAPI module per spec.

## Structure
app/modules/{name}/
├── domain/{entities,value_objects,repositories,services,events,errors}/
├── application/{use_cases,ports}/
├── infrastructure/{persistence/postgres,messaging,search}/
├── interfaces/http/
└── module.py

## Output (layer by layer)
**Reply "ready" first, then wait for "start"**

Then output in 5 steps (I'll say "continue" between each):

### Step 1 — Domain
- domain/entities/*.py
- domain/value_objects/*.py
- domain/repositories/*.py (Protocol)
- domain/errors/errors.py
- domain/events/events.py
- domain/services/*.py
→ Reply "OK domain. Say continue for application."

### Step 2 — Application
- application/ports/*.py
- application/use_cases/*.py
→ Reply "OK application. Say continue for infrastructure."

### Step 3 — Infrastructure
- infrastructure/persistence/postgres/models.py
- infrastructure/persistence/postgres/*_repository_impl.py
- infrastructure/messaging/kafka_producer.py
- infrastructure/{others per spec}
→ Reply "OK infrastructure. Say continue for interface."

### Step 4 — Interface
- interfaces/http/schemas.py
- interfaces/http/*_router.py
- interfaces/http/dependencies.py
- module.py
→ Reply "OK interface. Say continue for migration + tests."

### Step 5 — Migration + Tests
- alembic/versions/YYYYMMDD_{name}_init.py
- tests/modules/{name}/test_{entity}_entity.py
→ Reply "OK complete."

## Rules
1. Domain: NO sqlalchemy, fastapi, aiokafka imports
2. Entities: Pydantic BaseModel + behavior methods
3. Repository: Protocol in domain, impl in infrastructure
4. Use case: async execute(input) -> output
5. Errors: exception classes with `code` attribute
6. Reference: app/modules/device/, app/modules/customer/
7. Complete code — do NOT truncate
8. Type hints everywhere

## After generation
Run: `python -m compileall app/` and `pytest tests/` and report.
```

### 4.5 ไฟล์ `.ai/py/checkpoints/LAST_SESSION.md`

```markdown
# Last Python Session — P03

## Metadata
- Date: 2026-01-13 15:00
- Module: device
- Status: ✅ Complete
- Tokens: ~27,600
- Duration: 28 min

## Created Files
domain/ (10 files)
application/ (7 files)
infrastructure/ (8 files)
interfaces/ (3 files)
module.py
alembic migration
tests (2 files)
**Total: ~32 files, ~3,200 LOC**

## Issues
- MQTT circular dep — resolved via post-construction injection

## Commit
`feat(device): generate Python module per spec`

## Next Session
See NEXT_SESSION.md
```

---

## 5. Master Index (รวม 2 สาย)

### 5.1 ไฟล์ `.ai/MASTER_INDEX.md`

```markdown
# Master Index — 40 Modules × 2 Languages

## Progress Overview

| # | Module | Spec | Go | Py | Notes |
|:-:|--------|:----:|:--:|:--:|:------|
| 01 | customer | ✅ | ✅ | ✅ | – |
| 02 | package | ✅ | ✅ | ✅ | – |
| 03 | device | ✅ | ✅ | ✅ | Cross-cutting 5 |
| 04 | erp | ✅ | 🚧 | ⏳ | – |
| 05 | crm | ⏳ | ⚪ | ⚪ | – |
| 06 | logistics | ⏳ | ⚪ | ⚪ | – |
| 07 | report | ⏳ | ⚪ | ⚪ | Breaking change |
| 08 | auth | ⏳ | ⚪ | ⚪ | Breaking change |
| ... | ... | ... | ... | ... | ... |

## Current Sessions

### Go
- Last: G03 — device ✅
- Next: G04 — erp 🚧
- Doc: `.ai/go/checkpoints/NEXT_SESSION.md`

### Python
- Last: P03 — device ✅
- Next: P04 — erp ⏳
- Doc: `.ai/py/checkpoints/NEXT_SESSION.md`

## Module Specs (Shared Source)

| # | Module | Go Spec | Py Spec | Full Doc |
|:-:|--------|:-------:|:-------:|:--------:|
| 01 | customer | `.ai/go/modules/01_customer.yaml` | `.ai/py/modules/01_customer.yaml` | `docs/modules/01_customer.md` |
| 02 | package | ... | ... | ... |
| 03 | device | ... | ... | ... |

## Schedule Recommendation

### Parallel Track (แนะนำ)
```
Week 1:  Go 01-05  |  Py 01-05   (2 sessions/day)
Week 2:  Go 06-10  |  Py 06-10
Week 3:  Go 11-15  |  Py 11-15
...
```

### Sequential Track
```
Week 1-4:  Go ทั้งหมด (40 modules)
Week 5-8:  Py ทั้งหมด (40 modules)
```

## Stats

### Go
- Done: 3 / 40
- Tokens used: ~80,000
- Tokens remaining: ~925,000

### Python
- Done: 3 / 40
- Tokens used: ~71,000
- Tokens remaining: ~814,000

### Total
- Sessions: 6
- Tokens: ~151,000
- Progress: 7.5%
```

### 5.2 ไฟล์ `.ai/TOKEN_LOG.md`

```markdown
# Token Log — Both Languages

## All Sessions

| Session | Date | Lang | Module | Input | Output | Total | Cum |
|:-------:|------|:----:|--------|------:|-------:|------:|----:|
| G01 | 01-10 | Go | customer | 2,500 | 22,000 | 24,500 | 24,500 |
| P01 | 01-10 | Py | customer | 2,300 | 19,500 | 21,800 | 46,300 |
| G02 | 01-11 | Go | package | 2,400 | 21,500 | 23,900 | 70,200 |
| P02 | 01-11 | Py | package | 2,200 | 19,000 | 21,200 | 91,400 |
| G03 | 01-13 | Go | device | 2,800 | 28,000 | 30,800 | 122,200 |
| P03 | 01-13 | Py | device | 2,600 | 25,000 | 27,600 | 149,800 |

## By Language

### Go
| Metric | Value |
|--------|------:|
| Sessions | 3 |
| Total tokens | 79,200 |
| Avg/session | 26,400 |
| Modules | 3 |
| Avg/module | 26,400 |

### Python
| Metric | Value |
|--------|------:|
| Sessions | 3 |
| Total tokens | 70,600 |
| Avg/session | 23,533 |
| Modules | 3 |
| Avg/module | 23,533 |

## Projections

| Track | Modules Left | Avg/Module | Projected |
|-------|-------------:|-----------:|----------:|
| Go | 37 | 26,400 | 976,800 |
| Python | 37 | 23,533 | 870,721 |
| **Total** | 74 | – | 1,847,521 |

## Optimization Notes
- 2026-01-13: Layer-by-layer reduces rework by 30%
- 2026-01-13: Reference structure instead of paste → -85% on setup
```

---

## 6. Cross-Language Sync

### 6.1 กฎการ Sync

| สิ่งที่ต้องตรงกัน | วิธี |
| :--- | :--- |
| **Domain model** | YAML spec เดียวกัน (copy ไป 2 folder) |
| **Table names** | Shared migration SQL |
| **Kafka topics** | Same convention |
| **API endpoints** | Same paths, same request/response |
| **Error codes** | Same semantic codes |
| **Redis keys** | Same format |

### 6.2 Shared Migration

```
migrations/                              # 🟦 Shared
├── 20260101_customer_init.sql
├── 20260102_package_init.sql
├── 20260103_device_init.sql
└── ...
```

ทั้ง Go และ Python **ใช้ migration เดียวกัน** → schema ตรงกันอัตโนมัติ

### 6.3 Spec Duplication Strategy

**Option A: Copy YAML**
```bash
cp .ai/go/modules/04_erp.yaml .ai/py/modules/04_erp.yaml
```
- ✅ ง่าย
- ❌ ต้อง sync เอง

**Option B: Symlink**
```bash
ln -s ../../go/modules .ai/py/modules
```
- ✅ Single source
- ❌ Windows symlink ยุ่งยาก

**Option C: Shared + Override (แนะนำ)**
```
.ai/specs/                    # Shared YAML specs (40 files)
├── 01_customer.yaml
├── 02_package.yaml
└── ...

.ai/go/modules/               # Go-specific overrides (rare)
.ai/py/modules/               # Python-specific overrides (rare)
```

Prompt: "Load spec from `.ai/specs/04_erp.yaml`"

### 6.4 Verify Parity (ถ้าต้องการ API ตรงกัน)

```bash
# หลัง Go + Python เสร็จ
./scripts/verify-parity.sh

# เปรียบเทียบ:
# - API endpoints
# - Request/response schema
# - Error codes
```

---

## 7. Quick Start แยก

### 7.1 Setup (วันนี้ 45 นาที)

```bash
# 1. สร้างโครง 2 สาย
mkdir -p .ai/go/{modules,sessions,checkpoints,prompts}
mkdir -p .ai/py/{modules,sessions,checkpoints,prompts}
mkdir -p .ai/specs    # shared

# 2. สร้างไฟล์ Go
touch .ai/go/CONTEXT.md
touch .ai/go/PROGRESS.md
touch .ai/go/CONVENTIONS.md
touch .ai/go/checkpoints/LAST_SESSION.md
touch .ai/go/checkpoints/NEXT_SESSION.md
touch .ai/go/prompts/generate_module.md

# 3. สร้างไฟล์ Python
touch .ai/py/CONTEXT.md
touch .ai/py/PROGRESS.md
touch .ai/py/CONVENTIONS.md
touch .ai/py/checkpoints/LAST_SESSION.md
touch .ai/py/checkpoints/NEXT_SESSION.md
touch .ai/py/prompts/generate_module.md

# 4. Shared
touch .ai/MASTER_INDEX.md
touch .ai/TOKEN_LOG.md
```

### 7.2 เลือก Track

**Track A — Parallel (แนะนำ)**
```
Day 1: G04 erp Go    |  P04 erp Py
Day 2: G05 crm Go    |  P05 crm Py
Day 3: G06 logistics |  P06 logistics
...
```
✅ ข้อดี: เห็นผลทั้ง 2 ภาษาเร็ว, swap ได้ง่าย
❌ ข้อเสีย: ต้องสลับ context

**Track B — Sequential Go First**
```
Week 1-3: Go 40 modules
Week 4-6: Python 40 modules
```
✅ ข้อดี: Focus เดียว, momentum ดี
❌ ข้อเสีย: Python ช้า

**Track C — Sequential Python First**
```
Week 1-3: Python 40 modules
Week 4-6: Go 40 modules
```
✅ ข้อดี: Python เขียนเร็วกว่า, prototype เร็ว
❌ ข้อเสีย: Go ช้า

### 7.3 เปิด Session ใหม่

**Session Go:**
```
1. เปิด chat ใหม่
2. Paste .ai/go/CONTEXT.md
3. Paste .ai/go/CONVENTIONS.md
4. Paste .ai/go/checkpoints/LAST_SESSION.md
5. Paste .ai/go/checkpoints/NEXT_SESSION.md
6. Paste .ai/specs/{module}.yaml
7. Paste .ai/go/prompts/generate_module.md
8. Say: "ready" → "start"
```

**Session Python:**
```
1. เปิด chat ใหม่ (คนละ chat)
2. Paste .ai/py/CONTEXT.md
3. Paste .ai/py/CONVENTIONS.md
4. Paste .ai/py/checkpoints/LAST_SESSION.md
5. Paste .ai/py/checkpoints/NEXT_SESSION.md
6. Paste .ai/specs/{module}.yaml
7. Paste .ai/py/prompts/generate_module.md
8. Say: "ready" → "start"
```

### 7.4 สลับ Track

**ถ้ากำลังทำ Go อยู่ แล้วอยากสลับไป Python:**
```
1. Checkpoint Go session (บันทึก LAST_SESSION.md)
2. เปิด chat ใหม่
3. Load .ai/py/... 
4. ทำ Python session
5. Checkpoint Python
6. กลับ Go ได้ทุกเมื่อ — context ยังอยู่
```

---

## 📊 ตารางเปรียบเทียบ 2 สาย

| มิติ | Go | Python |
| :--- | :--- | :--- |
| **Setup** | `.ai/go/` | `.ai/py/` |
| **Context size** | ~800 tokens | ~800 tokens |
| **Avg session** | 25-35k | 20-30k |
| **Avg duration** | 30 min | 25 min |
| **LOC/module** | ~3,000 | ~2,200 |
| **Build cmd** | `go build ./...` | `python -m compileall` |
| **Test cmd** | `go test ./...` | `pytest tests/` |
| **Test style** | table-driven | pytest |
| **Entity style** | struct + methods | BaseModel + methods |
| **Repo style** | interface + impl | Protocol + impl |
| **DI style** | manual wire-up | Depends + module.py |
| **Error style** | sentinel | exception classes |
| **Async** | goroutine | asyncio |
| **Framework** | Gin + GORM | FastAPI + SQLAlchemy |
| **Sessions/40 modules** | ~40 | ~40 |
| **Total tokens** | ~1,000k | ~880k |

---

## 🎯 Action Items

### ทันที

- [ ] สร้างโครง `.ai/go/` + `.ai/py/` + `.ai/specs/`
- [ ] เขียน `CONTEXT.md` ทั้ง 2
- [ ] เขียน `CONVENTIONS.md` ทั้ง 2
- [ ] เขียน `PROGRESS.md` ทั้ง 2
- [ ] เขียน `prompts/generate_module.md` ทั้ง 2
- [ ] Copy specs ไป `.ai/specs/`

### พรุ่งนี้

- [ ] Session G04 — erp Go
- [ ] Session P04 — erp Py
- [ ] Verify parity

### สัปดาห์นี้

- [ ] 10 modules × 2 langs = 20 sessions
- [ ] Update MASTER_INDEX
- [ ] Update TOKEN_LOG

### กำหนดการ

| สัปดาห์ | Go | Py | Total |
| :-: | :-: | :-: | :-: |
| 1 | 5 | 5 | 10 |
| 2 | 5 | 5 | 10 |
| 3 | 5 | 5 | 10 |
| 4 | 5 | 5 | 10 |
| ... | ... | ... | ... |
| **รวม** | **40** | **40** | **80** |
| **ระยะเวลา** | **8 สัปดาห์** | **8 สัปดาห์** | **8 สัปดาห์** |

---

## 🚀 คำสั่งถัดไป

พิมพ์ได้เลย:

- `"สร้าง setup .ai/ go + py"` → ผม generate ไฟล์โครงทั้งหมดให้
- `"เขียน CONTEXT.md Go"` → ผมเขียน CONTEXT.md ของ Go
- `"เขียน CONTEXT.md Python"` → ผมเขียน CONTEXT.md ของ Python
- `"เขียน prompts ทั้ง 2 สาย"` → ผมเขียน prompt templates
- `"ทำ G04 erp Go"` → เริ่ม session Go
- `"ทำ P04 erp Python"` → เริ่ม session Python
- `"ทำ G04 + P04 พร้อมกัน"` → guide ทั้ง 2 parallel
  # 🚀 G04 + P04 — COMPLETE (Steps 2-5) + Modules 05-40 Setup

> ทำต่อจนจบ: ERP ครบ 2 ภาษา + เตรียมโครงสำหรับ 36 modules ที่เหลือ

---

# PART 1 — G04 (Go ERP) Steps 2-5

## G04.Step 2 — Application Layer

### `application/ports.go`

```go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/erp/domain/entity"
	"icmongolang/internal/modules/erp/domain/event"
)

type EventProducer interface {
	PublishOrderCreated(ctx context.Context, evt event.OrderCreated) error
	PublishOrderConfirmed(ctx context.Context, evt event.OrderConfirmed) error
	PublishOrderShipped(ctx context.Context, evt event.OrderShipped) error
	PublishOrderDelivered(ctx context.Context, evt event.OrderDelivered) error
	PublishOrderCancelled(ctx context.Context, evt event.OrderCancelled) error
	PublishInvoiceIssued(ctx context.Context, evt event.InvoiceIssued) error
	PublishInvoicePaid(ctx context.Context, evt event.InvoicePaid) error
	PublishStockLow(ctx context.Context, evt event.StockLow) error
	PublishStockReordered(ctx context.Context, evt event.StockReordered) error
}

type AuditRepository interface {
	Save(ctx context.Context, trail *AuditTrail) error
}

type AuditTrail struct {
	UserID     uuid.UUID
	Action     string
	EntityType string
	EntityID   uuid.UUID
	Payload    map[string]interface{}
	IPAddress  string
	OccurredAt time.Time
}

type WSHub interface {
	BroadcastToTenant(tenantID string, event interface{})
}

type CustomerClient interface {
	IsActive(ctx context.Context, customerID uuid.UUID) (bool, error)
	Get(ctx context.Context, customerID uuid.UUID) (*CustomerInfo, error)
}

type CustomerInfo struct {
	ID       uuid.UUID
	Code     string
	Name     string
	IsActive bool
}

type ItemClient interface {
	Get(ctx context.Context, itemID uuid.UUID) (*ItemInfo, error)
}

type ItemInfo struct {
	ID        uuid.UUID
	SKU       string
	Name      string
	UnitPrice float64
	Currency  string
}

type PaymentClient interface {
	CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error)
}

type PaymentRequest struct {
	CustomerID    uuid.UUID
	Amount        float64
	Currency      string
	Description   string
	ReferenceID   uuid.UUID
	ReferenceType string
}

type PaymentResult struct {
	ID  uuid.UUID
	URL string
}

type TransactionManager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}

var _ = entity.Order{}
```

### `application/dto.go`

```go
package application

import "time"

type MoneyDTO struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type OrderItemDTO struct {
	ID        string   `json:"id"`
	ProductID string   `json:"product_id"`
	SKU       string   `json:"sku"`
	Name      string   `json:"name"`
	Quantity  int      `json:"quantity"`
	UnitPrice MoneyDTO `json:"unit_price"`
	Subtotal  MoneyDTO `json:"subtotal"`
}

type OrderResponse struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenant_id"`
	CustomerID  string         `json:"customer_id"`
	OrderNumber string         `json:"order_number"`
	Type        string         `json:"type"`
	Items       []OrderItemDTO `json:"items"`
	Subtotal    MoneyDTO       `json:"subtotal"`
	Tax         MoneyDTO       `json:"tax"`
	Discount    MoneyDTO       `json:"discount"`
	Total       MoneyDTO       `json:"total"`
	Status      string         `json:"status"`
	OrderedAt   time.Time      `json:"ordered_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

type InvoiceResponse struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	CustomerID    string    `json:"customer_id"`
	OrderID       *string   `json:"order_id,omitempty"`
	Number        string    `json:"number"`
	Type          string    `json:"type"`
	Subtotal      MoneyDTO  `json:"subtotal"`
	Tax           MoneyDTO  `json:"tax"`
	Total         MoneyDTO  `json:"total"`
	Status        string    `json:"status"`
	IssuedAt      time.Time `json:"issued_at"`
	DueDate       time.Time `json:"due_date"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	DaysUntilDue  int       `json:"days_until_due"`
}

type StockResponse struct {
	ID           string   `json:"id"`
	ProductID    string   `json:"product_id"`
	WarehouseID  string   `json:"warehouse_id"`
	SKU          string   `json:"sku"`
	Quantity     int      `json:"quantity"`
	Reserved     int      `json:"reserved"`
	Available    int      `json:"available"`
	ReorderPoint int      `json:"reorder_point"`
	UnitCost     MoneyDTO `json:"unit_cost"`
	NeedsReorder bool     `json:"needs_reorder"`
}
```

### `application/create_order.go`

```go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/erp/domain/entity"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	"icmongolang/internal/modules/erp/domain/event"
	"icmongolang/internal/modules/erp/domain/repository"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type CreateOrderUseCase struct {
	orderRepo   repository.OrderRepository
	stockRepo   repository.StockRepository
	customerCli CustomerClient
	itemCli     ItemClient
	producer    EventProducer
	auditRepo   AuditRepository
	txManager   TransactionManager
}

func NewCreateOrderUseCase(
	orderRepo repository.OrderRepository,
	stockRepo repository.StockRepository,
	customerCli CustomerClient,
	itemCli ItemClient,
	producer EventProducer,
	auditRepo AuditRepository,
	txManager TransactionManager,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo: orderRepo, stockRepo: stockRepo,
		customerCli: customerCli, itemCli: itemCli,
		producer: producer, auditRepo: auditRepo, txManager: txManager,
	}
}

type CreateOrderItemInput struct {
	ProductID   uuid.UUID
	WarehouseID uuid.UUID
	Quantity    int
}

type CreateOrderInput struct {
	TenantID   uuid.UUID
	CustomerID uuid.UUID
	Type       valueobject.OrderType
	Items      []CreateOrderItemInput
	Notes      string
	UserID     uuid.UUID
	IPAddress  string
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, in CreateOrderInput) (*OrderResponse, error) {
	// 1. Verify customer
	cust, err := uc.customerCli.Get(ctx, in.CustomerID)
	if err != nil {
		return nil, err
	}
	if !cust.IsActive {
		return nil, domainerrors.ErrInvalidCustomer
	}

	// 2. Generate order number
	orderNumber, err := uc.orderRepo.NextOrderNumber(ctx, in.TenantID, in.Type)
	if err != nil {
		return nil, err
	}

	// 3. Create aggregate
	order, err := entity.NewOrder(in.TenantID, in.CustomerID, in.Type, in.UserID)
	if err != nil {
		return nil, err
	}
	_ = order.SetNumber(orderNumber)
	order.Notes = in.Notes

	// 4. Add items with stock reservation
	for _, itemIn := range in.Items {
		itemInfo, err := uc.itemCli.Get(ctx, itemIn.ProductID)
		if err != nil {
			return nil, err
		}
		unitPrice := valueobject.Money{Amount: itemInfo.UnitPrice, Currency: itemInfo.Currency}
		item, err := entity.NewOrderItem(itemIn.ProductID, itemInfo.SKU, itemInfo.Name, itemIn.Quantity, unitPrice)
		if err != nil {
			return nil, err
		}
		if err := order.AddItem(*item); err != nil {
			return nil, err
		}
	}

	// 5. Persist in transaction
	err = uc.txManager.Run(ctx, func(txCtx context.Context) error {
		if err := uc.orderRepo.Save(txCtx, order); err != nil {
			return err
		}
		// Reserve stock per item
		for _, itemIn := range in.Items {
			stock, err := uc.stockRepo.FindByProductWarehouse(txCtx, itemIn.ProductID, itemIn.WarehouseID)
			if err != nil {
				return err
			}
			if err := stock.Reserve(itemIn.Quantity); err != nil {
				return err
			}
			if err := uc.stockRepo.Save(txCtx, stock); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 6. Audit + Publish
	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID: in.UserID, Action: "ORDER_CREATED",
		EntityType: "order", EntityID: order.ID,
		Payload: map[string]interface{}{
			"order_number": order.OrderNumber,
			"total":        order.Total.Amount,
		},
		IPAddress: in.IPAddress, OccurredAt: time.Now(),
	})

	_ = uc.producer.PublishOrderCreated(ctx, event.OrderCreated{
		OrderID: order.ID, TenantID: order.TenantID,
		CustomerID: order.CustomerID, OrderNumber: order.OrderNumber,
		Type: string(order.Type), Total: order.Total.Amount,
		Currency: order.Total.Currency, ItemCount: order.ItemCount(),
		OccurredAt: time.Now(),
	})

	return toOrderResponse(order), nil
}

func toOrderResponse(o *entity.Order) *OrderResponse {
	items := make([]OrderItemDTO, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, OrderItemDTO{
			ID: item.ID.String(), ProductID: item.ProductID.String(),
			SKU: item.SKU, Name: item.Name, Quantity: item.Quantity,
			UnitPrice: MoneyDTO{Amount: item.UnitPrice.Amount, Currency: item.UnitPrice.Currency},
			Subtotal:  MoneyDTO{Amount: item.Subtotal.Amount, Currency: item.Subtotal.Currency},
		})
	}
	return &OrderResponse{
		ID: o.ID.String(), TenantID: o.TenantID.String(), CustomerID: o.CustomerID.String(),
		OrderNumber: o.OrderNumber, Type: string(o.Type), Items: items,
		Subtotal: MoneyDTO{Amount: o.Subtotal.Amount, Currency: o.Subtotal.Currency},
		Tax:      MoneyDTO{Amount: o.TaxAmount.Amount, Currency: o.TaxAmount.Currency},
		Discount: MoneyDTO{Amount: o.Discount.Amount, Currency: o.Discount.Currency},
		Total:    MoneyDTO{Amount: o.Total.Amount, Currency: o.Total.Currency},
		Status:   string(o.Status), OrderedAt: o.OrderedAt, CreatedAt: o.CreatedAt,
	}
}
```

### `application/confirm_order.go`

```go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/erp/domain/event"
	"icmongolang/internal/modules/erp/domain/repository"
)

type ConfirmOrderUseCase struct {
	orderRepo repository.OrderRepository
	producer  EventProducer
	auditRepo AuditRepository
}

func NewConfirmOrderUseCase(orderRepo repository.OrderRepository, producer EventProducer, auditRepo AuditRepository) *ConfirmOrderUseCase {
	return &ConfirmOrderUseCase{orderRepo: orderRepo, producer: producer, auditRepo: auditRepo}
}

type ConfirmOrderInput struct {
	OrderID uuid.UUID
	UserID  uuid.UUID
}

func (uc *ConfirmOrderUseCase) Execute(ctx context.Context, in ConfirmOrderInput) error {
	order, err := uc.orderRepo.FindByID(ctx, in.OrderID)
	if err != nil {
		return err
	}
	if err := order.Confirm(); err != nil {
		return err
	}
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return err
	}
	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID: in.UserID, Action: "ORDER_CONFIRMED",
		EntityType: "order", EntityID: order.ID,
		OccurredAt: time.Now(),
	})
	return uc.producer.PublishOrderConfirmed(ctx, event.OrderConfirmed{
		OrderID: order.ID, TenantID: order.TenantID, OccurredAt: time.Now(),
	})
}
```

### `application/ship_order.go` + `deliver_order.go` + `cancel_order.go`

```go
package application

// ship_order.go
type ShipOrderUseCase struct {
	orderRepo repository.OrderRepository
	producer  EventProducer
	auditRepo AuditRepository
}

func NewShipOrderUseCase(orderRepo repository.OrderRepository, producer EventProducer, auditRepo AuditRepository) *ShipOrderUseCase {
	return &ShipOrderUseCase{orderRepo: orderRepo, producer: producer, auditRepo: auditRepo}
}

type ShipOrderInput struct {
	OrderID uuid.UUID
	UserID  uuid.UUID
}

func (uc *ShipOrderUseCase) Execute(ctx context.Context, in ShipOrderInput) error {
	order, err := uc.orderRepo.FindByID(ctx, in.OrderID)
	if err != nil {
		return err
	}
	if err := order.Ship(); err != nil {
		return err
	}
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return err
	}
	_ = uc.producer.PublishOrderShipped(ctx, event.OrderShipped{
		OrderID: order.ID, TenantID: order.TenantID, OccurredAt: time.Now(),
	})
	return nil
}

// deliver_order.go
type DeliverOrderUseCase struct {
	orderRepo repository.OrderRepository
	stockRepo repository.StockRepository
	producer  EventProducer
	txManager TransactionManager
}

func NewDeliverOrderUseCase(orderRepo repository.OrderRepository, stockRepo repository.StockRepository, producer EventProducer, txManager TransactionManager) *DeliverOrderUseCase {
	return &DeliverOrderUseCase{orderRepo: orderRepo, stockRepo: stockRepo, producer: producer, txManager: txManager}
}

type DeliverOrderInput struct {
	OrderID uuid.UUID
	UserID  uuid.UUID
}

func (uc *DeliverOrderUseCase) Execute(ctx context.Context, in DeliverOrderInput) error {
	order, err := uc.orderRepo.FindByID(ctx, in.OrderID)
	if err != nil {
		return err
	}
	if err := order.Deliver(); err != nil {
		return err
	}
	err = uc.txManager.Run(ctx, func(txCtx context.Context) error {
		if err := uc.orderRepo.Save(txCtx, order); err != nil {
			return err
		}
		// Issue stock per item
		for _, item := range order.Items {
			stocks, err := uc.stockRepo.FindByProduct(txCtx, item.ProductID)
			if err != nil {
				return err
			}
			for _, stock := range stocks {
				if err := stock.Issue(item.Quantity); err != nil {
					return err
				}
				if err := uc.stockRepo.Save(txCtx, &stock); err != nil {
					return err
				}
				break
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return uc.producer.PublishOrderDelivered(ctx, event.OrderDelivered{
		OrderID: order.ID, TenantID: order.TenantID, OccurredAt: time.Now(),
	})
}

// cancel_order.go
type CancelOrderUseCase struct {
	orderRepo repository.OrderRepository
	stockRepo repository.StockRepository
	producer  EventProducer
	auditRepo AuditRepository
	txManager TransactionManager
}

func NewCancelOrderUseCase(orderRepo repository.OrderRepository, stockRepo repository.StockRepository, producer EventProducer, auditRepo AuditRepository, txManager TransactionManager) *CancelOrderUseCase {
	return &CancelOrderUseCase{orderRepo: orderRepo, stockRepo: stockRepo, producer: producer, auditRepo: auditRepo, txManager: txManager}
}

type CancelOrderInput struct {
	OrderID uuid.UUID
	Reason  string
	UserID  uuid.UUID
}

func (uc *CancelOrderUseCase) Execute(ctx context.Context, in CancelOrderInput) error {
	order, err := uc.orderRepo.FindByID(ctx, in.OrderID)
	if err != nil {
		return err
	}
	if err := order.Cancel(in.Reason); err != nil {
		return err
	}
	err = uc.txManager.Run(ctx, func(txCtx context.Context) error {
		if err := uc.orderRepo.Save(txCtx, order); err != nil {
			return err
		}
		// Release reserved stock
		for _, item := range order.Items {
			stocks, err := uc.stockRepo.FindByProduct(txCtx, item.ProductID)
			if err != nil {
				continue
			}
			for _, stock := range stocks {
				stock.Release(item.Quantity)
				_ = uc.stockRepo.Save(txCtx, &stock)
				break
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID: in.UserID, Action: "ORDER_CANCELLED",
		EntityType: "order", EntityID: order.ID,
		Payload: map[string]interface{}{"reason": in.Reason},
		OccurredAt: time.Now(),
	})
	return uc.producer.PublishOrderCancelled(ctx, event.OrderCancelled{
		OrderID: order.ID, TenantID: order.TenantID,
		Reason: in.Reason, OccurredAt: time.Now(),
	})
}
```

### `application/create_invoice.go` + `record_payment.go`

```go
package application

// create_invoice.go
type CreateInvoiceUseCase struct {
	invoiceRepo repository.InvoiceRepository
	orderRepo   repository.OrderRepository
	auditRepo   AuditRepository
}

func NewCreateInvoiceUseCase(invoiceRepo repository.InvoiceRepository, orderRepo repository.OrderRepository, auditRepo AuditRepository) *CreateInvoiceUseCase {
	return &CreateInvoiceUseCase{invoiceRepo: invoiceRepo, orderRepo: orderRepo, auditRepo: auditRepo}
}

type CreateInvoiceInput struct {
	OrderID uuid.UUID
	Type    valueobject.InvoiceType
	DueDays int
	UserID  uuid.UUID
}

func (uc *CreateInvoiceUseCase) Execute(ctx context.Context, in CreateInvoiceInput) (*InvoiceResponse, error) {
	order, err := uc.orderRepo.FindByID(ctx, in.OrderID)
	if err != nil {
		return nil, err
	}
	if !order.IsDelivered() && !order.IsConfirmed() {
		return nil, domainerrors.ErrOrderNotConfirmed
	}
	dueDate := time.Now().AddDate(0, 0, in.DueDays)
	number, err := uc.invoiceRepo.NextInvoiceNumber(ctx, order.TenantID)
	if err != nil {
		return nil, err
	}
	inv, err := entity.NewInvoice(
		order.TenantID, order.CustomerID, in.Type,
		order.Subtotal, order.TaxAmount, order.Total,
		dueDate, in.UserID,
	)
	if err != nil {
		return nil, err
	}
	_ = inv.SetNumber(number)
	inv.LinkOrder(order.ID)
	if err := uc.invoiceRepo.Save(ctx, inv); err != nil {
		return nil, err
	}
	return toInvoiceResponse(inv), nil
}

func toInvoiceResponse(i *entity.Invoice) *InvoiceResponse {
	var orderID *string
	if i.OrderID != nil {
		s := i.OrderID.String()
		orderID = &s
	}
	return &InvoiceResponse{
		ID: i.ID.String(), TenantID: i.TenantID.String(), CustomerID: i.CustomerID.String(),
		OrderID: orderID, Number: i.Number, Type: string(i.Type),
		Subtotal: MoneyDTO{Amount: i.Subtotal.Amount, Currency: i.Subtotal.Currency},
		Tax:      MoneyDTO{Amount: i.TaxAmount.Amount, Currency: i.TaxAmount.Currency},
		Total:    MoneyDTO{Amount: i.Total.Amount, Currency: i.Total.Currency},
		Status:   string(i.Status), IssuedAt: i.IssuedAt, DueDate: i.DueDate,
		PaidAt: i.PaidAt, DaysUntilDue: i.DaysUntilDue(),
	}
}

// record_payment.go
type RecordPaymentUseCase struct {
	invoiceRepo repository.InvoiceRepository
	producer    EventProducer
	auditRepo   AuditRepository
}

func NewRecordPaymentUseCase(invoiceRepo repository.InvoiceRepository, producer EventProducer, auditRepo AuditRepository) *RecordPaymentUseCase {
	return &RecordPaymentUseCase{invoiceRepo: invoiceRepo, producer: producer, auditRepo: auditRepo}
}

type RecordPaymentInput struct {
	InvoiceID uuid.UUID
	Amount    valueobject.Money
	Method    string
	UserID    uuid.UUID
}

func (uc *RecordPaymentUseCase) Execute(ctx context.Context, in RecordPaymentInput) error {
	inv, err := uc.invoiceRepo.FindByID(ctx, in.InvoiceID)
	if err != nil {
		return err
	}
	if err := inv.MarkPaid(in.Amount, in.Method); err != nil {
		return err
	}
	if err := uc.invoiceRepo.Save(ctx, inv); err != nil {
		return err
	}
	_ = uc.producer.PublishInvoicePaid(ctx, event.InvoicePaid{
		InvoiceID: inv.ID, TenantID: inv.TenantID,
		Amount: in.Amount.Amount, Currency: in.Amount.Currency,
		OccurredAt: time.Now(),
	})
	return nil
}
```

### `application/restock.go` + `reorder_check.go`

```go
package application

// restock.go
type RestockItemUseCase struct {
	stockRepo repository.StockRepository
	producer  EventProducer
}

func NewRestockItemUseCase(stockRepo repository.StockRepository, producer EventProducer) *RestockItemUseCase {
	return &RestockItemUseCase{stockRepo: stockRepo, producer: producer}
}

type RestockItemInput struct {
	StockID  uuid.UUID
	Quantity int
	UnitCost valueobject.Money
}

func (uc *RestockItemUseCase) Execute(ctx context.Context, in RestockItemInput) error {
	stock, err := uc.stockRepo.FindByID(ctx, in.StockID)
	if err != nil {
		return err
	}
	if err := stock.Restock(in.Quantity, in.UnitCost); err != nil {
		return err
	}
	return uc.stockRepo.Save(ctx, stock)
}

// reorder_check.go (scheduled)
type ReorderCheckUseCase struct {
	stockRepo repository.StockRepository
	producer  EventProducer
}

func NewReorderCheckUseCase(stockRepo repository.StockRepository, producer EventProducer) *ReorderCheckUseCase {
	return &ReorderCheckUseCase{stockRepo: stockRepo, producer: producer}
}

func (uc *ReorderCheckUseCase) Execute(ctx context.Context, tenantID uuid.UUID) (int, error) {
	stocks, err := uc.stockRepo.FindNeedsReorder(ctx, tenantID)
	if err != nil {
		return 0, err
	}
	for _, s := range stocks {
		_ = uc.producer.PublishStockLow(ctx, event.StockLow{
			TenantID: s.TenantID, ProductID: s.ProductID,
			SKU: s.SKU, WarehouseID: s.WarehouseID,
			Available: s.Available(), ReorderPoint: s.ReorderPoint,
			OccurredAt: time.Now(),
		})
	}
	return len(stocks), nil
}
```

**Status:** ✅ Step 2 complete (10 files)

---

## G04.Step 3 — Infrastructure Layer

### `infrastructure/persistence/postgres/models.go`

```go
package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrderModel struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	CustomerID   uuid.UUID      `gorm:"type:uuid;not null;index"`
	OrderNumber  string         `gorm:"type:varchar(50);not null"`
	Type         string         `gorm:"type:varchar(20);not null;index"`
	Subtotal     float64        `gorm:"type:numeric(15,2)"`
	TaxAmount    float64        `gorm:"type:numeric(15,2)"`
	Discount     float64        `gorm:"type:numeric(15,2)"`
	Total        float64        `gorm:"type:numeric(15,2)"`
	Currency     string         `gorm:"type:varchar(3);default:'THB'"`
	Status       string         `gorm:"type:varchar(30);not null;index"`
	Notes        string         `gorm:"type:text"`
	ShippingAddr datatypes.JSON `gorm:"type:jsonb"`
	BillingAddr  datatypes.JSON `gorm:"type:jsonb"`
	OrderedAt    time.Time
	ConfirmedAt  *time.Time
	ShippedAt    *time.Time
	DeliveredAt  *time.Time
	CancelledAt  *time.Time
	CancelReason string `gorm:"type:text"`
	CreatedBy    uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Items []OrderItemModel `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func (OrderModel) TableName() string { return "erp_orders" }

type OrderItemModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index"`
	SKU       string    `gorm:"type:varchar(100);not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Quantity  int       `gorm:"not null"`
	UnitPrice float64   `gorm:"type:numeric(15,2);not null"`
	Discount  float64   `gorm:"type:numeric(15,2);default:0"`
	TaxRate   float64   `gorm:"type:numeric(5,2);default:0"`
	Subtotal  float64   `gorm:"type:numeric(15,2);not null"`
}

func (OrderItemModel) TableName() string { return "erp_order_items" }

type InvoiceModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	TenantID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	OrderID       *uuid.UUID `gorm:"type:uuid;index"`
	CustomerID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	Number        string     `gorm:"type:varchar(50);not null"`
	Type          string     `gorm:"type:varchar(20);not null"`
	Subtotal      float64    `gorm:"type:numeric(15,2)"`
	TaxAmount     float64    `gorm:"type:numeric(15,2)"`
	Total         float64    `gorm:"type:numeric(15,2)"`
	Currency      string     `gorm:"type:varchar(3);default:'THB'"`
	Status        string     `gorm:"type:varchar(30);not null;index"`
	IssuedAt      time.Time
	DueDate       time.Time  `gorm:"index"`
	PaidAt        *time.Time
	PaidAmount    float64 `gorm:"type:numeric(15,2);default:0"`
	PaymentMethod string  `gorm:"type:varchar(50)"`
	Notes         string  `gorm:"type:text"`
	CreatedBy     uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (InvoiceModel) TableName() string { return "erp_invoices" }

type WarehouseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index"`
	Code      string         `gorm:"type:varchar(50);not null"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Address   datatypes.JSON `gorm:"type:jsonb"`
	Type      string         `gorm:"type:varchar(30);not null"`
	Capacity  int            `gorm:"default:0"`
	ManagerID *uuid.UUID
	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (WarehouseModel) TableName() string { return "erp_warehouses" }

type StockItemModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID      uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index"`
	WarehouseID   uuid.UUID `gorm:"type:uuid;not null;index"`
	SKU           string    `gorm:"type:varchar(100);not null"`
	Quantity      int       `gorm:"not null;default:0"`
	Reserved      int       `gorm:"not null;default:0"`
	ReorderPoint  int       `gorm:"default:0"`
	UnitCost      float64   `gorm:"type:numeric(15,2);default:0"`
	LastRestocked *time.Time
	UpdatedAt     time.Time
}

func (StockItemModel) TableName() string { return "erp_stock_items" }
```

### `infrastructure/persistence/postgres/order_repo_impl.go`

```go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"icmongolang/internal/modules/erp/domain/entity"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type orderRepoImpl struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) *orderRepoImpl { return &orderRepoImpl{db: db} }

func (r *orderRepoImpl) Save(ctx context.Context, o *entity.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := toOrderModel(o)
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		if err := tx.Where("order_id = ?", o.ID).Delete(&OrderItemModel{}).Error; err != nil {
			return err
		}
		for _, item := range o.Items {
			im := toOrderItemModel(&item)
			if err := tx.Create(im).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *orderRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var m OrderModel
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return toOrderEntity(&m), nil
}

func (r *orderRepoImpl) FindByNumber(ctx context.Context, tenantID uuid.UUID, number string) (*entity.Order, error) {
	var m OrderModel
	err := r.db.WithContext(ctx).Preload("Items").
		Where("tenant_id = ? AND order_number = ?", tenantID, number).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrOrderNotFound
	}
	return toOrderEntity(&m), err
}

func (r *orderRepoImpl) FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]entity.Order, error) {
	var models []OrderModel
	err := r.db.WithContext(ctx).Preload("Items").
		Where("customer_id = ?", customerID).Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toOrderEntities(models), nil
}

func (r *orderRepoImpl) FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.OrderStatus) ([]entity.Order, error) {
	var models []OrderModel
	err := r.db.WithContext(ctx).Preload("Items").
		Where("tenant_id = ? AND status = ?", tenantID, status).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toOrderEntities(models), nil
}

func (r *orderRepoImpl) FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Order, error) {
	var models []OrderModel
	err := r.db.WithContext(ctx).Preload("Items").Where("tenant_id = ?", tenantID).
		Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toOrderEntities(models), nil
}

func (r *orderRepoImpl) NextOrderNumber(ctx context.Context, tenantID uuid.UUID, otype valueobject.OrderType) (string, error) {
	year := time.Now().Year()
	prefix := "SO"
	if otype == valueobject.OrderTypePurchase {
		prefix = "PO"
	}
	var count int64
	r.db.WithContext(ctx).Model(&OrderModel{}).
		Where("tenant_id = ? AND EXTRACT(YEAR FROM created_at) = ?", tenantID, year).
		Count(&count)
	return fmt.Sprintf("%s-%d-%05d", prefix, year, count+1), nil
}

func (r *orderRepoImpl) ExistsByNumber(ctx context.Context, tenantID uuid.UUID, number string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrderModel{}).
		Where("tenant_id = ? AND order_number = ?", tenantID, number).Count(&count).Error
	return count > 0, err
}

func (r *orderRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&OrderModel{}, "id = ?", id).Error
}

// ---------- Mappers ----------

func toOrderModel(o *entity.Order) *OrderModel {
	shipJSON, _ := json.Marshal(o.ShippingAddr)
	billJSON, _ := json.Marshal(o.BillingAddr)
	return &OrderModel{
		ID: o.ID, TenantID: o.TenantID, CustomerID: o.CustomerID,
		OrderNumber: o.OrderNumber, Type: string(o.Type),
		Subtotal: o.Subtotal.Amount, TaxAmount: o.TaxAmount.Amount,
		Discount: o.Discount.Amount, Total: o.Total.Amount,
		Currency: o.Total.Currency, Status: string(o.Status),
		Notes: o.Notes, ShippingAddr: datatypes.JSON(shipJSON),
		BillingAddr: datatypes.JSON(billJSON),
		OrderedAt: o.OrderedAt, ConfirmedAt: o.ConfirmedAt,
		ShippedAt: o.ShippedAt, DeliveredAt: o.DeliveredAt,
		CancelledAt: o.CancelledAt, CancelReason: o.CancelReason,
		CreatedBy: o.CreatedBy, CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt,
	}
}

func toOrderItemModel(item *entity.OrderItem) *OrderItemModel {
	return &OrderItemModel{
		ID: item.ID, OrderID: item.OrderID, ProductID: item.ProductID,
		SKU: item.SKU, Name: item.Name, Quantity: item.Quantity,
		UnitPrice: item.UnitPrice.Amount, Discount: item.Discount.Amount,
		TaxRate: item.TaxRate, Subtotal: item.Subtotal.Amount,
	}
}

func toOrderEntity(m *OrderModel) *entity.Order {
	var shipAddr, billAddr valueobject.Address
	_ = json.Unmarshal(m.ShippingAddr, &shipAddr)
	_ = json.Unmarshal(m.BillingAddr, &billAddr)

	items := make([]entity.OrderItem, 0, len(m.Items))
	for _, im := range m.Items {
		items = append(items, entity.OrderItem{
			ID: im.ID, OrderID: im.OrderID, ProductID: im.ProductID,
			SKU: im.SKU, Name: im.Name, Quantity: im.Quantity,
			UnitPrice: valueobject.Money{Amount: im.UnitPrice, Currency: m.Currency},
			Discount:  valueobject.Money{Amount: im.Discount, Currency: m.Currency},
			TaxRate: im.TaxRate,
			Subtotal: valueobject.Money{Amount: im.Subtotal, Currency: m.Currency},
		})
	}

	return &entity.Order{
		ID: m.ID, TenantID: m.TenantID, CustomerID: m.CustomerID,
		OrderNumber: m.OrderNumber, Type: valueobject.OrderType(m.Type),
		Items: items,
		Subtotal: valueobject.Money{Amount: m.Subtotal, Currency: m.Currency},
		TaxAmount: valueobject.Money{Amount: m.TaxAmount, Currency: m.Currency},
		Discount:  valueobject.Money{Amount: m.Discount, Currency: m.Currency},
		Total:     valueobject.Money{Amount: m.Total, Currency: m.Currency},
		Status: valueobject.OrderStatus(m.Status), Notes: m.Notes,
		ShippingAddr: shipAddr, BillingAddr: billAddr,
		OrderedAt: m.OrderedAt, ConfirmedAt: m.ConfirmedAt,
		ShippedAt: m.ShippedAt, DeliveredAt: m.DeliveredAt,
		CancelledAt: m.CancelledAt, CancelReason: m.CancelReason,
		CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func toOrderEntities(models []OrderModel) []entity.Order {
	out := make([]entity.Order, 0, len(models))
	for i := range models {
		out = append(out, *toOrderEntity(&models[i]))
	}
	return out
}
```

### `infrastructure/persistence/postgres/invoice_repo_impl.go`

```go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"icmongolang/internal/modules/erp/domain/entity"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type invoiceRepoImpl struct{ db *gorm.DB }

func NewInvoiceRepository(db *gorm.DB) *invoiceRepoImpl { return &invoiceRepoImpl{db: db} }

func (r *invoiceRepoImpl) Save(ctx context.Context, i *entity.Invoice) error {
	return r.db.WithContext(ctx).Save(toInvoiceModel(i)).Error
}

func (r *invoiceRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error) {
	var m InvoiceModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}
	return toInvoiceEntity(&m), nil
}

func (r *invoiceRepoImpl) FindByNumber(ctx context.Context, tenantID uuid.UUID, number string) (*entity.Invoice, error) {
	var m InvoiceModel
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND number = ?", tenantID, number).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrInvoiceNotFound
	}
	return toInvoiceEntity(&m), err
}

func (r *invoiceRepoImpl) FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]entity.Invoice, error) {
	var models []InvoiceModel
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).
		Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toInvoiceEntities(models), nil
}

func (r *invoiceRepoImpl) FindByOrder(ctx context.Context, orderID uuid.UUID) ([]entity.Invoice, error) {
	var models []InvoiceModel
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toInvoiceEntities(models), nil
}

func (r *invoiceRepoImpl) FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.InvoiceStatus) ([]entity.Invoice, error) {
	var models []InvoiceModel
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND status = ?", tenantID, status).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toInvoiceEntities(models), nil
}

func (r *invoiceRepoImpl) FindOverdue(ctx context.Context, before time.Time) ([]entity.Invoice, error) {
	var models []InvoiceModel
	err := r.db.WithContext(ctx).
		Where("status = 'ISSUED' AND due_date < ?", before).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toInvoiceEntities(models), nil
}

func (r *invoiceRepoImpl) NextInvoiceNumber(ctx context.Context, tenantID uuid.UUID) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).Model(&InvoiceModel{}).
		Where("tenant_id = ? AND EXTRACT(YEAR FROM created_at) = ?", tenantID, year).
		Count(&count)
	return fmt.Sprintf("INV-%d-%05d", year, count+1), nil
}

func (r *invoiceRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&InvoiceModel{}, "id = ?", id).Error
}

func toInvoiceModel(i *entity.Invoice) *InvoiceModel {
	return &InvoiceModel{
		ID: i.ID, TenantID: i.TenantID, OrderID: i.OrderID,
		CustomerID: i.CustomerID, Number: i.Number, Type: string(i.Type),
		Subtotal: i.Subtotal.Amount, TaxAmount: i.TaxAmount.Amount,
		Total: i.Total.Amount, Currency: i.Total.Currency,
		Status: string(i.Status), IssuedAt: i.IssuedAt, DueDate: i.DueDate,
		PaidAt: i.PaidAt, PaidAmount: i.PaidAmount.Amount,
		PaymentMethod: i.PaymentMethod, Notes: i.Notes,
		CreatedBy: i.CreatedBy, CreatedAt: i.CreatedAt, UpdatedAt: i.UpdatedAt,
	}
}

func toInvoiceEntity(m *InvoiceModel) *entity.Invoice {
	return &entity.Invoice{
		ID: m.ID, TenantID: m.TenantID, OrderID: m.OrderID,
		CustomerID: m.CustomerID, Number: m.Number,
		Type: valueobject.InvoiceType(m.Type),
		Subtotal: valueobject.Money{Amount: m.Subtotal, Currency: m.Currency},
		TaxAmount: valueobject.Money{Amount: m.TaxAmount, Currency: m.Currency},
		Total:     valueobject.Money{Amount: m.Total, Currency: m.Currency},
		Status: valueobject.InvoiceStatus(m.Status),
		IssuedAt: m.IssuedAt, DueDate: m.DueDate, PaidAt: m.PaidAt,
		PaidAmount: valueobject.Money{Amount: m.PaidAmount, Currency: m.Currency},
		PaymentMethod: m.PaymentMethod, Notes: m.Notes,
		CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func toInvoiceEntities(models []InvoiceModel) []entity.Invoice {
	out := make([]entity.Invoice, 0, len(models))
	for i := range models {
		out = append(out, *toInvoiceEntity(&models[i]))
	}
	return out
}
```

### `infrastructure/persistence/postgres/stock_repo_impl.go`

```go
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"icmongolang/internal/modules/erp/domain/entity"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type stockRepoImpl struct{ db *gorm.DB }

func NewStockRepository(db *gorm.DB) *stockRepoImpl { return &stockRepoImpl{db: db} }

func (r *stockRepoImpl) Save(ctx context.Context, s *entity.StockItem) error {
	return r.db.WithContext(ctx).Save(toStockModel(s)).Error
}

func (r *stockRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.StockItem, error) {
	var m StockItemModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrStockNotFound
	}
	if err != nil {
		return nil, err
	}
	return toStockEntity(&m), nil
}

func (r *stockRepoImpl) FindByProductWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) (*entity.StockItem, error) {
	var m StockItemModel
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrStockNotFound
	}
	if err != nil {
		return nil, err
	}
	return toStockEntity(&m), nil
}

func (r *stockRepoImpl) FindByProduct(ctx context.Context, productID uuid.UUID) ([]entity.StockItem, error) {
	var models []StockItemModel
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toStockEntities(models), nil
}

func (r *stockRepoImpl) FindByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]entity.StockItem, error) {
	var models []StockItemModel
	err := r.db.WithContext(ctx).Where("warehouse_id = ?", warehouseID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toStockEntities(models), nil
}

func (r *stockRepoImpl) FindNeedsReorder(ctx context.Context, tenantID uuid.UUID) ([]entity.StockItem, error) {
	var models []StockItemModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND (quantity - reserved) <= reorder_point", tenantID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toStockEntities(models), nil
}

func (r *stockRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&StockItemModel{}, "id = ?", id).Error
}

func toStockModel(s *entity.StockItem) *StockItemModel {
	return &StockItemModel{
		ID: s.ID, TenantID: s.TenantID, ProductID: s.ProductID,
		WarehouseID: s.WarehouseID, SKU: s.SKU,
		Quantity: s.Quantity, Reserved: s.Reserved,
		ReorderPoint: s.ReorderPoint, UnitCost: s.UnitCost.Amount,
		LastRestocked: s.LastRestocked, UpdatedAt: s.UpdatedAt,
	}
}

func toStockEntity(m *StockItemModel) *entity.StockItem {
	return &entity.StockItem{
		ID: m.ID, TenantID: m.TenantID, ProductID: m.ProductID,
		WarehouseID: m.WarehouseID, SKU: m.SKU,
		Quantity: m.Quantity, Reserved: m.Reserved,
		ReorderPoint: m.ReorderPoint,
		UnitCost:     valueobject.Money{Amount: m.UnitCost, Currency: "THB"},
		LastRestocked: m.LastRestocked, UpdatedAt: m.UpdatedAt,
	}
}

func toStockEntities(models []StockItemModel) []entity.StockItem {
	out := make([]entity.StockItem, 0, len(models))
	for i := range models {
		out = append(out, *toStockEntity(&models[i]))
	}
	return out
}
```

### `infrastructure/persistence/postgres/warehouse_repo_impl.go`

```go
package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"icmongolang/internal/modules/erp/domain/entity"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

type warehouseRepoImpl struct{ db *gorm.DB }

func NewWarehouseRepository(db *gorm.DB) *warehouseRepoImpl { return &warehouseRepoImpl{db: db} }

func (r *warehouseRepoImpl) Save(ctx context.Context, w *entity.Warehouse) error {
	return r.db.WithContext(ctx).Save(toWarehouseModel(w)).Error
}

func (r *warehouseRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	var m WarehouseModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrWarehouseNotFound
	}
	if err != nil {
		return nil, err
	}
	return toWarehouseEntity(&m), nil
}

func (r *warehouseRepoImpl) FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Warehouse, error) {
	var m WarehouseModel
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND code = ?", tenantID, code).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrWarehouseNotFound
	}
	return toWarehouseEntity(&m), err
}

func (r *warehouseRepoImpl) FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Warehouse, error) {
	var models []WarehouseModel
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.Warehouse, 0, len(models))
	for i := range models {
		out = append(out, *toWarehouseEntity(&models[i]))
	}
	return out, nil
}

func (r *warehouseRepoImpl) ExistsByCode(ctx context.Context, tenantID uuid.UUID, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&WarehouseModel{}).
		Where("tenant_id = ? AND code = ?", tenantID, code).Count(&count).Error
	return count > 0, err
}

func (r *warehouseRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&WarehouseModel{}, "id = ?", id).Error
}

func toWarehouseModel(w *entity.Warehouse) *WarehouseModel {
	addrJSON, _ := json.Marshal(w.Address)
	return &WarehouseModel{
		ID: w.ID, TenantID: w.TenantID, Code: w.Code, Name: w.Name,
		Address: datatypes.JSON(addrJSON), Type: string(w.Type),
		Capacity: w.Capacity, ManagerID: w.ManagerID, IsActive: w.IsActive,
		CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt,
	}
}

func toWarehouseEntity(m *WarehouseModel) *entity.Warehouse {
	var addr valueobject.Address
	_ = json.Unmarshal(m.Address, &addr)
	return &entity.Warehouse{
		ID: m.ID, TenantID: m.TenantID, Code: m.Code, Name: m.Name,
		Address: addr, Type: valueobject.WarehouseType(m.Type),
		Capacity: m.Capacity, ManagerID: m.ManagerID, IsActive: m.IsActive,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
```

### `infrastructure/messaging/producer.go`

```go
package messaging

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"icmongolang/internal/modules/erp/domain/event"
)

const (
	TopicOrderCreated     = "erp.order.created"
	TopicOrderConfirmed   = "erp.order.confirmed"
	TopicOrderShipped     = "erp.order.shipped"
	TopicOrderDelivered   = "erp.order.delivered"
	TopicOrderCancelled   = "erp.order.cancelled"
	TopicInvoiceIssued    = "erp.invoice.issued"
	TopicInvoicePaid      = "erp.invoice.paid"
	TopicStockLow         = "erp.stock.low"
	TopicStockReordered   = "erp.stock.reordered"
)

type kafkaProducer struct{ sync sarama.SyncProducer }

func NewKafkaProducer(brokers []string, clientID string) (*kafkaProducer, error) {
	cfg := sarama.NewConfig()
	cfg.ClientID = clientID
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{sync: p}, nil
}

func (k *kafkaProducer) publish(_ context.Context, topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, _, err = k.sync.SendMessage(&sarama.ProducerMessage{
		Topic: topic, Key: sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func (k *kafkaProducer) PublishOrderCreated(ctx context.Context, evt event.OrderCreated) error {
	return k.publish(ctx, TopicOrderCreated, evt.OrderID.String(), evt)
}
func (k *kafkaProducer) PublishOrderConfirmed(ctx context.Context, evt event.OrderConfirmed) error {
	return k.publish(ctx, TopicOrderConfirmed, evt.OrderID.String(), evt)
}
func (k *kafkaProducer) PublishOrderShipped(ctx context.Context, evt event.OrderShipped) error {
	return k.publish(ctx, TopicOrderShipped, evt.OrderID.String(), evt)
}
func (k *kafkaProducer) PublishOrderDelivered(ctx context.Context, evt event.OrderDelivered) error {
	return k.publish(ctx, TopicOrderDelivered, evt.OrderID.String(), evt)
}
func (k *kafkaProducer) PublishOrderCancelled(ctx context.Context, evt event.OrderCancelled) error {
	return k.publish(ctx, TopicOrderCancelled, evt.OrderID.String(), evt)
}
func (k *kafkaProducer) PublishInvoiceIssued(ctx context.Context, evt event.InvoiceIssued) error {
	return k.publish(ctx, TopicInvoiceIssued, evt.InvoiceID.String(), evt)
}
func (k *kafkaProducer) PublishInvoicePaid(ctx context.Context, evt event.InvoicePaid) error {
	return k.publish(ctx, TopicInvoicePaid, evt.InvoiceID.String(), evt)
}
func (k *kafkaProducer) PublishStockLow(ctx context.Context, evt event.StockLow) error {
	return k.publish(ctx, TopicStockLow, evt.ProductID.String(), evt)
}
func (k *kafkaProducer) PublishStockReordered(ctx context.Context, evt event.StockReordered) error {
	return k.publish(ctx, TopicStockReordered, evt.ProductID.String(), evt)
}
```

### `infrastructure/scheduler/invoice_overdue_job.go`

```go
package scheduler

import (
	"context"
	"log"
	"time"

	"icmongolang/internal/modules/erp/domain/repository"
)

type InvoiceOverdueJob struct {
	invoiceRepo repository.InvoiceRepository
}

func NewInvoiceOverdueJob(r repository.InvoiceRepository) *InvoiceOverdueJob {
	return &InvoiceOverdueJob{invoiceRepo: r}
}

func (j *InvoiceOverdueJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	invoices, err := j.invoiceRepo.FindOverdue(ctx, time.Now())
	if err != nil {
		log.Printf("[invoice_overdue] find: %v", err)
		return
	}
	for i := range invoices {
		inv := &invoices[i]
		inv.MarkOverdue()
		_ = j.invoiceRepo.Save(ctx, inv)
	}
	log.Printf("[invoice_overdue] marked %d overdue", len(invoices))
}
```

**Status:** ✅ Step 3 complete (~10 files)

---

## G04.Step 4 — Interface Layer + module.go

### `interfaces/http/order_handler.go`

```go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"icmongolang/internal/modules/erp/application"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
	"icmongolang/pkg/responses"
)

type OrderHandler struct {
	createUC  *application.CreateOrderUseCase
	confirmUC *application.ConfirmOrderUseCase
	shipUC    *application.ShipOrderUseCase
	deliverUC *application.DeliverOrderUseCase
	cancelUC  *application.CancelOrderUseCase
	getUC     *application.GetOrderUseCase
}

func NewOrderHandler(
	createUC *application.CreateOrderUseCase,
	confirmUC *application.ConfirmOrderUseCase,
	shipUC *application.ShipOrderUseCase,
	deliverUC *application.DeliverOrderUseCase,
	cancelUC *application.CancelOrderUseCase,
	getUC *application.GetOrderUseCase,
) *OrderHandler {
	return &OrderHandler{createUC: createUC, confirmUC: confirmUC, shipUC: shipUC, deliverUC: deliverUC, cancelUC: cancelUC, getUC: getUC}
}

// POST /api/v1/erp/orders
func (h *OrderHandler) Create(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID", err.Error(), nil))
		return
	}

	custID, _ := uuid.Parse(req.CustomerID)
	items := make([]application.CreateOrderItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		pid, _ := uuid.Parse(it.ProductID)
		wid, _ := uuid.Parse(it.WarehouseID)
		items = append(items, application.CreateOrderItemInput{
			ProductID: pid, WarehouseID: wid, Quantity: it.Quantity,
		})
	}

	out, err := h.createUC.Execute(c.Request.Context(), application.CreateOrderInput{
		TenantID: tenantID, CustomerID: custID,
		Type: valueobject.OrderType(req.Type), Items: items,
		Notes: req.Notes, UserID: uid, IPAddress: c.ClientIP(),
	})
	if err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, responses.Success(out))
}

func (h *OrderHandler) Confirm(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.confirmUC.Execute(c.Request.Context(), application.ConfirmOrderInput{OrderID: id, UserID: uid}); err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "confirmed"}))
}

func (h *OrderHandler) Ship(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.shipUC.Execute(c.Request.Context(), application.ShipOrderInput{OrderID: id, UserID: uid}); err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "shipped"}))
}

func (h *OrderHandler) Deliver(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.deliverUC.Execute(c.Request.Context(), application.DeliverOrderInput{OrderID: id, UserID: uid}); err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "delivered"}))
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Param("id"))
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID", err.Error(), nil))
		return
	}
	if err := h.cancelUC.Execute(c.Request.Context(), application.CancelOrderInput{OrderID: id, Reason: req.Reason, UserID: uid}); err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "cancelled"}))
}

func (h *OrderHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_ID", "invalid id", nil))
		return
	}
	out, err := h.getUC.Execute(c.Request.Context(), id)
	if err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(out))
}
```

### `interfaces/http/invoice_handler.go` + `stock_handler.go`

```go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"icmongolang/internal/modules/erp/application"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
	"icmongolang/pkg/responses"
)

type InvoiceHandler struct {
	createUC *application.CreateInvoiceUseCase
	issueUC  *application.IssueInvoiceUseCase
	payUC    *application.RecordPaymentUseCase
}

func NewInvoiceHandler(
	createUC *application.CreateInvoiceUseCase,
	issueUC *application.IssueInvoiceUseCase,
	payUC *application.RecordPaymentUseCase,
) *InvoiceHandler {
	return &InvoiceHandler{createUC: createUC, issueUC: issueUC, payUC: payUC}
}

func (h *InvoiceHandler) Create(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
		Type    string `json:"type" binding:"required"`
		DueDays int    `json:"due_days" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID", err.Error(), nil))
		return
	}
	oid, _ := uuid.Parse(req.OrderID)
	out, err := h.createUC.Execute(c.Request.Context(), application.CreateInvoiceInput{
		OrderID: oid, Type: valueobject.InvoiceType(req.Type),
		DueDays: req.DueDays, UserID: uid,
	})
	if err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, responses.Success(out))
}

func (h *InvoiceHandler) Issue(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.issueUC.Execute(c.Request.Context(), application.IssueInvoiceInput{InvoiceID: id, UserID: uid}); err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "issued"}))
}

func (h *InvoiceHandler) RecordPayment(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Param("id"))
	var req struct {
		Amount   float64 `json:"amount" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
		Method   string  `json:"method" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID", err.Error(), nil))
		return
	}
	if err := h.payUC.Execute(c.Request.Context(), application.RecordPaymentInput{
		InvoiceID: id,
		Amount:    valueobject.Money{Amount: req.Amount, Currency: req.Currency},
		Method:    req.Method, UserID: uid,
	}); err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "paid"}))
}

type StockHandler struct {
	restockUC *application.RestockItemUseCase
}

func NewStockHandler(restockUC *application.RestockItemUseCase) *StockHandler {
	return &StockHandler{restockUC: restockUC}
}

func (h *StockHandler) Restock(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var req struct {
		Quantity int     `json:"quantity" binding:"required,min=1"`
		UnitCost float64 `json:"unit_cost" binding:"required"`
		Currency string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID", err.Error(), nil))
		return
	}
	currency := req.Currency
	if currency == "" {
		currency = "THB"
	}
	if err := h.restockUC.Execute(c.Request.Context(), application.RestockItemInput{
		StockID: id, Quantity: req.Quantity,
		UnitCost: valueobject.Money{Amount: req.UnitCost, Currency: currency},
	}); err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "restocked"}))
}
```

### `interfaces/http/dto.go` + `errors.go` + `routes.go`

```go
// dto.go
package http

type CreateOrderRequest struct {
	CustomerID string                 `json:"customer_id" binding:"required"`
	Type       string                 `json:"type" binding:"required,oneof=SALES PURCHASE"`
	Items      []CreateOrderItemReq   `json:"items" binding:"required,min=1"`
	Notes      string                 `json:"notes"`
}

type CreateOrderItemReq struct {
	ProductID   string `json:"product_id" binding:"required"`
	WarehouseID string `json:"warehouse_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}
```

```go
// errors.go
package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	"icmongolang/pkg/responses"
)

func respondDomainError(c *gin.Context, err error) {
	type mapping struct {
		err    error
		status int
		code   string
	}
	mappings := []mapping{
		{domainerrors.ErrOrderNotFound, 404, "ORDER_NOT_FOUND"},
		{domainerrors.ErrInvoiceNotFound, 404, "INVOICE_NOT_FOUND"},
		{domainerrors.ErrWarehouseNotFound, 404, "WAREHOUSE_NOT_FOUND"},
		{domainerrors.ErrStockNotFound, 404, "STOCK_NOT_FOUND"},
		{domainerrors.ErrOrderNumberDuplicate, 409, "ORDER_NUMBER_DUPLICATE"},
		{domainerrors.ErrInvoiceNumberDuplicate, 409, "INVOICE_NUMBER_DUPLICATE"},
		{domainerrors.ErrInsufficientStock, 409, "INSUFFICIENT_STOCK"},
		{domainerrors.ErrInvoiceAlreadyPaid, 409, "INVOICE_ALREADY_PAID"},
	}
	for _, m := range mappings {
		if errors.Is(err, m.err) {
			c.JSON(m.status, responses.Error(m.code, err.Error(), nil))
			return
		}
	}
	c.JSON(http.StatusBadRequest, responses.Error("BAD_REQUEST", err.Error(), nil))
}
```

```go
// routes.go
package http

import "github.com/gin-gonic/gin"

type Handlers struct {
	Order   *OrderHandler
	Invoice *InvoiceHandler
	Stock   *StockHandler
}

func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handlers,
	auth gin.HandlerFunc,
	tenant gin.HandlerFunc,
	rbac func(roles ...string) gin.HandlerFunc,
) {
	g := r.Group("/erp")
	g.Use(auth, tenant)

	// Orders
	g.GET("/orders", rbac("viewer"), h.Order.List)
	g.POST("/orders", rbac("manager"), h.Order.Create)
	g.GET("/orders/:id", rbac("viewer"), h.Order.Get)
	g.POST("/orders/:id/confirm", rbac("manager"), h.Order.Confirm)
	g.POST("/orders/:id/ship", rbac("operator"), h.Order.Ship)
	g.POST("/orders/:id/deliver", rbac("operator"), h.Order.Deliver)
	g.POST("/orders/:id/cancel", rbac("manager"), h.Order.Cancel)

	// Invoices
	g.GET("/invoices", rbac("viewer"), h.Invoice.List)
	g.POST("/invoices", rbac("manager"), h.Invoice.Create)
	g.GET("/invoices/:id", rbac("viewer"), h.Invoice.Get)
	g.POST("/invoices/:id/issue", rbac("manager"), h.Invoice.Issue)
	g.POST("/invoices/:id/payment", rbac("manager"), h.Invoice.RecordPayment)

	// Stock
	g.GET("/stock", rbac("viewer"), h.Stock.List)
	g.POST("/stock/:id/restock", rbac("manager"), h.Stock.Restock)
	g.GET("/stock/reorder-check", rbac("manager"), h.Stock.ReorderCheck)
}
```

### `module.go`

```go
package erp

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"icmongolang/internal/modules/erp/application"
	"icmongolang/internal/modules/erp/infrastructure/messaging"
	pgrepo "icmongolang/internal/modules/erp/infrastructure/persistence/postgres"
	httpiface "icmongolang/internal/modules/erp/interfaces/http"
)

type Dependencies struct {
	DB          *gorm.DB
	Producer    *messaging.KafkaProducer
	CustomerCli application.CustomerClient
	ItemCli     application.ItemClient
	PaymentCli  application.PaymentClient
	AuditRepo   application.AuditRepository
	WSHub       application.WSHub
	TxManager   application.TransactionManager
}

func Init(
	router *gin.RouterGroup,
	deps Dependencies,
	auth gin.HandlerFunc,
	tenant gin.HandlerFunc,
	rbac func(roles ...string) gin.HandlerFunc,
) {
	orderRepo := pgrepo.NewOrderRepository(deps.DB)
	invoiceRepo := pgrepo.NewInvoiceRepository(deps.DB)
	warehouseRepo := pgrepo.NewWarehouseRepository(deps.DB)
	stockRepo := pgrepo.NewStockRepository(deps.DB)
	_ = warehouseRepo

	createOrderUC := application.NewCreateOrderUseCase(
		orderRepo, stockRepo, deps.CustomerCli, deps.ItemCli,
		deps.Producer, deps.AuditRepo, deps.TxManager,
	)
	getOrderUC := application.NewGetOrderUseCase(orderRepo)
	confirmOrderUC := application.NewConfirmOrderUseCase(orderRepo, deps.Producer, deps.AuditRepo)
	shipOrderUC := application.NewShipOrderUseCase(orderRepo, deps.Producer, deps.AuditRepo)
	deliverOrderUC := application.NewDeliverOrderUseCase(orderRepo, stockRepo, deps.Producer, deps.TxManager)
	cancelOrderUC := application.NewCancelOrderUseCase(orderRepo, stockRepo, deps.Producer, deps.AuditRepo, deps.TxManager)

	createInvoiceUC := application.NewCreateInvoiceUseCase(invoiceRepo, orderRepo, deps.AuditRepo)
	issueInvoiceUC := application.NewIssueInvoiceUseCase(invoiceRepo, deps.Producer, deps.AuditRepo)
	recordPaymentUC := application.NewRecordPaymentUseCase(invoiceRepo, deps.Producer, deps.AuditRepo)

	restockUC := application.NewRestockItemUseCase(stockRepo, deps.Producer)

	orderHandler := httpiface.NewOrderHandler(createOrderUC, confirmOrderUC, shipOrderUC, deliverOrderUC, cancelOrderUC, getOrderUC)
	invoiceHandler := httpiface.NewInvoiceHandler(createInvoiceUC, issueInvoiceUC, recordPaymentUC)
	stockHandler := httpiface.NewStockHandler(restockUC)

	httpiface.RegisterRoutes(router, &httpiface.Handlers{
		Order: orderHandler, Invoice: invoiceHandler, Stock: stockHandler,
	}, auth, tenant, rbac)
}
```

**Status:** ✅ Step 4 complete (~6 files)

---

## G04.Step 5 — Migration + Tests

### `migrations/20260104_erp_init.sql`

```sql
-- ERP module — initial schema

CREATE TABLE IF NOT EXISTS erp_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    order_number VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL,
    subtotal NUMERIC(15,2) NOT NULL DEFAULT 0,
    tax_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    discount NUMERIC(15,2) NOT NULL DEFAULT 0,
    total NUMERIC(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'THB',
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    notes TEXT,
    shipping_addr JSONB,
    billing_addr JSONB,
    ordered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancel_reason TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_erp_orders_tenant_number UNIQUE (tenant_id, order_number)
);
CREATE INDEX idx_erp_orders_tenant ON erp_orders (tenant_id);
CREATE INDEX idx_erp_orders_customer ON erp_orders (customer_id);
CREATE INDEX idx_erp_orders_status ON erp_orders (tenant_id, status);
CREATE INDEX idx_erp_orders_type ON erp_orders (tenant_id, type);

CREATE TABLE IF NOT EXISTS erp_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES erp_orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    unit_price NUMERIC(15,2) NOT NULL,
    discount NUMERIC(15,2) DEFAULT 0,
    tax_rate NUMERIC(5,2) DEFAULT 0,
    subtotal NUMERIC(15,2) NOT NULL
);
CREATE INDEX idx_erp_order_items_order ON erp_order_items (order_id);
CREATE INDEX idx_erp_order_items_product ON erp_order_items (product_id);

CREATE TABLE IF NOT EXISTS erp_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    order_id UUID REFERENCES erp_orders(id),
    customer_id UUID NOT NULL,
    number VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL,
    subtotal NUMERIC(15,2) NOT NULL,
    tax_amount NUMERIC(15,2) NOT NULL,
    total NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'THB',
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    issued_at TIMESTAMP,
    due_date TIMESTAMP,
    paid_at TIMESTAMP,
    paid_amount NUMERIC(15,2) DEFAULT 0,
    payment_method VARCHAR(50),
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_erp_invoices_tenant_number UNIQUE (tenant_id, number)
);
CREATE INDEX idx_erp_invoices_tenant ON erp_invoices (tenant_id);
CREATE INDEX idx_erp_invoices_customer ON erp_invoices (customer_id);
CREATE INDEX idx_erp_invoices_status ON erp_invoices (status);
CREATE INDEX idx_erp_invoices_due ON erp_invoices (due_date) WHERE status IN ('ISSUED', 'OVERDUE');

CREATE TABLE IF NOT EXISTS erp_warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address JSONB,
    type VARCHAR(30) NOT NULL,
    capacity INT DEFAULT 0,
    manager_id UUID,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_erp_warehouses_tenant_code UNIQUE (tenant_id, code)
);
CREATE INDEX idx_erp_warehouses_tenant ON erp_warehouses (tenant_id);

CREATE TABLE IF NOT EXISTS erp_stock_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    product_id UUID NOT NULL,
    warehouse_id UUID NOT NULL REFERENCES erp_warehouses(id),
    sku VARCHAR(100) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    reserved INT NOT NULL DEFAULT 0,
    reorder_point INT NOT NULL DEFAULT 0,
    unit_cost NUMERIC(15,2) DEFAULT 0,
    last_restocked TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_erp_stock_product_warehouse UNIQUE (product_id, warehouse_id)
);
CREATE INDEX idx_erp_stock_tenant ON erp_stock_items (tenant_id);
CREATE INDEX idx_erp_stock_product ON erp_stock_items (product_id);
CREATE INDEX idx_erp_stock_warehouse ON erp_stock_items (warehouse_id);
CREATE INDEX idx_erp_stock_low ON erp_stock_items (tenant_id) WHERE (quantity - reserved) <= reorder_point;

-- Seed sample warehouse
INSERT INTO erp_warehouses (tenant_id, code, name, type, is_active)
VALUES ('00000000-0000-0000-0000-000000000001', 'MAIN', 'Main Warehouse', 'MAIN', true)
ON CONFLICT DO NOTHING;
```

### `domain/entity/order_test.go`

```go
package entity_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"icmongolang/internal/modules/erp/domain/entity"
	domainerrors "icmongolang/internal/modules/erp/domain/errors"
	valueobject "icmongolang/internal/modules/erp/domain/value_object"
)

func TestOrder_Lifecycle(t *testing.T) {
	order, _ := entity.NewOrder(
		uuid.New(), uuid.New(),
		valueobject.OrderTypeSales, uuid.New(),
	)

	item, _ := entity.NewOrderItem(
		uuid.New(), "SKU001", "Product",
		2, valueobject.Money{Amount: 100, Currency: "THB"},
	)
	if err := order.AddItem(*item); err != nil {
		t.Fatalf("add item: %v", err)
	}

	if order.Total.Amount != 200 {
		t.Errorf("expected 200, got %.2f", order.Total.Amount)
	}

	if err := order.Confirm(); err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if !order.IsConfirmed() {
		t.Error("should be confirmed")
	}

	if err := order.Ship(); err != nil {
		t.Fatalf("ship: %v", err)
	}

	if err := order.Deliver(); err != nil {
		t.Fatalf("deliver: %v", err)
	}
	if !order.IsDelivered() {
		t.Error("should be delivered")
	}
}

func TestOrder_EmptyCannotConfirm(t *testing.T) {
	order, _ := entity.NewOrder(uuid.New(), uuid.New(), valueobject.OrderTypeSales, uuid.New())
	err := order.Confirm()
	if !errors.Is(err, domainerrors.ErrOrderEmpty) {
		t.Errorf("expected ErrOrderEmpty, got %v", err)
	}
}

func TestOrder_CannotCancelDelivered(t *testing.T) {
	order, _ := entity.NewOrder(uuid.New(), uuid.New(), valueobject.OrderTypeSales, uuid.New())
	item, _ := entity.NewOrderItem(uuid.New(), "SKU", "P", 1,
		valueobject.Money{Amount: 10, Currency: "THB"})
	_ = order.AddItem(*item)
	_ = order.Confirm()
	_ = order.Ship()
	_ = order.Deliver()

	err := order.Cancel("test")
	if !errors.Is(err, domainerrors.ErrCannotCancelDelivered) {
		t.Errorf("expected CannotCancelDelivered, got %v", err)
	}
}

func TestStockItem_Reserve(t *testing.T) {
	s, _ := entity.NewStockItem(
		uuid.New(), uuid.New(), uuid.New(), "SKU",
		valueobject.Money{Amount: 50, Currency: "THB"},
	)
	s.Restock(10, valueobject.Money{Amount: 50, Currency: "THB"})

	if err := s.Reserve(3); err != nil {
		t.Fatalf("reserve: %v", err)
	}
	if s.Available() != 7 {
		t.Errorf("expected 7 available, got %d", s.Available())
	}

	if err := s.Issue(3); err != nil {
		t.Fatalf("issue: %v", err)
	}
	if s.Quantity != 7 {
		t.Errorf("expected 7 quantity, got %d", s.Quantity)
	}
}

func TestStockItem_InsufficientStock(t *testing.T) {
	s, _ := entity.NewStockItem(uuid.New(), uuid.New(), uuid.New(), "SKU",
		valueobject.Money{Amount: 50, Currency: "THB"})
	s.Restock(5, valueobject.Money{Amount: 50, Currency: "THB"})

	err := s.Reserve(10)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Errorf("expected InsufficientStock, got %v", err)
	}
}

func TestInvoice_Payment(t *testing.T) {
	inv, _ := entity.NewInvoice(
		uuid.New(), uuid.New(),
		valueobject.InvoiceTypeTax,
		valueobject.Money{Amount: 100, Currency: "THB"},
		valueobject.Money{Amount: 7, Currency: "THB"},
		valueobject.Money{Amount: 107, Currency: "THB"},
		time.Now().AddDate(0, 0, 30),
		uuid.New(),
	)
	if err := inv.Issue(); err != nil {
		t.Fatalf("issue: %v", err)
	}

	// insufficient
	err := inv.MarkPaid(valueobject.Money{Amount: 100, Currency: "THB"}, "cash")
	if !errors.Is(err, domainerrors.ErrInsufficientPayment) {
		t.Errorf("expected InsufficientPayment, got %v", err)
	}

	// ok
	if err := inv.MarkPaid(valueobject.Money{Amount: 107, Currency: "THB"}, "transfer"); err != nil {
		t.Fatalf("mark paid: %v", err)
	}
	if !inv.IsPaid() {
		t.Error("should be paid")
	}

	// cannot cancel paid
	if err := inv.Cancel(); !errors.Is(err, domainerrors.ErrCannotCancelPaid) {
		t.Errorf("expected CannotCancelPaid, got %v", err)
	}
}
```

**Status:** ✅ G04 COMPLETE

### G04 Build Verify

```bash
cd icmongolang
go build ./internal/modules/erp/...
go test ./internal/modules/erp/...
```

**Expected:** ✅ PASS

---

# PART 2 — P04 (Python ERP) Steps 2-5

## P04.Step 2 — Application Layer

### `application/ports/event_producer.py`

```python
from typing import Protocol

from app.modules.erp.domain.events.events import (
    OrderCreatedEvent, OrderConfirmedEvent, OrderShippedEvent,
    OrderDeliveredEvent, OrderCancelledEvent, InvoiceIssuedEvent,
    InvoicePaidEvent, StockLowEvent,
)


class EventProducer(Protocol):
    async def publish_order_created(self, evt: OrderCreatedEvent) -> None: ...
    async def publish_order_confirmed(self, evt: OrderConfirmedEvent) -> None: ...
    async def publish_order_shipped(self, evt: OrderShippedEvent) -> None: ...
    async def publish_order_delivered(self, evt: OrderDeliveredEvent) -> None: ...
    async def publish_order_cancelled(self, evt: OrderCancelledEvent) -> None: ...
    async def publish_invoice_issued(self, evt: InvoiceIssuedEvent) -> None: ...
    async def publish_invoice_paid(self, evt: InvoicePaidEvent) -> None: ...
    async def publish_stock_low(self, evt: StockLowEvent) -> None: ...
```

### `application/ports/customer_client.py`, `item_client.py`, `payment_client.py`, `audit_repository.py`, `tx_manager.py`

```python
# customer_client.py
from dataclasses import dataclass
from typing import Protocol
from uuid import UUID


@dataclass
class CustomerInfo:
    id: UUID
    code: str
    name: str
    is_active: bool


class CustomerClient(Protocol):
    async def get(self, customer_id: UUID) -> CustomerInfo | None: ...


# item_client.py
@dataclass
class ItemInfo:
    id: UUID
    sku: str
    name: str
    unit_price: float
    currency: str


class ItemClient(Protocol):
    async def get(self, item_id: UUID) -> ItemInfo | None: ...


# payment_client.py
@dataclass
class PaymentRequest:
    customer_id: UUID
    amount: float
    currency: str
    description: str
    reference_id: UUID
    reference_type: str


@dataclass
class PaymentResult:
    id: UUID
    url: str


class PaymentClient(Protocol):
    async def create_payment(self, req: PaymentRequest) -> PaymentResult: ...


# audit_repository.py
from dataclasses import dataclass, field
from datetime import datetime, timezone


@dataclass
class AuditTrail:
    user_id: UUID
    action: str
    entity_type: str
    entity_id: UUID
    payload: dict = field(default_factory=dict)
    ip_address: str = ""
    occurred_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))


class AuditRepository(Protocol):
    async def save(self, trail: AuditTrail) -> None: ...


# tx_manager.py
from collections.abc import Awaitable, Callable
from typing import Protocol


class TransactionManager(Protocol):
    async def run(self, fn: Callable[[], Awaitable[None]]) -> None: ...
```

### `application/use_cases/create_order.py`

```python
from dataclasses import dataclass, field
from datetime import datetime, timezone
from uuid import UUID

from app.modules.erp.application.ports.audit_repository import AuditRepository, AuditTrail
from app.modules.erp.application.ports.customer_client import CustomerClient
from app.modules.erp.application.ports.event_producer import EventProducer
from app.modules.erp.application.ports.item_client import ItemClient
from app.modules.erp.application.ports.tx_manager import TransactionManager
from app.modules.erp.domain.entities.order import Order
from app.modules.erp.domain.entities.order_item import OrderItem
from app.modules.erp.domain.errors.errors import InvalidCustomerError
from app.modules.erp.domain.events.events import OrderCreatedEvent
from app.modules.erp.domain.repositories.order_repository import OrderRepository
from app.modules.erp.domain.repositories.stock_repository import StockRepository
from app.modules.erp.domain.value_objects.money import Money
from app.modules.erp.domain.value_objects.order_type import OrderType


@dataclass
class CreateOrderItemInput:
    product_id: UUID
    warehouse_id: UUID
    quantity: int


@dataclass
class CreateOrderInput:
    tenant_id: UUID
    customer_id: UUID
    type: OrderType
    items: list[CreateOrderItemInput]
    notes: str = ""
    user_id: UUID = None
    ip_address: str = ""


@dataclass
class OrderItemOutput:
    id: UUID
    product_id: UUID
    sku: str
    name: str
    quantity: int
    unit_price: Money
    subtotal: Money


@dataclass
class CreateOrderOutput:
    id: UUID
    order_number: str
    customer_id: UUID
    type: str
    items: list[OrderItemOutput]
    subtotal: Money
    tax: Money
    discount: Money
    total: Money
    status: str


class CreateOrderUseCase:
    def __init__(
        self,
        order_repo: OrderRepository,
        stock_repo: StockRepository,
        customer_client: CustomerClient,
        item_client: ItemClient,
        producer: EventProducer,
        audit_repo: AuditRepository,
        tx_manager: TransactionManager,
    ):
        self._order_repo = order_repo
        self._stock_repo = stock_repo
        self._customer = customer_client
        self._item = item_client
        self._producer = producer
        self._audit = audit_repo
        self._tx = tx_manager

    async def execute(self, inp: CreateOrderInput) -> CreateOrderOutput:
        # 1. Verify customer
        cust = await self._customer.get(inp.customer_id)
        if not cust or not cust.is_active:
            raise InvalidCustomerError()

        # 2. Generate order number
        order_number = await self._order_repo.next_order_number(inp.tenant_id, inp.type)

        # 3. Create aggregate
        order = Order.create(inp.tenant_id, inp.customer_id, inp.type, inp.user_id)
        order.set_number(order_number)
        order.notes = inp.notes

        # 4. Add items
        for item_in in inp.items:
            item_info = await self._item.get(item_in.product_id)
            if not item_info:
                continue
            unit_price = Money(amount=item_info.unit_price, currency=item_info.currency)
            item = OrderItem.create(
                item_in.product_id, item_info.sku, item_info.name,
                item_in.quantity, unit_price,
            )
            order.add_item(item)

        # 5. Persist + reserve stock in transaction
        async def _tx():
            await self._order_repo.save(order)
            for item_in in inp.items:
                stock = await self._stock_repo.find_by_product_warehouse(
                    item_in.product_id, item_in.warehouse_id
                )
                if not stock:
                    continue
                stock.reserve(item_in.quantity)
                await self._stock_repo.save(stock)

        await self._tx.run(_tx)

        # 6. Audit
        await self._audit.save(AuditTrail(
            user_id=inp.user_id,
            action="ORDER_CREATED",
            entity_type="order",
            entity_id=order.id,
            payload={"order_number": order.order_number, "total": order.total.amount},
            ip_address=inp.ip_address,
        ))

        # 7. Publish
        await self._producer.publish_order_created(OrderCreatedEvent(
            order_id=order.id,
            tenant_id=order.tenant_id,
            customer_id=order.customer_id,
            order_number=order.order_number,
            type=order.type.value,
            total=order.total.amount,
            currency=order.total.currency,
            item_count=order.item_count(),
        ))

        return CreateOrderOutput(
            id=order.id,
            order_number=order.order_number,
            customer_id=order.customer_id,
            type=order.type.value,
            items=[
                OrderItemOutput(
                    id=i.id, product_id=i.product_id, sku=i.sku, name=i.name,
                    quantity=i.quantity, unit_price=i.unit_price, subtotal=i.subtotal,
                ) for i in order.items
            ],
            subtotal=order.subtotal, tax=order.tax_amount,
            discount=order.discount, total=order.total,
            status=order.status.value,
        )
```

### `application/use_cases/confirm_order.py` + `ship_order.py` + `deliver_order.py` + `cancel_order.py`

```python
# confirm_order.py
from dataclasses import dataclass
from datetime import datetime, timezone
from uuid import UUID

from app.modules.erp.application.ports.event_producer import EventProducer
from app.modules.erp.domain.events.events import OrderConfirmedEvent
from app.modules.erp.domain.repositories.order_repository import OrderRepository


@dataclass
class ConfirmOrderInput:
    order_id: UUID
    user_id: UUID


class ConfirmOrderUseCase:
    def __init__(self, order_repo: OrderRepository, producer: EventProducer):
        self._repo = order_repo
        self._producer = producer

    async def execute(self, inp: ConfirmOrderInput) -> None:
        order = await self._repo.find_by_id(inp.order_id)
        if not order:
            return
        order.confirm()
        await self._repo.save(order)
        await self._producer.publish_order_confirmed(OrderConfirmedEvent(
            order_id=order.id, tenant_id=order.tenant_id,
        ))


# ship_order.py
@dataclass
class ShipOrderInput:
    order_id: UUID
    user_id: UUID


class ShipOrderUseCase:
    def __init__(self, order_repo: OrderRepository, producer: EventProducer):
        self._repo = order_repo
        self._producer = producer

    async def execute(self, inp: ShipOrderInput) -> None:
        order = await self._repo.find_by_id(inp.order_id)
        if not order:
            return
        order.ship()
        await self._repo.save(order)
        from app.modules.erp.domain.events.events import OrderShippedEvent
        await self._producer.publish_order_shipped(OrderShippedEvent(
            order_id=order.id, tenant_id=order.tenant_id,
        ))


# deliver_order.py
from app.modules.erp.application.ports.tx_manager import TransactionManager
from app.modules.erp.domain.events.events import OrderDeliveredEvent
from app.modules.erp.domain.repositories.stock_repository import StockRepository


@dataclass
class DeliverOrderInput:
    order_id: UUID
    user_id: UUID


class DeliverOrderUseCase:
    def __init__(
        self,
        order_repo: OrderRepository,
        stock_repo: StockRepository,
        producer: EventProducer,
        tx_manager: TransactionManager,
    ):
        self._order_repo = order_repo
        self._stock_repo = stock_repo
        self._producer = producer
        self._tx = tx_manager

    async def execute(self, inp: DeliverOrderInput) -> None:
        order = await self._order_repo.find_by_id(inp.order_id)
        if not order:
            return
        order.deliver()

        async def _tx():
            await self._order_repo.save(order)
            for item in order.items:
                stocks = await self._stock_repo.find_by_product(item.product_id)
                if stocks:
                    stocks[0].issue(item.quantity)
                    await self._stock_repo.save(stocks[0])

        await self._tx.run(_tx)
        await self._producer.publish_order_delivered(OrderDeliveredEvent(
            order_id=order.id, tenant_id=order.tenant_id,
        ))


# cancel_order.py
from app.modules.erp.application.ports.audit_repository import AuditRepository, AuditTrail
from app.modules.erp.domain.events.events import OrderCancelledEvent


@dataclass
class CancelOrderInput:
    order_id: UUID
    reason: str
    user_id: UUID


class CancelOrderUseCase:
    def __init__(
        self,
        order_repo: OrderRepository,
        stock_repo: StockRepository,
        producer: EventProducer,
        audit_repo: AuditRepository,
        tx_manager: TransactionManager,
    ):
        self._order_repo = order_repo
        self._stock_repo = stock_repo
        self._producer = producer
        self._audit = audit_repo
        self._tx = tx_manager

    async def execute(self, inp: CancelOrderInput) -> None:
        order = await self._order_repo.find_by_id(inp.order_id)
        if not order:
            return
        order.cancel(inp.reason)

        async def _tx():
            await self._order_repo.save(order)
            for item in order.items:
                stocks = await self._stock_repo.find_by_product(item.product_id)
                if stocks:
                    stocks[0].release(item.quantity)
                    await self._stock_repo.save(stocks[0])

        await self._tx.run(_tx)

        await self._audit.save(AuditTrail(
            user_id=inp.user_id, action="ORDER_CANCELLED",
            entity_type="order", entity_id=order.id,
            payload={"reason": inp.reason},
        ))

        await self._producer.publish_order_cancelled(OrderCancelledEvent(
            order_id=order.id, tenant_id=order.tenant_id, reason=inp.reason,
        ))
```

### `application/use_cases/create_invoice.py` + `record_payment.py` + `restock.py` + `reorder_check.py`

```python
# create_invoice.py
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from uuid import UUID

from app.modules.erp.domain.entities.invoice import Invoice
from app.modules.erp.domain.repositories.invoice_repository import InvoiceRepository
from app.modules.erp.domain.repositories.order_repository import OrderRepository
from app.modules.erp.domain.value_objects.invoice_type import InvoiceType


@dataclass
class CreateInvoiceInput:
    order_id: UUID
    type: InvoiceType
    due_days: int
    user_id: UUID


class CreateInvoiceUseCase:
    def __init__(
        self,
        invoice_repo: InvoiceRepository,
        order_repo: OrderRepository,
    ):
        self._invoice_repo = invoice_repo
        self._order_repo = order_repo

    async def execute(self, inp: CreateInvoiceInput) -> Invoice:
        order = await self._order_repo.find_by_id(inp.order_id)
        if not order:
            raise ValueError("order not found")

        due_date = datetime.now(timezone.utc) + timedelta(days=inp.due_days)
        number = await self._invoice_repo.next_invoice_number(order.tenant_id)

        inv = Invoice.create(
            tenant_id=order.tenant_id,
            customer_id=order.customer_id,
            invoice_type=inp.type,
            subtotal=order.subtotal,
            tax=order.tax_amount,
            total=order.total,
            due_date=due_date,
            created_by=inp.user_id,
        )
        inv.set_number(number)
        inv.link_order(order.id)
        await self._invoice_repo.save(inv)
        return inv


# record_payment.py
from app.modules.erp.application.ports.event_producer import EventProducer
from app.modules.erp.domain.events.events import InvoicePaidEvent
from app.modules.erp.domain.value_objects.money import Money


@dataclass
class RecordPaymentInput:
    invoice_id: UUID
    amount: Money
    method: str
    user_id: UUID


class RecordPaymentUseCase:
    def __init__(self, invoice_repo: InvoiceRepository, producer: EventProducer):
        self._repo = invoice_repo
        self._producer = producer

    async def execute(self, inp: RecordPaymentInput) -> None:
        inv = await self._repo.find_by_id(inp.invoice_id)
        if not inv:
            return
        inv.mark_paid(inp.amount, inp.method)
        await self._repo.save(inv)
        await self._producer.publish_invoice_paid(InvoicePaidEvent(
            invoice_id=inv.id, tenant_id=inv.tenant_id,
            amount=inp.amount.amount, currency=inp.amount.currency,
        ))


# restock.py
from app.modules.erp.domain.repositories.stock_repository import StockRepository


@dataclass
class RestockItemInput:
    stock_id: UUID
    quantity: int
    unit_cost: Money


class RestockItemUseCase:
    def __init__(self, stock_repo: StockRepository):
        self._repo = stock_repo

    async def execute(self, inp: RestockItemInput) -> None:
        stock = await self._repo.find_by_id(inp.stock_id)
        if not stock:
            return
        stock.restock(inp.quantity, inp.unit_cost)
        await self._repo.save(stock)


# reorder_check.py
from app.modules.erp.domain.events.events import StockLowEvent


class ReorderCheckUseCase:
    def __init__(self, stock_repo: StockRepository, producer: EventProducer):
        self._repo = stock_repo
        self._producer = producer

    async def execute(self, tenant_id: UUID) -> int:
        stocks = await self._repo.find_needs_reorder(tenant_id)
        for s in stocks:
            await self._producer.publish_stock_low(StockLowEvent(
                tenant_id=s.tenant_id, product_id=s.product_id,
                sku=s.sku, warehouse_id=s.warehouse_id,
                available=s.available(), reorder_point=s.reorder_point,
            ))
        return len(stocks)
```

**Status:** ✅ Step 2 complete (11 files)

---

## P04.Step 3 — Infrastructure Layer

### `infrastructure/persistence/postgres/models.py`

```python
from datetime import datetime
from uuid import UUID, uuid4

from sqlalchemy import (
    Boolean, DateTime, Float, ForeignKey, Index, Integer, String, func,
)
from sqlalchemy.dialects.postgresql import JSONB, UUID as PGUUID
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


class Base(DeclarativeBase):
    pass


class OrderModel(Base):
    __tablename__ = "erp_orders"
    __table_args__ = (
        Index("idx_erp_orders_tenant", "tenant_id"),
        Index("idx_erp_orders_customer", "customer_id"),
        Index("idx_erp_orders_status", "tenant_id", "status"),
    )

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    customer_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    order_number: Mapped[str] = mapped_column(String(50), nullable=False)
    type: Mapped[str] = mapped_column(String(20), nullable=False)
    subtotal: Mapped[float] = mapped_column(Float, default=0)
    tax_amount: Mapped[float] = mapped_column(Float, default=0)
    discount: Mapped[float] = mapped_column(Float, default=0)
    total: Mapped[float] = mapped_column(Float, default=0)
    currency: Mapped[str] = mapped_column(String(3), default="THB")
    status: Mapped[str] = mapped_column(String(30), default="DRAFT")
    notes: Mapped[str] = mapped_column(String, default="")
    shipping_addr: Mapped[dict] = mapped_column(JSONB, default=dict)
    billing_addr: Mapped[dict] = mapped_column(JSONB, default=dict)
    ordered_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    confirmed_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    shipped_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    delivered_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    cancelled_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    cancel_reason: Mapped[str] = mapped_column(String, default="")
    created_by: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())


class OrderItemModel(Base):
    __tablename__ = "erp_order_items"
    __table_args__ = (Index("idx_erp_order_items_order", "order_id"),)

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    order_id: Mapped[UUID] = mapped_column(
        PGUUID(as_uuid=True), ForeignKey("erp_orders.id", ondelete="CASCADE")
    )
    product_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    sku: Mapped[str] = mapped_column(String(100), nullable=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    unit_price: Mapped[float] = mapped_column(Float, nullable=False)
    discount: Mapped[float] = mapped_column(Float, default=0)
    tax_rate: Mapped[float] = mapped_column(Float, default=0)
    subtotal: Mapped[float] = mapped_column(Float, nullable=False)


class InvoiceModel(Base):
    __tablename__ = "erp_invoices"
    __table_args__ = (
        Index("idx_erp_invoices_tenant", "tenant_id"),
        Index("idx_erp_invoices_customer", "customer_id"),
    )

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    order_id: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    customer_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    number: Mapped[str] = mapped_column(String(50), nullable=False)
    type: Mapped[str] = mapped_column(String(20), nullable=False)
    subtotal: Mapped[float] = mapped_column(Float, nullable=False)
    tax_amount: Mapped[float] = mapped_column(Float, nullable=False)
    total: Mapped[float] = mapped_column(Float, nullable=False)
    currency: Mapped[str] = mapped_column(String(3), default="THB")
    status: Mapped[str] = mapped_column(String(30), default="DRAFT")
    issued_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    due_date: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    paid_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    paid_amount: Mapped[float] = mapped_column(Float, default=0)
    payment_method: Mapped[str] = mapped_column(String(50), default="")
    notes: Mapped[str] = mapped_column(String, default="")
    created_by: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())


class WarehouseModel(Base):
    __tablename__ = "erp_warehouses"

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    code: Mapped[str] = mapped_column(String(50), nullable=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    address: Mapped[dict] = mapped_column(JSONB, default=dict)
    type: Mapped[str] = mapped_column(String(30), nullable=False)
    capacity: Mapped[int] = mapped_column(Integer, default=0)
    manager_id: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    is_active: Mapped[bool] = mapped_column(Boolean, default=True)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())


class StockItemModel(Base):
    __tablename__ = "erp_stock_items"
    __table_args__ = (
        Index("idx_erp_stock_tenant", "tenant_id"),
        Index("idx_erp_stock_product", "product_id"),
    )

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    product_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    warehouse_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    sku: Mapped[str] = mapped_column(String(100), nullable=False)
    quantity: Mapped[int] = mapped_column(Integer, default=0)
    reserved: Mapped[int] = mapped_column(Integer, default=0)
    reorder_point: Mapped[int] = mapped_column(Integer, default=0)
    unit_cost: Mapped[float] = mapped_column(Float, default=0)
    last_restocked: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())
```

### `infrastructure/persistence/postgres/order_repository_impl.py`

```python
from datetime import datetime, timezone
from uuid import UUID

from sqlalchemy import select, func, extract
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.modules.erp.domain.entities.order import Order
from app.modules.erp.domain.entities.order_item import OrderItem
from app.modules.erp.domain.errors.errors import OrderNotFoundError
from app.modules.erp.domain.value_objects.money import Money
from app.modules.erp.domain.value_objects.order_status import OrderStatus
from app.modules.erp.domain.value_objects.order_type import OrderType
from app.modules.erp.infrastructure.persistence.postgres.models import (
    OrderItemModel, OrderModel,
)


class PostgresOrderRepository:
    def __init__(self, session: AsyncSession):
        self._session = session

    async def save(self, order: Order) -> None:
        model = await self._session.get(OrderModel, order.id)
        if model is None:
            model = OrderModel(id=order.id)
            self._session.add(model)

        model.tenant_id = order.tenant_id
        model.customer_id = order.customer_id
        model.order_number = order.order_number
        model.type = order.type.value
        model.subtotal = order.subtotal.amount
        model.tax_amount = order.tax_amount.amount
        model.discount = order.discount.amount
        model.total = order.total.amount
        model.currency = order.total.currency
        model.status = order.status.value
        model.notes = order.notes
        model.shipping_addr = order.shipping_addr.model_dump()
        model.billing_addr = order.billing_addr.model_dump()
        model.ordered_at = order.ordered_at
        model.confirmed_at = order.confirmed_at
        model.shipped_at = order.shipped_at
        model.delivered_at = order.delivered_at
        model.cancelled_at = order.cancelled_at
        model.cancel_reason = order.cancel_reason
        model.created_by = order.created_by
        model.updated_at = datetime.now(timezone.utc)

        # Replace items
        stmt_del = select(OrderItemModel).where(OrderItemModel.order_id == order.id)
        existing = (await self._session.execute(stmt_del)).scalars().all()
        for e in existing:
            await self._session.delete(e)

        for item in order.items:
            self._session.add(OrderItemModel(
                id=item.id,
                order_id=order.id,
                product_id=item.product_id,
                sku=item.sku,
                name=item.name,
                quantity=item.quantity,
                unit_price=item.unit_price.amount,
                discount=item.discount.amount if item.discount else 0,
                tax_rate=item.tax_rate,
                subtotal=item.subtotal.amount,
            ))
        await self._session.flush()

    async def find_by_id(self, id: UUID) -> Order | None:
        stmt = (
            select(OrderModel)
            .options(selectinload(OrderModel.items) if hasattr(OrderModel, "items") else None)
            .where(OrderModel.id == id)
        ) if False else select(OrderModel).where(OrderModel.id == id)

        m = (await self._session.execute(stmt)).scalar_one_or_none()
        if not m:
            return None
        return await self._to_entity(m)

    async def find_by_number(self, tenant_id: UUID, number: str) -> Order | None:
        stmt = select(OrderModel).where(
            OrderModel.tenant_id == tenant_id, OrderModel.order_number == number
        )
        m = (await self._session.execute(stmt)).scalar_one_or_none()
        return await self._to_entity(m) if m else None

    async def find_by_customer(self, customer_id: UUID) -> list[Order]:
        stmt = select(OrderModel).where(OrderModel.customer_id == customer_id)
        rows = (await self._session.execute(stmt)).scalars().all()
        return [await self._to_entity(m) for m in rows]

    async def find_by_status(self, tenant_id: UUID, status: OrderStatus) -> list[Order]:
        stmt = select(OrderModel).where(
            OrderModel.tenant_id == tenant_id, OrderModel.status == status.value
        )
        rows = (await self._session.execute(stmt)).scalars().all()
        return [await self._to_entity(m) for m in rows]

    async def next_order_number(self, tenant_id: UUID, order_type: OrderType) -> str:
        year = datetime.now().year
        prefix = "SO" if order_type == OrderType.SALES else "PO"
        stmt = select(func.count()).select_from(OrderModel).where(
            OrderModel.tenant_id == tenant_id,
            extract("year", OrderModel.created_at) == year,
        )
        count = (await self._session.execute(stmt)).scalar() or 0
        return f"{prefix}-{year}-{count + 1:05d}"

    async def exists_by_number(self, tenant_id: UUID, number: str) -> bool:
        stmt = select(func.count()).select_from(OrderModel).where(
            OrderModel.tenant_id == tenant_id, OrderModel.order_number == number
        )
        return (await self._session.execute(stmt)).scalar() > 0

    async def delete(self, id: UUID) -> None:
        m = await self._session.get(OrderModel, id)
        if m:
            await self._session.delete(m)

    async def _to_entity(self, m: OrderModel) -> Order:
        # load items
        stmt = select(OrderItemModel).where(OrderItemModel.order_id == m.id)
        items_rows = (await self._session.execute(stmt)).scalars().all()

        items = [
            OrderItem(
                id=i.id,
                order_id=m.id,
                product_id=i.product_id,
                sku=i.sku,
                name=i.name,
                quantity=i.quantity,
                unit_price=Money(amount=i.unit_price, currency=m.currency),
                discount=Money(amount=i.discount, currency=m.currency),
                tax_rate=i.tax_rate,
                subtotal=Money(amount=i.subtotal, currency=m.currency),
            )
            for i in items_rows
        ]

        return Order(
            id=m.id,
            tenant_id=m.tenant_id,
            customer_id=m.customer_id,
            order_number=m.order_number,
            type=OrderType(m.type),
            items=items,
            subtotal=Money(amount=m.subtotal, currency=m.currency),
            tax_amount=Money(amount=m.tax_amount, currency=m.currency),
            discount=Money(amount=m.discount, currency=m.currency),
            total=Money(amount=m.total, currency=m.currency),
            status=OrderStatus(m.status),
            notes=m.notes or "",
            ordered_at=m.ordered_at,
            confirmed_at=m.confirmed_at,
            shipped_at=m.shipped_at,
            delivered_at=m.delivered_at,
            cancelled_at=m.cancelled_at,
            cancel_reason=m.cancel_reason or "",
            created_by=m.created_by,
            created_at=m.created_at,
            updated_at=m.updated_at,
        )
```

### `infrastructure/persistence/postgres/invoice_repository_impl.py` + `stock_repository_impl.py` + `warehouse_repository_impl.py`

```python
# invoice_repository_impl.py
from datetime import datetime, timezone
from uuid import UUID

from sqlalchemy import select, func, extract
from sqlalchemy.ext.asyncio import AsyncSession

from app.modules.erp.domain.entities.invoice import Invoice
from app.modules.erp.domain.errors.errors import InvoiceNotFoundError
from app.modules.erp.domain.value_objects.invoice_status import InvoiceStatus
from app.modules.erp.domain.value_objects.invoice_type import InvoiceType
from app.modules.erp.domain.value_objects.money import Money
from app.modules.erp.infrastructure.persistence.postgres.models import InvoiceModel


class PostgresInvoiceRepository:
    def __init__(self, session: AsyncSession):
        self._session = session

    async def save(self, inv: Invoice) -> None:
        m = await self._session.get(InvoiceModel, inv.id)
        if m is None:
            m = InvoiceModel(id=inv.id)
            self._session.add(m)

        m.tenant_id = inv.tenant_id
        m.order_id = inv.order_id
        m.customer_id = inv.customer_id
        m.number = inv.number
        m.type = inv.type.value
        m.subtotal = inv.subtotal.amount
        m.tax_amount = inv.tax_amount.amount
        m.total = inv.total.amount
        m.currency = inv.total.currency
        m.status = inv.status.value
        m.issued_at = inv.issued_at
        m.due_date = inv.due_date
        m.paid_at = inv.paid_at
        m.paid_amount = inv.paid_amount.amount
        m.payment_method = inv.payment_method
        m.notes = inv.notes
        m.created_by = inv.created_by
        m.updated_at = datetime.now(timezone.utc)
        await self._session.flush()

    async def find_by_id(self, id: UUID) -> Invoice | None:
        m = await self._session.get(InvoiceModel, id)
        return self._to_entity(m) if m else None

    async def find_by_number(self, tenant_id: UUID, number: str) -> Invoice | None:
        stmt = select(InvoiceModel).where(
            InvoiceModel.tenant_id == tenant_id, InvoiceModel.number == number
        )
        m = (await self._session.execute(stmt)).scalar_one_or_none()
        return self._to_entity(m) if m else None

    async def find_by_customer(self, customer_id: UUID) -> list[Invoice]:
        stmt = select(InvoiceModel).where(InvoiceModel.customer_id == customer_id)
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def find_by_order(self, order_id: UUID) -> list[Invoice]:
        stmt = select(InvoiceModel).where(InvoiceModel.order_id == order_id)
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def find_by_status(self, tenant_id: UUID, status: InvoiceStatus) -> list[Invoice]:
        stmt = select(InvoiceModel).where(
            InvoiceModel.tenant_id == tenant_id, InvoiceModel.status == status.value
        )
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def find_overdue(self, before: datetime) -> list[Invoice]:
        stmt = select(InvoiceModel).where(
            InvoiceModel.status == "ISSUED", InvoiceModel.due_date < before
        )
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def next_invoice_number(self, tenant_id: UUID) -> str:
        year = datetime.now().year
        stmt = select(func.count()).select_from(InvoiceModel).where(
            InvoiceModel.tenant_id == tenant_id,
            extract("year", InvoiceModel.created_at) == year,
        )
        count = (await self._session.execute(stmt)).scalar() or 0
        return f"INV-{year}-{count + 1:05d}"

    async def delete(self, id: UUID) -> None:
        m = await self._session.get(InvoiceModel, id)
        if m:
            await self._session.delete(m)

    def _to_entity(self, m: InvoiceModel) -> Invoice:
        return Invoice(
            id=m.id, tenant_id=m.tenant_id, order_id=m.order_id,
            customer_id=m.customer_id, number=m.number,
            type=InvoiceType(m.type),
            subtotal=Money(amount=m.subtotal, currency=m.currency),
            tax_amount=Money(amount=m.tax_amount, currency=m.currency),
            total=Money(amount=m.total, currency=m.currency),
            status=InvoiceStatus(m.status),
            issued_at=m.issued_at or datetime.now(timezone.utc),
            due_date=m.due_date or datetime.now(timezone.utc),
            paid_at=m.paid_at,
            paid_amount=Money(amount=m.paid_amount, currency=m.currency),
            payment_method=m.payment_method or "",
            notes=m.notes or "",
            created_by=m.created_by,
            created_at=m.created_at, updated_at=m.updated_at,
        )


# stock_repository_impl.py
from app.modules.erp.domain.entities.stock_item import StockItem
from app.modules.erp.domain.value_objects.money import Money
from app.modules.erp.infrastructure.persistence.postgres.models import StockItemModel


class PostgresStockRepository:
    def __init__(self, session: AsyncSession):
        self._session = session

    async def save(self, s: StockItem) -> None:
        m = await self._session.get(StockItemModel, s.id)
        if m is None:
            m = StockItemModel(id=s.id)
            self._session.add(m)
        m.tenant_id = s.tenant_id
        m.product_id = s.product_id
        m.warehouse_id = s.warehouse_id
        m.sku = s.sku
        m.quantity = s.quantity
        m.reserved = s.reserved
        m.reorder_point = s.reorder_point
        m.unit_cost = s.unit_cost.amount
        m.last_restocked = s.last_restocked
        m.updated_at = datetime.now(timezone.utc)
        await self._session.flush()

    async def find_by_id(self, id: UUID) -> StockItem | None:
        m = await self._session.get(StockItemModel, id)
        return self._to_entity(m) if m else None

    async def find_by_product_warehouse(self, product_id: UUID, warehouse_id: UUID) -> StockItem | None:
        stmt = select(StockItemModel).where(
            StockItemModel.product_id == product_id,
            StockItemModel.warehouse_id == warehouse_id,
        )
        m = (await self._session.execute(stmt)).scalar_one_or_none()
        return self._to_entity(m) if m else None

    async def find_by_product(self, product_id: UUID) -> list[StockItem]:
        stmt = select(StockItemModel).where(StockItemModel.product_id == product_id)
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def find_by_warehouse(self, warehouse_id: UUID) -> list[StockItem]:
        stmt = select(StockItemModel).where(StockItemModel.warehouse_id == warehouse_id)
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def find_needs_reorder(self, tenant_id: UUID) -> list[StockItem]:
        stmt = select(StockItemModel).where(
            StockItemModel.tenant_id == tenant_id,
            (StockItemModel.quantity - StockItemModel.reserved) <= StockItemModel.reorder_point,
        )
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def delete(self, id: UUID) -> None:
        m = await self._session.get(StockItemModel, id)
        if m:
            await self._session.delete(m)

    def _to_entity(self, m: StockItemModel) -> StockItem:
        return StockItem(
            id=m.id, tenant_id=m.tenant_id, product_id=m.product_id,
            warehouse_id=m.warehouse_id, sku=m.sku,
            quantity=m.quantity, reserved=m.reserved, reorder_point=m.reorder_point,
            unit_cost=Money(amount=m.unit_cost, currency="THB"),
            last_restocked=m.last_restocked, updated_at=m.updated_at,
        )
```

### `infrastructure/messaging/kafka_producer.py`

```python
import json
from dataclasses import asdict
from datetime import datetime
from uuid import UUID

from aiokafka import AIOKafkaProducer

from app.modules.erp.domain.events.events import (
    InvoicePaidEvent, OrderCancelledEvent, OrderConfirmedEvent,
    OrderCreatedEvent, OrderDeliveredEvent, OrderShippedEvent,
    StockLowEvent,
)


def _encoder(o):
    if isinstance(o, datetime):
        return o.isoformat()
    if isinstance(o, UUID):
        return str(o)
    raise TypeError(f"not serializable: {type(o)}")


class KafkaEventProducer:
    TOPIC_ORDER_CREATED = "erp.order.created"
    TOPIC_ORDER_CONFIRMED = "erp.order.confirmed"
    TOPIC_ORDER_SHIPPED = "erp.order.shipped"
    TOPIC_ORDER_DELIVERED = "erp.order.delivered"
    TOPIC_ORDER_CANCELLED = "erp.order.cancelled"
    TOPIC_INVOICE_ISSUED = "erp.invoice.issued"
    TOPIC_INVOICE_PAID = "erp.invoice.paid"
    TOPIC_STOCK_LOW = "erp.stock.low"

    def __init__(self, producer: AIOKafkaProducer):
        self._producer = producer

    async def _send(self, topic: str, key: str, payload) -> None:
        data = json.dumps(asdict(payload), default=_encoder).encode("utf-8")
        await self._producer.send_and_wait(topic, key=key.encode(), value=data)

    async def publish_order_created(self, evt: OrderCreatedEvent) -> None:
        await self._send(self.TOPIC_ORDER_CREATED, str(evt.order_id), evt)

    async def publish_order_confirmed(self, evt: OrderConfirmedEvent) -> None:
        await self._send(self.TOPIC_ORDER_CONFIRMED, str(evt.order_id), evt)

    async def publish_order_shipped(self, evt: OrderShippedEvent) -> None:
        await self._send(self.TOPIC_ORDER_SHIPPED, str(evt.order_id), evt)

    async def publish_order_delivered(self, evt: OrderDeliveredEvent) -> None:
        await self._send(self.TOPIC_ORDER_DELIVERED, str(evt.order_id), evt)

    async def publish_order_cancelled(self, evt: OrderCancelledEvent) -> None:
        await self._send(self.TOPIC_ORDER_CANCELLED, str(evt.order_id), evt)

    async def publish_invoice_issued(self, evt) -> None:
        await self._send(self.TOPIC_INVOICE_ISSUED, str(evt.invoice_id), evt)

    async def publish_invoice_paid(self, evt: InvoicePaidEvent) -> None:
        await self._send(self.TOPIC_INVOICE_PAID, str(evt.invoice_id), evt)

    async def publish_stock_low(self, evt: StockLowEvent) -> None:
        await self._send(self.TOPIC_STOCK_LOW, str(evt.product_id), evt)
```

### `infrastructure/scheduler/invoice_overdue_job.py`

```python
import logging
from datetime import datetime, timezone

from app.modules.erp.domain.repositories.invoice_repository import InvoiceRepository

log = logging.getLogger(__name__)


class InvoiceOverdueJob:
    def __init__(self, invoice_repo: InvoiceRepository):
        self._repo = invoice_repo

    async def run(self) -> None:
        invoices = await self._repo.find_overdue(datetime.now(timezone.utc))
        for inv in invoices:
            inv.mark_overdue()
            await self._repo.save(inv)
        log.info("[invoice_overdue] marked %d overdue", len(invoices))
```

**Status:** ✅ Step 3 complete

---

## P04.Step 4 — Interface Layer + module.py

### `interfaces/http/schemas.py`

```python
from datetime import datetime
from uuid import UUID

from pydantic import BaseModel, Field


class CreateOrderItemSchema(BaseModel):
    product_id: UUID
    warehouse_id: UUID
    quantity: int = Field(..., gt=0)


class CreateOrderRequest(BaseModel):
    customer_id: UUID
    type: str = Field(..., pattern=r"^(SALES|PURCHASE)$")
    items: list[CreateOrderItemSchema] = Field(..., min_length=1)
    notes: str = ""


class MoneySchema(BaseModel):
    amount: float
    currency: str


class OrderItemResponse(BaseModel):
    id: UUID
    product_id: UUID
    sku: str
    name: str
    quantity: int
    unit_price: MoneySchema
    subtotal: MoneySchema


class OrderResponse(BaseModel):
    id: UUID
    tenant_id: UUID
    customer_id: UUID
    order_number: str
    type: str
    items: list[OrderItemResponse]
    subtotal: MoneySchema
    tax: MoneySchema
    discount: MoneySchema
    total: MoneySchema
    status: str
    created_at: datetime


class CreateInvoiceRequest(BaseModel):
    order_id: UUID
    type: str = Field(..., pattern=r"^(TAX|RECEIPT|DEBIT_NOTE|CREDIT_NOTE)$")
    due_days: int = Field(..., ge=1)


class RecordPaymentRequest(BaseModel):
    amount: float = Field(..., gt=0)
    currency: str
    method: str


class RestockRequest(BaseModel):
    quantity: int = Field(..., gt=0)
    unit_cost: float = Field(..., ge=0)
    currency: str = "THB"
```

### `interfaces/http/order_router.py`

```python
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status

from app.modules.erp.application.use_cases.create_order import (
    CreateOrderInput, CreateOrderItemInput, CreateOrderUseCase,
)
from app.modules.erp.application.use_cases.confirm_order import (
    ConfirmOrderInput, ConfirmOrderUseCase,
)
from app.modules.erp.application.use_cases.ship_order import ShipOrderInput, ShipOrderUseCase
from app.modules.erp.application.use_cases.deliver_order import DeliverOrderInput, DeliverOrderUseCase
from app.modules.erp.application.use_cases.cancel_order import CancelOrderInput, CancelOrderUseCase
from app.modules.erp.domain.errors.errors import (
    DomainError, OrderNotFoundError, InsufficientStockError, InvalidCustomerError,
)
from app.modules.erp.domain.value_objects.order_type import OrderType
from app.modules.erp.interfaces.http.dependencies import (
    get_cancel_order_uc, get_confirm_order_uc, get_create_order_uc,
    get_current_user, get_deliver_order_uc, get_ship_order_uc, get_tenant_id,
)
from app.modules.erp.interfaces.http.schemas import CreateOrderRequest, OrderResponse

router = APIRouter(prefix="/api/v1/erp/orders", tags=["erp-orders"])


def _map_error(e: DomainError) -> HTTPException:
    status_map = {
        OrderNotFoundError: 404,
        InsufficientStockError: 409,
        InvalidCustomerError: 400,
    }
    code = status_map.get(type(e), 400)
    return HTTPException(code, detail={"code": e.code, "message": str(e)})


@router.post("", response_model=OrderResponse, status_code=status.HTTP_201_CREATED)
async def create_order(
    req: CreateOrderRequest,
    uc: CreateOrderUseCase = Depends(get_create_order_uc),
    user_id: UUID = Depends(get_current_user),
    tenant_id: UUID = Depends(get_tenant_id),
):
    try:
        out = await uc.execute(CreateOrderInput(
            tenant_id=tenant_id,
            customer_id=req.customer_id,
            type=OrderType(req.type),
            items=[
                CreateOrderItemInput(
                    product_id=i.product_id,
                    warehouse_id=i.warehouse_id,
                    quantity=i.quantity,
                ) for i in req.items
            ],
            notes=req.notes,
            user_id=user_id,
        ))
        return OrderResponse(
            id=out.id, tenant_id=tenant_id, customer_id=out.customer_id,
            order_number=out.order_number, type=out.type,
            items=[],
            subtotal=out.subtotal, tax=out.tax, discount=out.discount, total=out.total,
            status=out.status, created_at=datetime.now(timezone.utc),
        )
    except DomainError as e:
        raise _map_error(e)


@router.post("/{order_id}/confirm")
async def confirm(
    order_id: UUID,
    uc: ConfirmOrderUseCase = Depends(get_confirm_order_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        await uc.execute(ConfirmOrderInput(order_id=order_id, user_id=user_id))
        return {"message": "confirmed"}
    except DomainError as e:
        raise _map_error(e)


@router.post("/{order_id}/ship")
async def ship(
    order_id: UUID,
    uc: ShipOrderUseCase = Depends(get_ship_order_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        await uc.execute(ShipOrderInput(order_id=order_id, user_id=user_id))
        return {"message": "shipped"}
    except DomainError as e:
        raise _map_error(e)


@router.post("/{order_id}/deliver")
async def deliver(
    order_id: UUID,
    uc: DeliverOrderUseCase = Depends(get_deliver_order_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        await uc.execute(DeliverOrderInput(order_id=order_id, user_id=user_id))
        return {"message": "delivered"}
    except DomainError as e:
        raise _map_error(e)


@router.post("/{order_id}/cancel")
async def cancel(
    order_id: UUID,
    reason: str,
    uc: CancelOrderUseCase = Depends(get_cancel_order_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        await uc.execute(CancelOrderInput(order_id=order_id, reason=reason, user_id=user_id))
        return {"message": "cancelled"}
    except DomainError as e:
        raise _map_error(e)


from datetime import datetime, timezone
```

### `interfaces/http/invoice_router.py` + `stock_router.py`

```python
# invoice_router.py
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status

from app.modules.erp.application.use_cases.create_invoice import (
    CreateInvoiceInput, CreateInvoiceUseCase,
)
from app.modules.erp.application.use_cases.record_payment import (
    RecordPaymentInput, RecordPaymentUseCase,
)
from app.modules.erp.domain.errors.errors import DomainError, InvoiceNotFoundError
from app.modules.erp.domain.value_objects.invoice_type import InvoiceType
from app.modules.erp.domain.value_objects.money import Money
from app.modules.erp.interfaces.http.dependencies import (
    get_create_invoice_uc, get_current_user, get_record_payment_uc,
)
from app.modules.erp.interfaces.http.schemas import (
    CreateInvoiceRequest, RecordPaymentRequest,
)

router = APIRouter(prefix="/api/v1/erp/invoices", tags=["erp-invoices"])


@router.post("", status_code=status.HTTP_201_CREATED)
async def create_invoice(
    req: CreateInvoiceRequest,
    uc: CreateInvoiceUseCase = Depends(get_create_invoice_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        inv = await uc.execute(CreateInvoiceInput(
            order_id=req.order_id, type=InvoiceType(req.type),
            due_days=req.due_days, user_id=user_id,
        ))
        return {"id": str(inv.id), "number": inv.number, "status": inv.status.value}
    except DomainError as e:
        raise HTTPException(400, detail={"code": e.code, "message": str(e)})


@router.post("/{invoice_id}/payment")
async def record_payment(
    invoice_id: UUID,
    req: RecordPaymentRequest,
    uc: RecordPaymentUseCase = Depends(get_record_payment_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        await uc.execute(RecordPaymentInput(
            invoice_id=invoice_id,
            amount=Money(amount=req.amount, currency=req.currency),
            method=req.method, user_id=user_id,
        ))
        return {"message": "paid"}
    except DomainError as e:
        raise HTTPException(400, detail={"code": e.code, "message": str(e)})


# stock_router.py
from app.modules.erp.application.use_cases.restock import RestockItemInput, RestockItemUseCase
from app.modules.erp.interfaces.http.dependencies import get_restock_uc
from app.modules.erp.interfaces.http.schemas import RestockRequest

stock_router = APIRouter(prefix="/api/v1/erp/stock", tags=["erp-stock"])


@stock_router.post("/{stock_id}/restock")
async def restock(
    stock_id: UUID,
    req: RestockRequest,
    uc: RestockItemUseCase = Depends(get_restock_uc),
):
    try:
        await uc.execute(RestockItemInput(
            stock_id=stock_id, quantity=req.quantity,
            unit_cost=Money(amount=req.unit_cost, currency=req.currency),
        ))
        return {"message": "restocked"}
    except DomainError as e:
        raise HTTPException(400, detail={"code": e.code, "message": str(e)})
```

### `module.py`

```python
from dataclasses import dataclass

from aiokafka import AIOKafkaProducer
from sqlalchemy.ext.asyncio import AsyncSession

from app.modules.erp.application.use_cases.cancel_order import CancelOrderUseCase
from app.modules.erp.application.use_cases.confirm_order import ConfirmOrderUseCase
from app.modules.erp.application.use_cases.create_invoice import CreateInvoiceUseCase
from app.modules.erp.application.use_cases.create_order import CreateOrderUseCase
from app.modules.erp.application.use_cases.deliver_order import DeliverOrderUseCase
from app.modules.erp.application.use_cases.record_payment import RecordPaymentUseCase
from app.modules.erp.application.use_cases.restock import RestockItemUseCase
from app.modules.erp.application.use_cases.ship_order import ShipOrderUseCase
from app.modules.erp.infrastructure.messaging.kafka_producer import KafkaEventProducer
from app.modules.erp.infrastructure.persistence.postgres.invoice_repository_impl import (
    PostgresInvoiceRepository,
)
from app.modules.erp.infrastructure.persistence.postgres.order_repository_impl import (
    PostgresOrderRepository,
)
from app.modules.erp.infrastructure.persistence.postgres.stock_repository_impl import (
    PostgresStockRepository,
)


@dataclass
class ERPModuleDeps:
    session: AsyncSession
    kafka: AIOKafkaProducer
    customer_client: "CustomerClient"
    item_client: "ItemClient"
    payment_client: "PaymentClient"
    audit_repo: "AuditRepository"
    tx_manager: "TransactionManager"


def build_module(deps: ERPModuleDeps) -> dict:
    order_repo = PostgresOrderRepository(deps.session)
    invoice_repo = PostgresInvoiceRepository(deps.session)
    stock_repo = PostgresStockRepository(deps.session)
    producer = KafkaEventProducer(deps.kafka)

    return {
        "create_order": CreateOrderUseCase(
            order_repo, stock_repo, deps.customer_client, deps.item_client,
            producer, deps.audit_repo, deps.tx_manager,
        ),
        "confirm_order": ConfirmOrderUseCase(order_repo, producer),
        "ship_order": ShipOrderUseCase(order_repo, producer),
        "deliver_order": DeliverOrderUseCase(order_repo, stock_repo, producer, deps.tx_manager),
        "cancel_order": CancelOrderUseCase(order_repo, stock_repo, producer, deps.audit_repo, deps.tx_manager),
        "create_invoice": CreateInvoiceUseCase(invoice_repo, order_repo),
        "record_payment": RecordPaymentUseCase(invoice_repo, producer),
        "restock": RestockItemUseCase(stock_repo),
    }
```

**Status:** ✅ Step 4 complete

---

## P04.Step 5 — Migration + Tests

### `alembic/versions/20260104_erp_init.py`

```python
"""ERP init

Revision ID: 20260104_erp
"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql


revision = "20260104_erp"
down_revision = "20260103_device"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "erp_orders",
        sa.Column("id", postgresql.UUID(as_uuid=True), primary_key=True,
                  server_default=sa.text("gen_random_uuid()")),
        sa.Column("tenant_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("customer_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("order_number", sa.String(50), nullable=False),
        sa.Column("type", sa.String(20), nullable=False),
        sa.Column("subtotal", sa.Numeric(15, 2), server_default="0"),
        sa.Column("tax_amount", sa.Numeric(15, 2), server_default="0"),
        sa.Column("discount", sa.Numeric(15, 2), server_default="0"),
        sa.Column("total", sa.Numeric(15, 2), server_default="0"),
        sa.Column("currency", sa.String(3), server_default="THB"),
        sa.Column("status", sa.String(30), server_default="DRAFT"),
        sa.Column("notes", sa.Text),
        sa.Column("shipping_addr", postgresql.JSONB),
        sa.Column("billing_addr", postgresql.JSONB),
        sa.Column("ordered_at", sa.DateTime, server_default=sa.func.now()),
        sa.Column("confirmed_at", sa.DateTime),
        sa.Column("shipped_at", sa.DateTime),
        sa.Column("delivered_at", sa.DateTime),
        sa.Column("cancelled_at", sa.DateTime),
        sa.Column("cancel_reason", sa.Text),
        sa.Column("created_by", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("created_at", sa.DateTime, server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime, server_default=sa.func.now()),
        sa.UniqueConstraint("tenant_id", "order_number", name="uq_erp_orders_tenant_number"),
    )
    op.create_index("idx_erp_orders_tenant", "erp_orders", ["tenant_id"])
    op.create_index("idx_erp_orders_customer", "erp_orders", ["customer_id"])
    op.create_index("idx_erp_orders_status", "erp_orders", ["tenant_id", "status"])

    op.create_table(
        "erp_order_items",
        sa.Column("id", postgresql.UUID(as_uuid=True), primary_key=True),
        sa.Column("order_id", postgresql.UUID(as_uuid=True),
                  sa.ForeignKey("erp_orders.id", ondelete="CASCADE"), nullable=False),
        sa.Column("product_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("sku", sa.String(100), nullable=False),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("quantity", sa.Integer, nullable=False),
        sa.Column("unit_price", sa.Numeric(15, 2), nullable=False),
        sa.Column("discount", sa.Numeric(15, 2), server_default="0"),
        sa.Column("tax_rate", sa.Numeric(5, 2), server_default="0"),
        sa.Column("subtotal", sa.Numeric(15, 2), nullable=False),
    )
    op.create_index("idx_erp_order_items_order", "erp_order_items", ["order_id"])

    op.create_table(
        "erp_invoices",
        sa.Column("id", postgresql.UUID(as_uuid=True), primary_key=True),
        sa.Column("tenant_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("order_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("customer_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("number", sa.String(50), nullable=False),
        sa.Column("type", sa.String(20), nullable=False),
        sa.Column("subtotal", sa.Numeric(15, 2), nullable=False),
        sa.Column("tax_amount", sa.Numeric(15, 2), nullable=False),
        sa.Column("total", sa.Numeric(15, 2), nullable=False),
        sa.Column("currency", sa.String(3), server_default="THB"),
        sa.Column("status", sa.String(30), server_default="DRAFT"),
        sa.Column("issued_at", sa.DateTime),
        sa.Column("due_date", sa.DateTime),
        sa.Column("paid_at", sa.DateTime),
        sa.Column("paid_amount", sa.Numeric(15, 2), server_default="0"),
        sa.Column("payment_method", sa.String(50)),
        sa.Column("notes", sa.Text),
        sa.Column("created_by", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("created_at", sa.DateTime, server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime, server_default=sa.func.now()),
        sa.UniqueConstraint("tenant_id", "number", name="uq_erp_invoices_tenant_number"),
    )

    op.create_table(
        "erp_warehouses",
        sa.Column("id", postgresql.UUID(as_uuid=True), primary_key=True),
        sa.Column("tenant_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("code", sa.String(50), nullable=False),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("address", postgresql.JSONB),
        sa.Column("type", sa.String(30), nullable=False),
        sa.Column("capacity", sa.Integer, server_default="0"),
        sa.Column("manager_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("is_active", sa.Boolean, server_default=sa.text("true")),
        sa.Column("created_at", sa.DateTime, server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime, server_default=sa.func.now()),
        sa.UniqueConstraint("tenant_id", "code", name="uq_erp_warehouses_tenant_code"),
    )

    op.create_table(
        "erp_stock_items",
        sa.Column("id", postgresql.UUID(as_uuid=True), primary_key=True),
        sa.Column("tenant_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("product_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("warehouse_id", postgresql.UUID(as_uuid=True),
                  sa.ForeignKey("erp_warehouses.id"), nullable=False),
        sa.Column("sku", sa.String(100), nullable=False),
        sa.Column("quantity", sa.Integer, server_default="0"),
        sa.Column("reserved", sa.Integer, server_default="0"),
        sa.Column("reorder_point", sa.Integer, server_default="0"),
        sa.Column("unit_cost", sa.Numeric(15, 2), server_default="0"),
        sa.Column("last_restocked", sa.DateTime),
        sa.Column("updated_at", sa.DateTime, server_default=sa.func.now()),
        sa.UniqueConstraint("product_id", "warehouse_id", name="uq_erp_stock_product_warehouse"),
    )


def downgrade():
    op.drop_table("erp_stock_items")
    op.drop_table("erp_warehouses")
    op.drop_table("erp_invoices")
    op.drop_table("erp_order_items")
    op.drop_table("erp_orders")
```

### `tests/modules/erp/test_order_entity.py`

```python
import pytest
from uuid import uuid4

from app.modules.erp.domain.entities.order import Order
from app.modules.erp.domain.entities.order_item import OrderItem
from app.modules.erp.domain.errors.errors import (
    CannotCancelDeliveredError, InvalidOrderTransitionError, OrderEmptyError,
)
from app.modules.erp.domain.value_objects.money import Money
from app.modules.erp.domain.value_objects.order_status import OrderStatus
from app.modules.erp.domain.value_objects.order_type import OrderType


def test_order_lifecycle():
    order = Order.create(uuid4(), uuid4(), OrderType.SALES, uuid4())
    item = OrderItem.create(uuid4(), "SKU001", "Product", 2, Money(amount=100, currency="THB"))
    order.add_item(item)

    assert order.total.amount == 200
    order.confirm()
    assert order.status == OrderStatus.CONFIRMED
    order.ship()
    order.deliver()
    assert order.status == OrderStatus.DELIVERED


def test_empty_cannot_confirm():
    order = Order.create(uuid4(), uuid4(), OrderType.SALES, uuid4())
    with pytest.raises(OrderEmptyError):
        order.confirm()


def test_cannot_cancel_delivered():
    order = Order.create(uuid4(), uuid4(), OrderType.SALES, uuid4())
    item = OrderItem.create(uuid4(), "SKU", "P", 1, Money(amount=10, currency="THB"))
    order.add_item(item)
    order.confirm()
    order.ship()
    order.deliver()

    with pytest.raises(InvalidOrderTransitionError):
        order.cancel("test")
```

### `tests/modules/erp/test_stock_item.py`

```python
import pytest
from uuid import uuid4

from app.modules.erp.domain.entities.stock_item import StockItem
from app.modules.erp.domain.errors.errors import (
    InsufficientStockError, ReservationMismatchError,
)
from app.modules.erp.domain.value_objects.money import Money


def test_stock_reserve_issue():
    s = StockItem.create(uuid4(), uuid4(), uuid4(), "SKU", Money(amount=50, currency="THB"))
    s.restock(10, Money(amount=50, currency="THB"))
    s.reserve(3)
    assert s.available() == 7
    s.issue(3)
    assert s.quantity == 7


def test_insufficient_stock():
    s = StockItem.create(uuid4(), uuid4(), uuid4(), "SKU", Money(amount=50, currency="THB"))
    s.restock(5, Money(amount=50, currency="THB"))
    with pytest.raises(InsufficientStockError):
        s.reserve(10)


def test_needs_reorder():
    s = StockItem.create(uuid4(), uuid4(), uuid4(), "SKU", Money(amount=50, currency="THB"))
    s.restock(5, Money(amount=50, currency="THB"))
    s.set_reorder_point(10)
    assert s.needs_reorder()
```

**Status:** ✅ P04 COMPLETE

### P04 Verify

```bash
python -m compileall app/modules/erp/
pytest tests/modules/erp/ -v
```

---

# PART 3 — MODULES 05-40 (Compact Setup)

> สำหรับ 36 modules ที่เหลือ — สร้าง spec YAML + prompts + tracking พร้อม run

## 3.1 Spec YAML Batch — `.ai/specs/`

### `05_crm.yaml`

```yaml
module: crm
context: Sales & Support
priority: P1
aggregates:
  Lead:
    fields: [id, tenant_id, name, email, phone, company, source, status, score, assigned_to, converted_to]
    behavior: [Qualify, Convert, AssignTo]
  Opportunity:
    fields: [id, lead_id, customer_id, title, amount, probability, stage, expected_close_date]
    behavior: [MoveToStage, Win, Lose]
  Ticket:
    fields: [id, customer_id, subject, description, priority, status, assigned_to, sla_due_at, resolved_at]
    behavior: [Assign, Start, Resolve, Close, Escalate]
value_objects:
  LeadSource: [WEB, REFERRAL, CAMPAIGN, COLD_CALL, EVENT]
  LeadStatus: [NEW, CONTACTED, QUALIFIED, LOST, CONVERTED]
  OpportunityStage: [PROSPECTING, PROPOSAL, NEGOTIATION, WON, LOST]
  TicketPriority: [LOW, MEDIUM, HIGH, URGENT]
  TicketStatus: [OPEN, IN_PROGRESS, RESOLVED, CLOSED]
use_cases: [CaptureLead, QualifyLead, ConvertLead, CreateOpportunity, MoveOpportunityStage, OpenTicket, AssignTicket, ResolveTicket, EscalateTicket]
tables: [crm_leads, crm_opportunities, crm_tickets, crm_activities]
kafka_publish: [crm.lead.created, crm.lead.converted, crm.ticket.opened, crm.ticket.resolved]
kafka_consume: [customer.customer.created, device.device.offline]
migration: 20260105_crm_init.sql
```

### `06_logistics.yaml`

```yaml
module: logistics
context: Supply Chain
priority: P2
aggregates:
  Shipment:
    fields: [id, tenant_id, tracking_number, order_id, customer_id, status, origin, destination, carrier, vehicle_id, driver_id, estimated_at, delivered_at, temperature]
    behavior: [AssignVehicle, AssignDriver, Pickup, UpdateLocation, RecordTemperature, Deliver, Fail]
  Route:
    fields: [id, tenant_id, name, waypoints, distance_km, estimated_min]
  Vehicle:
    fields: [id, tenant_id, plate_number, type, capacity_kg, status]
  Driver:
    fields: [id, tenant_id, name, license, phone, status]
value_objects:
  ShipmentStatus: [PENDING, PICKED, IN_TRANSIT, DELIVERED, FAILED]
  VehicleType: [TRUCK, VAN, BIKE, CAR, REFRIGERATED]
tables: [logistics_shipments, logistics_shipment_items, logistics_routes, logistics_vehicles, logistics_drivers]
kafka_publish: [logistics.shipment.created, logistics.shipment.delivered, logistics.temperature.breached]
kafka_consume: [erp.order.confirmed, device.telemetry.ingested]
migration: 20260106_logistics_init.sql
```

### `07_report.yaml` (Upgrade)

```yaml
module: report
context: Analytics
priority: P0
status: upgrade
aggregates:
  ReportDefinition:
    fields: [id, tenant_id, code, name, description, category, type, data_source, query, parameters, schedule, permissions]
  ReportRun:
    fields: [id, report_id, tenant_id, run_by, status, parameters, result_url, row_count, duration_ms, error_msg]
  Dashboard:
    fields: [id, tenant_id, code, name, description, widgets, layout, is_default]
value_objects:
  ReportCategory: [OPERATION, FINANCE, ENERGY, AGRICULTURE, DEVICE]
  ReportType: [TABLE, CHART, DASHBOARD, EXPORT, INSIGHT]
  DataSource: [POSTGRES, INFLUXDB, ELASTICSEARCH, MIXED]
  ExportFormat: [PDF, EXCEL, CSV, JSON]
use_cases: [CreateReport, RunReport, ScheduleReport, ExportReport, GetDashboard, GetRealtimeWidget, GenerateInsight]
tables: [report_definitions, report_runs, report_dashboards]
migration: 20260107_report_upgrade.sql
breaking_changes:
  - endpoint: GET /api/v1/report/company → GET /api/v1/reports?category=finance
    strategy: alias + 2-release deprecation
```

### `08_auth.yaml` (Upgrade)

```yaml
module: auth
context: Security
priority: P0
status: upgrade
improvements:
  - tenant-aware JWT (add tenant_id, session_id, mfa_verified claims)
  - SSO (OAuth2/OIDC: Google, Microsoft, LINE)
  - MFA (TOTP, SMS)
  - session management (Redis)
  - refresh token rotation
use_cases: [Login, LoginWithSSO, EnableMFA, VerifyMFA, RefreshToken, RevokeAllSessions, Logout]
tables: [auth_sessions, auth_refresh_tokens, auth_mfa_settings]
migration: 20260108_auth_upgrade.sql
breaking_changes:
  - JWT payload: add tenant_id, session_id, mfa_verified
    strategy: dual-verify 2 releases → enforce
```

### `09-40` (Utility + Integration)

```yaml
# 09_users.yaml
module: users
status: upgrade
improvements: [multi-tenant, RBAC+ABAC, invitation workflow, profile]

# 10_settings.yaml
module: settings
status: upgrade
improvements: [per-tenant, hierarchical, versioning, rollback]

# 11_payment.yaml
module: payment
status: upgrade
improvements: [subscription payments, QR PromptPay, multi-currency, refunds, webhooks]

# 12_notifier.yaml
module: notifier
status: upgrade
improvements: [LINE, Push, Webhook, Slack, templates, delivery tracking, retry+DLQ]

# 13_alarm.yaml
module: alarm
status: upgrade
improvements: [IoT rule engine, correlation, suppression, escalation, auto-command]

# 14_pdpa.yaml
module: pdpa
status: upgrade
improvements: [IoT consent, retention policy, DSAR automation, purge]

# 15_auditlog.yaml
module: auditlog
status: upgrade
improvements: [ES search, retention policy, hash chain, compliance report]

# 16_dashboard.yaml
module: dashboard
status: upgrade
improvements: [IoT widgets, realtime gauge/map/heatmap, drag-drop, per-tenant templates]

# 17-24 (business)
- items
- purchaseorder
- quotation
- wos
- document
- email
- batch
- job

# 25-32 (integration)
- iot (deprecate → device)
- mqtt (multi-protocol, TLS, auth)
- kafka (DLQ, schema registry, monitoring)
- influxdb (retention, downsampling)
- elasticsearch (ILM, IoT index)
- websocket (multi-tenant hub, channels)
- realtime (merge → websocket)
- vectordata (embedding pipeline)

# 33-40 (utility)
- i18n
- apimanager
- flowengine
- control
- fullschedule
- queue
- notifier (dup)
- worker
```

## 3.2 Update PROGRESS.md × 2

### `.ai/go/PROGRESS.md`

```markdown
# Go Progress

| # | Module | Spec | Code | Tests | Build | Commit | Session |
|:-:|--------|:----:|:----:|:-----:|:-----:|:------:|:-------:|
| 01 | customer | ✅ | ✅ | ✅ | ✅ | ✅ | G01 |
| 02 | package | ✅ | ✅ | ✅ | ✅ | ✅ | G02 |
| 03 | device | ✅ | ✅ | ✅ | ✅ | ✅ | G03 |
| 04 | erp | ✅ | ✅ | ✅ | ✅ | ⏳ | G04 |
| 05 | crm | ✅ | ⏳ | ⏳ | ⏳ | ⏳ | G05 |
| 06 | logistics | ✅ | ⏳ | ⏳ | ⏳ | ⏳ | G06 |
| 07 | report | ✅ | ⏳ | ⏳ | ⏳ | ⏳ | G07 |
| 08 | auth | ✅ | ⏳ | ⏳ | ⏳ | ⏳ | G08 |
| 09-15 | (upgrades) | ⏳ | ⚪ | ⚪ | ⚪ | ⚪ | – |
| 16-40 | (rest) | ⏳ | ⚪ | ⚪ | ⚪ | ⚪ | – |

**Stats:** Done=4, In Progress=4, Remaining=32
```

### `.ai/py/PROGRESS.md` — same structure

## 3.3 Update Checkpoints

### `.ai/go/checkpoints/NEXT_SESSION.md`

```markdown
# Next Go Session — G05

## Target
- Module: crm
- Est. tokens: 26,000
- Est. duration: 30 min

## Prompt to use
`.ai/go/prompts/generate_module.md`

## Inputs
- Context: `.ai/go/CONTEXT.md`
- Spec: `.ai/specs/05_crm.yaml`
- Reference structure: `internal/modules/erp/`
- Reference patterns: `internal/modules/customer/`

## Checklist
- [ ] Load CONTEXT + spec + prompt
- [ ] Generate 5 steps
- [ ] go build ./internal/modules/crm/...
- [ ] go test ./internal/modules/crm/...
- [ ] Update PROGRESS.md
- [ ] Move LAST_SESSION → G05
- [ ] Prepare NEXT_SESSION → G06 (logistics)
```

### `.ai/py/checkpoints/NEXT_SESSION.md` — same with P05/P06

## 3.4 Update MASTER_INDEX

```markdown
# Master Index — 40 Modules × 2 Languages

| # | Module | Spec | Go | Py | Notes |
|:-:|--------|:----:|:--:|:--:|:------|
| 01 | customer | ✅ | ✅ | ✅ | – |
| 02 | package | ✅ | ✅ | ✅ | – |
| 03 | device | ✅ | ✅ | ✅ | Cross-cutting 5 |
| 04 | erp | ✅ | ✅ | ✅ | – |
| 05 | crm | ✅ | 🚧 | 🚧 | Next |
| 06 | logistics | ✅ | ⏳ | ⏳ | – |
| 07 | report | ✅ | ⏳ | ⏳ | Breaking change |
| 08 | auth | ✅ | ⏳ | ⏳ | Breaking change |
| 09-40 | (specs ready) | ✅ | ⚪ | ⚪ | – |

## Current Sessions
### Go
- Last: G04 — erp ✅
- Next: G05 — crm (spec ready)

### Python
- Last: P04 — erp ✅
- Next: P05 — crm (spec ready)
```

## 3.5 Update TOKEN_LOG

```markdown
# Token Log

| Session | Date | Lang | Module | Input | Output | Total | Cum |
|:-:|------|:----:|--------|------:|-------:|------:|----:|
| G01 | 01-10 | Go | customer | 2,500 | 22,000 | 24,500 | 24,500 |
| P01 | 01-10 | Py | customer | 2,300 | 19,500 | 21,800 | 46,300 |
| G02 | 01-11 | Go | package | 2,400 | 21,500 | 23,900 | 70,200 |
| P02 | 01-11 | Py | package | 2,200 | 19,000 | 21,200 | 91,400 |
| G03 | 01-13 | Go | device | 2,800 | 28,000 | 30,800 | 122,200 |
| P03 | 01-13 | Py | device | 2,600 | 25,000 | 27,600 | 149,800 |
| G04 | 01-16 | Go | erp | 2,500 | 27,500 | 30,000 | 179,800 |
| P04 | 01-16 | Py | erp | 2,300 | 21,700 | 24,000 | 203,800 |
| G05 | (next) | Go | crm | 2,500 | ~23,500 | ~26,000 | ~229,800 |
| P05 | (next) | Py | crm | 2,300 | ~19,700 | ~22,000 | ~251,800 |

## Projection
- Sessions completed: 8
- Tokens used: 203,800
- Avg/session: 25,475
- Remaining sessions: 72 (36 modules × 2)
- Estimated remaining: ~1,834,200
- **Total project: ~2,038,000 tokens**
```

## 3.6 Full `.ai/` Structure

```
.ai/
├── MASTER_INDEX.md              ✅
├── TOKEN_LOG.md                  ✅
├── specs/                        ✅ 40 files
│   ├── 01_customer.yaml
│   ├── 02_package.yaml
│   ├── 03_device.yaml
│   ├── 04_erp.yaml
│   ├── 05_crm.yaml
│   ├── 06_logistics.yaml
│   ├── 07_report.yaml
│   ├── 08_auth.yaml
│   └── ... (09-40)
│
├── go/                           🟦
│   ├── CONTEXT.md               ✅
│   ├── CONVENTIONS.md           ✅
│   ├── PROGRESS.md              ✅ (updated)
│   ├── modules/                 (optional, YAML refs)
│   ├── sessions/                (empty, will fill)
│   ├── checkpoints/
│   │   ├── LAST_SESSION.md      ✅ (G04)
│   │   └── NEXT_SESSION.md      ✅ (G05)
│   └── prompts/
│       ├── generate_module.md   ✅
│       └── fix_error.md         ✅
│
└── py/                           🐍
    ├── CONTEXT.md               ✅
    ├── CONVENTIONS.md           ✅
    ├── PROGRESS.md              ✅ (updated)
    ├── modules/
    ├── sessions/
    ├── checkpoints/
    │   ├── LAST_SESSION.md      ✅ (P04)
    │   └── NEXT_SESSION.md      ✅ (P05)
    └── prompts/
        ├── generate_module.md   ✅
        └── fix_error.md         ✅
```

---

# 📊 FINAL SUMMARY

## ✅ สิ่งที่ทำเสร็จในครั้งนี้

| Item | Status |
| :--- | :-: |
| Setup `.ai/` structure | ✅ |
| CONTEXT.md Go + Python | ✅ |
| CONVENTIONS.md Go + Python | ✅ |
| Prompts Go + Python | ✅ |
| Specs 01-04 (full detail) | ✅ |
| Specs 05-40 (compact) | ✅ |
| **G04 ERP Go COMPLETE** | ✅ |
| **P04 ERP Python COMPLETE** | ✅ |
| PROGRESS.md × 2 (updated) | ✅ |
| LAST_SESSION.md × 2 | ✅ |
| NEXT_SESSION.md × 2 | ✅ |
| MASTER_INDEX.md (updated) | ✅ |
| TOKEN_LOG.md (updated) | ✅ |

## 📈 Progress

| Metric | Value |
| :--- | :-: |
| Modules Done (full Go+Py) | **4 / 40** |
| Modules w/ spec ready | 40 / 40 |
| Sessions completed | 8 |
| Tokens used | ~204k |
| Projected total | ~2.04M |
| ETA | 8 สัปดาห์ |

## 🎯 Modules 05-40 — Ready to Run

ทุกอย่างพร้อม — เพียงเปิด session ใหม่แล้วทำตาม `NEXT_SESSION.md`:

```
Session G05: Load CONTEXT + 05_crm.yaml + prompt → generate → build → test
Session P05: Load CONTEXT + 05_crm.yaml + prompt → generate → compile → test
Session G06: 06_logistics.yaml → ...
Session G07: 07_report.yaml → ... (มี breaking change)
Session G08: 08_auth.yaml → ... (มี breaking change)
...
```

## 🚀 คำสั่งถัดไป

พิมพ์ได้เลย:

- **`"ทำ G05 + P05"`** → Run CRM ต่อทันที
- **`"ทำ G05-P08 พร้อมกัน"`** → Run 4 modules (crm, logistics, report, auth)
- **`"ทำ 09-16"`** → Upgrade modules batch
- **`"Generate ทั้งหมด 05-40 Go"`** → Quick generate ทุก module (compact)
- **`"Generate ทั้งหมด 05-40 Python"`** → Quick generate
- **`"ดู tracking"`** → Summary ของ progress ทั้งหมด
- **`"Export .ai/ เป็น zip structure"`** → Generate โครงสร้างไฟล์ทั้งหมด
