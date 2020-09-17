import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import {
  Divider,
  GridList,
  GridListTile,
  Paper,
  Typography
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  paper: {
    height: "100%",
    width: "100%",
  },
}));
const h = {
  "PluginHelp": {
    "approve": {
      "Description": "The approve plugin implements a pull request approval process that manages the 'approved' label and an approval notification comment. Approval is achieved when the set of users that have approved the PR is capable of approving every file changed by the PR. A user is able to approve a file if their username or an alias they belong to is listed in the 'approvers' section of an OWNERS file in the directory of the file or higher in the directory tree.\n\u003cbr\u003e\n\u003cbr\u003ePer-repo configuration may be used to require that PRs link to an associated issue before approval is granted. It may also be used to specify that the PR authors implicitly approve their own PRs.\n\u003cbr\u003eFor more information see \u003ca href=\"https://git.k8s.io/test-infra/prow/plugins/approve/approvers/README.md\"\u003ehere\u003c/a\u003e.",
      "Config": {
        "tektoncd/catalog": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/chains": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/cli": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/community": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/dashboard": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/experimental": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/friends": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/homebrew-tools": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/hub": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/operator": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/pipeline": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/plumbing": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/triggers": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively.",
        "tektoncd/website": "Pull requests do not  require an associated issue.\u003cbr\u003ePull request authors do not  implicitly approve their own PRs.\u003cbr\u003eThe /lgtm [cancel] command(s) will not  act as approval.\u003cbr\u003eA GitHub approved or changes requested review will  act as approval or cancel respectively."
      },
      "Events": [
        "pull_request",
        "pull_request_review",
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/approve [no-issue|cancel]",
          "Featured": true,
          "Description": "Approves a pull request",
          "Examples": [
            "/approve",
            "/approve no-issue"
          ],
          "WhoCanUse": "Users listed as 'approvers' in appropriate OWNERS files."
        }
      ]
    },
    "assign": {
      "Description": "The assign plugin assigns or requests reviews from users. Specific users can be assigned with the command '/assign @user1' or have reviews requested of them with the command '/cc @user1'. If no user is specified the commands default to targeting the user who created the command. Assignments and requested reviews can be removed in the same way that they are added by prefixing the commands with 'un'.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/[un]assign [[@]\u003cusername\u003e...]",
          "Featured": true,
          "Description": "Assigns an assignee to the PR",
          "Examples": [
            "/assign",
            "/unassign",
            "/assign @k8s-ci-robot"
          ],
          "WhoCanUse": "Anyone can use the command, but the target user must be an org member, a repo collaborator, or should have previously commented on the issue or PR."
        },
        {
          "Usage": "/[un]cc [[@]\u003cusername\u003e...]",
          "Featured": true,
          "Description": "Requests a review from the user(s).",
          "Examples": [
            "/cc",
            "/uncc",
            "/cc @k8s-ci-robot"
          ],
          "WhoCanUse": "Anyone can use the command, but the target user must be a member of the org that owns the repository."
        }
      ]
    },
    "blockade": {
      "Description": "The blockade plugin blocks pull requests from merging if they touch specific files. The plugin applies the 'do-not-merge/blocked-paths' label to pull requests that touch files that match a blockade's block regular expression and none of the corresponding exception regular expressions.",
      "Config": {

      },
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "blunderbuss": {
      "Description": "The blunderbuss plugin automatically requests reviews from reviewers when a new PR is created. The reviewers are selected based on the reviewers specified in the OWNERS files that apply to the files modified by the PR.",
      "Config": {
        "": "Blunderbuss is currently configured to request reviews from 2 reviewers."
      },
      "Events": [
        "pull_request",
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/auto-cc",
          "Featured": false,
          "Description": "Manually request reviews from reviewers for a PR. Useful if OWNERS file were updated since the PR was opened.",
          "Examples": [
            "/auto-cc"
          ],
          "WhoCanUse": "Anyone"
        }
      ]
    },
    "branchcleaner": {
      "Description": "The branchcleaner plugin automatically deletes source branches for merged PRs between two branches on the same repository. This is helpful to keep repos that don't allow forking clean.",
      "Config": null,
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "bugzilla": {
      "Description": "The bugzilla plugin ensures that pull requests reference a valid Bugzilla bug in their title.",
      "Config": {

      },
      "Events": [
        "pull_request",
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/bugzilla refresh",
          "Featured": false,
          "Description": "Check Bugzilla for a valid bug referenced in the PR title",
          "Examples": [
            "/bugzilla refresh"
          ],
          "WhoCanUse": "Anyone"
        }
      ]
    },
    "buildifier": {
      "Description": "The buildifier plugin runs buildifier on changes made to Bazel files in a PR. It then creates a new review on the pull request and leaves warnings at the appropriate lines of code.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/buildif(y|ier)",
          "Featured": false,
          "Description": "Runs buildifier on changes made to Bazel files in a PR",
          "Examples": [
            "/buildify",
            "/buildifier"
          ],
          "WhoCanUse": "Anyone can trigger this command on a PR."
        }
      ]
    },
    "cat": {
      "Description": "The cat plugin adds a cat image to an issue or PR in response to the `/meow` command.",
      "Config": {
        "": "The cat plugin uses an api key for thecatapi.com stored in ."
      },
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/meow(vie) [CATegory]",
          "Featured": false,
          "Description": "Add a cat image to the issue or PR",
          "Examples": [
            "/meow",
            "/meow caturday",
            "/meowvie clothes"
          ],
          "WhoCanUse": "Anyone"
        }
      ]
    },
    "cherry-pick-unapproved": {
      "Description": "Label PRs against a release branch which do not have the `cherry-pick-approved` label with the `do-not-merge/cherry-pick-not-approved` label.",
      "Config": {
        "": "The cherry-pick-unapproved plugin treats PRs against branch names satisfying the regular expression `^release-.*$` as cherry-pick PRs and adds the following comment:\nThis PR is not for the master branch but does not have the `cherry-pick-approved`  label. Adding the `do-not-merge/cherry-pick-not-approved`  label.\n\nTo approve the cherry-pick, please ping the *kubernetes/patch-release-team* in a comment when ready.\n\nSee also [Kubernetes Patch Releases](https://github.com/kubernetes/sig-release/blob/master/releases/patch-releases.md)"
      },
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "cla": {
      "Description": "The cla plugin manages the application and removal of the 'cncf-cla' prefixed labels on pull requests as a reaction to the cla/linuxfoundation github status context. It is also responsible for warning unauthorized PR authors that they need to sign the CNCF CLA before their PR will be merged.",
      "Config": null,
      "Events": [
        "status",
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/check-cla",
          "Featured": true,
          "Description": "Forces rechecking of the CLA status.",
          "Examples": [
            "/check-cla"
          ],
          "WhoCanUse": "Anyone"
        }
      ]
    },
    "config-updater": {
      "Description": "The config-updater plugin automatically redeploys configuration and plugin configuration files when they change. The plugin watches for pull request merges that modify either of the config files and updates the cluster's configmap resources in response.",
      "Config": null,
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "dco": {
      "Description": "The dco plugin checks pull request commits for 'DCO sign off' and maintains the 'dco' status context, as well as the 'dco' label.",
      "Config": {

      },
      "Events": [
        "pull_request",
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/check-dco",
          "Featured": true,
          "Description": "Forces rechecking of the DCO status.",
          "Examples": [
            "/check-dco"
          ],
          "WhoCanUse": "Anyone"
        }
      ]
    },
    "docs-no-retest": {
      "Description": "The docs-no-retest plugin applies the 'retest-not-required-docs-only' label to pull requests that only touch documentation type files and thus do not need to be retested against the latest master commit before merging.\n\u003cbr\u003eFiles extensions '.md', '.png', '.svg', and '.dia' are considered documentation.",
      "Config": null,
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "dog": {
      "Description": "The dog plugin adds a dog image to an issue or PR in response to the `/woof` command.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/(woof|bark|this-is-{fine|not-fine|unbearable})",
          "Featured": false,
          "Description": "Add a dog image to the issue or PR",
          "Examples": [
            "/woof",
            "/bark",
            "this-is-{fine|not-fine|unbearable}"
          ],
          "WhoCanUse": "Anyone"
        }
      ]
    },
    "golint": {
      "Description": "The golint plugin runs golint on changes made to *.go files in a PR. It then creates a new review on the pull request and leaves golint warnings at the appropriate lines of code.",
      "Config": {
        "": "The golint plugin will report problems with a minimum confidence of 0.800000."
      },
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/lint",
          "Featured": false,
          "Description": "Runs golint on changes made to *.go files in a PR",
          "Examples": [
            "/lint"
          ],
          "WhoCanUse": "Anyone can trigger this command on a PR."
        }
      ]
    },
    "heart": {
      "Description": "The heart plugin celebrates certain GitHub actions with the reaction emojis. Emojis are added to pull requests that make additions to OWNERS or OWNERS_ALIASES files and to comments left by specified \"adorees\".",
      "Config": {
        "": "The heart plugin is configured to react to comments,  satisfying the regular expression , left by the following GitHub users: ."
      },
      "Events": [
        "issue_comment",
        "pull_request"
      ],
      "Commands": null
    },
    "help": {
      "Description": "The help plugin provides commands that add or remove the 'help wanted' and the 'good first issue' labels from issues.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/[remove-](help|good-first-issue)",
          "Featured": false,
          "Description": "Applies or removes the 'help wanted' and 'good first issue' labels to an issue.",
          "Examples": [
            "/help",
            "/remove-help",
            "/good-first-issue",
            "/remove-good-first-issue"
          ],
          "WhoCanUse": "Anyone can trigger this command on a PR."
        }
      ]
    },
    "hold": {
      "Description": "The hold plugin allows anyone to add or remove the 'do-not-merge/hold' Label from a pull request in order to temporarily prevent the PR from merging without withholding approval.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/hold [cancel]",
          "Featured": false,
          "Description": "Adds or removes the `do-not-merge/hold` Label which is used to indicate that the PR should not be automatically merged.",
          "Examples": [
            "/hold",
            "/hold cancel"
          ],
          "WhoCanUse": "Anyone can use the /hold command to add or remove the 'do-not-merge/hold' Label."
        }
      ]
    },
    "invalidcommitmsg": {
      "Description": "The invalidcommitmsg plugin applies the 'do-not-merge/invalid-commit-message' label to pull requests whose commit messages contain @ mentions or keywords which can automatically close issues.",
      "Config": null,
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "label": {
      "Description": "The label plugin provides commands that add or remove certain types of labels. Labels of the following types can be manipulated: 'area/*', 'committee/*', 'kind/*', 'language/*', 'priority/*', 'sig/*', 'triage/*', and 'wg/*'. More labels can be configured to be used via the /label command.",
      "Config": {
        "": "The label plugin will work on \"kind/*\", \"priority/*\" and \"area/*\" labels."
      },
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/[remove-](area|committee|kind|language|priority|sig|triage|wg|label) \u003ctarget\u003e",
          "Featured": false,
          "Description": "Applies or removes a label from one of the recognized types of labels.",
          "Examples": [
            "/kind bug",
            "/remove-area prow",
            "/sig testing",
            "/language zh"
          ],
          "WhoCanUse": "Anyone can trigger this command on a PR."
        }
      ]
    },
    "lgtm": {
      "Description": "The lgtm plugin manages the application and removal of the 'lgtm' (Looks Good To Me) label which is typically used to gate merging.",
      "Config": {

      },
      "Events": [
        "pull_request",
        "pull_request_review",
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/lgtm [cancel] or GitHub Review action",
          "Featured": true,
          "Description": "Adds or removes the 'lgtm' label which is typically used to gate merging.",
          "Examples": [
            "/lgtm",
            "/lgtm cancel",
            "\u003ca href=\"https://help.github.com/articles/about-pull-request-reviews/\"\u003e'Approve' or 'Request Changes'\u003c/a\u003e"
          ],
          "WhoCanUse": "Collaborators on the repository. '/lgtm cancel' can be used additionally by the PR author."
        }
      ]
    },
    "lifecycle": {
      "Description": "Close, reopen, flag and/or unflag an issue or PR as frozen/stale/rotten",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/close",
          "Featured": false,
          "Description": "Closes an issue or PR.",
          "Examples": [
            "/close"
          ],
          "WhoCanUse": "Authors and collaborators on the repository can trigger this command."
        },
        {
          "Usage": "/reopen",
          "Featured": false,
          "Description": "Reopens an issue or PR",
          "Examples": [
            "/reopen"
          ],
          "WhoCanUse": "Authors and collaborators on the repository can trigger this command."
        },
        {
          "Usage": "/[remove-]lifecycle \u003cfrozen|stale|rotten\u003e",
          "Featured": false,
          "Description": "Flags an issue or PR as frozen/stale/rotten",
          "Examples": [
            "/lifecycle frozen",
            "/remove-lifecycle stale"
          ],
          "WhoCanUse": "Anyone can trigger this command."
        }
      ]
    },
    "mergecommitblocker": {
      "Description": "The merge commit blocker plugin adds the do-not-merge/contains-merge-commits label to pull requests that contain merge commits",
      "Config": null,
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "milestone": {
      "Description": "The milestone plugin allows members of a configurable GitHub team to set the milestone on an issue or pull request.",
      "Config": {
        "": "The milestone maintainers team is the GitHub team \"\" with ID: 0."
      },
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/milestone \u003cversion\u003e or /milestone clear",
          "Featured": false,
          "Description": "Updates the milestone for an issue or PR",
          "Examples": [
            "/milestone v1.10",
            "/milestone v1.9",
            "/milestone clear"
          ],
          "WhoCanUse": "Members of the milestone maintainers GitHub team can use the '/milestone' command."
        }
      ]
    },
    "milestoneapplier": {
      "Description": "The milestoneapplier plugin automatically applies the configured milestone for the base branch after a PR is merged. If a PR targets a non-default branch, it also adds the milestone when the PR is opened.",
      "Config": {

      },
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "milestonestatus": {
      "Description": "The milestonestatus plugin allows members of the milestone maintainers GitHub team to specify the 'status/*' label that should apply to a pull request.",
      "Config": {
        "": "The milestone maintainers team is the GitHub team \"\" with ID: 0."
      },
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/status (approved-for-milestone|in-progress|in-review)",
          "Featured": false,
          "Description": "Applies the 'status/' label to a PR.",
          "Examples": [
            "/status approved-for-milestone",
            "/status in-progress",
            "/status in-review"
          ],
          "WhoCanUse": "Members of the milestone maintainers GitHub team can use the '/status' command. This team is specified in the config by providing the GitHub team's ID."
        }
      ]
    },
    "override": {
      "Description": "The override plugin allows repo admins to force a github status context to pass",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/override [context]",
          "Featured": false,
          "Description": "Forces a github status context to green (one per line).",
          "Examples": [
            "/override pull-repo-whatever",
            "/override ci/circleci",
            "/override deleted-job"
          ],
          "WhoCanUse": "Repo administrators"
        }
      ]
    },
    "owners-label": {
      "Description": "The owners-label plugin automatically adds labels to PRs based on the files they touch. Specifically, the 'labels' sections of OWNERS files are used to determine which labels apply to the changes.",
      "Config": null,
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "pony": {
      "Description": "The pony plugin adds a pony image to an issue or PR in response to the `/pony` command.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/(pony) [pony]",
          "Featured": false,
          "Description": "Add a little pony image to the issue or PR. A particular pony can optionally be named for a picture of that specific pony.",
          "Examples": [
            "/pony",
            "/pony Twilight Sparkle"
          ],
          "WhoCanUse": "Anyone"
        }
      ]
    },
    "project": {
      "Description": "The project plugin allows members of a GitHub team to set the project and column on an issue or pull request.",
      "Config": {
        "tektoncd/catalog": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/chains": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/cli": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/community": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/dashboard": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/experimental": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/friends": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/homebrew-tools": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/hub": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/operator": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/pipeline": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/plumbing": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/triggers": "There are no maintainer team specified for this repo or its org.",
        "tektoncd/website": "There are no maintainer team specified for this repo or its org."
      },
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/project \u003cboard\u003e, /project \u003cboard\u003e \u003ccolumn\u003e, or /project clear \u003cboard\u003e",
          "Featured": false,
          "Description": "Add an issue or PR to a project board and column",
          "Examples": [
            "/project 0.5.0",
            "/project 0.5.0 To do",
            "/project clear 0.4.0"
          ],
          "WhoCanUse": "Members of the project maintainer GitHub team can use the '/project' command."
        }
      ]
    },
    "release-note": {
      "Description": "The releasenote plugin implements a release note process that uses a markdown 'releasenote' code block to associate a release note with a pull request. Until the 'releasenote' block in the pull request body is populated the PR will be assigned the 'do-not-merge/release-note-label-needed' label.\n\u003cbr\u003eThere are three valid types of release notes that can replace this label:\n\u003col\u003e\u003cli\u003ePRs with a normal release note in the 'releasenote' block are given the label 'release-note'.\u003c/li\u003e\n\u003cli\u003ePRs that have a release note of 'none' in the block are given the label 'release-note-none' to indicate that the PR does not warrant a release note.\u003c/li\u003e\n\u003cli\u003ePRs that contain 'action required' in their 'releasenote' block are given the label 'release-note-action-required' to indicate that the PR introduces potentially breaking changes that necessitate user action before upgrading to the release.\u003c/li\u003e\u003c/ol\u003e\nTo use the plugin, in the pull request body text:\n\n```releasenote\n\u003crelease note content\u003e\n```",
      "Config": null,
      "Events": [
        "issue_comment",
        "pull_request"
      ],
      "Commands": [
        {
          "Usage": "/release-note-none",
          "Featured": false,
          "Description": "Adds the 'release-note-none' label to indicate that the PR does not warrant a release note. This is deprecated and ideally \u003ca href=\"https://git.k8s.io/community/contributors/guide/release-notes.md\"\u003ethe release note process\u003c/a\u003e should be followed in the PR body instead.",
          "Examples": [
            "/release-note-none"
          ],
          "WhoCanUse": "PR Authors and Org Members."
        }
      ]
    },
    "require-matching-label": {
      "Description": "The require-matching-label plugin is a configurable plugin that applies a label to issues and/or PRs that do not have any labels matching a regular expression. An example of this is applying a 'needs-sig' label to all issues that do not have a 'sig/*' label. This plugin can have multiple configurations to provide this kind of behavior for multiple different label sets. The configuration allows issue type, PR branch, and an optional explanation comment to be specified.",
      "Config": {
        "": "The plugin has the following configurations:\n\u003cul\u003e\u003cli\u003e\u003c/li\u003e\u003c/ul\u003e"
      },
      "Events": [
        "issue",
        "pull_request"
      ],
      "Commands": null
    },
    "require-sig": {
      "Description": "When a new issue is opened the require-sig plugin adds the \"needs-sig\" label and leaves a comment requesting that a SIG (Special Interest Group) label be added to the issue. SIG labels are labels that have one of the following prefixes: [\"sig/\" \"committee/\" \"wg/\"].\n\u003cbr\u003eOnce a SIG label has been added to an issue, this plugin removes the \"needs-sig\" label and deletes the comment it made previously.",
      "Config": {
        "": "The comment the plugin creates includes this link to a list of the existing groups: \u003cno url provided\u003e"
      },
      "Events": [
        "issue"
      ],
      "Commands": null
    },
    "retitle": {
      "Description": "The retitle plugin allows users to re-title pull requests and issues where GitHub permissions don't allow them to.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/retitle \u003ctitle\u003e",
          "Featured": true,
          "Description": "Edits the pull request or issue title.",
          "Examples": [
            "/retitle New Title"
          ],
          "WhoCanUse": "Collaborators on the repository."
        }
      ]
    },
    "shrug": {
      "Description": "¯\\_(ツ)_/¯",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/[un]shrug",
          "Featured": false,
          "Description": "¯\\_(ツ)_/¯",
          "Examples": [
            "/shrug",
            "/unshrug"
          ],
          "WhoCanUse": "Anyone, ¯\\_(ツ)_/¯"
        }
      ]
    },
    "sigmention": {
      "Description": "The sigmention plugin responds to SIG (Special Interest Group) GitHub team mentions like '@kubernetes/sig-testing-bugs'. The plugin responds in two ways:\n\u003col\u003e\u003cli\u003e The appropriate 'sig/*' and 'kind/*' labels are applied to the issue or pull request. In this case 'sig/testing' and 'kind/bug'.\u003c/li\u003e\n\u003cli\u003e If the user who mentioned the GitHub team is not a member of the organization that owns the repository the bot will create a comment that repeats the mention. This is necessary because non-member mentions do not trigger GitHub notifications.\u003c/li\u003e\u003c/ol\u003e",
      "Config": {
        "": "Labels added by the plugin are triggered by mentions of GitHub teams matching the following regexp:\n(?m)@kubernetes/sig-([\\w-]*)-(misc|test-failures|bugs|feature-requests|proposals|pr-reviews|api-reviews)"
      },
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": null
    },
    "size": {
      "Description": "The size plugin manages the 'size/*' labels, maintaining the appropriate label on each pull request as it is updated. Generated files identified by the config file '.generated_files' at the repo root are ignored. Labels are applied based on the total number of lines of changes (additions and deletions).",
      "Config": {
        "": "The plugin has the following thresholds:\u003cul\u003e\n\u003cli\u003esize/XS:  0-9\u003c/li\u003e\n\u003cli\u003esize/S:   10-29\u003c/li\u003e\n\u003cli\u003esize/M:   30-99\u003c/li\u003e\n\u003cli\u003esize/L:   100-499\u003c/li\u003e\n\u003cli\u003esize/XL:  500-999\u003c/li\u003e\n\u003cli\u003esize/XXL: 1000+\u003c/li\u003e\n\u003c/ul\u003e"
      },
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "skip": {
      "Description": "The skip plugin allows users to clean up GitHub stale commit statuses for non-blocking jobs on a PR.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/skip",
          "Featured": false,
          "Description": "Cleans up GitHub stale commit statuses for non-blocking jobs on a PR.",
          "Examples": [
            "/skip"
          ],
          "WhoCanUse": "Anyone can trigger this command on a PR."
        }
      ]
    },
    "slackevents": {
      "Description": "The slackevents plugin reacts to various GitHub events by commenting in Slack channels.\n\u003col\u003e\u003cli\u003eThe plugin can create comments to alert on manual merges. Manual merges are merges made by a normal user instead of a bot or trusted user.\u003c/li\u003e\n\u003cli\u003eThe plugin can create comments to reiterate SIG mentions like '@kubernetes/sig-testing-bugs' from GitHub.\u003c/li\u003e\u003c/ol\u003e",
      "Config": {
        "": "SIG mentions on GitHub are reiterated for the following SIG Slack channels: ."
      },
      "Events": [
        "push",
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": null
    },
    "stage": {
      "Description": "Label the stage of an issue as alpha/beta/stable",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/[remove-]stage \u003calpha|beta|stable\u003e",
          "Featured": false,
          "Description": "Labels the stage of an issue as alpha/beta/stable",
          "Examples": [
            "/stage alpha",
            "/remove-stage alpha"
          ],
          "WhoCanUse": "Anyone can trigger this command."
        }
      ]
    },
    "trigger": {
      "Description": "The trigger plugin starts tests in reaction to commands and pull request events. It is responsible for ensuring that test jobs are only run on trusted PRs. A PR is considered trusted if the author is a member of the 'trusted organization' for the repository or if such a member has left an '/ok-to-test' command on the PR.\n\u003cbr\u003eTrigger starts jobs automatically when a new trusted PR is created or when an untrusted PR becomes trusted, but it can also be used to start jobs manually via the '/test' command.\n\u003cbr\u003eThe '/retest' command can be used to rerun jobs that have reported failure.",
      "Config": {
        "tektoncd/catalog": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/chains": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/cli": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/community": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/dashboard": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/experimental": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/friends": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/homebrew-tools": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/hub": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/operator": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/pipeline": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/plumbing": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/triggers": "The trusted GitHub organization for this repository is \"tektoncd\".",
        "tektoncd/website": "The trusted GitHub organization for this repository is \"tektoncd\"."
      },
      "Events": [
        "pull_request",
        "push",
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/ok-to-test",
          "Featured": false,
          "Description": "Marks a PR as 'trusted' and starts tests.",
          "Examples": [
            "/ok-to-test"
          ],
          "WhoCanUse": "Members of the trusted organization for the repo."
        },
        {
          "Usage": "/test (\u003cjob name\u003e|all)",
          "Featured": true,
          "Description": "Manually starts a/all test job(s).",
          "Examples": [
            "/test all",
            "/test pull-bazel-test"
          ],
          "WhoCanUse": "Anyone can trigger this command on a trusted PR."
        },
        {
          "Usage": "/retest",
          "Featured": true,
          "Description": "Rerun test jobs that have failed.",
          "Examples": [
            "/retest"
          ],
          "WhoCanUse": "Anyone can trigger this command on a trusted PR."
        }
      ]
    },
    "verify-owners": {
      "Description": "The verify-owners plugin validates OWNERS and OWNERS_ALIASES files and ensures that they always contain collaborators of the org, if they are modified in a PR. On validation failure it automatically adds the 'do-not-merge/invalid-owners-file' label to the PR, and a review comment on the incriminating file(s).",
      "Config": {
        "": "The verify-owners plugin will complain if OWNERS files contain any of the following blacklisted labels: approved, lgtm."
      },
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "welcome": {
      "Description": "The welcome plugin posts a welcoming message when it detects a user's first contribution to a repo.",
      "Config": {

      },
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "wip": {
      "Description": "The wip (Work In Progress) plugin applies the 'do-not-merge/work-in-progress' Label to pull requests whose title starts with 'WIP' or are in the 'draft' stage, and removes it from pull requests when they remove the title prefix or become ready for review. The 'do-not-merge/work-in-progress' Label is typically used to block a pull request from merging while it is still in progress.",
      "Config": null,
      "Events": [
        "pull_request"
      ],
      "Commands": null
    },
    "yuks": {
      "Description": "The yuks plugin comments with jokes in response to the `/joke` command.",
      "Config": null,
      "Events": [
        "GenericCommentEvent (any event for user text)"
      ],
      "Commands": [
        {
          "Usage": "/joke",
          "Featured": false,
          "Description": "Tells a joke.",
          "Examples": [
            "/joke"
          ],
          "WhoCanUse": "Anyone can use the `/joke` command."
        }
      ]
    }
  }
};

function Plugins() {
  const classes = useStyles();

  return (
    <GridList cellHeight={160} cols={6} spacing={20}>
      {Object.entries(h.PluginHelp).map(([key, value]) => (
        <GridListTile cols={1}>
          <Paper className={classes.paper}>
            <Typography>{key}</Typography>
            <Divider />
            <Typography>{value.Description.split(".")[0]}</Typography>
            <Divider />
            <Typography>Details</Typography>
          </Paper>
        </GridListTile>
      ))}
    </GridList>
  );
}

export default Plugins;
