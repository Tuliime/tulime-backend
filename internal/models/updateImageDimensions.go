package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Tuliime/tulime-backend/internal/packages"
)

func UpdateImageDimensions() {
	log.Println("====== Start update chatroomFile Dimensions ======")
	chatroomFile := ChatroomFile{}

	files, err := chatroomFile.FindAll(1000)
	if err != nil {
		log.Println("Error fetching chatroom files:", err)
		return
	}

	for _, file := range files {
		if file.Dimensions != nil {
			continue
		}
		width, height, err := packages.GetRemoteImageDimensions(file.URL)
		if err != nil {
			log.Println("Error getting image dimensions:", err)
			return
		}

		dimensions := struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		}{
			Width:  width,
			Height: height,
		}

		dimensionJson, err := json.Marshal(dimensions)
		if err != nil {
			log.Println("Error Marshalling dimensions:", err)
			return
		}

		file.Dimensions = JSONB(dimensionJson)

		updatedChatroomFile, err := file.Update()
		if err != nil {
			log.Println("Error updating chatroomFile ", err)
			return
		}
		log.Println("updatedChatroomFile: ", updatedChatroomFile)
	}
	log.Println("====== End update chatroomFile Dimensions ======")
}

func init() {
	// Schedule UpdateImageDimensions to run 10 seconds after the server starts
	time.AfterFunc(10*time.Second, UpdateImageDimensions)
}
