package packages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type FirebaseStorage struct {
	StorageBucketBucket string `json:"storageBucket"`
	BaseURL             string `json:"baseURL"`
	DownloadURL         string `json:"url"`
	FilePath            string `json:"path"`
}

var firebaseManager = NewFirebaseManager("serviceAccountKey.json")

func (fs *FirebaseStorage) initStorageBucket() (*storage.BucketHandle, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	currentDirPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fs.BaseURL = "https://firebasestorage.googleapis.com/v0/b/"
	storageBucket := os.Getenv("FIREBASE_STORAGE_BUCKET")
	fs.StorageBucketBucket = storageBucket

	configStorage := &firebase.Config{
		StorageBucket: storageBucket,
	}

	if err := firebaseManager.CreateFile(); err != nil {
		firebaseManager.DeleteFile()
		log.Println("Error creating file:", err)
	}

	opt := option.WithCredentialsFile(currentDirPath + "/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), configStorage, opt)
	if err != nil {
		firebaseManager.DeleteFile()
		return nil, err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		firebaseManager.DeleteFile()
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		firebaseManager.DeleteFile()
		return nil, err
	}

	return bucket, nil

}

func (fs *FirebaseStorage) Add(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	filePath := fs.FilePath
	if filePath == "" {
		firebaseManager.DeleteFile()
		return "", errors.New("no file path provided")
	}

	bucket, err := fs.initStorageBucket()
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	wc := bucket.Object(filePath).NewWriter(context.Background())
	_, err = io.Copy(wc, file)
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	err = wc.Close()
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	url, err := fs.getDownloadURL()
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	return url, nil
}

func (fs *FirebaseStorage) Update(file multipart.File, fileHeader *multipart.FileHeader, savedFilePath string) (string, error) {

	filePath := fs.FilePath
	if filePath == "" {
		return "", errors.New("no file path provided")
	}

	bucket, err := fs.initStorageBucket()
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	wc := bucket.Object(filePath).NewWriter(context.Background())
	_, err = io.Copy(wc, file)
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	err = wc.Close()
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	url, err := fs.getDownloadURL()
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	if url != "" {
		if err := fs.Delete(savedFilePath); err != nil {
			firebaseManager.DeleteFile()
			return "", err
		}

		fmt.Println("file deleted from storage using path ==>", savedFilePath)
	}

	return url, nil
}

// TODO: to debug deletion functionality
func (fs *FirebaseStorage) Delete(filePath string) error {

	if filePath == "" {
		firebaseManager.DeleteFile()
		return errors.New("no file path provided")
	}

	bucket, err := fs.initStorageBucket()
	if err != nil {
		firebaseManager.DeleteFile()
		return err
	}
	obj := bucket.Object(filePath)

	if err := obj.Delete(context.Background()); err != nil {
		firebaseManager.DeleteFile()
		return err
	}

	return nil
}

func (fs *FirebaseStorage) transformFilePath() (string, error) {
	path := fs.FilePath

	if path == "" {
		return "", errors.New("no file path provided")
	}

	path = strings.ReplaceAll(path, "/", "%2F")
	path = strings.ReplaceAll(path, " ", "%20")
	path = strings.ReplaceAll(path, "?", "%3F")
	path = strings.ReplaceAll(path, "&", "%26")
	path = strings.ReplaceAll(path, "=", "%3D")
	path = strings.ReplaceAll(path, ":", "%3A")
	path = strings.ReplaceAll(path, ",", "%2C")

	return path, nil
}

func (fs *FirebaseStorage) getDownloadURL() (string, error) {
	start := time.Now()

	transformedFilePath, err := fs.transformFilePath()
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	FIREBASE_STORAGE_URL := fs.BaseURL + fs.StorageBucketBucket + "/o/" + transformedFilePath

	req, err := http.NewRequest(http.MethodGet, FIREBASE_STORAGE_URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		firebaseManager.DeleteFile()
		return "", err
	}

	type Response struct {
		Name               string `json:"name"`
		Bucket             string `json:"bucket"`
		Generation         string `json:"generation"`
		Metageneration     string `json:"metageneration"`
		ContentType        string `json:"contentType"`
		TimeCreated        string `json:"timeCreated"`
		Updated            string `json:"updated"`
		StorageClass       string `json:"storageClass"`
		Size               string `json:"size"`
		Md5Hash            string `json:"md5Hash"`
		ContentEncoding    string `json:"contentEncoding"`
		ContentDisposition string `json:"contentDisposition"`
		Crc32c             string `json:"crc32c"`
		Etag               string `json:"etag"`
		DownloadTokens     string `json:"downloadTokens"`
	}

	if res.StatusCode != http.StatusOK {
		_, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New("request to firebase storage failed")
	}

	fmt.Printf("firebase-storage-service: status code: %d\n", res.StatusCode)
	rBody, _ := io.ReadAll(res.Body)

	response := Response{}
	json.NewDecoder(strings.NewReader(string(rBody))).Decode(&response)

	downloadURL := FIREBASE_STORAGE_URL + "?alt=media&token=" + response.DownloadTokens
	fs.DownloadURL = downloadURL

	fmt.Printf(
		"Firebase storage request duration : %s\n",
		time.Since(start),
	)
	firebaseManager.DeleteFile()

	return downloadURL, nil
}
