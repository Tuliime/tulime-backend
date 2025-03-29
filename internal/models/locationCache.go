package models

import (
	"encoding/json"
	"log"
	"time"
)

func (l *Location) WriteToCache(location Location, ip string) error {

	locationJson, err := json.Marshal(&location)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 1 * time.Hour

	if err = redisClient.Set(ctx, ip, locationJson, expiration).Err(); err != nil {
		log.Println("Error saving location to Redis:", err)
		return err
	}

	return nil
}

func (l *Location) ReadFromCache(ip string) (Location, error) {
	location := Location{}

	locationString, err := redisClient.Get(ctx, ip).Result()
	if err != nil {
		log.Println("Error fetching data from Redis:", err)
		return location, nil
	}

	err = json.Unmarshal([]byte(locationString), &location)
	if err != nil {
		log.Println("Error un-marshalling JSON:", err)
		return location, nil
	}

	return location, nil
}

func (l *Location) DeleteFromCache(ip string) error {
	err := redisClient.Del(ctx, ip).Err()
	if err != nil {
		log.Println("Error deleting location from Redis:", err)
		return err
	}

	return nil
}
