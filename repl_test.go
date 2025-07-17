package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/thulio63/pokecli/internal"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "   hello    world ",
			expected: []string{"hello","world"},
		},
		//more cases
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		//fails test if lengths don't match
		if len(actual) != len(c.expected) {
			t.Errorf("Length of slice should be %d, but returned %d.", len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			//fails test if slices don't match
			if word != expectedWord {
				t.Errorf("Word number %d in slice should be %s, but returned %s.", i, expectedWord, word)
			}
		}
	}
}

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
		//more cases
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := internal.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := internal.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
	}

	time.Sleep(waitTime)
	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
	}
}