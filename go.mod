module github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions

go 1.25.1

require (
	github.com/99designs/gqlgen v0.17.81
	github.com/gorilla/websocket v1.5.3
	// ... existing code ...
	github.com/vektah/gqlparser/v2 v2.5.30
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	// ... existing code ...
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	golang.org/x/time v0.13.0 // indirect
)

require (
	github.com/labstack/echo/v4 v4.13.4
	github.com/redis/go-redis/v9 v9.14.1
	github.com/thanhpk/randstr v1.0.6
)

require (
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
)

// Pin golang.org/x/tools to a version compatible with Go 1.23.
replace golang.org/x/tools => golang.org/x/tools v0.38.0
