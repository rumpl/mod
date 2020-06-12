package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/github/hub/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	gitignore = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

cov

# Test binary, built with "go test -c\"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out
# The binary
%s
`
	readme = `# %s
`
	makefile = `all: cli

cli:
	@go build .

.PHONY: cli
`
	mainfile = `package main

import (
	"fmt"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fmt.Println("Hello world")

	return nil
}
`
)

func main() {
	app := cli.App{
		Usage: "mod",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "init",
				Action: func(ctx *cli.Context) error {
					repo := ctx.Args().First()
					if repo == "" {
						return errors.New("please provide a repo name")
					}
					return run(repo)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(repo string) error {
	if _, err := os.Stat(repo); !os.IsNotExist(err) {
		return errors.New("directory already exists")
	}

	cc := github.CurrentConfig()
	user := cc.PromptForUser("github.com")

	if user == "" {
		host, err := cc.PromptForHost("github.com")
		if err != nil {
			return err
		}
		user = host.User
	}

	inits := []initFunc{
		mkdir,
		chdir,
		modInit,
		gitInit,
	}

	for _, init := range inits {
		if err := init(user, repo); err != nil {
			return err
		}
	}

	if err := addFiles(repo); err != nil {
		return err
	}

	fmt.Println("All done, happy hacking!")

	return nil
}

type initFunc func(string, string) error

func mkdir(user string, repo string) error {
	return os.Mkdir(repo, 0755)
}

func chdir(user string, repo string) error {
	return os.Chdir(repo)
}

func modInit(user string, repo string) error {
	return execute("go", "mod", "init", fmt.Sprintf("github.com/%s/%s", user, repo))
}

func gitInit(_ string, _ string) error {
	return execute("git", "init")
}

func execute(args ...string) error {
	logrus.Debugf("Executing %s", strings.Join(args, " "))
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return errors.Wrap(err, string(out))
	}
	return nil
}

func addFiles(repo string) error {
	actions := []action{
		writeGitIgnore,
		writeMainGo,
		writeMakefile,
		writeReadme,
	}

	for _, action := range actions {
		if err := action(repo); err != nil {
			return err
		}
	}

	return nil
}

type action func(string) error

func writeGitIgnore(repo string) error {
	return writeFile(".gitignore", []byte(fmt.Sprintf(gitignore, repo)))
}

func writeMainGo(_ string) error {
	return writeFile("main.go", []byte(mainfile))
}

func writeMakefile(_ string) error {
	return writeFile("Makefile", []byte(makefile))
}

func writeReadme(repo string) error {
	return writeFile("README.md", []byte(fmt.Sprintf(readme, repo)))
}

func writeFile(file string, contents []byte) error {
	logrus.Debugf("Writing file %s", file)
	return ioutil.WriteFile(file, contents, 0644)
}
