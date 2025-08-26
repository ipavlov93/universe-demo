package env_test

import (
	"os"
	"testing"

	"github.com/ipavlov93/universe-demo/universe-pkg/env"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentVariable(t *testing.T) {
	const key = "TEST_ENV_STRING"
	const fallback = "default"

	t.Run("should return value when env is set", func(t *testing.T) {
		_ = os.Setenv(key, "hello")
		defer os.Unsetenv(key)

		got := env.EnvironmentVariable(key, fallback)
		if got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})

	t.Run("should return fallback when env is not set", func(t *testing.T) {
		_ = os.Unsetenv(key)

		got := env.EnvironmentVariable(key, fallback)
		if got != fallback {
			t.Errorf("got %q, want fallback %q", got, fallback)
		}
	})
}

func TestParseIntEnv(t *testing.T) {
	const key = "TEST_ENV_INT"
	const fallback = 42

	t.Run("should return parsed int when env is valid", func(t *testing.T) {
		_ = os.Setenv(key, "123")
		defer os.Unsetenv(key)

		expected := 123
		assert.EqualValues(t, expected, env.ParseIntEnv(key, fallback))
	})

	t.Run("should return fallback when env is not set", func(t *testing.T) {
		_ = os.Unsetenv(key)

		assert.Equal(t, fallback, env.ParseIntEnv(key, fallback))
	})

	t.Run("should return fallback when env is invalid", func(t *testing.T) {
		_ = os.Setenv(key, "not-an-int")
		defer os.Unsetenv(key)

		assert.Equal(t, fallback, env.ParseIntEnv(key, fallback))
	})
}
