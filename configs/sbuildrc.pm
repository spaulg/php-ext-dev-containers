# Name to use as override in .changes files for the Maintainer: field
# (mandatory, no default!).
$maintainer_name = 'Simon Paulger <spaulger@codezen.co.uk>';

# Default distribution to build.
$distribution = "bullseye";

# Build arch-all by default.
$build_arch_all = 1;

# When to purge the build directory afterwards
$purge_build_directory = 'successful';
$purge_session = 'successful';
$purge_build_deps = 'successful';

# Logging options
$verbose = 1;
$nolog = 0;
$log_dir = $ENV{'SBUILD_LOGS_PATH'} ? $ENV{'SBUILD_LOGS_PATH'} : $ENV{'HOME'};

# Disable apt updates/upgrades when building
$apt_update = 0;
$apt_upgrade = 0;

1;
