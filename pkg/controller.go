package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (ctrl *FileController) GetFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	log.Infof("GetFile, filepath: %s", filepath)
	file, err := ctrl.Storage.ReadFile(filepath)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprint(w, "Error during file opening, %s", err.Error())
		return
	}
	defer file.Close()

	/*
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Error during writing to file")
		return
	}
	 */

	content, err := ioutil.ReadAll(file)
	fileNames := url.QueryEscape(filepath) // In case of Chinese wrong code
	w.Header().Add("Content-Type", "application/octet-stream")
	segs := strings.Split(fileNames, "/")
	downloadFileName := segs[len(segs)-1]
	w.Header().Add("Content-Disposition", "attachment; filename=\""+downloadFileName+"\"")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("server error: %s", err)
		return
	}

	w.Write(content)
	log.Infof("file %s downloaded!", downloadFileName)
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
