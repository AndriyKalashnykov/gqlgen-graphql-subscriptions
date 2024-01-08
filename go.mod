module github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions

go 1.21

require (
	github.com/99designs/gqlgen v0.17.42
	github.com/gorilla/websocket v1.5.1
	github.com/labstack/echo/v4 v4.11.4
	github.com/redis/go-redis/v9 v9.4.0
	github.com/thanhpk/randstr v1.0.6
	github.com/vektah/gqlparser/v2 v2.5.10
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/sosodev/duration v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)

replace github.com/redis/go-redis/v9 => github.com/go-redis/redis/v9 v9.4.0

replace github.com/go-redis/redis/v9 => github.com/redis/go-redis/v9 v9.4.0
