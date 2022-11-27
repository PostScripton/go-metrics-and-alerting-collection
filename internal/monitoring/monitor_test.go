package monitoring

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {
	memoryStorage := (&factory.StorageFactory{}).CreateStorage()
	clientToServer := client.NewClient("https://test.com", 60*time.Second, "secret")

	monitor := NewMonitor(memoryStorage, clientToServer)

	assert.Implements(t, (*Monitorer)(nil), monitor)
}
