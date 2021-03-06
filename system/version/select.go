package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/ironstar-io/tokaido/system/version/goos"
	"github.com/ironstar-io/tokaido/utils"
)

const minimumTokVersion = "1.12.0"

// Select - Change the users' Tokaido to their selection
func Select(selection string) {
	v := Get().Version
	cv := strings.Replace(v, "v", "", 0)
	cs, _ := semver.Parse(cv)

	sv, err := semver.Parse(selection)
	if err != nil {
		fmt.Println("Invalid semver selection supplied. Exiting...")
		return
	}

	// Checks if the current version is Equal to the selected version and exits if so
	if sv.EQ(cs) == true {
		fmt.Println("Selected version (" + sv.String() + ") is the same as the currently active version. Exiting...")
		return
	}

	mv, _ := semver.Parse(minimumTokVersion)
	// Checks if the selected version is Lesser Than the minimum version
	if sv.LT(mv) == true {
		fmt.Println("Selected version (" + sv.String() + ") is less than the minimum allowed version (" + minimumTokVersion + "). Exiting...")
		return
	}

	confirmChange := utils.ConfirmationPrompt("This will change your Tokaido version to "+sv.String()+".\n\nAre you sure?", "y")
	if confirmChange == false {
		fmt.Println("Exiting...")
		return
	}

	ip := GetInstallPath(sv.String())
	// Empty string if not installed, in which case, download and symlink
	if ip == "" {
		// Download & install selected release from GH
		p, err := DownloadAndInstall(sv.String())
		if err != nil {
			fmt.Println("Tokaido was unable to upgrade you to the selected version: ", err.Error())
			os.Exit(1)
		}

		ip = p
	}

	// This running instance is saved, now we activate it as the default 'installed' version
	goos.ActivateSavedVersion(sv.String())

	fmt.Println("Successfully changed Tokaido to version " + sv.String())
}
