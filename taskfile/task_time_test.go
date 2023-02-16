package taskfile

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTaskTimeMarshal(t *testing.T) {
	startDateStr := "17 Jan 23 17:13:05 UTC"
	startDate, _ := time.Parse(TimeFormat, startDateStr)
	start := &TaskTime{startDate}

	b, err := json.Marshal(start)
	if err != nil {
		t.Error(err)
	}

	if string(b) != fmt.Sprintf("\"%s\"", startDateStr) {
		t.Errorf("expected time to marshal to %s, got %s", startDateStr, string(b))
	}
}
