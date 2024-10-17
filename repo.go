package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lesiw.io/defers"
	"lesiw.io/flag"
	repo "lesiw.io/repo/lib"
)

var repodir = "."
var errParse = errors.New("parse error")

var (
	flags   = flag.NewSet(os.Stderr, "repo URL")
	version = flags.Bool("V,version", "print version and exit")
	force   = flags.Bool("f,force", "delete and re-clone repository")
)

//go:embed version.txt
var versionfile string

func main() {
	if err := run(); err != nil {
		if !errors.Is(err, errParse) {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		defers.Exit(1)
	}
	defers.Exit(0)
}

func run() error {
	defers.Add(func() { fmt.Println(repodir) })

	if err := flags.Parse(os.Args[1:]...); err != nil {
		return errParse
	}
	if *version {
		return errors.New(strings.TrimSpace(versionfile))
	}
	if len(flags.Args) == 0 {
		flags.PrintError("no URL given")
		return errParse
	}

	url := flags.Args[0]
	prefix := os.Getenv("REPOPREFIX")
	if prefix == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not get home directory: %s", err)
		}
		prefix = filepath.Join(home, ".local", "src")
	}
	var err error
	if repodir, err = repo.Clone(prefix, url, *force); err != nil {
		return err
	}
	return nil
}
