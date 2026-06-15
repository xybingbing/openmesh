package wg

import (
	"context"
	"reflect"
	"testing"

	"github.com/xybingbing/openmesh/internal/system"
)

func TestRuntimeUpDryRun(t *testing.T) {
	runner := &system.DryRunner{}
	rt := Runtime{Runner: runner}
	err := rt.Up(context.Background(), RuntimeConfig{Interface: "om0", Config: "/tmp/om.conf", Address: "100.64.0.2/32", MTU: 1280})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"ip link add om0 type wireguard",
		"ip link set dev om0 mtu 1280",
		"wg setconf om0 /tmp/om.conf",
		"ip addr add 100.64.0.2/32 dev om0",
		"ip link set up dev om0",
	}
	if !reflect.DeepEqual(runner.Commands, want) {
		t.Fatalf("commands mismatch:\nwant=%v\n got=%v", want, runner.Commands)
	}
}

func TestRuntimeDownDryRun(t *testing.T) {
	runner := &system.DryRunner{}
	rt := Runtime{Runner: runner}
	if err := rt.Down(context.Background(), "om0"); err != nil {
		t.Fatal(err)
	}
	want := []string{"ip link del om0"}
	if !reflect.DeepEqual(runner.Commands, want) {
		t.Fatalf("commands mismatch: %v", runner.Commands)
	}
}
