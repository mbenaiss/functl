package gcp

import (
	"os"
	"os/exec"

	"github.com/mbenaiss/functl/api"
	"github.com/mbenaiss/functl/provider"
)

type Client struct {
}

func New() provider.Provider {
	return &Client{}
}

func (c *Client) List() ([]string, error) {
	return []string{}, nil
}

func (c *Client) Deploy(fname string, cfg *api.Config) (string, error) {
	args := []string{"functions", "deploy", fname, "--runtime", "go113", "--trigger-http", "--allow-unauthenticated"}
	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return "", nil
}

func (c *Client) Delete(fname string) error {
	return nil
}
