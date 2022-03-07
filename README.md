# docker-sbom-cli-plugin

Plugin for Docker CLI to support SBOM creation using Syft.

**Note: this is a proof of concept / work in progress**

## Getting started

```
# install the docker-sbom plugin
curl -sSfL https://raw.githubusercontent.com/anchore/docker-sbom-cli-plugin/main/install.sh | sh -s --

# use the sbom plugin
docker sbom <my-image>
```

## Configuration

Configuration search paths:

- `.syft.yaml`
- `.syft/config.yaml`
- `~/.syft.yaml`
- `<XDG_CONFIG_HOME>/syft/config.yaml`

Configuration options (example values are the default):

```yaml
# the output format(s) of the SBOM report (options: table, text, json, spdx, ...)
# same as -o, --output, and SYFT_OUTPUT env var
# to specify multiple output files in differing formats, use a list:
# output:
#   - "json=<syft-json-output-file>"
#   - "spdx-json=<spdx-json-output-file>"
output: "table"

# suppress all output (except for the SBOM report)
# same as -q ; SYFT_QUIET env var
quiet: false

# same as --file; write output report to a file (default is to write to stdout)
file: ""

# a list of globs to exclude from scanning. same as --exclude ; for example:
# exclude:
#   - "/etc/**"
#   - "./out/**/*.json"
exclude: []

# override the OS or architecture options (e.g. "windows/arm/v7", "linux/arm64", "arm64") similar to "docker --platform"
# same as --platform; SYFT_PLATFORM env var
platform: ""

# cataloging packages is exposed through the packages and power-user subcommands
package:

  # search within archives that do contain a file index to search against (zip)
  # note: for now this only applies to the java package cataloger
  # SYFT_PACKAGE_SEARCH_INDEXED_ARCHIVES env var
  search-indexed-archives: true

  # search within archives that do not contain a file index to search against (tar, tar.gz, tar.bz2, etc)
  # note: enabling this may result in a performance impact since all discovered compressed tars will be decompressed
  # note: for now this only applies to the java package cataloger
  # SYFT_PACKAGE_SEARCH_UNINDEXED_ARCHIVES env var
  search-unindexed-archives: false

  cataloger:
    # enable/disable cataloging of packages
    # SYFT_PACKAGE_CATALOGER_ENABLED env var
    enabled: true

    # the search space to look for packages (options: all-layers, squashed)
    # same as -s ; SYFT_PACKAGE_CATALOGER_SCOPE env var
    scope: "squashed"

# cataloging file metadata is exposed through the power-user subcommand
file-metadata:
  cataloger:
    # enable/disable cataloging of file metadata
    # SYFT_FILE_METADATA_CATALOGER_ENABLED env var
    enabled: true

    # the search space to look for file metadata (options: all-layers, squashed)
    # SYFT_FILE_METADATA_CATALOGER_SCOPE env var
    scope: "squashed"

  # the file digest algorithms to use when cataloging files (options: "sha256", "md5", "sha1")
  # SYFT_FILE_METADATA_DIGESTS env var
  digests: ["sha256"]

log:
  # use structured logging
  # same as SYFT_LOG_STRUCTURED env var
  structured: false

  # the log level; note: detailed logging suppress the ETUI
  # same as SYFT_LOG_LEVEL env var
  level: "error"

  # location to write the log file (default is not to have a log file)
  # same as SYFT_LOG_FILE env var
  file: ""

```

