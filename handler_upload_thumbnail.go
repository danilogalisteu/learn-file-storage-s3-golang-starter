package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	dbVideo, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get video", err)
		return
	}
	if dbVideo.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "You are not allowed to upload a thumbnail for this video", err)
		return
	}

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	// "thumbnail" should match the HTML form input name
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	fileType := header.Header.Get("Content-Type")
	if fileType == "" {
		respondWithError(w, http.StatusBadRequest, "Unable to get file type", err)
		return
	}

	fileData := make([]byte, header.Size)
	_, err = file.Read(fileData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to read file", err)
		return
	}

	fileString := fmt.Sprintf("data:%s;base64,%s", fileType, base64.StdEncoding.EncodeToString(fileData))
	dbVideo.ThumbnailURL = &fileString
	dbVideo.UpdatedAt = time.Now()
	err = cfg.db.UpdateVideo(dbVideo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update video", err)
		return
	}
	fmt.Println(dbVideo)

	respondWithJSON(w, http.StatusOK, dbVideo)
}
