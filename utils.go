package main

import(
	"os"
	"log"
	"errors"
	"runtime"
	"path/filepath"
	
	"github.com/flytam/filenamify"
)

const APPID = "com.sunshine.gopod"

func getAppSupportFolder() string {
	var saveDir string
	homeDir, _ := os.UserHomeDir()
	switch runtime.GOOS {
		case "windows":
			saveDir = filepath.Join(filepath.Join(filepath.Join(homeDir, "AppData"), "Roaming"), APPID)
			_, err := os.Stat(saveDir)
			if os.IsNotExist(err) {
				dirErr := os.Mkdir(saveDir, 0755)
				if dirErr != nil {
					panic("settings: log file - " + dirErr.Error())
				}
			}
		case "darwin":			
			saveDir = filepath.Join(filepath.Join(filepath.Join(homeDir, "Library"), "Application Support"), APPID)
			_, statErr := os.Stat(saveDir)
			if os.IsNotExist(statErr) {
				dirErr := os.Mkdir(saveDir, 0755)
				if dirErr != nil {
					panic("settings: " + dirErr.Error())
				}
			}
		case "linux":
			saveDir = "/var/lib/" + APPID
			_, statErr := os.Stat(saveDir)
			if os.IsNotExist(statErr) {
				dirErr := os.Mkdir(saveDir, 0755)
				if dirErr != nil {
					panic("settings: log file - " + dirErr.Error())
				}
			}
	}
	return saveDir
}

func checkIfDownloaded(filename string) bool {
	path := getAppSupportFolder() + "/" + filename + ".mp3"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
  		return false
	} else {
		return true
	}
}

func makeFileNameSafe(title string) string {
	filename, err := filenamify.Filenamify(title, filenamify.Options{})
	if err != nil {
		log.Println("utils.go - Cannot make safe: " + title)
	}
	return filename
}