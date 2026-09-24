package config

import (
	"os"
	"testing"
)

// TestVerifyBrokerEnvOverride verifies MQTT broker env handling:
//   - a plain ambient MQTT_BROKER (e.g. a docker-published port) must NOT
//     override the app config;
//   - only the ICMON_-prefixed vars may shape the broker / port.
//
// Every phase pins all three keys because gotenv.Load(".env") (called from
// inside LoadConfig) re-injects the repo .env's ICMON_MQTT_* values into the
// process environment on each call, so an "unset" phase would silently pick
// them back up.
func TestVerifyBrokerEnvOverride(t *testing.T) {
	// go test runs with cwd = this package dir, but LoadConfig resolves the
	// config file relative to the repo root.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd) //nolint:errcheck

	saved := map[string]string{}
	for _, k := range []string{"MQTT_BROKER", "ICMON_MQTT_BROKER", "ICMON_MQTT_PORT"} {
		saved[k] = os.Getenv(k)
	}
	t.Cleanup(func() {
		for k, val := range saved {
			if val == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, val)
			}
		}
	})

	// pinEnv sets the plain env key plus the two ICMON_-prefixed keys so the
	// ambient environment / .env cannot leak into a phase.
	pinEnv := func(icmonBroker, icmonPort string) {
		os.Setenv("MQTT_BROKER", "tcp://localhost:1885")
		os.Setenv("ICMON_MQTT_BROKER", icmonBroker)
		os.Setenv("ICMON_MQTT_PORT", icmonPort)
	}

	load := func() string {
		v, err := LoadConfig()
		if err != nil {
			t.Fatal(err)
		}
		c, err := ParseConfig(v)
		if err != nil {
			t.Fatal(err)
		}
		return c.MQTT.Broker
	}

	// Ambient plain MQTT_BROKER is ignored; only the pinned ICMON_* pair
	// shapes the broker (mqtt://localhost + port 1883).
	pinEnv("mqtt://localhost", "1883")
	got := load()
	t.Logf("ambient MQTT_BROKER=1885 + pin 1883 -> effective broker=%s", got)
	if got != "mqtt://localhost:1883" {
		t.Fatalf("expected mqtt://localhost:1883, got %s", got)
	}

	// A broker with its own explicit port is honored.
	pinEnv("tcp://localhost:1999", "1999")
	got = load()
	t.Logf("ICMON_MQTT_BROKER=1999 + ICMON_MQTT_PORT=1999 -> effective broker=%s", got)
	if got != "tcp://localhost:1999" {
		t.Fatalf("expected tcp://localhost:1999, got %s", got)
	}

	// An explicit ICMON_MQTT_PORT overrides the broker's embedded port.
	pinEnv("tcp://localhost:1999", "1884")
	got = load()
	t.Logf("ICMON_MQTT_BROKER=1999 + ICMON_MQTT_PORT=1884 -> effective broker=%s", got)
	if got != "tcp://localhost:1884" {
		t.Fatalf("expected tcp://localhost:1884, got %s", got)
	}
}