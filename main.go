package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

const MinimalAPIVersion = "1.12"

func main() {
	plugin.Run(func(_ command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Short: "Docker Nuke",
			Long:  `A tool to remove !!!EVERYTHING!!! from your runtime and registry.`,
			Use:   "nuke",
			Run:   nuke,
		}
		originalPreRun := cmd.PersistentPreRunE
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			if originalPreRun != nil {
				return originalPreRun(cmd, args)
			}
			return nil
		}
		return cmd
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Ulysses Souza",
		Version:       "1.0.0",
	})
}

func nuke(_ *cobra.Command, _ []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion(MinimalAPIVersion))
	if err != nil {
		panic(err)
	}

	cs, err := listContainers(ctx, cli, false)
	if err != nil {
		panic(err)
	}
	if len(cs) > 0 && ask4Confirm("Remove all containers?") {
		err = removeContainers(ctx, cli, cs)
		if err != nil {
			panic(err)
		}
	}

	is, err := listImages(ctx, cli, false)
	if err != nil {
		panic(err)
	}
	originalIsCount := len(is)

	okForRemove := false
	for len(is) > 0 && (okForRemove || ask4Confirm("Remove all images?")) {
		okForRemove = true
		err = removeImages(ctx, cli, is)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		is, err = listImages(ctx, cli, true)
		if err != nil {
			panic(err)
		}
	}

	if len(cs)+originalIsCount == 0 {
		fmt.Println("Nothing to remove...")
	}
}

func listContainers(ctx context.Context, cli client.APIClient, quiet bool) ([]types.Container, error) {
	listOpts := types.ContainerListOptions{All: true}
	cs, err := cli.ContainerList(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	if !quiet {
		fmt.Printf("Container count: %d\n", len(cs))
		for _, c := range cs {
			fmt.Printf("Container ID(%s), NAMES(%v)\n", c.ID, c.Names)
		}
	}
	return cs, nil
}

func listImages(ctx context.Context, cli client.APIClient, quiet bool) ([]types.ImageSummary, error) {
	listOpts := types.ImageListOptions{All: true}
	is, err := cli.ImageList(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	if !quiet {
		fmt.Printf("Image count: %d\n", len(is))
		for _, i := range is {
			fmt.Printf("Image ID(%s)\n", i.ID)
		}
	}
	return is, nil
}

func removeContainers(ctx context.Context, cli client.APIClient, cs []types.Container) error {
	fmt.Println("Removing containers...")
	removeOpts := types.ContainerRemoveOptions{Force: true}
	for _, c := range cs {
		err := cli.ContainerRemove(ctx, c.ID, removeOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeImages(ctx context.Context, cli client.APIClient, is []types.ImageSummary) error {
	removeOpts := types.ImageRemoveOptions{Force: true, PruneChildren: true}
	for _, i := range is {
		_, err := cli.ImageRemove(ctx, i.ID, removeOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func ask4Confirm(q string) bool {
	var s string
	fmt.Printf("%s (y/N): ", q)
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "y" || s == "yes"
}
