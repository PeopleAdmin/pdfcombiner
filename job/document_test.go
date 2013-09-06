package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocValidation(t *testing.T) {
	invalidDoc := &Document{}
	validDoc := &Document{Key: "Something"}
	assert.False(t, invalidDoc.isValid(), "Invalid Doc should not validate")
	assert.True(t, validDoc.isValid(), "Valid Doc should validate")
}
