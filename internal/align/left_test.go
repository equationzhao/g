package align

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLeft(t *testing.T) {
	Register("name")
	assert.Equal(t, true, IsLeft("name"))
	assert.Equal(t, false, IsLeft("name1"))
}

func TestIsLeftHeaderFooter(t *testing.T) {
	RegisterHeaderFooter("name")
	assert.Equal(t, true, IsLeftHeaderFooter("name"))
	assert.Equal(t, false, IsLeftHeaderFooter("name1"))
}
