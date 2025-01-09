package packages

import (
	"github.com/Tuliime/tulime-backend/internal/models"
)

// HasImagePath is temporary function to check if an agroproduct
// has imagePath value
func HasImagePath(agroProduct models.Agroproduct) bool {
	return agroProduct.ImagePath != ""
}
