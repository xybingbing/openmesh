package agent

import (
	"context"
	"reflect"
	"testing"

	"github.com/xybingbing/openmesh/internal/system"
)

func TestExtractAddress(t *testing.T) {
	got := extractAddress("[Interface]\nAddress = 100.64.0.2/32\nListenPort = 51820\n")
	if got != "100.64.0.2/32" {
		t.Fatalf("unexpected address: %q", got)
	}
}

func TestDownUsesConfiguredInterface(t *testing.T) {
	runner := &system.DryRunner{}
	cfg := LocalConfig{WGInterface: "omtest"}
	if err := Down(context.Background(), cfg, runner); err != nil {
		t.Fatal(err)
	}
	want := []string{"ip link del omtest"}
	if !reflect.DeepEqual(runner.Commands, want) {
		t.Fatalf("commands mismatch: %v", runner.Commands)
	}
}
