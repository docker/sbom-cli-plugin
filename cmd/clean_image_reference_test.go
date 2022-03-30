package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_cleanImageReference(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{
			input: "alpine:latest",
			want:  "alpine:latest",
		},
		{
			input: "alpine",
			want:  "alpine:latest",
		},
		{
			input: "docker",
			want:  "docker:latest",
		},
		{
			input: "anchore/syft:latest",
			want:  "anchore/syft:latest",
		},
		{
			input: "anchore/syft:v1.4.5",
			want:  "anchore/syft:v1.4.5",
		},
		{
			input: "anchore/syft",
			want:  "anchore/syft:latest",
		},
		{
			input: "docker.io/anchore/syft",
			want:  "docker.io/anchore/syft:latest",
		},
		{
			input: "registry.upbound.io/crossplane/provider-gcp:stable",
			want:  "registry.upbound.io/crossplane/provider-gcp:stable",
		},
		{
			input: "registry.upbound.io/crossplane/provider-gcp",
			want:  "registry.upbound.io/crossplane/provider-gcp:latest",
		},
		{
			input: "anchore/syft@sha256:dba09c285770f58d6685b25a0606d72420b0a7525a2338080807d138a258c671",
			want:  "anchore/syft@sha256:dba09c285770f58d6685b25a0606d72420b0a7525a2338080807d138a258c671",
		},
		{
			// mix tag and digest
			input: "anchore/syft:latest@sha256:8bbaebbd4bfc3fed46227eba1d49643fc1bb79b23378956f96cff4c5d69dd42b",
			want:  "anchore/syft:latest@sha256:8bbaebbd4bfc3fed46227eba1d49643fc1bb79b23378956f96cff4c5d69dd42b",
		},
		{
			// mix tag and digest
			input: "registry.upbound.io/crossplane/provider-gcp:v0.2.0@sha256:8bbaebbd4bfc3fed46227eba1d49643fc1bb79b23378956f96cff4c5d69dd42b",
			want:  "registry.upbound.io/crossplane/provider-gcp:v0.2.0@sha256:8bbaebbd4bfc3fed46227eba1d49643fc1bb79b23378956f96cff4c5d69dd42b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if tt.wantErr == nil {
				tt.wantErr = require.NoError
			}
			got, err := cleanImageReference(tt.input)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
