package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type FileController struct {
	Storage Storage
}

func NewFileController(storage Storage) *FileController {
	return &FileController{
		Storage: storage,
	}
}

func check(s Storage, filepath string) (bool, error) {
	// Principle:
	// 1. single file will be considered if there is no key of this filepath ?
	// 2. single file with key of this filepath but size 0 at the meanwhile ?
	// 3. single file with key of this filepath but size 0 and no suffix of .7z ?
	// 4. single file end up with /

	/*
	if strings.EqualFold(s.GetObjectKey(filepath), "") {
		return false, errors.New("Key not found")
	}
	 */

	// it's directory
	/*
		if  0 == s.GetObjectSize(filepath) && 1 == len(strings.Split(filepath, ".")) {
			return true, nil
		}
	*/
	// using principle 4
	if strings.HasSuffix(filepath, "/") {
		return true, nil
	}

	// TODO: put directory into object metadata
	return false, nil
}

func (ctrl *FileController) GetFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	if strings.EqualFold(filepath, "/") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "liveness get %d goroutines", runtime.NumGoroutine())
		return
	}

	filepath = strings.TrimPrefix(filepath, "/")

	// get or list will be determinded by filepath
	directory, err := check(ctrl.Storage, filepath)
	log.Infof("directory=%v, err=%s", directory, err)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error during file opening, %s", err.Error())
		log.Errorf("check error: %s", err)
		return
	}

	if directory {
		log.Infof("GetFile, list directory filepath: %s", filepath)
		objects, err := ctrl.Storage.ListDirectory(filepath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error during directory opening , %s", err.Error())
			log.Errorf("list directory error: %s", err)
			return
		}

		data, err := json.Marshal(objects)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error during directory opening , %s", err.Error())
			log.Errorf("list directory json marshal error: %s", err)
			return
		}

		w.Write(data)
		log.Infof("directory %s listed! data=%s", filepath, string(data))
	} else {
		log.Infof("GetFile, read file filepath: %s", filepath)
		file, err := ctrl.Storage.ReadFile(filepath)
		defer file.Close()

		if err != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Error during file opening, %s", err.Error())
			log.Errorf("read file error: %s", err)
			return
		}

		content, err := ioutil.ReadAll(file)
		fileNames := url.QueryEscape(filepath) // In case of Chinese wrong code
		w.Header().Add("Content-Type", "application/octet-stream")
		segs := strings.Split(fileNames, "/")
		downloadFileName := segs[len(segs)-1]
		w.Header().Add("Content-Disposition", "attachment; filename=\""+downloadFileName+"\"")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("read file server error: %s", err)
			return
		}

		w.Write(content)
		log.Infof("file %s downloaded!", downloadFileName)
	}

	/*
		_, err = io.Copy(w, file)
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Error during writing to file")
			return
		}
	*/
}

func (ctrl *FileController) PutFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	log.Infof("PutFile, filepath: %s\n", filepath)
	err := ctrl.Storage.WriteFile(filepath, r.Body)
	if err != nil {
		log.Errorf("Error during file saving: %s", err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.WriteHeader(200)
}

func (ctrl *FileController) DeleteFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	filepath = strings.TrimPrefix(filepath, "/")
	log.Infof("DeleteFile, filepath: %s\n", filepath)
	err := ctrl.Storage.DeleteFile(filepath)
	if err != nil {
		log.Errorf("Error during file deleting: %s", err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.WriteHeader(200)
}

type CopyParam struct {
	Src string `json:"src"`
	TargetRepo string `json:"targetRepo"`
	Dist string `json:"dist"`
	Recursive bool `json:"recursive"`
}

func (ctrl *FileController) CopyFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	log.Infof("CopyFile, filepath: %s\n", filepath)
	filepath = strings.TrimPrefix(filepath, "/")

	param := &CopyParam{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error during file copying: %s", err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	err = json.Unmarshal(data, &param)
	if err != nil {
		log.Errorf("Error during file copying: %s", err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	log.Infof("copy params: %+v", param)

	param.Src = filepath

	err = ctrl.Storage.Copy(param.TargetBucket, param.Src, param.Dist, param.Recursive)
	if err != nil {
		log.Errorf("Error during file copying: %s", err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.WriteHeader(200)
}
