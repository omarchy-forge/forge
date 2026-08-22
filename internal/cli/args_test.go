package cli

import (
	"reflect"
	"testing"
)

func TestReorderInitArgsSupportsOptionsAfterDirectory(t *testing.T) {
	got, err := reorderInitArgs([]string{"project", "--id", "dev.example.project", "--dry-run", "--section=left"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"--id", "dev.example.project", "--dry-run", "--section=left", "project"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("reorderInitArgs() = %v, want %v", got, want)
	}
}
