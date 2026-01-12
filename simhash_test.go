package main

import (
	"testing"
)

func TestNewSimHash_Basic(t *testing.T) {
	text := "hello world"
	hash := NewSimHash(text)
	
	if hash == 0 {
		t.Errorf("hash should not be zero")
	}
}

func TestNewSimHash_Deterministic(t *testing.T) {
	text := "the quick brown fox"
	hash1 := NewSimHash(text)
	hash2 := NewSimHash(text)
	
	if hash1 != hash2 {
		t.Errorf("same text should produce same hash: got %d and %d", hash1, hash2)
	}
}

func TestNewSimHash_DifferentTexts(t *testing.T) {
	text1 := "hello world"
	text2 := "goodbye universe"
	hash1 := NewSimHash(text1)
	hash2 := NewSimHash(text2)
	
	if hash1 == hash2 {
		t.Errorf("different texts should produce different hashes")
	}
}

func TestNewSimHash_SimilarTexts(t *testing.T) {
	text1 := "the cat sat on the mat"
	text2 := "the cat sat on the rug"
	hash1 := NewSimHash(text1)
	hash2 := NewSimHash(text2)
	
	distance := HammingDistance(hash1, hash2)
	if distance > 30 {
		t.Errorf("similar texts should have small hamming distance, got %d", distance)
	}
}

func TestNewSimHash_WithStopWords(t *testing.T) {
	text := "the cat sat on the mat"
	hash1 := NewSimHash(text)
	hash2 := NewSimHash(text, WithStopWords("the", "on"))
	
	if hash1 == hash2 {
		t.Errorf("stop words should affect hash")
	}
}

func TestNewSimHash_WithMinTokenLength(t *testing.T) {
	text := "a an the cat sat"
	hash1 := NewSimHash(text)
	hash2 := NewSimHash(text, WithMinTokenLength(3))
	
	if hash1 == hash2 {
		t.Errorf("min token length should affect hash")
	}
}

func TestNewSimHash_WithNgram(t *testing.T) {
	text := "hello world"
	hash1 := NewSimHash(text)
	hash2 := NewSimHash(text, WithNgram(4))
	
	if hash1 == hash2 {
		t.Errorf("ngram mode should produce different hash")
	}
}

func TestNewSimHash_NgramWidth(t *testing.T) {
	text := "hello world test"
	hash3 := NewSimHash(text, WithNgram(3))
	hash4 := NewSimHash(text, WithNgram(4))
	hash5 := NewSimHash(text, WithNgram(5))
	
	if hash3 == hash4 || hash4 == hash5 || hash3 == hash5 {
		t.Errorf("different ngram widths should produce different hashes")
	}
}

func TestNewSimHash_CombinedOptions(t *testing.T) {
	text := "the cat sat on the mat"
	hash1 := NewSimHash(text)
	hash2 := NewSimHash(text,
		WithStopWords("the", "on"),
		WithMinTokenLength(3),
	)
	
	if hash1 == hash2 {
		t.Errorf("combined options should affect hash")
	}
}

func TestNewSimHash_EmptyString(t *testing.T) {
	hash := NewSimHash("")
	
	if hash != 0 {
		t.Errorf("empty string should produce zero hash, got %d", hash)
	}
}

func TestNewSimHash_WhitespaceOnly(t *testing.T) {
	hash := NewSimHash("   \n\t  ")
	
	if hash != 0 {
		t.Errorf("whitespace only should produce zero hash, got %d", hash)
	}
}

func TestNewSimHash_CaseInsensitive(t *testing.T) {
	text1 := "Hello World"
	text2 := "hello world"
	hash1 := NewSimHash(text1)
	hash2 := NewSimHash(text2)
	
	if hash1 != hash2 {
		t.Errorf("case should not affect hash: got %d and %d", hash1, hash2)
	}
}

func TestHammingDistance_Identical(t *testing.T) {
	hash := uint64(12345)
	distance := HammingDistance(hash, hash)
	
	if distance != 0 {
		t.Errorf("identical hashes should have distance 0, got %d", distance)
	}
}

func TestHammingDistance_Different(t *testing.T) {
	hash1 := uint64(0)
	hash2 := uint64(^uint64(0))
	distance := HammingDistance(hash1, hash2)
	
	if distance != 64 {
		t.Errorf("opposite hashes should have distance 64, got %d", distance)
	}
}

func TestHammingDistance_RealWorld(t *testing.T) {
	text1 := "the cat sat on the mat"
	text2 := "the cat sat on the rug"
	hash1 := NewSimHash(text1)
	hash2 := NewSimHash(text2)
	distance := HammingDistance(hash1, hash2)
	
	if distance < 0 || distance > 64 {
		t.Errorf("hamming distance should be between 0 and 64, got %d", distance)
	}
	
	text3 := "completely different text about something else"
	hash3 := NewSimHash(text3)
	distance13 := HammingDistance(hash1, hash3)
	distance23 := HammingDistance(hash2, hash3)
	
	if distance >= distance13 || distance >= distance23 {
		t.Errorf("similar texts should have smaller distance than different texts")
	}
}

func TestNewSimHash_SpecialCharacters(t *testing.T) {
	text1 := "hello world"
	text2 := "hello, world!"
	hash1 := NewSimHash(text1)
	hash2 := NewSimHash(text2)
	
	if hash1 != hash2 {
		t.Errorf("special characters should be ignored, got %d and %d", hash1, hash2)
	}
}

func TestNewSimHash_Numbers(t *testing.T) {
	text1 := "test 123"
	text2 := "test 123"
	hash1 := NewSimHash(text1)
	hash2 := NewSimHash(text2)
	
	if hash1 != hash2 {
		t.Errorf("same text with numbers should produce same hash, got %d and %d", hash1, hash2)
	}
}

func TestNewSimHash_LongText(t *testing.T) {
	text := "this is a very long text that contains many words and should produce a valid hash even though it is quite lengthy and has many tokens"
	hash := NewSimHash(text)
	
	if hash == 0 {
		t.Errorf("long text should produce non-zero hash")
	}
}

func TestNewSimHash_RepeatedWords(t *testing.T) {
	text1 := "cat cat cat dog"
	text2 := "cat dog"
	hash1 := NewSimHash(text1)
	hash2 := NewSimHash(text2)
	
	distance := HammingDistance(hash1, hash2)
	if distance == 0 {
		t.Errorf("repeated words should affect hash weight, distance should not be zero")
	}
}

func TestNewSimHash_NgramShortText(t *testing.T) {
	text := "abc"
	hash := NewSimHash(text, WithNgram(4))
	
	if hash != 0 {
		t.Errorf("text shorter than ngram width should produce zero hash, got %d", hash)
	}
	
	text2 := "abcdef"
	hash2 := NewSimHash(text2, WithNgram(4))
	if hash2 == 0 {
		t.Errorf("text longer than ngram width should produce non-zero hash")
	}
}

func BenchmarkNewSimHash(b *testing.B) {
	text := "the quick brown fox jumps over the lazy dog"
	for i := 0; i < b.N; i++ {
		NewSimHash(text)
	}
}

func BenchmarkNewSimHash_WithOptions(b *testing.B) {
	text := "the quick brown fox jumps over the lazy dog"
	for i := 0; i < b.N; i++ {
		NewSimHash(text,
			WithStopWords("the", "a", "an"),
			WithMinTokenLength(3),
		)
	}
}

func BenchmarkNewSimHash_Ngram(b *testing.B) {
	text := "the quick brown fox jumps over the lazy dog"
	for i := 0; i < b.N; i++ {
		NewSimHash(text, WithNgram(4))
	}
}

func BenchmarkHammingDistance(b *testing.B) {
	hash1 := NewSimHash("the cat sat on the mat")
	hash2 := NewSimHash("the cat sat on the rug")
	for i := 0; i < b.N; i++ {
		HammingDistance(hash1, hash2)
	}
}

