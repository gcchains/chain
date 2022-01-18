package main

import (
	"testing"

	"math/big"

	"github.com/stretchr/testify/assert"
)

func TestIsInvalidAddress(t *testing.T) {
	assert.Equal(t, isInvalidAddress("e94b7b6c5a0e521a4d97f9768ad6097bd125c62a"), false)
	assert.Equal(t, isInvalidAddress("0xe94b7b6c1a0e526a4197f9768ad6097bde25c62a"), false)
	assert.Equal(t, isInvalidAddress("0e14b7b6c5a0e526a4d97f9768ad6097bde21c62a"), true)
	assert.Equal(t, isInvalidAddress("1e94b7b1c5a0e526a4d97f9768ad60971de25c62a"), true)
}

func TestFormatNumber1(t *testing.T) {
	assert.Equal(t, "100,200,300", formatNumber(big.NewInt(100200300)))
	assert.Equal(t, "100,200,300,123", formatNumber(big.NewInt(100200300123)))
}

func TestFormatNumber2(t *testing.T) {
	assert.Equal(t, "12,345", formatNumber(big.NewInt(12345)))
}

func TestFormatNumber3(t *testing.T) {
	assert.Equal(t, "9,999,244,000", formatNumber(big.NewInt(9999244000)))
}

func TestFormatNumber4(t *testing.T) {
	assert.Equal(t, "82,999,244,000", formatNumber(big.NewInt(82999244000)))
}
