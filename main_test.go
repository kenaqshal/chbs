package main

import (
	"strings"
	"testing"
)

// All lowercase ASCII so capitalization and separator handling stay observable.
var testWords = []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel"}

func TestGenerateHonoursMinimums(t *testing.T) {
	for range 300 {
		got := generate(testWords, 4, 20, "-", true, true)
		parts := strings.Split(got, "-")

		if len(parts) < 5 {
			t.Fatalf("want at least 4 words plus the digit, got %q", got)
		}
		digit := parts[len(parts)-1]
		if len(digit) != 1 || digit[0] < '0' || digit[0] > '9' {
			t.Fatalf("want a single trailing digit, got %q in %q", digit, got)
		}
		for _, p := range parts[:len(parts)-1] {
			if p == "" || p[0] < 'A' || p[0] > 'Z' {
				t.Fatalf("word %q not capitalized in %q", p, got)
			}
		}
		if body := strings.TrimSuffix(got, "-"+digit); len(body) < 20 {
			t.Fatalf("want a body of at least 20 chars, got %d in %q", len(body), got)
		}
	}
}

func TestGenerateWithoutDigitsOrUpper(t *testing.T) {
	got := generate(testWords, 3, 1, ".", false, false)
	if got != strings.ToLower(got) {
		t.Fatalf("unexpected capitals in %q", got)
	}
	if strings.ContainsAny(got, "0123456789") {
		t.Fatalf("unexpected digit in %q", got)
	}
	if parts := strings.Split(got, "."); len(parts) != 3 {
		t.Fatalf("want exactly 3 words at -l 1, got %d in %q", len(parts), got)
	}
}

// -l is a floor, so a long phrase must pull in more words than -n asks for.
func TestGenerateLengthRaisesWordCount(t *testing.T) {
	got := generate(testWords, 1, 60, "-", false, false)
	if n := len(strings.Split(got, "-")); n < 8 {
		t.Fatalf("want the length floor to add words, got %d in %q", n, got)
	}
}

func TestGenerateEmptySeparator(t *testing.T) {
	got := generate(testWords, 4, 20, "", false, false)
	if strings.ContainsAny(got, "-.") {
		t.Fatalf("unexpected separator in %q", got)
	}
	if len(got) < 20 {
		t.Fatalf("want a phrase of at least 20 chars, got %q", got)
	}
}

func TestLoadWords(t *testing.T) {
	primary, err := loadWords("")
	if err != nil {
		t.Fatal(err)
	}
	if len(primary) < 1000 {
		t.Fatalf("primary list looks truncated: %d words", len(primary))
	}

	extra, err := loadWords(" Jargon , science ")
	if err != nil {
		t.Fatal(err)
	}
	if len(extra) <= len(primary) {
		t.Fatalf("extra dictionaries added no words: %d vs %d", len(extra), len(primary))
	}

	if _, err := loadWords("nope"); err == nil {
		t.Fatal("want an error for an unknown dictionary")
	}
}

func TestRandIntStaysInRange(t *testing.T) {
	for range 1000 {
		if got := randInt(10); got < 0 || got >= 10 {
			t.Fatalf("randInt(10) = %d, out of range", got)
		}
	}
}
