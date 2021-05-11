package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionType_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		at   ActionType
		want error
	}{
		{at: ActionNone, want: ErrInvalidActionType},
		{at: ActionHTTP, want: nil},
		{at: ActionType("invalid"), want: ErrInvalidActionType},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			err := tt.at.Validate()
			assert.ErrorIs(t, err, tt.want)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		c    *Config
		want error
	}{
		{c: &Config{}, want: errors.New("version is required")},
		{
			c: &Config{
				Version: "1",
			},
			want: errors.New("handler should be defined one or more"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			err := tt.c.Validate()
			if tt.want == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.want.Error(), err.Error())
			}
		})
	}
}
