package main

import (
	"context"
	"dagger.io/dagger"
	"errors"
	"log"
	"os"
)

func main() {
	os.Exit(RunBuild())
}

func RunBuild() int {
	// Parse command arguments to capture build information
	buildParameters := parseArguments()

	// Connect Dagger client
	ctx, client := connectDaggerClient()
	defer client.Close()

	// Build package
	container, err := build(buildParameters, ctx, client)

	if err != nil {
		log.Println(err)
	}

	return ExportArtifacts(buildParameters, ctx, container)
}

func ExportArtifacts(buildParameters *BuildParameters, ctx context.Context, container *dagger.Container) int {
	// Create the output directory if it does not exist
	if _, err := os.Stat("output"); errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir("output", 0750); err != nil {
			log.Println(err)
			return 1
		}
	}

	// Export build log
	exitCode := 0
	_, err := container.File("/home/build/packages/build.log").Export(ctx, "output/"+buildParameters.LogFileName)

	if err != nil {
		exitCode = 1
	}

	// Export debian package files
	directory := container.Directory("/home/build/packages")
	files, err := directory.Glob(ctx, "**.deb")

	if err != nil {
		exitCode = 1
	} else {
		for _, file := range files {
			_, err = container.
				Directory("/").
				File(file).
				Export(ctx, "output", dagger.FileExportOpts{AllowParentDirPath: true})

			if err != nil {
				log.Println(err)
				exitCode = 1
			}
		}
	}

	return exitCode
}
