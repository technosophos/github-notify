package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Environment variable names.
const (
	// EnvState is one of pending, failure, error, success.
	EnvState = "GH_STATE"
	// EnvDescription is the description part of the status notification. This will
	// be visible in the GitHub UI
	EnvDescription = "GH_DESCRIPTION"
	// EnvContext is the "status context", which is a freeform string that GitHub uses to group notifications.
	// When in doubt, leave blank or set to "brigade"
	EnvContext = "GH_CONTEXT"
	// EnvTargetURL is the URL that the status message should point back to. An empty string is allowed.
	EnvTargetURL = "GH_TARGET_URL"
	// EnvToken is the OAuth2 static token
	EnvToken = "GH_TOKEN"
	// EnvRepo captures the GitHub repository as org/project, e.g. foo/bar
	EnvRepo = "GH_REPO"
	// EnvCommit is the commit to apply this to. A commit can be just about any refish.
	EnvCommit = "GH_COMMIT"
)

// Status States
const (
	StatePending = "pending"
	StateFailure = "failure"
	StateError   = "error"
	StateSuccess = "success"
)

func main() {
	commit := envOr(EnvCommit, "master")
	state := strings.ToLower(envOr(EnvState, StateSuccess))
	repo := envOr(EnvRepo, "")
	desc := envOr(EnvDescription, "updated")
	ctx := envOr(EnvContext, "github-notify")
	u := envOr(EnvTargetURL, "")

	tok := os.Getenv(EnvToken)
	if tok == "" {
		fmt.Fprintf(os.Stderr, "Environment variable %s is required\n", EnvToken)
		os.Exit(1)
	}

	if !isValidState(state) {
		fmt.Fprintf(os.Stderr, "Invalid state: %q\n", state)
	}

	status := &github.RepoStatus{
		State:       &state,
		Description: &desc,
		Context:     &ctx,
		TargetURL:   &u,
	}

	if err := sendNotification(tok, commit, repo, status); err != nil {
		fmt.Fprintf(os.Stderr, "Failed status updated: %s\n", err)
		os.Exit(2)
	}
}

func envOr(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func isValidState(s string) bool {
	switch s {
	case StatePending, StateFailure, StateError, StateSuccess:
		return true
	}
	return false
}

func sendNotification(token, commit, repo string, status *github.RepoStatus) error {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return fmt.Errorf("expected repository for form org/repo, got %q", repo)
	}

	t := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	c := context.Background()
	tc := oauth2.NewClient(c, t)
	client := github.NewClient(tc)

	_, _, err := client.Repositories.CreateStatus(
		c,
		parts[0],
		parts[1],
		commit,
		status,
	)

	return err
}
