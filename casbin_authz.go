package main

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"log"
	"net/http"
)

const (
	CasbinAuthHeader = "X-Casbin-Authorization"
)

type Config struct {
	ModelPath  string `json:"modelPath,omitempty" toml:"modelPath,omitempty" yaml:"modelPath,omitempty"`
	PolicyPath string `json:"policyPath,omitempty" toml:"policyPath,omitempty" yaml:"policyPath,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type casbinAuth struct {
	next     http.Handler
	enforcer *casbin.Enforcer
	name     string
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	enforcer, err := casbin.NewEnforcer(config.ModelPath, config.PolicyPath)
	if err != nil {
		return nil, err
	}

	return &casbinAuth{
		next:     next,
		enforcer: enforcer,
		name:     name,
	}, nil
}

func (c *casbinAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), c.name, casbinTypeName))
	user, path, method := getParam(req)

	ok, err := c.enforcer.Enforce(user, path, method)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "could not operate authorization: %v", err)
		return
	}

	if !ok {
		log.Println("Authorization failed")

		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(rw, "authorization faile: %s, %s, %s", user, path, method)
		return
	}

	log.Println("Authorization succeeded")

	c.next.ServeHTTP(rw, req)
}

func getParam(req *http.Request) (string, string, string) {
	user := req.Header.Get(CasbinAuthHeader)
	path := req.URL.Path
	method := req.Method
	return user, path, method
}
