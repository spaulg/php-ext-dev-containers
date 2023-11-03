VERSION 0.7

build:
    ARG version=""
    ARG suffix=""
    ARG architecture="amd64"
    ARG distribution="bullseye"
    ARG build_number="1"

    FROM github.com/spaulg/earthly-debuilder+image --distribution=${distribution}

    # Parse the full version to get the short version for the packagen ame
    ENV short_version="$(echo "${version}" | awk -F \. {'print $1"."$2'})"
    ENV package_name="php${short_version}"

    # Copy debian packaging sources
    COPY packages/${short_version} /home/build/packages/${package_name}_${version}${suffix}/debian
    WORKDIR /home/build/packages/${package_name}_${version}${suffix}

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

    SAVE ARTIFACT /tmp/build.${architecture}.log AS LOCAL build/${package_name}.build.${architecture}.log
    SAVE ARTIFACT /tmp/build.${architecture}.status AS LOCAL build/${package_name}.build.${architecture}.status
    SAVE ARTIFACT /home/build/packages/*.deb AS LOCAL build/

    #RUN [ "$(cat /tmp/build.status)" -eq 0 ] || exit 1
