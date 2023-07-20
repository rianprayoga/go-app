module example.com/chi-restful

go 1.20

require github.com/go-chi/chi/v5 v5.0.10

require github.com/google/uuid v1.3.0

require github.com/lib/pq v1.10.9

replace example.com/chi-restful/pkg/handlers => ../pkg/handlers
