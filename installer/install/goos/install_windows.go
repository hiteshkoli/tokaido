package goos

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/system/fs"
)

var baseInstallPath = "/AppData/Local/Ironstar/Tokaido"
var binaryName = "tok-windows.exe"

// InstallTokBinary - Install a selected tok version and returns install path
func InstallTokBinary(version string) (string, error) {
	p := filepath.Join(fs.HomeDir(), baseInstallPath, version, "tok")
	b := filepath.Join(p, binaryName)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	sb := filepath.Join(pwd, "tokaido", binaryName)

	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fmt.Println("There was an error creating the install directory")

		log.Fatal(err)
	}

	fs.Copy(sb, b)
	// Change file permission bit
	err = os.Chmod(b, 0755)
	if err != nil {
		panic(err)
	}

	CreateSymlink(b)

	return b, nil
}