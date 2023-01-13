package monitoring

import (
	"testing"
	"time"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"
	"github.com/stretchr/testify/assert"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory"
)

func TestNewMonitor(t *testing.T) {
	memoryStorage := (&factory.StorageFactory{}).CreateStorage()
	clientToServer := client.NewClient(server.HTTPType, "https://test.com", 60*time.Second, "secret", "", "")

	monitor := NewMonitor(memoryStorage, clientToServer)

	assert.Implements(t, (*IMonitor)(nil), monitor)
}
