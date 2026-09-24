# 📁 โครงสร้างโปรเจกต์ Spring Boot DDD Template
 
## 🏗️ โครงสร้างหลัก (Root Structure)

```
spring-boot-ddd-template/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/
│   │   │       └── template/
│   │   │           └── app/
│   │   │               ├── Application.java
│   │   │               ├── _shared/
│   │   │               ├── modules/
│   │   │               ├── configuration/
│   │   │               ├── exception/
│   │   │               ├── logging/
│   │   │               └── utils/
│   │   └── resources/
│   │       ├── application.yml
│   │       ├── application-dev.yml
│   │       ├── application-prod.yml
│   │       ├── application-test.yml
│   │       ├── static/
│   │       │   └── template/
│   │       │       └── jrxml/
│   │       ├── i18n/
│   │       └── db/
│   │           └── migration/
│   └── test/
│       └── java/
│           └── com/
│               └── template/
│                   └── app/
│                       ├── _shared/
│                       └── modules/
├── docker/
│   ├── docker-compose.yml
│   ├── Dockerfile
│   └── .env.example
├── pom.xml
├── README.md
├── ARCHITECTURE.md
└── LICENSE
```

---

## 📦 โมดูลที่ 1: Authentication & Permission (`modules/auth`)

โมดูลนี้เป็นโมดูลพื้นฐานสำหรับการจัดการความปลอดภัยของระบบ ประกอบด้วยการเข้าสู่ระบบ การจัดการผู้ใช้ และการควบคุมสิทธิ์

```
modules/auth/
├── application/
│   ├── interfaces/
│   │   ├── AuthService.java
│   │   ├── UserService.java
│   │   └── PermissionService.java
│   ├── impl/
│   │   ├── AuthServiceImpl.java
│   │   ├── UserServiceImpl.java
│   │   └── PermissionServiceImpl.java
│   └── usecase/
│       ├── LoginUseCase.java
│       ├── LogoutUseCase.java
│       ├── RefreshTokenUseCase.java
│       ├── CreateUserUseCase.java
│       ├── UpdateUserUseCase.java
│       ├── DeleteUserUseCase.java
│       ├── GetUserUseCase.java
│       └── ValidatePermissionUseCase.java
├── domain/
│   ├── MUser.java
│   ├── MUserMenu.java
│   ├── MUserJobRole.java
│   ├── MUserToken.java
│   ├── enums/
│   │   ├── UserStatus.java
│   │   ├── RoleType.java
│   │   └── PermissionType.java
│   └── valueobjects/
│       ├── Email.java
│       └── PasswordHash.java
├── infrastructure/
│   ├── repository/
│   │   ├── UserRepository.java
│   │   ├── UserMenuRepository.java
│   │   ├── UserJobRoleRepository.java
│   │   ├── UserTokenRepository.java
│   │   └── impl/
│   │       ├── UserRepositoryImpl.java
│   │       ├── UserMenuRepositoryImpl.java
│   │       ├── UserJobRoleRepositoryImpl.java
│   │       └── UserTokenRepositoryImpl.java
│   ├── security/
│   │   ├── JwtTokenFilter.java
│   │   ├── JwtTokenProvider.java
│   │   ├── PermissionInterceptor.java
│   │   ├── CustomUserDetailsService.java
│   │   ├── SecurityConfig.java
│   │   ├── CustomAccessDeniedHandler.java
│   │   └── CustomAuthenticationEntryPoint.java
│   ├── entity/
│   │   ├── UserEntity.java
│   │   ├── UserMenuEntity.java
│   │   ├── UserJobRoleEntity.java
│   │   └── UserTokenEntity.java
│   └── mapper/
│       ├── UserMapper.java
│       ├── UserMenuMapper.java
│       ├── UserJobRoleMapper.java
│       └── UserTokenMapper.java
└── presentation/
    ├── controller/
    │   ├── AuthController.java
    │   ├── UserController.java
    │   ├── PermissionController.java
    │   └── RoleController.java
    ├── dto/
    │   ├── request/
    │   │   ├── LoginRequestDTO.java
    │   │   ├── RegisterRequestDTO.java
    │   │   ├── UserCreateDTO.java
    │   │   ├── UserUpdateDTO.java
    │   │   └── PermissionRequestDTO.java
    │   └── response/
    │       ├── LoginResponseDTO.java
    │       ├── UserResponseDTO.java
    │       ├── PermissionResponseDTO.java
    │       └── TokenResponseDTO.java
    └── validator/
        └── UserValidator.java
```

### 🔍 คำอธิบายแต่ละส่วน

| ส่วนประกอบ | คำอธิบาย |
|-----------|----------|
| **application/** | ชั้น Application ประกอบด้วย Service และ Use Case |
| **application/interfaces/** | Interface ของ Service แต่ละตัว |
| **application/impl/** | Implementation ของ Service |
| **application/usecase/** | Use Case แต่ละอย่าง (Single Responsibility) |
| **domain/** | ชั้น Domain ประกอบด้วย Entity และ Value Objects |
| **domain/enums/** | Enum ต่างๆ เช่น สถานะผู้ใช้, บทบาท |
| **domain/valueobjects/** | Value Objects เช่น Email, PasswordHash |
| **infrastructure/** | ชั้น Infrastructure ประกอบด้วย Repository, Security, Entity, Mapper |
| **infrastructure/repository/** | Interface และ Implementation ของ Repository |
| **infrastructure/security/** | เรื่องความปลอดภัย JWT, Filter, Interceptor |
| **infrastructure/entity/** | JPA Entity สำหรับฐานข้อมูล |
| **infrastructure/mapper/** | MapStruct Mapper |
| **presentation/** | ชั้น Presentation ประกอบด้วย Controller และ DTO |
| **presentation/dto/request/** | Request DTO |
| **presentation/dto/response/** | Response DTO |
| **presentation/validator/** | Validator เฉพาะของโมดูล |

### 📝 ไฟล์สำคัญตัวอย่าง

#### `application/interfaces/AuthService.java`
```java
package com.template.app.modules.auth.application.interfaces;

import com.template.app.modules.auth.presentation.dto.request.LoginRequestDTO;
import com.template.app.modules.auth.presentation.dto.response.LoginResponseDTO;
import com.template.app.exception.SystemGlobalException;

public interface AuthService {
    LoginResponseDTO login(LoginRequestDTO request) throws SystemGlobalException;
    void logout(String token) throws SystemGlobalException;
    LoginResponseDTO refreshToken(String refreshToken) throws SystemGlobalException;
    boolean validateToken(String token);
}
```

#### `domain/MUser.java`
```java
package com.template.app.modules.auth.domain;

import com.template.app._shared.domain.GenericBusinessClass;
import com.template.app.modules.auth.domain.valueobjects.Email;
import com.template.app.modules.auth.domain.enums.UserStatus;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@EqualsAndHashCode(callSuper = true)
public class MUser extends GenericBusinessClass {
    private String username;
    private Email email;
    private String passwordHash;
    private String fullName;
    private UserStatus status;
    private String phoneNumber;
    private String profileImageUrl;
    private RoleType role;
    
    public void activate() {
        this.status = UserStatus.ACTIVE;
    }
    
    public void deactivate() {
        this.status = UserStatus.INACTIVE;
    }
    
    public boolean isActive() {
        return this.status == UserStatus.ACTIVE;
    }
}
```

#### `infrastructure/entity/UserEntity.java`
```java
package com.template.app.modules.auth.infrastructure.entity;

import com.template.app._shared.infrastructure.GenericBusinessEntity;
import com.template.app.modules.auth.domain.enums.UserStatus;
import com.template.app.modules.auth.domain.enums.RoleType;
import jakarta.persistence.*;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "m_user")
@EqualsAndHashCode(callSuper = true)
public class UserEntity extends GenericBusinessEntity {
    
    @Column(nullable = false, unique = true)
    private String username;
    
    @Column(nullable = false, unique = true)
    private String email;
    
    @Column(nullable = false)
    private String passwordHash;
    
    @Column(name = "full_name", nullable = false)
    private String fullName;
    
    @Enumerated(EnumType.STRING)
    private UserStatus status;
    
    @Column(name = "phone_number")
    private String phoneNumber;
    
    @Column(name = "profile_image_url")
    private String profileImageUrl;
    
    @Enumerated(EnumType.STRING)
    private RoleType role;
}
```

#### `presentation/controller/AuthController.java`
```java
package com.template.app.modules.auth.presentation.controller;

import com.template.app.modules.auth.application.interfaces.AuthService;
import com.template.app.modules.auth.presentation.dto.request.LoginRequestDTO;
import com.template.app.modules.auth.presentation.dto.response.LoginResponseDTO;
import com.template.app.exception.SystemGlobalException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/auth")
@Tag(name = "Authentication", description = "Authentication and Authorization APIs")
@RequiredArgsConstructor
public class AuthController {

    private final AuthService authService;

    @PostMapping("/login")
    @Operation(summary = "Login to system")
    public ResponseEntity<LoginResponseDTO> login(@Valid @RequestBody LoginRequestDTO request) 
            throws SystemGlobalException {
        return ResponseEntity.ok(authService.login(request));
    }

    @PostMapping("/logout")
    @Operation(summary = "Logout from system")
    public ResponseEntity<Void> logout(@RequestHeader("Authorization") String token) 
            throws SystemGlobalException {
        authService.logout(token);
        return ResponseEntity.ok().build();
    }

    @PostMapping("/refresh")
    @Operation(summary = "Refresh JWT token")
    public ResponseEntity<LoginResponseDTO> refresh(@RequestParam String refreshToken) 
            throws SystemGlobalException {
        return ResponseEntity.ok(authService.refreshToken(refreshToken));
    }
}
```

---
 