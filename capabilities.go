package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/syndtr/gocapability/capability"
)

const allCapabilityTypes = capability.CAPS | capability.BOUNDS | capability.AMBS

var capabilityMap map[string]capability.Cap

func init() {
	capabilityMap = make(map[string]capability.Cap)
	for _, cap := range capability.List() {
		if cap > capability.CAP_LAST_CAP {
			continue
		}
		key := fmt.Sprintf("CAP_%s", strings.ToUpper(cap.String()))
		capabilityMap[key] = cap
	}
}

// listCaps() returns a list of strings representing valid capabilities.
// eg: "CAP_SYS_ADMIN", "CAP_SYS_PTRACE", etc
func listCaps() []string {
	keys := []string{}
	for key := range capabilityMap {
		keys = append(keys, key)
	}
	return keys
}

// clearCapabilities clears all capability sets
func clearCapabilities() error {
	pid, err := capability.NewPid(os.Getpid())
	if err != nil {
		return err
	}
	pid.Clear(allCapabilityTypes)
	return pid.Apply(allCapabilityTypes)
}

// setCapabilities() resets all capabilities to empty and then applies the
// capability BOUNDING set specified in bounds. Pass an empty slice to set
// and empty bounding set.
func setCapabilities(capType capability.CapType, caps []capability.Cap) error {
	pid, err := capability.NewPid(os.Getpid())
	if err != nil {
		return err
	}
	pid.Clear(capType)
	pid.Set(capType, caps...)
	return pid.Apply(capType)
}
