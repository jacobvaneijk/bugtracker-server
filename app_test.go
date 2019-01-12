package main

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
    app := NewApp()

    assert.NotNil(t, app)
}
