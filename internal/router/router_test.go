package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph"
	graphqlserver "github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/graphql"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/testutil"
)

func TestNewRouter(t *testing.T) {
	mock := &testutil.MockRedisClient{}
	resolver := graph.NewResolver(mock)
	srv := graphqlserver.NewGraphQLServer(resolver)

	e := NewRouter(echo.New(), srv)

	if e == nil {
		t.Fatal("expected echo instance, got nil")
	}
}

func TestNewRouter_WelcomeRoute(t *testing.T) {
	mock := &testutil.MockRedisClient{}
	resolver := graph.NewResolver(mock)
	srv := graphqlserver.NewGraphQLServer(resolver)

	e := NewRouter(echo.New(), srv)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "Welcome!" {
		t.Errorf("expected 'Welcome!', got %s", rec.Body.String())
	}
}

func TestNewRouter_PlaygroundRoute(t *testing.T) {
	mock := &testutil.MockRedisClient{}
	resolver := graph.NewResolver(mock)
	srv := graphqlserver.NewGraphQLServer(resolver)

	e := NewRouter(echo.New(), srv)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/playground", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestGetAllowedOrigins_Default(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "")

	origins := getAllowedOrigins()
	if len(origins) != 1 || origins[0] != "*" {
		t.Errorf("expected ['*'], got %v", origins)
	}
}

func TestGetAllowedOrigins_Custom(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,https://example.com")

	origins := getAllowedOrigins()
	if len(origins) != 2 {
		t.Fatalf("expected 2 origins, got %d", len(origins))
	}
	if origins[0] != "http://localhost:3000" {
		t.Errorf("expected first origin 'http://localhost:3000', got %s", origins[0])
	}
}
