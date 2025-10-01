# :warning: Discontinued

The `docker sbom` command has been removed, please use the [`docker scout sbom` command](https://docs.docker.com/reference/cli/docker/scout/sbom/)
instead.

# sbom-cli-plugin

Plugin for Docker CLI to support viewing and creating SBOMs for Docker images using Syft.

## Getting started

```
# install the docker-sbom plugin
curl -sSfL https://raw.githubusercontent.com/docker/sbom-cli-plugin/main/install.sh | sh -s --

# use the sbom plugin
docker sbom <my-image>
```
