โครงสร้างที่สร้างตาม modules_0.md + modules_1.md ใช้ package `com.icmon.module.auth`:

| Layer | ไฟล์ที่สร้าง |
|---|---|
| **domain/** | `MUser`, `MRole`, `MPermission`, `MUserToken`, enums 4 ตัว, valueobjects 2 ตัว |
| **application/interfaces/** | `AuthService`, `UserService`, `PermissionService` |
| **application/impl/** | `AuthServiceImpl`, `UserServiceImpl`, `PermissionServiceImpl` |
| **application/usecase/** | 8 use cases (Login, Logout, Refresh, CRUD User, ValidatePermission) |
| **infrastructure/entity/** | `UserEntity`, `RoleEntity`, `PermissionEntity`, `UserTokenEntity`, `UserRoleEntity`, `RolePermissionEntity` |
| **infrastructure/repository/** | 6 JPA repositories + 6 impl wrappers |
| **infrastructure/mapper/** | 6 MapStruct mappers |
| **infrastructure/security/** | `JwtTokenProvider`, `JwtTokenFilter`, `CustomUserDetailsService`, `SecurityConfig`, `PermissionInterceptor` |
| **infrastructure/cache/** | `CacheConfig`, `UserPermissionCacheService`, `TokenCacheService` |
| **infrastructure/ratelimit/** | `@RateLimit`, `RateLimiterInterceptor`, `RateLimitExceededException` |
| **presentation/controller/** | `AuthController`, `UserController`, `PermissionController` |
| **presentation/dto/** | 5 request DTOs + 4 response DTOs |
| **presentation/validator/** | `UserValidator` |
| **configuration/web/** | `WebConfig` (register RateLimiterInterceptor) |
| **db/migration/** | V1__auth_schema.sql (ตาราง 6 ตาราง) |

Updated todo list

Made changes.