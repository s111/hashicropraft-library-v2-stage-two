package raft

import (
	"regexp"
	"testing"
	"time"
)

func TestRandomTimeout(t *testing.T) {
	start := time.Now()
	timeout := randomTimeout(time.Millisecond)

	select {
	case <-timeout:
		diff := time.Now().Sub(start)
		if diff < time.Millisecond {
			t.Fatalf("fired early")
		}
	case <-time.After(3 * time.Millisecond):
		t.Fatalf("timeout")
	}
}

func TestNewSeed(t *testing.T) {
	vals := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		seed := newSeed()
		if _, exists := vals[seed]; exists {
			t.Fatal("newSeed() return a value it'd previously returned")
		}
		vals[seed] = true
	}
}

func TestRandomTimeout_NoTime(t *testing.T) {
	timeout := randomTimeout(0)
	if timeout != nil {
		t.Fatalf("expected nil channel")
	}
}

func TestMin(t *testing.T) {
	if min(1, 1) != 1 {
		t.Fatalf("bad min")
	}
	if min(2, 1) != 1 {
		t.Fatalf("bad min")
	}
	if min(1, 2) != 1 {
		t.Fatalf("bad min")
	}
}

func TestMax(t *testing.T) {
	if max(1, 1) != 1 {
		t.Fatalf("bad max")
	}
	if max(2, 1) != 2 {
		t.Fatalf("bad max")
	}
	if max(1, 2) != 2 {
		t.Fatalf("bad max")
	}
}

func TestGenerateUUID(t *testing.T) {
	prev := generateUUID()
	for i := 0; i < 100; i++ {
		id := generateUUID()
		if prev == id {
			t.Fatalf("Should get a new ID!")
		}

		matched, err := regexp.MatchString(
			`[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}`, id)
		if !matched || err != nil {
			t.Fatalf("expected match %s %v %s", id, matched, err)
		}
	}
}

type backoffTest struct {
	round    uint64
	base     time.Duration
	limit    time.Duration
	expected time.Duration
}

const ms = time.Millisecond

var backoffTests = []backoffTest{
	{0, 10 * ms, 80 * ms, 0 * ms},
	{1, 10 * ms, 80 * ms, 10 * ms},
	{2, 10 * ms, 80 * ms, 20 * ms},
	{3, 10 * ms, 80 * ms, 40 * ms},
	{4, 10 * ms, 80 * ms, 80 * ms},
	{5, 10 * ms, 80 * ms, 80 * ms},
	{3, 2 * ms, 9 * ms, 8 * ms},
	{4, 2 * ms, 9 * ms, 9 * ms},
}

func TestBackoff(t *testing.T) {
	for i, test := range backoffTests {
		actual := backoff(test.round, test.base, test.limit)
		if actual != test.expected {
			t.Errorf("backoff(%v, %v, %v) = %v, expected %v (test %v)",
				test.round, test.base, test.limit, actual, test.expected, i)
		}
	}
}

func TestEnsureClosed_basic(t *testing.T) {
	ch := make(chan struct{})
	ensureClosed(ch)
	ensureClosed(ch)
	ensureClosed(ch)
	select {
	case <-ch:
		// ok
	default:
		t.Errorf("Not closed")
	}
}

func TestEnsureClosed_panicOnNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("ensureClosed(nil) should panic but didn't")
		}
	}()
	ensureClosed(nil)
}
