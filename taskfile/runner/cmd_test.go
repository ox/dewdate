package runner

import "testing"

func TestCmdRun(t *testing.T) {
	cmd := &Cmd{}
	t.Run("bad options", func(t *testing.T) {
		badOptions := make(map[string]interface{})
		if err := cmd.Run(badOptions); err == nil {
			t.Errorf("expected Cmd.Run to fail")
			t.FailNow()
		}
	})

	t.Run("good options", func(t *testing.T) {
		goodOptions := make(map[string]interface{})
		goodOptions["Exec"] = "echo hello world"
		goodOptions["DryRun"] = true
		if err := cmd.Run(goodOptions); err != nil {
			t.Errorf("did not expect Cmd.Run to fail with: %s", err)
			t.FailNow()
		}
	})
}
