
# KLoops plugins

The list below lists the KLoops plugins.

Every plugin receive events from scm providers and react to those events in a specific way.

In addition, plugins can react to comments (on pull requests, issues or reviews). The list of supported commands is document [here](./COMMANDS.md).

- [assign](#assign)
- [branchcleaner](#branchcleaner)
- [cat](#cat)
- [dog](#dog)
- [goose](#goose)
- [help](#help)
- [hold](#hold)
- [label](#label)
- [pony](#pony)
- [shrug](#shrug)
- [size](#size)
- [stage](#stage)
- [welcome](#welcome)
- [wip](#wip)
- [yuks](#yuks)


## assign

The assign plugin assigns or requests reviews from users. Specific users can be assigned with the command '/assign @user1' or have reviews requested of them with the command '/cc @user1'. If no user is specified the commands default to targeting the user who created the command. Assignments and requested reviews can be removed in the same way that they are added by prefixing the commands with 'un'.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/[un]cc`
- `/[un]assign`

## branchcleaner

The branchcleaner plugin automatically deletes source branches for merged PRs between two branches on the same repository. This is helpful to keep repos that don't allow forking clean.

This plugin reacts to the following events:
- `pull_request`


This plugins supports all scm providers.

This plugins has no commands.

## cat

The cat plugin adds a cat image to an issue or PR in response to the `/meow` command.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/meow`
- `/meowvie`

## dog

The dog plugin adds a dog image to an issue or PR in response to the `/woof` command.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/woof`
- `/bark`
- `/this-is-fine`
- `/this-is-not-fine`
- `/this-is-unbearable`

## goose

The goose plugin adds a goose image to an issue or PR in response to the `/honk` command.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/honk`

## help

The help plugin provides commands that add or remove the 'help wanted' and the 'good first issue' labels from issues.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/[remove-]help`
- `/[remove-]good-first-issue`

## hold

The hold plugin allows anyone to add or remove the 'do-not-merge/hold' Label from a pull request in order to temporarily prevent the PR from merging without withholding approval.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/hold`

## label

The label plugin provides commands that add or remove certain types of labels. Labels of the following types can be manipulated: 'area/*', 'committee/*', 'kind/*', 'language/*', 'priority/*', 'sig/*', 'triage/*', and 'wg/*'. More labels can be configured to be used via the /label command.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/[remove-]area`
- `/[remove-]committee`
- `/[remove-]kind`
- `/[remove-]language`
- `/[remove-]priority`
- `/[remove-]sig`
- `/[remove-]triage`
- `/[remove-]wg`
- `/[remove-]label`

## pony

The pony plugin adds a pony image to an issue or PR in response to the `/pony` command.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/pony`

## shrug

`¯\_(ツ)_/¯`

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/[un]shrug`

## size

The size plugin manages the 'size/*' labels, maintaining the appropriate label on each pull request as it is updated. Generated files identified by the config file '.generated_files' at the repo root are ignored. Labels are applied based on the total number of lines of changes (additions and deletions).

This plugin reacts to the following events:
- `pull_request`


This plugins supports all scm providers.

This plugins has no commands.

## stage

Label the stage of an issue as alpha/beta/stable

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/stage`
- `/remove-stage`

## welcome

The welcome plugin posts a welcoming message when it detects a user's first contribution to a repo.

This plugin reacts to the following events:
- `pull_request`


This plugins supports all scm providers.

This plugins has no commands.

## wip

The wip (Work In Progress) plugin applies the 'do-not-merge/work-in-progress' Label to pull requests whose title starts with 'WIP' or are in the 'draft' stage, and removes it from pull requests when they remove the title prefix or become ready for review. The 'do-not-merge/work-in-progress' Label is typically used to block a pull request from merging while it is still in progress.

This plugin reacts to the following events:
- `pull_request`


This plugins supports all scm providers.

This plugins has no commands.

## yuks

The yuks plugin comments with jokes in response to the `/joke` command.

This plugin reacts to the following events:
- `none`


This plugins supports all scm providers.

This plugin has the following commands:
- `/joke`

