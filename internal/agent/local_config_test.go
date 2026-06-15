package agent

import "testing"

func TestSaveLoadLocalConfig(t *testing.T) {
	path := t.TempDir() + "/agent.json"
	in := LocalConfig{ControllerURL: "http://controller", Token: "token", NodeID: "node-1", WGConfigPath: "wg.conf", SyncCommand: "true"}
	if err := SaveLocalConfig(path, in); err != nil {
		t.Fatal(err)
	}
	out, err := LoadLocalConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if out.NodeID != in.NodeID || out.ControllerURL != in.ControllerURL || out.WGConfigPath != in.WGConfigPath {
		t.Fatalf("unexpected config: %#v", out)
	}
}
