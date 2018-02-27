# vim: set ft=sh:
# -*- mode: sh -*-

BASE_WHITELIST="execve,exit,read,open,close,mmap,mmap2,fstat,fstat64,access,mprotect,set_thread_area,brk,openat,exit_group,prctl,arch_prctl"

# echo needs the basic set plus "write" syscall to succeed
@test "wrap echo with the correct seccomp whitelist" {
  run ./go-jail -user "daemon" -group "daemon" -syscall-allow="$BASE_WHITELIST,write" -- echo "hello world"
	echo "output = ${output}"
	[[ $status -eq 0 ]]
	[[ "$output" = *"hello world"* ]]
}

# remove 'write' and it should fail
@test "wrap echo with incorrect seccomp whitelist for it to succeed" {
  run ./go-jail -user "daemon" -group "daemon" -syscall-allow="$BASE_WHITELIST" -- echo "hello world"
  echo "output = ${output}"
  [[ $status -ne 0 ]]
}

@test "denying the write syscall should cause echo to fail" {
  run ./go-jail -user "daemon" -group "daemon" -syscall-deny="write" -- echo "hello world"
  echo "output = ${output}"
  [[ $status -ne 0 ]]
}

