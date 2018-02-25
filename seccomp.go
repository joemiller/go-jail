package main

import (
	"fmt"

	libseccomp "github.com/seccomp/libseccomp-golang"
	"golang.org/x/sys/unix"
)

var (
	seccompAllow = libseccomp.ActAllow
	seccompDeny  = libseccomp.ActErrno.SetReturnCode(int16(unix.EPERM))
)

func initSeccompWhitelist(syscalls []string) error {
	filter, err := libseccomp.NewFilter(seccompDeny)
	if err != nil {
		return fmt.Errorf("failed to initialize seccomp filter: %s", err)
	}
	for _, name := range syscalls {
		scmpSyscall, err := libseccomp.GetSyscallFromName(name)
		if err != nil {
			return fmt.Errorf("failed to resolve syscall '%s': %s", name, err)
		}
		if err := filter.AddRule(scmpSyscall, seccompAllow); err != nil {
			return fmt.Errorf("failed to add seccomp rule: %s", err)
		}
	}
	return filter.Load()
}

func initSeccompBlacklist(syscalls []string) error {
	filter, err := libseccomp.NewFilter(seccompAllow)
	if err != nil {
		return fmt.Errorf("failed to initialize seccomp filter: %s", err)
	}
	for _, name := range syscalls {
		scmpSyscall, err := libseccomp.GetSyscallFromName(name)
		if err != nil {
			return fmt.Errorf("failed to resolve syscall '%s': %s", name, err)
		}
		if err := filter.AddRule(scmpSyscall, seccompDeny); err != nil {
			return fmt.Errorf("failed to add seccomp rule: %s", err)
		}
	}
	return filter.Load()
}
