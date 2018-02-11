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
