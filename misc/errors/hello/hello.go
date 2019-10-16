package hello

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var ErrMissingName = fmt.Errorf("missing name")

func Hi(name string) (string, error) {
	saneName, err := sanitizeName(name)
	if err != nil {
		return "", errors.Wrap(err, "cannot sanitize name")
	}
	return fmt.Sprintf("Hi, %s!", saneName), nil
}

func sanitizeName(name string) (string, error) {
	if name == "" {
		return "", ErrMissingName
	}
	return strings.TrimSpace(name), nil
}
