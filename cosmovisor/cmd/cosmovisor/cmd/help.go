package cmd

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/cosmovisor"
	"github.com/cosmos/cosmos-sdk/cosmovisor/errors"
)

// ShouldGiveHelp checks the env and provided args to see if help is needed or being requested.
// Help is needed if at least one of the following is true:
// * the cosmovisor.EnvName env var isn't set.
// * the cosmovisor.EnvHome env var isn't set.
// Help is requested if one of the following is true:
// * the first arg is "help"
// * any args are "-h"
// * any args are "--help"
func ShouldGiveHelp(args []string) bool {
	if len(os.Getenv(cosmovisor.EnvName)) == 0 || len(os.Getenv(cosmovisor.EnvHome)) == 0 {
		return true
	}
	if len(args) == 0 {
		return false
	}
	if args[0] == "help" {
		return true
	}
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

// DoHelp outputs help text, config info, and attempts to run the binary with the --help flag.
func DoHelp() {
	// Output the help text
	fmt.Println(GetHelpText())
	// If the config isn't valid, say what's wrong and we're done.
	cfg, cerr := cosmovisor.GetConfigFromEnv()
	switch err := cerr.(type) {
	case nil:
		// Nothing to do. Move on.
	case *errors.MultiError:
		fmt.Fprintf(os.Stderr, "[cosmovisor] multiple configuration errors found:\n")
		for i, e := range err.GetErrors() {
			fmt.Fprintf(os.Stderr, "  %d: %v\n", i+1, e)
		}
		return
	default:
		fmt.Fprintf(os.Stderr, "[cosmovisor] %v\n", err)
		return
	}
	// If the config's legit, output what we see it as.
	fmt.Println("[cosmovisor] config is valid:")
	fmt.Println(cfg.DetailString())
	// Attempt to run the configured binary with the --help flag.
	if err := cosmovisor.RunHelp(cfg, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "[cosmovisor] %v\n", err)
	}
}

// GetHelpText creates the help text multi-line string.
func GetHelpText() string {
	return fmt.Sprintf(`Cosmosvisor - A process manager for Cosmos SDK application binaries.

Cosmovisor is a wrapper for a Cosmos SDK based App (set using the required %s env variable).
It starts the App by passing all provided arguments and monitors the %s/data/upgrade-info.json
file to perform an update. The upgrade-info.json file is created by the App x/upgrade module
when the blockchain height reaches an approved upgrade proposal. The file includes data from
the proposal. Cosmovisor interprets that data to perform an update: switch a current binary
and restart the App.

Configuration of Cosmovisor is done through environment variables, which are
documented in: https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor/README.md
`, cosmovisor.EnvName, cosmovisor.EnvHome)
}
