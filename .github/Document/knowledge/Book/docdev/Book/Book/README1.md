# Go Restful API dev

An API dev written in Golang with chi-route and Gorm. Write restful API with fast development and developer friendly.

## Architecture

In this project use 3 layer architecture

- Models
- Repository
- Usecase
- Delivery

## Features

- CRUD
- Jwt, refresh token saved in redis
- Cached user in redis
- Email verification
- Forget/reset password, send email

## Technical

- `chi`: router and middleware
- `viper`: configuration
- `cobra`: CLI features
- `gorm`: orm
- `validator`: data validation
- `jwt`: jwt authentication
- `zap`: logger
- `gomail`: email
- `hermes`: generate email body
- `air`: hot-reload

## Start Application

### Generate the Private and Public Keys

- Generate the private and public keys: [travistidwell.com/jsencrypt/demo/](https://travistidwell.com/jsencrypt/demo/)
- Copy the generated private key and visit this Base64 encoding website to convert it to base64
- Copy the base64 encoded key and add it to the `config/config-local.yml` file as `jwt`
- Similar for public key

### Stmp mail config

- Create [mailtrap](https://mailtrap.io/) account
- Create new inboxes
- Update smtp config `config/config-local.yml` file as `smtpEmail`

### Run

- `docker-compose up`
- Swagger: [localhost:5000/swagger/](http://localhost:5000/swagger/)

## TODO

- Traefik
- Config using .env
- Linter
- Jaeger
- Production docker file version
- Mock database using gomock

## Acknowledgements

- [github.com/dhax/go-base](https://github.com/dhax/go-base)
- [github.com/akmamun/go-fication](https://github.com/akmamun/go-fication)
- [github.com/wpcodevo/golang-fiber-jwt](https://github.com/wpcodevo/golang-fiber-jwt)
- [github.com/wpcodevo/golang-fiber](https://github.com/wpcodevo/golang-fiber)
- [github.com/kienmatu/togo](https://github.com/kienmatu/togo)
- [github.com/AleksK1NG/Go-Clean-Architecture-REST-API](https://github.com/AleksK1NG/Go-Clean-Architecture-REST-API)
- [github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [codevoweb.com/golang-and-gorm-user-registration-email-verification/](https://codevoweb.com/golang-and-gorm-user-registration-email-verification/)
- [codevoweb.com/golang-gorm-postgresql-user-registration-with-refresh-tokens/](https://codevoweb.com/golang-gorm-postgresql-user-registration-with-refresh-tokens/)
- [codevoweb.com/how-to-implement-google-oauth2-in-golang/](https://codevoweb.com/how-to-implement-google-oauth2-in-golang/)
- [codevoweb.com/how-to-upload-single-and-multiple-files-in-golang/](https://codevoweb.com/how-to-upload-single-and-multiple-files-in-golang/)
- [codevoweb.com/forgot-reset-passwords-in-golang-with-html-email/](https://codevoweb.com/forgot-reset-passwords-in-golang-with-html-email/)
- [techmaster.vn/posts/34577/kien-truc-sach-voi-golang](https://techmaster.vn/posts/34577/kien-truc-sach-voi-golang)



### Installation

Perfect! You're setting up an existing Go project (gorestapi). Here's how to properly set it up and run it:

## Complete Setup Steps for Your gorestapi Project

```bash
# 1. Clone the repository
git clone https://gorestapi
cd gorestapi

# 2. Download and tidy up dependencies
go mod tidy

# 3. Verify the module is set up correctly
go mod verify


# Inside container or locally
go mod tidy
go mod download
go mod verify

# 4. Run the application


# Or if the main file is in cmd directory:
go run cmd/gorestapi/main.go

# Project Structure Check

# Check if these exist
ls docker-compose.yml
ls Dockerfile
ls main.go
ls go.mod

```

## Common Next Steps Based on Project Type

### If it's a REST API (likely based on the name):

```bash
# Run the server (usually on port 5000)
go run main.go

# In another terminal, test the API endpoints
curl http://localhost:5000/health
curl http://localhost:5000/api/users
```

### If you need to build the binary:

```bash
# Build for current OS
go build -o gorestapi

# Run the binary
./gorestapi

# Build for different platforms (optional)
# Windows:
GOOS=windows GOARCH=amd64 go build -o gorestapi.exe
# Linux:
GOOS=linux GOARCH=amd64 go build -o gorestapi-linux
# macOS:
GOOS=darwin GOARCH=arm64 go build -o gorestapi-mac
```

## Troubleshooting Common Issues

### If `go mod tidy` fails:

```bash
# Check Go version (need Go 1.11+ for modules)
go version

# If Go is not installed, install it first:
# On Ubuntu/Debian:
sudo apt install golang-go
# On macOS:
brew install go
# On Windows: Download from https://golang.org/dl

# Clean module cache if needed
go clean -modcache
go mod tidy
```

### If there are dependency conflicts:

```bash
# Update all dependencies to compatible versions
go get -u ./...
go mod tidy
```

### If the project uses environment variables:

```bash
# Check for .env file or example
ls -la | grep env

# If .env.example exists, copy it
cp .env.example .env

# Load environment variables (on Linux/macOS)
source .env
# Or export them directly
export PORT=5000
export DB_URL="postgres://..."
```

## Project Structure Check

After cloning, look for these common patterns:

```bash
# Check the project structure
ls -la

# Typical gorestapi structure might look like:
# .
# ├── main.go
# ├── go.mod
# ├── go.sum
# ├── handlers/
# ├── models/
# ├── routes/
# ├── database/
# └── config/
```

## Quick Testing Commands

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestFunctionName ./handlers

# Run tests with coverage
go test -cover ./...
```

## Development with Hot Reload (Optional)

For better development experience:

```bash
# Install air (hot reload tool)
go install github.com/cosmtrek/air@latest

# Run with air
air

# Or use nodemon if you have it installed
nodemon --exec go run main.go --signal SIGTERM
```

## Check What the API Does

Look at the code to understand the endpoints:

```bash
# Look at main.go
cat main.go

# Search for route definitions
grep -r "router\." . --include="*.go"
grep -r "HandleFunc" . --include="*.go"
```

**Need more specific help?** Let me know:
1. What error messages you're seeing (if any)
2. What database the project uses (look for postgres, mysql, sqlite in go.mod)
3. Whether there's a README with specific setup instructions

The `go mod tidy` command you mentioned is exactly the right first step after cloning!


This error indicates that your `go.mod` file has a dependency on `github.com/hiennguyen9874/go-dev@v0.0.1`, which doesn't exist or is not accessible. This appears to be an incorrectly referenced module that has been pulled in indirectly through other dependencies.

Here's how to fix this issue:

## Solution 1: Clean and Reset Dependencies

```powershell
# Clear the module cache
go clean -modcache

# Remove go.sum and try again
Remove-Item go.sum -ErrorAction SilentlyContinue

# Run go mod tidy again
go mod tidy
```

## Solution 2: Check and Fix go.mod File

Open your `go.mod` file and look for any reference to `github.com/hiennguyen9874/go-dev`. If found, remove it manually:

```powershell
# Check if the bad dependency is directly in go.mod
Select-String -Path go.mod -Pattern "hiennguyen9874"
```

If found, edit `go.mod` and remove that line, then run:
```powershell
go mod tidy
```

## Solution 3: Replace the Bad Dependency (if it's a transitive dependency)

Add a replace directive to your `go.mod` file to exclude the problematic module:

```powershell
# Add replace directive to go.mod
go mod edit -replace github.com/hiennguyen9874/go-dev=github.com/hiennguyen9874/go-dev@none
```

Then run:
```powershell
go mod tidy
```

## Solution 4: Update Your Dependencies

Try updating all dependencies to their latest versions:

```powershell
# Get latest versions of all direct dependencies
go get -u ./...

# Then tidy up
go mod tidy
```

## Solution 5: Check Your Code for Incorrect Imports

Search your codebase for any incorrect imports that might reference this module:

```powershell
# Search for the problematic import in your code
Get-ChildItem -Recurse -Include "*.go" | Select-String "hiennguyen9874"
```

If you find any, correct those imports to use the proper packages.

## Most Likely Fix

The most common cause is a corrupted module cache. Try this sequence:

```powershell
# Step 1: Clean everything
go clean -modcache
go clean -cache

# Step 2: Remove go.sum
Remove-Item go.sum -ErrorAction SilentlyContinue

# Step 3: Download dependencies again
go mod download

# Step 4: Tidy up
go mod tidy
```
 
# FIX
Great! Your application is running successfully. The output shows that your Go REST API application is working and displaying the Cobra CLI help menu.

## What's happening:

Your application has several CLI commands available:
- **`serve`** - Starts the HTTP server with your API
- **`migrate`** - Runs database migrations
- **`initdata`** - Initializes seed data
- **`worker`** - Starts background workers (likely for asynq tasks)
- **`completion`** - Generates shell autocompletion scripts

## To start your API server:

Run one of these commands:

```bash
# Start the API server (default configuration)
go run main.go serve

# Or with specific port (if your app supports flags)
go run main.go serve --port=5000

# Run database migrations first (if you haven't already)
go run main.go migrate

# Initialize data after migrations
go run main.go initdata
```

## Common next steps:

1. **Check your configuration** - Make sure you have a config file (likely `.env` or `config.yaml`) with database, Redis, and other service configurations.

2. **Start the server**:
   ```bash
   go run main.go serve
   ```

3. **Check if the server starts successfully** - You should see log messages indicating the server is running (using your zap logger).

4. **Access your API** - Once running, your API should be available at something like `http://localhost:5000`

5. **Swagger documentation** - If you've set up swagger, you might be able to access it at `http://localhost:5000/swagger/index.html`

Is the server starting correctly when you run `go run main.go serve`, or are you seeing any errors?


go run main.go serve