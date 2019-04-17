package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestAnnotationList(t *testing.T) {
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
    "annotations": {
        "cmd": {
            "k8s.annotation": "testing",
            "k8s.another/annotation": "anotherone",
            "k8s.json/annotation": "{\"hello\":\"world\"}"
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

	err = cmdr.AnnotationList("foo", "")
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `=== cmd Annotations
k8s.annotation              testing
k8s.another/annotation      anotherone
k8s.json/annotation         {"hello":"world"}
`, "output")
	b.Reset()

	err = cmdr.AnnotationList("foo", "oneline")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(),"cmd: k8s.annotation=testing k8s.another/annotation=anotherone k8s.json/annotation={\"hello\":\"world\"}\n", "output")

	b.Reset()

	err = cmdr.AnnotationList("foo", "diff")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(),`cmd:
    k8s.annotation=testing
    k8s.another/annotation=anotherone
    k8s.json/annotation={"hello":"world"}
`, "output")
}

func TestAnnotationsSet(t *testing.T) {
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
				Annotations: map[string]api.Annotation{
					"cmd": map[string]interface{} {
						"k8s.annotation": "testing",
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
    "annotations": {
        "cmd": {
            "k8s.annotation": "testing"
        }
    }
}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.AnnotationSet("foo", "cmd", []string{"k8s.annotation=testing"})
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Creating Annotations... done

=== cmd Annotations
k8s.annotation      testing
`, "output")
}

func TestAnnotationUnset(t *testing.T) {
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
				Annotations: map[string]api.Annotation{
					"cmd": map[string]interface{}{
						"k8s.annotation": nil,
					},
				},
			}, r)
		}

		fmt.Fprintf(w, `{
	"owner": "jkirk",
	"app": "foo",
    "annotations": {
        "cmd": {
            "k8s.another/annotation": "another"
        }
    }
}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.AnnotationUnset("foo", "cmd", []string{"k8s.annotation"})
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Removing Annotations... done

=== cmd Annotations
k8s.another/annotation      another
`, "output")
}

