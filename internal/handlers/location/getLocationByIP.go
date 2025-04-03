package location

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Tuliime/tulime-backend/internal/models"
)

func GetUserLocationByIP(userID, ip string) (models.Location, error) {
	startTime := time.Now()
	location := models.Location{}
	savedLocation, err := location.FindByIP(ip)
	if err != nil {
		return location, err
	}
	location.ID = savedLocation.ID

	if savedLocation.ID != "" {
		log.Println("GetIPInfo Duration:", time.Since(startTime))
		log.Println("IPInfo Already exists")
		return savedLocation, nil
	}
	// ip := "197.239.8.162"
	info, err := GetIPInfo(ip)
	if err != nil {
		fmt.Printf("Error getting ip info:%+v\n", err)
		// return location, err
	}
	fmt.Printf("Found ip location:%+v\n", info)

	ipinfoJson, err := json.Marshal(&info)
	if err != nil {
		fmt.Printf("Error Marshalling ipinfo:%+v\n", err)
		// return location, err
	}
	newLocation, err := location.Create(models.Location{UserID: userID,
		Info: models.JSONB(ipinfoJson)}, ip)
	if err != nil {
		return location, err
	}

	log.Println("GetIPInfo Duration:", time.Since(startTime))
	log.Println("New IPInfo")

	return newLocation, nil

}
