package cache

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func cacheTailwind(cacheDir string) (bin string, err error) {
	bin = filepath.Join(cacheDir, "bin/tailwindcss")

	twOS := runtime.GOOS
	if twOS == "darwin" {
		twOS = "macos"
	}

	twURL := fmt.Sprintf(
		"https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-%s-%s",
		twOS,
		runtime.GOARCH,
	)

	_, err = os.Stat(bin)
	if err != nil {
		// we do not have the binary; download

		res, err := http.Get(twURL)
		if err != nil {
			return "", err
		}

		f, err := os.OpenFile(bin, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return "", err
		}
		defer f.Close()

		_, err = io.Copy(f, res.Body)
		if err != nil {
			return "", err
		}
	}

	return bin, nil
}
