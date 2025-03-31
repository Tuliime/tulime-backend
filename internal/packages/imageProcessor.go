package packages

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/h2non/bimg"
)

// ImageProcessor handles image compression and conversion
type ImageProcessor struct{}

func (ip *ImageProcessor) CompressMultipartFile(file multipart.File, quality int) ([]byte, error) {
	imgData, err := ip.FileToBytes(file)
	if err != nil {
		return nil, err
	}

	imgBuffer, err := bimg.NewImage(imgData).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}

	return ip.Compress(imgBuffer, quality)
}

// CompressToLessThanOneKB compresses an image to ensure its size is below 1KB
func (ip *ImageProcessor) CompressToLessThanOneKB(file multipart.File) ([]byte, error) {
	imgData, err := ip.FileToBytes(file)
	if err != nil {
		return nil, err
	}

	// Convert image to JPEG using bimg
	imgData, err = bimg.NewImage(imgData).Convert(bimg.JPEG)
	if err != nil {
		return nil, err
	}

	quality := 100
	for {
		// Compress using standard Go image/jpeg package
		compressedData, err := ip.compressJPEG(imgData, quality)
		if err != nil {
			return nil, err
		}
		log.Println("image size in kb: ", len(compressedData)/1024)
		log.Println("quality: ", quality)
		if len(compressedData) < 1024 {
			log.Println("inside threshold image size in kb: ", len(compressedData)/1024)
			return compressedData, nil
		}

		if quality > 20 {
			quality -= 10
		} else if quality > 10 {
			quality -= 5
		} else if quality > 5 {
			quality -= 2
		} else if quality >= 1 {
			log.Println("Quality is:", quality)
			quality -= 1
		} else {
			quality -= 1
			log.Println("breaking the loop...")
			return compressedData, nil
		}
	}
}

// compressJPEG compresses image data using Go's standard image/jpeg package
func (ip *ImageProcessor) compressJPEG(imgData []byte, quality int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	opt := jpeg.Options{Quality: quality}
	if err := jpeg.Encode(&buf, img, &opt); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Compress reduces image size by adjusting
// quality and stripping metadata
func (ip *ImageProcessor) Compress(imgData []byte, quality int) ([]byte, error) {
	startTime := time.Now()
	compressedBuf, err := bimg.NewImage(imgData).Process(bimg.Options{
		Quality:       quality,
		StripMetadata: true,
		Compression:   9,
	})
	if err != nil {
		return nil, err
	}
	log.Println("File Compression Duration:", time.Since(startTime))

	return compressedBuf, nil
}

// Convert changes image format (e.g., from PNG to JPEG)
func (ip *ImageProcessor) Convert(imgData []byte, format bimg.ImageType) ([]byte, error) {
	image := bimg.NewImage(imgData)
	return image.Convert(format)
}

// FileToBytes converts []byte to Base64 string
func (ip *ImageProcessor) FileToBytes(file multipart.File) ([]byte, error) {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// ToBase64 converts []byte to Base64 string
func (ip *ImageProcessor) ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
