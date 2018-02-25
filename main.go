package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	osuser "os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/opencontainers/runc/libcontainer/system"
	"github.com/syndtr/gocapability/capability"
	"golang.org/x/sys/unix"
)

var (
	version = "development"

	uid int
	gid int

	user         = flag.String("user", "", "Switch to this user, or uid if prefixed with #")
	group        = flag.String("group", "", "Switch to this group, or gid if prefixed with #")
	caps         = flag.String("caps", "", "Comma-separated list of capabilities to allow. Default is \"\": no capabilities allowed.")
	syscallAllow = flag.String("syscall-allow", "", "Whitelist mode: Block all syscalls except those listed in a comma-separated list.")
	syscallDeny  = flag.String("syscall-deny", "", "Blacklist mode: Allow all syscalls except those listed in a comma-separated list.")
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

	// after parsing args, the remainder is our command (and its args) to wrap.
	if flag.NArg() == 0 {
		log.Fatalf("error: missing command. Run with -h for usage.")
	}

	// get full path to command to wrap, fail if it does not exist
	cmd, err := exec.LookPath(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	// require -user and -group args
	if *user == "" {
		log.Fatal("error: missing -user. Run with -h for usage.")
	}
	if *group == "" {
		log.Fatal("error: missing -group. Run with -h for usage.")
	}

	// parse user flag. Either an exact direct UID (if prefixed with #) or lookup an existing user by name
	if strings.HasPrefix(*user, "#") {
		parsed, err := strconv.Atoi(strings.Trim(*user, "#"))
		if err != nil {
			log.Fatalf("Failed to parse user ID: %s", err)
		}
		uid = parsed
	} else {
		u, err := osuser.Lookup(*user)
		if err != nil {
			log.Fatal(err)
		}
		uid, _ = strconv.Atoi(u.Uid)
	}

	// parse group flag. Either an exact direct GID (if prefixed with #) or lookup an existing group by name
	if strings.HasPrefix(*group, "#") {
		parsed, err := strconv.Atoi(strings.Trim(*group, "#"))
		if err != nil {
			log.Fatalf("Failed to parse group ID: %s", err)
		}
		gid = parsed
	} else {
		g, err := osuser.LookupGroup(*group)
		if err != nil {
			log.Fatal(err)
		}
		gid, _ = strconv.Atoi(g.Gid)
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
	// set capability bounding set. Do this before setuid
	if err := setCapabilities(capability.BOUNDS, bounding); err != nil {
		log.Fatal(err)
	}

	// setuid/setgid (must setgid first)
	if err := system.Setgid(gid); err != nil {
		log.Fatalf("failed to setgid: %s", err)
	}
	if err := system.Setuid(uid); err != nil {
		log.Fatalf("failed to setuid: %s", err)
	}

	// after setuid we can drop all effective and inherited caps
	if err := clearCapabilities(); err != nil {
		log.Fatal(err)
	}

	// setup seccomp (syscall) filters
	{
		if *syscallAllow != "" && *syscallDeny != "" {
			log.Fatal("Cannot specify both -syscall-allow and -syscall-deny. Use one or the other. -h for help.")
		}
		if *syscallAllow != "" {
			if err := initSeccompWhitelist(strings.Split(*syscallAllow, ",")); err != nil {
				log.Fatal(err)
			}
		}
		if *syscallDeny != "" {
			if err := initSeccompBlacklist(strings.Split(*syscallDeny, ",")); err != nil {
				log.Fatal(err)
			}
		}
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
