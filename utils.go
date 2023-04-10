package main

import(
	"io"
	"os"
	"log"
	"errors"
	"runtime"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	
	"fyne.io/fyne/v2/storage"
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
					log.Fatal("App Support: log file - " + dirErr.Error())
				}
			}
		case "darwin":			
			saveDir = filepath.Join(filepath.Join(filepath.Join(homeDir, "Library"), "Application Support"), APPID)
			_, statErr := os.Stat(saveDir)
			if os.IsNotExist(statErr) {
				dirErr := os.Mkdir(saveDir, 0755)
				if dirErr != nil {
					log.Fatal("App Support: " + dirErr.Error())
				}
			}
		case "linux":
			saveDir = "/var/lib/" + APPID
			_, statErr := os.Stat(saveDir)
			if os.IsNotExist(statErr) {
				dirErr := os.Mkdir(saveDir, 0755)
				if dirErr != nil {
					log.Fatal("App Support: log file - " + dirErr.Error())
				}
			}
		case "android":
			saveDir = filepath.Join(homeDir, "Podcasts")
			_, statErr := os.Stat(saveDir)
			if os.IsNotExist(statErr) {
				dirErr := os.Mkdir(saveDir, 0755)
				if dirErr != nil {
					log.Fatal("App Support: log file - " + dirErr.Error())
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
	filename, err := filenamify.Filenamify(title, filenamify.Options{
    	Replacement:" -",
    })
	if err != nil {
		log.Println("utils.go - Cannot make safe: " + title)
	}
	return filename
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func downloadToFolder(url, folder, name string) {
	path := getAppSupportFolder() + "/" + folder
	createDirIfNotExist(path)
	
	newFile := path + "/" + name + ".mp3"
	out, createErr := os.Create(newFile)
	if createErr != nil {
		log.Fatal(createErr)
	}
	defer out.Close()
	
	
	resp, downloadErr := http.Get(url)
	if downloadErr != nil {
		log.Fatal(downloadErr)
	}
	defer resp.Body.Close()
	
	_, copyErr := io.Copy(out, resp.Body)
	if copyErr != nil {
		log.Fatal(copyErr)
	}
}

func deleteFile(folder, name string) {
	path := getAppSupportFolder() + "/" + folder + "/" + name + ".mp3"
	deleteErr := os.Remove(path)
	
	if deleteErr != nil {
		log.Fatal(deleteErr)
	}
}

func loadShowsFromJSON() []*Show {
	path := filepath.Join(getAppSupportFolder(), "/DATA.sun")
	
	var shows []*Show
	var data []byte
	
	if runtime.GOOS != "android" {
		data = loadJSON(path)
	} else {
		data = loadAndroidJSON()
		data = []byte("[{\"Type\":\"Omny\",\"Name\":\"Behind The Bastards\",\"Slug\":\"behind-the-bastards\"}]")
	}
	
	jsonErr := json.Unmarshal(data, &shows)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	
	return shows
	
}

func loadJSON(path string) []byte {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
  		os.Create(path)
  		ioutil.WriteFile(path, []byte("[{\"Type\":\"Omny\",\"Name\":\"Behind The Bastards\",\"Slug\":\"behind-the-bastards\"}]"), 0600)
	}
	
	data, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return data
}

func loadAndroidJSON() []byte {
	path := fyneApp.Storage().RootURI().Path() + "/DATA.sun"
	uri, _ := storage.ParseURI(path)
	
	_, readErr := storage.Reader(uri)
	
	if readErr != nil {
		log.Fatal(readErr, " " + path)
	}
	
	return []byte("")
}