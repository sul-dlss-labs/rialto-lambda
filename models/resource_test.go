package models

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	title, _ := rdf.NewLiteral("Hello world!")

	data[RdfTypePredicate] = []rdf.Term{document}
	data[TitlePredicate] = []rdf.Term{title}

	resource := NewResource("http://example.com/record1", data)

	assert.Equal(t, "Hello world!", resource.Titles()[0].String())
	assert.Equal(t, "http://purl.org/ontology/bibo/Document", resource.ResourceTypes()[0].String())
	assert.True(t, resource.IsPublication())

}