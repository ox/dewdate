package runner

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Cmd struct{}

func (c *Cmd) Run(options map[string]interface{}) error {
	execStr, ok := options["Exec"]
	if !ok {
		return fmt.Errorf("execStr expects Exec option")
	}

	execParts := strings.Split(execStr.(string), " ")

	log.Printf("[cmd]: $ %s\n", execStr)
	if dry, ok := options["DryRun"]; ok && dry.(bool) {
		return nil
	}

	cmd := exec.Command(execParts[0], execParts[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[cmd]: error reading output: %s", err)
		return err
	}

	fmt.Printf("%s\n", string(output))

	return nil
}
