package suite

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	porter "github.com/xoticdsign/porter2"
)

type suite struct {
	T *testing.T

	Porter porter.M
	Temp   porter.Temp
}

type mockClient struct{}

func (m mockClient) CreateIndex(ctx context.Context, name string, body []byte) error {
	return nil
}

func (m mockClient) CreateDocuments(ctx context.Context, name string, documents []byte) error {
	return nil
}

func (m mockClient) DeleteIndex(ctx context.Context, name string) error {
	return nil
}

func (m mockClient) DeleteDocuments(ctx context.Context, name string, query string) error {
	return nil
}

func New(t *testing.T, offline bool) (*suite, error) {
	t.Helper()
	t.Parallel()

	if !offline {
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
			return nil, err
		}

		host, err := container.Host(context.Background())
		if err != nil {
			return nil, err
		}

		port, err := container.MappedPort(context.Background(), "9200")
		if err != nil {
			return nil, err
		}

		cc, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{fmt.Sprintf("http://%s", net.JoinHostPort(host, port.Port()))},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		})
		if err != nil {
			return nil, err
		}

		p := porter.New(cc)

		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			err := container.Terminate(ctx)
			if err != nil {
				panic(err)
			}
		})

		return &suite{
			T: t,

			Porter: p,
			Temp: porter.Temp{
				Config: porter.Config{Name: "porter"},
			},
		}, nil
	}
	p := porter.New(nil)
	p.Client = mockClient{}

	return &suite{
		T: t,

		Porter: p,
		Temp: porter.Temp{
			Config: porter.Config{Name: "porter"},
		},
	}, nil
}
