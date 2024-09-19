package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
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
}
