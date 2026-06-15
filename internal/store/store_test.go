package store

import "testing"

import "github.com/xybingbing/openmesh/internal/model"

func TestRegisterAllocatesMeshIP(t *testing.T) {
	dir := t.TempDir()
	st, err := Open(dir + "/state.json")
	if err != nil {
		t.Fatal(err)
	}
	n, err := st.Register(model.RegisterRequest{Name: "node-a", PublicKey: "pub-a"})
	if err != nil {
		t.Fatal(err)
	}
	if n.ID != "node-1" {
		t.Fatalf("unexpected node id: %s", n.ID)
	}
	if n.MeshIP != "100.64.0.2/32" {
		t.Fatalf("unexpected mesh ip: %s", n.MeshIP)
	}
}
