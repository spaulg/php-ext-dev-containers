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

	if err != nil && buildParameters.Interactive {
		log.Println(err)

		RunInteractiveShell(buildParameters, ctx, container)
	}

	return ExportArtifacts(buildParameters, ctx, container)
}

func RunInteractiveShell(buildParameters *BuildParameters, ctx context.Context, container *dagger.Container) {
	//// todo: export the image in to a tar file
	//_, err := container.Export(context, ".image.tar")
	//
	//if err != nil {
	//	// todo: handle error
	//}
	//
	//// todo: load the tar using docker/podman
	//if err = exec.Command("docker", "load", "-qi", ".image.tar").Run(); err != nil {
	//	// todo: handle error
	//}
	//
	//// todo: run the image using docker run -it <image> sh
	//// todo: wait for the process to exit
	//if err = exec.Command("docker", "run", "-it", "--rm", "tagged:image", "sh").Run(); err != nil {
	//	// todo: handle error
	//}
	//
	//// todo: remove the image from docker
	//if err = exec.Command("docker", "rmi", "-i", "tagged:image").Run(); err != nil {
	//	// todo: handle error
	//}
	//
	//// todo: remove the file from disk
	//if err = os.Remove(".image.tar"); err != nil {
	//	// todo: handle error
	//}
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
