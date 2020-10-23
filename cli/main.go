package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/rakyll/statik/fs"

	"github.com/mbenaiss/functl/api"
	"github.com/mbenaiss/functl/cli/provider"
	"github.com/mbenaiss/functl/cli/provider/gcp"
	_ "github.com/mbenaiss/functl/cli/statik"
)

const usageText = `functl [cmd] <options>
Commands:
  - new Create a new api config file
  - deploy  Creates a new lambda function from a zip.
  - list  Updates the existing lambda function with a zip.
  - delete  Updates the existing lambda function with a zip.
`

var (
	fnFiles = map[string]string{
		"/api.go":    "api.go",
		"/config.go": "config.go",
		"/go.mod":    "go.mod",
		"/go.sum":    "go.sum",
	}
)

func main() {
	if len(os.Args) < 2 {
		printUsage(1)
	}
	flag.Usage = func() {
		printUsage(1)
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "new":
		err := newConfigFile(statikFS)
		if err != nil {
			panic(err)
		}
	case "deploy":
		c, err := api.LoadConfig("./api.yaml")
		if err != nil {
			panic(err)
		}
		currentProvider, err := getCurrentProvider(os.Args[2])
		if err != nil {
			panic(err)
		}
		err = copyFnFile(statikFS)
		if err != nil {
			panic(err)
		}
		_, err = currentProvider.Deploy("Handler", c)
		if err != nil {
			panic(err)
		}
		err = deleteFnFile()
		if err != nil {
			panic(err)
		}
	case "list":
		currentProvider, err := getCurrentProvider(os.Args[2])
		if err != nil {
			panic(err)
		}
		currentProvider.List()
	case "delete":
		currentProvider, err := getCurrentProvider(os.Args[2])
		if err != nil {
			panic(err)
		}
		currentProvider.Delete("")
	default:
		printUsage(1)
	}
}

func getCurrentProvider(p string) (provider.Provider, error) {
	err := api.ValidProvider(p)
	if err != nil {
		return nil, err
	}
	switch p {
	case "aws":
	case "gcp":
		return gcp.New(), nil
	case "kubernetes":
	case "vercel":
	}
	return nil, nil
}

func printUsage(code int) {
	fmt.Print(usageText)
	os.Exit(code)
}

func copyFile(f http.FileSystem, srcName string, dstName string) error {
	tem, err := f.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open template file %w", err)
	}
	dst, err := os.Create(dstName)
	if err != nil {
		return fmt.Errorf("failed to create config file %w", err)
	}
	_, err = io.Copy(dst, tem)
	if err != nil {
		return fmt.Errorf("failed to add config file %w", err)
	}
	err = tem.Close()
	if err != nil {
		return fmt.Errorf("failed to close template file %w", err)
	}
	err = dst.Close()
	if err != nil {
		return fmt.Errorf("failed to close config file %w", err)
	}
	return nil
}

func newConfigFile(f http.FileSystem) error {
	return copyFile(f, "/template.yaml", "api.yaml")
}

func copyFnFile(f http.FileSystem) error {
	wg := &sync.WaitGroup{}
	for src, dst := range fnFiles {
		wg.Add(1)
		go func(s, d string) {
			defer wg.Done()
			err := copyFile(f, s, d)
			if err != nil {
				log.Fatalf("an error occurred while copying file %+v", err)
			}
		}(src, dst)
	}
	wg.Wait()
	return nil
}

func deleteFnFile() error {
	wg := &sync.WaitGroup{}
	for _, dst := range fnFiles {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			err := os.Remove(d)
			if err != nil {
				log.Fatalf("an error occurred while deleting file %+v", err)
			}
		}(dst)
	}
	wg.Wait()
	return nil
}
