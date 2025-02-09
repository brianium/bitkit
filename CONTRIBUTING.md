# Conventions

## Branch Names

Every branch should be named `<type>/<short-topic-description>` where type is one of:

- `feature` - new functionality that adds value for the user
- `bugfix` - corrects or restores functionality to an existing feature
- `plumbing` - development chores

The default base for every branch is `master`. Every PR into `master` should only contain changes that are safe to deploy to production.

# Workflows

## GitHub Labels

| Label | Type | Description |
| :-: |:-:| :-- |
| **feature** | Category | New functionality / `feature` branch |
| **bug** | Category | Broken functionality / `bugfix` branch |
| **chore** | Category | Non-functional change / `plumbing` branch |
| **question** | Category | Asynchronous discussion threads |
| **discuss** | Flag | Flagged for next in-person discussion |


## Reviews

As a PR reviewer:

1. **DON'T** click "Start a review" unless you want GitHub to include all your comments in your review summary
2. Test the app using the "View deployment" button on the "Conversation" tab
3. Browse the "Files changed" tab and leave comments
4. Try to use wording that distinguishes general/positive comments from actionable ones
5. If changes are required, select "Request changes" from the "Review changes" menu, else "Approve"
6. Optionally leave a summary comment when submitting the review

As the PR owner:

1. Give all reviewers adequate time to review the PR
2. Wait for at least one positive review (optional for chores)
3. **DON'T** merge the PR if there are negative reviews
4. Request a second review after making significant changes

## Migrations

Testing migrations is sort of a pain. Follow the conventions laid down by pgmigrate of a file name following the convention
of `V<version>__<description>.sql`.

Migrations will only be run on a merge into master, and deploys are contingent on their success.

## Deploying

Deployment occurs when a pull request is merged into master. It is done through the CI environment

**production**

Current production instance of the Go API is running at:

[https://api.bitkit.live](https://api.bitkit.live)

The root API Gateway URL provided by amazon is

[https://u5ganepu2f.execute-api.us-west-1.amazonaws.com/production/](https://u5ganepu2f.execute-api.us-west-1.amazonaws.com/production/)

**staging**

Review branches are available only for the UI at this time
