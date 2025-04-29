package search

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"sync"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var Search = func(c *fiber.Ctx) error {
	news := models.News{}
	advert := models.Advert{}
	parameterListEncoding := c.Query("parameters")
	searchQuery := c.Query("query")

	if searchQuery == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Search query can't be empty!")
	}

	device := c.Get("X-Device")

	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("Invalid userID type!")
	}

	clientIP, ok := c.Locals("clientIP").(string)
	if !ok {
		log.Println("Invalid client type!")
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(parameterListEncoding)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var parameterList []string

	parameterListStr := string(decodedBytes)
	if parameterListStr == "" {
		parameterList = []string{"*"}
	}

	if parameterListStr != "" {
		err = json.Unmarshal(decodedBytes, &parameterList)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid format! encoded data must array of strings.")
		}
	}

	var searchNews, searchAdverts bool

	for _, param := range parameterList {
		if param == "*" {
			searchNews = true
			searchAdverts = true
			break
		}
	}

	// If not searching all, check for specific search types
	if !searchNews && !searchAdverts {
		for _, param := range parameterList {
			switch param {
			case "news":
				searchNews = true
			case "advert", "adverts":
				searchAdverts = true
			}
		}
	}

	type searchResults struct {
		News    []models.News   `json:"news,omitempty"`
		Adverts []models.Advert `json:"adverts,omitempty"`
		Errors  []string        `json:"errors,omitempty"`
	}

	results := searchResults{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	if searchNews {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newsResults, err := news.Search(searchQuery)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				results.Errors = append(results.Errors, "News search error: "+err.Error())
			} else {
				results.News = newsResults
			}
		}()
	}

	if searchAdverts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			advertResults, err := advert.Search(searchQuery)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				results.Errors = append(results.Errors, "Advert search error: "+err.Error())
			} else {
				results.Adverts = advertResults
			}
		}()
	}

	wg.Wait()

	if len(results.Errors) > 0 {
		log.Printf("results.Errors : %+v", results.Errors)
		return fiber.NewError(fiber.StatusInternalServerError, results.Errors[0])
	}

	events.EB.Publish("searchQuery", SearchQueryEvent{Query: searchQuery,
		UserID: userID, ClientIP: clientIP, Device: device})

	response := fiber.Map{
		"status": "success",
		"data":   results,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
