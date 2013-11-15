package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultStates(t *testing.T) {
	j, _ := newFromString(ValidJSON)
	assert.True(t, j.Source.IsDefault())
	assert.False(t, j.Source.IsOverflow())
}

func TestSourceToggling(t *testing.T) {
	j, _ := newFromString(ValidJSON)
	assert.True(t, j.Source.IsDefault())
	assert.False(t, j.Source.IsOverflow())
	j.Source.SetOverflow()
	assert.False(t, j.Source.IsDefault())
	assert.True(t, j.Source.IsOverflow())
	j.Source.SetDefault()
	assert.True(t, j.Source.IsDefault())
	assert.False(t, j.Source.IsOverflow())
}
