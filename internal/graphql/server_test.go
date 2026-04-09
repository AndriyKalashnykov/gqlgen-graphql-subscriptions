package graphql

import (
	"testing"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/testutil"
)

func TestNewGraphQLServer(t *testing.T) {
	mock := &testutil.MockRedisClient{}
	resolver := graph.NewResolver(mock)
	srv := NewGraphQLServer(resolver)

	if srv == nil {
		t.Fatal("expected server to be created, got nil")
	}
}

func TestNewGraphQLServer_IntrospectionDisabled(t *testing.T) {
	t.Setenv("DISABLE_INTROSPECTION", "true")

	mock := &testutil.MockRedisClient{}
	resolver := graph.NewResolver(mock)
	srv := NewGraphQLServer(resolver)

	if srv == nil {
		t.Fatal("expected server to be created, got nil")
	}
}

func TestGetAllowedOrigins_Default(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "")

	origins := getAllowedOrigins()
	if origins != nil {
		t.Errorf("expected nil for unset env, got %v", origins)
	}
}

func TestGetAllowedOrigins_Single(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")

	origins := getAllowedOrigins()
	if len(origins) != 1 {
		t.Fatalf("expected 1 origin, got %d", len(origins))
	}
	if origins[0] != "http://localhost:3000" {
		t.Errorf("expected 'http://localhost:3000', got %s", origins[0])
	}
}

func TestGetAllowedOrigins_Multiple(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,https://example.com")

	origins := getAllowedOrigins()
	if len(origins) != 2 {
		t.Fatalf("expected 2 origins, got %d", len(origins))
	}
}
