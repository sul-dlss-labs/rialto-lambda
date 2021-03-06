package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestPublicationResourceToDoc(t *testing.T) {
	indexer := &PublicationIndexer{}
	doi := "https://doi.org/10.1109/5.771073"
	abstract := "Some very compelling research"
	description := "Super interesting"
	publisher := "http://example.com/publishing_house"
	author1 := &models.Author{URI: "http://example.com/person1", Label: "Harry Potter"}
	author2 := &models.Author{URI: "http://example.com/person2", Label: "Hermione Granger"}
	concept := &models.Concept{URI: "http://example.com/concept1", Label: "Magic"}
	grant := &models.Labeled{URI: "http://example.com/grant1", Label: "Learning Magic"}
	identifier := "123456"

	authors := []*models.Author{author1, author2}
	concepts := []*models.Concept{concept}
	grants := []*models.Labeled{grant}
	identifiers := []string{identifier}
	resource := &models.Publication{
		DOI:         &doi,
		Title:       "Hello world!",
		Created:     "2004-06-11?",
		Identifiers: identifiers,
		Abstract:    &abstract,
		Publisher:   &publisher,
		Description: &description,
		Authors:     authors,
		Concepts:    concepts,
		Grants:      grants,
		CreatedYear: 2004,
	}
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, "Hello world!", doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
	assert.Equal(t, "2004-06-11?", doc.Get("created_ssim"))
	assert.Equal(t, []string{"123456"}, doc.Get("identifiers_ssim"))
	assert.Equal(t, []string{"http://example.com/person1", "http://example.com/person2"}, doc.Get("authors_ssim"))
	assert.Equal(t, []string{"Harry Potter", "Hermione Granger"}, doc.Get("author_labels_tsim"))
	assert.Equal(t, []string{"http://example.com/concept1"}, doc.Get("concepts_ssim"))
	assert.Equal(t, []string{"Magic"}, doc.Get("concept_labels_ssim"))
	assert.Equal(t, []string{"http://example.com/grant1"}, doc.Get("grants_ssim"))
	assert.Equal(t, []string{"Learning Magic"}, doc.Get("grant_labels_ssim"))

	assert.Equal(t, doi, doc.Get("doi_ssim"))
	assert.Equal(t, abstract, doc.Get("abstract_tesim"))
	assert.Equal(t, publisher, doc.Get("publisher_ssim"))
	assert.Equal(t, description, doc.Get("description_tesim"))
	assert.Equal(t, 2004, doc.Get("created_year_isim"))
}
