package main

import (
	"strings"
	"testing"
)

func TestManifestTestFunctionLinks(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, declaration string
		valid             bool
	}{
		{"TestOutput", "func TestOutput(t *testing.T) {}", true},
		{"Test", "func Test(*testing.T) {}", true},
		{"Test_case", "func Test_case(t *testing.T) {}", true},
		{"Test1", "func Test1(t *testing.T) {}", true},
		{"TestÉ", "func TestÉ(t *testing.T) {}", true},
		{"TestAlias", "func TestAlias(t *alias.T) {}", true},
		{"TestDot", "func TestDot(t *T) {}", true},
		{"TestMain", "func TestMain(t *testing.T) {}", true},
		{"TestMethodShadow", "type helper struct{}\nfunc (helper) TestMethodShadow(t *testing.T) {}\nfunc TestMethodShadow(t *testing.T) {}", true},
		{"Testhelper", "func Testhelper(t *testing.T) {}", false},
		{"Testé", "func Testé(t *testing.T) {}", false},
		{"TestMissingArgument", "func TestMissingArgument() {}", false},
		{"TestValueArgument", "func TestValueArgument(t testing.T) {}", false},
		{"TestWrongArgument", "func TestWrongArgument(t *testing.B) {}", false},
		{"TestMultipleArguments", "func TestMultipleArguments(t, u *testing.T) {}", false},
		{"TestResult", "func TestResult(t *testing.T) bool { return true }", false},
		{"TestGeneric", "func TestGeneric[P any](t *testing.T) {}", false},
		{"TestMethod", "type helper struct{}\nfunc (helper) TestMethod(t *testing.T) {}", false},
		{"TestMain", "func TestMain(m *testing.M) {}", false},
	} {
		for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
			t.Run(tt.name+"/"+mode, func(t *testing.T) {
				t.Parallel()
				root, m := testManifest(t)
				imports := "import \"testing\"\n"
				switch tt.name {
				case "TestAlias":
					imports = "import alias \"testing\"\n"
				case "TestDot":
					imports = "import . \"testing\"\n"
				}
				writeArtifact(t, root, "output_test.go", "package example\n"+imports+tt.declaration+"\n")
				m.Probes[0].Implementation.State = "verified"
				m.Probes[0].Implementation.Tests = []string{"output_test.go#" + tt.name}
				err := validate(root, m, mode, nil)
				if tt.valid && err != nil {
					t.Fatal(err)
				}
				if !tt.valid && (err == nil || !strings.Contains(err.Error(), "test link")) {
					t.Fatalf("error = %v, want invalid test link", err)
				}
			})
		}
	}
}
