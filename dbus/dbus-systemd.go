package main

import (
	"fmt"
	"os"
	"strings"

	"log"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
		os.Exit(1)
	}

	// list all unit
	var (
		s [][]interface{}
	)
	err = conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1").Call("org.freedesktop.systemd1.Manager.ListUnits", 0).Store(&s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get list of owned names:", err)
		os.Exit(1)
	}

	fmt.Println("Currently owned names on the system bus:")
	for _, v := range s {
		fmt.Println(v)
	}

	// get object path
	var path dbus.ObjectPath
	err = conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1").Call("org.freedesktop.systemd1.Manager.GetUnit", 0, "nginx.service").Store(&path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get object path:", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", path)

	// get load status
	var loadStatus string
	err = conn.Object("org.freedesktop.systemd1", path).Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "LoadState").Store(&loadStatus)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get load status:", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", loadStatus)

	// get active status
	var activeStatus string
	err = conn.Object("org.freedesktop.systemd1", path).Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "ActiveState").Store(&activeStatus)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get load status:", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", activeStatus)

	// get main pid
	var pid uint32
	err = conn.Object("org.freedesktop.systemd1", path).Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Service", "MainPID").Store(&pid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get main pid:", err)
		os.Exit(1)
	}
	fmt.Printf("%d\n", pid)

	var f func(*introspect.Node) string
	f = func(node *introspect.Node) string {
		for _, i := range node.Interfaces {
			if strings.Compare(i.Name, "org.fedoraproject.FirewallD1") != 0 {
				continue
			}
			for _, m := range i.Methods {
				if strings.Compare(m.Name, "getZoneSettings") == 0 {
					for _, a := range m.Args {
						if a.Direction == "out" {
							return a.Type
						}
					}
				}
			}
		}
		return ""
	}

	intro, err := introspect.Call(conn.Object("org.freedesktop.systemd1", path))
	if err != nil {
		log.Fatal("can't introspect because ", err)
	}
	signature := f(intro)
	if signature == "" {
		log.Fatal("did not find signature. sad.")
	}

	log.Println("found signature ", signature)

	type FirewallZone struct {
		Version     string
		Short       string
		Description string
		Unused      bool
		Target      string
		Services    []string
		Ports       []struct {
			Port     string
			Protocol string
		}
		ICMP         []string
		Masquerade   bool
		ForwardPorts []struct {
			Port     string
			Protocol string
			ToPort   string
			ToAddr   string
		}
		Interfaces  []interface{}
		Sources     []string
		Rich        []string
		Protocols   []string
		SourcePorts []struct {
			Port     string
			Protocol string
		}
		BlockInversion bool //
	}

	// FirewallZoneBis will appear (at a guess) firewalld > 0.7.4
	// Adds Forward.
	type FirewallZoneBis struct {
		FirewallZone
		Forward bool //  b
	}

	var store interface{}

	switch signature {
	case "(sssbsasa(ss)asba(ssss)asasasasa(ss)b)":
		store = &FirewallZone{}
	case "(sssbsasa(ss)asba(ssss)asasasasa(ss)bb)":
		store = &FirewallZoneBis{}
	}

	obj := conn.Object("org.freedesktop.systemd1", path).Call("getZoneSettings", 0, "public")
	if err := obj.Err; err != nil {
		log.Fatal("call getZoneSettings failed ", err)
	}

	if err := obj.Store(store); err != nil {
		log.Println("store failed", err)
	}
}
