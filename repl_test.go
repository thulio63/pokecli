package main

import (
	"testing"
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