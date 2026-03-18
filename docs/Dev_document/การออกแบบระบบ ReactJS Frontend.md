<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## การออกแบบระบบ ReactJS Frontend

ReactJS เป็น component-based library สำหรับสร้าง user interface ที่มีประสิทธิภาพ โดยใช้ Virtual DOM และ unidirectional data flow  สำหรับระบบศูนย์บริการรถยนต์ออนไลน์ ระบบ frontend จะประกอบด้วย authentication, state management, routing และ API integration[^1][^2][^3][^4]

## โครงสร้างโปรเจกต์ (Project Structure)

### Folder Structure

```
car-service-frontend/
├── public/
│   ├── index.html
│   ├── favicon.ico
│   └── assets/
│       └── images/
│
├── src/
│   ├── assets/                    # Static resources
│   │   ├── images/
│   │   ├── icons/
│   │   └── styles/
│   │       └── global.css
│   │
│   ├── components/                # Reusable components
│   │   ├── common/
│   │   │   ├── Button/
│   │   │   │   ├── Button.jsx
│   │   │   │   ├── Button.test.jsx
│   │   │   │   └── Button.module.css
│   │   │   ├── Input/
│   │   │   ├── Modal/
│   │   │   ├── Loading/
│   │   │   └── ErrorBoundary/
│   │   │
│   │   └── layout/
│   │       ├── Header/
│   │       ├── Sidebar/
│   │       ├── Footer/
│   │       └── Layout.jsx
│   │
│   ├── features/                  # Feature-based modules
│   │   ├── auth/
│   │   │   ├── components/
│   │   │   │   ├── LoginForm.jsx
│   │   │   │   └── RegisterForm.jsx
│   │   │   ├── hooks/
│   │   │   │   └── useAuth.js
│   │   │   ├── services/
│   │   │   │   └── authService.js
│   │   │   └── context/
│   │   │       └── AuthContext.jsx
│   │   │
│   │   ├── bookings/
│   │   │   ├── components/
│   │   │   │   ├── BookingList.jsx
│   │   │   │   ├── BookingForm.jsx
│   │   │   │   └── BookingCard.jsx
│   │   │   ├── hooks/
│   │   │   │   └── useBookings.js
│   │   │   └── services/
│   │   │       └── bookingService.js
│   │   │
│   │   ├── repairs/
│   │   ├── vehicles/
│   │   └── payments/
│   │
│   ├── pages/                     # Page components
│   │   ├── HomePage.jsx
│   │   ├── LoginPage.jsx
│   │   ├── DashboardPage.jsx
│   │   ├── BookingsPage.jsx
│   │   ├── RepairStatusPage.jsx
│   │   └── NotFoundPage.jsx
│   │
│   ├── hooks/                     # Global custom hooks
│   │   ├── useApi.js
│   │   ├── useForm.js
│   │   └── useDebounce.js
│   │
│   ├── services/                  # API services
│   │   ├── api.js                # Axios instance
│   │   └── endpoints.js          # API endpoints
│   │
│   ├── store/                     # State management
│   │   ├── slices/
│   │   │   ├── authSlice.js
│   │   │   ├── bookingSlice.js
│   │   │   └── repairSlice.js
│   │   └── store.js
│   │
│   ├── utils/                     # Utility functions
│   │   ├── formatters.js
│   │   ├── validators.js
│   │   └── constants.js
│   │
│   ├── routes/                    # Routing configuration
│   │   ├── ProtectedRoute.jsx
│   │   └── routes.jsx
│   │
│   ├── App.jsx                    # Root component
│   ├── main.jsx                   # Entry point
│   └── index.css
│
├── package.json
├── vite.config.js
└── .env
```


## Core Implementation

### 1. Authentication System

**Auth Context** (`src/features/auth/context/AuthContext.jsx`)

```jsx
import { createContext, useState, useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import authService from '../services/authService';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  // Check if user is logged in on mount
  useEffect(() => {
    const initAuth = async () => {
      if (token) {
        try {
          const userData = await authService.getProfile();
          setUser(userData);
        } catch (error) {
          console.error('Failed to fetch user profile:', error);
          logout();
        }
      }
      setLoading(false);
    };
    
    initAuth();
  }, [token]);

  // Login function
  const login = async (email, password) => {
    try {
      const response = await authService.login(email, password);
      const { token: accessToken, user: userData } = response;
      
      setToken(accessToken);
      setUser(userData);
      localStorage.setItem('token', accessToken);
      
      navigate('/dashboard');
      return { success: true };
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Login failed' 
      };
    }
  };

  // Register function
  const register = async (userData) => {
    try {
      const response = await authService.register(userData);
      const { token: accessToken, user: newUser } = response;
      
      setToken(accessToken);
      setUser(newUser);
      localStorage.setItem('token', accessToken);
      
      navigate('/dashboard');
      return { success: true };
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Registration failed' 
      };
    }
  };

  // Logout function
  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
    navigate('/login');
  };

  const value = {
    user,
    token,
    loading,
    isAuthenticated: !!token,
    login,
    register,
    logout,
  };

  return (
    <AuthContext.Provider value={value}>
      {!loading && children}
    </AuthContext.Provider>
  );
};

// Custom hook to use auth context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};

export default AuthContext;
```

**Auth Service** (`src/features/auth/services/authService.js`)

```javascript
import api from '../../../services/api';

const authService = {
  // Login user
  login: async (email, password) => {
    const response = await api.post('/auth/login', {
      email,
      password,
    });
    return response.data;
  },

  // Register new user
  register: async (userData) => {
    const response = await api.post('/auth/register', userData);
    return response.data;
  },

  // Get user profile
  getProfile: async () => {
    const response = await api.get('/auth/profile');
    return response.data;
  },

  // Refresh token
  refreshToken: async () => {
    const response = await api.post('/auth/refresh');
    return response.data;
  },

  // Logout
  logout: async () => {
    const response = await api.post('/auth/logout');
    return response.data;
  },
};

export default authService;
```


### 2. Axios Configuration

**API Instance** (`src/services/api.js`)

```javascript
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - Add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor - Handle errors and token refresh
api.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;

    // If error is 401 and we haven't retried yet
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Try to refresh token
        const refreshResponse = await axios.post(
          `${API_BASE_URL}/auth/refresh`,
          {},
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
          }
        );

        const { token } = refreshResponse.data;
        localStorage.setItem('token', token);

        // Retry original request with new token
        originalRequest.headers.Authorization = `Bearer ${token}`;
        return api(originalRequest);
      } catch (refreshError) {
        // Refresh failed, logout user
        localStorage.removeItem('token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    // Handle other errors
    if (error.response?.status === 403) {
      console.error('Access forbidden');
    }

    if (error.response?.status === 500) {
      console.error('Server error');
    }

    return Promise.reject(error);
  }
);

export default api;
```


### 3. Protected Routes

**Protected Route Component** (`src/routes/ProtectedRoute.jsx`)

```jsx
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../features/auth/context/AuthContext';
import Loading from '../components/common/Loading/Loading';

const ProtectedRoute = ({ children, requiredRoles = [] }) => {
  const { isAuthenticated, user, loading } = useAuth();
  const location = useLocation();

  if (loading) {
    return <Loading />;
  }

  if (!isAuthenticated) {
    // Redirect to login and save the attempted location
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Check if user has required roles
  if (requiredRoles.length > 0) {
    const hasRequiredRole = requiredRoles.some((role) =>
      user?.roles?.includes(role)
    );

    if (!hasRequiredRole) {
      return <Navigate to="/unauthorized" replace />;
    }
  }

  return children;
};

export default ProtectedRoute;
```

**Router Configuration** (`src/routes/routes.jsx`)

```jsx
import { createBrowserRouter } from 'react-router-dom';
import Layout from '../components/layout/Layout';
import ProtectedRoute from './ProtectedRoute';

// Pages
import HomePage from '../pages/HomePage';
import LoginPage from '../pages/LoginPage';
import RegisterPage from '../pages/RegisterPage';
import DashboardPage from '../pages/DashboardPage';
import BookingsPage from '../pages/BookingsPage';
import RepairStatusPage from '../pages/RepairStatusPage';
import AdminPage from '../pages/AdminPage';
import NotFoundPage from '../pages/NotFoundPage';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: 'login',
        element: <LoginPage />,
      },
      {
        path: 'register',
        element: <RegisterPage />,
      },
      {
        path: 'dashboard',
        element: (
          <ProtectedRoute>
            <DashboardPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'bookings',
        element: (
          <ProtectedRoute>
            <BookingsPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'repairs/:repairId',
        element: (
          <ProtectedRoute>
            <RepairStatusPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin',
        element: (
          <ProtectedRoute requiredRoles={['admin']}>
            <AdminPage />
          </ProtectedRoute>
        ),
      },
      {
        path: '*',
        element: <NotFoundPage />,
      },
    ],
  },
]);

export default router;
```


### 4. Custom Hooks

**useApi Hook** (`src/hooks/useApi.js`)

```javascript
import { useState, useEffect, useCallback } from 'react';
import api from '../services/api';

const useApi = (endpoint, options = {}) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const { 
    method = 'GET', 
    body = null, 
    immediate = true,
    onSuccess,
    onError 
  } = options;

  const execute = useCallback(
    async (overrideBody = null) => {
      setLoading(true);
      setError(null);

      try {
        const config = {
          method,
          url: endpoint,
        };

        if (overrideBody || body) {
          config.data = overrideBody || body;
        }

        const response = await api(config);
        setData(response.data);

        if (onSuccess) {
          onSuccess(response.data);
        }

        return response.data;
      } catch (err) {
        const errorMessage = err.response?.data?.message || err.message;
        setError(errorMessage);

        if (onError) {
          onError(errorMessage);
        }

        throw err;
      } finally {
        setLoading(false);
      }
    },
    [endpoint, method, body, onSuccess, onError]
  );

  useEffect(() => {
    if (immediate && method === 'GET') {
      execute();
    }
  }, [immediate, method, execute]);

  return { data, loading, error, execute, refetch: execute };
};

export default useApi;
```

**useForm Hook** (`src/hooks/useForm.js`)

```javascript
import { useState } from 'react';

const useForm = (initialValues, validationRules = {}) => {
  const [values, setValues] = useState(initialValues);
  const [errors, setErrors] = useState({});
  const [touched, setTouched] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Handle input change
  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    const newValue = type === 'checkbox' ? checked : value;

    setValues((prev) => ({
      ...prev,
      [name]: newValue,
    }));

    // Clear error when user starts typing
    if (errors[name]) {
      setErrors((prev) => ({
        ...prev,
        [name]: '',
      }));
    }
  };

  // Handle input blur
  const handleBlur = (e) => {
    const { name } = e.target;
    setTouched((prev) => ({
      ...prev,
      [name]: true,
    }));

    // Validate field on blur
    validateField(name, values[name]);
  };

  // Validate single field
  const validateField = (name, value) => {
    if (!validationRules[name]) return true;

    const rules = validationRules[name];
    let error = '';

    if (rules.required && !value) {
      error = rules.requiredMessage || `${name} is required`;
    } else if (rules.minLength && value.length < rules.minLength) {
      error = `${name} must be at least ${rules.minLength} characters`;
    } else if (rules.pattern && !rules.pattern.test(value)) {
      error = rules.patternMessage || `Invalid ${name} format`;
    } else if (rules.custom) {
      error = rules.custom(value, values);
    }

    if (error) {
      setErrors((prev) => ({
        ...prev,
        [name]: error,
      }));
      return false;
    }

    return true;
  };

  // Validate all fields
  const validate = () => {
    const newErrors = {};
    let isValid = true;

    Object.keys(validationRules).forEach((name) => {
      const fieldValid = validateField(name, values[name]);
      if (!fieldValid) {
        isValid = false;
      }
    });

    return isValid;
  };

  // Handle form submit
  const handleSubmit = (onSubmit) => async (e) => {
    e.preventDefault();
    setIsSubmitting(true);

    // Mark all fields as touched
    const allTouched = Object.keys(values).reduce((acc, key) => {
      acc[key] = true;
      return acc;
    }, {});
    setTouched(allTouched);

    // Validate
    if (validate()) {
      try {
        await onSubmit(values);
      } catch (error) {
        console.error('Form submission error:', error);
      }
    }

    setIsSubmitting(false);
  };

  // Reset form
  const reset = () => {
    setValues(initialValues);
    setErrors({});
    setTouched({});
    setIsSubmitting(false);
  };

  return {
    values,
    errors,
    touched,
    isSubmitting,
    handleChange,
    handleBlur,
    handleSubmit,
    reset,
    setValues,
    setErrors,
  };
};

export default useForm;
```


### 5. Feature Components

**Login Form** (`src/features/auth/components/LoginForm.jsx`)

```jsx
import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import useForm from '../../../hooks/useForm';
import Button from '../../../components/common/Button/Button';
import Input from '../../../components/common/Input/Input';
import './LoginForm.css';

const LoginForm = () => {
  const { login } = useAuth();
  const [apiError, setApiError] = useState('');

  const validationRules = {
    email: {
      required: true,
      pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
      patternMessage: 'Please enter a valid email',
    },
    password: {
      required: true,
      minLength: 6,
    },
  };

  const {
    values,
    errors,
    touched,
    isSubmitting,
    handleChange,
    handleBlur,
    handleSubmit,
  } = useForm(
    {
      email: '',
      password: '',
    },
    validationRules
  );

  const onSubmit = async (formValues) => {
    setApiError('');
    
    const result = await login(formValues.email, formValues.password);
    
    if (!result.success) {
      setApiError(result.error);
    }
  };

  return (
    <div className="login-form">
      <h2>เข้าสู่ระบบ</h2>

      <form onSubmit={handleSubmit(onSubmit)}>
        <Input
          label="อีเมล"
          type="email"
          name="email"
          value={values.email}
          onChange={handleChange}
          onBlur={handleBlur}
          error={touched.email && errors.email}
          placeholder="example@email.com"
          required
        />

        <Input
          label="รหัสผ่าน"
          type="password"
          name="password"
          value={values.password}
          onChange={handleChange}
          onBlur={handleBlur}
          error={touched.password && errors.password}
          placeholder="••••••••"
          required
        />

        {apiError && (
          <div className="error-message">{apiError}</div>
        )}

        <Button
          type="submit"
          fullWidth
          loading={isSubmitting}
          disabled={isSubmitting}
        >
          {isSubmitting ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ'}
        </Button>
      </form>

      <div className="form-footer">
        <p>
          ยังไม่มีบัญชี? <Link to="/register">สมัครสมาชิก</Link>
        </p>
      </div>
    </div>
  );
};

export default LoginForm;
```

**Booking List** (`src/features/bookings/components/BookingList.jsx`)

```jsx
import { useEffect, useState } from 'react';
import useApi from '../../../hooks/useApi';
import BookingCard from './BookingCard';
import Loading from '../../../components/common/Loading/Loading';
import './BookingList.css';

const BookingList = () => {
  const [filter, setFilter] = useState('all');
  
  const { 
    data: bookings, 
    loading, 
    error, 
    refetch 
  } = useApi('/bookings', {
    method: 'GET',
    immediate: true,
  });

  const filteredBookings = bookings?.filter((booking) => {
    if (filter === 'all') return true;
    return booking.status === filter;
  });

  const handleStatusChange = async () => {
    // Refresh list after status change
    await refetch();
  };

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return (
      <div className="error-container">
        <p>เกิดข้อผิดพลาด: {error}</p>
        <button onClick={refetch}>ลองใหม่</button>
      </div>
    );
  }

  return (
    <div className="booking-list">
      <div className="list-header">
        <h2>รายการจองของฉัน</h2>
        
        <div className="filter-buttons">
          <button
            className={filter === 'all' ? 'active' : ''}
            onClick={() => setFilter('all')}
          >
            ทั้งหมด
          </button>
          <button
            className={filter === 'pending' ? 'active' : ''}
            onClick={() => setFilter('pending')}
          >
            รอดำเนินการ
          </button>
          <button
            className={filter === 'confirmed' ? 'active' : ''}
            onClick={() => setFilter('confirmed')}
          >
            ยืนยันแล้ว
          </button>
          <button
            className={filter === 'completed' ? 'active' : ''}
            onClick={() => setFilter('completed')}
          >
            เสร็จสิ้น
          </button>
        </div>
      </div>

      <div className="booking-grid">
        {filteredBookings?.length > 0 ? (
          filteredBookings.map((booking) => (
            <BookingCard
              key={booking.id}
              booking={booking}
              onStatusChange={handleStatusChange}
            />
          ))
        ) : (
          <p className="no-bookings">ไม่มีรายการจอง</p>
        )}
      </div>
    </div>
  );
};

export default BookingList;
```


### 6. Main App Setup

**Main Entry** (`src/main.jsx`)

```jsx
import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';
import { AuthProvider } from './features/auth/context/AuthContext';
import router from './routes/routes';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <RouterProvider router={router}>
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </RouterProvider>
  </React.StrictMode>
);
```

**App Component** (`src/App.jsx`)

```jsx
import { Outlet } from 'react-router-dom';
import { AuthProvider } from './features/auth/context/AuthContext';
import ErrorBoundary from './components/common/ErrorBoundary/ErrorBoundary';
import './App.css';

function App() {
  return (
    <ErrorBoundary>
      <AuthProvider>
        <div className="app">
          <Outlet />
        </div>
      </AuthProvider>
    </ErrorBoundary>
  );
}

export default App;
```


## Application Flow Diagram

### User Authentication Flow

```
1. User visits app
   ↓
2. Check localStorage for token
   ↓
3. If token exists → Fetch user profile from API
   ├─ Success → Set user state → Navigate to Dashboard
   └─ Fail → Clear token → Navigate to Login
   ↓
4. User fills login form
   ↓
5. Form validation (useForm hook)
   ├─ Invalid → Show errors
   └─ Valid → Continue
   ↓
6. Submit to API (POST /auth/login)
   ↓
7. API Response
   ├─ Success (200) → Save token → Set user → Navigate to Dashboard
   └─ Error (401) → Show error message
```


### Data Fetching Flow

```
1. Component mounts
   ↓
2. useApi hook executes
   ↓
3. Check if immediate fetch required
   ↓
4. Axios interceptor adds Bearer token
   ↓
5. Send request to API
   ↓
6. API Response
   ├─ Success → Update state with data
   ├─ 401 Unauthorized → Refresh token → Retry request
   └─ Other errors → Set error state
   ↓
7. Component renders with data/error/loading
```


### Protected Route Flow

```
1. User navigates to protected route
   ↓
2. ProtectedRoute component checks authentication
   ↓
3. Check isAuthenticated
   ├─ No → Redirect to /login (save attempted location)
   └─ Yes → Continue
   ↓
4. Check required roles (if any)
   ├─ No required role → Redirect to /unauthorized
   └─ Has required role → Render children components
```


## สรุป Key Libraries

| Library | Purpose | Installation |
| :-- | :-- | :-- |
| `react-router-dom` | Routing and navigation | `npm install react-router-dom` |
| `axios` | HTTP client for API calls | `npm install axios` |
| `@reduxjs/toolkit` | State management (optional) | `npm install @reduxjs/toolkit react-redux` |
| `react-hook-form` | Form management (alternative) | `npm install react-hook-form` |
| `react-query` | Data fetching/caching (alternative) | `npm install @tanstack/react-query` |

ระบบนี้ใช้ modern React patterns รวมถึง custom hooks, context API, protected routes และ axios interceptors เพื่อสร้าง scalable และ maintainable frontend application[^2][^3][^4]
<span style="display:none">[^10][^11][^12][^13][^14][^15][^16][^17][^18][^19][^20][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://www.geeksforgeeks.org/reactjs/react-architecture-pattern-and-best-practices/

[^2]: https://www.creolestudios.com/reactjs-architecture-pattern/

[^3]: https://strapi.io/blog/react-and-nextjs-in-2025-modern-best-practices

[^4]: https://www.robinwieruch.de/react-router-authentication/

[^5]: https://dev.to/holasoymalva/why-i-decided-to-stop-working-with-reactjs-in-2025-4d1l

[^6]: https://react.dev

[^7]: https://dev.to/kinsflow/effective-patterns-for-shared-state-management-in-react-4ha1

[^8]: https://stackoverflow.com/questions/47188513/how-can-i-set-up-a-react-router-auth-flow

[^9]: https://www.bacancytechnology.com/blog/react-architecture-patterns-and-best-practices

[^10]: https://dionarodrigues.dev/blog/redux-toolkit-fundamentals-simplifying-state-management

[^11]: https://dev.to/pramod_boda/recommended-folder-structure-for-react-2025-48mc

[^12]: https://www.robinwieruch.de/react-folder-structure/

[^13]: https://www.linkedin.com/pulse/production-grade-react-project-structure-from-setup-gmokc

[^14]: https://javascript.plainenglish.io/react-best-practices-for-folder-structure-system-design-architecture-8fc2f09e3fff

[^15]: https://www.reddit.com/r/reactjs/comments/10nakgt/react_project_folder_structure/

[^16]: https://www.patterns.dev/react/hooks-pattern/

[^17]: https://www.youtube.com/watch?v=X3qyxo_UTR4

[^18]: https://namastedev.com/blog/best-practices-for-folder-structure-in-react-6/

[^19]: https://www.youtube.com/watch?v=I2Bgi0Qcdvc

[^20]: https://www.digitalocean.com/community/tutorials/react-axios-react

