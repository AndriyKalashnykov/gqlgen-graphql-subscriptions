package graphql

import (
	"net/http"
	"time"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"

	"github.com/99designs/gqlgen/graphql/handler"
)

func NewGraphQLServer(resolver *graph.Resolver) *handler.Server {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}
