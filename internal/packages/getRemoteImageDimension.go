package packages

import (
	"fmt"
	"io"
	"net/http"

	"github.com/h2non/bimg"
)

// GetRemoteImageDimensions returns integers
// of width and height sizes of the fetched image
// respectively
func GetRemoteImageDimensions(url string) (int, int, error) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch image, status code:", resp.StatusCode)
		return 0, 0, err

	}
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading image data:", err)
		return 0, 0, err

	}

	image := bimg.NewImage(imageData)
	metadata, err := image.Metadata()
	if err != nil {
		fmt.Println("Error getting image metadata:", err)
		return 0, 0, err

	}

	fmt.Printf("Image Width: %dpx\n", metadata.Size.Width)
	fmt.Printf("Image Height: %dpx\n", metadata.Size.Height)

	return metadata.Size.Width, metadata.Size.Height, nil
}
