package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		fmt.Printf("Error connecting to Dagger Engine: %s", err)
		os.Exit(1)
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	// mount cloned repository into `golang` image
	golang := client.Container().From("golang:latest")
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src").WithEnvVariable("CGO_ENABLED", "0")

	// define the application build command
	path := "build/"
	golang = golang.WithExec([]string{"go", "build", "-o", path})

	// get reference to build output directory in container
	output := golang.Directory(path)

	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, path)
	if err != nil {
		fmt.Printf("Error writing build directory to the host: %s", err)
		os.Exit(1)
	}

	cn, _ := client.Container().Build(src).Publish(ctx, "mwittek/dagger-example:latest")
	if err != nil {
		fmt.Printf("Error building and publishing container: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully built and published container: %s", cn)
}
