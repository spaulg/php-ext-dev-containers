"""These are constants used in detecting build versions to report for building"""


# PHP releases API endpoint
PHP_RELEASE_API = "https://www.php.net/releases/index.php?json"

# Minimum PHP version number to build, ignoring
# anything older
PHP_MIN_MAJOR_VERSION = 5
PHP_MIN_MINOR_VERSION = 6

# Maximum age for a build in days before it must
# be rebuilt regardless of whether a newer version
# is available
MAX_AGE_IN_DAYS = 0
