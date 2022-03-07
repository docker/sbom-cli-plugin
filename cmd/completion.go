package cmd

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func dockerImageValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// since we use ValidArgsFunction, Cobra will call this AFTER having parsed all flags and arguments provided
	dockerImageRepoTags, err := listLocalDockerImages(toComplete)
	if err != nil {
		// indicates that an error occurred and completions should be ignored
		return []string{"completion failed"}, cobra.ShellCompDirectiveError
	}
	if len(dockerImageRepoTags) == 0 {
		return []string{"no docker images found"}, cobra.ShellCompDirectiveError
	}
	// ShellCompDirectiveDefault indicates that the shell will perform its default behavior after completions have
	// been provided (without implying other possible directives)
	return dockerImageRepoTags, cobra.ShellCompDirectiveDefault
}

func listLocalDockerImages(prefix string) ([]string, error) {
	var repoTags = make([]string, 0)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return repoTags, err
	}

	// only want to return tagged images
	imageListArgs := filters.NewArgs()
	imageListArgs.Add("dangling", "false")
	images, err := cli.ImageList(ctx, types.ImageListOptions{All: false, Filters: imageListArgs})
	if err != nil {
		return repoTags, err
	}

	for _, image := range images {
		// image may have multiple tags
		for _, tag := range image.RepoTags {
			if strings.HasPrefix(tag, prefix) {
				repoTags = append(repoTags, tag)
			}
		}
	}
	return repoTags, nil
}
