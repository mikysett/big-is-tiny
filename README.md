# Big is Tiny (BiT) - Big changes made simple

BiT is a basic and simple tool to split your big branches into smaller PR to improve development speed, reduce reviews times and improve their quality.

## Key features

- Automatically split a big branch in multiple sub-branches and PRs
- Cleanup mode to delete the created branches/PRs
- Dry-run mode to evaluate the changes before to apply them
- Create PRs as draft to refine them before asking reviews
- Templates for commit messages, PRs and branch names
- Supported Platforms: `GitHub`, `Azure`
- Customizable with a `config.json` file

## How to install it

Available `make` targets at the root of the repo:

- `make install`: install `bit` globally on your machine
- `make build`: create `bit` binary in `./bin`
- `make test`: runs the unit tests and create coverage reports files in `./src`
- `make clean`: deletes the local binary and the coverage reports files

## How to use it

## Prerequisites

- [Download and install Golang](https://go.dev/doc/install)
- [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- Depending on the chosen platform for the Pull Requests
  - [GitHub CLI](https://cli.github.com/)
  - [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli)

## Limits and known issues

- BiT has only been tested on Linux

## License

MIT license

## Contributing

Contributions are welcome, please create issues or open PRs following usual best practices and common sense.