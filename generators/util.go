package generators

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
)

var (
	templates *rice.Box

	// ErrorSkip is a pseudo-error to indicate code generation is skipped
	ErrorSkip = errors.New("skipped generation")
)

// GenerationWork represents generation work
type GenerationWork struct {
	Waitgroup *sync.WaitGroup
	Done      chan<- GeneratedCode
}

// GeneratedCode represents code that has been generated and the intended file
type GeneratedCode struct {
	Generator string
	Filename  string
	Code      string
	Overwrite bool
	Error     error
}

// SetTemplates injects template box
func SetTemplates(box *rice.Box) {
	templates = box
}

//GetAbsPath returns absolute path
func GetAbsPath(src string) (string, error) {
	gopath := os.Getenv("GOPATH")
	location := strings.TrimRight(strings.Replace(src, "$GOPATH", gopath, -1), "\\/")

	if !path.IsAbs(location) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		location = path.Clean(wd + "/" + location)
	}

	if !FileExists(location) {
		return "", fmt.Errorf("file not found: %s", src)
	}

	return location, nil
}

// FileExists returns true is f (absolute path) exists; false otherwise
func FileExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// ExecuteTemplate applies templating a text/template template given data and returns the string output
func ExecuteTemplate(name string, data interface{}) (string, error) {
	var output bytes.Buffer

	raw, err := templates.String(name)

	if err != nil {
		return "", err
	}

	tpl, err := template.New(name).Parse(raw)
	if err != nil {
		return "", err
	}

	err = tpl.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
