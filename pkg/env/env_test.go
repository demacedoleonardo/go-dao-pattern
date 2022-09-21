package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvironmentDefault(t *testing.T) {
	assert.NotNil(t, "dev", Environment())
}

func TestGetEnvironment(t *testing.T) {
	empty := ""
	prod := "prod"
	assert.Nil(t, os.Setenv("ENV", prod))
	assert.NotNil(t, prod, Environment())
	assert.Nil(t, os.Setenv("ENV", empty))
}
