package messaging

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetUp()

	result := m.Run()

	TearDown()
	os.Exit(result)
}
