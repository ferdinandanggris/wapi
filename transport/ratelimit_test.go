package transport_test

import (
	"testing"
	"time"

	"github.com/ferdinandanggris/wapi/transport"
)

func TestTokenBucket_AllowsWithinLimit(t *testing.T) {
	tb := transport.NewTokenBucket(10, 5)

	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("expected allow at token %d", i)
		}
	}
}

func TestTokenBucket_BlocksWhenEmpty(t *testing.T) {
	tb := transport.NewTokenBucket(1, 1)

	if !tb.Allow() {
		t.Fatal("expected first call to allow")
	}
	if tb.Allow() {
		t.Error("expected second call to block immediately")
	}
}

func TestTokenBucket_RefillsOverTime(t *testing.T) {
	tb := transport.NewTokenBucket(100, 1)

	if !tb.Allow() {
		t.Fatal("expected initial token")
	}
	if tb.Allow() {
		t.Fatal("expected block after using burst")
	}

	time.Sleep(15 * time.Millisecond)

	if !tb.Allow() {
		t.Error("expected token after refill")
	}
}

func TestTokenBucket_BurstCapacity(t *testing.T) {
	tb := transport.NewTokenBucket(1, 5)

	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("expected allow for burst token %d", i)
		}
	}
	if tb.Allow() {
		t.Error("expected block after burst exhausted")
	}
}
