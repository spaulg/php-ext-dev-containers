#!/bin/bash
set -eux

# Fix degraded systemd preventing cgroups status in pbuilder
systemctl reset-failed

# Install build tools
apt update -y
DEBIAN_FRONTEND=noninteractive apt install -y \
    build-essential \
    devscripts \
    debhelper \
    dh-make \
    vim \
    quilt \
    debootstrap \
    dh-apache2 \
    debian-archive-keyring \
    jq \
    ant \
    ubuntu-dev-tools \
    sbuild

# Add vagrant user to sbuild group
sudo usermod -a -G sbuild vagrant

## Add configuration files
cp /home/vagrant/php-ext-dev-containers/configs/quiltrc.conf /home/vagrant/.quiltrc
cp /home/vagrant/php-ext-dev-containers/configs/sbuildrc.pm /home/vagrant/.sbuildrc
