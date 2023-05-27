module pulse-service

go 1.16

require (
	github.com/gin-gonic/gin v1.9.0
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-redis/redis/v8 v8.11.5
	github.com/hashicorp/go-retryablehttp v0.7.2
	github.com/rs/zerolog v1.29.1
	go.uber.org/zap v1.24.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.1
)
