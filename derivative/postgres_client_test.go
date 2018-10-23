package derivative

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// MockRepository is a mocked object that implements the repository interface
type MockRepository struct {
	mock.Mock
}

func (f *MockRepository) SubjectsToResources(ids []string) ([]models.Resource, error) {
	args := f.Called(ids)
	return args.Get(0).([]models.Resource), args.Error(1)
}

func (f *MockRepository) AllResources(fun func([]models.Resource) error) error {
	return nil
}

func addPerson(client *PostgresClient, id string, subtype string, firstName string, lastName string) error {
	data := make(map[string]rdf.Term)
	// Id
	idIRI, _ := rdf.NewIRI(id)
	data["id"] = idIRI
	// Type
	resourceTypeIRI, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Person")
	data["type"] = resourceTypeIRI

	// First
	if firstName != "" {
		firstNameLiteral, _ := rdf.NewLiteral(firstName)
		data["firstname"] = firstNameLiteral
	}

	// Last
	if lastName != "" {
		lastNameLiteral, _ := rdf.NewLiteral(lastName)
		data["lastname"] = lastNameLiteral
	}

	// Subtype
	if subtype != "" {
		subtypeIRI, _ := rdf.NewIRI(subtype)
		data["subtype"] = subtypeIRI
	}

	resource := models.NewResource(data)

	return client.addPerson(resource.(*models.Person))

}

func TestPostgresAddPerson(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	repo := new(MockRepository)

	client := NewPostgresClient(conf, repo)
	client.RemoveAll()

	err := addPerson(client, "http://example.com/record1", "http://vivoweb.org/ontology/core#Student", "Barbara", "Liskov")
	assert.Nil(t, err)

	person, err := client.retrieveOnePerson("http://example.com/record1")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"name": "Barbara Liskov", "countries": null, "departments": [], "institutionalAffiliations": []}`, person)

}

func TestPostgresUpdatePerson(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	repo := new(MockRepository)

	client := NewPostgresClient(conf, repo)
	client.RemoveAll()

	err := addPerson(client, "http://example.com/record1", "http://vivoweb.org/ontology/core#Student", "Barbara", "Liskov")
	assert.Nil(t, err)

	err = addPerson(client, "http://example.com/record1", "http://vivoweb.org/ontology/core#Student", "B.", "Liskov")
	assert.Nil(t, err)

	person, err := client.retrieveOnePerson("http://example.com/record1")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"name": "B. Liskov", "countries": null, "departments": [], "institutionalAffiliations": []}`, person)

}

func TestPostgresAddOrganization(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	repo := new(MockRepository)

	client := NewPostgresClient(conf, repo)
	client.RemoveAll()

	data := make(map[string]rdf.Term)
	name, _ := rdf.NewLiteral("School of Engineering")
	school, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#School")
	id, _ := rdf.NewIRI("http://example.com/record1")
	resourceType, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")

	data["id"] = id
	data["type"] = resourceType
	data["name"] = name
	data["subtype"] = school

	resource := models.NewResource(data)

	err := client.addOrganization(resource.(*models.Organization))
	assert.Nil(t, err)

	org, err := client.retrieveOneOrganization("http://example.com/record1")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"name": "School of Engineering", "type": "http://vivoweb.org/ontology/core#School"}`, org)
}

func TestPostgresAddPublication(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	repo := new(MockRepository)

	client := NewPostgresClient(conf, repo)
	client.RemoveAll()

	// Add authors
	err := addPerson(client, "http://example.com/record1", "", "Barbara", "Liskov")
	assert.Nil(t, err)

	err = addPerson(client, "http://example.com/record2", "", "Barry", "Liskovich")
	assert.Nil(t, err)

	data := make(map[string]rdf.Term)

	titleLiteral, _ := rdf.NewLiteral("New developments in the management of narcolepsy")
	idIRI, _ := rdf.NewIRI("http://example.com/publication1")
	resourceType, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	createdLiteral, _ := rdf.NewLiteral("Right now")
	identifierLiteral, _ := rdf.NewLiteral("publication1")

	data["id"] = idIRI
	data["title"] = titleLiteral
	data["type"] = resourceType
	data["created"] = createdLiteral
	data["identifier"] = identifierLiteral

	resource := models.NewResource(data).(*models.Publication)

	// Add authors
	resource.Authors = append(resource.Authors, &models.Author{URI: "http://example.com/record1", Label: "Barbara Liskov"})
	resource.Authors = append(resource.Authors, &models.Author{URI: "http://example.com/record2", Label: "Barry Liskovich"})

	err = client.addPublication(resource)
	assert.Nil(t, err)

	pub, err := client.retrieveOnePublication("http://example.com/publication1")
	assert.Nil(t, err)
	assert.Equal(t, `{"title": "New developments in the management of narcolepsy"}`, pub)

	uris, err := client.retrievePeoplePublicationRelationship("http://example.com/publication1")
	assert.Nil(t, err)
	assert.Len(t, *uris, 2)

}
