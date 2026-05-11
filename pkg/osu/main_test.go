package osu

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func runTest(num int) func(t *testing.T) {
	filename := fmt.Sprintf("testfiles/v%d.osu", num)
	outfile := fmt.Sprintf("testfiles/v%d_out.json", num)
	return func(t *testing.T) {
		if _, err := os.Stat(filename); err != nil {
			return
		}
		b, err := ParseFile(filename)
		if err != nil {
			t.Error(err)
		}
		if bytes, err := json.MarshalIndent(b, "", "\t"); err != nil {
			t.Error(err)
		} else if os.Getenv("UPDATE_GOLDEN") != "" {
			if err := os.WriteFile(outfile, bytes, 0644); err != nil {
				t.Error(err)
			}
		} else if ans, err := os.ReadFile(outfile); err != nil {
			t.Error(err)
		} else if string(bytes) != string(ans) {
			t.Errorf("Test v%d: Output are not the same", num)
		}
	}
}

func TestAll(t *testing.T) {
	for i := 3; i <= 14; i++ {
		t.Run(fmt.Sprintf("Test v%d", i), runTest(i))
	}
}
