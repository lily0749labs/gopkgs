package jwtUtil

import (
	"errors"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func TestParseDuration(t *testing.T) {
	testCases := map[string]time.Duration{
		"2h":      2 * time.Hour,
		"7d":      7 * 24 * time.Hour,
		"1d12h":   36 * time.Hour,
		"1000000": time.Millisecond,
	}
	for input, want := range testCases {
		got, err := ParseDuration(input)
		if err != nil || got != want {
			t.Fatalf("ParseDuration(%q) = %v, %v; want %v", input, got, err, want)
		}
	}
}

func TestTokenRoundTrip(t *testing.T) {
	SigningKey = "test-signing-key"
	ExpiresTime = "2h"
	BufferTime = "15m"
	Issuer = "test-issuer"

	fixedNow := time.Now().Add(-time.Minute).Truncate(time.Second)
	manager := NewJWT()
	manager.now = func() time.Time { return fixedNow }
	claims := manager.CreateClaims(BaseClaims{ID: 42, Account: "andy", Roles: "admin"})

	token, err := manager.CreateToken(claims)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := manager.ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.BaseClaims.ID != 42 || parsed.Account != "andy" || parsed.Roles != "admin" {
		t.Fatalf("parsed claims = %+v", parsed.BaseClaims)
	}
	if parsed.BufferTime != int64((15*time.Minute)/time.Second) {
		t.Fatalf("BufferTime = %d", parsed.BufferTime)
	}
	if !parsed.ExpiresAt.Time.Equal(fixedNow.Add(2 * time.Hour)) {
		t.Fatalf("ExpiresAt = %v", parsed.ExpiresAt.Time)
	}
}

func TestExpiredToken(t *testing.T) {
	SigningKey = "test-signing-key"
	manager := NewJWT()
	claims := CustomClaims{RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
	}}
	token, err := manager.CreateToken(claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := manager.ParseToken(token); !errors.Is(err, TokenExpired) {
		t.Fatalf("ParseToken() error = %v, want TokenExpired", err)
	}
}
