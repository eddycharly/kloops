
# KLoops commands

The table below lists the commands KLoops understands.

In addition, you can prefix commands with `kl-`. For example, `/meow` and `/kl-meow` are equivalent.

This is because Gitlab hijacks some slash commands for its own quick actions, and we never get notified.
In practice, we don’t need the `kl-` prefix for everything, just commands that also are quick actions,
but we opted to play it safe and add `kl-` prefixes for every command just in case Gitlab eventually adds conflicting quick actions.

There is a [PR open](https://gitlab.com/gitlab-org/gitlab/-/issues/215934) for them to send webhook events for quick actions,
at which point we wouldn't need to worry about it any more, but for now, we do need the `/kl-(command)`.


| Command | Argument | Description | Examples | Plugin |
|---------|----------|-------------|----------|--------|
| `/[un]cc` | Optional | Assigns an assignee to the PR or issue or requests a review from the user(s) | <ul><li>`/cc`</li><li>`/uncc`</li></ul> | assign |
| `/[un]assign` | Optional | Assigns an assignee to the PR or issue or requests a review from the user(s) | <ul><li>`/assign`</li><li>`/unassign`</li></ul> | assign |
| `/meow` | Optional | Add a cat image to the issue or PR | <ul><li>`/meow`</li></ul> | cat |
| `/meowvie` | Optional | Add a cat image to the issue or PR | <ul><li>`/meowvie`</li></ul> | cat |
| `/woof` |  | Add a dog image to the issue or PR | <ul><li>`/woof`</li></ul> | dog |
| `/bark` |  | Add a dog image to the issue or PR | <ul><li>`/bark`</li></ul> | dog |
| `/this-is-fine` |  | Add a dog image to the issue or PR | <ul><li>`/this-is-fine`</li></ul> | dog |
| `/this-is-not-fine` |  | Add a dog image to the issue or PR | <ul><li>`/this-is-not-fine`</li></ul> | dog |
| `/this-is-unbearable` |  | Add a dog image to the issue or PR | <ul><li>`/this-is-unbearable`</li></ul> | dog |
| `/honk` |  | Add a goose image to the issue or PR | <ul><li>`/honk`</li></ul> | goose |
| `/[remove-]help` |  | Applies or removes the 'help wanted' and 'good first issue' labels to an issue. | <ul><li>`/help`</li><li>`/remove-help`</li></ul> | help |
| `/[remove-]good-first-issue` |  | Applies or removes the 'help wanted' and 'good first issue' labels to an issue. | <ul><li>`/good-first-issue`</li><li>`/remove-good-first-issue`</li></ul> | help |
| `/hold` | Optional | Adds or removes the `do-not-merge/hold` Label which is used to indicate that the PR should not be automatically merged. | <ul><li>`/hold`</li></ul> | hold |
| `/[remove-]area` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/area`</li><li>`/remove-area`</li></ul> | label |
| `/[remove-]committee` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/committee`</li><li>`/remove-committee`</li></ul> | label |
| `/[remove-]kind` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/kind`</li><li>`/remove-kind`</li></ul> | label |
| `/[remove-]language` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/language`</li><li>`/remove-language`</li></ul> | label |
| `/[remove-]priority` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/priority`</li><li>`/remove-priority`</li></ul> | label |
| `/[remove-]sig` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/sig`</li><li>`/remove-sig`</li></ul> | label |
| `/[remove-]triage` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/triage`</li><li>`/remove-triage`</li></ul> | label |
| `/[remove-]wg` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/wg`</li><li>`/remove-wg`</li></ul> | label |
| `/[remove-]label` | Mandatory | Applies or removes a label from one of the recognized types of labels. | <ul><li>`/label`</li><li>`/remove-label`</li></ul> | label |
| `/pony` | Optional | Add a little pony image to the issue or PR. A particular pony can optionally be named for a picture of that specific pony. | <ul><li>`/pony`</li></ul> | pony |
| `/[un]shrug` |  | Adds or removes the ¯\_(ツ)_/¯ label | <ul><li>`/shrug`</li><li>`/unshrug`</li></ul> | shrug |
| `/stage` | Mandatory | Labels the stage of an issue as alpha/beta/stable | <ul><li>`/stage`</li></ul> | stage |
| `/remove-stage` | Mandatory | Removes the stage label of an issue as alpha/beta/stable | <ul><li>`/remove-stage`</li></ul> | stage |
| `/joke` |  | Tells a joke. | <ul><li>`/joke`</li></ul> | yuks |