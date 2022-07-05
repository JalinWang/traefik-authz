package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCasbinAuth(t *testing.T) {

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, "traefik")
	})
	cfg := &Config{
		ModelPath:  "./examples/rbac_with_pattern_model.conf",
		PolicyPath: "./examples/rbac_with_pattern_policy.csv",
	}

	authHandkler, err := New(ctx, next, cfg, "casbin_auth")
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/pen/2", nil)
	require.Nil(t, err)
	req.Header.Add(CasbinAuthHeader, "alice")
	authHandkler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Result().StatusCode, "they should be equal")

	recorder = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/pen/1", nil)
	require.Nil(t, err)
	req.Header.Add(CasbinAuthHeader, "alice")
	authHandkler.ServeHTTP(recorder, req)
	res := recorder.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode, "they should be equal")

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, "traefik\n", string(body), "they should be equal")
}
