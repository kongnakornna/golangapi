package main

import "icmongolang/cmd"

// @title Go IoT API
// @version 1.0
// @description API for ICMON project
// @BasePath /api

// OAuth2Password (optional – keep if you need password flow)
// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl /api/auth/login
// @scope.read Grants read access
// @scope.write Grants write access

// API Key (Bearer token) – recommended
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by your access token. Example: "Bearer eyJhbGciOiJSUzI1NiIs..."

func main() {
	cmd.Execute()
}
