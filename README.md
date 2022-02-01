# docker-sbom-cli-plugin

Plugin for Docker CLI to support SBOM creation using Syft.

**Note: this is a proof of concept / work in progress**

## Getting started

```
# install syft
curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin

# build
make snapshot

# install
cp snapshot/<path/to/your>/docker-sbom ~/.docker/cli-plugins

# use
docker sbom <my-image>
```
