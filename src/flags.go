package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = `Usage: bit [-v | --verbose] [-cleanup] [-p | --platform] [-m | --markdown] [-o | --output] [-h | --help] <path to config file>

If not specified the default path to the config file is './bit_config.json'

  -cleanup
        delete branches and PRs
  -v, --verbose
        set logs to DEBUG level
  -p, --platform
        platform used for PRs, can be "github" (default) or "azure"
  -m, --markdown
        format the created PRs in markdown (default json)
  -o, --output
        writes the results in the specified file
  -d, --allow-deletions
        also updates file deletions from the source branch (git --no-overlay flag)
  -h, --help
        print this help information
`

func getFlags(progName string, args []string) (*Flags, error) {
	rawFlags := flag.NewFlagSet(progName, flag.ExitOnError)

	var verbose, cleanup, markdownOut, allowDeletions bool
	var rawPlatform, fileOut string
	var platform Platform
	rawFlags.BoolVar(&cleanup, "cleanup", false, "delete branches and PRs")
	rawFlags.BoolVar(&verbose, "verbose", false, "set logs to DEBUG level")
	rawFlags.BoolVar(&verbose, "v", false, "set logs to DEBUG level")
	rawFlags.BoolVar(&markdownOut, "markdown", false, "format the created PRs in markdown (default json)")
	rawFlags.BoolVar(&markdownOut, "m", false, "format the created PRs in markdown (default json)")
	rawFlags.BoolVar(&allowDeletions, "d", false, "writes the results in the specified file")
	rawFlags.BoolVar(&allowDeletions, "allow-deletions", false, "writes the results in the specified file")
	rawFlags.StringVar(&rawPlatform, "platform", "github", "platform used for PRs, can be `github` (default) or `azure`")
	rawFlags.StringVar(&rawPlatform, "p", "github", "platform used for PRs, can be `github` (default) or `azure`")
	rawFlags.StringVar(&fileOut, "output", "", "writes the results in the specified file")
	rawFlags.StringVar(&fileOut, "o", "", "writes the results in the specified file")
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
		Cleanup:        cleanup,
		Verbose:        verbose,
		Platform:       platform,
		MarkdownOut:    markdownOut,
		FileOut:        fileOut,
		AllowDeletions: allowDeletions,
	}
	if configPath := rawFlags.Arg(0); configPath != "" {
		flags.ConfigPath = configPath
	} else {
		flags.ConfigPath = "bit_config.json"
	}

	return flags, nil
}
