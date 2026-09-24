<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## การแยก Business Logic: Custom Hooks vs Services

การแยก business logic ออกจาก UI components เป็นหลักการสำคัญในการพัฒนา React applications ที่มีความ scalable และ maintainable  มีสองแนวทางหลัก: **Custom Hooks** สำหรับ logic ที่เกี่ยวข้องกับ React lifecycle และ **Services** สำหรับ pure business logic และ API calls[^1][^2][^3][^4]

## ความแตกต่างระหว่าง Custom Hooks และ Services

| Aspect | Custom Hooks | Services |
| :-- | :-- | :-- |
| **Purpose** | React-specific logic, state management | Pure business logic, API calls |
| **Can use React features** | ✅ useState, useEffect, useContext | ❌ Plain JavaScript functions |
| **Where to use** | Inside React components only | Anywhere (components, hooks, tests) |
| **Rerenders** | Can trigger component rerenders | No direct rerender capability |
| **Testability** | Requires React Testing Library | Easy to test with Jest alone |
| **Reusability** | Across React components | Across entire application |

[^2][^5]

## Architecture Pattern: 3-Layer Structure

```
┌─────────────────────────────────────┐
│       UI Layer (Components)         │  ← Presentation only
├─────────────────────────────────────┤
│    Business Logic (Custom Hooks)    │  ← React-aware logic
├─────────────────────────────────────┤
│    Data Layer (Services)            │  ← API calls, pure functions
└─────────────────────────────────────┘
```


## ตัวอย่างที่ 1: Booking Management

### Layer 1: Services (Data Layer)

**Booking Service** (`src/features/bookings/services/bookingService.js`)

```javascript
import api from '../../../services/api';

/**
 * Pure functions for API calls
 * No React dependencies, no state management
 * Easy to test, reusable anywhere
 */
const bookingService = {
  // Get all bookings
  getAllBookings: async () => {
    const response = await api.get('/bookings');
    return response.data;
  },

  // Get single booking by ID
  getBookingById: async (bookingId) => {
    const response = await api.get(`/bookings/${bookingId}`);
    return response.data;
  },

  // Create new booking
  createBooking: async (bookingData) => {
    const response = await api.post('/bookings', bookingData);
    return response.data;
  },

  // Update booking
  updateBooking: async (bookingId, updates) => {
    const response = await api.put(`/bookings/${bookingId}`, updates);
    return response.data;
  },

  // Cancel booking
  cancelBooking: async (bookingId, reason) => {
    const response = await api.post(`/bookings/${bookingId}/cancel`, {
      reason,
      cancelledAt: new Date().toISOString(),
    });
    return response.data;
  },

  // Check availability
  checkAvailability: async (date, serviceType) => {
    const response = await api.get('/bookings/availability', {
      params: { date, serviceType },
    });
    return response.data;
  },

  // Calculate price
  calculatePrice: (serviceType, vehicleType, additionalServices = []) => {
    const basePrices = {
      'oil-change': 800,
      'general-maintenance': 1500,
      'tire-change': 2000,
      'full-service': 3500,
    };

    const vehicleMultipliers = {
      'sedan': 1.0,
      'suv': 1.3,
      'truck': 1.5,
    };

    let total = basePrices[serviceType] || 0;
    total *= vehicleMultipliers[vehicleType] || 1.0;

    additionalServices.forEach((service) => {
      total += service.price;
    });

    return {
      basePrice: basePrices[serviceType],
      multiplier: vehicleMultipliers[vehicleType],
      additionalServicesTotal: additionalServices.reduce(
        (sum, s) => sum + s.price,
        0
      ),
      total,
    };
  },

  // Validate booking data
  validateBookingData: (data) => {
    const errors = {};

    if (!data.customerId) {
      errors.customerId = 'Customer ID is required';
    }

    if (!data.vehicleId) {
      errors.vehicleId = 'Vehicle ID is required';
    }

    if (!data.bookingDate) {
      errors.bookingDate = 'Booking date is required';
    } else {
      const bookingDate = new Date(data.bookingDate);
      const today = new Date();
      today.setHours(0, 0, 0, 0);

      if (bookingDate < today) {
        errors.bookingDate = 'Cannot book in the past';
      }
    }

    if (!data.serviceType) {
      errors.serviceType = 'Service type is required';
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors,
    };
  },
};

export default bookingService;
```


### Layer 2: Custom Hooks (Business Logic Layer)

**useBookings Hook** (`src/features/bookings/hooks/useBookings.js`)

```javascript
import { useState, useEffect, useCallback } from 'react';
import bookingService from '../services/bookingService';
import { useAuth } from '../../auth/context/AuthContext';

/**
 * Custom hook for booking management
 * Handles React state, side effects, and lifecycle
 * Consumes bookingService for data operations
 */
const useBookings = () => {
  const { user } = useAuth();
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Fetch all bookings
  const fetchBookings = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await bookingService.getAllBookings();
      setBookings(data);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch bookings');
      console.error('Error fetching bookings:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Load bookings on mount
  useEffect(() => {
    if (user) {
      fetchBookings();
    }
  }, [user, fetchBookings]);

  // Create booking with validation
  const createBooking = useCallback(async (bookingData) => {
    setLoading(true);
    setError(null);

    // Validate before sending
    const validation = bookingService.validateBookingData(bookingData);
    if (!validation.isValid) {
      setError('Please fix validation errors');
      setLoading(false);
      return { success: false, errors: validation.errors };
    }

    try {
      const newBooking = await bookingService.createBooking(bookingData);
      setBookings((prev) => [newBooking, ...prev]);
      return { success: true, data: newBooking };
    } catch (err) {
      const errorMessage = err.response?.data?.message || 'Failed to create booking';
      setError(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  }, []);

  // Cancel booking
  const cancelBooking = useCallback(async (bookingId, reason) => {
    setLoading(true);
    setError(null);

    try {
      await bookingService.cancelBooking(bookingId, reason);
      
      // Update local state
      setBookings((prev) =>
        prev.map((booking) =>
          booking.id === bookingId
            ? { ...booking, status: 'cancelled', cancelReason: reason }
            : booking
        )
      );

      return { success: true };
    } catch (err) {
      const errorMessage = err.response?.data?.message || 'Failed to cancel booking';
      setError(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  }, []);

  // Filter bookings by status
  const getBookingsByStatus = useCallback(
    (status) => {
      return bookings.filter((booking) => booking.status === status);
    },
    [bookings]
  );

  // Get upcoming bookings (within next 7 days)
  const getUpcomingBookings = useCallback(() => {
    const now = new Date();
    const sevenDaysLater = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);

    return bookings.filter((booking) => {
      const bookingDate = new Date(booking.bookingDate);
      return (
        bookingDate >= now &&
        bookingDate <= sevenDaysLater &&
        booking.status !== 'cancelled'
      );
    });
  }, [bookings]);

  return {
    bookings,
    loading,
    error,
    fetchBookings,
    createBooking,
    cancelBooking,
    getBookingsByStatus,
    getUpcomingBookings,
  };
};

export default useBookings;
```

**useBookingForm Hook** (`src/features/bookings/hooks/useBookingForm.js`)

```javascript
import { useState, useEffect } from 'react';
import bookingService from '../services/bookingService';
import { useAuth } from '../../auth/context/AuthContext';

/**
 * Custom hook for booking form logic
 * Handles form state, validation, and availability checking
 */
const useBookingForm = (initialData = {}) => {
  const { user } = useAuth();
  
  const [formData, setFormData] = useState({
    customerId: user?.id || '',
    vehicleId: '',
    bookingDate: '',
    serviceType: '',
    vehicleType: '',
    additionalServices: [],
    notes: '',
    ...initialData,
  });

  const [errors, setErrors] = useState({});
  const [priceBreakdown, setPriceBreakdown] = useState(null);
  const [availability, setAvailability] = useState(null);
  const [checkingAvailability, setCheckingAvailability] = useState(false);

  // Calculate price when relevant fields change
  useEffect(() => {
    if (formData.serviceType && formData.vehicleType) {
      const breakdown = bookingService.calculatePrice(
        formData.serviceType,
        formData.vehicleType,
        formData.additionalServices
      );
      setPriceBreakdown(breakdown);
    }
  }, [formData.serviceType, formData.vehicleType, formData.additionalServices]);

  // Check availability when date and service type change
  useEffect(() => {
    const checkAvailability = async () => {
      if (formData.bookingDate && formData.serviceType) {
        setCheckingAvailability(true);
        try {
          const result = await bookingService.checkAvailability(
            formData.bookingDate,
            formData.serviceType
          );
          setAvailability(result);
        } catch (error) {
          console.error('Failed to check availability:', error);
          setAvailability(null);
        } finally {
          setCheckingAvailability(false);
        }
      }
    };

    const debounceTimer = setTimeout(checkAvailability, 500);
    return () => clearTimeout(debounceTimer);
  }, [formData.bookingDate, formData.serviceType]);

  // Update form field
  const updateField = (field, value) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));

    // Clear error for this field
    if (errors[field]) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  // Add additional service
  const addAdditionalService = (service) => {
    setFormData((prev) => ({
      ...prev,
      additionalServices: [...prev.additionalServices, service],
    }));
  };

  // Remove additional service
  const removeAdditionalService = (serviceId) => {
    setFormData((prev) => ({
      ...prev,
      additionalServices: prev.additionalServices.filter(
        (s) => s.id !== serviceId
      ),
    }));
  };

  // Validate form
  const validate = () => {
    const validation = bookingService.validateBookingData(formData);
    setErrors(validation.errors);
    return validation.isValid;
  };

  // Reset form
  const reset = () => {
    setFormData({
      customerId: user?.id || '',
      vehicleId: '',
      bookingDate: '',
      serviceType: '',
      vehicleType: '',
      additionalServices: [],
      notes: '',
    });
    setErrors({});
    setPriceBreakdown(null);
    setAvailability(null);
  };

  return {
    formData,
    errors,
    priceBreakdown,
    availability,
    checkingAvailability,
    updateField,
    addAdditionalService,
    removeAdditionalService,
    validate,
    reset,
  };
};

export default useBookingForm;
```


### Layer 3: UI Components

**BookingForm Component** (`src/features/bookings/components/BookingForm.jsx`)

```jsx
import { useState } from 'react';
import useBookings from '../hooks/useBookings';
import useBookingForm from '../hooks/useBookingForm';
import Button from '../../../components/common/Button/Button';
import Input from '../../../components/common/Input/Input';
import Select from '../../../components/common/Select/Select';
import './BookingForm.css';

/**
 * Pure UI component
 * No business logic, just presentation and event handling
 * Consumes custom hooks for all logic
 */
const BookingForm = ({ onSuccess }) => {
  const { createBooking, loading } = useBookings();
  const {
    formData,
    errors,
    priceBreakdown,
    availability,
    checkingAvailability,
    updateField,
    addAdditionalService,
    removeAdditionalService,
    validate,
    reset,
  } = useBookingForm();

  const [submitError, setSubmitError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitError('');

    if (!validate()) {
      setSubmitError('กรุณากรอกข้อมูลให้ครบถ้วน');
      return;
    }

    if (!availability?.isAvailable) {
      setSubmitError('ช่วงเวลาที่เลือกไม่ว่าง กรุณาเลือกวันอื่น');
      return;
    }

    const result = await createBooking(formData);

    if (result.success) {
      reset();
      onSuccess?.(result.data);
    } else {
      setSubmitError(result.error || 'ไม่สามารถสร้างการจองได้');
    }
  };

  return (
    <form className="booking-form" onSubmit={handleSubmit}>
      <h2>จองบริการ</h2>

      <Select
        label="ประเภทบริการ"
        value={formData.serviceType}
        onChange={(e) => updateField('serviceType', e.target.value)}
        error={errors.serviceType}
        required
      >
        <option value="">เลือกบริการ</option>
        <option value="oil-change">เปลี่ยนถ่ายน้ำมันเครื่อง</option>
        <option value="general-maintenance">ตรวจเช็คทั่วไป</option>
        <option value="tire-change">เปลี่ยนยาง</option>
        <option value="full-service">บริการเต็มรูปแบบ</option>
      </Select>

      <Select
        label="ประเภทรถยนต์"
        value={formData.vehicleType}
        onChange={(e) => updateField('vehicleType', e.target.value)}
        error={errors.vehicleType}
        required
      >
        <option value="">เลือกประเภทรถ</option>
        <option value="sedan">รถเก๋ง</option>
        <option value="suv">รถ SUV</option>
        <option value="truck">รถกระบะ</option>
      </Select>

      <Input
        type="date"
        label="วันที่จอง"
        value={formData.bookingDate}
        onChange={(e) => updateField('bookingDate', e.target.value)}
        error={errors.bookingDate}
        min={new Date().toISOString().split('T')[^0]}
        required
      />

      {/* Availability indicator */}
      {checkingAvailability && (
        <div className="availability-checking">กำลังตรวจสอบความพร้อม...</div>
      )}

      {availability && (
        <div
          className={`availability-status ${
            availability.isAvailable ? 'available' : 'unavailable'
          }`}
        >
          {availability.isAvailable ? (
            <>
              ✓ ช่วงเวลานี้ว่าง - มีช่องว่าง {availability.availableSlots} ช่อง
            </>
          ) : (
            <>✗ ช่วงเวลานี้เต็มแล้ว กรุณาเลือกวันอื่น</>
          )}
        </div>
      )}

      {/* Price breakdown */}
      {priceBreakdown && (
        <div className="price-breakdown">
          <h3>รายละเอียดราคา</h3>
          <div className="price-row">
            <span>ราคาพื้นฐาน:</span>
            <span>{priceBreakdown.basePrice.toLocaleString()} บาท</span>
          </div>
          <div className="price-row">
            <span>ตัวคูณตามประเภทรถ:</span>
            <span>×{priceBreakdown.multiplier}</span>
          </div>
          {priceBreakdown.additionalServicesTotal > 0 && (
            <div className="price-row">
              <span>บริการเสริม:</span>
              <span>
                {priceBreakdown.additionalServicesTotal.toLocaleString()} บาท
              </span>
            </div>
          )}
          <div className="price-row total">
            <span>รวมทั้งหมด:</span>
            <span>{priceBreakdown.total.toLocaleString()} บาท</span>
          </div>
        </div>
      )}

      <Input
        type="textarea"
        label="หมายเหตุ (ถ้ามี)"
        value={formData.notes}
        onChange={(e) => updateField('notes', e.target.value)}
        rows={4}
      />

      {submitError && <div className="error-message">{submitError}</div>}

      <div className="form-actions">
        <Button type="button" variant="secondary" onClick={reset}>
          ล้างข้อมูล
        </Button>
        <Button
          type="submit"
          loading={loading}
          disabled={loading || !availability?.isAvailable}
        >
          {loading ? 'กำลังจอง...' : 'ยืนยันการจอง'}
        </Button>
      </div>
    </form>
  );
};

export default BookingForm;
```


## ตัวอย่างที่ 2: Repair Tracking

### Service Layer

**Repair Service** (`src/features/repairs/services/repairService.js`)

```javascript
import api from '../../../services/api';

const repairService = {
  // Get repair details
  getRepairById: async (repairId) => {
    const response = await api.get(`/repairs/${repairId}`);
    return response.data;
  },

  // Get repairs by booking ID
  getRepairsByBooking: async (bookingId) => {
    const response = await api.get(`/repairs`, {
      params: { bookingId },
    });
    return response.data;
  },

  // Subscribe to repair status updates (WebSocket/SSE simulation)
  subscribeToUpdates: (repairId, callback) => {
    // In real app, this would use WebSocket or Server-Sent Events
    const eventSource = new EventSource(`/api/v1/repairs/${repairId}/stream`);
    
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      callback(data);
    };

    eventSource.onerror = (error) => {
      console.error('SSE Error:', error);
      eventSource.close();
    };

    // Return cleanup function
    return () => eventSource.close();
  },

  // Calculate repair progress percentage
  calculateProgress: (repair) => {
    const statusWeights = {
      'pending': 0,
      'diagnosed': 20,
      'parts-ordered': 40,
      'in-progress': 60,
      'testing': 80,
      'completed': 100,
    };

    return statusWeights[repair.status] || 0;
  },

  // Estimate completion time
  estimateCompletion: (repair) => {
    const now = new Date();
    const startTime = new Date(repair.startedAt);
    const estimatedDuration = repair.estimatedDuration; // in minutes

    const elapsedMinutes = (now - startTime) / (1000 * 60);
    const remainingMinutes = Math.max(0, estimatedDuration - elapsedMinutes);

    return {
      elapsed: Math.round(elapsedMinutes),
      remaining: Math.round(remainingMinutes),
      completionTime: new Date(now.getTime() + remainingMinutes * 60 * 1000),
    };
  },

  // Format status for display
  formatStatus: (status) => {
    const statusMap = {
      'pending': 'รอดำเนินการ',
      'diagnosed': 'ตรวจสอบแล้ว',
      'parts-ordered': 'สั่งอะไหล่แล้ว',
      'in-progress': 'กำลังซ่อม',
      'testing': 'ทดสอบ',
      'completed': 'เสร็จสิ้น',
    };

    return statusMap[status] || status;
  },
};

export default repairService;
```


### Custom Hook Layer

**useRepairTracking Hook** (`src/features/repairs/hooks/useRepairTracking.js`)

```javascript
import { useState, useEffect, useCallback, useRef } from 'react';
import repairService from '../services/repairService';

/**
 * Custom hook for real-time repair tracking
 * Manages WebSocket connection and live updates
 */
const useRepairTracking = (repairId) => {
  const [repair, setRepair] = useState(null);
  const [progress, setProgress] = useState(0);
  const [timeEstimate, setTimeEstimate] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isLive, setIsLive] = useState(false);
  
  const unsubscribeRef = useRef(null);

  // Fetch initial repair data
  const fetchRepair = useCallback(async () => {
    if (!repairId) return;

    setLoading(true);
    setError(null);

    try {
      const data = await repairService.getRepairById(repairId);
      setRepair(data);
      
      // Calculate derived data
      const progressPercent = repairService.calculateProgress(data);
      setProgress(progressPercent);

      if (data.status !== 'completed') {
        const estimate = repairService.estimateCompletion(data);
        setTimeEstimate(estimate);
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch repair');
    } finally {
      setLoading(false);
    }
  }, [repairId]);

  // Subscribe to live updates
  useEffect(() => {
    if (!repairId || !repair) return;

    // Only subscribe if repair is not completed
    if (repair.status === 'completed') {
      setIsLive(false);
      return;
    }

    setIsLive(true);

    // Subscribe to updates
    unsubscribeRef.current = repairService.subscribeToUpdates(
      repairId,
      (update) => {
        setRepair((prev) => ({
          ...prev,
          ...update,
        }));

        // Recalculate progress
        const newProgress = repairService.calculateProgress(update);
        setProgress(newProgress);

        // Update time estimate
        if (update.status !== 'completed') {
          const estimate = repairService.estimateCompletion({
            ...repair,
            ...update,
          });
          setTimeEstimate(estimate);
        } else {
          setIsLive(false);
        }
      }
    );

    // Cleanup subscription on unmount
    return () => {
      if (unsubscribeRef.current) {
        unsubscribeRef.current();
      }
    };
  }, [repairId, repair]);

  // Update time estimate every minute
  useEffect(() => {
    if (!repair || repair.status === 'completed') return;

    const interval = setInterval(() => {
      const estimate = repairService.estimateCompletion(repair);
      setTimeEstimate(estimate);
    }, 60000); // Update every minute

    return () => clearInterval(interval);
  }, [repair]);

  // Load initial data
  useEffect(() => {
    fetchRepair();
  }, [fetchRepair]);

  return {
    repair,
    progress,
    timeEstimate,
    loading,
    error,
    isLive,
    refetch: fetchRepair,
  };
};

export default useRepairTracking;
```


### UI Component

**RepairStatusCard Component**

```jsx
import useRepairTracking from '../hooks/useRepairTracking';
import repairService from '../services/repairService';
import ProgressBar from '../../../components/common/ProgressBar/ProgressBar';
import './RepairStatusCard.css';

const RepairStatusCard = ({ repairId }) => {
  const {
    repair,
    progress,
    timeEstimate,
    loading,
    error,
    isLive,
  } = useRepairTracking(repairId);

  if (loading) return <div>กำลังโหลด...</div>;
  if (error) return <div>เกิดข้อผิดพลาด: {error}</div>;
  if (!repair) return <div>ไม่พบข้อมูลการซ่อม</div>;

  return (
    <div className="repair-status-card">
      {isLive && (
        <div className="live-indicator">
          <span className="pulse"></span>
          อัพเดทแบบ Real-time
        </div>
      )}

      <h3>สถานะการซ่อม</h3>
      
      <div className="status-badge">
        {repairService.formatStatus(repair.status)}
      </div>

      <ProgressBar value={progress} />

      {timeEstimate && (
        <div className="time-estimate">
          <p>เวลาที่ผ่านไป: {timeEstimate.elapsed} นาที</p>
          <p>เวลาที่เหลือ: {timeEstimate.remaining} นาที</p>
          <p>
            คาดว่าจะเสร็จ:{' '}
            {timeEstimate.completionTime.toLocaleTimeString('th-TH')}
          </p>
        </div>
      )}

      <div className="repair-details">
        <p><strong>ช่างผู้รับผิดชอบ:</strong> {repair.technician?.name}</p>
        <p><strong>ประเภทงาน:</strong> {repair.serviceType}</p>
        <p><strong>ค่าใช้จ่ายโดยประมาณ:</strong> {repair.estimatedCost} บาท</p>
      </div>
    </div>
  );
};

export default RepairStatusCard;
```


## สรุป Best Practices

### เมื่อใดควรใช้ Services

- ✅ API calls และ HTTP requests
- ✅ Pure business logic calculations
- ✅ Data transformation และ formatting
- ✅ Validation rules
- ✅ Constants และ configuration
- ✅ Logic ที่ต้องการใช้นอก React (tests, Node.js scripts)


### เมื่อใดควรใช้ Custom Hooks

- ✅ State management
- ✅ Side effects (useEffect)
- ✅ Lifecycle management
- ✅ Context consumption
- ✅ Combining multiple hooks
- ✅ React-specific logic
- ✅ Component rerender coordination


### ข้อดีของการแยก Logic

**Testability**: Services ทดสอบง่ายด้วย Jest, Hooks ทดสอบด้วย React Testing Library แยกกันชัดเจน[^1][^4]

**Reusability**: Services ใช้ได้ทุกที่, Hooks ใช้ได้ทุก React component[^6][^3]

**Maintainability**: แก้ไข logic ในที่เดียว ไม่กระทบ UI components[^7][^3]

**Separation of Concerns**: UI, Business Logic และ Data Layer แยกกันชัดเจน[^8][^1]

การออกแบบแบบนี้ทำให้ codebase มีความเป็นระเบียบ ง่ายต่อการ debug และ scale ได้ดี[^3][^7]
<span style="display:none">[^10][^9]</span>

<div align="center">⁂</div>

[^1]: https://profy.dev/article/react-architecture-business-logic-and-dependency-injection

[^2]: https://dev.to/chiangs/custom-react-hooks-vs-services-mcm/comments

[^3]: https://emoosavi.com/blog/decoupling-business-logic-from-ui-with-custom-react-hooks

[^4]: https://dev.to/chema/services-layer-approach-in-reactjs-1eo2

[^5]: https://stackoverflow.com/questions/76106495/why-now-devs-use-custom-hooks-instead-services-in-react-perspective-of-pre-hoo

[^6]: https://react.dev/learn/reusing-logic-with-custom-hooks

[^7]: https://www.developerway.com/posts/react-project-structure

[^8]: https://www.reddit.com/r/reactjs/comments/1g8i778/the_case_for_writing_business_logic_in_custom/

[^9]: https://stackoverflow.com/questions/79035369/splitting-component-logic-into-hooks

[^10]: https://www.youtube.com/watch?v=oi1amMsTD70

