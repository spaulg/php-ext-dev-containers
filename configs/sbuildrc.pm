# Name to use as override in .changes files for the Maintainer: field
# (mandatory, no default!).
$maintainer_name='Your Name <user@example.org>';

# Default distribution to build.
$distribution = "bullseye";

# Build arch-all by default.
$build_arch_all = 1;

# When to purge the build directory afterwards
$purge_build_directory = 'successful';
$purge_session = 'successful';
$purge_build_deps = 'successful';

# Directory for writing build logs to
$log_dir=$ENV{HOME}."/php-ext-dev-containers/logs";

# Disable apt updates/upgrades when building
$apt_update = 0;
$apt_upgrade = 0;

1;
