package main

import (
	"context"
	"dagger.io/dagger"
	"github.com/dchest/uniuri"
	"main/internal/pkg/build"
)

func buildOutput(buildParameters *build.BuildParameters, ctx context.Context, client *dagger.Client) (*dagger.Container, error) {
	var container *dagger.Container

	// Start container
	container, err := client.Container().
		From(buildParameters.ContainerImage).
		Sync(ctx)

	if err != nil {
		return container, err
	}

	// Bust cache if required
	if buildParameters.NoCache {
		container, err = container.
			WithEnvVariable("BURST_CACHE", uniuri.New()).
			Sync(ctx)

		if err != nil {
			return container, err
		}
	}

	return container, nil
}
