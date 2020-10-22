package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	validMethods = map[string]bool{
		http.MethodGet:    true,
		http.MethodPost:   true,
		http.MethodPut:    true,
		http.MethodPatch:  true,
		http.MethodDelete: true,
	}
)

type Config struct {
	Routes []Route
}

type Route struct {
	Path    string   `yaml:"path"`
	Methods []Method `yaml:"methods"`
}

type Method struct {
	Method   string    `yaml:"method"`
	Headers  *Headers  `yaml:"headers"`
	Response *Response `yaml:"response"`
}

type Headers struct {
	ContentType string `yaml:"contentType"`
}

type Response struct {
	ContentType string `yaml:"contentType"`
	File        string `yaml:"file"`
	StatusCode  int    `yaml:"statusCode"`
}

//LoadConfig load and parse config file
func LoadConfig() (*Config, error) {
	b, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("unable to read config file %w", err)
	}
	var routes = []Route{}
	err = yaml.Unmarshal(b, &routes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config file %w", err)
	}
	c := &Config{
		Routes: routes,
	}
	return validateConfig(c)
}

func validateConfig(c *Config) (*Config, error) {
	paths := map[string]int{}
	for _, r := range c.Routes {
		if r.Path == "" {
			return nil, fmt.Errorf("unexpected path %s", r.Path)
		}
		paths[r.Path] += 1
		if paths[r.Path] > 1 {
			return nil, fmt.Errorf("redandant path %s", r.Path)
		}
		for _, m := range r.Methods {
			if _, ok := validMethods[m.Method]; !ok {
				return nil, fmt.Errorf("unexpected method: %s \n valid methods are %+v", m.Method, validMethodsString())
			}
			if m.Headers == nil {
				m.Headers = &Headers{
					ContentType: "application/json",
				}
			}
			if m.Response == nil {
				m.Response = &Response{
					ContentType: "application/json",
					StatusCode:  http.StatusOK,
				}
			}
		}
	}
	return c, nil
}

func validMethodsString() string {
	var valid []string
	for v := range validMethods {
		valid = append(valid, v)
	}
	return strings.Join(valid, ", ")
}
