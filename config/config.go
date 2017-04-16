package config

import (
	"log"
	"os"
)

type Config struct {
	PublicKey     string
	PrivateKey    string
	DatabaseURL   string
	URL           string
	Debug         bool
	MailgunAPIKey string
	MailgunDomain string
	MailgunURL    string
	MailerSubject string
	MailerFrom    string
}

func New() Config {
	return Config{
		PublicKey:   mustGetenv("MCI_PUBLIC_KEY"),
		PrivateKey:  mustGetenv("MCI_PRIVATE_KEY"),
		DatabaseURL: mustGetenv("DATABASE_URL"),
		URL:         defaultGetenv("MCI_URL", "https://gateway.marvel.com"),
		Debug:       boolGetenv("DEBUG"),

		MailgunAPIKey: mustGetenv("MAILGUN_API_KEY"),
		MailgunDomain: mustGetenv("MAILGUN_DOMAIN"),
		MailgunURL:    defaultGetenv("MAILGUN_URL", "https://api.mailgun.net"),

		MailerSubject: mustGetenv("MAILER_SUBJECT"),
		MailerFrom:    mustGetenv("MAILER_FROM"),
	}
}

func mustGetenv(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}

	log.Panicf("%q not found in environment", k)
	return ""
}

func boolGetenv(k string) bool {
	if os.Getenv(k) == "" {
		return false
	}

	return true
}

func defaultGetenv(k, dv string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}

	return dv
}
