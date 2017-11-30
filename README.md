# GitHub Notify: Simply Send GitHub Status Updates

This tool sets status updates on GitHub repository commits (or commit-like things).
It is a simple command line for sending such updates.

## Usage

`github-notify` uses environment variables to run. It does this because it is
intended for use within Docker containers.

To test locally, you can do something like this:

```bash
#!/bin/bash
# REQUIRED: The GitHub repository
export GH_REPO="technosophos/github-notify"
# REQUIRED: One of: success, failure, error, pending
export GH_STATE="success"
# REQUIRED: Your GitHub OAuth token.
# Go to your account, Settings -> Developer settings -> Personal access tokens
# to generate one.
export GH_TOKEN="REDACTED"
# RECOMMENDED: The commit-ish thing to update. Default is master.
export GH_COMMIT="62727511e4c87b2d8a5f0ea9d0288bb74ce7dc2d"
# OPTIONAL: A short description displayed to GH users
export GH_DESCRIPTION="doing hard stuff"
# OPTIONAL: A short word used for grouping notifications
export GH_CONTEXT="ci"
# OPTIONAL: A URL that the user can click to learn more
export GH_TARGET_URL="http://technosophos.com"

./github-notify
```

## With Brigade

The primary use-case envisioned for `github-notify` is Brigade scripts. You can
use it in Brigade scripts like this:

```javascript
const { events, Job } = require("brigadier");

events.on("exec", (e, p) => {
  var gh = new Job("gh", "technosophos/github-notify:latest");
  gh.env = {
    GH_REPO: p.repo.name,
    GH_STATE: "success",
    GH_DESCRIPTION: "brigade says YES!",
    GH_CONTEXT: "brigade"

    // We get the token from the project's secrets section.
    GH_TOKEN: p.secrets.githubToken, // YOU MUST SET THIS IN YOUR PROJECT

    // We set this because 'exec' doesn't have a real commit. Normally you
    // would use e.commit.
    GH_COMMIT: "68f9a8f8efab12e6bb8fc4c2b3f2ac6b7051df8f",
  };
  gh.run().then(() => console.log("Status updated."));
});
```

Make sure that in your project configuration, you set:

```
secrets:
  githubToken: YOUR_ACTUAL_TOKEN
```
