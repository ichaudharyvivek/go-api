# DDD Folder structure:
> This is generally used if you have a large-scale apps and you need modularity
```bash
myapp/
├── cmd/
│   ├── api/
│   │   └── main.go               # Entry point for HTTP server
│   └── migrate/
│       └── main.go               # DB migrations
│
├── internal/
│   ├── config/
│   │   ├── config.go             # Main config struct
│   │   ├── db.go                 # DB-specific config
│   │   └── server.go             # Server-specific config
│
│   ├── domain/
│   │   └── post/
│   │       ├── model.go          # Post entity: ID, Title, Content, etc.
│   │       ├── form.go           # DTO: request/response formatting + validation
│   │       └── service.go        # Interface + implementation of business logic
│
│   ├── app/
│   │   └── post/
│   │       └── service.go        # Optional separation if you want domain = pure interfaces only
│
│   ├── repository/
│   │   └── post/
│   │       └── repository.go     # GORM/Postgres implementation of post.Repository
│
│   ├── interfaces/
│   │   └── api/
│   │       └── v1/
│   │           └── post/
│   │               └── handler.go  # HTTP handler for Post (uses post.Service)
│
│   └── router/
│       └── router.go             # Sets up global router and plugs v1/post handler
│
├── go.mod
└── go.sum

```  
  
# Flattened Folder structure
```bash
myapp/
├── cmd/
│   ├── api/
│   │   └── main.go               # HTTP server entrypoint
│   └── migrate/
│       └── main.go               # DB migrations
│
├── internal/
│   ├── config/
│   │   ├── config.go             # Root config struct
│   │   ├── db.go                 # DB config logic
│   │   └── server.go             # Server config
│
│   ├── domain/
│   │   └── post/
│   │       ├── model.go          # Post entity (pure domain struct) + DTO (input/output + mapping)
│   │       ├── service.go        # Service interface + implementation
│   │       └── repository.go     # Repository interface
│
│   ├── repository/
│   │   └── post_repository.go    # GORM/Postgres implementation of domain.post.Repository
│
│   ├── handler/
│   ├── v1/
│   │   ├── post.go               # HTTP handler, routes for /api/v1/posts
│   │   ├── user.go
│   │   ├── comment.go
│   │   └── ... etc.
│   └── v2/
│       ├── post.go
│       └── ... (when needed)
│
│   └── router/
│   │   └── router.go             # chi mux setup and handler registration
│
├── go.mod
└── go.sum
```  


