package main_test

import (
	"testing"

	main "gitlab.local/hassio-addons/mq-lightshow"
)

func TestGetLogger(t *testing.T) {
	t.Parallel()

	testLogLevel := "error"

	logger := main.GetLogger(testLogLevel)
	if logger == nil {
		t.Errorf("GetLogger did not return a valid logger")
	}
}
