package generators

import (
	"fmt"

	"github.com/fluxynet/gocipe/util"
)

// GenerateBootstrap returns bootstrap generated code
func GenerateBootstrap(work util.GenerationWork, opts util.BootstrapOpts) error {
	if !opts.Generate {
		util.DeleteIfExists("app/bootstrap.gocipe.go")
		work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Error: util.ErrorSkip}
	}

	work.Waitgroup.Add(1)

	if opts.HTTPPort == "" {
		opts.HTTPPort = "7000"
	}

	if opts.GRPCPort == "" {
		opts.GRPCPort = "4000"
	}

	code, err := util.ExecuteTemplate("bootstrap.go.tmpl", struct {
		Bootstrap util.BootstrapOpts
	}{opts})

	if err != nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Error: fmt.Errorf("failed to execute template: %s", err)}
		return err
	}

	env, err := util.ExecuteTemplate("bootstrap_env.tmpl", struct {
		Bootstrap util.BootstrapOpts
	}{opts})

	if err != nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Error: fmt.Errorf("failed to execute template: %s", err)}
		return err
	}

	work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Filename: "app/bootstrap.gocipe.go", Code: code}
	work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Filename: ".env.dist", Code: env, GeneratedHeaderFormat: "# %s"}
	return nil
}
