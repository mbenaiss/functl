package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mbenaiss/functl/config"
	"github.com/mbenaiss/functl/provider"
	"github.com/mbenaiss/functl/provider/gcp"
)

const usageText = `functl [cmd] <options>
Commands:
  - new Create a new api config file
  - deploy  Creates a new lambda function from a zip.
  - list  Updates the existing lambda function with a zip.
  - delete  Updates the existing lambda function with a zip.
`

func main() {

	if len(os.Args) < 2 {
		printUsage(1)
	}
	flag.Usage = func() {
		printUsage(1)
	}

	switch os.Args[1] {
	case "new":
		err := newConfigFile()
		if err != nil {
			panic(err)
		}
	case "deploy":
		c, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}
		currentProvider, err := getCurrentProvider(os.Args[2])
		if err != nil {
			panic(err)
		}
		currentProvider.Deploy("test", c)
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
	err := config.ValidProvider(p)
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

func newConfigFile() error {
	tem, err := os.OpenFile("./template/template.yaml", os.O_RDONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open template file %w", err)
	}
	dst, err := os.Create("api.yaml")
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
