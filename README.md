# docker-sbom-cli-plugin

Plugin for Docker CLI to support SBOM creation using Syft.

**Note: this is a work in progress**

## Getting started

```
# build
make snapshot

# install
cp snapshot/<path/to/your>/docker-sbom ~/.docker/cli-plugins

# use
docker sbom <my-image>
```
