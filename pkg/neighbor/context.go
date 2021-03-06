package neighbor

import (
	// stdlib
	"context"
	"fmt"
	"os"
	"strings"

	// external
	log "github.com/sirupsen/logrus"

	// internal
	"github.com/mccurdyc/neighbor/pkg/config"
)

// Ctx is an object that contains information that will be used throughout the
// life of the neighor command. The idea was taken from the dep tool (https://github.com/golang/dep/blob/master/context.go#L23).
// This does NOT satisfy the context.Context interface (https://golang.org/pkg/context/#Context),
// therefore, it cannot be used as a context for methods or functions requiring a context.Context.
type Ctx struct {
	Config       *config.Config  // the query config created by the user
	Context      context.Context // a context object required by the GitHub SDK
	GitHub       GitHubDetails
	Logger       *log.Logger // the logger to be used throughout the project
	NeighborDir  string
	ExtResultDir string   // where the external projects and test results will be stored
	ExternalCmd  []string // external project command and args
}

// GitHubDetails are GitHub-specifc details necessary throughout the project
type GitHubDetails struct {
	AccessToken string
	// SearchType is the GitHub search type https://developer.github.com/v3/search/#search
	SearchType string
	// Query is the GitHub search query to execute
	Query string
}

// NewCtx creates a pointer to a new neighbor context.
func NewCtx() *Ctx {
	return &Ctx{}
}

// CreateExternalResultDir creates the external projects and results directory if
// it doesn't exist.
func (ctx *Ctx) CreateExternalResultDir() error {
	_, err := os.Stat(ctx.ExtResultDir)
	if os.IsNotExist(err) {
		return os.Mkdir(ctx.ExtResultDir, os.ModePerm)
	}
	return nil
}

// Validate ensures that all of the required configuration attributes are set
// and valid.
func (ctx *Ctx) Validate() error {
	// @TODO: for now, we are just making sure that they are not empty, but in the future
	// we could actually perform some validation (e.g., SearchType is in the list of valid
	// search types, AccessToken fits the format of a token, etc.)
	var failed []string

	if len(ctx.GitHub.AccessToken) == 0 {
		failed = append(failed, "access token")
	}

	if len(ctx.GitHub.SearchType) == 0 {
		failed = append(failed, "search type")
	}

	if len(ctx.GitHub.Query) == 0 {
		failed = append(failed, "query")
	}

	if len(ctx.ExternalCmd) == 0 {
		failed = append(failed, "external command")
	}

	if len(failed) != 0 {
		fstr := strings.Join(failed, ", ")
		// need to append
		return fmt.Errorf("%s cannot be empty", fstr)
	}

	return nil
}
