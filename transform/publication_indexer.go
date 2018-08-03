package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// PublicationIndexer transforms publication resource types into solr Documents
type PublicationIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *PublicationIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	mapping := map[string]string{
		"title":            "title_tesi",
		"abstract":         "abstract_tesim",
		"cites":            "cites_ssim",
		"doi":              "doi_ssim",
		"identifier":       "identifier_ssim",
		"link":             "link_ssim",
		"description":      "description_tesim",
		"fundedBy":         "funded_by_ssim",
		"publisher":        "publisher_label_tsim", // TODO: Needs URI lookup
		"sponsor":          "sponsor_label_tsim",   // TODO: Needs URI lookup
		"hasInstrument":    "has_instrument_ssim",
		"sameAs":           "same_as_ssim",
		"journalIssue":     "journal_issue_ssim",
		"subject":          "subject_label_ssim", // TODO: Needs URI
		"alternativeTitle": "alternative_title_tesim",
	}

	doc = indexMapping(resource, doc, mapping)

	// TODO: complex lookups
	// author 	vivo:relatedBy vivo:Authorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Author of the publication.
	// Profiles confirmed 	vivo:relatedBy vivo:Authorship dcterms:source 	"Profiles" string-literal 	[0,1] 	If the authorship relationship has been confirmed by the Author in Profiles. Can be reused for any relationship needed (i.e. Editorship, Advising Relationship, etc.)
	// editor 	vivo:relatedBy vivo:Editorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Editor of the publication.

	date := resource.ValueOf("dateOfCreation")
	if date != nil {
		// TODO: This may be a ETDF in the resource, but we need it to be a Solr date
		doc.Set("date_created_dtsi", resource.ValueOf("dateOfCreation")[0].String())
	}

	return doc
}