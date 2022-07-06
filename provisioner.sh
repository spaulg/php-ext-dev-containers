#!/bin/bash
set -eux

# Fix degraded systemd preventing cgroups status in pbuilder
systemctl reset-failed

# Enable source repositories
#sed -i -e 's/#\s*\(deb-src.*\)/\1/' /etc/apt/sources.list

# Install build tools
apt update -y
DEBIAN_FRONTEND=noninteractive apt install -y \
    build-essential \
    devscripts \
    debhelper \
    dh-make \
    pbuilder \
    vim \
    quilt \
    debootstrap \
    dh-apache2 \
    debian-archive-keyring \
    jq \
    ant

## Add bullseye debian repositories
#cat <<EOF > /etc/apt/sources.list.d/debian-sources.list
#deb [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://deb.debian.org/debian bullseye main
#deb-src [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://deb.debian.org/debian bullseye main
#
#deb [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://security.debian.org/debian-security bullseye-security main
#deb-src [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://security.debian.org/debian-security bullseye-security main
#
#deb [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://deb.debian.org/debian bullseye-updates main
#deb-src [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://deb.debian.org/debian bullseye-updates main
#EOF



# Add quilt configuration file
cat <<-'EOF' > /etc/quilt.quiltrc
# Diff options
QUILT_DIFF_OPTS="--show-c-function"
QUILT_NO_DIFF_INDEX=1
QUILT_NO_DIFF_TIMESTAMPS=1

# Patch options
QUILT_PATCHES_PREFIX=yes
QUILT_PATCHES=debian/patches
QUILT_PATCH_OPTS="--reject-format=unified"

# Options to pass to commands (QUILT_${COMMAND}_ARGS)
QUILT_PUSH_ARGS="--color=auto"
QUILT_SERIES_ARGS="--color=auto"
QUILT_PATCHES_ARGS="--color=auto"
QUILT_REFRESH_ARGS="-p ab"
QUILT_DIFF_ARGS="--color=auto"

# When non-default less options are used, add the -R option so that less outputs
# ANSI color escape codes "raw".
[ -n "$LESS" -a -z "${QUILT_PAGER+x}" ] && QUILT_PAGER="less -FRX"

QUILT_COLORS="diff_hdr=1;32:diff_add=1;34:diff_rem=1;31:diff_hunk=1;33:diff_ctx=35:diff_cctx=33"
EOF



# Add pbuilder hook script
mkdir -p /var/cache/pbuilder/hook.d

cat <<EOF > /var/cache/pbuilder/hook.d/C10shell
#!/bin/sh
# invoke shell if build fails.

apt-get install -y --force-yes vim less bash
cd /tmp/buildd/*/debian/..
/bin/bash < /dev/tty > /dev/tty 2> /dev/tty
EOF

chmod a+x /var/cache/pbuilder/hook.d/C10shell


# Add pbuilder configuration file
cat <<EOF >> /etc/pbuilderrc
PBUILDERSATISFYDEPENDSCMD="/usr/lib/pbuilder/pbuilder-satisfydepends-apt"
HOOKDIR="/var/cache/pbuilder/hook.d"
EOF

touch /root/.pbuilderrc
