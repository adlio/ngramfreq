package main

import "testing"

func TestNGramFreqString(t *testing.T) {
	ngf := &NGramFreq{Text: "message number one", Freq: 45}
	s := ngf.String()
	expected := "45 - message number one"
	if s != expected {
		t.Errorf("Expected '%s', got '%s'", expected, s)
	}
}
