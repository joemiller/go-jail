package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/syndtr/gocapability/capability"
	"golang.org/x/sys/unix"
)

var (
	version = "development"

	user         = flag.String("user", "", "Switch to this user, or uid if prefixed with #")
	group        = flag.String("group", "", "Switch to this group, or gid if prefixed with #")
	caps         = flag.String("caps", "", "Comma-separated list of capabilities to allow. Default is \"\": no capabilities allowed.")
	syscallAllow = flag.String("syscall-allow", "", "TODO")
	syscallDeny  = flag.String("syscall-deny", "", "TODO")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [--] command [args]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Valid capabilities: %s\n\n", strings.Join(listCaps(), ","))
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	// log.Println("DEBUG", *user)         // TODO
	// log.Println("DEBUG", flag.Args())   // TODO
	// log.Println("DEBUG", capabilityMap) // TODO

	// remaining args specify the command (and its args) to wrap
	if flag.NArg() == 0 {
		log.Fatalf("error: missing command. Run with -h for usage.")
	}

	// get full path to command to wrap, fail if it does not exist
	cmd, err := exec.LookPath(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	// set capability bounding set
	bounding := []capability.Cap{}
	if *caps != "" {
		for _, capName := range strings.Split(*caps, ",") {
			v, ok := capabilityMap[strings.ToUpper(capName)]
			if !ok {
				log.Fatalf("Unknown capability '%s'", capName)
			}
			bounding = append(bounding, v)
		}
	}
	if err := setCapabilities(bounding); err != nil {
		log.Fatal(err)
	}

	// setup seccomp (syscall) filters
	if *syscallAllow != "" && *syscallDeny != "" {
		log.Fatal("Cannot specify both -syscall-allow and -syscall-deny. Use one or the other. -h for help.")
	}
	if *syscallAllow != "" {
		log.Println("TODO: setup syscall whitelist")

	}
	if *syscallDeny != "" {
		log.Println("TODO: setup syscall blacklist")
	}

	// set NoNewPrivileges
	if err := unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0); err != nil {
		log.Fatal(err)
	}

	if err := syscall.Exec(cmd, flag.Args(), os.Environ()); err != nil {
		log.Fatal(err)
	}
	// never reached
}
