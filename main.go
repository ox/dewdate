package main

import (
	"dewdate/taskfile"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var taskfilePath string

func init() {
	flag.StringVar(&taskfilePath, "tasks", "", "Path to taskfile")
	flag.Parse()

	if taskfilePath == "" {
		log.Fatalf("--tasks parameter is required")
	}
}

func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func tick(tf *taskfile.TaskFile) {
	nextTask, err := tf.GetNextTask()
	if err != nil {
		log.Fatalf("Could not get next task: %s", err)
	}
	log.Printf("Next task %s to run at %s", nextTask.Name, nextTask.NextExecutionTime().Format(taskfile.TimeFormat))

	// Sleep until it's time for the next task
	durationToNextTask := time.Until(nextTask.NextExecutionTime())
	log.Printf("Sleeping %s until next task (%s)\n", durationToNextTask.Round(time.Second).String(), nextTask.Name)
	time.Sleep(durationToNextTask)

	nextTask.Exec()

	log.Printf("Saving task file at %s\n", taskfilePath)
	if err = tf.Save(taskfilePath); err != nil {
		log.Fatalf("Could not save taskfile at '%s': %s", taskfilePath, err)
	}
}

func main() {
	setupCloseHandler()

	log.Printf("Reading taskfile at '%s'\n", taskfilePath)
	tf, err := taskfile.NewFromFile(taskfilePath)
	if err != nil {
		log.Fatalf("Error reading taskfile: %s", err)
	}

	log.Printf("Loaded %d tasks\n", len(tf.Tasks))
	for i, task := range tf.Tasks {
		log.Printf("%d: %s", i, task)
	}

	for {
		tick(tf)
	}
}
