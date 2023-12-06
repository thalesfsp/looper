package looper_test

import (
	"context"
	"testing"
	"time"

	"github.com/thalesfsp/looper/looper"
)

func TestNew(t *testing.T) {
	// Create a string channel. We need that so we can read the metrics from the
	// looper and hold the test until the metrics are written.
	ch := make(chan string)

	// Close the channel when the test is done.
	defer close(ch)

	// Create a simple looper, just for testing.
	l, err := looper.New(
		"test",
		1*time.Second,
		false,
		"none",
		"9000", // required even if metrics isn't enabled. To be improved later.
		func(metrics looper.IMetrics) {
			// Write to the channel.
			ch <- "worked"
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Stop the looper when the test is done.
	defer l.Stop()

	// Start the looper in a goroutine.
	l.StartAsync()

	// Create a 5 seconds context to hold the test.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// Cancel the context when the test is done.
	defer cancel()

	// Reads from the channel until the context is done.
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			if msg != "worked" {
				t.Fatal("expected worked")
			}
		}
	}
}
