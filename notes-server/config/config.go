package config

import "os"

type Config struct {
	Addr       string
	NotesPath  string
	StaticPath string
	SchemaPath string
}

func Load() *Config {
	// 优先使用环境变量，方便不同环境部署
	notesPath := os.Getenv("NOTES_PATH")
	if notesPath == "" {
		notesPath = "../changex-notes"
	}

	staticPath := os.Getenv("STATIC_PATH")
	if staticPath == "" {
		staticPath = "../notes-web/dist"
	}

	schemaPath := os.Getenv("SCHEMA_PATH")
	if schemaPath == "" {
		schemaPath = "../schema/note-schema.yaml"
	}

	return &Config{
		Addr:       ":8080",
		NotesPath:  notesPath,
		StaticPath: staticPath,
		SchemaPath: schemaPath,
	}
}
