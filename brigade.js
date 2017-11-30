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
