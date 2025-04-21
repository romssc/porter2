package suite

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/xoticdsign/porter"
)

type suite struct {
	*testing.T

	Elasticsearch *es
	Migrator      migrator
}

type es struct {
	client

	Container testcontainers.Container
}

type client struct {
	*elasticsearch.Client
}

func (c client) IsIndexExist(ctx context.Context, name string) (*esapi.Response, error) {
	return nil, nil
}

func (c client) CreateIndex(ctx context.Context, name string, body []byte) (*esapi.Response, error) {
	return nil, nil
}

func (c client) CreateDocuments(ctx context.Context, name string, documents []byte) (*esapi.Response, error) {
	return nil, nil
}

func (c client) DeleteIndex(ctx context.Context, name string) (*esapi.Response, error) {
	return nil, nil
}

func (c client) DeleteDocuments(ctx context.Context, name string, query string) (*esapi.Response, error) {
	return nil, nil
}

type migrator struct {
	Components porter.C
}

func New(t *testing.T, d bool) *suite {
	t.Helper()
	t.Parallel()

	components := porter.GetComponents()

	if !d {
		return &suite{
			T: t,

			Elasticsearch: &es{
				client: client{},
			},
			Migrator: migrator{
				Components: components,
			},
		}
	} else {
		container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "docker.io/elasticsearch:8.16.6",
				ExposedPorts: []string{"9200/tcp"},
				Env: map[string]string{
					"xpack.security.enabled": "false",
					"discovery.type":         "single-node",
					"ES_JAVA_OPTS":           "-Xms512m -Xmx512m",
					"logger.level":           "WARN",
					"network.host":           "0.0.0.0",
					"cluster.name":           "test-cluster",
					"bootstrap.memory_lock":  "true",
					"path.repo":              "/tmp",
					"node.name":              "test-node",
				},
				WaitingFor: wait.ForHTTP("/").WithPort("9200/tcp").WithStartupTimeout(60 * time.Second),
			},
			Started: true,
		})
		if err != nil {
			panic(err)
		}

		host, err := container.Host(context.Background())
		if err != nil {
			panic(err)
		}

		port, err := container.MappedPort(context.Background(), "9200")
		if err != nil {
			panic(err)
		}

		cc, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{fmt.Sprintf("http://%s:%s", host, port.Port())},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		})
		if err != nil {
			panic(err)
		}

		return &suite{
			T: t,

			Elasticsearch: &es{
				client:    client{cc},
				Container: container,
			},
			Migrator: migrator{
				Components: components,
			},
		}
	}
}
