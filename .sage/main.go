package main

import (
	"context"
	"fmt"
	"os"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sgbuf"
	"go.einride.tech/sage/tools/sgconvco"
	"go.einride.tech/sage/tools/sggit"
	"go.einride.tech/sage/tools/sggo"
	"go.einride.tech/sage/tools/sggolangcilint"
	"go.einride.tech/sage/tools/sgmdformat"
	"go.einride.tech/sage/tools/sgsqlc"
	"go.einride.tech/sage/tools/sgyamlfmt"
)

func main() {
	sg.GenerateMakefiles(
		sg.Makefile{
			Path:          sg.FromGitRoot("Makefile"),
			DefaultTarget: Default,
		},
	)
}

func Default(ctx context.Context) error {
	sg.Deps(ctx, ConvcoCheck, FormatMarkdown, FormatYaml)
	sg.Deps(ctx, GoLint)
	sg.Deps(ctx, GoTest)
	sg.Deps(ctx, GoModTidy)
	sg.Deps(ctx, GitVerifyNoDiff)
	return nil
}

func GenerateSQL(ctx context.Context) error {
	sg.Logger(ctx).Println("generating SQL files...")
	cmd := sgsqlc.Command(ctx, "generate")
	cmd.Dir = sg.FromGitRoot("internal/server/db")
	return cmd.Run()
}

func GoModTidy(ctx context.Context) error {
	sg.Logger(ctx).Println("tidying Go module files...")
	return sg.Command(ctx, "go", "mod", "tidy", "-v").Run()
}

func GoTest(ctx context.Context) error {
	sg.Logger(ctx).Println("running Go tests...")
	return sggo.TestCommand(ctx).Run()
}

func GoLint(ctx context.Context) error {
	sg.Logger(ctx).Println("linting Go files...")
	return sggolangcilint.Run(ctx)
}

func FormatMarkdown(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting Markdown files...")
	return sgmdformat.Command(ctx).Run()
}

func FormatYaml(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting Yaml files...")
	return sgyamlfmt.Run(ctx)
}

func ConvcoCheck(ctx context.Context) error {
	sg.Logger(ctx).Println("checking git commits...")
	return sgconvco.Command(ctx, "check", "origin/master..HEAD").Run()
}

func GitVerifyNoDiff(ctx context.Context) error {
	sg.Logger(ctx).Println("verifying that git has no diff...")
	return sggit.VerifyNoDiff(ctx)
}

const protoFolder = "proto"

func BufLint(ctx context.Context) error {
	sg.Logger(ctx).Println("linting protobuf files...")
	cmd := sgbuf.Command(ctx, "lint")
	cmd.Dir = sg.FromGitRoot(protoFolder)
	return cmd.Run()
}

func BufGenerate(ctx context.Context) error {
	sg.Logger(ctx).Println("generating protobuf files...")
	cmd := sgbuf.Command(ctx, "generate")
	cmd.Dir = sg.FromGitRoot(protoFolder)
	return cmd.Run()
}

func BufModUpdate(ctx context.Context) error {
	sg.Logger(ctx).Println("updating buf modules...")
	cmd := sgbuf.Command(ctx, "mod", "update")
	cmd.Dir = sg.FromGitRoot(protoFolder)
	return cmd.Run()
}

func Proto(ctx context.Context) error {
	sg.Deps(ctx, BufModUpdate, BufLint)
	sg.Deps(ctx, BufGenerate)
	return nil
}

func dirty(ctx context.Context) bool {
	return sg.Output(
		sggit.Command(ctx, "status", "--porcelain"),
	) != ""
}

func lddFlags(ctx context.Context) string {
	ver := sggit.ShortSHA(ctx)
	tags := sggit.Tags(ctx)
	if len(tags) > 0 {
		ver = tags[0]
		if dirty(ctx) {
			ver += "-dirty"
		}
	}
	return fmt.Sprintf("-ldflags=-X 'main.LDDVersion=%s'", ver)
}

func Build(ctx context.Context) error {
	sg.Logger(ctx).Println("building...")
	return sg.Command(ctx, "go", "build", lddFlags(ctx), "./cmd/profzf").Run()
}

func Install(ctx context.Context) error {
	sg.Logger(ctx).Printf("installing profzf to %s/bin...", os.Getenv("GOPATH"))
	return sg.Command(ctx, "go", "install", lddFlags(ctx), "./cmd/profzf").Run()
}
