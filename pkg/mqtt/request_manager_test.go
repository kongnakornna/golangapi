package mqtt

import (
	"context"
	"fmt"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type testLogger struct{}

func (testLogger) InitLogger()                                        {}
func (testLogger) Debug(args ...interface{})                          {}
func (testLogger) Debugf(template string, args ...interface{})        {}
func (testLogger) Info(args ...interface{})                           {}
func (testLogger) Infof(template string, args ...interface{})         {}
func (testLogger) Warn(args ...interface{})                           {}
func (testLogger) Warnf(template string, args ...interface{})         {}
func (testLogger) Error(args ...interface{})                          {}
func (testLogger) Errorf(template string, args ...interface{})        {}
func (testLogger) DPanic(args ...interface{})                         {}
func (testLogger) DPanicf(template string, args ...interface{})       {}
func (testLogger) Fatal(args ...interface{})                          {}
func (testLogger) Fatalf(template string, args ...interface{})        {}
func (testLogger) Sync() error                                        { return nil }

// TestGetTopicResolvesAfterWildcardRouteReplacement reproduces the bug where the IoT
// ingester's Subscribe("#", 0, handler) replaces the requestManager's global "#" route
// (paho keeps a single callback per topic). Previously getTopicInternal subscribed the
// wait topic with a nil callback, so the incoming message only reached the ingester and
// every pending request timed out. The wait topic must now be subscribed with the
// requestManager's own handler so resolution no longer depends on the "#" route.
func TestGetTopicResolvesAfterWildcardRouteReplacement(t *testing.T) {
	opts := paho.NewClientOptions().
		AddBroker("tcp://localhost:1885").
		SetClientID(fmt.Sprintf("icm-repro-%d", time.Now().UnixNano())).
		SetConnectTimeout(2 * time.Second)
	c := paho.NewClient(opts)
	tok := c.Connect()
	if !tok.WaitTimeout(3 * time.Second) {
		t.Skip("broker not reachable at localhost:1885")
	}
	if tok.Error() != nil {
		t.Skipf("broker not reachable at localhost:1885: %v", tok.Error())
	}
	defer c.Disconnect(250)

	rm := newRequestManager(c, false, testLogger{})

	// Simulate StartIngest: subscribing "#" with a callback replaces the requestManager's
	// global "#" route inside paho's router.
	if tok := c.Subscribe("#", 0, func(paho.Client, paho.Message) {}); !tok.WaitTimeout(5*time.Second) || tok.Error() != nil {
		t.Fatalf("subscribe # failed: %v", tok.Error())
	}

	topic := fmt.Sprintf("REPRO-%d/DATA", time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	type result struct {
		payload []byte
		err     error
	}
	done := make(chan result, 1)
	go func() {
		data, err := rm.getTopicRaw(ctx, topic, 2000)
		done <- result{data, err}
	}()

	// Give the async SUBSCRIBE time to reach the broker before publishing.
	time.Sleep(400 * time.Millisecond)
	if tok := c.Publish(topic, 0, false, "hello-world"); !tok.WaitTimeout(2*time.Second) || tok.Error() != nil {
		t.Fatalf("publish failed: %v", tok.Error())
	}

	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("getTopicRaw failed: %v", r.err)
		}
		if string(r.payload) != "hello-world" {
			t.Fatalf("unexpected payload: %q", r.payload)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("request timed out: pending request never resolved after # route was replaced")
	}
}
