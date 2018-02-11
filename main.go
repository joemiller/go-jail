package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/syndtr/gocapability/capability"
	"golang.org/x/sys/unix"
)

const allCapabilityTypes = capability.CAPS | capability.BOUNDS | capability.AMBS

func main() {

	// config := &configs.Config{
	// 	Capabilities: &configs.Capabilities{
	// 		Bounding: []string{
	// 			"CAP_CHOWN",
	// 			"CAP_DAC_OVERRIDE",
	// 			"CAP_FSETID",
	// 			"CAP_FOWNER",
	// 			"CAP_MKNOD",
	// 			"CAP_NET_RAW",
	// 			"CAP_SETGID",
	// 			"CAP_SETUID",
	// 			"CAP_SETFCAP",
	// 			"CAP_SETPCAP",
	// 			"CAP_NET_BIND_SERVICE",
	// 			"CAP_SYS_CHROOT",
	// 			"CAP_KILL",
	// 			"CAP_AUDIT_WRITE",
	// 		},
	// 	},
	// }

	// clear capabilities
	bound := []capability.Cap{}
	pid, err := capability.NewPid(os.Getpid())
	if err != nil {
		log.Fatal(err)
	}
	pid.Clear(allCapabilityTypes)
	pid.Set(capability.BOUNDS, bound...)
	if err := pid.Apply(allCapabilityTypes); err != nil {
		log.Fatal(err)
	}
	log.Println(pid)

	// set NoNewPrivileges
	if err := unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0); err != nil {
		log.Fatal(err)
	}

	// get full path to command to wrap, fail if it does not exist
	cmd, err := exec.LookPath(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if err := syscall.Exec(cmd, os.Args[1:], os.Environ()); err != nil {
		log.Fatal(err)
	}
}
