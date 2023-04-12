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

//File Download

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

func downloadPodcast(url, folder, name string) {
	if runtime.GOOS != "android" {
		desktopDownload(url, folder, name)
	} else {
		androidDownload(url, folder, name)
	}
}

func desktopDownload(url, folder, name string) {
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

func androidDownload(url, folder, name string) {
	folderPath := "/sdcard/Podcasts" + "/" + folder
	filePath := folderPath + "/" + name
	
	//Folder
	createDirIfNotExist(folderPath)
	
	//File
	file, createErr := os.Create(filePath)
	if createErr != nil {
		log.Fatal(createErr)
	}
	defer file.Close()
	
	/*
	log.Println("Folder")
	//Check if Folder Exists
	folderURI, folderErr := storage.ParseURI("file://" + folderPath)
	if folderErr != nil {
		log.Fatal(folderErr)
	}
	folderExists, existsErr := storage.Exists(folderURI)
	if existsErr != nil {
		log.Fatal(existsErr)
	}
	if !folderExists {
		log.Println("Folder Does Not Exist At: " + folderURI.Path())
		folderCreateErr := storage.CreateListable(folderURI)
		if folderCreateErr != nil {
			log.Fatal(folderCreateErr)
		}
		err := os.Mkdir("/sdcard/Podcasts", os.ModePerm)
		log.Fatal(err)
	}
	
	log.Println("File URI")
	//Make File URI
	uri, uriErr := storage.ParseURI("file://" + filePath)
	if uriErr != nil {
		log.Fatal(uriErr)
	}
	canWrite, canWriteErr := storage.CanWrite(uri)
	if canWriteErr != nil {
		log.Fatal(canWriteErr)
	}
	if !canWrite {
		log.Fatal("Can not write to this location: " + uri.Path())
	}
	*/
	
	log.Println("Download")
	//Download
	resp, downloadErr := http.Get(url)
	if downloadErr != nil {
		log.Fatal(downloadErr)
	}
	defer resp.Body.Close()
	
	/*
	log.Println("Writer")
	//Make Writer
	writer, writeErr := storage.Writer(uri)
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	
	log.Println("Data")
	//Get Data & Save
	data, dataErr := ioutil.ReadAll(resp.Body)
	if dataErr != nil {
		log.Fatal(dataErr)
	}
	log.Println("Save")
	_, saveErr := writer.Write(data)
	if writeErr != nil {
		log.Fatal(saveErr)
	}
	*/
	
	_, copyErr := io.Copy(file, resp.Body)
	if copyErr != nil {
		log.Fatal(copyErr)
	}
}

func deleteFile(folder, name string) {
	if runtime.GOOS != "android" {
		desktopDelete(folder, name)
	} else {
		androidDelete(folder, name)
	}
}

func desktopDelete(folder, name string) {
	path := getAppSupportFolder() + "/" + folder + "/" + name + ".mp3"
	deleteErr := os.Remove(path)
	
	if deleteErr != nil {
		log.Fatal(deleteErr)
	}
}

func androidDelete(folder, name string) {

}

//Load Shows

func loadShowsFromJSON() []*Show {
	path := filepath.Join(getAppSupportFolder(), "/DATA.sun")
	
	var shows []*Show
	var data []byte
	
	if runtime.GOOS != "android" {
		data = loadJSON(path)
	} else {
		data = loadAndroidJSON()
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
	//Fyne has a special way of handling android files... 
	
	//Get APPDATA path
	path := fyneApp.Storage().RootURI().Path() + "/DATA.sun"
	
	//Create URI
	uri, uriErr := storage.ParseURI("file://" + path)
	if uriErr != nil {
		log.Fatal(uriErr)
	}
	
	//Create Fyne Reader
	r, readerErr := storage.Reader(uri)	
	if readerErr != nil {
		log.Fatal(readerErr, " " + path)
	}
	
	//Read data using regular io package
	data, readErr := io.ReadAll(r)
	if readErr != nil {
		log.Fatal(readErr)
	}
	
	//Return read data
	return data
}