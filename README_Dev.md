# golang-rest-api-template

[![license](https://img.shields.io/badge/license-MIT-green)](https://raw.githubusercontent.com/araujo88/golang-rest-api-template/main/LICENSE)
[![build](https://github.com/araujo88/golang-rest-api-template//actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/araujo88/golang-rest-api-template/actions/workflows/go.yml)

## Overview

This repository provides a template for building a RESTful API using Go with features like JWT Authentication, rate limiting, Swagger documentation, and database operations using GORM. The application uses the Gin Gonic web framework and is containerized using Docker.

## Features

- RESTful API endpoints for CRUD operations.
- JWT Authentication.
- Rate Limiting.
- Swagger Documentation.
- PostgreSQL database integration using GORM.
- Redis cache.
- MongoDB for logging storage.
- Dockerized application for easy setup and deployment.

## Folder structure

```
golang-rest-api-template/
├── bin
│  └── server
├── cmd
│  └── server
│     └── main.go
├── docker-compose.yml
├── Dockerfile
├── docs
│  ├── docs.go
│  ├── swagger.json
│  └── swagger.yaml
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── pkg
│  ├── api
│  │  ├── books.go
│  │  ├── books_test.go
│  │  ├── router.go
│  │  └── user.go
│  ├── auth
│  │  ├── auth.go
│  │  └── auth_test.go
│  ├── cache
│  │  ├── cache.go
│  │  ├── cache_mock.go
│  │  └── cache_test.go
│  ├── database
│  │  ├── db.go
│  │  ├── db_mock.go
│  │  └── db_test.go
│  ├── middleware
│  │  ├── api_key.go
│  │  ├── authenticateJWT.go
│  │  ├── cors.go
│  │  ├── rate_limit.go
│  │  ├── security.go
│  │  └── xss.go
│  └── models
│     ├── book.go
│     └── user.go
├── README.md
├── scripts
│  ├── generate_key
│  └── generate_key.go
└── vendor
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker
- Docker Compose

### Installation

1. Clone the repository

```bash
git clone https://github.com/araujo88/golang-rest-api-template
```

2. Navigate to the directory

```bash
cd golang-rest-api-template
```

3. Build and run the Docker containers

```bash
make up
```

Please refer to the [Makefile](./Makefile) if you need to build in the local environment.

### Environment Variables

You can set the environment variables in the `.env` file. Here are some important variables:

- `POSTGRES_HOST`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `JWT_SECRET`
- `API_SECRET_KEY`

### API Documentation

The API is documented using Swagger and can be accessed at:

```
http://localhost:8001/swagger/index.html
```

## Usage

### Endpoints

- `GET /api/v1/books`: Get all books.
- `GET /api/v1/books/:id`: Get a single book by ID.
- `POST /api/v1/books`: Create a new book.
- `PUT /api/v1/books/:id`: Update a book.
- `DELETE /api/v1/books/:id`: Delete a book.
- `POST /api/v1/login`: Login.
- `POST /api/v1/register`: Register a new user.

### Authentication

To use authenticated routes, you must include the `Authorization` header with the JWT token.

```bash
curl -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8001/api/v1/books
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## End-to-End (E2E) Tests

This project contains end-to-end (E2E) tests to verify the functionality of the API. The tests are written in Python using the `pytest` framework.

### Prerequisites

Before running the tests, ensure you have the following:

- Python 3.x installed
- `pip` (Python package manager)
- The API service running locally or on a staging server
- API key available

### Setup

#### 1. Create a virtual environment (optional but recommended):

```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

#### 2. Install dependencies:

```bash
pip install -r requirements.txt
```

The main dependency is `requests`, but you may need to include it in your `requirements.txt` file if it's not already listed.

#### 3. Set up the environment variables:

You need to set the `BASE_URL` and `API_KEY` as environment variables before running the tests.

For a **local** API service:

```bash
export BASE_URL=http://localhost:8001/api/v1
export API_KEY=your-api-key-here
```

For a **staging** server:

```bash
export BASE_URL=https://staging-server-url.com/api/v1
export API_KEY=your-api-key-here
```

On **Windows**, you can use:

```bash
set BASE_URL=http://localhost:8001/api/v1
set API_KEY=your-api-key-here
```

#### 4. Run the tests:

Once the environment variables are set, you can run the tests using `pytest`:

```bash
pytest test_e2e.py
```

### Test Structure

The tests will perform the following actions:

1. Register a new user and obtain a JWT token.
2. Create a new book in the system.
3. Retrieve all books and verify the created book is present.
4. Retrieve a specific book by its ID.
5. Update the book's details.
6. Delete the book and verify it is no longer accessible.

Each test includes assertions to ensure that the API behaves as expected.
ปัญหาที่พบคือ Go ไม่สามารถหา package ในโปรเจคของคุณได้ เพราะมันไปค้นหาใน `GOROOT` (`C:\Program Files\Go\src`) แทนที่จะมองหาใน module ของคุณ ซึ่งแสดงว่า Go กำลังทำงานใน **GOPATH mode** ไม่ใช่ **module mode** ทั้งที่คุณใช้ `-mod=mod` แล้วก็ตาม

## สาเหตุ
- `GO111MODULE` อาจถูกตั้งเป็น `off` หรือ `auto` แต่เนื่องจากโปรเจคของคุณ **ไม่ได้อยู่ใน `GOPATH`** (`%USERPROFILE%\go\src` โดยปกติ) และอาจไม่มีไฟล์ `go.mod` ในขณะที่รัน หรือ Go version เก่ากว่า 1.16 ที่ default เป็น `auto` และอาจสลับไป GOPATH mode ได้

## วิธีแก้ไข

### 1. ตรวจสอบและเปิดใช้งาน Go module mode
รันคำสั่งนี้ใน PowerShell (ที่อยู่ในโฟลเดอร์โปรเจค):
```powershell
$env:GO111MODULE="on"
go env GO111MODULE   # ควรแสดง "on"
```

จากนั้นลองรันอีกครั้ง:
```powershell
go run cmd/server/main.go
```

หรือถ้าต้องการตั้งค่าให้ถาวร:
```powershell
go env -w GO111MODULE=on
```

### 2. ตรวจสอบไฟล์ go.mod
ตรวจสอบว่าไฟล์ `go.mod` อยู่ใน root ของโปรเจค (`D:\git\golangrestapi\go.mod`) และมีเนื้อหาขึ้นต้นด้วย:
```
module golang-rest-api-template

go 1.24.0
```
และไม่มี syntax error

### 3. ล้าง vendor และ dependencies
บางครั้ง vendor ที่ไม่สมบูรณ์หรือค้างทำให้ module mode สับสน:
```powershell
Remove-Item -Recurse -Force vendor
go mod tidy
go mod download
```

### 4. ลองรันด้วยคำสั่งที่ชัดเจน
ใช้ path แบบเต็ม (ใช้ backslash หรือ forward slash ก็ได้):
```powershell
go run -mod=mod ./cmd/server/main.go
# หรือ
go run -mod=mod cmd/server/main.go
```

### 5. ตรวจสอบ environment อื่น ๆ
รัน `go env` ดูว่า `GOPATH`, `GOROOT`, `GO111MODULE` ถูกต้องหรือไม่
โดยเฉพาะ `GOPATH` ไม่ควรเป็นโฟลเดอร์เดียวกับโปรเจคของคุณ (ควรเป็นคนละที่ เช่น `C:\Users\<คุณ>\go`)

### 6. กรณีเร่งด่วน (ถ้ายังไม่สำเร็จ)
ลอง build เป็น executable เพื่อทดสอบ:
```powershell
go build -o app.exe ./cmd/server
./app.exe
```

## สรุปขั้นตอนที่ควรทำตามลำดับ
1. ตั้งค่า `GO111MODULE=on` (ชั่วคราวใน PowerShell นั้น)
2. ลบโฟลเดอร์ `vendor`
3. รัน `go mod tidy`
4. รัน `go run cmd/server/main.go`

ถ้าทำตามแล้วยังไม่ได้ผล กรุณาแจ้งผลลัพธ์ของคำสั่ง `go env` และแสดงโครงสร้างโฟลเดอร์ปัจจุบัน (ใช้ `ls` หรือ `dir`) เพื่อตรวจสอบเพิ่มเติมครับ