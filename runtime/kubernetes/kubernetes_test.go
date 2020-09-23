// +build kubernetes

package kubernetes

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/micro/go-micro/v3/util/kubernetes/client"
)

func setupClient(t *testing.T) {
	files := []string{"token", "ca.crt"}
	for _, f := range files {
		cmd := exec.Command("kubectl", "get", "secrets", "-o",
			fmt.Sprintf(`jsonpath="{.items[?(@.metadata.annotations['kubernetes\.io/service-account\.name']=='micro-runtime')].data.%s}"`,
				strings.ReplaceAll(f, ".", "\\.")))
		if outp, err := cmd.Output(); err != nil {
			t.Fatalf("Failed to set k8s token %s", err)
		} else {
			outq := outp[1 : len(outp)-1]
			decoded, err := base64.StdEncoding.DecodeString(string(outq))
			if err != nil {
				t.Fatalf("Failed to set k8s token %s '%s'", err, outq)
			}
			if err := ioutil.WriteFile("/var/run/secrets/kubernetes.io/serviceaccount/"+f, decoded, 0755); err != nil {
				t.Fatalf("Error setting up k8s %s", err)
			}
		}

	}
	outp, err := exec.Command("kubectl", "config", "view", "-o", `jsonpath='{.clusters[?(@.name=="kind-kind")].cluster.server}'`).Output()
	if err != nil {
		t.Fatalf("Cannot find server for kind %s", err)
	}
	serverHost := string(outp)

	split := strings.Split(serverHost[9:len(serverHost)-1], ":")
	os.Setenv("KUBERNETES_SERVICE_HOST", split[0])
	os.Setenv("KUBERNETES_SERVICE_PORT", split[1])

}

func TestNamespaceCreateDelete(t *testing.T) {
	defer func() {
		exec.Command("kubectl", "delete", "namespace", "foobar").Run()
	}()
	setupClient(t)
	r := NewRuntime()
	if err := r.CreateNamespace("foobar"); err != nil {
		t.Fatalf("Unexpected error creating namespace %s", err)
	}

	if !namespaceExists(t, "foobar") {
		t.Fatalf("Namespace foobar not found")
	}
	if err := r.DeleteNamespace("foobar"); err != nil {
		t.Fatalf("Unexpected error deleting namespace %s", err)
	}
	if namespaceExists(t, "foobar") {
		t.Fatalf("Namespace foobar still exists")
	}
}

// This tests the generic CreateResource and DeleteResource methods with a networkpolicy:
func TestResourceCreateDelete(t *testing.T) {
	defer func() {
		exec.Command("kubectl", "-n", "baz", "delete", "networkpolicy", "ingress").Run()
		exec.Command("kubectl", "delete", "namespace", "baz").Run()
	}()
	setupClient(t)
	r := NewRuntime()
	if err := r.CreateResource(&client.Resource{
		Kind: "namespace",
		Value: client.Namespace{
			Metadata: &client.Metadata{
				Name: "baz",
			},
		},
	}); err != nil {
		t.Fatalf("Unexpected error creating namespace %s", err)
	}

	if !namespaceExists(t, "baz") {
		t.Fatalf("Namespace baz not found")
	}

	if err := r.CreateResource(&client.Resource{
		Kind: "networkpolicy",
		Name: "ingress",
		Value: client.NetworkPolicy{
			AllowedLabels: map[string]string{
				"owner": "test",
			},
			Metadata: &client.Metadata{
				Name:      "ingress",
				Namespace: "baz",
			},
		},
	}); err != nil {
		t.Fatalf("Unexpected error creating networkpolicy %s", err)
	}

	if !networkPolicyExists(t, "baz", "ingress") {
		t.Fatalf("NetworkPolicy baz.ingress not found")
	}

	if err := r.DeleteResource(&client.Resource{
		Kind: "networkpolicy",
		Name: "ingress",
		Value: client.NetworkPolicy{
			AllowedLabels: map[string]string{
				"owner": "test",
			},
			Metadata: &client.Metadata{
				Name:      "ingress",
				Namespace: "baz",
			},
		},
	}); err != nil {
		t.Fatalf("Unexpected error creating networkpolicy %s", err)
	}

	if networkPolicyExists(t, "baz", "ingress") {
		t.Fatalf("NetworkPolicy baz.ingress still exists")
	}

	if err := r.DeleteResource(&client.Resource{
		Kind: "namespace",
		Value: client.Namespace{
			Metadata: &client.Metadata{
				Name: "baz",
			},
		},
	}); err != nil {
		t.Fatalf("Unexpected error deleting namespace %s", err)
	}
	if namespaceExists(t, "baz") {
		t.Fatalf("Namespace baz still exists")
	}
}

func namespaceExists(t *testing.T, ns string) bool {
	cmd := exec.Command("kubectl", "get", "namespaces")
	outp, err := cmd.Output()
	if err != nil {
		t.Fatalf("Unexpected error listing namespaces %s", err)
	}
	exists, err := regexp.Match(ns+"\\s+Active", outp)
	if err != nil {
		t.Fatalf("Error listing namespaces %s", err)
	}
	return exists
}

func networkPolicyExists(t *testing.T, ns, np string) bool {
	cmd := exec.Command("kubectl", "-n", ns, "get", "networkpolicy")
	outp, err := cmd.Output()
	if err != nil {
		t.Fatalf("Unexpected error listing networkpolicies %s", err)
	}
	exists, err := regexp.Match(np, outp)
	if err != nil {
		t.Fatalf("Error listing networkpolicies %s", err)
	}
	return exists
}
