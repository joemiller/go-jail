go-jail
=======

Simple wrapper for executing processes under a sandbox. Sandboxing is implemented
with capabilities(7) filtering, and seccomp2 (syscall) filtering.

WARNING: Consider this alpha quality software. It is an experiment at this stage.
Don't use in production, I haven't (but if you do, please send feedback via
github issues).

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

Development & Testing
---------------------

This project is Linux specific and must be built and tested within Linux. You can
still do "local development" on macOS with Docker installed. Use the `make devshell`
command to create an interactive container suitable for build and test tasks.

```
$ make devshell
...
root@f9962fe0a031:/go/src/github.com/joemiller/go-jail# make deps
root@f9962fe0a031:/go/src/github.com/joemiller/go-jail# make test
root@f9962fe0a031:/go/src/github.com/joemiller/go-jail# make build
root@f9962fe0a031:/go/src/github.com/joemiller/go-jail# exit
```

### Dependencies

You will need libseccomp-dev (debian) or libseccomp-devel (redhat) packages installed
to build a binary. The `build` and `test` make targets will build and run within a docker
container that has these dependencies. This makes it easy to develop and build/test
on platforms other than Linux (ie: macOS).

To run tests, install `bats` (available in most distros).

Run `make deps` to install Go dependencies.

The `./Dockerfile.build` is also a good reference for package dependencies.

### Tests

Run `make test`. You must run `make build` first to create the go-jail binary. Tests
are performed using `bats` utilizing the binary.

### Build

Run `make build`


TODO:
----

- [x] decide on UI:
  - [x] minimal -- take a list of capabilities to add/drop, and a list of syscalls to blacklist or whitelist
  - [ ] oci/docker compatibility -- take a config.json file and use the capabilities and seccomp
      settings. This would allow more granular seccomp policies such as filtering on the args
      to syscalls
- [x] CI/CD pipeline. build linux amd64 binary and push to github-releases
- [ ] add license file
