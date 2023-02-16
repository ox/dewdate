package taskfile

import (
	"fmt"
	"time"
)

const TimeFormat = "02 Jan 06 15:04:05 MST"

type TaskTime struct {
	time.Time
}

func (t TaskTime) String() string {
	return t.Time.Format(TimeFormat)
}

func (t TaskTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(TimeFormat))), nil
}

func (t *TaskTime) UnmarshalJSON(data []byte) error {
	nt, err := time.Parse(fmt.Sprintf("\"%s\"", TimeFormat), string(data))
	t.Time = nt
	return err
}
