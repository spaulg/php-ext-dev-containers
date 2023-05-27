
# Introduction

This repository is a suite of scripts/tools for building docker containers
of debug builds of each major/minor version combination of PHP.

The docker containers are pushed to Docker Hub and refreshed on a weekly basis
to update to the latest PHP patch releases and obtain the latest updates from
the base images (Debian Bullseye).

Objectives of the project are to:

1. create builds of each major/minor version of PHP
2. use the same base OS version and toolchain to build each version
3. use the `--debug` flag for each build
4. produce builds for both arm64 and amd64 architectures
5. automate the process

### Local development environment

The local development environment uses Vagrant with Mutagen for synchronisation of files.
This is limited to Mac or Linux.

#### Installation

1. Install Vagrant
2. Install a hypervisor (Virtualbox or Parallels)
3. Install Mutagen

#### Starting the environment

```
vagrant up
vagrant ssh

ant init
```

### Building

Use the Apache ant build script, passing the desired PHP major/minor version combination.

This will download and build the latest patch release for that version, create the .deb
files for both arm64 and amd64 architectures and package in to a container image. 

```
cd php-ext-dev-containers
ant build -Dversion=7.4
```
