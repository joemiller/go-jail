go-jail
=======

Simple wrapper for executing processes under a sandbox. Sandboxing is implemented
with capabilities(7) filtering, and seccomp2 (syscall) filtering.

WARNING: Consider this alpha quality software. It is an experiment at this stage.
Don't use in production.

Usage
-----

The wrapper should be invoked as root.

Required arguments:

- `-user`: User name or UID prefixed with `#` to execute command as.
- `-group`: Group name or GID prefixed with `#` to execute command as.
- `command [args]`: the command and args to wrap.

Optional arguments:

- `-caps="..."`: Comma separated list of `capabilities(7)` to include in the
  capability bounding set. Empty string to drop all capabilities (default if not specified)
- `-syscall-allow="..."`: Comma separated listed of system calls to allow. This is effectively
  a whitelist mode. All system calls not listed here will be denied with `EPERM` error code.
- `-syscall-block="..."`: Comma separated listed of system calls to block. This is effectively
  a blacklist mode. All system calls listed here will be denied with `EPERM` error code, and the
  rest will be allowed.

Examples:

1. Run a process as user `daemon`, group `daemon`:

```
$ go-jail -user daemon -group daemon -- whoami
```

Note that by default an empty capability set is applied:

```
$ go-jail -user daemon -group daemon -- grep CapBnd /proc/self/status
CapBnd: 0000000000000000
```

2. Run with limited capability bounding set:

```
$ go-jail -user "daemon" -group "daemon" --caps="CAP_CHOWN,CAP_KILL,CAP_DAC_OVERRIDE" -- grep CapBnd /proc/self/status
CapBnd: 0000000000000023
```

3. Run with a restricted list of allowed syscalls:

```
$ go-jail -user daemon -group daemon -syscall-allow="execve,exit,read,open,close,mmap,mmap2,fstat,fstat64,access,mprotect,set_thread_area,brk,openat,exit_group,prctl,arch_prctl,write" \
    -- echo "just enough system calls for echo to succeed"
just enough system calls for echo to succeed
```

When running with `-syscall-allow` you will need a minimum set of syscalls for most
programs to execute correctly at all. The list above list is a start.

4. Run with a list of system calls to deny:

```
$ go-jail -user daemon -group daemon -syscall-deny="write" \
-- echo "just enough system calls for echo to succeed"
$ echo $?
1
```


TODO:
- [x] decide on UI:
  - [x] minimal -- take a list of capabilities to add/drop, and a list of syscalls to blacklist or whitelist
  - [ ] oci/docker compatibility -- take a config.json file and use the capabilities and seccomp
      settings. This would allow more granular seccomp policies such as filtering on the args
      to syscalls
