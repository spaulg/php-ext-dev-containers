package main

import (
	"context"
	"dagger.io/dagger"
	"strconv"
)

func build(buildParameters *BuildParameters, context context.Context, client *dagger.Client) (*dagger.Container, error) {
	var container *dagger.Container

	// Start container
	container = client.Container().
		From(buildParameters.ContainerImage).
		WithDirectory("/home/build/source", client.Host().Directory("assets/source")).
		WithExec([]string{"mkdir", "-p", buildParameters.BuildDirectoryPath}).
		WithWorkdir(buildParameters.BuildDirectoryPath)

	// Download source archive
	sourceArchiveFileName, err := downloadSourceArchive(buildParameters)

	if err != nil {
		return container, err
	}

	buildDepsFileName := buildParameters.PackageName + "-build-deps" + "_" + buildParameters.Version + "-" + strconv.Itoa(buildParameters.BuildNumber) + "_" + buildParameters.Architecture
	buildDepsInfoFileName := buildDepsFileName + ".buildinfo"
	buildDepsChangesFileName := buildDepsFileName + ".changes"
	buildDepsDebFileName := buildDepsFileName + ".deb"

	// Build package
	container = container.
		WithExec([]string{"cp", "/home/build/source/" + sourceArchiveFileName, "/home/build/packages/" + sourceArchiveFileName}).
		WithExec([]string{"tar", "-xzf", "/home/build/packages/" + sourceArchiveFileName, "--strip-components=1", "--exclude", "debian"}).
		WithExec([]string{"cp", "-R", "/home/build/source/" + buildParameters.ShortVersion, buildParameters.BuildDirectoryPath + "/debian"}).
		WithExec([]string{"rm", "-f", "debian/changelog"}).
		WithExec([]string{"debchange", "--create", "--package", buildParameters.PackageName, "--distribution", "stable", "-v", buildParameters.Version + "-" + strconv.Itoa(buildParameters.BuildNumber), buildParameters.Version + "-" + strconv.Itoa(buildParameters.BuildNumber) + " automated build"}).
		WithExec([]string{"make", "-f", "debian/rules", "prepare"}).
		WithExec([]string{"sudo", "dpkg", "--add-architecture", buildParameters.Architecture}).
		WithExec([]string{"sudo", "apt", "update", "-y"}).
		WithExec([]string{"sudo", "mk-build-deps", "-i", "-t", "apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends -y", "--host-arch", buildParameters.Architecture}).
		WithExec([]string{"rm", "-f", buildDepsInfoFileName, buildDepsChangesFileName, buildDepsDebFileName}).
		WithExec([]string{"debuild", "-us", "-uc", "-a" + buildParameters.Architecture}, dagger.ContainerWithExecOpts{
			RedirectStdout: "/home/build/packages/build.log",
			RedirectStderr: "/home/build/packages/build.log",
		})

	return container.Sync(context)
}
