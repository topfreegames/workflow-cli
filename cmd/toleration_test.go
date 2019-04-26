package cmd

import (
	"bytes"
	"fmt"
	"k8s.io/api/core/v1"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestTolerationList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "owner": "jkirk",
    "app": "foo",
    "tolerations": {
        "cmd": {
			"toleration-test": {
				"key": "somekey",
				"operator": "Equal",
				"value": "somevalue",
				"effect": "NoSchedule",
				"tolerationSeconds": 300
			}
        }
    },
    "values": {},
    "memory": {},
    "cpu": {},
    "tags": {},
    "registry": {},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.TolerationList("foo", "")
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `=== cmd Tolerations
---- toleration-test
Effect                  NoSchedule
Key                     somekey
Operator                Equal
Toleration Seconds      300
Value                   somevalue
`, "output")
	b.Reset()

	err = cmdr.TolerationList("foo", "oneline")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(),"cmd: toleration-test|Key=somekey,Operator=Equal,Value=somevalue,Effect=NoSchedule,TolerationSeconds=300;\n", "output")

	b.Reset()

	err = cmdr.TolerationList("foo", "diff")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `cmd:
---- toleration-test
    Key=somekey
    Operator=Equal
    Value=somevalue
    Effect=NoSchedule
    TolerationSeconds=300
`, "output")
}

func TestTolerationSet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	var testSeconds int64 = 300
	var mockToleration = v1.Toleration{
		Key:               "somekey",
		Value:             "somevalue",
		Operator:          v1.TolerationOpEqual,
		TolerationSeconds: &testSeconds,
		Effect:            v1.TaintEffectNoSchedule,
	}
	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Tolerations: map[string]map[string]*v1.Toleration{
					"cmd": {
						"annotation-test": &mockToleration,
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
    "tolerations": {
        "cmd": {
            "annotation-test": {
				"key": "somekey",
				"value": "somevalue",
				"operator": "Equal",
				"tolerationSeconds": 300,
				"effect": "NoSchedule",
			}
        }
    }
}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.TolerationSet("foo", "cmd", "toleration-test", mockToleration)
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Creating Annotations... done

=== cmd Tolerations
---- toleration-test
Effect                  NoSchedule
Key                     somekey
Operator                Equal
Toleration Seconds      300
Value                   somevalue
`, "output")
}

func TestTolerationUnset(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Tolerations: map[string]map[string]*v1.Toleration{
					"cmd": {
						"toleration": nil,
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
    "tolerations": {
        "cmd": {
            "annotation-test": {
				"key": "somekey",
				"value": "somevalue",
				"operator": "Equal",
				"tolerationSeconds": 300,
				"effect": "NoSchedule",
			}
        }
    }
}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.TolerationUnset("foo", "cmd", []string{"toleration"})
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Removing Annotations... done

=== cmd Tolerations
---- toleration-test
Effect                  NoSchedule
Key                     somekey
Operator                Equal
Toleration Seconds      300
Value                   somevalue
`, "output")
}

