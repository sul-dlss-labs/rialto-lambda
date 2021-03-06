package derivative

import "github.com/sul-dlss/rialto-derivatives/models"

// Writer writes a derivative document
type Writer interface {
	RemoveAll() error
	Add(docs []models.Resource) error
}
