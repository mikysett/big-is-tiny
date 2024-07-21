package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = `Usage: bit [-v | --verbose] [-cleanup] [-dryrun] [-p | --platform] [-h | --help] <path to config file>

If not specified the default path to the config file is './bit_config.json'

  -cleanup
        delete branches and PRs
  -v, --verbose
        set logs to DEBUG level
  -dryrun
        do not create branches or PRs
  -p, --platform
		platform used for PRs, can be "github" (default) or "azure"
  -h, --help
        print this help information
`

func getFlags(progName string, args []string) (*Flags, error) {
	rawFlags := flag.NewFlagSet(progName, flag.ExitOnError)

	var verbose, cleanup, dryRun bool
	var rawPlatform string
	var platform Platform
	rawFlags.BoolVar(&cleanup, "cleanup", false, "delete branches and PRs")
	rawFlags.BoolVar(&verbose, "verbose", false, "set logs to DEBUG level")
	rawFlags.BoolVar(&verbose, "v", false, "set logs to DEBUG level")
	rawFlags.BoolVar(&dryRun, "dryrun", false, "do not create branches or PRs")
	rawFlags.StringVar(&rawPlatform, "platform", "github", "platform used for PRs, can be `github` (default) or `azure`")
	rawFlags.StringVar(&rawPlatform, "p", "github", "platform used for PRs, can be `github` (default) or `azure`")
	rawFlags.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	rawFlags.Parse(args)

	switch rawPlatform {
	case "github":
		platform = Platform(GitHub)
	case "azure":
		platform = Platform(Azure)
	default:
		return nil, fmt.Errorf("platform '%s' is not supported", rawPlatform)
	}

	flags := &Flags{
		Cleanup:  cleanup,
		Verbose:  verbose,
		DryRun:   dryRun,
		Platform: platform,
	}
	if configPath := rawFlags.Arg(0); configPath != "" {
		flags.ConfigPath = configPath
	} else {
		flags.ConfigPath = "bit_config.json"
	}

	return flags, nil
}
