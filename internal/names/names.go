// Package names creates random names within bounded length ranges.
package names

import (
	crand "crypto/rand"
	"math/big"
	mrand "math/rand"
	"time"
)

// Generator produces random ASCII names within a bounded length range.
type Generator struct {
	minLen  int
	maxLen  int
	charset []rune
	src     *mrand.Rand
}

// New creates a Generator. Names will be between min and max inclusive.
func New(minLen, maxLen int, charset []rune) Generator {
	seed, err := crand.Int(crand.Reader, big.NewInt(1<<62))
	if err != nil {
		seed = big.NewInt(time.Now().UnixNano())
	}

	return Generator{
		minLen:  minLen,
		maxLen:  maxLen,
		charset: charset,
		src:     mrand.New(mrand.NewSource(seed.Int64())), //nolint:gosec // deterministic pseudo-random is sufficient
	}
}

// String returns a random string.
func (g Generator) String() string {
	length := g.minLen
	if g.maxLen > g.minLen {
		length += g.src.Intn(g.maxLen - g.minLen + 1)
	}

	runes := make([]rune, length)
	for i := range runes {
		runes[i] = g.charset[g.src.Intn(len(g.charset))]
	}

	return string(runes)
}

// AllowedCharset returns a broad ASCII charset that includes alphanumerics
// and common punctuation to satisfy the requirement of not avoiding
// characters beyond ASCII alphanumerics.
func AllowedCharset() []rune {
	chars := []rune{}
	for r := '0'; r <= '9'; r++ {
		chars = append(chars, r)
	}
	for r := 'a'; r <= 'z'; r++ {
		chars = append(chars, r)
	}
	for r := 'A'; r <= 'Z'; r++ {
		chars = append(chars, r)
	}
	punctuation := "-_+.@#" // stay filesystem-friendly while exceeding alphanumerics
	for _, r := range punctuation {
		chars = append(chars, r)
	}
	return chars
}
