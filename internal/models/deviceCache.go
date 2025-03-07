package models

import (
	"encoding/json"
	"fmt"
	"log"
)

var ALL_USER_DEVICE_PATTERN = "user-device:*"
var USER_DEVICE_PATTERN = "user-device:"

func (d *Device) writeToCache(devices []Device) error {
	if len(devices) == 0 {
		return nil
	}
	devicesJson, err := json.Marshal(&devices)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	userDeviceKey := d.buildUserDeviceKey(devices[0].UserID)
	log.Printf("userDeviceKey: %v", userDeviceKey)

	if err = redisClient.Set(ctx, userDeviceKey, devicesJson, 0).Err(); err != nil {
		log.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (d *Device) readFromCache(userID string) ([]Device, error) {
	devices := []Device{}
	userDeviceKey := d.buildUserDeviceKey(userID)

	devicesString, err := redisClient.Get(ctx, userDeviceKey).Result()
	if err != nil {
		log.Println("Error fetching data from Redis:", err)
		return devices, nil
	}

	err = json.Unmarshal([]byte(devicesString), &devices)
	if err != nil {
		log.Println("Error un-marshalling JSON:", err)
		return devices, nil
	}

	return devices, nil
}

func (d *Device) readAllFromCache() ([]Device, error) {
	devices := []Device{}
	var cursor uint64
	var keys []string
	var err error

	// Scan for keys matching a pattern
	for {
		var batch []string
		batch, cursor, err = redisClient.Scan(ctx, cursor, ALL_USER_DEVICE_PATTERN, 10).Result()
		if err != nil {
			log.Printf("Error scanning keys: %v", err)
			return devices, nil
		}

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return devices, nil
	}

	// Fetch values for the found keys
	devicesData, err := redisClient.MGet(ctx, keys...).Result()
	if err != nil {
		log.Printf("Error getting values: %v", err)
		return devices, nil
	}

	for _, device := range devicesData {
		deviceData, ok := device.(Device)
		if !ok {
			log.Printf("Invalid Device type received: %T", device)
			return devices, nil
		}
		devices = append(devices, deviceData)
	}

	return devices, nil
}

func (d *Device) deleteFromCache(userID string) error {
	userDeviceKey := d.buildUserDeviceKey(userID)

	err := redisClient.Del(ctx, userDeviceKey).Err()
	if err != nil {
		log.Println("Error deleting data from Redis:", err)
		return err
	}

	return nil
}

func (d *Device) buildUserDeviceKey(userID string) string {
	return fmt.Sprintf("%s%s", USER_DEVICE_PATTERN, userID)
}
