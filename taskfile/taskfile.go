package taskfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var ErrorNoTasks = fmt.Errorf("no tasks")

type TaskFile struct {
	Tasks []*Task
}

func NewTaskFile() *TaskFile {
	return &TaskFile{Tasks: make([]*Task, 0)}
}

func (tf *TaskFile) AddTask(task *Task) {
	tf.Tasks = append(tf.Tasks, task)
}

func (tf *TaskFile) GetNextTask() (*Task, error) {
	if len(tf.Tasks) == 0 {
		return nil, ErrorNoTasks
	}

	var soonestTask *Task
	for _, task := range tf.Tasks {
		if soonestTask == nil {
			soonestTask = task
			continue
		}

		nextTime := task.NextExecutionTime().Round(time.Second)
		soonestTime := soonestTask.NextExecutionTime().Round(time.Second)
		if nextTime.Before(soonestTime) {
			soonestTask = task
		}
	}

	return soonestTask, nil
}

func (tf *TaskFile) Save(path string) error {
	bytes, err := json.MarshalIndent(tf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0755)
}

func NewFromFile(path string) (*TaskFile, error) {
	tf := NewTaskFile()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("taskfile at %s does not exist", path)
	}

	err = json.Unmarshal(data, tf)
	if err != nil {
		return nil, fmt.Errorf("could not parse taskfile: %w", err)
	}

	return tf, nil
}
