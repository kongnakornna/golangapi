# Authentication, Users & System Modules API

## Authentication — `/api/auth`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/auth/login` | `Login` | public (rate-limited) |
| POST | `/auth/signin` | `SignIn` | public (rate-limited) |
| GET | `/auth/publickey` | `GetPublicKey` | public |
| GET | `/auth/verifyemail` | `VerifyEmail` | public |
| POST | `/auth/forgotpassword` | `ForgotPassword` | public |
| PATCH | `/auth/resetpassword` | `ResetPassword` | public |
| GET | `/auth/refresh` | `RefreshToken` | auth |
| GET | `/auth/logout` | `Logout` | auth |
| GET | `/auth/logoutall` | `LogoutAllToken` | auth |
| POST | `/auth/logoutall` | `LogoutAllToken` | auth |

### POST `/auth/login`
Body: `{ "username": "kongnakornna", "password": "password" }`
Response: `{ "access_token": "eyJ...", "refresh_token": "eyJ...", "token_type": "Bearer" }`
Error: `400` `401` `404`

### POST `/auth/signin`
Body: `{ "email": "kongnakornna@gmail.com", "password": "password" }`

### GET `/auth/publickey`
Response:
```json
{
  "public_key_access_token": "-----BEGIN PUBLIC KEY-----\n...",
  "public_key_refresh_token": "-----BEGIN PUBLIC KEY-----\n..."
}
```

### GET `/auth/verifyemail`
Query: `code` (string, required)

### POST `/auth/forgotpassword`
Body: `{ "email": "kongnakornna@gmail.com" }`

### PATCH `/auth/resetpassword`
Body: `{ "code": "...", "password": "newpassword" }` (ตาม handler)

---

## Users — `/api` + `/api/user`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/register` | `Register` | public |
| POST | `/signin` | `SignInEmail` | public |
| POST | `/login` | `SignInUsername` | public |
| GET | `/user/me` | `Me` | auth |
| PUT | `/user/me` | `UpdateMe` | auth |
| PATCH | `/user/me/updatepass` | `UpdatePasswordMe` | auth |
| POST | `/user/me/avatar` | `UploadAvatar` | auth |
| GET | `/user/profile/{id}` | `Profile` | auth |
| GET | `/user/` | `GetMulti` | auth + superuser |
| POST | `/user/` | `Create` | auth + superuser |
| PATCH | `/user/{id}/role` | `UpdateRole` | auth + superuser |
| GET | `/user/list` | `ListUsers` | auth + superuser |
| GET | `/user/statistics` | `Statistics` | auth + superuser |
| GET | `/user/notify/{channel}` | `NotifyList` | auth + superuser |
| PATCH | `/user/{id}/activestatus` | `UpdateActiveStatus` | auth + superuser |
| GET | `/user/{id}` | `Get` | auth |
| DELETE | `/user/{id}` | `Delete` | auth + superuser |
| PUT | `/user/{id}` | `Update` | auth + superuser |
| PATCH | `/user/{id}/updatepass` | `UpdatePassword` | auth + superuser |
| GET | `/user/{id}/logoutall` | `LogoutAllAdmin` | auth + superuser |

### POST `/register`
Body: `{ "username", "email", "password" }` (ตาม handler)

### GET `/user/me`
Response: profile ของ user ที่ login

---

## Items — `/api/item`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| GET | `/item/` | `GetMulti` | auth |
| POST | `/item/` | `Create` | auth |
| GET | `/item/{id}` | `Get` | auth |
| DELETE | `/item/{id}` | `Delete` | auth |
| PUT | `/item/{id}` | `Update` | auth |

---

## Orders (Kafka) — `/api/orders`

| Method | Path | Handler | Access |
|--------|------|---------|--------|
| POST | `/orders` | (Kafka producer) | auth |

---

## WebSocket

| Method | Path | Access |
|--------|------|--------|
| GET | `/api/ws` | public (handshake) |

---

## System & Health

| Method | Path | Access |
|--------|------|--------|
| GET | `/` | public |
| GET | `/health` | public |
| GET | `/api/ping` | public |
| GET | `/api/health` | public |
| GET | `/metrics` | public (Prometheus) |
| GET | `/apimetric` | public |
| GET | `/uploads/avatar/*` | public (static) |
| GET | `/swagger/*` | public (Swagger UI) |
