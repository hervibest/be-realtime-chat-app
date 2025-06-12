package repository

import (
	"net/http"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
)

type DB interface {
	DiscoverNodes() error
	InstrumentationEnabled() elastictransport.Instrumentation
	Metrics() (elastictransport.Metrics, error)
	Perform(req *http.Request) (*http.Response, error)
	doProductCheck(f func() error) error
}
