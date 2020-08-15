package server

import (
	"encoding/json"
	"fmt"
	"image_server/config"
	"image_server/model"
	"image_server/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var album model.CreateAlbum

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse("CreateAlbum", w, http.StatusBadRequest, "Bad Request: empty body", fmt.Sprintf("Failed to read body in request: %v", err))
		return
	}

	err = json.Unmarshal(body, &album)
	if err != nil {
		sendErrorResponse("CreateAlbum", w, http.StatusBadRequest, "Bad Request: invalid JSON", fmt.Sprintf("unmarshaling failed for request body: %v", err))
		return
	}

	file_path := "data/" + album.Name
	if isFileExist(file_path) {
		sendErrorResponse("CreateAlbum", w, http.StatusConflict, "conflict: album already exists", fmt.Sprintf("album already exists: %s", album.Name))
		return
	}

	os.Mkdir(file_path, 0777)
	sendSuccessResponse("CreateAlbum", w, http.StatusOK, "success", "") //don't log for GETs
}

func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	album_name := vars["album_name"]

	if len(album_name) == 0 {
		sendErrorResponse("DeleteAlbum", w, http.StatusBadRequest, "Bad Request: empty album name", fmt.Sprintf("empty album name"))
		return
	}

	file_path := "data/" + album_name
	if !isFileExist(file_path) {
		sendErrorResponse("DeleteAlbum", w, http.StatusNotFound, "album not found: "+album_name, "album not found: "+album_name)
		return
	}

	os.RemoveAll(file_path)
	os.Remove(file_path)
	sendSuccessResponse("DeleteAlbum", w, http.StatusOK, "success", "") //don't log for GETs
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	album_name := vars["album_name"]

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	image_name := handler.Filename
	dir_path := "./data/" + album_name
	file_path := "./data/" + album_name + "/" + handler.Filename

	if !isFileExist(dir_path) {
		sendErrorResponse("UploadImage", w, http.StatusNotFound, "album doesn't exist: "+album_name, "album doesn't exist: "+album_name)
		return
	}

	if isFileExist(file_path) {
		sendErrorResponse("UploadImage", w, http.StatusConflict, "image already exists in album: "+album_name+", image: "+image_name, "image already exists in album: "+album_name+", image: "+image_name)
		return
	}

	f, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
		sendErrorResponse("UploadImage", w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	defer f.Close()
	io.Copy(f, file)

	fmt.Fprintf(w, "%v", handler.Header)

	kafka_event := model.KafkaEvent{Event: "create_image", Album: album_name, Name: image_name, Operation: "create"}
	event_msg, _ := json.Marshal(kafka_event)
	utils.PublishKafkaEvent(string(event_msg)) //publish kafka message for image creation
}

func GetAllImages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	album_name := vars["album_name"]
	file_path := "data/" + album_name

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if !isFileExist(file_path) {
		sendErrorResponse("GetAllImages", w, http.StatusNotFound, "album doesn't exist: "+album_name, "album doesn't exist: "+album_name)
		return
	}

	list := make([]model.AllImagesList, 0)

	files, err := ioutil.ReadDir(file_path)
	if err != nil {
		sendErrorResponse("GetAllImages", w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	// loop through all the files inside a dir
	for _, file := range files {
		image := file.Name()
		list = append(list, model.AllImagesList{Image: image, URL: fmt.Sprintf("http://%s/albums/%s/images/%s", config.Config["host"], album_name, image)})
	}

	resp, _ := json.Marshal(list)
	w.Write(resp)
}

func DownloadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	album_name, image_name := vars["album_name"], vars["image_name"]
	file_path := "data/" + album_name + "/" + image_name

	if !isFileExist(file_path) {
		sendErrorResponse("DownloadImage", w, http.StatusNotFound, "image doesn't exist in album: "+album_name+", image: "+image_name, "image doesn't exist in album: "+album_name+", image: "+image_name)
		return
	}

	img, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img) //serve image
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	album_name, image_name := vars["album_name"], vars["image_name"]
	file_path := "data/" + album_name + "/" + image_name

	if !isFileExist(file_path) {
		sendErrorResponse("DeleteImage", w, http.StatusNotFound, "image doesn't exist in album: "+album_name+", image: "+image_name, "image doesn't exist in album: "+album_name+", image: "+image_name)
		return
	}

	os.Remove(file_path) //delete image
	w.WriteHeader(http.StatusOK)

	kafka_event := model.KafkaEvent{Event: "delete_image", Album: album_name, Name: image_name, Operation: "delete"}
	event_msg, _ := json.Marshal(kafka_event)
	utils.PublishKafkaEvent(string(event_msg)) //publish kafka message for image deletion

	sendSuccessResponse("DeleteImage", w, http.StatusOK, "success", "")
}

func sendSuccessResponse(func_name string, w http.ResponseWriter, http_code int, resp_body string, log_message string) {
	if len(log_message) > 0 {
		utils.Log(fmt.Sprintf("%s: %s", func_name, log_message))
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http_code)
	w.Write([]byte(resp_body))
	return
}

// logs and sends error response
func sendErrorResponse(func_name string, w http.ResponseWriter, http_code int, resp_body string, log_message string) {
	utils.Log(fmt.Sprintf("%s: %s", func_name, log_message))
	w.WriteHeader(http_code)
	w.Write([]byte(resp_body))
	return
}

// checks if file or directory exists
func isFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
