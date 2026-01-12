package main

import (
	"crypto/md5"
	"encoding/binary"
	"math/bits"
	"slices"
	"strings"
	"unicode"
)

const fingerPrintSize = 64

type simHash struct {
	stopWords      []string
	minTokenLength int
	useNgram       bool
	ngramWidth     int
}

type feature struct {
	value  string
	weight int
}

type Option func(*simHash)

func WithStopWords(words ...string) Option {
	return func(s *simHash) {
		s.stopWords = words
	}
}

func WithMinTokenLength(length int) Option {
	return func(s *simHash) {
		s.minTokenLength = length
	}
}

func WithNgram(width int) Option {
	return func(s *simHash) {
		s.useNgram = true
		s.ngramWidth = width
	}
}

func (f *feature) hash() uint64 {
	h := md5.New()
	h.Write([]byte(f.value))
	sum := h.Sum(nil)
	return binary.BigEndian.Uint64(sum[:8])
}

func (s *simHash) tokenize(text string) []string {
	var tokens []string
	var current strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (s *simHash) extractTokens(text string) []feature {
	tokens := s.tokenize(text)

	m := make(map[string]int)
	for _, token := range tokens {
		token = strings.ToLower(token)

		if !slices.Contains(s.stopWords, token) &&
			len(token) > s.minTokenLength {
			m[token]++
		}
	}

	f := make([]feature, 0, len(m))
	for token, count := range m {
		f = append(f, feature{
			value:  token,
			weight: count,
		})
	}

	return f
}

func (s *simHash) extractNgrams(text string) []feature {
	width := s.ngramWidth
	if width <= 0 {
		width = 4
	}

	m := make(map[string]int)
	for i := 0; i <= len(text)-width; i++ {
		ngram := text[i : i+width]
		m[ngram]++
	}

	f := make([]feature, 0, len(m))
	for ngram, count := range m {
		f = append(f, feature{
			value:  ngram,
			weight: count,
		})
	}

	return f
}

func (s *simHash) extractor(text string) []feature {
	if s.useNgram {
		return s.extractNgrams(text)
	}
	return s.extractTokens(text)
}

func (s *simHash) compute(features []feature) uint64 {
	v := make([]int, fingerPrintSize)

	for _, f := range features {
		hash := f.hash()
		for i := 0; i < fingerPrintSize; i++ {
			bit := (hash >> i) & 1
			if bit == 1 {
				v[i] += f.weight
			} else {
				v[i] -= f.weight
			}
		}
	}

	var sh uint64
	for i := 0; i < fingerPrintSize; i++ {
		if v[i] > 0 {
			sh |= (1 << i)
		}
	}

	return sh
}

func NewSimHash(text string, opts ...Option) uint64 {
	s := &simHash{
		minTokenLength: 2,
		ngramWidth:     4,
	}

	for _, opt := range opts {
		opt(s)
	}

	features := s.extractor(text)
	return s.compute(features)
}

func HammingDistance(a, b uint64) int {
	return bits.OnesCount64(a ^ b)
}

func AreDocumentsSimilar(text1, text2 string, k int) bool {
	hash1 := NewSimHash(text1)
	hash2 := NewSimHash(text2)

	return HammingDistance(hash1, hash2) <= k
}
