package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

func V5_3_0(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	// Insert button appearance settings.
	_, err := db.Exec(`
		INSERT INTO settings (key, value) VALUES
			('appearance.button.bg_color', '"#0055d4"'),
			('appearance.button.text_color', '"#ffffff"'),
			('appearance.button.hover_bg_color', '"#222222"'),
			('appearance.button.hover_text_color', '"#ffffff"'),
			('appearance.button.border_radius', '"3px"')
		ON CONFLICT (key) DO NOTHING
	`)
	if err != nil {
		return err
	}

	return nil
}

