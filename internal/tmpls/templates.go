package tmpls

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

type Cache map[string]*template.Template

func LoadCache(fsys fs.FS) (Cache, error) {
	cache := make(Cache)

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		if path == "base.html" || path == "base_auth.html" {
			return nil
		}

		base := "base.html"
		if strings.HasPrefix(path, "auth/") {
			base = "base_auth.html"
		}

		t, err := template.ParseFS(fsys, base, path)
		if err != nil {
			return err
		}

		name := path
		cache[name] = t

		return nil
	})
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func (c Cache) Render(w http.ResponseWriter, r *http.Request, name string, data any) error {
	tmpl, ok := c[name]
	if !ok {
		return fmt.Errorf("template not found in cache: %s", name)
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	isBoosted := r.Header.Get("HX-Boosted") == "true"

	var err error
	var buf bytes.Buffer
	var layout string

	if isHTMX && !isBoosted {
		layout = "content"
		err = tmpl.ExecuteTemplate(&buf, layout, data)
	} else {
		layout = "base.html"
		if strings.HasPrefix(name, "auth/") {
			layout = "base_auth.html"
		}
		err = tmpl.ExecuteTemplate(&buf, layout, data)
	}

	if err != nil {
		return fmt.Errorf("error executing %s: %w", name, err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err = buf.WriteTo(w); err != nil {
		return fmt.Errorf("error writing to response writer: %w", err)
	}

	return nil
}
