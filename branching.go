package branching

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	bexpr "github.com/hashicorp/go-bexpr"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/config/runtime"
	"github.com/traefik/traefik/v2/pkg/server/middleware"
)

// Logger Main logger
var (
	Logger = log.New(os.Stdout, "Branching: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// Config the plugin configuration.
type Config struct {
	Condition string
	Chain     map[string]*dynamic.Middleware
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// Branching plugin
type Branching struct {
	name    string
	next    http.Handler
	branch  http.Handler
	matcher *bexpr.Evaluator
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	eval, err := bexpr.CreateEvaluator(config.Condition)
	if err != nil {
		return nil, fmt.Errorf("failed to create evaluator for expression %q: %w", config.Condition, err)
	}

	var midChain []string
	for name := range config.Chain {
		midChain = append(midChain, name)
	}

	rtConf := runtime.NewConfig(dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Middlewares: config.Chain,
		},
	})
	builder := middleware.NewBuilder(rtConf.Middlewares, nil, nil)

	chain := builder.BuildChain(ctx, midChain)
	branchHandler, err := chain.Then(next)
	if err != nil {
		return nil, fmt.Errorf("failed to create middleware chain %w", err)
	}

	Logger.Printf("%s created, matching %q", name, config.Condition)
	return &Branching{
		name:    name,
		next:    next,
		matcher: eval,
		branch:  branchHandler,
	}, nil
}

func (e *Branching) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	match, err := e.matcher.Evaluate(req)
	if err != nil {
		Logger.Printf("ignoring branch, unable to match request: %v", err)
	}

	if match {
		e.branch.ServeHTTP(rw, req)
		return
	}

	e.next.ServeHTTP(rw, req)
}
