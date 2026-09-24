package main

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

import (
	"icmongolang/cmd"
	"icmongolang/config"
	"icmongolang/pkg/helpers"
	"log"
)

func main() {
	// Load config
	v, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.ParseConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	// ✅ Inject the timezone from config into helpers
	helpers.SetDefaultTimezone(cfg.Server.Timezone)

	// ... start your application (routes, MQTT, etc.)
	cmd.Execute()
}
