package main

import (
	"dagger.io/dagger"
	"log"
)

func main() {
	// Parse command arguments to capture build information
	buildParameters := parseArguments()

	// Connect Dagger client
	context, client := connectDaggerClient()
	defer client.Close()

	// Build package
	container, err := build(buildParameters, context, client)

	if err != nil {
		log.Println(err)
		return
	}

	// Export build log
	_, err = container.File("/home/build/packages/build.log").Export(context, "output", dagger.FileExportOpts{AllowParentDirPath: true})

	if err != nil {
		log.Println(err)
		return
	}

	// Export debian package files
	directory := container.Directory("/home/build/packages")
	files, err := directory.Glob(context, "**.deb")

	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range files {
		_, err = container.
			Directory("/").
			File(file).
			Export(context, "output", dagger.FileExportOpts{AllowParentDirPath: true})

		if err != nil {
			log.Println(err)
			return
		}
	}
}
