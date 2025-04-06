package goelasticmigrator

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const (
	noAPI     string = ""
	builkAPI  string = "_bulk"
	deleteAPI string = "_delete_by_query"
)

// Error messages to indicate various failure scenarios.
var (
	ErrCantEstablishConnection = fmt.Errorf("error: connection to elasticsearch can't be established")
	ErrMappingsDontMatch       = fmt.Errorf("error: index exists, but mappings don't match")
	ErrIndexDoesntExists       = fmt.Errorf("error: index doesn't exists")
	ErrCantCreateIndex         = fmt.Errorf("error: can't create an index")
	ErrCantCreateDocuments     = fmt.Errorf("error: can't create documents")
	ErrCantDeleteDocuments     = fmt.Errorf("error: can't delete documents")
	ErrUnexpected              = fmt.Errorf("error: unexpected")
)

// Migrator struct encapsulates the migration logic.
type Migrator struct {
	MigratorConfig
}

// MigratorConfig holds configurations for the migrator.
type MigratorConfig struct {
	Client        http.Client
	Log           *slog.Logger
	ElasticSearch ElasticSearch
}

// ElasticSearch struct holds the address, index details, and credentials for Elasticsearch.
type ElasticSearch struct {
	Address     string
	Index       Index
	Credentials Credentials
}

// Index holds the index name and its definition (mappings).
type Index struct {
	Name       string
	Definition map[string]interface{}
}

// Credentials struct stores Elasticsearch credentials (username & password).
type Credentials struct {
	Username string
	Password string
}

// New initializes a new Migrator instance with the provided configuration.
func New(cfg MigratorConfig) *Migrator {
	cfg.Log.Info(
		"creating migrator",
		slog.Any("config", cfg),
	)

	return &Migrator{
		cfg,
	}
}

// MigrateUp() performs an "up" migration, either creating an index or adding documents.
func (m *Migrator) MigrateUp(path string) error {
	m.Log.Info(
		"migrating up",
	)

	m.Log.Info(
		"checking index",
	)

	err := checkIndex(m)
	if err != nil {
		if errors.Is(err, ErrIndexDoesntExists) {
			m.Log.Info(
				"creating index",
			)

			err := createIndex(m)
			if err != nil {
				m.Log.Error(
					"can't create index",
					slog.String("error", err.Error()),
				)

				return err
			}
		} else {
			m.Log.Error(
				"index check failed",
				slog.String("error", err.Error()),
			)

			return err
		}
	}

	if path == "" {
		m.Log.Info(
			"migration completed, documents ommited",
		)

		return nil
	} else {
		m.Log.Info(
			"creating documents",
		)

		err := createDocuments(m, path)
		if err != nil {
			m.Log.Error(
				"can't create documents",
				slog.String("error", err.Error()),
			)

			return err
		}
		m.Log.Info(
			"migration completed, documents created",
		)

		return nil
	}
}

// MigrateDown() performs a "down" migration, either deleting documents or the entire index.
func (m *Migrator) MigrateDown(documentsOnly bool) error {
	m.Log.Info(
		"migrating down",
	)

	m.Log.Info(
		"checking index",
	)

	err := checkIndex(m)
	if err != nil {
		m.Log.Error(
			"index check failed",
			slog.String("error", err.Error()),
		)

		return err
	}
	if documentsOnly {
		m.Log.Info(
			"deleting document only",
		)

		err := deleteDocuments(m)
		if err != nil {
			m.Log.Error(
				"can't delete documents",
				slog.String("error", err.Error()),
			)

			return err
		}
		m.Log.Info(
			"migration completed, documents deleted",
		)

		return nil
	} else {
		m.Log.Info(
			"deleting whole index",
		)

		err = deleteIndex(m)
		if err != nil {
			m.Log.Error(
				"can't delete index",
				slog.String("error", err.Error()),
			)

			return err
		}
		m.Log.Info(
			"migration completed, index deleted",
		)

		return nil
	}
}

// createDocuments() handles the creation of documents in the Elasticsearch index.
func createDocuments(m *Migrator, path string) error {
	documents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	req, err := buildRequest(
		http.MethodPost,
		m.ElasticSearch.Credentials.Username,
		m.ElasticSearch.Credentials.Password,
		buildUrl(m.ElasticSearch.Address, m.ElasticSearch.Index.Name, builkAPI),
		bytes.NewBuffer(documents),
	)
	if err != nil {
		return err
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return ErrCantEstablishConnection
	}
	defer resp.Body.Close()

	body, err := decode(resp.Body)
	if err != nil {
		return err
	}
	_, ok := body["error"].(map[string]interface{})
	if ok {
		return ErrCantCreateDocuments
	}

	return nil
}

// deleteDocuments() handles the deletion of documents in the Elasticsearch index.
func deleteDocuments(m *Migrator) error {
	req, err := buildRequest(
		http.MethodPost,
		m.ElasticSearch.Credentials.Username,
		m.ElasticSearch.Credentials.Password,
		buildUrl(m.ElasticSearch.Address, m.ElasticSearch.Index.Name, deleteAPI),
		bytes.NewBuffer(marshal(map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		})),
	)
	if err != nil {
		return err
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return ErrCantEstablishConnection
	}
	defer resp.Body.Close()

	body, err := decode(resp.Body)
	if err != nil {
		return err
	}
	_, ok := body["error"].(map[string]interface{})
	if ok {
		return ErrCantDeleteDocuments
	}
	return nil
}

// checkIndex() checks if the specified index exists in Elasticsearch and if the mappings match.
func checkIndex(m *Migrator) error {
	req, err := buildRequest(
		http.MethodGet,
		m.ElasticSearch.Credentials.Username,
		m.ElasticSearch.Credentials.Password,
		buildUrl(m.ElasticSearch.Address, m.ElasticSearch.Index.Name, noAPI),
		nil,
	)
	if err != nil {
		return err
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return ErrCantEstablishConnection
	}
	defer resp.Body.Close()

	body, err := decode(resp.Body)
	if err != nil {
		return err
	}
	_, ok := body["error"].(map[string]interface{})
	if ok {
		return ErrIndexDoesntExists
	}
	mappings := body[m.ElasticSearch.Index.Name].(map[string]interface{})["mappings"].(map[string]interface{})

	ok = bytes.Equal(marshal(mappings), marshal(m.ElasticSearch.Index.Definition["mappings"].(map[string]interface{})))
	if !ok {
		return ErrMappingsDontMatch
	}

	return nil
}

// createIndex() creates a new index in Elasticsearch.
func createIndex(m *Migrator) error {
	req, err := buildRequest(
		http.MethodPut,
		m.ElasticSearch.Credentials.Username,
		m.ElasticSearch.Credentials.Password,
		buildUrl(m.ElasticSearch.Address, m.ElasticSearch.Index.Name, noAPI),
		bytes.NewBuffer(marshal(m.ElasticSearch.Index.Definition)),
	)
	if err != nil {
		return err
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return ErrCantEstablishConnection
	}
	defer resp.Body.Close()

	body, err := decode(resp.Body)
	if err != nil {
		return err
	}
	_, ok := body["error"].(map[string]interface{})
	if ok {
		return ErrCantCreateIndex
	}

	return nil
}

// deleteIndex() deletes the specified index from Elasticsearch.
func deleteIndex(m *Migrator) error {
	req, err := buildRequest(
		http.MethodDelete,
		m.ElasticSearch.Credentials.Username,
		m.ElasticSearch.Credentials.Password,
		buildUrl(m.ElasticSearch.Address, m.ElasticSearch.Index.Name, noAPI),
		nil,
	)
	if err != nil {
		return err
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return ErrCantEstablishConnection
	}
	defer resp.Body.Close()

	return nil
}
