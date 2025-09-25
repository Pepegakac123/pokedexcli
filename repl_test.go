package main

import (
	"testing"
)



func TestCleanInput(t *testing.T){
cases := []struct {
	input    string
	expected []string
}{
	{
		input:    "  hello  world  ",
		expected: []string{"hello", "world"},
	},
	{
		input:    "  Charmander Bulbasaur PIKACHU  ",
		expected: []string{"charmander", "bulbasaur", "pikachu"},
	},
	{
		input:    "",
		expected: []string{},
	},
	{
		input:    "   ",
		expected: []string{},
	},
	{
		input:    "Go\tis\nAwesome",
		expected: []string{"go", "is", "awesome"},
	},
	{
		input:    "SingleWord",
		expected: []string{"singleword"},
	},
	{
		input:    "  Multiple    Spaces     Between    Words  ",
		expected: []string{"multiple", "spaces", "between", "words"},
	},
	// add more cases here
}

	for _, c := range cases{
		actual := cleanInput(c.input)
		if(len(actual) != len(c.expected)){
			t.Errorf("Expected %v, got %v", c.expected, actual)
		}
		for i := range actual{
			word:= actual[i]
			expectedWord := c.expected[i]
			if(word != expectedWord){
				t.Errorf("Expected %v, got %v", expectedWord, word)
			}
		}
	}
}