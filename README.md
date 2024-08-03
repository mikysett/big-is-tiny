# Big is Tiny (BiT) - Big changes made simple

[![Unit Tests](https://github.com/mikysett/big-is-tiny/actions/workflows/general.yml/badge.svg)](https://github.com/mikysett/big-is-tiny/actions/workflows/general.yml)

BiT is a basic and simple tool to split your big branches into smaller PRs to improve development speed, reduce reviews times and improve their quality.

## Key features

- Automatically split a big branch in multiple sub-branches and PRs based on domains paths
- Cleanup mode to delete the created branches/PRs
- Dry-run mode to evaluate the changes before to apply them
- Create PRs as draft to refine them before asking reviews
- Templates for domain based commit messages, PRs and branch names
- Supported Platforms: `GitHub`, `Azure`
- Customizable with a `config.json` file

## How to install it

Clone the repository locally: `git clone git@github.com:mikysett/big-is-tiny.git`.
Available `make` targets at the root of the repo:

- `make install`: install `bit` globally on your machine
- `make build`: create `bit` binary in `./bin`
- `make test`: runs the unit tests and create coverage reports files in `./src`
- `make clean`: deletes the local binary and the coverage reports files

## How to use it

- Install BiT globally
- `cd` at the root of the repository concerned by the change
- Run `bit 'path/to/config.json'`
- For all available flags run `bit --help`

## Hints

- It's important that the repo do not have any uncommitted changes at the moment of the execution as those could be added to the generated PRs if they match the paths
- BiT will fetch the changed files from the remote branch, so if you have commits on your local branches you should push them first
- If you want to create a miscellaneous "catch all" PR with all non-domain changes you can add a domain **at the end** of the config file with the path `./` (if you add it as first domain this one will include all changes as PRs as domains are evaluated from top to bottom)
- At the end of the execution if there is files that were not included in any PR they will still be there as uncommitted changes, you may want to `git stash` them or `git reset --hard` in order to remove them

### Example of a configuration file

- You will find example configs in `/example_config` directory
- A dummy repository [bit_test_repo](https://github.com/mikysett/bit_test_repo) can be forked and used as a playground with those config files
- Mandatory fields are:
  - `settings.mainBranch`
  - `settings.remote`
  - `settings.branchToSplit`
  - At least one domain
  - Domains should always have at least the `path` field
- Templates placeholders:
  - `{{change_id}}`: `id`
  - `{{domain_id}}`: `domain.id`
  - `{{domain_name}}`: `domain.name`
  - `{{team_name_1}}`: `domain.teams[0].name` (notice the template counting starts with `1`)
  - `{{team_url_1}}`: `domain.Teams[0].Url`

## Prerequisites

- [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Download and install Golang](https://go.dev/doc/install)
- Depending on the chosen platform for the Pull Requests:
  - [GitHub CLI](https://cli.github.com/)
  - [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli)

## Limits and known issues

- BiT has only been tested on Linux and MacOS
- Under the hood vanilla `git` commands are called, this made it faster to implement but brings limitations in performance and stability (if `git` changes some of its returned values BiT may break)
- Paths are plain strings, this limits portability
- The changes are not done in a transaction style, which means if the operation fails mid-way you may find the repository in an unwanted state and you may need to do manual cleanup or run `bit -cleanup path/to/your/config.json`
- Wildcards are not supported for domains paths

## License

MIT license

## Contributing

Contributions are welcome, please create issues or open PRs following usual best practices and common sense.
