package output

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/fluxynet/gocipe/util"
)

// toolset represents go tools used by the generators
type toolset struct {
	GoImports string
	GoFmt     string
	Protoc    string
	Dep       string
}

const (
	logWritten = "✅ [Wrote]"
	logSkipped = "💤 [Skipped]"
	logError   = "❌ [Error]"
)

var (
	_recipePath                 string
	_log, _gofiles              []string
	_tools                      toolset
	_written, _skipped, _failed int
	_verbose                    bool
)

func init() {
	initToolset()
}

// Inject gets path injected into this package
func Inject(path string) {
	_recipePath = path
}

// SetVerbose can be used to toggle verbosity
func SetVerbose(verbose bool) {
	_verbose = verbose
}

// Log outputs to log file
func Log(message string, a ...interface{}) {
	_log = append(_log, fmt.Sprintf(message, a...))
}

// AddGoFile adds a file in the queue to be gofmt'ed and goimport'ed
func AddGoFile(filename string) {
	filename, _ = util.GetAbsPath(filename)
	_gofiles = append(_gofiles, filename)
}

// WriteLog write logs in log-file
func WriteLog() {
	err := ioutil.WriteFile(_recipePath+".log", []byte(strings.Join(_log, "\n")), os.FileMode(0755))
	if err != nil {
		fmt.Printf("failed to write file log file %s.log: %s", _recipePath, err)
		return
	}

	if _verbose {
		fmt.Println(strings.Join(_log, "\n"))
	}

	_log = []string{}
}

// Process listens to generators and processes generated files
func Process(waitgroup *sync.WaitGroup, work util.GenerationWork, noSkip bool) {

	var (
		output []string
	)

	aggregates := make(map[string][]util.GeneratedCode)

	for generated := range work.Done {
		if generated.Error == util.ErrorSkip {
			Log(logSkipped+" Generation skipped [%s]", generated.Generator)
			_skipped++
		} else if generated.Error != nil {
			Log(logError+" Generation failed [%s]: %s", generated.Generator, generated.Error)
			_failed++
		} else if generated.Aggregate {
			a := aggregates[generated.Filename]
			aggregates[generated.Filename] = append(a, generated)
		} else {
			fname, l, err := saveGenerated(generated, noSkip)
			Log(l)

			if err == nil {
				if strings.HasSuffix(fname, ".go") {
					AddGoFile(fname)
				} else if strings.HasSuffix(fname, ".sh") {
					os.Chmod(fname, 0755)
				}

				_written++
			} else if err == util.ErrorSkip {
				_skipped++
			} else {
				_failed++
			}
		}

		work.Waitgroup.Done()
	}

	for _, generated := range aggregates {
		fname, l, err := saveAggregate(generated, noSkip)
		Log(l)

		if err == nil {
			if strings.HasSuffix(fname, ".go") {
				AddGoFile(fname)
			}

			_written++
		} else if err == util.ErrorSkip {
			_skipped++
		} else {
			_failed++
		}
	}

	// WriteLog()

	if _skipped > 0 {
		output = append(output, fmt.Sprintf("Skipped %d files.", _skipped))
	}

	if _written > 0 {
		output = append(output, fmt.Sprintf("Wrote %d files.", _written))
	}

	if _failed > 0 {
		output = append(output, fmt.Sprintf("%d errors occurred during recipe generation.", _failed))
	}

	output = append(output, fmt.Sprintf("See log file %s.log for details.", _recipePath))
	fmt.Println(strings.Join(output, " "))
	waitgroup.Done()
}

// saveGenerated saves a generated file and returns absolute filename, log entry and error
func saveGenerated(generated util.GeneratedCode, noSkip bool) (string, string, error) {
	filename, err := util.GetAbsPath(generated.Filename)
	if err != nil {
		return "", fmt.Sprintf(logError+" cannot resolve path [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	if !noSkip && util.FileExists(filename) && generated.NoOverwrite {
		return "", fmt.Sprintf(logSkipped+" skipping existing file [%s] %s", generated.Generator, generated.Filename), util.ErrorSkip
	}

	var mode os.FileMode = 0755
	if err = os.MkdirAll(path.Dir(filename), mode); err != nil {
		return "", fmt.Sprintf(logError+" directory creation failed [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	var code []byte
	if generated.NoOverwrite || generated.GeneratedHeaderFormat == util.NoHeaderFormat {
		code = []byte(generated.Code)
	} else {
		var generatedHeaderFormat string
		if generated.GeneratedHeaderFormat == "" {
			generatedHeaderFormat = "// %s"
		} else {
			generatedHeaderFormat = generated.GeneratedHeaderFormat
		}

		generatedHeaderFormat = fmt.Sprintf(generatedHeaderFormat, `generated by gocipe; DO NOT EDIT`)

		code = []byte(generatedHeaderFormat + "\n\n" + generated.Code)
	}

	err = ioutil.WriteFile(filename, code, mode)
	if err != nil {
		return "", fmt.Sprintf(logError+" failed to write file [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	return filename, fmt.Sprintf(logWritten+" %s", filename), nil
}

// GenerateAndSave saves a generated file and returns error
func GenerateAndSave(component string, template string, filename string, data interface{}, noOverwrite bool) error {
	var (
		code     string
		err      error
		isString bool
		mode     os.FileMode = 0755
	)

	filename, err = util.GetAbsPath(filename)
	if err != nil {
		Log(logError+" Generate (%s) %s failed: %s", component, filename, err)
		_failed++
		return err
	}

	if noOverwrite && util.FileExists(filename) {
		Log(logSkipped+" %s (%s)", filename, component)
		_skipped++
		return util.ErrorSkip
	}

	if err = os.MkdirAll(path.Dir(filename), mode); err != nil {
		Log(logError+" Generate (%s) %s failed: %s", component, filename, err)
		_failed++
		return err
	}

	if code, isString = data.(string); !isString {
		code, err = util.ExecuteTemplate(template, data)
		if err != nil {
			Log(logError+" Generate (%s) %s failed: %s", component, filename, err)
			_failed++
			return err
		}
	}

	err = ioutil.WriteFile(filename, []byte(code), mode)
	if err != nil {
		Log(logError+" Generate (%s) %s failed: %s", component, filename, err)
		_failed++
		return err
	}

	Log(logWritten+" %s", filename)
	_written++

	if strings.HasSuffix(filename, ".go") {
		AddGoFile(filename)
	} else if strings.HasSuffix(filename, ".sh") {
		os.Chmod(filename, 0755)
	}

	return nil
}

// saveAggregate saves aggregated files and returns absolute filename, log entry and Error
func saveAggregate(aggregate []util.GeneratedCode, noSkip bool) (string, string, error) {
	var generated util.GeneratedCode

	generated.Filename = aggregate[0].Filename
	generated.Generator = aggregate[0].Generator
	generated.GeneratedHeaderFormat = aggregate[0].GeneratedHeaderFormat

	for _, g := range aggregate {
		generated.NoOverwrite = generated.NoOverwrite || g.NoOverwrite
		generated.Code += g.Code + "\n"
	}

	return saveGenerated(generated, noSkip)
}

// initToolset check if all required tools are present
func initToolset() {
	var (
		err error
		ok  = true
	)

	_tools.GoImports, err = exec.LookPath("goimports")
	if err != nil {
		fmt.Println("Required tool goimports not found: ", err)
		ok = false
	}

	_tools.GoFmt, err = exec.LookPath("gofmt")
	if err != nil {
		fmt.Println("Required tool gofmt not found: ", err)
		ok = false
	}

	_tools.Protoc, err = exec.LookPath("protoc")
	if err != nil {
		fmt.Println("Required tool protoc not found: ", err)
		ok = false
	}

	_, err = exec.LookPath("protoc-gen-go")
	if err != nil {
		fmt.Println("Required tool protoc-gen-go not found: ", err)
		fmt.Println("Install using go get -u github.com/golang/protobuf/protoc-gen-go")
		ok = false
	}

	_tools.Dep, err = exec.LookPath("dep")
	if err != nil {
		fmt.Println("Required tool dep not found: ", err)
		fmt.Println("Install using go get -u github.com/golang/dep/cmd/dep")
		ok = false
	}

	if !ok {
		log.Fatalln("Please install above tools before continuing.")
	}
}

// PostProcessGoFiles executes goimports and gofmt on go files that have been generated
func PostProcessGoFiles() {
	if len(_gofiles) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(_gofiles))

	for _, file := range _gofiles {
		go func(file string) {
			defer wg.Done()

			cmd := exec.Command(_tools.GoImports, "-w", file)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				fmt.Printf("Error running %s on %s: %s\n", _tools.GoImports, file, err)
				return
			}

			cmd = exec.Command(_tools.GoFmt, "-w", file)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()

			if err != nil {
				fmt.Printf("Error running %s on %s: %s\n", _tools.GoFmt, file, err)
			}
		}(file)
	}

	wg.Wait()

	var mode string
	if util.FileExists(util.WorkingDir + "/Gopkg.toml") {
		mode = "ensure"
	} else {
		mode = "init"
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		cmd := exec.Command(_tools.Dep, mode)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			fmt.Printf("Error running dep %s: %s\n", mode, err)
		}
	}()

	fmt.Printf("dep %s in progress...", mode)
	wg.Wait()
}

// ProcessProto executes protoc to generate go files from protobuf files
func ProcessProto() {
	var (
		cmd    *exec.Cmd
		err    error
		mode   os.FileMode = 0755
		gopath             = os.Getenv("GOPATH") + "/src/"
	)

	Log("[Protobuf] Executing protoc to generate go files...")

	// models.proto
	if !util.FileExists(util.WorkingDir + "/models") {
		if err = os.MkdirAll(util.WorkingDir+"/models", mode); err != nil {
			Log("✗ could not create folder: %s", util.WorkingDir+"/models")
			fmt.Printf("Error creating folder %s: %s\n", util.WorkingDir+"/models", err)
			return
		}

		Log("✓ created folder: %s", util.WorkingDir+"/models")
	}
	cmd = exec.Command(
		_tools.Protoc,
		`-I=`+util.WorkingDir+`/proto`,
		util.WorkingDir+`/proto/models.proto`,
		`--go_out=plugins=grpc:`+gopath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		Log("✗ protoc execution error (%s): %s", "models.proto", err)
		fmt.Printf("Error running %s: %s\n", _tools.Protoc, err)
		return
	}

	Log("✓ protoc generated go files from: %s", "models.proto")

	// service_bread.proto, if bread service is to be generated
	if util.FileExists(util.WorkingDir + `/proto/service_bread.proto`) {
		if !util.FileExists(util.WorkingDir + "/services/bread") {
			if err = os.MkdirAll(util.WorkingDir+"/services/bread", mode); err != nil {
				Log("✗ could not create folder: %s", util.WorkingDir+"/services/bread")
				fmt.Printf("Error creating folder %s: %s\n", util.WorkingDir+"/services/bread", err)
				return
			}

			Log("✓ created folder: %s", util.WorkingDir+"/services/bread")
		}
		cmd = exec.Command(
			_tools.Protoc,
			`-I=`+util.WorkingDir+`/proto`,
			util.WorkingDir+`/proto/service_bread.proto`,
			`--go_out=plugins=grpc:`+gopath,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			Log("✗ protoc execution error (%s): %s", "service_bread.proto", err)
			fmt.Printf("Error running %s: %s\n", _tools.Protoc, err)
			return
		}

		Log("✓ protoc generated go files from: %s", "service_bread.proto")
	}

	// cmd = exec.Command(
	// 	toolset.Protoc,
	// 	`-I=proto`,
	// 	`--plugin="protoc-gen-ts=`+util.WorkingDir+`/web/node_modules/.bin/protoc-gen-ts"`,
	// 	`--js_out="binary:`+util.WorkingDir+`/web/src/services"`,
	// 	`--ts_out="`+util.WorkingDir+`/web/src/services"`,
	// 	`proto/models.proto`,
	// )

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Dir = util.WorkingDir
	// err = cmd.Run()

	// if err != nil {
	// 	fmt.Printf("Error running %s: %s\n", toolset.Protoc, err)
	// 	return
	// }
}
