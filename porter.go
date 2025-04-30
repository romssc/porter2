package porter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/xoticdsign/porter/internal/utils"
)

var (
	ErrClientBadConnection     = fmt.Errorf("elasticsearch client: failed to establish a connection")
	ErrClientCreatingIndex     = fmt.Errorf("elasticsearch client: failed to create index")
	ErrClientCreatingDocuments = fmt.Errorf("elasticsearch client: bulk insert operation failed")
	ErrClientDeletingIndex     = fmt.Errorf("elasticsearch client: failed to delete index")
	ErrClientDeletingDocuments = fmt.Errorf("elasticsearch client: failed to delete documents by query")

	ErrMigratorMigratingIndex = fmt.Errorf("migrator: index operation failed during migration process")
	ErrMigratorDocuments      = fmt.Errorf("migrator: document operation failed during migration process")

	ErrPorterMigratingUp   = fmt.Errorf("porter: failed to perform 'up' migration")
	ErrPorterMigratingDown = fmt.Errorf("porter: failed to perform 'down' migration")
)

// Constants for migration direction
const (
	directionUp   = 1
	directionDown = 2
)

// M{} represents the migration object that holds information about index and document migration.
type M struct {
	Index     index
	Documents documents

	Client searcher
}

// index{} represents the settings and mappings of the index
type index struct {
	Settings settings
	Mappings mappings
}

// settings{} represents the analysis settings of the index, including analyzers and normalizers
type settings struct {
	Analysis analysis
}

// analysis{} defines the analyzers and normalizers for the index.
type analysis struct {
	Analyzer   analyzer
	Normalizer normalizer
}

// mappings{} defines the properties of the index, i.e., the fields in the documents.
type mappings struct {
	Properties properties
}

// properties{} defines various field types in the index.
type properties struct {
	fields
}

// documents{} represents the documents that need to be migrated.
type documents struct {
	Origin origin
}

// origin{} represents the location or source of the documents.
type origin struct {
	location
}

// searcher{} defines methods for interacting with Elasticsearch for index and document operations.
type searcher interface {
	CreateIndex(ctx context.Context, name string, body []byte) error
	CreateDocuments(ctx context.Context, name string, documents []byte) error
	DeleteIndex(ctx context.Context, name string) error
	DeleteDocuments(ctx context.Context, name string, query string) error
}

// client{} wraps the Elasticsearch client and provides convenience methods for interacting with Elasticsearch.
type client struct {
	*elasticsearch.Client
}

func (c client) CreateIndex(ctx context.Context, name string, body []byte) error {
	resp, err := c.Indices.Create(
		name,
		c.Indices.Create.WithContext(ctx),
		c.Indices.Create.WithBody(bytes.NewBuffer(body)),
		c.Indices.Create.WithPretty(),
	)
	if err != nil {
		return fmt.Errorf("%w [%s]", ErrClientBadConnection, err)
	}
	defer resp.Body.Close()

	r, ok := utils.ExtractError(resp.Body)
	if ok {
		return fmt.Errorf("%w [%s]", ErrClientCreatingIndex, r)
	}

	return nil
}

func (c client) CreateDocuments(ctx context.Context, name string, documents []byte) error {
	resp, err := c.Bulk(
		bytes.NewBuffer(documents),
		c.Bulk.WithContext(context.Background()),
		c.Bulk.WithIndex(name),
		c.Bulk.WithPretty(),
	)
	if err != nil {
		return fmt.Errorf("%w [%s]", ErrClientBadConnection, err)
	}
	defer resp.Body.Close()

	var r map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&r)

	ok, _ := r["errors"].(bool)
	if !ok {
		return nil
	}

	items, ok := r["items"].([]interface{})
	if !ok {
		return nil
	}

	var errors []string

	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		for _, v := range m {
			doc, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			err, ok := doc["error"].(map[string]interface{})
			if ok {
				errors = append(errors, err["reason"].(string))
			}
		}
	}

	return fmt.Errorf("%w [%s]", ErrClientCreatingDocuments, strings.Join(errors, " + "))
}

func (c client) DeleteIndex(ctx context.Context, name string) error {
	resp, err := c.Indices.Delete(
		[]string{name},
		c.Indices.Delete.WithContext(context.Background()),
		c.Indices.Delete.WithPretty(),
	)
	if err != nil {
		return fmt.Errorf("%w [%s]", ErrClientBadConnection, err)
	}
	defer resp.Body.Close()

	r, ok := utils.ExtractError(resp.Body)
	if ok {
		return fmt.Errorf("%w [%s]", ErrClientDeletingIndex, r)
	}
	return nil
}

func (c client) DeleteDocuments(ctx context.Context, name string, query string) error {
	resp, err := c.DeleteByQuery(
		[]string{name},
		strings.NewReader(query),
		c.DeleteByQuery.WithContext(context.Background()),
		c.DeleteByQuery.WithPretty(),
	)
	if err != nil {
		return fmt.Errorf("%w [%s]", ErrClientBadConnection, err)
	}
	defer resp.Body.Close()

	r, ok := utils.ExtractError(resp.Body)
	if ok {
		return fmt.Errorf("%w [%s]", ErrClientDeletingDocuments, r)
	}
	return nil
}

// New() initializes and returns a new migration object.
func New(cc *elasticsearch.Client) M {
	return M{
		Index: index{
			Settings: settings{
				Analysis: analysis{
					Analyzer: analyzer{
						Standart:    newAnalyzerStandard(),
						Simple:      newAnalyzerSimple(),
						Whitespace:  newAnalyzerWhitespace(),
						Stop:        newAnalyzerStop(),
						Keyword:     newAnalyzerKeyword(),
						Pattern:     newAnalyzerPattern(),
						Language:    newAnalyzerLanguage(),
						Fingerprint: newAnalyzerFingerprint(),
						Custom:      newAnalyzerCustom(),
					},
					Normalizer: normalizer{
						Custom: newNormalizerCustom(),
					},
				},
			},
			Mappings: mappings{
				Properties: properties{
					fields: fields{
						Keyword:     newFieldKeyword(),
						Text:        newFieldText(),
						Integer:     newFieldInteger(),
						Long:        newFieldLong(),
						Float:       newFieldFloat(),
						Double:      newFieldDouble(),
						Short:       newFieldShort(),
						Byte:        newFieldByte(),
						HalfFloat:   newFieldHalfFloat(),
						ScaledFloat: newFieldScaledFloat(),
						Date:        newFieldDate(),
						Boolean:     newFieldBoolean(),
						IP:          newFieldIP(),
					},
				},
			},
		},
		Documents: documents{
			Origin: origin{
				location: location{
					FromFile: newLocationFromFile(),
					Generate: newLocationGenerate(),
				},
			},
		},

		Client: client{cc},
	}
}

// temp{} represents temporary data during the migration process, including the direction (up or down).
type temp struct {
	direction int

	config Config
	client searcher
}

// MigrateUp() performs the "up" migration, which includes creating/updating the index and migrating documents.
func (m M) MigrateUp(config Config, index indexFunc, documents documentsFunc) error {
	t := temp{
		direction: directionUp,

		config: config,
		client: m.Client,
	}

	err := index(t)
	if err != nil {
		return fmt.Errorf("%w\n%v", ErrPorterMigratingUp, err)
	}

	err = documents(t)
	if err != nil {
		return fmt.Errorf("%w\n%v", ErrPorterMigratingUp, err)
	}

	return nil
}

// MigrateDown() performs the "down" migration, which includes deleting documents and the index.
func (m M) MigrateDown(config Config, documents documentsFunc, index indexFunc) error {
	t := temp{
		direction: directionDown,

		config: config,
		client: m.Client,
	}

	err := documents(t)
	if err != nil {
		return fmt.Errorf("%w\n%v", ErrPorterMigratingDown, err)
	}

	err = index(t)
	if err != nil {
		return fmt.Errorf("%w\n%v", ErrPorterMigratingDown, err)
	}

	return nil
}

type indexFunc func(t temp) error

// NoIndex() represents no operation for the index during migration (used for down migrations).
func (i index) NoIndex() indexFunc {
	return func(t temp) error {
		return nil
	}
}

// MigrateIndex() migrates the index up or down, depending on the migration direction.
func (i index) MigrateIndex() indexFunc {
	return func(t temp) error {
		if t.direction == directionUp {
			err := t.client.CreateIndex(context.Background(), t.config.Name, utils.MarshalJSON(t.config.Definition))
			if err != nil {
				return fmt.Errorf("%w\n%s", ErrMigratorMigratingIndex, err)
			}
			return nil
		} else {
			err := t.client.DeleteIndex(context.Background(), t.config.Name)
			if err != nil {
				return fmt.Errorf("%w\n%s", ErrMigratorMigratingIndex, err)
			}
			return nil
		}
	}
}

type documentsFunc func(t temp) error

// NoDocuments() represents no operation for documents during migration (used for down migrations).
func (d documents) NoDocuments() documentsFunc {
	return func(t temp) error {
		return nil
	}
}

// MigrateDocuments migrates the documents up or down, depending on the migration direction.
func (d documents) MigrateDocuments(origin OriginFunc) documentsFunc {
	return func(t temp) error {
		if t.direction == directionUp {
			docs, err := origin(t)
			if err != nil {
				return fmt.Errorf("%w\n%v", ErrMigratorDocuments, err)
			}

			err = t.client.CreateDocuments(context.Background(), t.config.Name, docs)
			if err != nil {
				return fmt.Errorf("%w\n%v", ErrMigratorDocuments, err)
			}
			return nil
		} else {
			err := t.client.DeleteDocuments(context.Background(), t.config.Name, `{"query": {"match_all": {}}}`)
			if err != nil {
				return fmt.Errorf("%w\n%v", ErrMigratorDocuments, err)
			}
			return nil
		}
	}
}
