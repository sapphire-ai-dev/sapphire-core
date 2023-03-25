package text

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertEmptyNotNil(t *testing.T, a any) {
	assert.Empty(t, a)
	assert.NotNil(t, a)
}
