package provider

import (
	"github.com/mbenaiss/functl/api"
)

type Provider interface {
	List() ([]string, error)
	Deploy(fname string, c *api.Config) (string, error)
	Delete(fname string) error
}
