package main

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"log"
	"strconv"
)

func build(buildParameters *BuildParameters, ctx context.Context, client *dagger.Client) (*dagger.Container, error) {
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

	// Prepare package
	container, err = container.
		WithExec([]string{"false"}).
		WithExec([]string{"cp", "/home/build/source/" + sourceArchiveFileName, "/home/build/packages/" + sourceArchiveFileName}).
		WithExec([]string{"tar", "-xzf", "/home/build/packages/" + sourceArchiveFileName, "--strip-components=1", "--exclude", "debian"}).
		WithExec([]string{"cp", "-R", "/home/build/source/" + buildParameters.ShortVersion, buildParameters.BuildDirectoryPath + "/debian"}).
		WithExec([]string{"rm", "-f", "debian/changelog"}).
		WithExec([]string{"debchange", "--create", "--package", buildParameters.PackageName, "--distribution", "stable", "-v", buildParameters.Version + "-" + strconv.Itoa(buildParameters.BuildNumber), buildParameters.Version + "-" + strconv.Itoa(buildParameters.BuildNumber) + " automated build"}).
		WithExec([]string{"make", "-f", "debian/rules", "prepare"}).
		WithExec([]string{"sudo", "dpkg", "--add-architecture", buildParameters.Architecture}).
		WithExec([]string{"sudo", "apt", "update", "-y"}).
		WithExec([]string{"sudo", "mk-build-deps", "-i", "-t", "apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends -y", "--host-arch", buildParameters.Architecture}).
		Sync(ctx)

	if err != nil {
		return container, err
	}

	// Clean mk-build-deps files and delete
	buildDirectory := container.Directory(buildParameters.BuildDirectoryPath)
	var removeFiles []string

	for _, globPattern := range []string{"**.deb", "**.changes", "**.buildinfo"} {
		globFiles, err := buildDirectory.Glob(ctx, globPattern)

		if err != nil {
			return container, fmt.Errorf("unable to list glob files for cleanup: %v", err)
		}

		for _, file := range globFiles {
			file = "/" + file

			log.Println("Removing file: " + file)
			removeFiles = append(removeFiles, file)
		}
	}

	if len(removeFiles) > 0 {
		container, err = container.
			WithExec(append([]string{"rm", "-f"}, removeFiles...)).
			Sync(ctx)

		if err != nil {
			return container, err
		}
	}

	// Final build
	return container.
		WithExec([]string{"debuild", "-us", "-uc", "-a" + buildParameters.Architecture}).
		Sync(ctx)
}
