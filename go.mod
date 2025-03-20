module github.com/docker/sbom-cli-plugin

go 1.18
toolchain go1.24.1

require (
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	github.com/anchore/stereoscope v0.0.0-20220518185348-c97a3c6ffc67
	github.com/anchore/syft v0.46.3
	github.com/containerd/containerd v1.5.10 // indirect
	github.com/docker/cli v20.10.12+incompatible
	github.com/docker/docker v25.0.6+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/fvbommel/sortorder v1.0.2 // indirect
	github.com/gookit/color v1.4.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.10.0
	github.com/theupdateframework/notary v0.7.0 // indirect
	github.com/wagoodman/go-partybus v0.0.0-20210627031916-db1f5573bbc5
	github.com/wagoodman/jotframe v0.0.0-20211129225309-56b0d0a4aebb
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	golang.org/x/term v0.29.0
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.1.0 // indirect
)

require github.com/google/go-containerregistry v0.8.1-0.20220209165246-a44adc326839

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/CycloneDX/cyclonedx-go v0.5.2 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/acobaugh/osrelease v0.1.0 // indirect
	github.com/anchore/go-macholibre v0.0.0-20220308212642-53e6d0aaf6fb // indirect
	github.com/anchore/go-rpmdb v0.0.0-20210914181456-a9c52348da63 // indirect
	github.com/anchore/go-testutils v0.0.0-20200925183923-d5f45b0d3c04 // indirect
	github.com/anchore/packageurl-go v0.1.1-0.20220428202044-a072fa3cb6d7 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.0.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.10.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/distribution v2.8.0+incompatible // indirect
	github.com/docker/go v1.5.1-1.0.20160303222718-d30aec9fd63c // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/facebookincubator/nvdtools v0.1.4 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-restruct/restruct v1.2.0-alpha // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jinzhu/copier v0.3.2 // indirect
	github.com/klauspost/compress v1.14.2 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/logrusorgru/aurora v0.0.0-20200102142835-e9ef32dff381 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mholt/archiver/v3 v3.5.1 // indirect
	github.com/miekg/pkcs11 v1.1.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pelletier/go-toml/v2 v2.0.0-beta.8 // indirect
	github.com/pierrec/lz4/v4 v4.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/scylladb/go-set v1.0.3-0.20200225121959-cc7b2070d91e // indirect
	github.com/spdx/tools-golang v0.2.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vbatts/tar-split v0.11.2 // indirect
	github.com/vifraa/gopom v0.1.0 // indirect
	github.com/wagoodman/go-progress v0.0.0-20200731105512-1020f39e6240 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.60.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/genproto v0.0.0-20220414192740-2d67ff6cf2b4 // indirect
	google.golang.org/grpc v1.71.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
