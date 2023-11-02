VERSION 0.7

build:
    ARG version="8.2.6"
    ARG suffix=""
    ARG arch="amd64"
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
    COPY build/${package_name}_${version}.orig.tar.gz /home/build/packages/${package_name}_${version}.orig.tar.gz
    RUN tar -xzf /home/build/packages/${package_name}_${version}.orig.tar.gz --strip-components=1 --exclude debian

    # Regenerate changelog with version passed
    RUN rm -f debian/changelog
    RUN debchange --create --package "${package_name}" --distribution stable -v "${version}-${build_number}" "${version}-${build_number} automated build"

    # Add the target architecture and install build dependencies
    RUN sudo dpkg --add-architecture ${arch}
    RUN sudo apt update -y

    RUN sudo mk-build-deps -i -r -t "apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends -y" --host-arch ${arch}
    RUN rm -f php*-build-deps_*.buildinfo php*-build-deps_*.changes

    # Generate missing build files from templates
    RUN make -f debian/rules prepare

    # Build the package
    RUN debuild -us -uc -a${arch} > /tmp/build.log 2>&1; echo $? > /tmp/build.status

    SAVE ARTIFACT /tmp/build.log AS LOCAL build/
    SAVE ARTIFACT /tmp/build.statis AS LOCAL build/
    SAVE ARTIFACT /home/build/packages/*.deb AS LOCAL build/
