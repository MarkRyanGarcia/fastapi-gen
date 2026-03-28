package generator

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed all:templates
var templateFS embed.FS

type ProjectConfig struct {
	ProjectName       string
	Database          string
	IncludeSQLAlchemy bool
	IncludeMongoDB    bool
}

// fileMap maps destination path -> template path (executed as Go templates)
func fileMap() map[string]string {
	return map[string]string{
		".gitignore":                        "templates/.gitignore.tmpl",
		"conftest.py":                       "templates/conftest.py.tmpl",
		"alembic.ini":                       "templates/alembic.ini.tmpl",
		"migrations/env.py":                 "templates/migrations/env.py.tmpl",
		"migrations/versions/.gitkeep":      "templates/migrations/versions/.gitkeep.tmpl",
		"app/main.py":                       "templates/app/main.py.tmpl",
		"app/__init__.py":                   "templates/app/__init__.py.tmpl",
		"app/api/v1/__init__.py":            "templates/app/api/v1/__init__.py.tmpl",
		"app/api/v1/routers/users.py":       "templates/app/api/v1/routers/users.py.tmpl",
		"app/api/v1/routers/items.py":       "templates/app/api/v1/routers/items.py.tmpl",
		"app/core/config.py":                "templates/app/core/config.py.tmpl",
		"app/core/security.py":              "templates/app/core/security.py.tmpl",
		"app/db/session.py":                 "templates/app/db/session.py.tmpl",
		"app/db/base.py":                    "templates/app/db/base.py.tmpl",
		"app/models/user.py":                "templates/app/models/user.py.tmpl",
		"app/models/item.py":                "templates/app/models/item.py.tmpl",
		"app/schemas/user.py":               "templates/app/schemas/user.py.tmpl",
		"app/schemas/item.py":               "templates/app/schemas/item.py.tmpl",
		"app/services/user_service.py":      "templates/app/services/user_service.py.tmpl",
		"app/services/item_service.py":      "templates/app/services/item_service.py.tmpl",
		"tests/test_users.py":               "templates/tests/test_users.py.tmpl",
		"tests/test_items.py":               "templates/tests/test_items.py.tmpl",
		"requirements.txt":                  "templates/requirements.txt.tmpl",
		".env":                              "templates/.env.tmpl",
		"README.md":                         "templates/README.md.tmpl",
	}
}

// rawFileMap maps destination path -> embed path for files copied verbatim (no template execution)
func rawFileMap() map[string]string {
	return map[string]string{
		"migrations/script.py.mako": "templates/migrations/script.py.mako.tmpl",
	}
}

func CreateProject(cfg ProjectConfig) error {
	for dest, tmplPath := range fileMap() {
		if err := writeTemplate(cfg, dest, tmplPath); err != nil {
			return err
		}
	}
	for dest, src := range rawFileMap() {
		if err := copyRaw(cfg.ProjectName, dest, src); err != nil {
			return err
		}
	}
	return nil
}

func writeTemplate(cfg ProjectConfig, dest, tmplPath string) error {
	outPath := filepath.Join(cfg.ProjectName, filepath.FromSlash(dest))
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	tmpl, err := template.ParseFS(templateFS, tmplPath)
	if err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, cfg)
}

func copyRaw(projectName, dest, src string) error {
	outPath := filepath.Join(projectName, filepath.FromSlash(dest))
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	in, err := templateFS.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
