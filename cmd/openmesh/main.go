package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/xybingbing/openmesh/internal/agent"
	"github.com/xybingbing/openmesh/internal/controller"
	"github.com/xybingbing/openmesh/internal/faketcp"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	ctx := context.Background()
	var err error

	switch os.Args[1] {
	case "controller":
		err = runController(ctx, os.Args[2:])
	case "agent":
		err = runAgent(ctx, os.Args[2:])
	case "faketcp":
		err = runFakeTCP(ctx, os.Args[2:])
	case "version":
		fmt.Println("openmesh dev")
	default:
		usage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`openmesh

Usage:
  openmesh controller --listen :8080 --data ./openmesh.json --token dev-token
  openmesh agent register --controller http://127.0.0.1:8080 --token dev-token --name node-a --public-key <key>
  openmesh agent config --controller http://127.0.0.1:8080 --token dev-token --node-id <id>
  openmesh agent save-config --controller http://127.0.0.1:8080 --token dev-token --node-id <id> --config /etc/openmesh/agent.json
  openmesh agent daemon --config /etc/openmesh/agent.json
  openmesh agent up --config /etc/openmesh/agent.json
  openmesh agent down --config /etc/openmesh/agent.json

  openmesh faketcp server --listen :9000
  openmesh faketcp client --remote 127.0.0.1:9000

  openmesh version`)
}

func runController(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("controller", flag.ContinueOnError)
	listen := fs.String("listen", ":8080", "HTTP listen address")
	data := fs.String("data", "openmesh.json", "data file path")
	token := fs.String("token", "", "API token")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return controller.Run(ctx, controller.Config{Listen: *listen, DataPath: *data, Token: *token})
}

func runFakeTCP(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing faketcp subcommand")
	}

	switch args[0] {
	case "server":
		fs := flag.NewFlagSet("faketcp server", flag.ContinueOnError)
		listen := fs.String("listen", ":9000", "udp listen address")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}

		t, err := faketcp.ListenUDP(*listen)
		if err != nil {
			return err
		}
		defer t.Close()

		fmt.Println("faketcp server listening on", *listen)
		for {
			pkt, addr, err := t.Recv(ctx, 65535)
			if err != nil {
				return err
			}
			fmt.Printf("recv from %v: %s\n", addr, string(pkt.Payload))
		}

	case "client":
		fs := flag.NewFlagSet("faketcp client", flag.ContinueOnError)
		remote := fs.String("remote", "127.0.0.1:9000", "udp remote address")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}

		t, err := faketcp.ListenUDP(":0")
		if err != nil {
			return err
		}
		defer t.Close()

		pkt := faketcp.Packet{
			SrcIP:   faketcp.ParseIPv4("10.0.0.2"),
			DstIP:   faketcp.ParseIPv4("10.0.0.1"),
			SrcPort: 40000,
			DstPort: 443,
			Seq:     1,
			Flags:   faketcp.FlagACK | faketcp.FlagPSH,
			Payload: []byte("hello-openmesh"),
		}

		fmt.Println("faketcp client sending to", *remote)
		return t.SendTo(ctx, *remote, pkt)

	default:
		return fmt.Errorf("unknown faketcp subcommand %q", args[0])
	}
}

func runAgent(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing agent subcommand")
	}
	switch args[0] {
	case "register":
		fs := flag.NewFlagSet("agent register", flag.ContinueOnError)
		controllerURL := fs.String("controller", "", "controller URL")
		token := fs.String("token", "", "API token")
		name := fs.String("name", "", "node name")
		publicKey := fs.String("public-key", "", "WireGuard public key")
		endpoint := fs.String("endpoint", "", "optional public endpoint")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		return agent.Register(ctx, agent.RegisterConfig{ControllerURL: *controllerURL, Token: *token, Name: *name, PublicKey: *publicKey, Endpoint: *endpoint}, os.Stdout)
	case "config":
		fs := flag.NewFlagSet("agent config", flag.ContinueOnError)
		controllerURL := fs.String("controller", "", "controller URL")
		token := fs.String("token", "", "API token")
		nodeID := fs.String("node-id", "", "node id")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		return agent.Config(ctx, agent.ConfigConfig{ControllerURL: *controllerURL, Token: *token, NodeID: *nodeID}, os.Stdout)
	case "save-config":
		fs := flag.NewFlagSet("agent save-config", flag.ContinueOnError)
		controllerURL := fs.String("controller", "", "controller URL")
		token := fs.String("token", "", "API token")
		nodeID := fs.String("node-id", "", "node id")
		path := fs.String("config", "/etc/openmesh/agent.json", "agent config path")
		wgPath := fs.String("wg-config", "/etc/wireguard/openmesh.conf", "WireGuard config path")
		wgInterface := fs.String("wg-interface", "openmesh", "WireGuard interface name")
		wgAddress := fs.String("wg-address", "", "WireGuard interface address, optional")
		wgMTU := fs.Int("wg-mtu", 1280, "WireGuard interface MTU")
		syncCommand := fs.String("sync-command", "", "optional sync command")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		return agent.SaveLocalConfig(*path, agent.LocalConfig{ControllerURL: *controllerURL, Token: *token, NodeID: *nodeID, WGConfigPath: *wgPath, WGInterface: *wgInterface, WGAddress: *wgAddress, WGMTU: *wgMTU, SyncCommand: *syncCommand})
	case "daemon":
		fs := flag.NewFlagSet("agent daemon", flag.ContinueOnError)
		path := fs.String("config", "/etc/openmesh/agent.json", "agent config path")
		interval := fs.Duration("interval", 30*time.Second, "sync interval")
		once := fs.Bool("once", false, "run once and exit")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		return agent.Daemon(ctx, agent.DaemonConfig{ConfigPath: *path, Interval: *interval, Once: *once})
	case "up":
		fs := flag.NewFlagSet("agent up", flag.ContinueOnError)
		path := fs.String("config", "/etc/openmesh/agent.json", "agent config path")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		cfg, err := agent.LoadLocalConfig(*path)
		if err != nil {
			return err
		}
		return agent.Up(ctx, cfg, nil)
	case "down":
		fs := flag.NewFlagSet("agent down", flag.ContinueOnError)
		path := fs.String("config", "/etc/openmesh/agent.json", "agent config path")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		cfg, err := agent.LoadLocalConfig(*path)
		if err != nil {
			return err
		}
		return agent.Down(ctx, cfg, nil)
	default:
		return fmt.Errorf("unknown agent subcommand %q", args[0])
	}
}
