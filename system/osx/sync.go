package osx

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"text/template"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/osx/templates"
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

type service struct {
	ProjectName string
	ProjectPath string
	Username    string
}

func createSyncFile() {
	c := conf.GetConfig()
	u, uErr := user.Current()
	if uErr != nil {
		log.Fatal(uErr)
	}

	s := service{
		ProjectName: c.Project,
		ProjectPath: c.Path,
		Username:    u.Username,
	}

	serviceFilename := "tokaido.sync." + s.ProjectName + ".plist"

	tmpl := template.New(serviceFilename)
	tmpl, err := tmpl.Parse(synctmpl.SyncTemplateStr)

	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, s); err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	writeSyncFile(tpl.String(), c.LaunchdPath, serviceFilename)
}

func writeSyncFile(body string, path string, filename string) {
	mkdErr := os.MkdirAll(path, os.ModePerm)
	if mkdErr != nil {
		log.Fatal("Mkdir: ", mkdErr)
	}

	var file, err = os.Create(path + filename)
	if err != nil {
		log.Fatal("Create: ", err)
	}

	_, _ = file.WriteString(body)

	defer file.Close()
}

func loadSyncService() {
	c := conf.GetConfig()
	_, err := utils.CommandSubSplitOutput("launchctl", "load", c.LaunchdPath+"/tokaido.sync."+c.Project+".plist")
	if err != nil {
		log.Fatal("Unable to load sync service: ", err)
	}
}

func unloadSyncService() {
	c := conf.GetConfig()
	_, err := utils.CommandSubSplitOutput("launchctl", "unload", c.LaunchdPath+"/tokaido.sync."+c.Project+".plist")
	if err != nil {
		log.Fatal("Unable to unload sync service: ", err)
	}
}

func startSyncService() {
	c := conf.GetConfig()
	_, err := utils.CommandSubSplitOutput("launchctl", "start", "tokaido.sync."+c.Project+".plist")
	if err != nil {
		log.Fatal("Unable to start the sync service: ", err)
	}
}

func stopSyncService() {
	c := conf.GetConfig()
	_, err := utils.CommandSubSplitOutput("launchctl", "stop", "tokaido.sync."+c.Project+".plist")
	if err != nil {
		log.Fatal("Unable to stop the sync service: ", err)
	}
}

func deleteSyncService() {
	c := conf.GetConfig()
	os.Remove(c.LaunchdPath + "/tokaido.sync." + c.Project + ".plist")
}

// RegisterLaunchdService Register the unison sync service for launchd
func RegisterLaunchdService() {
	createSyncFile()
	loadSyncService()
}

// StartLaunchdService Start the launchd service after it is created
func StartLaunchdService() {
	startSyncService()
	CheckSyncService()
}

// CheckSyncService checks if the unison background process is running
func CheckSyncService() string {
	c := conf.GetConfig()

	u, uErr := user.Current()
	if uErr != nil {
		log.Fatal(uErr)
	}

	o, _ := utils.CommandSubSplitOutput("launchctl", "print", "gui/"+u.Uid+"/tokaido.sync."+c.Project+".plist")

	if c.Debug == true {
		fmt.Printf("\033[33m%s\033[0m\n", o)
	}

	if strings.Contains(o, "state = running") == true {
		return "running"
	}

	return "stopped"
}

// StopLaunchdService ...
func StopLaunchdService() {
	c := conf.GetConfig()
	ps, _ := utils.CommandSubSplitOutput("launchctl", "list", "tokaido.sync."+c.Project+".plist")
	if ps != "" {
		fmt.Println(`
🔄  Removing the background sync process
	`)
		stopSyncService()
		unloadSyncService()
		deleteSyncService()
	}
}