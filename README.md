go-jail
=======

Simple wrapper for executing processes under a sandbox. Sandboxing is implemented
with capabilities(7) filtering, and seccomp2 (syscall) filtering.


TODO:
- [ ] decide on UI:
  - [ ] minimal -- take a list of capabilities to add/drop, and a list of syscalls to blacklist or whitelist
  - [ ] oci/docker compatibility -- take a config.json file and use the capabilities and seccomp
      settings. This would allow more granular seccomp policies such as filtering on the args
      to syscalls


UI sketches
-----------

- whitelist/blacklist syscall, cap-set:

```
    go-jail -u user -g group --caps="cap_sys_admin,cap_net_admin" --sc-allow="foo,bar,baz" -- path//somecmd --and args
    go-jail -u user -g group --caps="cap_sys_admin,cap_net_admin" --sc-deny="foo,bar,baz" -- path/somecmd --and args

    # --sc-allow and --sc-deny are mutually exclusive. If --sc-allow is provided, the default
    # action for syscalls is "deny" (whitelist). If --sc-deny is provided the default
    # action is "allow" (blacklist mode)
```

Maybe later:

- using capabilities,user/group,seccomp settings from an OCI config.json spec:

```
    go-jail -s config.json -- /path/to/somecmd --and-its args

```
