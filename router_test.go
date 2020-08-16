package dgc

import (
	"fmt"
	"testing"
)

func Test_buildCheckPrefixes(t *testing.T) {
	cmd := Command{
		Name: "help",
		Aliases: []string {
			"h",
			"hep",
		},
	}

	toCheck := buildCheckPrefixes(&cmd)
	if len(toCheck) != 3 {
		t.Fatal(fmt.Sprintf("expected 3 aliases, got %d", len(toCheck)))
	}
}
