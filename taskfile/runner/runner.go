package runner

import "fmt"

var runners map[string]Runner
var ErrorNoRunner = fmt.Errorf("no runners")

type Runner interface {
	Run(options map[string]interface{}) error
}

func init() {
	runners = make(map[string]Runner)
	runners["cmd"] = &Cmd{}
}

func Exec(key string, options map[string]interface{}) error {
	if runner, ok := runners[key]; ok {
		return runner.Run(options)
	}
	return ErrorNoRunner
}
