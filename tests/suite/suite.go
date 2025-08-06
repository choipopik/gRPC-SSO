package suite

import (
	"testing"

	"github.com/choipopik/gRPC-SSO/internal/config"
)

type Suite struct {
	*testing.T
	Cfg *config.Config
}
