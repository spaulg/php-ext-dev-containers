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

### Building

Earthly is used to compile PHP in to a debian package and build the docker image 
in a two stage process.

Non zts:

```
earthly +build --platform=amd64 --version=8.2.6
earthly +package --platform=amd64 --version=8.2.6
```

ZTS:

```
earthly +build --platform=amd64 --version=8.2.6 --suffix=zts
earthly +package --platform=amd64 --version=8.2.6 --suffix=zts
```
