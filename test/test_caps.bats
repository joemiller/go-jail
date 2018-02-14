# vim: set ft=sh:
# -*- mode: sh -*-

@test "empty capabilitiy bounding set" {
  run ./go-jail --caps="" -- grep CapBnd /proc/self/status
	echo "output = ${output}"
	[[ "$output" =~ ^CapBnd:[[:space:]]*0000000000000000$ ]]
}

@test "allow single CAP" {
  run ./go-jail --caps="CAP_CHOWN" -- grep CapBnd /proc/self/status
	echo "output = ${output}"
	[[ "$output" =~ ^CapBnd:[[:space:]]*0000000000000001$ ]]
}

@test "allow multiple CAPs" {
  run ./go-jail --caps="CAP_CHOWN,CAP_KILL,CAP_DAC_OVERRIDE" -- grep CapBnd /proc/self/status
	echo "output = ${output}"
	[[ "$output" =~ ^CapBnd:[[:space:]]*0000000000000023$ ]]
}

@test "fail on unknown capability" {
  run ./go-jail --caps="CAP_FOO" -- grep CapBnd /proc/self/status
	echo "output = ${output}"
	[[ $status -eq 1 ]]
	[[ "$output" == *"Unknown capability 'CAP_FOO'" ]]
}
