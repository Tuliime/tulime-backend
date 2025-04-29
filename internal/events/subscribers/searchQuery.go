package subscribers

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/handlers/location"
	"github.com/Tuliime/tulime-backend/internal/handlers/search"
	"github.com/Tuliime/tulime-backend/internal/models"
)

func saveSearchQuery(sQuery search.SearchQueryEvent) {
	searchQuery := models.SearchQuery{}
	user := models.User{}

	device := sQuery.Device
	userID := sQuery.UserID

	if userID == "" {
		anonymousUser, err := user.FindByTelNumber(constants.AnonymousTelNumber)
		if err != nil {
			log.Printf("Error getting Anonymous user :%+v", err)
		}
		userID = anonymousUser.ID
		log.Println("Using Anonymous userID : ", userID)
	}

	clientIP := sQuery.ClientIP

	location, err := location.GetUserLocationByIP(userID, clientIP)
	if err != nil {
		log.Printf("Error getting location :%+v", err)
	}

	if _, err := searchQuery.Create(models.SearchQuery{
		Query: sQuery.Query, UserID: userID, LocationID: location.ID,
		Device: device}); err != nil {
		log.Printf("Error saving SearchQuery :%+v", err)
	}
	log.Println("SearchQuery saved")
}

func SearchQueryEventListener() {
	SearchQueryChan := make(chan events.DataEvent)
	events.EB.Subscribe("searchQuery", SearchQueryChan)

	for {
		searchQueryEvent := <-SearchQueryChan
		searchQuery, ok := searchQueryEvent.Data.(search.SearchQueryEvent)
		if !ok {
			log.Printf("Invalid searchQuery type received: %+v", searchQueryEvent.Data)
			return
		}
		log.Printf("saving search query...")
		saveSearchQuery(searchQuery)
	}
}
