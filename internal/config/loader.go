package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

const (
	varDir    = "_vars"
	delimeter = "."
)

var supportedExts = []string{".yaml", ".yml", ".json", ".toml"}

var ErrUnsupportedFileExt = errors.New("unsupported config file extension")

type loaderOption func(*Loader)

type Loader struct {
	fsys    fs.FS
	rootDir string
	varsDir string
	k       *koanf.Koanf
	vars    *koanf.Koanf
}

func NewLoader(root string, opts ...loaderOption) *Loader {
	loader := &Loader{
		fsys:    os.DirFS("."),
		rootDir: root,
		varsDir: path.Join(root, varDir),
		k:       koanf.New(delimeter),
		vars:    koanf.New(delimeter),
	}

	for _, opt := range opts {
		opt(loader)
	}

	return loader
}

func WithFS(fsys fs.FS) loaderOption {
	return func(l *Loader) {
		l.fsys = fsys
	}
}

func WithVarsDir(varsDir string) loaderOption {
	return func(l *Loader) {
		l.varsDir = path.Join(l.rootDir, varsDir)
	}
}

func (l *Loader) Load() (*Config, error) {
	if err := l.loadFiles(l.varsDir, l.vars, ""); err != nil {
		return nil, fmt.Errorf("loading variables: %w", err)
	}

	if err := l.loadFiles(l.rootDir, l.k, l.varsDir); err != nil {
		return nil, fmt.Errorf("loading config files: %w", err)
	}

	var cfg Config
	if err := l.k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}
	return &cfg, nil
}

func (l *Loader) loadFiles(dir string, k *koanf.Koanf, skipDir string) error {
	if _, err := fs.Stat(l.fsys, dir); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	files, err := collectConfigFiles(l.fsys, dir, skipDir)
	if err != nil {
		return fmt.Errorf("collecting files from %s: %w", dir, err)
	}

	for _, file := range files {
		if err := l.loadFile(file, k); err != nil {
			return fmt.Errorf("loading file %s: %w", file, err)
		}
	}
	return nil
}

func (l *Loader) loadFile(file string, k *koanf.Koanf) error {
	rendered, p, err := l.readAndRenderFile(file)
	if err != nil {
		return fmt.Errorf("reading file %q: %w", file, err)
	}

	if err := k.Load(rawbytes.Provider(rendered), p); err != nil {
		return fmt.Errorf("loading into config map: %w", err)
	}
	return nil
}

func (l *Loader) readAndRenderFile(file string) ([]byte, koanf.Parser, error) {
	content, err := fs.ReadFile(l.fsys, file)
	if err != nil {
		return nil, nil, fmt.Errorf("reading file: %w", err)
	}

	rendered, err := l.renderTemplate(file, content)
	if err != nil {
		return nil, nil, fmt.Errorf("rendering tempalte: %w", err)
	}

	p, err := parser(file)
	if err != nil {
		return nil, nil, fmt.Errorf("getting parser: %w", err)
	}
	return rendered, p, nil
}

func (l *Loader) renderTemplate(name string, content []byte) ([]byte, error) {
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, l.vars.All()); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}
	return buf.Bytes(), nil
}

func collectConfigFiles(fsys fs.FS, dir, skipDir string) ([]string, error) {
	var files []string
	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && path == skipDir {
			return fs.SkipDir
		}

		if isConfigFile(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	slices.SortFunc(files, func(a, b string) int {
		depthA, depthB := strings.Count(a, "/"), strings.Count(b, "/")
		if depthA != depthB {
			return depthA - depthB
		}
		return strings.Compare(a, b)
	})
	return files, nil
}

func parser(file string) (koanf.Parser, error) {
	switch filepath.Ext(strings.ToLower(file)) {
	case ".yaml", ".yml":
		return yaml.Parser(), nil
	case ".json":
		return json.Parser(), nil
	case ".toml":
		return toml.Parser(), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileExt, filepath.Ext(file))
}

func isConfigFile(file string) bool {
	return slices.Contains(supportedExts, filepath.Ext(strings.ToLower(file)))
}
