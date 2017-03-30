# packer-pkg

A universal package management provisioner.

Supported operating systems:

* Fedora
* Ubuntu
* Debian

We aim to support all stable versions of a operating system.

## How it works

Pkg attems to automatically install a package based on the information given. For example:

```javascript
{
  "type": "pkg",
  "name": "docker",
  "file": "/usr/bin/docker",
}
```

This will

1. Determine the operating system and package manager
2. Update the package manager cache if needed
3. Attempt to use the file path to install the package (not supported by all package managers)
4. Fall back to the name if file installation was not successful

For cases where all other methods fail, you can set the package name manually.

```javascript
{
  // applies to all debian versions
  "debian": "docker.io",
  // applies to Ubuntu 16.04 only
  "ubuntu:16.04": "docker.io",
}
```