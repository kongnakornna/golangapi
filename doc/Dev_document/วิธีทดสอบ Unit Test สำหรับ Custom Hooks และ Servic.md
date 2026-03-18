<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## วิธีทดสอบ Unit Test สำหรับ Custom Hooks และ Services

การทดสอบ custom hooks และ services ต้องใช้เครื่องมือและวิธีการที่แตกต่างกัน โดย Services ทดสอบง่ายกว่าเพราะเป็น pure JavaScript functions ส่วน Custom Hooks ต้องใช้ React Testing Library[^1][^2]

## การติดตั้ง Testing Libraries

```bash
# Install Jest and React Testing Library
npm install --save-dev jest @testing-library/react @testing-library/jest-dom @testing-library/user-event

# Install axios mock adapter (optional)
npm install --save-dev axios-mock-adapter jest-mock-extended
```


## ตัวอย่างที่ 1: Unit Testing Services

### Service Code

**bookingService.js**

```javascript
import api from '../../../services/api';

const bookingService = {
  getAllBookings: async () => {
    const response = await api.get('/bookings');
    return response.data;
  },

  createBooking: async (bookingData) => {
    const response = await api.post('/bookings', bookingData);
    return response.data;
  },

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

  validateBookingData: (data) => {
    const errors = {};

    if (!data.customerId) {
      errors.customerId = 'Customer ID is required';
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

    return {
      isValid: Object.keys(errors).length === 0,
      errors,
    };
  },
};

export default bookingService;
```


### Service Unit Tests

**bookingService.test.js**

```javascript
import bookingService from './bookingService';
import api from '../../../services/api';

// Mock the API module
jest.mock('../../../services/api');

describe('bookingService', () => {
  // Reset mocks before each test
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('getAllBookings', () => {
    it('should fetch all bookings successfully', async () => {
      // Arrange: Setup mock data
      const mockBookings = [
        { id: '1', customerId: 'c1', status: 'confirmed' },
        { id: '2', customerId: 'c2', status: 'pending' },
      ];

      // Mock axios response
      api.get.mockResolvedValue({ data: mockBookings });

      // Act: Call the service
      const result = await bookingService.getAllBookings();

      // Assert: Verify results
      expect(result).toEqual(mockBookings);
      expect(api.get).toHaveBeenCalledWith('/bookings');
      expect(api.get).toHaveBeenCalledTimes(1);
    });

    it('should handle errors when fetching bookings fails', async () => {
      // Arrange: Mock error response
      const errorMessage = 'Network Error';
      api.get.mockRejectedValue(new Error(errorMessage));

      // Act & Assert: Expect rejection
      await expect(bookingService.getAllBookings()).rejects.toThrow(
        errorMessage
      );
    });
  });

  describe('createBooking', () => {
    it('should create booking successfully', async () => {
      // Arrange
      const bookingData = {
        customerId: 'c1',
        vehicleId: 'v1',
        bookingDate: '2025-12-01',
        serviceType: 'oil-change',
      };

      const mockResponse = {
        id: 'b1',
        ...bookingData,
        status: 'pending',
      };

      api.post.mockResolvedValue({ data: mockResponse });

      // Act
      const result = await bookingService.createBooking(bookingData);

      // Assert
      expect(result).toEqual(mockResponse);
      expect(api.post).toHaveBeenCalledWith('/bookings', bookingData);
    });

    it('should handle validation errors from API', async () => {
      // Arrange
      const bookingData = { customerId: 'c1' };
      const errorResponse = {
        response: {
          status: 400,
          data: { message: 'Validation failed' },
        },
      };

      api.post.mockRejectedValue(errorResponse);

      // Act & Assert
      await expect(bookingService.createBooking(bookingData)).rejects.toEqual(
        errorResponse
      );
    });
  });

  describe('calculatePrice', () => {
    it('should calculate price for sedan oil change', () => {
      // Act
      const result = bookingService.calculatePrice('oil-change', 'sedan', []);

      // Assert
      expect(result).toEqual({
        basePrice: 800,
        multiplier: 1.0,
        additionalServicesTotal: 0,
        total: 800,
      });
    });

    it('should calculate price for SUV with multiplier', () => {
      // Act
      const result = bookingService.calculatePrice('oil-change', 'suv', []);

      // Assert
      expect(result).toEqual({
        basePrice: 800,
        multiplier: 1.3,
        additionalServicesTotal: 0,
        total: 1040, // 800 * 1.3
      });
    });

    it('should include additional services in total', () => {
      // Arrange
      const additionalServices = [
        { id: 's1', name: 'Air Filter', price: 200 },
        { id: 's2', name: 'Brake Check', price: 300 },
      ];

      // Act
      const result = bookingService.calculatePrice(
        'oil-change',
        'sedan',
        additionalServices
      );

      // Assert
      expect(result).toEqual({
        basePrice: 800,
        multiplier: 1.0,
        additionalServicesTotal: 500,
        total: 1300, // 800 + 200 + 300
      });
    });

    it('should handle unknown service type', () => {
      // Act
      const result = bookingService.calculatePrice('unknown', 'sedan', []);

      // Assert
      expect(result.total).toBe(0);
    });
  });

  describe('validateBookingData', () => {
    it('should validate correct booking data', () => {
      // Arrange
      const validData = {
        customerId: 'c1',
        vehicleId: 'v1',
        bookingDate: '2025-12-31',
        serviceType: 'oil-change',
      };

      // Act
      const result = bookingService.validateBookingData(validData);

      // Assert
      expect(result.isValid).toBe(true);
      expect(result.errors).toEqual({});
    });

    it('should return error when customerId is missing', () => {
      // Arrange
      const invalidData = {
        bookingDate: '2025-12-31',
      };

      // Act
      const result = bookingService.validateBookingData(invalidData);

      // Assert
      expect(result.isValid).toBe(false);
      expect(result.errors.customerId).toBe('Customer ID is required');
    });

    it('should return error when booking date is in the past', () => {
      // Arrange
      const yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1);
      
      const invalidData = {
        customerId: 'c1',
        bookingDate: yesterday.toISOString().split('T')[^0],
      };

      // Act
      const result = bookingService.validateBookingData(invalidData);

      // Assert
      expect(result.isValid).toBe(false);
      expect(result.errors.bookingDate).toBe('Cannot book in the past');
    });

    it('should return multiple errors for multiple issues', () => {
      // Arrange
      const invalidData = {};

      // Act
      const result = bookingService.validateBookingData(invalidData);

      // Assert
      expect(result.isValid).toBe(false);
      expect(Object.keys(result.errors).length).toBeGreaterThan(1);
    });
  });
});
```


## ตัวอย่างที่ 2: Unit Testing Custom Hooks

### Custom Hook Code

**useBookings.js**

```javascript
import { useState, useEffect, useCallback } from 'react';
import bookingService from '../services/bookingService';

const useBookings = () => {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchBookings = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await bookingService.getAllBookings();
      setBookings(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchBookings();
  }, [fetchBookings]);

  const createBooking = useCallback(async (bookingData) => {
    setLoading(true);
    setError(null);

    try {
      const newBooking = await bookingService.createBooking(bookingData);
      setBookings((prev) => [newBooking, ...prev]);
      return { success: true, data: newBooking };
    } catch (err) {
      setError(err.message);
      return { success: false, error: err.message };
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    bookings,
    loading,
    error,
    fetchBookings,
    createBooking,
  };
};

export default useBookings;
```


### Custom Hook Unit Tests

**useBookings.test.js**

```javascript
import { renderHook, waitFor, act } from '@testing-library/react';
import useBookings from './useBookings';
import bookingService from '../services/bookingService';

// Mock the service
jest.mock('../services/bookingService');

describe('useBookings', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should initialize with empty state', () => {
    // Arrange
    bookingService.getAllBookings.mockResolvedValue([]);

    // Act
    const { result } = renderHook(() => useBookings());

    // Assert
    expect(result.current.bookings).toEqual([]);
    expect(result.current.loading).toBe(true); // Initially loading
    expect(result.current.error).toBe(null);
  });

  it('should fetch bookings on mount', async () => {
    // Arrange
    const mockBookings = [
      { id: '1', customerId: 'c1', status: 'confirmed' },
      { id: '2', customerId: 'c2', status: 'pending' },
    ];
    bookingService.getAllBookings.mockResolvedValue(mockBookings);

    // Act
    const { result } = renderHook(() => useBookings());

    // Assert - Wait for async operation to complete
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.bookings).toEqual(mockBookings);
    expect(result.current.error).toBe(null);
    expect(bookingService.getAllBookings).toHaveBeenCalledTimes(1);
  });

  it('should handle fetch errors', async () => {
    // Arrange
    const errorMessage = 'Failed to fetch bookings';
    bookingService.getAllBookings.mockRejectedValue(
      new Error(errorMessage)
    );

    // Act
    const { result } = renderHook(() => useBookings());

    // Assert
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.bookings).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should create booking successfully', async () => {
    // Arrange
    bookingService.getAllBookings.mockResolvedValue([]);
    
    const newBookingData = {
      customerId: 'c1',
      vehicleId: 'v1',
      bookingDate: '2025-12-31',
    };

    const mockNewBooking = {
      id: 'b1',
      ...newBookingData,
      status: 'pending',
    };

    bookingService.createBooking.mockResolvedValue(mockNewBooking);

    // Act
    const { result } = renderHook(() => useBookings());

    // Wait for initial fetch
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let createResult;
    await act(async () => {
      createResult = await result.current.createBooking(newBookingData);
    });

    // Assert
    expect(createResult.success).toBe(true);
    expect(createResult.data).toEqual(mockNewBooking);
    expect(result.current.bookings).toHaveLength(1);
    expect(result.current.bookings[^0]).toEqual(mockNewBooking);
    expect(bookingService.createBooking).toHaveBeenCalledWith(newBookingData);
  });

  it('should handle create booking errors', async () => {
    // Arrange
    bookingService.getAllBookings.mockResolvedValue([]);
    const errorMessage = 'Validation failed';
    bookingService.createBooking.mockRejectedValue(new Error(errorMessage));

    // Act
    const { result } = renderHook(() => useBookings());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let createResult;
    await act(async () => {
      createResult = await result.current.createBooking({});
    });

    // Assert
    expect(createResult.success).toBe(false);
    expect(createResult.error).toBe(errorMessage);
    expect(result.current.error).toBe(errorMessage);
    expect(result.current.bookings).toHaveLength(0);
  });

  it('should refetch bookings when fetchBookings is called', async () => {
    // Arrange
    const initialBookings = [{ id: '1', status: 'pending' }];
    const updatedBookings = [
      { id: '1', status: 'confirmed' },
      { id: '2', status: 'pending' },
    ];

    bookingService.getAllBookings
      .mockResolvedValueOnce(initialBookings)
      .mockResolvedValueOnce(updatedBookings);

    // Act
    const { result } = renderHook(() => useBookings());

    await waitFor(() => {
      expect(result.current.bookings).toEqual(initialBookings);
    });

    // Call refetch
    await act(async () => {
      await result.current.fetchBookings();
    });

    // Assert
    expect(result.current.bookings).toEqual(updatedBookings);
    expect(bookingService.getAllBookings).toHaveBeenCalledTimes(2);
  });
});
```


## ตัวอย่างที่ 3: Testing Authentication with localStorage

### Auth Service

**authService.js**

```javascript
import api from '../../../services/api';

const authService = {
  login: async (email, password) => {
    const response = await api.post('/auth/login', { email, password });
    const { token, user } = response.data;
    
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    
    return response.data;
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  getToken: () => {
    return localStorage.getItem('token');
  },

  isAuthenticated: () => {
    return !!localStorage.getItem('token');
  },
};

export default authService;
```


### Auth Service Tests with localStorage Mock

**authService.test.js**

```javascript
import authService from './authService';
import api from '../../../services/api';

// Mock API
jest.mock('../../../services/api');

// Mock localStorage
const localStorageMock = (() => {
  let store = {};

  return {
    getItem: jest.fn((key) => store[key] || null),
    setItem: jest.fn((key, value) => {
      store[key] = value.toString();
    }),
    removeItem: jest.fn((key) => {
      delete store[key];
    }),
    clear: jest.fn(() => {
      store = {};
    }),
  };
})();

// Replace global localStorage with mock
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

describe('authService', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks();
    localStorageMock.clear();
  });

  describe('login', () => {
    it('should login successfully and store token', async () => {
      // Arrange
      const mockResponse = {
        token: 'mock-jwt-token',
        user: { id: 'u1', email: 'test@example.com', name: 'Test User' },
      };

      api.post.mockResolvedValue({ data: mockResponse });

      // Act
      const result = await authService.login('test@example.com', 'password123');

      // Assert
      expect(result).toEqual(mockResponse);
      expect(api.post).toHaveBeenCalledWith('/auth/login', {
        email: 'test@example.com',
        password: 'password123',
      });
      
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'token',
        'mock-jwt-token'
      );
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'user',
        JSON.stringify(mockResponse.user)
      );
    });

    it('should handle login errors', async () => {
      // Arrange
      const errorResponse = {
        response: {
          status: 401,
          data: { message: 'Invalid credentials' },
        },
      };

      api.post.mockRejectedValue(errorResponse);

      // Act & Assert
      await expect(
        authService.login('test@example.com', 'wrongpassword')
      ).rejects.toEqual(errorResponse);

      expect(localStorageMock.setItem).not.toHaveBeenCalled();
    });
  });

  describe('logout', () => {
    it('should remove token and user from localStorage', () => {
      // Arrange
      localStorageMock.setItem('token', 'mock-token');
      localStorageMock.setItem('user', JSON.stringify({ id: 'u1' }));

      // Act
      authService.logout();

      // Assert
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('token');
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('user');
    });
  });

  describe('getToken', () => {
    it('should return token from localStorage', () => {
      // Arrange
      const mockToken = 'mock-jwt-token';
      localStorageMock.setItem('token', mockToken);

      // Act
      const result = authService.getToken();

      // Assert
      expect(result).toBe(mockToken);
      expect(localStorageMock.getItem).toHaveBeenCalledWith('token');
    });

    it('should return null when no token exists', () => {
      // Act
      const result = authService.getToken();

      // Assert
      expect(result).toBeNull();
    });
  });

  describe('isAuthenticated', () => {
    it('should return true when token exists', () => {
      // Arrange
      localStorageMock.setItem('token', 'mock-token');

      // Act
      const result = authService.isAuthenticated();

      // Assert
      expect(result).toBe(true);
    });

    it('should return false when no token exists', () => {
      // Act
      const result = authService.isAuthenticated();

      // Assert
      expect(result).toBe(false);
    });
  });
});
```


## Best Practices สำหรับการทดสอบ

### 1. AAA Pattern (Arrange-Act-Assert)

```javascript
it('should do something', () => {
  // Arrange: Set up test data and mocks
  const input = { /* ... */ };
  service.method.mockResolvedValue(expectedOutput);
  
  // Act: Execute the code under test
  const result = functionUnderTest(input);
  
  // Assert: Verify the results
  expect(result).toEqual(expectedOutput);
  expect(service.method).toHaveBeenCalledWith(input);
});
```


### 2. ใช้ describe เพื่อจัดกลุ่ม Tests

```javascript
describe('bookingService', () => {
  describe('getAllBookings', () => {
    it('should handle success case', () => { /* ... */ });
    it('should handle error case', () => { /* ... */ });
  });
  
  describe('createBooking', () => {
    it('should create successfully', () => { /* ... */ });
  });
});
```


### 3. Clean Up After Each Test

```javascript
beforeEach(() => {
  jest.clearAllMocks(); // Clear mock call history
  localStorage.clear(); // Clear localStorage
});

afterEach(() => {
  jest.restoreAllMocks(); // Restore original implementations
});
```


### 4. Test Edge Cases

```javascript
it('should handle empty array', () => { /* ... */ });
it('should handle null input', () => { /* ... */ });
it('should handle undefined values', () => { /* ... */ });
it('should handle network timeout', () => { /* ... */ });
```


### 5. ใช้ waitFor สำหรับ Async Operations

```javascript
await waitFor(() => {
  expect(result.current.loading).toBe(false);
});

await waitFor(() => {
  expect(result.current.data).toBeDefined();
}, { timeout: 3000 });
```


## Package.json Configuration

```json
{
  "scripts": {
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage",
    "test:verbose": "jest --verbose"
  },
  "jest": {
    "testEnvironment": "jsdom",
    "setupFilesAfterEnv": ["<rootDir>/src/setupTests.js"],
    "moduleNameMapper": {
      "\\.(css|less|scss|sass)$": "identity-obj-proxy"
    },
    "collectCoverageFrom": [
      "src/**/*.{js,jsx}",
      "!src/index.js",
      "!src/reportWebVitals.js"
    ],
    "coverageThreshold": {
      "global": {
        "branches": 80,
        "functions": 80,
        "lines": 80,
        "statements": 80
      }
    }
  }
}
```


## สรุปความแตกต่างการทดสอบ

| Aspect | Services | Custom Hooks |
| :-- | :-- | :-- |
| **Testing Tool** | Jest only | Jest + React Testing Library |
| **Render Required** | ❌ No | ✅ Yes (renderHook) |
| **Use act()** | ❌ No | ✅ Yes for state updates |
| **Use waitFor()** | ❌ No (use async/await) | ✅ Yes for async operations |
| **Mock Dependencies** | API, localStorage | Services, Context |
| **Complexity** | Simple | More complex |

การทดสอบ Services ง่ายกว่าเพราะเป็น pure functions ไม่มี React dependencies  ส่วน Custom Hooks ต้องใช้ `renderHook()` และ `act()` เพื่อจำลอง React environment[^3][^1][^2]
<span style="display:none">[^10][^11][^12][^13][^14][^15][^16][^17][^18][^19][^20][^4][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://www.builder.io/blog/test-custom-hooks-react-testing-library

[^2]: https://kentcdodds.com/blog/how-to-test-custom-react-hooks

[^3]: https://www.dhiwise.com/blog/design-converter/testing-library-renderhook-a-quick-start-guide

[^4]: https://dev.to/manuartero/testing-a-custom-hook-like-a-pro-1b19

[^5]: https://stackoverflow.com/questions/67173251/testing-component-with-custom-hook

[^6]: https://tillitsdone.com/blogs/test-custom-hooks-in-react-apps/

[^7]: https://semaphore.io/blog/unit-tests-nodejs-jest

[^8]: https://www.youtube.com/watch?v=Ru4V8yCR6jQ

[^9]: https://blog.appsignal.com/2024/11/27/unit-testing-in-nodejs-with-jest.html

[^10]: https://stackoverflow.com/questions/64287448/react-hook-testing-with-renderhook

[^11]: https://www.csrhymes.com/2022/03/09/mocking-axios-with-jest-and-typescript.html

[^12]: https://stackoverflow.com/questions/45016033/how-do-i-test-axios-in-jest

[^13]: https://jestjs.io/docs/mock-functions

[^14]: https://github.com/axios/axios/issues/6112

[^15]: https://www.browserstack.com/guide/mock-api-calls-in-jest

[^16]: https://www.testim.io/blog/react-testing-library-waitfor/

[^17]: https://github.com/jestjs/jest/issues/2098

[^18]: https://www.maxpou.fr/blog/jest-mock-axios-calls/

[^19]: https://stackoverflow.com/questions/65725591/react-testing-library-how-to-use-waitfor

[^20]: https://stackoverflow.com/questions/32911630/how-do-i-deal-with-localstorage-in-jest-tests

