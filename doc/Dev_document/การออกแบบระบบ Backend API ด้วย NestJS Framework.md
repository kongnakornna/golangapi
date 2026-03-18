<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## การออกแบบระบบ ReactJS Frontend

NestJS เป็น progressive Node.js framework ที่ใช้ TypeScript และออกแบบตามหลัก SOLID principles เหมาะสำหรับสร้าง scalable server-side applications  สำหรับระบบศูนย์บริการรถยนต์ออนไลน์ ระบบจะประกอบด้วย JWT authentication, query caching, และโครงสร้าง modular ที่มีประสิทธิภาพ[^1][^2][^3][^4]

## โครงสร้างโปรเจกต์ (Project Structure)
### Modular Architecture

```
  car-service-api/
  ├── src/
  │   ├── config/                    # Configuration files
  │   │   ├── database.config.ts
  │   │   ├── redis.config.ts
  │   │   └── jwt.config.ts
  │   │
  │   ├── common/                    # Global shared resources
  │   │   ├── decorators/           # Custom decorators
  │   │   │   ├── roles.decorator.ts
  │   │   │   └── public.decorator.ts
  │   │   ├── guards/               # Global guards
  │   │   │   ├── jwt-auth.guard.ts
  │   │   │   └── roles.guard.ts
  │   │   ├── interceptors/         # Global interceptors
  │   │   │   ├── logging.interceptor.ts
  │   │   │   └── transform.interceptor.ts
  │   │   ├── filters/              # Exception filters
  │   │   │   └── http-exception.filter.ts
  │   │   ├── pipes/                # Validation pipes
  │   │   │   └── validation.pipe.ts
  │   │   └── middleware/           # Middleware
  │   │       └── logger.middleware.ts
  │   │
  │   ├── modules/     # Feature modules
  │   │   ├── auth/
  │   │   │   ├── dto/
  │   │   │   │   ├── login.dto.ts
  │   │   │   │   └── register.dto.ts
  │   │   │   ├── strategies/
  │   │   │   │   ├── jwt.strategy.ts
  │   │   │   │   └── local.strategy.ts
  │   │   │   ├── auth.controller.ts
  │   │   │   ├── auth.service.ts
  │   │   │   └── auth.module.ts
  │   │   │
  │   │   ├── users/
  │   │   │   ├── entities/
  │   │   │   │   └── user.entity.ts
  │   │   │   ├── dto/
  │   │   │   │   ├── create-user.dto.ts
  │   │   │   │   └── update-user.dto.ts
  │   │   │   ├── users.controller.ts
  │   │   │   ├── users.service.ts
  │   │   │   ├── users.repository.ts
  │   │   │   └── users.module.ts
  │   │   │
  │   │   ├── bookings/
  │   │   │   ├── entities/
  │   │   │   │   └── booking.entity.ts
  │   │   │   ├── dto/
  │   │   │   ├── bookings.controller.ts
  │   │   │   ├── bookings.service.ts
  │   │   │   └── bookings.module.ts
  │   │   │
  │   │   ├── repairs/
  │   │   ├── vehicles/
  │   │   ├── payments/
  │   │   └── notifications/
  │   │
  │   ├── database/      # Database related
  │   │   ├── migrations/
  │   │   ├── seeds/
  │   │   └── factories/
  │   │
  │   ├── app.module.ts             # Root module
  │   └── main.ts                   # Entry point
  │
  ├── test/                         # E2E tests
  ├── .env                          # Environment variables
  └── package.json
```


## Core Modules และ Functions
### 1. Authentication Module (JWT)

**Installation**

```bash
npm install @nestjs/jwt @nestjs/passport passport passport-jwt bcrypt
npm install -D @types/passport-jwt @types/bcrypt
```

**JWT Configuration** (`src/config/jwt.config.ts`)

```typescript
import { ConfigService } from '@nestjs/config';
import { JwtModuleOptions } from '@nestjs/jwt';

export const getJwtConfig = (configService: ConfigService): JwtModuleOptions => ({
  secret: configService.get('JWT_SECRET'),
  signOptions: {
    expiresIn: configService.get('JWT_EXPIRATION', '7d'),
    algorithm: 'HS256',
  },
});
```

**Auth Module** (`src/modules/auth/auth.module.ts`)

```typescript
import { Module } from '@nestjs/common';
import { JwtModule } from '@nestjs/jwt';
import { PassportModule } from '@nestjs/passport';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { AuthService } from './auth.service';
import { AuthController } from './auth.controller';
import { JwtStrategy } from './strategies/jwt.strategy';
import { LocalStrategy } from './strategies/local.strategy';
import { UsersModule } from '../users/users.module';
import { getJwtConfig } from '../../config/jwt.config';

@Module({
  imports: [
    UsersModule,
    PassportModule.register({ defaultStrategy: 'jwt' }),
    JwtModule.registerAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: getJwtConfig,
    }),
  ],
  controllers: [AuthController],
  providers: [AuthService, JwtStrategy, LocalStrategy],
  exports: [AuthService, JwtModule],
})
export class AuthModule {}
```

**Auth Service** (`src/modules/auth/auth.service.ts`)

```typescript
import { Injectable, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { UsersService } from '../users/users.service';
import * as bcrypt from 'bcrypt';
import { LoginDto } from './dto/login.dto';
import { RegisterDto } from './dto/register.dto';

@Injectable()
export class AuthService {
  constructor(
    private usersService: UsersService,
    private jwtService: JwtService,
  ) {}

  // Validate user credentials
  async validateUser(email: string, password: string): Promise<any> {
    const user = await this.usersService.findByEmail(email);
    
    if (!user) {
      throw new UnauthorizedException('Invalid credentials');
    }

    const isPasswordValid = await bcrypt.compare(password, user.password);
    
    if (!isPasswordValid) {
      throw new UnauthorizedException('Invalid credentials');
    }

    const { password: _, ...result } = user;
    return result;
  }

  // Register new user
  async register(registerDto: RegisterDto) {
    const hashedPassword = await bcrypt.hash(registerDto.password, 10);
    
    const user = await this.usersService.create({
      ...registerDto,
      password: hashedPassword,
    });

    const token = await this.generateToken(user);
    
    return {
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
      },
      token,
    };
  }

  // Login user
  async login(loginDto: LoginDto) {
    const user = await this.validateUser(loginDto.email, loginDto.password);
    const token = await this.generateToken(user);

    return {
      user,
      token,
    };
  }

  // Generate JWT token
  private async generateToken(user: any) {
    const payload = { 
      sub: user.id, 
      email: user.email,
      roles: user.roles || ['user']
    };

    return {
      accessToken: this.jwtService.sign(payload),
      expiresIn: '7d',
    };
  }

  // Refresh token
  async refreshToken(userId: string) {
    const user = await this.usersService.findById(userId);
    return this.generateToken(user);
  }
}
```

**JWT Strategy** (`src/modules/auth/strategies/jwt.strategy.ts`)

```typescript
import { Injectable, UnauthorizedException } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';
import { ConfigService } from '@nestjs/config';
import { UsersService } from '../../users/users.service';

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(
    private configService: ConfigService,
    private usersService: UsersService,
  ) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: configService.get('JWT_SECRET'),
    });
  }

  async validate(payload: any) {
    const user = await this.usersService.findById(payload.sub);
    
    if (!user) {
      throw new UnauthorizedException();
    }

    return {
      id: user.id,
      email: user.email,
      roles: user.roles,
    };
  }
}
```

**Auth Controller** (`src/modules/auth/auth.controller.ts`)

```typescript
import { Controller, Post, Body, UseGuards, Get, Request } from '@nestjs/common';
import { AuthService } from './auth.service';
import { LoginDto } from './dto/login.dto';
import { RegisterDto } from './dto/register.dto';
import { JwtAuthGuard } from '../../common/guards/jwt-auth.guard';
import { Public } from '../../common/decorators/public.decorator';

@Controller('api/v1/auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Public()
  @Post('register')
  async register(@Body() registerDto: RegisterDto) {
    return this.authService.register(registerDto);
  }

  @Public()
  @Post('login')
  async login(@Body() loginDto: LoginDto) {
    return this.authService.login(loginDto);
  }

  @UseGuards(JwtAuthGuard)
  @Get('profile')
  async getProfile(@Request() req) {
    return req.user;
  }

  @UseGuards(JwtAuthGuard)
  @Post('refresh')
  async refreshToken(@Request() req) {
    return this.authService.refreshToken(req.user.id);
  }
}
```


### 2. Database Configuration (TypeORM + PostgreSQL)

**Installation**

```bash
  npm install @nestjs/typeorm typeorm pg typeorm-naming-strategies
```

**Database Config** (`src/config/database.config.ts`)

```typescript
  import { TypeOrmModuleOptions } from '@nestjs/typeorm';
  import { ConfigService } from '@nestjs/config';
  import { SnakeNamingStrategy } from 'typeorm-naming-strategies';

  export const getDatabaseConfig = (
    configService: ConfigService,
  ): TypeOrmModuleOptions => ({
    type: 'postgres',
    host: configService.get('DB_HOST', 'localhost'),
    port: configService.get('DB_PORT', 5432),
    username: configService.get('DB_USERNAME'),
    password: configService.get('DB_PASSWORD'),
    database: configService.get('DB_NAME'),
    entities: ['dist/**/*.entity{.ts,.js}'],
    migrations: ['dist/database/migrations/*{.ts,.js}'],
    synchronize: configService.get('NODE_ENV') === 'development',
    logging: configService.get('NODE_ENV') === 'development',
    namingStrategy: new SnakeNamingStrategy(),
    cache: {
      type: 'redis',
      options: {
        host: configService.get('REDIS_HOST', 'localhost'),
        port: configService.get('REDIS_PORT', 6379),
      },
      duration: 30000, // 30 seconds
    },
  });
```

**Entity Example** (`src/modules/users/entities/user.entity.ts`)

```typescript
  import {
    Entity,
    PrimaryGeneratedColumn,
    Column,
    CreateDateColumn,
    UpdateDateColumn,
    DeleteDateColumn,
    Index,
  } from 'typeorm';

  @Entity('users')
  @Index(['email'], { unique: true })
  export class User {
    @PrimaryGeneratedColumn('uuid')
    id: string;

    @Column({ unique: true })
    email: string;

    @Column()
    password: string;

    @Column()
    name: string;

    @Column({ nullable: true })
    phone: string;

    @Column('simple-array', { default: '' })
    roles: string[];

    @Column({ default: true })
    isActive: boolean;

    @CreateDateColumn()
    createdAt: Date;

    @UpdateDateColumn()
    updatedAt: Date;

    @DeleteDateColumn()
    deletedAt: Date;
  }
```


### 3. Redis Caching

**Installation**

```bash
  npm install @nestjs/cache-manager cache-manager cache-manager-redis-store
  npm install -D @types/cache-manager
```

**Redis Config** (`src/config/redis.config.ts`)

```typescript
    import { CacheModuleOptions } from '@nestjs/cache-manager';
    import { ConfigService } from '@nestjs/config';
    import * as redisStore from 'cache-manager-redis-store';

    export const getRedisConfig = (
      configService: ConfigService,
    ): CacheModuleOptions => ({
      store: redisStore,
      host: configService.get('REDIS_HOST', 'localhost'),
      port: configService.get('REDIS_PORT', 6379),
      ttl: configService.get('CACHE_TTL', 300), // 5 minutes default
      max: 100, // maximum number of items in cache
    });
```

**Cache Module Setup** (`app.module.ts`)

```typescript
    import { Module } from '@nestjs/common';
    import { CacheModule } from '@nestjs/cache-manager';
    import { ConfigModule, ConfigService } from '@nestjs/config';
    import { getRedisConfig } from './config/redis.config';

    @Module({
      imports: [
        ConfigModule.forRoot({ isGlobal: true }),
        CacheModule.registerAsync({
          imports: [ConfigModule],
          inject: [ConfigService],
          useFactory: getRedisConfig,
          isGlobal: true,
        }),
        // Other modules...
      ],
    })
    export class AppModule {}
```

**Service with Caching** (`src/modules/bookings/bookings.service.ts`)

```typescript
        import { Injectable, Inject } from '@nestjs/common';
        import { InjectRepository } from '@nestjs/typeorm';
        import { Repository } from 'typeorm';
        import { CACHE_MANAGER } from '@nestjs/cache-manager';
        import { Cache } from 'cache-manager';
        import { Booking } from './entities/booking.entity';

        @Injectable()
        export class BookingsService {
          constructor(
            @InjectRepository(Booking)
            private bookingsRepository: Repository<Booking>,
            @Inject(CACHE_MANAGER)
            private cacheManager: Cache,
          ) {}

          async findAll(): Promise<Booking[]> {
            const cacheKey = 'bookings:all';
            
            // Try to get from cache first
            const cachedData = await this.cacheManager.get<Booking[]>(cacheKey);
            
            if (cachedData) {
              return cachedData;
            }

            // If not in cache, fetch from database
            const bookings = await this.bookingsRepository.find();
            
            // Store in cache for future requests
            await this.cacheManager.set(cacheKey, bookings, 300); // 5 minutes
            
            return bookings;
          }

          async findById(id: string): Promise<Booking> {
            const cacheKey = `booking:${id}`;
            
            const cached = await this.cacheManager.get<Booking>(cacheKey);
            if (cached) return cached;

            const booking = await this.bookingsRepository.findOne({ where: { id } });
            
            if (booking) {
              await this.cacheManager.set(cacheKey, booking, 600); // 10 minutes
            }
            
            return booking;
          }

          async create(createBookingDto: any): Promise<Booking> {
            const booking = this.bookingsRepository.create(createBookingDto);
            const saved = await this.bookingsRepository.save(booking);
            
            // Invalidate list cache when new booking is created
            await this.cacheManager.del('bookings:all');
            
            return saved;
          }

          async update(id: string, updateBookingDto: any): Promise<Booking> {
            await this.bookingsRepository.update(id, updateBookingDto);
            
            // Invalidate both specific and list caches
            await Promise.all([
              this.cacheManager.del(`booking:${id}`),
              this.cacheManager.del('bookings:all'),
            ]);
            
            return this.findById(id);
          }
        }
```

**Controller with Cache Interceptor**

```typescript
      import { 
        Controller, 
        Get, 
        UseInterceptors 
      } from '@nestjs/common';
      import { CacheInterceptor, CacheTTL } from '@nestjs/cache-manager';
      import { BookingsService } from './bookings.service';

      @Controller('api/v1/bookings')
      export class BookingsController {
        constructor(private readonly bookingsService: BookingsService) {}

        @Get()
        @UseInterceptors(CacheInterceptor)
        @CacheTTL(300) // Override default TTL to 5 minutes
        async findAll() {
          return this.bookingsService.findAll();
        }

        @Get(':id')
        @UseInterceptors(CacheInterceptor)
        @CacheTTL(600) // 10 minutes for individual items
        async findOne(@Param('id') id: string) {
          return this.bookingsService.findById(id);
        }
      }
```


## Request Flow และ Lifecycle

### Request Processing Order

```
  Incoming Request
      ↓
  1. Middleware (LoggerMiddleware)
      ↓
  2. Guards (JwtAuthGuard, RolesGuard)
      ↓
  3. Interceptors (Before - LoggingInterceptor, CacheInterceptor)
      ↓
  4. Pipes (ValidationPipe)
      ↓
  5. Controller Method
      ↓
  6. Service Layer
      ↓
  7. Repository/Database
      ↓
  8. Interceptors (After - TransformInterceptor)
      ↓
  9. Filters (ExceptionFilter if error)
      ↓
  Response to Client
```

### Guards Implementation
**JWT Auth Guard** (`src/common/guards/jwt-auth.guard.ts`)

```typescript
import { Injectable, ExecutionContext } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { AuthGuard } from '@nestjs/passport';
import { IS_PUBLIC_KEY } from '../decorators/public.decorator';

@Injectable()
export class JwtAuthGuard extends AuthGuard('jwt') {
  constructor(private reflector: Reflector) {
    super();
  }

  canActivate(context: ExecutionContext) {
    const isPublic = this.reflector.getAllAndOverride<boolean>(IS_PUBLIC_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);
    
    if (isPublic) {
      return true;
    }
    
    return super.canActivate(context);
  }
}
```
**Roles Guard** (`src/common/guards/roles.guard.ts`)

```typescript
import { Injectable, CanActivate, ExecutionContext } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { ROLES_KEY } from '../decorators/roles.decorator';

@Injectable()
export class RolesGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  canActivate(context: ExecutionContext): boolean {
    const requiredRoles = this.reflector.getAllAndOverride<string[]>(
      ROLES_KEY,
      [context.getHandler(), context.getClass()],
    );
    
    if (!requiredRoles) {
      return true;
    }
    
    const { user } = context.switchToHttp().getRequest();
    return requiredRoles.some((role) => user.roles?.includes(role));
  }
}
```


### Interceptors

**Logging Interceptor** (`src/common/interceptors/logging.interceptor.ts`)

```typescript
import {
  Injectable,
  NestInterceptor,
  ExecutionContext,
  CallHandler,
  Logger,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

@Injectable()
export class LoggingInterceptor implements NestInterceptor {
  private readonly logger = new Logger(LoggingInterceptor.name);

  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    const request = context.switchToHttp().getRequest();
    const { method, url } = request;
    const now = Date.now();

    return next.handle().pipe(
      tap(() => {
        const response = context.switchToHttp().getResponse();
        const delay = Date.now() - now;
        this.logger.log(
          `${method} ${url} ${response.statusCode} - ${delay}ms`,
        );
      }),
    );
  }
}
```

**Transform Interceptor** (`src/common/interceptors/transform.interceptor.ts`)

```typescript
import {
  Injectable,
  NestInterceptor,
  ExecutionContext,
  CallHandler,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export interface Response<T> {
  data: T;
  statusCode: number;
  message: string;
  timestamp: string;
}

@Injectable()
export class TransformInterceptor<T>
  implements NestInterceptor<T, Response<T>>
{
  intercept(
    context: ExecutionContext,
    next: CallHandler,
  ): Observable<Response<T>> {
    return next.handle().pipe(
      map((data) => ({
        data,
        statusCode: context.switchToHttp().getResponse().statusCode,
        message: 'Success',
        timestamp: new Date().toISOString(),
      })),
    );
  }
}
```


### Pipes

**Validation Pipe Setup** (`main.ts`)

```typescript
import { NestFactory } from '@nestjs/core';
import { ValidationPipe } from '@nestjs/common';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  // Global validation pipe
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true, // Strip properties that don't have decorators
      forbidNonWhitelisted: true, // Throw error if non-whitelisted values exist
      transform: true, // Automatically transform payloads to DTO instances
      transformOptions: {
        enableImplicitConversion: true,
      },
    }),
  );
  
  await app.listen(3000);
}
bootstrap();
```


## Environment Configuration

**.env File**

```env
    # Application
    NODE_ENV=development
    PORT=3000

    # Database
    DB_HOST=localhost
    DB_PORT=5432
    DB_USERNAME=postgres
    DB_PASSWORD=password
    DB_NAME=car_service_db
    DB_SYNCHRONIZE=false

    # JWT
    JWT_SECRET=your-super-secret-jwt-key-change-in-production
    JWT_EXPIRATION=7d

    # Redis
    REDIS_HOST=localhost
    REDIS_PORT=6379
    CACHE_TTL=300

    # API
    API_PREFIX=api/v1
```


## สรุป Libraries และ Modules

| Library/Module | Purpose | Installation |
| :-- | :-- | :-- |
| `@nestjs/jwt` | JWT token generation | Required for authentication |
| `@nestjs/passport` | Authentication strategies | Works with passport-jwt |
| `@nestjs/typeorm` | Database ORM | PostgreSQL integration |
| `@nestjs/cache-manager` | Caching layer | Redis caching support |
| `@nestjs/config` | Environment configuration | Manages .env files |
| `class-validator` | DTO validation | Input validation pipes |
| `class-transformer` | Object transformation | DTO transformation |

ระบบนี้ใช้ modular architecture ที่แยก concerns ชัดเจน มี authentication ที่แข็งแกร่งด้วย JWT, caching layer ด้วย Redis และ lifecycle management ที่มีประสิทธิภาพผ่าน guards, interceptors และ pipes

