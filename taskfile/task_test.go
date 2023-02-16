package taskfile

import (
	"testing"
	"time"
)

func TestNextExecutionTime(t *testing.T) {
	start := time.Now()
	yesterday := start.Add(-24 * time.Hour)
	checkpointStart := start.Add(-2 * time.Minute)
	checkpointPast := start.Add(-2 * time.Hour)
	roundingDuration := 5 * time.Second
	tests := []struct {
		Name       string
		Start      time.Time
		Interval   time.Duration
		Expected   time.Time
		Checkpoint *time.Time
	}{
		{Name: "Now", Start: start, Interval: time.Hour, Expected: start},
		{Name: "Basic future", Start: start.Add(1 * time.Hour), Interval: 2 * time.Hour, Expected: start.Add(1 * time.Hour)},
		{Name: "Past due", Start: yesterday, Interval: 2 * time.Hour, Expected: start},
		{Name: "Recent checkpoint", Start: checkpointStart, Interval: time.Hour, Expected: checkpointStart.Add(time.Hour), Checkpoint: &checkpointStart},
		{Name: "Past due checkpoint", Start: start.Add(-3 * time.Hour), Interval: time.Hour, Expected: start, Checkpoint: &checkpointPast},
	}

	for _, test := range tests {
		task := Task{Name: test.Name, Start: TaskTime{test.Start}, Interval: test.Interval}
		if test.Checkpoint != nil {
			task.Checkpoint = &TaskTime{*test.Checkpoint}
		}

		nextExecutionTime := task.NextExecutionTime().Round(roundingDuration)
		if !nextExecutionTime.Equal(test.Expected.Round(roundingDuration)) {
			t.Errorf("expected next execution of %s to be %s, got %s", test.Name, test.Expected.Format(TimeFormat), nextExecutionTime.Format(TimeFormat))
		}
	}
}
