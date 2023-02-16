package taskfile

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	startDateStr := "18 Jan 23 17:13:05 UTC"
	startDate, _ := time.Parse(TimeFormat, startDateStr)
	checkpointStr := "18 Jan 23 17:13:05 UTC"
	checkpoint, _ := time.Parse(TimeFormat, checkpointStr)
	tf := NewTaskFile()
	task := &Task{Name: "test", Start: TaskTime{startDate}, Interval: time.Hour * 2, Checkpoint: &TaskTime{checkpoint}}
	tf.AddTask(task)

	b, err := json.Marshal(tf)
	if err != nil {
		t.Error(err)
	}

	bexp := []byte(`{"Tasks":[{"Name":"test","Start":"18 Jan 23 17:13:05 UTC","Interval":"2h0m0s","Checkpoint":"18 Jan 23 17:13:05 UTC"}]}`)
	if !bytes.Equal(b, bexp) {
		t.Errorf("unexpected marshalled json: %v", string(b))
	}

	newTf := NewTaskFile()
	if err = json.Unmarshal(b, newTf); err != nil {
		t.Fatalf("could not unmarshal: %s", err)
	}

	if !reflect.DeepEqual(tf, newTf) {
		t.Errorf("marshaling and unmarshaling does not create idential taskfiles. before: %v, after: %v", tf, newTf)
	}
}

func TestMarshalAfterExec(t *testing.T) {
	startDateStr := "18 Jan 23 17:13:05 UTC"
	startDate, _ := time.Parse(TimeFormat, startDateStr)

	tf := NewTaskFile()
	task := &Task{Name: "test", Start: TaskTime{startDate}, Interval: time.Hour * 2}
	tf.AddTask(task)

	b, err := json.Marshal(tf)
	if err != nil {
		t.Error(err)
	}

	bexp := []byte(`{"Tasks":[{"Name":"test","Start":"18 Jan 23 17:13:05 UTC","Interval":"2h0m0s"}]}`)
	if !bytes.Equal(b, bexp) {
		t.Errorf("unexpected marshalled json: %v", string(b))
	}

	// Simulate Exec()
	checkpointStr := "18 Jan 23 17:13:05 UTC"
	checkpoint, _ := time.Parse(TimeFormat, checkpointStr)
	task.Checkpoint = &TaskTime{checkpoint}

	b2, err := json.Marshal(tf)
	if err != nil {
		t.Error(err)
	}

	bexp2 := []byte(`{"Tasks":[{"Name":"test","Start":"18 Jan 23 17:13:05 UTC","Interval":"2h0m0s","Checkpoint":"18 Jan 23 17:13:05 UTC"}]}`)
	if !bytes.Equal(b2, bexp2) {
		t.Errorf("unexpected marshalled json: %v", string(b))
	}
}

func TestGetNextTask(t *testing.T) {
	start := time.Now()
	tf := NewTaskFile()
	_, err := tf.GetNextTask()
	if err != ErrorNoTasks {
		t.Errorf("expected ErrorNoTasks, got: %s", err)
		t.FailNow()
	}

	// task1 should run twice before task2
	task1 := &Task{Name: "task1", Start: TaskTime{start}, Interval: 5 * time.Second}
	task2 := &Task{Name: "task2", Start: TaskTime{start.Add(12 * time.Second)}, Interval: 12 * time.Second}
	tf.AddTask(task1)
	tf.AddTask(task2)

	// task1 should be next
	nextTask, err := tf.GetNextTask()
	if err != nil {
		t.Fatalf("taskfile should have tasks: %s", err)
	}
	if nextTask != task1 {
		t.Errorf("next task should have been %s, got %s", task1.Name, nextTask.Name)
		t.FailNow()
	}
	task1.Checkpoint = &TaskTime{task1.Start.Add(task1.Interval)}

	// task1 should still have another execution before Task2 even begins
	nextTask, _ = tf.GetNextTask()
	if nextTask != task1 {
		t.Errorf("next task should have been %s, got %s", task1.Name, nextTask.Name)
		t.FailNow()
	}
	task1.Checkpoint.Time = task1.Checkpoint.Time.Add(task1.Interval)

	// Now task2 should now be next
	nextTask, _ = tf.GetNextTask()
	if nextTask != task2 {
		t.Errorf("next task should have been %s, got %s", task2.Name, nextTask.Name)
		t.FailNow()
	}
}
