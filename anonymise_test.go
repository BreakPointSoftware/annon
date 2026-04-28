package annon

import (
	"bytes"
	"strings"
	"testing"
)

type customer struct {
	ID     string `json:"id" yaml:"id"`
	Email  string `json:"email" yaml:"email"`
	Secret string `json:"secret" yaml:"secret" anonymise:"remove"`
	Note   string `json:"note" yaml:"note" anonymise:"false"`
}

func TestCopyJSONYAMLAndBlobs(t *testing.T) {
	a, err := New()
	if err != nil { t.Fatal(err) }
	input := customer{ID: "123", Email: "greg@example.com", Secret: "secret", Note: "greg@example.com"}

	copyAny, err := a.Copy(input)
	if err != nil { t.Fatal(err) }
	copyValue := copyAny.(customer)
	if copyValue.Email != "g***@example.com" || copyValue.Secret != "" || copyValue.Note != "greg@example.com" {
		t.Fatalf("unexpected copy result: %+v", copyValue)
	}

	jsonBlob, err := a.JSON(input)
	if err != nil { t.Fatal(err) }
	if !bytes.Contains(jsonBlob, []byte(`"email":"g***@example.com"`)) || bytes.Contains(jsonBlob, []byte(`"secret"`)) {
		t.Fatalf("unexpected json output: %s", jsonBlob)
	}

	yamlBlob, err := a.YAML(input)
	if err != nil { t.Fatal(err) }
	if !strings.Contains(string(yamlBlob), "email: g***@example.com") || strings.Contains(string(yamlBlob), "secret:") {
		t.Fatalf("unexpected yaml output: %s", yamlBlob)
	}

	rawJSON := []byte(`{"email":"greg@example.com","vehicle":{"reg":"AB12 CDE"}}`)
	safeJSON, err := a.FromJSON(rawJSON)
	if err != nil { t.Fatal(err) }
	if !bytes.Contains(safeJSON, []byte(`"reg":"AB12 ***"`)) { t.Fatalf("unexpected raw json output: %s", safeJSON) }

	rawYAML := []byte("email: greg@example.com\nvehicle:\n  reg: AB12 CDE\n")
	safeYAML, err := a.FromYAML(rawYAML)
	if err != nil { t.Fatal(err) }
	if !strings.Contains(string(safeYAML), "reg: AB12 ***") { t.Fatalf("unexpected raw yaml output: %s", safeYAML) }
}

func TestFieldRulesAndValueDetection(t *testing.T) {
	type alias struct { Alias string `json:"customerAlias"` }
	jsonBlob, err := JSON(alias{Alias: "secret"}, WithFieldRules(StrongRule(RedactStrategy, "customerAlias")))
	if err != nil { t.Fatal(err) }
	if !bytes.Contains(jsonBlob, []byte(`"customerAlias":"[REDACTED]"`)) {
		t.Fatalf("unexpected field-rule json: %s", jsonBlob)
	}

	safe, err := FromJSON([]byte(`{"note":"greg@example.com"}`), WithFieldDetection(false), WithValueDetection(true))
	if err != nil { t.Fatal(err) }
	if !bytes.Contains(safe, []byte(`"note":"g***@example.com"`)) {
		t.Fatalf("unexpected value-detection output: %s", safe)
	}
}
