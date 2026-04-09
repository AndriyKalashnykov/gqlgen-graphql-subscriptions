package graphql

import (
	"net/http"
	"os"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph/generated"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/constants"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"

	"github.com/99designs/gqlgen/graphql/handler"
)

func NewGraphQLServer(resolver *graph.Resolver) *handler.Server {
	allowedOrigins := getAllowedOrigins()

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				if len(allowedOrigins) == 0 {
					return true
				}
				origin := r.Header.Get("Origin")
				for _, allowed := range allowedOrigins {
					if allowed == "*" || allowed == origin {
						return true
					}
				}
				return false
			},
			ReadBufferSize:  constants.WebSocketReadBufferSize,
			WriteBufferSize: constants.WebSocketWriteBufferSize,
		},
		KeepAlivePingInterval: constants.WebSocketKeepAlivePing,
	})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](constants.QueryCacheSize))

	if os.Getenv("DISABLE_INTROSPECTION") != "true" {
		srv.Use(extension.Introspection{})
	}
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](constants.APQCacheSize),
	})

	return srv
}

// getAllowedOrigins reads CORS origins from ALLOWED_ORIGINS env var.
// Returns empty slice if unset (allows all origins for local development).
func getAllowedOrigins() []string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		return nil
	}
	return strings.Split(origins, ",")
}
