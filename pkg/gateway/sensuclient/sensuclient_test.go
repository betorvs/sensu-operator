package sensuclient

import (
	"strings"
	"testing"
)

const (
	fakeSensuURL = "http://sensu-api.svc.cluster.local:8080"
)

func TestSensuURLGenerator(t *testing.T) {
	test1a := "checks"
	test1b := "development"
	test1c := "test"
	result1 := sensuURLGenerator(fakeSensuURL, test1a, test1b, test1c)
	if !strings.Contains(result1, "namespaces/development/checks") {
		t.Fatalf("Invalid 1.1  TestSensuURLGenerator %s", result1)
	}

	test2a := "namespaces"
	test2b := "default"
	test2c := "all"
	result2 := sensuURLGenerator(fakeSensuURL, test2a, test2b, test2c)
	if !strings.Contains(result2, "v2/namespaces") {
		t.Fatalf("Invalid 1.2  TestSensuURLGenerator %s", result2)
	}
}

func TestSensuPostURLGenerator(t *testing.T) {
	test1a := "checks"
	test1b := "development"
	result1 := sensuPostURLGenerator(fakeSensuURL, test1a, test1b)
	if !strings.Contains(result1, "namespaces/development/checks") {
		t.Fatalf("Invalid 2.1  TestSensuPostURLGenerator %s", result1)
	}

	test2a := "namespaces"
	test2b := "default"
	result2 := sensuPostURLGenerator(fakeSensuURL, test2a, test2b)
	if !strings.Contains(result2, "v2/namespaces") {
		t.Fatalf("Invalid 2.2  TestSensuPostURLGenerator %s", result2)
	}
}
