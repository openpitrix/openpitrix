// Copyright confd. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE-confd file.

// +build ignore

package libconfd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"text/template"
)

// createTempDirs is a helper function which creates temporary directories
// required by confd. createTempDirs returns the path name representing the
// confd confDir.
// It returns an error if any.
func createTempDirs() (string, error) {
	confDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	err = os.Mkdir(filepath.Join(confDir, "templates"), 0755)
	if err != nil {
		return "", err
	}
	err = os.Mkdir(filepath.Join(confDir, "conf.d"), 0755)
	if err != nil {
		return "", err
	}
	return confDir, nil
}

var templateResourceConfigTmpl = `
[template]
src = "{{.src}}"
dest = "{{.dest}}"
keys = [
  "foo",
]
`

func TestProcessTemplateResources(t *testing.T) {
	// Setup temporary conf, config, and template directories.
	tempConfDir, err := createTempDirs()
	if err != nil {
		t.Errorf("Failed to create temp dirs: %s", err)
	}
	defer os.RemoveAll(tempConfDir)

	// Create the src template.
	srcTemplateFile := filepath.Join(tempConfDir, "templates", "foo.tmpl")
	err = ioutil.WriteFile(srcTemplateFile, []byte(`foo = {{getv "/foo"}}`), 0644)
	if err != nil {
		t.Error(err)
	}

	// Create the dest.
	destFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Errorf("Failed to create destFile: %v", err)
	}
	defer os.Remove(destFile.Name())

	// Create the template resource configuration file.
	templateResourcePath := filepath.Join(tempConfDir, "conf.d", "foo.toml")
	templateResourceFile, err := os.Create(templateResourcePath)
	if err != nil {
		t.Errorf("%v", err)
	}
	tmpl, err := template.New("templateResourceConfig").Parse(templateResourceConfigTmpl)
	if err != nil {
		t.Errorf("Unable to parse template resource template: %v", err)
	}
	data := make(map[string]string)
	data["src"] = "foo.tmpl"
	data["dest"] = destFile.Name()
	err = tmpl.Execute(templateResourceFile, data)
	if err != nil {
		t.Errorf("%v", err)
	}

	os.Setenv("FOO", "bar")
	storeClient, err := NewEnvBackendClient()
	if err != nil {
		t.Errorf("%v", err)
	}
	c := &Config{
		ConfDir:     tempConfDir,
		ConfigDir:   filepath.Join(tempConfDir, "conf.d"),
		TemplateDir: filepath.Join(tempConfDir, "templates"),
	}
	// Process the test template resource.
	err = NewOnetimeProcessor(c).Process(storeClient)
	if err != nil {
		t.Error(err)
	}
	// Verify the results.
	expected := "foo = bar"
	results, err := ioutil.ReadFile(destFile.Name())
	if err != nil {
		t.Error(err)
	}
	if string(results) != expected {
		t.Errorf("Expected contents of dest == '%s', got %s", expected, string(results))
	}
}

func TestSameConfigTrue(t *testing.T) {
	var tr TemplateResourceProcessor

	src, err := ioutil.TempFile("", "src")
	defer os.Remove(src.Name())
	if err != nil {
		t.Errorf("%v", err)
	}
	_, err = src.WriteString("foo")
	if err != nil {
		t.Errorf("%v", err)
	}
	dest, err := ioutil.TempFile("", "dest")
	defer os.Remove(dest.Name())
	if err != nil {
		t.Errorf("%v", err)
	}
	_, err = dest.WriteString("foo")
	if err != nil {
		t.Errorf("%v", err)
	}
	status, err := tr.checkSameConfig(src.Name(), dest.Name())
	if err != nil {
		t.Errorf("%v", err)
	}
	if status != true {
		t.Errorf("Expected sameConfig(src, dest) to be %v, got %v", true, status)
	}
}

func TestSameConfigFalse(t *testing.T) {
	var tr TemplateResourceProcessor

	src, err := ioutil.TempFile("", "src")
	defer os.Remove(src.Name())
	if err != nil {
		t.Errorf("%v", err)
	}
	_, err = src.WriteString("src")
	if err != nil {
		t.Errorf("%v", err)
	}
	dest, err := ioutil.TempFile("", "dest")
	defer os.Remove(dest.Name())
	if err != nil {
		t.Errorf("%v", err)
	}
	_, err = dest.WriteString("dest")
	if err != nil {
		t.Errorf("%v", err)
	}
	status, err := tr.checkSameConfig(src.Name(), dest.Name())
	if err != nil {
		t.Errorf("%v", err)
	}
	if status != false {
		t.Errorf("Expected sameConfig(src, dest) to be %v, got %v", false, status)
	}
}
