VERSION 0.7

build:
    ARG version=""
    ARG suffix=""
    ARG architecture="amd64"
    ARG distribution="bullseye"
    ARG build_number="1"
    ARG log_filename=""
    ARG status_filename=""

    FROM github.com/spaulg/earthly-debuilder+image --distribution=${distribution}

    # Parse the full version to get the short version for the packagen ame
    ENV short_version="$(echo "${version}" | awk -F \. {'print $1"."$2'})"
    ENV package_name="php${short_version}${suffix}"
    ENV env_log_filename="$(test "${log_filename}" != "" && echo "${log_filename}" || echo "${package_name}.build.${architecture}.log")"
    ENV env_status_filename="$(test "${status_filename}" != "" && echo "${status_filename}" || echo "${package_name}.build.${architecture}.status")"

    # Copy debian packaging sources
    COPY packages/${short_version} /home/build/packages/${package_name}_${version}/debian
    WORKDIR /home/build/packages/${package_name}_${version}

    # Copy upstream sources
    COPY build/php-${version}.tar.gz /home/build/packages/${package_name}_${version}.orig.tar.gz
    RUN tar -xzf /home/build/packages/${package_name}_${version}.orig.tar.gz --strip-components=1 --exclude debian

    # Regenerate changelog with version passed
    RUN rm -f debian/changelog
    RUN debchange --create --package "${package_name}" --distribution stable -v "${version}-${build_number}" "${version}-${build_number} automated build"

    # Generate build files from templates
    RUN make -f debian/rules prepare

    # Add the target architecture and install build dependencies
    RUN sudo dpkg --add-architecture ${architecture}
    RUN sudo apt update -y

    RUN sudo mk-build-deps -i -r -t "apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends -y" --host-arch ${architecture}
    RUN rm -f php*-build-deps_*.buildinfo php*-build-deps_*.changes

    # Build the package
    RUN debuild -us -uc -a${architecture} > /tmp/build.${architecture}.log 2>&1; echo $? > /tmp/build.${architecture}.status

    SAVE ARTIFACT --if-exists /home/build/packages/*.deb AS LOCAL ./output/
    SAVE ARTIFACT /tmp/build.${architecture}.log AS LOCAL ./output/${env_log_filename}
    SAVE ARTIFACT /tmp/build.${architecture}.status AS LOCAL ./output/${env_status_filename}
