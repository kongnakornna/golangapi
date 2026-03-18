# คู่มือมาตรฐานการพัฒนาโครงการ (ฉบับสมบูรณ์)  
**Complete Development Standards Guide**

---

## 1. บทนำ (Introduction)

**ภาษาไทย**  
เอกสารนี้จัดทำขึ้นเพื่อกำหนดกรอบแนวทางและมาตรฐานในการพัฒนาโปรเจกต์ที่ใช้เทคโนโลยีหลากหลาย ครอบคลุมทั้งการออกแบบสถาปัตยกรรม การพัฒนาแบ็กเอนด์และฟรอนต์เอนด์ การจัดการโครงสร้างพื้นฐาน และการปรับใช้ระบบ AI/LLM โดยมีเป้าหมายเพื่อให้ทีมพัฒนาสามารถทำงานร่วมกันได้อย่างมีประสิทธิภาพ มีความสอดคล้องในแนวทางการเขียนโค้ด การควบคุมเวอร์ชัน การตรวจสอบโค้ด และการปรับใช้ระบบสู่สภาพแวดล้อมจริง

**English**  
This document establishes a framework and standards for developing projects that utilize multiple technologies, covering architectural design, backend and frontend development, infrastructure management, and deployment of AI/LLM systems. The goal is to enable development teams to collaborate efficiently, maintain consistency in coding practices, version control, code reviews, and deployment to production environments.

---

## 2. บทนิยาม (Definitions)

| คำศัพท์ (Term) | คำอธิบาย (Description) |
|----------------|-------------------------|
| **RESTful API** | รูปแบบการออกแบบเว็บเซอร์วิสที่ใช้ HTTP Method ในการดำเนินการกับทรัพยากร / A web service design pattern using HTTP methods to operate on resources. |
| **gRPC** | ระบบ RPC (Remote Procedure Call) ประสิทธิภาพสูง พัฒนาโดย Google ใช้ Protocol Buffers / High-performance RPC framework developed by Google, using Protocol Buffers. |
| **Domain-Driven Design (DDD)** | แนวทางการออกแบบซอฟต์แวร์ที่เน้นโมเดลโดเมนเป็นศูนย์กลาง / Software design approach focusing on domain models. |
| **Clean Architecture** | สถาปัตยกรรมซอฟต์แวร์ที่แยกส่วนประกอบออกเป็นชั้นๆ เพื่อให้ง่ายต่อการทดสอบและบำรุงรักษา / Software architecture that separates components into layers for testability and maintainability. |
| **LangChain** | เฟรมเวิร์กสำหรับพัฒนาแอปพลิเคชันที่ใช้ภาษาธรรมชาติและโมเดลภาษา (LLM) / Framework for developing applications with natural language and LLMs. |
| **RAG (Retrieval-Augmented Generation)** | เทคนิคการเพิ่มประสิทธิภาพการสร้างคำตอบโดยดึงข้อมูลจากฐานความรู้ภายนอก / Technique to enhance answer generation by retrieving from external knowledge bases. |
| **IaC (Infrastructure as Code)** | การจัดการโครงสร้างพื้นฐานผ่านไฟล์กำหนดค่า เช่น Terraform / Managing infrastructure through configuration files like Terraform. |
| **CI/CD** | กระบวนการบูรณาการและส่งมอบอย่างต่อเนื่อง / Continuous Integration and Continuous Delivery. |
| **MQTT** | โพรโทคอลส่งข้อความขนาดเล็ก เหมาะสำหรับ IoT / Lightweight messaging protocol for IoT. |
| **InfluxDB** | ฐานข้อมูลอนุกรมเวลาสำหรับจัดเก็บข้อมูล IoT / Time-series database for IoT data. |

---

## 3. สารบัญ (Table of Contents)

1. [บทนำ (Introduction)](#1-บทนำ-introduction)
2. [บทนิยาม (Definitions)](#2-บทนิยาม-definitions)
3. [สารบัญ (Table of Contents)](#3-สารบัญ-table-of-contents)
4. [เวิร์กโฟลว์การพัฒนา (Development Workflow)](#4-เวิร์กโฟลว์การพัฒนา-development-workflow)
5. [ผังการไหลของข้อมูล (Data Flow Architecture)](#5-ผังการไหลของข้อมูล-data-flow-architecture)
6. [แม่แบบโค้ด (Code Templates)](#6-แม่แบบโค้ด-code-templates)
    - 6.1 [แบ็กเอนด์ (Backend)](#61-แบ็กเอนด์-backend)
        - 6.1.1 [Node.js (NestJS + TypeORM)](#611-nodejs-nestjs--typeorm)
        - 6.1.2 [Python (FastAPI + SQLAlchemy)](#612-python-fastapi--sqlalchemy)
        - 6.1.3 [LangChain (FastAPI)](#613-langchain-fastapi)
        - 6.1.4 [Golang (DDD + Clean Architecture)](#614-golang-ddd--clean-architecture)
    - 6.2 [ฟรอนต์เอนด์ (Frontend)](#62-ฟรอนต์เอนด์-frontend)
        - 6.2.1 [React / Next.js](#621-react--nextjs)
    - 6.3 [โครงสร้างพื้นฐาน (Infrastructure as Code)](#63-โครงสร้างพื้นฐาน-infrastructure-as-code)
        - 6.3.1 [Docker](#631-docker)
        - 6.3.2 [Docker Compose](#632-docker-compose)
        - 6.3.3 [Terraform](#633-terraform)
        - 6.3.4 [Kubernetes](#634-kubernetes)
7. [กระบวนการตรวจสอบโค้ด (Code Review Process)](#7-กระบวนการตรวจสอบโค้ด-code-review-process)
8. [กลยุทธ์การจัดการ Git (Git Flow Strategy)](#8-กลยุทธ์การจัดการ-git-git-flow-strategy)
9. [โครงสร้างโปรเจกต์ (Project Layout)](#9-โครงสร้างโปรเจกต์-project-layout)
10. [ระบบ AI/LLM และ IoT (AI/LLM and IoT Systems)](#10-ระบบ-aillm-และ-iot-aillm-and-iot-systems)

---

## 4. เวิร์กโฟลว์การพัฒนา (Development Workflow)

**ภาษาไทย**  
เวิร์กโฟลว์การพัฒนาประกอบด้วยขั้นตอนหลักดังนี้:

1. **วางแผน (Planning)** – วิเคราะห์ความต้องการ กำหนดขอบเขต และออกแบบสถาปัตยกรรม รวบรวม User Stories และสร้าง Task ใน Jira/Trello
2. **พัฒนา (Development)** – เขียนโค้ดตามมาตรฐานที่กำหนด ใช้ Git Flow จัดการ Branch และ commit อย่างสม่ำเสมอ
3. **ทดสอบ (Testing)** – เขียน Unit test และ Integration test ครอบคลุมฟังก์ชันหลัก รัน automated tests ใน CI pipeline
4. **ตรวจสอบโค้ด (Code Review)** – สร้าง Pull Request และให้เพื่อนร่วมทีมตรวจสอบ ใช้หลักเกณฑ์ตามข้อ 7
5. **ปรับใช้ (Deployment)** – ใช้ CI/CD (Jenkins/GitHub Actions) เพื่อ build image และ deploy สู่ environment ต่างๆ (dev/staging/prod) ด้วย Infrastructure as Code
6. **ติดตามผล (Monitoring)** – ตรวจสอบประสิทธิภาพ ข้อผิดพลาด และ logs ผ่านเครื่องมือเช่น Grafana, Prometheus, ELK Stack

**English**  
The development workflow consists of the following main steps:

1. **Planning** – Analyze requirements, define scope, and design architecture. Gather User Stories and create tasks in Jira/Trello.
2. **Development** – Write code following defined standards, use Git Flow for branch management, and commit frequently.
3. **Testing** – Write unit tests and integration tests covering core functionality, run automated tests in CI pipeline.
4. **Code Review** – Create a Pull Request and have team members review it, following guidelines in Section 7.
5. **Deployment** – Use CI/CD (Jenkins/GitHub Actions) to build images and deploy to different environments (dev/staging/prod) with Infrastructure as Code.
6. **Monitoring** – Monitor performance, errors, and logs using tools like Grafana, Prometheus, ELK Stack.

---

## 5. ผังการไหลของข้อมูล (Data Flow Architecture)

**ภาษาไทย**  
แผนภาพการไหลของข้อมูลระหว่างส่วนประกอบต่างๆ ในระบบ:

```mermaid
graph TD
    User[ผู้ใช้] --> FE[Frontend: React/Next.js]
    FE --> API[API Gateway / Load Balancer]
    API --> Auth[Authentication Service]
    API --> UserSvc[User Service - Node.js]
    API --> ProductSvc[Product Service - Python]
    API --> AISvc[AI Service - LangChain/FastAPI]
    API --> IoTSvc[IoT Service - Go]
    
    Auth --> DB[(PostgreSQL)]
    UserSvc --> DB
    UserSvc --> Cache[(Redis)]
    ProductSvc --> DB
    ProductSvc --> Cache
    
    AISvc --> VectorDB[(Vector DB - pgvector)]
    AISvc --> LLM[LLM APIs - OpenAI/DeepSeek]
    
    IoTSvc --> MQTT[MQTT Broker]
    MQTT --> Devices[IoT Devices]
    IoTSvc --> TSDB[(InfluxDB)]
    
    UserSvc --> Queue[Redis Queue]
    Queue --> Worker[Background Workers]
    
    Logs --> ELK[Elasticsearch - Logstash - Kibana]
    Metrics --> Prom[Prometheus]
    Prom --> Grafana[Grafana]
```

**English**  
Data flow diagram between system components:

```mermaid
graph TD
    User[User] --> FE[Frontend: React/Next.js]
    FE --> API[API Gateway / Load Balancer]
    API --> Auth[Authentication Service]
    API --> UserSvc[User Service - Node.js]
    API --> ProductSvc[Product Service - Python]
    API --> AISvc[AI Service - LangChain/FastAPI]
    API --> IoTSvc[IoT Service - Go]
    
    Auth --> DB[(PostgreSQL)]
    UserSvc --> DB
    UserSvc --> Cache[(Redis)]
    ProductSvc --> DB
    ProductSvc --> Cache
    
    AISvc --> VectorDB[(Vector DB - pgvector)]
    AISvc --> LLM[LLM APIs - OpenAI/DeepSeek]
    
    IoTSvc --> MQTT[MQTT Broker]
    MQTT --> Devices[IoT Devices]
    IoTSvc --> TSDB[(InfluxDB)]
    
    UserSvc --> Queue[Redis Queue]
    Queue --> Worker[Background Workers]
    
    Logs --> ELK[Elasticsearch - Logstash - Kibana]
    Metrics --> Prom[Prometheus]
    Prom --> Grafana[Grafana]
```

---

## 6. แม่แบบโค้ด (Code Templates)

### 6.1 แบ็กเอนด์ (Backend)

#### 6.1.1 Node.js (NestJS + TypeORM)

**ภาษาไทย**  
ตัวอย่างโมดูล `User` ใน NestJS พร้อมโค้ดเต็มรูปแบบ

**English**  
Example `User` module in NestJS with full code.

##### โครงสร้างไฟล์ (File Structure)

```
src/
├── modules/
│   └── user/
│       ├── dto/
│       │   ├── create-user.dto.ts
│       │   └── update-user.dto.ts
│       ├── entities/
│       │   └── user.entity.ts
│       ├── repositories/
│       │   └── user.repository.ts
│       ├── services/
│       │   └── user.service.ts
│       ├── controllers/
│       │   └── user.controller.ts
│       └── user.module.ts
├── core/
│   ├── cache/
│   │   └── cache.service.ts
│   ├── logger/
│   │   └── logger.service.ts
│   ├── jwt/
│   │   └── jwt.strategy.ts
│   └── database/
│       └── database.module.ts
└── main.ts
```

##### 1. Entity: `user.entity.ts`

```typescript
import { Entity, Column, PrimaryGeneratedColumn, CreateDateColumn, UpdateDateColumn } from 'typeorm';

@Entity('users')
export class User {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ unique: true })
  email: string;

  @Column()
  password: string; // hashed

  @Column({ nullable: true })
  firstName: string;

  @Column({ nullable: true })
  lastName: string;

  @Column({ default: true })
  isActive: boolean;

  @Column({ type: 'simple-array', default: 'user' })
  roles: string[];

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;
}
```

##### 2. DTO: `create-user.dto.ts`

```typescript
import { IsEmail, IsString, MinLength, IsOptional } from 'class-validator';

export class CreateUserDto {
  @IsEmail()
  email: string;

  @IsString()
  @MinLength(6)
  password: string;

  @IsOptional()
  @IsString()
  firstName?: string;

  @IsOptional()
  @IsString()
  lastName?: string;
}
```

##### 3. Repository: `user.repository.ts`

```typescript
import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from '../entities/user.entity';

@Injectable()
export class UserRepository {
  constructor(
    @InjectRepository(User)
    private readonly repo: Repository<User>,
  ) {}

  async create(userData: Partial<User>): Promise<User> {
    const user = this.repo.create(userData);
    return this.repo.save(user);
  }

  async findById(id: string): Promise<User | null> {
    return this.repo.findOneBy({ id });
  }

  async findByEmail(email: string): Promise<User | null> {
    return this.repo.findOneBy({ email });
  }

  async update(id: string, data: Partial<User>): Promise<void> {
    await this.repo.update(id, data);
  }

  async delete(id: string): Promise<void> {
    await this.repo.delete(id);
  }

  async findAll(): Promise<User[]> {
    return this.repo.find();
  }
}
```

##### 4. Service: `user.service.ts`

```typescript
import { Injectable, ConflictException, NotFoundException } from '@nestjs/common';
import * as bcrypt from 'bcrypt';
import { UserRepository } from '../repositories/user.repository';
import { CreateUserDto } from '../dto/create-user.dto';
import { UpdateUserDto } from '../dto/update-user.dto';
import { User } from '../entities/user.entity';
import { LoggerService } from '../../../core/logger/logger.service';
import { CacheService } from '../../../core/cache/cache.service';

@Injectable()
export class UserService {
  constructor(
    private readonly userRepo: UserRepository,
    private readonly logger: LoggerService,
    private readonly cache: CacheService,
  ) {}

  async create(dto: CreateUserDto): Promise<User> {
    const existing = await this.userRepo.findByEmail(dto.email);
    if (existing) {
      throw new ConflictException('Email already exists');
    }

    const hashedPassword = await bcrypt.hash(dto.password, 10);
    const user = await this.userRepo.create({
      ...dto,
      password: hashedPassword,
    });

    this.logger.log(`User created: ${user.id}`);
    return user;
  }

  async findById(id: string): Promise<User> {
    // Try cache first
    const cached = await this.cache.get<User>(`user:${id}`);
    if (cached) return cached;

    const user = await this.userRepo.findById(id);
    if (!user) {
      throw new NotFoundException('User not found');
    }

    // Store in cache for 5 minutes
    await this.cache.set(`user:${id}`, user, 300);
    return user;
  }

  async update(id: string, dto: UpdateUserDto): Promise<void> {
    await this.findById(id); // ensure exists
    await this.userRepo.update(id, dto);
    await this.cache.del(`user:${id}`); // invalidate cache
  }

  async delete(id: string): Promise<void> {
    await this.findById(id);
    await this.userRepo.delete(id);
    await this.cache.del(`user:${id}`);
  }
}
```

##### 5. Controller: `user.controller.ts`

```typescript
import { Controller, Get, Post, Body, Param, Delete, Put, UseGuards } from '@nestjs/common';
import { UserService } from '../services/user.service';
import { CreateUserDto } from '../dto/create-user.dto';
import { UpdateUserDto } from '../dto/update-user.dto';
import { JwtAuthGuard } from '../../../core/jwt/jwt-auth.guard';
import { RolesGuard } from '../../../core/jwt/roles.guard';
import { Roles } from '../../../core/jwt/roles.decorator';

@Controller('users')
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Post()
  @UseGuards(JwtAuthGuard, RolesGuard)
  @Roles('admin')
  async create(@Body() dto: CreateUserDto) {
    return this.userService.create(dto);
  }

  @Get(':id')
  @UseGuards(JwtAuthGuard)
  async findOne(@Param('id') id: string) {
    return this.userService.findById(id);
  }

  @Put(':id')
  @UseGuards(JwtAuthGuard)
  async update(@Param('id') id: string, @Body() dto: UpdateUserDto) {
    await this.userService.update(id, dto);
    return { message: 'User updated' };
  }

  @Delete(':id')
  @UseGuards(JwtAuthGuard, RolesGuard)
  @Roles('admin')
  async delete(@Param('id') id: string) {
    await this.userService.delete(id);
    return { message: 'User deleted' };
  }
}
```

##### 6. Module: `user.module.ts`

```typescript
import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { UserController } from './controllers/user.controller';
import { UserService } from './services/user.service';
import { UserRepository } from './repositories/user.repository';
import { User } from './entities/user.entity';
import { LoggerModule } from '../../core/logger/logger.module';
import { CacheModule } from '../../core/cache/cache.module';

@Module({
  imports: [
    TypeOrmModule.forFeature([User]),
    LoggerModule,
    CacheModule,
  ],
  controllers: [UserController],
  providers: [UserService, UserRepository],
  exports: [UserService],
})
export class UserModule {}
```

---

#### 6.1.2 Python (FastAPI + SQLAlchemy)

**ภาษาไทย**  
ตัวอย่างโมดูล `User` ใน FastAPI พร้อมโค้ดเต็มรูปแบบ

**English**  
Example `User` module in FastAPI with full code.

##### โครงสร้างไฟล์ (File Structure)

```
app/
├── modules/
│   └── user/
│       ├── models.py
│       ├── schemas.py
│       ├── repositories.py
│       ├── services.py
│       ├── controllers.py
│       └── dependencies.py
├── core/
│   ├── database.py
│   ├── cache.py
│   ├── jwt.py
│   └── logger.py
├── main.py
└── config.py
```

##### 1. Model: `models.py`

```python
from sqlalchemy import Column, String, Boolean, DateTime
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.sql import func
import uuid
from app.core.database import Base

class User(Base):
    __tablename__ = "users"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    email = Column(String, unique=True, nullable=False)
    password = Column(String, nullable=False)  # hashed
    first_name = Column(String)
    last_name = Column(String)
    is_active = Column(Boolean, default=True)
    roles = Column(String, default="user")  # comma-separated or JSON
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
```

##### 2. Schema: `schemas.py`

```python
from pydantic import BaseModel, EmailStr, Field
from uuid import UUID
from datetime import datetime
from typing import Optional

class UserBase(BaseModel):
    email: EmailStr
    first_name: Optional[str] = None
    last_name: Optional[str] = None

class UserCreate(UserBase):
    password: str = Field(..., min_length=6)

class UserUpdate(BaseModel):
    first_name: Optional[str] = None
    last_name: Optional[str] = None
    password: Optional[str] = None

class UserResponse(UserBase):
    id: UUID
    is_active: bool
    roles: str
    created_at: datetime
    updated_at: Optional[datetime]

    class Config:
        from_attributes = True
```

##### 3. Repository: `repositories.py`

```python
from sqlalchemy.orm import Session
from app.modules.user import models, schemas
from app.core.database import get_db
from fastapi import Depends

class UserRepository:
    def __init__(self, db: Session):
        self.db = db

    def create(self, user: schemas.UserCreate, hashed_password: str) -> models.User:
        db_user = models.User(
            email=user.email,
            password=hashed_password,
            first_name=user.first_name,
            last_name=user.last_name,
        )
        self.db.add(db_user)
        self.db.commit()
        self.db.refresh(db_user)
        return db_user

    def get_by_email(self, email: str) -> models.User | None:
        return self.db.query(models.User).filter(models.User.email == email).first()

    def get_by_id(self, user_id: UUID) -> models.User | None:
        return self.db.query(models.User).filter(models.User.id == user_id).first()

    def update(self, user_id: UUID, data: dict) -> models.User | None:
        self.db.query(models.User).filter(models.User.id == user_id).update(data)
        self.db.commit()
        return self.get_by_id(user_id)

    def delete(self, user_id: UUID) -> None:
        self.db.query(models.User).filter(models.User.id == user_id).delete()
        self.db.commit()

def get_user_repository(db: Session = Depends(get_db)) -> UserRepository:
    return UserRepository(db)
```

##### 4. Service: `services.py`

```python
from fastapi import HTTPException, status
from app.modules.user import schemas, repositories
import bcrypt
from app.core.cache import redis_client
import json

class UserService:
    def __init__(self, repo: repositories.UserRepository):
        self.repo = repo

    def _hash_password(self, password: str) -> str:
        salt = bcrypt.gensalt()
        return bcrypt.hashpw(password.encode('utf-8'), salt).decode('utf-8')

    def _verify_password(self, plain: str, hashed: str) -> bool:
        return bcrypt.checkpw(plain.encode('utf-8'), hashed.encode('utf-8'))

    async def create_user(self, user_data: schemas.UserCreate) -> schemas.UserResponse:
        existing = self.repo.get_by_email(user_data.email)
        if existing:
            raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="Email already registered")
        hashed = self._hash_password(user_data.password)
        db_user = self.repo.create(user_data, hashed)
        return schemas.UserResponse.model_validate(db_user)

    async def get_user(self, user_id: UUID) -> schemas.UserResponse:
        # check cache
        cached = await redis_client.get(f"user:{user_id}")
        if cached:
            return schemas.UserResponse.model_validate(json.loads(cached))

        db_user = self.repo.get_by_id(user_id)
        if not db_user:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")

        # cache for 5 minutes
        await redis_client.setex(f"user:{user_id}", 300, json.dumps(schemas.UserResponse.model_validate(db_user).model_dump()))
        return schemas.UserResponse.model_validate(db_user)

    async def update_user(self, user_id: UUID, update_data: schemas.UserUpdate) -> schemas.UserResponse:
        db_user = self.repo.get_by_id(user_id)
        if not db_user:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")

        update_dict = update_data.model_dump(exclude_unset=True)
        if 'password' in update_dict:
            update_dict['password'] = self._hash_password(update_dict['password'])

        updated = self.repo.update(user_id, update_dict)
        await redis_client.delete(f"user:{user_id}")
        return schemas.UserResponse.model_validate(updated)

    async def delete_user(self, user_id: UUID) -> None:
        db_user = self.repo.get_by_id(user_id)
        if not db_user:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
        self.repo.delete(user_id)
        await redis_client.delete(f"user:{user_id}")
```

##### 5. Controller: `controllers.py`

```python
from fastapi import APIRouter, Depends, HTTPException, status
from uuid import UUID
from app.modules.user import schemas, services
from app.modules.user.repositories import get_user_repository, UserRepository
from app.core.jwt import get_current_user, require_roles

router = APIRouter(prefix="/users", tags=["users"])

@router.post("/", response_model=schemas.UserResponse)
async def create_user(
    user_data: schemas.UserCreate,
    repo: UserRepository = Depends(get_user_repository),
    current_user: dict = Depends(require_roles(["admin"]))
):
    service = services.UserService(repo)
    return await service.create_user(user_data)

@router.get("/{user_id}", response_model=schemas.UserResponse)
async def get_user(
    user_id: UUID,
    repo: UserRepository = Depends(get_user_repository),
    current_user: dict = Depends(get_current_user)
):
    service = services.UserService(repo)
    return await service.get_user(user_id)

@router.put("/{user_id}", response_model=schemas.UserResponse)
async def update_user(
    user_id: UUID,
    update_data: schemas.UserUpdate,
    repo: UserRepository = Depends(get_user_repository),
    current_user: dict = Depends(get_current_user)
):
    service = services.UserService(repo)
    return await service.update_user(user_id, update_data)

@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_user(
    user_id: UUID,
    repo: UserRepository = Depends(get_user_repository),
    current_user: dict = Depends(require_roles(["admin"]))
):
    service = services.UserService(repo)
    await service.delete_user(user_id)
```

##### 6. Main: `main.py`

```python
from fastapi import FastAPI
from app.modules.user.controllers import router as user_router
from app.core.database import engine, Base
from app.core.logger import setup_logging
from app.core.cache import redis_client
from app.core.jwt import setup_jwt

# Create tables
Base.metadata.create_all(bind=engine)

app = FastAPI(title="My API", version="1.0.0")

# Setup logging
setup_logging()

# Include routers
app.include_router(user_router)

@app.on_event("startup")
async def startup():
    await redis_client.initialize()

@app.on_event("shutdown")
async def shutdown():
    await redis_client.close()

@app.get("/health")
async def health():
    return {"status": "ok"}
```

---

#### 6.1.3 LangChain (FastAPI)

**ภาษาไทย**  
ตัวอย่างโมดูล RAG (Retrieval-Augmented Generation) โดยใช้ LangChain ร่วมกับ FastAPI

**English**  
Example RAG module using LangChain with FastAPI.

##### โครงสร้างไฟล์ (File Structure)

```
app/
├── ai/
│   ├── llm.py
│   ├── chains/
│   │   └── rag_chain.py
│   ├── vector_store/
│   │   └── pgvector.py
│   └── embeddings.py
├── modules/
│   └── rag/
│       ├── controllers.py
│       ├── services.py
│       └── schemas.py
├── core/
│   └── database.py
└── main.py
```

##### 1. LLM Configuration: `ai/llm.py`

```python
from langchain_openai import ChatOpenAI
from langchain_anthropic import ChatAnthropic
from langchain_deepseek import ChatDeepSeek
import os

def get_llm(provider: str = "openai", temperature: float = 0.7):
    if provider == "openai":
        return ChatOpenAI(
            model="gpt-4",
            temperature=temperature,
            api_key=os.getenv("OPENAI_API_KEY")
        )
    elif provider == "anthropic":
        return ChatAnthropic(
            model="claude-3-opus-20240229",
            temperature=temperature,
            api_key=os.getenv("ANTHROPIC_API_KEY")
        )
    elif provider == "deepseek":
        return ChatDeepSeek(
            model="deepseek-chat",
            temperature=temperature,
            api_key=os.getenv("DEEPSEEK_API_KEY")
        )
    else:
        raise ValueError(f"Unknown provider: {provider}")
```

##### 2. Vector Store: `ai/vector_store/pgvector.py`

```python
from langchain_community.vectorstores import PGVector
from langchain_community.embeddings import OpenAIEmbeddings
import os

CONNECTION_STRING = PGVector.connection_string_from_db_params(
    driver="psycopg2",
    host=os.getenv("DB_HOST", "localhost"),
    port=int(os.getenv("DB_PORT", "5432")),
    database=os.getenv("DB_NAME", "vectordb"),
    user=os.getenv("DB_USER", "user"),
    password=os.getenv("DB_PASSWORD", "password"),
)

def get_vector_store(collection_name: str = "documents"):
    embeddings = OpenAIEmbeddings()
    return PGVector(
        connection_string=CONNECTION_STRING,
        embedding_function=embeddings,
        collection_name=collection_name,
    )
```

##### 3. RAG Chain: `ai/chains/rag_chain.py`

```python
from langchain.chains import create_retrieval_chain
from langchain.chains.combine_documents import create_stuff_documents_chain
from langchain_core.prompts import ChatPromptTemplate
from app.ai.llm import get_llm
from app.ai.vector_store.pgvector import get_vector_store

def create_rag_chain():
    llm = get_llm()
    vector_store = get_vector_store()
    retriever = vector_store.as_retriever(search_kwargs={"k": 5})

    prompt = ChatPromptTemplate.from_messages([
        ("system", "You are a helpful assistant. Use the following context to answer the question. If you don't know, say you don't know."),
        ("system", "Context: {context}"),
        ("human", "{input}")
    ])

    document_chain = create_stuff_documents_chain(llm, prompt)
    return create_retrieval_chain(retriever, document_chain)
```

##### 4. Service: `modules/rag/services.py`

```python
from app.ai.chains.rag_chain import create_rag_chain

class RAGService:
    def __init__(self):
        self.chain = create_rag_chain()

    async def ask(self, question: str) -> str:
        # Run chain
        result = self.chain.invoke({"input": question})
        return result["answer"]
```

##### 5. Schema: `modules/rag/schemas.py`

```python
from pydantic import BaseModel

class QuestionRequest(BaseModel):
    question: str

class AnswerResponse(BaseModel):
    answer: str
```

##### 6. Controller: `modules/rag/controllers.py`

```python
from fastapi import APIRouter, HTTPException
from app.modules.rag.schemas import QuestionRequest, AnswerResponse
from app.modules.rag.services import RAGService

router = APIRouter(prefix="/rag", tags=["rag"])
rag_service = RAGService()

@router.post("/ask", response_model=AnswerResponse)
async def ask_question(request: QuestionRequest):
    try:
        answer = await rag_service.ask(request.question)
        return AnswerResponse(answer=answer)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
```

##### 7. Main: รวม router

```python
from fastapi import FastAPI
from app.modules.rag.controllers import router as rag_router

app = FastAPI()
app.include_router(rag_router)
```

---

#### 6.1.4 Golang (DDD + Clean Architecture)

**ภาษาไทย**  
ตัวอย่างโมดูล `User` ใน Go ตามหลัก DDD และ Clean Architecture พร้อมโค้ดเต็มรูปแบบ

**English**  
Example `User` module in Go following DDD and Clean Architecture with full code.

##### โครงสร้างไฟล์ (File Structure)

```
internal/
├── core/
│   └── user/
│       ├── entity/
│       │   └── user.go
│       ├── repository/
│       │   └── user_repository.go (interface)
│       ├── service/
│       │   └── user_service.go
│       └── handler/
│           └── user_handler.go
├── platform/
│   ├── db/
│   │   └── postgres.go
│   ├── cache/
│   │   └── redis.go
│   └── logger/
│       └── logger.go
└── transport/
    └── http/
        ├── middleware/
        │   └── auth.go
        └── router.go
cmd/
└── app/
    └── main.go
pkg/
└── utils/
    └── hash.go
```

##### 1. Entity: `internal/core/user/entity/user.go`

```go
package entity

import (
    "time"
    "google/uuid"
)

type User struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-"` // exclude from JSON
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    IsActive  bool      `json:"is_active" gorm:"default:true"`
    Roles     string    `json:"roles" gorm:"default:'user'"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

##### 2. Repository Interface: `internal/core/user/repository/user_repository.go`

```go
package repository

import (
    "context"
    "google/uuid"
    "yourproject/internal/core/user/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context) ([]entity.User, error)
}
```

##### 3. Service: `internal/core/user/service/user_service.go`

```go
package service

import (
    "context"
    "errors"
    "time"
    "google/uuid"
    "yourproject/internal/core/user/entity"
    "yourproject/internal/core/user/repository"
    "yourproject/pkg/utils"
    "yourproject/internal/platform/cache"
    "yourproject/internal/platform/logger"
)

type UserService struct {
    repo  repository.UserRepository
    cache cache.Cache
    log   logger.Logger
}

func NewUserService(repo repository.UserRepository, cache cache.Cache, log logger.Logger) *UserService {
    return &UserService{repo: repo, cache: cache, log: log}
}

func (s *UserService) Create(ctx context.Context, email, password, firstName, lastName string) (*entity.User, error) {
    // Check if user exists
    existing, _ := s.repo.GetByEmail(ctx, email)
    if existing != nil {
        return nil, errors.New("email already exists")
    }

    hashed, err := utils.HashPassword(password)
    if err != nil {
        return nil, err
    }

    user := &entity.User{
        ID:        uuid.New(),
        Email:     email,
        Password:  hashed,
        FirstName: firstName,
        LastName:  lastName,
        IsActive:  true,
        Roles:     "user",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }

    s.log.Info("User created", "id", user.ID)
    return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    // Try cache
    var user entity.User
    found, err := s.cache.Get(ctx, "user:"+id.String(), &user)
    if err == nil && found {
        return &user, nil
    }

    userPtr, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if userPtr == nil {
        return nil, errors.New("user not found")
    }

    // Set cache
    _ = s.cache.Set(ctx, "user:"+id.String(), userPtr, 5*time.Minute)
    return userPtr, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, firstName, lastName *string, password *string) error {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil || user == nil {
        return errors.New("user not found")
    }

    if firstName != nil {
        user.FirstName = *firstName
    }
    if lastName != nil {
        user.LastName = *lastName
    }
    if password != nil {
        hashed, err := utils.HashPassword(*password)
        if err != nil {
            return err
        }
        user.Password = hashed
    }
    user.UpdatedAt = time.Now()

    if err := s.repo.Update(ctx, user); err != nil {
        return err
    }

    // Invalidate cache
    _ = s.cache.Del(ctx, "user:"+id.String())
    return nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
    if err := s.repo.Delete(ctx, id); err != nil {
        return err
    }
    _ = s.cache.Del(ctx, "user:"+id.String())
    return nil
}
```

##### 4. Handler: `internal/core/user/handler/user_handler.go`

```go
package handler

import (
    "net/http"
    "gin-gonic/gin"
    "google/uuid"
    "yourproject/internal/core/user/service"
    "yourproject/internal/transport/http/middleware"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

// Request/Response DTOs
type CreateUserRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=6"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

type UpdateUserRequest struct {
    FirstName *string `json:"first_name"`
    LastName  *string `json:"last_name"`
    Password  *string `json:"password" binding:"omitempty,min=6"`
}

type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    IsActive  bool      `json:"is_active"`
    Roles     string    `json:"roles"`
    CreatedAt string    `json:"created_at"`
}

func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userService.Create(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, toResponse(user))
}

func (h *UserHandler) Get(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    user, err := h.userService.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, toResponse(user))
}

func (h *UserHandler) Update(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.userService.Update(c.Request.Context(), id, req.FirstName, req.LastName, req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    if err := h.userService.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}

func toResponse(u *entity.User) UserResponse {
    return UserResponse{
        ID:        u.ID,
        Email:     u.Email,
        FirstName: u.FirstName,
        LastName:  u.LastName,
        IsActive:  u.IsActive,
        Roles:     u.Roles,
        CreatedAt: u.CreatedAt.Format(time.RFC3339),
    }
}

// Register routes
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
    users := r.Group("/users")
    users.Use(authMiddleware)
    {
        users.POST("/", middleware.RequireRoles("admin"), h.Create)
        users.GET("/:id", h.Get)
        users.PUT("/:id", h.Update)
        users.DELETE("/:id", middleware.RequireRoles("admin"), h.Delete)
    }
}
```

##### 5. Main: `cmd/app/main.go` (รวม dependencies)

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "yourproject/internal/core/user/entity"
    "yourproject/internal/core/user/handler"
    "yourproject/internal/core/user/repository"
    "yourproject/internal/core/user/service"
    "yourproject/internal/platform/cache"
    "yourproject/internal/platform/logger"
    "yourproject/internal/transport/http/middleware"
)

func main() {
    // Logger
    log := logger.New()

    // Database
    dsn := "host=localhost user=user password=password dbname=mydb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database", "error", err)
    }
    db.AutoMigrate(&entity.User{})

    // Redis
    redisClient := cache.NewRedisClient("localhost:6379", "", 0)
    cache := cache.NewRedisCache(redisClient)

    // Repositories
    userRepo := repository.NewGormUserRepository(db)

    // Services
    userService := service.NewUserService(userRepo, cache, log)

    // Handlers
    userHandler := handler.NewUserHandler(userService)

    // Gin engine
    r := gin.Default()
    r.Use(middleware.Logger(log))
    r.Use(middleware.Recovery())

    // Auth middleware
    authMiddleware := middleware.JWTAuth(os.Getenv("JWT_SECRET"))

    // Register routes
    api := r.Group("/api/v1")
    userHandler.RegisterRoutes(api, authMiddleware)

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Start server
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("listen error", "error", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("server forced to shutdown", "error", err)
    }
}
```

##### 6. Repository Implementation (GORM): `internal/platform/db/gorm_user_repository.go`

```go
package repository

import (
    "context"
    "google/uuid"
    "gorm.io/gorm"
    "yourproject/internal/core/user/entity"
)

type gormUserRepository struct {
    db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *gormUserRepository {
    return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(ctx context.Context, user *entity.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
    if err == gorm.ErrRecordNotFound {
        return nil, nil
    }
    return &user, err
}

func (r *gormUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err == gorm.ErrRecordNotFound {
        return nil, nil
    }
    return &user, err
}

func (r *gormUserRepository) Update(ctx context.Context, user *entity.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *gormUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *gormUserRepository) List(ctx context.Context) ([]entity.User, error) {
    var users []entity.User
    err := r.db.WithContext(ctx).Find(&users).Error
    return users, err
}
```

---

### 6.2 ฟรอนต์เอนด์ (Frontend)

#### 6.2.1 React / Next.js

**ภาษาไทย**  
ตัวอย่างโมดูล `User` ใน Next.js (App Router) พร้อมการเรียก API

**English**  
Example `User` module in Next.js (App Router) with API calls.

##### โครงสร้างไฟล์ (File Structure)

```
src/
├── app/
│   ├── (dashboard)/
│   │   └── users/
│   │       ├── page.tsx
│   │       └── [id]/
│   │           └── page.tsx
├── modules/
│   └── user/
│       ├── components/
│       │   ├── UserList.tsx
│       │   ├── UserForm.tsx
│       │   └── UserCard.tsx
│       ├── services/
│       │   └── userService.ts
│       ├── hooks/
│       │   └── useUsers.ts
│       └── types/
│           └── user.types.ts
├── lib/
│   ├── api/
│   │   └── axios.ts
│   └── cache/
│       └── react-query.ts
```

##### 1. Types: `modules/user/types/user.types.ts`

```typescript
export interface User {
  id: string;
  email: string;
  firstName?: string;
  lastName?: string;
  isActive: boolean;
  roles: string;
  createdAt: string;
}

export interface CreateUserDto {
  email: string;
  password: string;
  firstName?: string;
  lastName?: string;
}

export interface UpdateUserDto {
  firstName?: string;
  lastName?: string;
  password?: string;
}
```

##### 2. Service: `modules/user/services/userService.ts`

```typescript
import api from '@/lib/api/axios';
import { User, CreateUserDto, UpdateUserDto } from '../types/user.types';

export const userService = {
  async getAll(): Promise<User[]> {
    const response = await api.get<User[]>('/users');
    return response.data;
  },

  async getById(id: string): Promise<User> {
    const response = await api.get<User>(`/users/${id}`);
    return response.data;
  },

  async create(data: CreateUserDto): Promise<User> {
    const response = await api.post<User>('/users', data);
    return response.data;
  },

  async update(id: string, data: UpdateUserDto): Promise<User> {
    const response = await api.put<User>(`/users/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/users/${id}`);
  },
};
```

##### 3. React Query Hooks: `modules/user/hooks/useUsers.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { userService } from '../services/userService';
import { CreateUserDto, UpdateUserDto } from '../types/user.types';

export const useUsers = () => {
  return useQuery({
    queryKey: ['users'],
    queryFn: userService.getAll,
  });
};

export const useUser = (id: string) => {
  return useQuery({
    queryKey: ['user', id],
    queryFn: () => userService.getById(id),
    enabled: !!id,
  });
};

export const useCreateUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateUserDto) => userService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserDto }) =>
      userService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      queryClient.invalidateQueries({ queryKey: ['user', variables.id] });
    },
  });
};

export const useDeleteUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => userService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
};
```

##### 4. Component: `modules/user/components/UserList.tsx`

```tsx
import React from 'react';
import { useUsers } from '../hooks/useUsers';
import { UserCard } from './UserCard';

export const UserList: React.FC = () => {
  const { data: users, isLoading, error } = useUsers();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading users</div>;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {users?.map((user) => (
        <UserCard key={user.id} user={user} />
      ))}
    </div>
  );
};
```

##### 5. Page: `app/(dashboard)/users/page.tsx`

```tsx
'use client';

import { UserList } from '@/modules/user/components/UserList';
import { Button } from '@/components/ui/button';
import Link from 'next/link';

export default function UsersPage() {
  return (
    <div className="container mx-auto py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Users</h1>
        <Link href="/users/new">
          <Button>Create User</Button>
        </Link>
      </div>
      <UserList />
    </div>
  );
}
```

##### 6. Form Component: `modules/user/components/UserForm.tsx`

```tsx
import React from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { CreateUserDto } from '../types/user.types';

const schema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
  firstName: z.string().optional(),
  lastName: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

interface UserFormProps {
  onSubmit: (data: CreateUserDto) => void;
  defaultValues?: Partial<CreateUserDto>;
}

export const UserForm: React.FC<UserFormProps> = ({ onSubmit, defaultValues }) => {
  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues,
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <label className="block text-sm font-medium">Email</label>
        <input type="email" {...register('email')} className="mt-1 block w-full border rounded p-2" />
        {errors.email && <p className="text-red-500 text-sm">{errors.email.message}</p>}
      </div>
      <div>
        <label className="block text-sm font-medium">Password</label>
        <input type="password" {...register('password')} className="mt-1 block w-full border rounded p-2" />
        {errors.password && <p className="text-red-500 text-sm">{errors.password.message}</p>}
      </div>
      <div>
        <label className="block text-sm font-medium">First Name</label>
        <input {...register('firstName')} className="mt-1 block w-full border rounded p-2" />
      </div>
      <div>
        <label className="block text-sm font-medium">Last Name</label>
        <input {...register('lastName')} className="mt-1 block w-full border rounded p-2" />
      </div>
      <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded">
        Submit
      </button>
    </form>
  );
};
```

---

### 6.3 โครงสร้างพื้นฐาน (Infrastructure as Code)

#### 6.3.1 Docker

**ภาษาไทย**  
ตัวอย่าง Dockerfile สำหรับแต่ละภาษา (ดังแสดงในข้อ 6.1.1-6.1.4 และ 6.2.1) สามารถใช้ได้ตามความเหมาะสม

**English**  
Example Dockerfiles for each language (as shown in sections 6.1.1-6.1.4 and 6.2.1) can be used accordingly.

#### 6.3.2 Docker Compose

**ภาษาไทย**  
ตัวอย่าง `docker-compose.yml` สำหรับรันทั้งระบบ (PostgreSQL, Redis, Backend services, Frontend)

**English**  
Example `docker-compose.yml` to run the entire system (PostgreSQL, Redis, Backend services, Frontend).

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: postgres
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    networks:
      - app-network

  api-node:
    build: ./backend/nodejs
    container_name: api-node
    ports:
      - "3000:3000"
    depends_on:
      - postgres
      - redis
    environment:
      DATABASE_URL: postgresql://user:password@postgres:5432/mydb
      REDIS_URL: redis://redis:6379
    networks:
      - app-network

  api-python:
    build: ./backend/python
    container_name: api-python
    ports:
      - "8000:8000"
    depends_on:
      - postgres
      - redis
    environment:
      DATABASE_URL: postgresql://user:password@postgres:5432/mydb
      REDIS_URL: redis://redis:6379
    networks:
      - app-network

  api-go:
    build: ./backend/go
    container_name: api-go
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    environment:
      DB_HOST: postgres
      DB_USER: user
      DB_PASSWORD: password
      DB_NAME: mydb
      REDIS_ADDR: redis:6379
    networks:
      - app-network

  frontend:
    build: ./frontend/nextjs
    container_name: frontend
    ports:
      - "3001:3000"
    depends_on:
      - api-node
      - api-python
      - api-go
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:3000
    networks:
      - app-network

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

#### 6.3.3 Terraform

**ภาษาไทย**  
ตัวอย่าง Terraform สำหรับสร้างโครงสร้างพื้นฐานบน AWS (VPC, EC2, RDS, Redis)

**English**  
Example Terraform to create infrastructure on AWS (VPC, EC2, RDS, Redis).

```hcl
# main.tf
provider "aws" {
  region = "ap-southeast-1"
}

# VPC
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = { Name = "myapp-vpc" }
}

# Subnets
resource "aws_subnet" "public" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "ap-southeast-1a"
  map_public_ip_on_launch = true
  tags = { Name = "public-subnet" }
}

# Internet Gateway
resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
  tags = { Name = "myapp-igw" }
}

# Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }
  tags = { Name = "public-rt" }
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

# Security Group for EC2
resource "aws_security_group" "ec2" {
  name        = "ec2-sg"
  description = "Allow HTTP, HTTPS, SSH"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# EC2 Instance
resource "aws_instance" "app" {
  ami           = "ami-0c55b159cbfafe1f0" # Amazon Linux 2
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.public.id
  vpc_security_group_ids = [aws_security_group.ec2.id]
  key_name               = "mykey" # Create key pair beforehand
  user_data = <<-EOF
              #!/bin/bash
              yum update -y
              yum install -y docker
              service docker start
              usermod -a -G docker ec2-user
              curl -L "https://docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
              chmod +x /usr/local/bin/docker-compose
              EOF
  tags = { Name = "app-server" }
}

# RDS PostgreSQL
resource "aws_db_instance" "postgres" {
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "postgres"
  engine_version       = "15.3"
  instance_class       = "db.t3.micro"
  db_name              = "mydb"
  username             = "admin"
  password             = "password123"
  parameter_group_name = "default.postgres15"
  skip_final_snapshot  = true
  publicly_accessible  = false
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
}

# ElastiCache Redis
resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "redis-cluster"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.main.name
  security_group_ids   = [aws_security_group.redis.id]
}
```

#### 6.3.4 Kubernetes

**ภาษาไทย**  
ตัวอย่าง Kubernetes manifests สำหรับ Deployments, Services, Ingress

**English**  
Example Kubernetes manifests for Deployments, Services, Ingress.

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: myapp
---
# deployment-api-node.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-node
  namespace: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-node
  template:
    metadata:
      labels:
        app: api-node
    spec:
      containers:
      - name: api-node
        image: myregistry/api-node:latest
        ports:
        - containerPort: 3000
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: redis-secret
              key: url
---
# service-api-node.yaml
apiVersion: v1
kind: Service
metadata:
  name: api-node-service
  namespace: myapp
spec:
  selector:
    app: api-node
  ports:
    - protocol: TCP
      port: 80
      targetPort: 3000
  type: ClusterIP
---
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp-ingress
  namespace: myapp
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: api.myapp.com
    http:
      paths:
      - path: /node
        pathType: Prefix
        backend:
          service:
            name: api-node-service
            port:
              number: 80
      - path: /python
        pathType: Prefix
        backend:
          service:
            name: api-python-service
            port:
              number: 80
      - path: /go
        pathType: Prefix
        backend:
          service:
            name: api-go-service
            port:
              number: 80
```

---

## 7. กระบวนการตรวจสอบโค้ด (Code Review Process)

**ภาษาไทย**  
หลักเกณฑ์การตรวจสอบโค้ด:

- **Pull Request (PR)** ต้องถูกสร้างสำหรับทุก feature/fix ก่อน merge
- PR ต้องผ่านการทดสอบอัตโนมัติ (CI) ก่อน review
- ตรวจสอบความถูกต้องของ logic, การจัดการ error, security, performance
- ตรวจสอบการตั้งชื่อ, โครงสร้าง, การจัดรูปแบบตามมาตรฐานของภาษา
- ตรวจสอบว่ามีการเขียน test ครอบคลุม (unit test, integration test)
- อย่างน้อย 1 คนในทีมต้อง approve ก่อน merge (สำหรับโปรเจกต์สำคัญควรมี 2 คน)
- ใช้เครื่องมือเช่น SonarQube ช่วยวิเคราะห์คุณภาพโค้ด
- ห้าม merge PR ที่มี unresolved comments

**Checklist การตรวจสอบ:**
- [ ] โค้ดทำงานตาม requirements
- [ ] ไม่มี security vulnerability (SQL injection, XSS, etc.)
- [ ] มี error handling ที่เหมาะสม
- [ ] มี logging ที่เพียงพอ
- [ ] performance ดี (ไม่ query ซ้ำซ้อน, ใช้ index)
- [ ] มีการเขียน test และ test ผ่าน
- [ ] โค้ดอ่านง่าย มี comment เมื่อจำเป็น
- [ ] ตั้งชื่อตัวแปร ฟังก์ชัน สื่อความหมาย
- [ ] ไม่มีโค้ดที่ถูก comment ทิ้งไว้

**English**  
Code review guidelines:

- **Pull Request (PR)** must be created for every feature/fix before merging.
- PR must pass automated tests (CI) before review.
- Verify logic correctness, error handling, security, performance.
- Check naming, structure, formatting according to language standards.
- Ensure tests are written and cover the changes (unit test, integration test).
- At least 1 team member must approve before merging (for critical projects, 2 approvals).
- Use tools like SonarQube to analyze code quality.
- Do not merge PR with unresolved comments.

**Review Checklist:**
- [ ] Code works as per requirements
- [ ] No security vulnerabilities (SQL injection, XSS, etc.)
- [ ] Proper error handling
- [ ] Sufficient logging
- [ ] Good performance (no redundant queries, indexes used)
- [ ] Tests written and passing
- [ ] Code readable, comments where necessary
- [ ] Meaningful variable/function names
- [ ] No commented-out code

---

## 8. กลยุทธ์การจัดการ Git (Git Flow Strategy)

**ภาษาไทย**  
ใช้ Git Flow เป็นหลัก:

- **main** – สาขาหลักสำหรับ production release ห้าม commit ตรง
- **develop** – สาขาสำหรับการพัฒนารวม (integration branch)
- **feature/*** – สาขาสำหรับพัฒนา feature ใหม่ (แตกจาก develop) ตั้งชื่อเช่น `feature/user-login`
- **release/*** – สาขาสำหรับเตรียม release (แก้ไขบั๊กสุดท้าย) ตั้งชื่อเช่น `release/v1.2.0`
- **hotfix/*** – สาขาสำหรับแก้ไขด่วนใน production (แตกจาก main) ตั้งชื่อเช่น `hotfix/critical-bug`

เมื่อเสร็จสิ้น feature ให้ merge กลับไปที่ develop และลบ feature branch  
เมื่อพร้อม release ให้สร้าง release branch จาก develop, ทดสอบและแก้ไขบั๊กเล็กน้อย, แล้ว merge ไปที่ main และ develop  
เมื่อมี hotfix ให้สร้างจาก main, แก้ไข, แล้ว merge ไปที่ main และ develop

**Commit Message Convention:** ใช้ Conventional Commits เช่น  
`feat: add user login`  
`fix: resolve password reset issue`  
`docs: update README`

**English**  
Use Git Flow as the main branching model:

- **main** – production releases, no direct commits.
- **develop** – integration branch for development.
- **feature/*** – new features (branched from develop), e.g., `feature/user-login`.
- **release/*** – release preparation (bug fixes), e.g., `release/v1.2.0`.
- **hotfix/*** – urgent fixes for production (branched from main), e.g., `hotfix/critical-bug`.

After completing a feature, merge back to develop and delete the feature branch.  
When ready for release, create a release branch from develop, test and fix minor bugs, then merge to main and develop.  
For hotfix, branch from main, fix, then merge to main and develop.

**Commit Message Convention:** Use Conventional Commits, e.g.,  
`feat: add user login`  
`fix: resolve password reset issue`  
`docs: update README`

---

## 9. โครงสร้างโปรเจกต์ (Project Layout)

**ภาษาไทย**  
สำหรับโปรเจกต์ที่มีหลายภาษาและหลายบริการ ควรใช้โครงสร้างแบบ monorepo เพื่อจัดการโค้ดทั้งหมดใน repository เดียว ตัวอย่างโครงสร้าง:

```
/
├── backend/
│   ├── nodejs/
│   │   ├── src/
│   │   ├── Dockerfile
│   │   └── package.json
│   ├── python/
│   │   ├── app/
│   │   ├── Dockerfile
│   │   └── requirements.txt
│   ├── go/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   └── langchain/
│       ├── app/
│       ├── Dockerfile
│       └── requirements.txt
├── frontend/
│   └── nextjs/
│       ├── src/
│       ├── Dockerfile
│       └── package.json
├── infrastructure/
│   ├── docker-compose/
│   │   └── docker-compose.yml
│   ├── terraform/
│   │   └── aws/
│   └── k8s/
│       ├── deployments/
│       └── services/
├── docs/
├── scripts/
└── README.md
```

**English**  
For projects with multiple languages and services, use a monorepo structure to manage all code in a single repository. Example structure:

```
/
├── backend/
│   ├── nodejs/
│   │   ├── src/
│   │   ├── Dockerfile
│   │   └── package.json
│   ├── python/
│   │   ├── app/
│   │   ├── Dockerfile
│   │   └── requirements.txt
│   ├── go/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   └── langchain/
│       ├── app/
│       ├── Dockerfile
│       └── requirements.txt
├── frontend/
│   └── nextjs/
│       ├── src/
│       ├── Dockerfile
│       └── package.json
├── infrastructure/
│   ├── docker-compose/
│   │   └── docker-compose.yml
│   ├── terraform/
│   │   └── aws/
│   └── k8s/
│       ├── deployments/
│       └── services/
├── docs/
├── scripts/
└── README.md
```

---

## 10. ระบบ AI/LLM และ IoT (AI/LLM and IoT Systems)

**ภาษาไทย**  
**AI/LLM Systems**  
- ใช้ LangChain ในการสร้าง pipelines สำหรับ RAG, agent, chain
- จัดเก็บ embeddings ใน vector database เช่น PostgreSQL with pgvector, Pinecone, Weaviate
- รองรับการเรียก LLM หลาย providers (OpenAI, DeepSeek, Anthropic) ผ่าน abstraction layer (ตามตัวอย่างใน 6.1.3)
- มีการ monitoring latency, cost, และ quality ของคำตอบ ผ่าน logging และ metrics

**ตัวอย่างเพิ่มเติม: การใช้ Ollama สำหรับ local LLMs**  
สามารถรันโมเดลในเครื่องด้วย Ollama และเรียกผ่าน LangChain:

```python
from langchain_community.llms import Ollama
llm = Ollama(model="llama2")
```

**IoT Systems**  
- ใช้ MQTT broker (เช่น Mosquitto, EMQX) สำหรับรับส่งข้อมูลจาก devices
- เก็บข้อมูลอนุกรมเวลาใน InfluxDB
- มีบริการ backend (เช่น Go service) สำหรับรับข้อมูล via MQTT subscriber และบันทึกลง InfluxDB
- จัดทำ API สำหรับ query ข้อมูลและแสดงผลผ่าน dashboard (Grafana)

**ตัวอย่าง MQTT Subscriber ใน Go โดยใช้ Eclipse Paho:**

```go
package main

import (
    "fmt"
    mqtt "eclipse/paho.mqtt.golang"
    "influxdata/influxdb-client-go/v2"
    "log"
    "time"
)

func main() {
    // MQTT
    opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
    opts.SetClientID("go_subscriber")
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }

    // InfluxDB
    influxClient := influxdb2.NewClient("http://localhost:8086", "my-token")
    writeAPI := influxClient.WriteAPI("my-org", "my-bucket")

    // Subscribe
    token := client.Subscribe("sensor/#", 1, func(client mqtt.Client, msg mqtt.Message) {
        // Parse message and write to InfluxDB
        point := influxdb2.NewPoint(
            "sensor",
            map[string]string{"device": msg.Topic()},
            map[string]interface{}{"value": string(msg.Payload())},
            time.Now(),
        )
        writeAPI.WritePoint(point)
        writeAPI.Flush()
        fmt.Printf("Received: %s from %s\n", msg.Payload(), msg.Topic())
    })
    token.Wait()
    select {} // block forever
}
```

**Production Deployment**  
- Deploy AI services บน Kubernetes หรือใช้ serverless (AWS Lambda) สำหรับบาง workloads
- ใช้ logging และ monitoring ครบวงจร (ELK, Grafana, Prometheus)
- ปรับแต่ง performance และ cost สำหรับ AI workloads เช่น การใช้ GPU, การ cache คำตอบ

**English**  
**AI/LLM Systems**  
- Use LangChain to build pipelines for RAG, agents, chains.
- Store embeddings in vector databases like PostgreSQL with pgvector, Pinecone, Weaviate.
- Support multiple LLM providers (OpenAI, DeepSeek, Anthropic) via an abstraction layer (as in 6.1.3).
- Monitor latency, cost, and answer quality through logging and metrics.

**Additional Example: Using Ollama for local LLMs**  
Run models locally with Ollama and call via LangChain:

```python
from langchain_community.llms import Ollama
llm = Ollama(model="llama2")
```

**IoT Systems**  
- Use MQTT broker (e.g., Mosquitto, EMQX) for device data ingestion.
- Store time-series data in InfluxDB.
- Provide backend service (e.g., Go) that subscribes to MQTT and writes to InfluxDB.
- Build APIs to query data and visualize with dashboards (Grafana).

**Example MQTT Subscriber in Go using Eclipse Paho:**

```go
package main

import (
    "fmt"
    mqtt "eclipse/paho.mqtt.golang"
    "influxdata/influxdb-client-go/v2"
    "log"
    "time"
)

func main() {
    // MQTT
    opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
    opts.SetClientID("go_subscriber")
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }

    // InfluxDB
    influxClient := influxdb2.NewClient("http://localhost:8086", "my-token")
    writeAPI := influxClient.WriteAPI("my-org", "my-bucket")

    // Subscribe
    token := client.Subscribe("sensor/#", 1, func(client mqtt.Client, msg mqtt.Message) {
        // Parse message and write to InfluxDB
        point := influxdb2.NewPoint(
            "sensor",
            map[string]string{"device": msg.Topic()},
            map[string]interface{}{"value": string(msg.Payload())},
            time.Now(),
        )
        writeAPI.WritePoint(point)
        writeAPI.Flush()
        fmt.Printf("Received: %s from %s\n", msg.Payload(), msg.Topic())
    })
    token.Wait()
    select {} // block forever
}
```

**Production Deployment**  
- Deploy AI services on Kubernetes or use serverless (AWS Lambda) for certain workloads.
- Implement comprehensive logging and monitoring (ELK, Grafana, Prometheus).
- Optimize performance and cost for AI workloads (e.g., GPU usage, answer caching).

---

**หมายเหตุ:** เอกสารนี้เป็นแนวทางที่สมบูรณ์ พร้อมโค้ดตัวอย่างที่สามารถนำไปปรับใช้ได้ทันที ควรปรับเปลี่ยนตามความเหมาะสมของแต่ละโปรเจกต์  
**Note:** This document is a complete guideline with ready-to-use code examples. Adapt as needed for each project.