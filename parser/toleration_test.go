package parser

import (
	"bytes"
	"errors"
	"k8s.io/api/core/v1"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDeisCmd) TolerationList(string, string) error {
	return errors.New("toleration:list")
}

func (d FakeDeisCmd) TolerationSet(string, string, string, v1.Toleration) error {
	return errors.New("toleration:set")
}

func (d FakeDeisCmd) TolerationUnset(string, string, []string) error {
	return errors.New("toleration:unset")
}

func TestToleration(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := FakeDeisCmd{WOut: &b, ConfigFile: cf}

	// cases defines the arguments and expected return of the call.
	// if expected is "", it defaults to args[0].
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"toleration:list"},
			expected: "",
		},
		{
			args:     []string{"toleration:set", "--type=cmd", "var=value"},
			expected: "",
		},
		{
			args:     []string{"toleration:unset", "--type=cmd", "var"},
			expected: "",
		},
		{
			args:     []string{"toleration"},
			expected: "toleration:list",
		},
	}

	// For each case, check that calling the route with the arguments
	// returns the expected error, which is args[0] if not provided.
	for _, c := range cases {
		var expected string
		if c.expected == "" {
			expected = c.args[0]
		} else {
			expected = c.expected
		}
		err = Toleration(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
