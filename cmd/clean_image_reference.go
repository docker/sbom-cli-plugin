package cmd

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

func cleanImageReference(userInput string) (string, error) {
	ref, err := name.ParseReference(userInput, name.WeakValidation, name.WithDefaultTag("latest"))
	if err != nil {
		return "", fmt.Errorf("unable to parse image reference: %w", err)
	}

	if t, ok := ref.(name.Tag); ok {
		if !strings.HasSuffix(userInput, t.Identifier()) {
			return userInput + ":" + t.Identifier(), nil
		}
		return userInput, nil
	}

	if d, ok := ref.(name.Digest); ok {
		if !strings.HasSuffix(userInput, d.Identifier()) {
			return userInput + "@" + d.Identifier(), nil
		}
		return userInput, nil
	}

	return ref.Name(), nil
}
