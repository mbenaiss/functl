package provider

import "github.com/mbenaiss/functl/config"

type Provider interface {
	List() ([]string, error)
	Deploy(fname string, c *config.Config) (string, error)
	Delete(fname string) error
}
