# Dewdate

Dewdate is a scheduled task runner like cron, except that it uses durations and start times to specify when tasks should be done. This is meant more for tasks that happen periodically.

A taskfile looks something like:

```go
{
  "Tasks": [
    {
      "Name": "greet the world",
      "Start": "18 Jan 23 17:13:05 UTC",
      "Interval": "5s",
      "Checkpoint": "03 Feb 23 12:15:51 EST",
      "Runners": {
        "cmd": {
          "Exec": "echo hello world"
        }
      }
    },
    {
      "Name": "test",
      "Start": "18 Jan 23 17:13:05 UTC",
      "Interval": "12s"
    }
  ]
}
```

A Task has a Name, a Start time, a period, and optionally a command that should be run. Checkpoints are used by Dewdate to know when the last time a task was run.

## Running Dewdate

Download the source, create a taskfile, and run `go run main.go --tasks <path to taskfile>`. You should see log messages start flowing.
