package main

import (
	"fmt"
	"os"

	"github.com/imroc/req/v3"
	"github.com/joho/godotenv"
)

// GithubClient is the go client for GitHub API.
type GitGuardianClient struct {
	*req.Client
	isLogged bool
}

type APIError struct {
	Detail string `json:"detail"`
}

func NewGitGuardianClient() *GitGuardianClient {
	loadErr := godotenv.Load()
	if loadErr != nil {
		fmt.Println("Error loading .env file")
	}
	c := req.C().
		// All GitHub API requests use the same base URL.
		SetBaseURL("https://api.gitguardian.com/v1").
		// Enable dump at the request level for each request, which dump content into
		// memory (not print to stdout), we can record dump content only when unexpected
		// exception occurs, it is helpful to troubleshoot problems in production.
		EnableDumpEachRequest().
		// Unmarshal all GitHub error response into struct.
		SetCommonErrorResult(&APIError{}).
		// Handle common exceptions in response middleware.
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
				return nil
			}
			if apiErr, ok := resp.ErrorResult().(*APIError); ok {
				// Server returns an error message, convert it to human-readable go error.
				resp.Err = apiErr
				return nil
			}
			// Corner case: neither an error state response nor a success state response,
			// dump content to help troubleshoot.
			if !resp.IsSuccessState() {
				return fmt.Errorf("bad response, raw dump:\n%s", resp.Dump())
			}
			return nil
		})

	return &GitGuardianClient{
		Client: c,
	}
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("API error: %s", e.Detail)
	return msg
}

// LoginWithToken login with GitGuardian API token
// GitGuardian API doc: https://api.gitguardian.com/docs#section/Authentication
func (c *GitGuardianClient) LoginWithToken(token string) *GitGuardianClient {
	c.SetCommonHeader("Authorization", "Token "+os.Getenv("GITGUARDIAN_API_KEY"))
	c.isLogged = true
	return c
}

// IsLogged return true is user is logged in, otherwise false.
func (c *GitGuardianClient) IsLogged() bool {
	return c.isLogged
}

// SetDebug enable debug if set to true, disable debug if set to false.
func (c *GitGuardianClient) SetDebug(enable bool) *GitGuardianClient {
	if enable {
		c.EnableDebugLog()
		c.EnableDumpAll()
	} else {
		c.DisableDebugLog()
		c.DisableDumpAll()
	}
	return c
}
