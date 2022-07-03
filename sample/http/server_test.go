package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dannykopping/errata/sample/errata"
	"github.com/dannykopping/errata/sample/exec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponses(t *testing.T) {
	server := NewServer()

	tests := []struct {
		name               string
		command            string
		args               []string
		expectedResponse   *exec.Result
		expectedStatusCode int
		expectedErrataCode string
	}{
		// expected
		{
			name:    "echo",
			command: "/usr/bin/echo",
			args:    []string{"-n", "hello world"},
			expectedResponse: &exec.Result{
				Stdout:   "hello world",
				Stderr:   "",
				ExitCode: 0,
			},
			expectedStatusCode: http.StatusOK,
			expectedErrataCode: "",
		},
		{
			name:               "missing command",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrataCode: errata.MissingCommandErrCode,
		},
		{
			name:               "invalid command path",
			command:            "/usr/bin/nope",
			expectedStatusCode: http.StatusNotFound,
			expectedErrataCode: errata.ScriptNotFoundErrCode,
		},

		// unexpected, wrapped
		{
			name:    "permission denied",
			command: "/usr/bin/kill",
			args:    []string{"1"},
			expectedResponse: &exec.Result{
				Stdout:   "",
				Stderr:   "kill: (1): Operation not permitted\n",
				ExitCode: 1,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrataCode: errata.ScriptExecutionFailedErrCode,
		},
		{
			name:    "error in script",
			command: "/usr/bin/ls",
			args:    []string{"/nonexistent"},
			expectedResponse: &exec.Result{
				Stdout:   "",
				Stderr:   "/usr/bin/ls: cannot access '/nonexistent': No such file or directory\n",
				ExitCode: 2,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrataCode: errata.ScriptExecutionFailedErrCode,
		},
	}

	for _, tt := range tests {
		form := url.Values{}
		form.Set("command", tt.command)
		form.Set("args", strings.Join(tt.args, " "))

		req := httptest.NewRequest("POST", "/exec", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := server.Test(req, 2000)
		require.NoError(t, err)

		assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
		assert.Equal(t, tt.expectedErrataCode, resp.Header.Get("X-Errata-Code"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		if tt.expectedResponse != nil {
			var res *exec.Result
			assert.NoError(t, json.Unmarshal(body, &res))

			assert.EqualValues(t, tt.expectedResponse, res)
		}
	}
}
