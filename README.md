# sbom-cli-plugin

Plugin for Docker CLI to support viewing and creating SBOMs for Docker images using Syft.

**Note: this is a proof of concept / work in progress**

## Getting started

```
# install the docker-sbom plugin
curl -sSfL https://raw.githubusercontent.com/docker/sbom-cli-plugin/main/install.sh | sh -s --

# use the sbom plugin
docker sbom <my-image>
```
