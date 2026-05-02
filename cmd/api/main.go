package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/auth"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/db"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/env"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/mailer"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store/cache"

	_ "github.com/lib/pq"
)

// @title Social Media API
// @version 1.0
// @description Social Media Backend API
// @host localhost:8080
// @BasePath /v1
func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		auth: authConfig{
			token: tokenConfig{
				exp:       time.Hour * 24 * 3,
				secretKey: env.GetString("JWT_SECRET", ""),
			},
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{

			exp: time.Hour * 24 * 3,
			mailTrap: mailTrapConfig{
				host:      env.GetString("MAILTRAP_HOST", "live.smtp.mailtrap.io"),
				port:      env.GetInt("MAILTRAP_PORT", 587),
				username:  env.GetString("MAILTRAP_USERNAME", ""),
				password:  env.GetString("MAILTRAP_PASSWORD", ""),
				fromEmail: env.GetString("FROM_EMAIL_MT", "kumarsahiljee19@gmail.com"),
			},
			sendGrid: sendGridConfig{
				apikey:    env.GetString("SENDGRID_API_KEY", ""),
				fromEmail: env.GetString("FROM_EMAIL_SG", "sk2000jee@gmail.com"),
			},
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {

		log.Fatal(err)
	}
	defer db.Close()

	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		fmt.Println("redis cache connection established")

		defer rdb.Close()
	}
	var mailClient mailer.Client
	switch {
	case cfg.mail.mailTrap.username != "" && cfg.mail.mailTrap.password != "":
		mailClient, err = mailer.NewMailTrapMailer(
			cfg.mail.mailTrap.fromEmail,
			cfg.mail.mailTrap.host,
			cfg.mail.mailTrap.port,
			cfg.mail.mailTrap.username,
			cfg.mail.mailTrap.password,
		)
	case cfg.mail.sendGrid.apikey != "":
		mailClient, err = mailer.NewSendgrid(cfg.mail.sendGrid.fromEmail, cfg.mail.sendGrid.apikey)
	default:
		log.Fatal("no mail provider configured")
	}
	if err != nil {
		log.Fatal(err)
	}
	jwtClient := auth.NewJWT(cfg.auth.token.secretKey, cfg.auth.token.exp)

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)
	app := &application{
		config:       cfg,
		store:        store,
		mailer:       mailClient,
		auth:         jwtClient,
		cacheStorage: cacheStorage,
	}

	mux := app.mount()
	srv := app.run(mux)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- srv.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	case <-ctx.Done():
		log.Println("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}

		log.Println("server shut down gracefully")
	}

}
