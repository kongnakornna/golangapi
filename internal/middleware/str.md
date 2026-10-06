icmongolang/
├── config/
│   └── config.go
├── internal/
│   ├── redis_repository.go
│   ├── models/
│   │   └── sd_user.go
│   ├── middleware/
│   │   ├── middleware.go          ⭐ แก้
│   │   ├── jwtauth.go             ⭐ แก้
│   │   ├── authorization.go       ⭐ แก้
│   │   ├── cors.go
│   │   ├── logging.go
│   │   ├── monitoring.go
│   │   ├── prometheus.go
│   │   ├── rate_limit.go
│   │   └── security.go
│   └── modules/
│       ├── auth/domain/
│       │   └── role.go            ⭐ แก้
│       ├── users/
│       │   └── usecase.go
│       └── keys/
│           ├── handlers.go        (interface)
│           ├── http/              ⭐ handlers+routes ย้ายมาที่นี่
│           │   ├── handlers.go
│           │   └── routes.go
│           ├── presenter/
│           │   └── apikey.go
│           └── usecase/
│               └── apikey.go