// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLowerCamel(t *testing.T) {
	assert.Equal(t, "", lowerCamel(""))
	assert.Equal(t, "a", lowerCamel("a"))
	assert.Equal(t, "a", lowerCamel("A"))
	assert.Equal(t, "aa", lowerCamel("AA"))
	assert.Equal(t, "asdfTest", lowerCamel("asdfTest"))
	assert.Equal(t, "asdfTest", lowerCamel("AsdfTest"))
	assert.Equal(t, "asdf", lowerCamel("asdf"))
}
