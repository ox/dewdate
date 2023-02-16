package taskfile

import (
	"dewdate/taskfile/runner"
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	Name       string
	Start      TaskTime
	Interval   time.Duration
	Checkpoint *TaskTime `json:",omitempty"`
	Runners    map[string]map[string]interface{}
}

type task_json struct {
	Name       string
	Start      TaskTime
	Interval   string
	Checkpoint *TaskTime `json:",omitempty"`
	Runners    map[string]map[string]interface{}
}

func (t *Task) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(&task_json{Name: t.Name, Start: t.Start, Interval: t.Interval.String(), Checkpoint: t.Checkpoint, Runners: t.Runners})
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (t *Task) UnmarshalJSON(data []byte) error {
	nt := &task_json{}
	if err := json.Unmarshal(data, nt); err != nil {
		return err
	}

	var err error
	t.Name = nt.Name
	t.Start = nt.Start
	t.Interval, err = time.ParseDuration(nt.Interval)
	t.Checkpoint = nt.Checkpoint
	t.Runners = nt.Runners

	return err
}

func (t *Task) String() string {
	checkpoint := ""
	if t.Checkpoint != nil {
		checkpoint = fmt.Sprintf(" last checkpointed at '%s'", t.Checkpoint.Format(TimeFormat))
	}
	return fmt.Sprintf("<Task '%s' every '%s' starting '%s'%s>", t.Name, t.Interval, t.Start.Format(TimeFormat), checkpoint)
}

func (t *Task) NextExecutionTime() time.Time {
	// If the task has never run...
	if t.Checkpoint == nil {
		// And it's past due, start now
		if time.Now().After(t.Start.Time) {
			return time.Now()
		} else {
			// Otherwise start at the start time
			return t.Start.Time
		}
	}

	// Calculate the next time it should run
	// TODO: task.ForceCadence, make tasks happen at exactly the interval even if missed
	nextCheckpoint := t.Checkpoint.Time.Add(t.Interval)

	// If the next checkpoint is past due, run now
	if nextCheckpoint.Before(time.Now()) {
		return time.Now()
	}
	return nextCheckpoint
}

func (t *Task) Exec() {
	for r, opts := range t.Runners {
		go runner.Exec(r, opts)
	}

	if t.Checkpoint != nil {
		t.Checkpoint.Time = time.Now()
	} else {
		t.Checkpoint = &TaskTime{time.Now()}
	}
}
