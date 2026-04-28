package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/auth"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/db"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/env"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/mailer"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"

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
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {

		log.Fatal(err)
	}
	defer db.Close()

	store := store.NewStorage(db)
	fmt.Print(store)

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
	jwtClient := auth.NewJWT(cfg.auth.token.secretKey)
	app := &application{
		config: cfg,
		store:  store,
		mailer: mailClient,
		auth:   jwtClient,
	}

	mux := app.mount()
	err = app.run(mux)
	fmt.Print("bye")

	if err != nil {
		log.Fatal(err)
	}

}
