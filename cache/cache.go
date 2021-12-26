package cache

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

var cmdLookup = map[string]func(cacheDir string) (bin string, err error){
	"tailwindcss": cacheTailwind,
}

var usage = errors.New("Usage: go run blake.io/cache@latest run command")

func Main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.SetPrefix("cache: ")
	log.SetFlags(0)

	flag.Parse()

	if flag.NArg() < 2 || flag.Arg(0) != "run" {
		return usage
	}

	cmd := flag.Arg(1)
	f := cmdLookup[cmd]
	if f == nil {
		return fmt.Errorf("cache: unknown command %q; currently only tailwindcss is supported", cmd)
	}

	cacheDir, err := filepath.Abs(".blake.io.cache")
	if err != nil {
		return err
	}
	cacheDir = filepath.Join(cacheDir, cmd)

	binDir := filepath.Join(cacheDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	bin, err := f(cacheDir)
	if err != nil {
		return err
	}

	args := flag.Args()[2:]
	// exec out vs wrap with exec.Cmd so we become tailwindcss
	return syscall.Exec(
		bin,
		append([]string{bin}, args...),
		os.Environ(),
	)
}
