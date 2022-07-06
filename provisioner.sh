#!/bin/bash
set -eux

# Fix degraded systemd preventing cgroups status
systemctl reset-failed

# Add docker upstream repository
apt update -y
DEBIAN_FRONTEND=noninteractive apt install -y \
    ca-certificates \
    curl \
    gnupg

install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

apt update -y

# Install build tools
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
    sbuild \
    docker-ce \
    docker-ce-cli \
    containerd.io \
    docker-buildx-plugin \
    docker-compose-plugin \
    git \
    bindfs

# Add vagrant user to secondary groups
sudo usermod -a -G sbuild vagrant
sudo usermod -a -G docker vagrant

## Add configuration files
cp /home/vagrant/php-ext-dev-containers/configs/quiltrc.conf /home/vagrant/.quiltrc
cp /home/vagrant/php-ext-dev-containers/configs/sbuildrc.pm /home/vagrant/.sbuildrc

cat <<'EOF' | tee -a /home/vagrant/.bashrc
export SBUILD_LOGS_PATH="$HOME/php-ext-dev-containers/logs"
export VISUAL=vim
export EDITOR="$VISUAL"

if [ -e "$HOME/php-ext-dev-containers" ]; then
  cd php-ext-dev-containers
fi

git config --global alias.st status
EOF

rm -f /etc/localtime
ln -s /usr/share/zoneinfo/Europe/London /etc/localtime
