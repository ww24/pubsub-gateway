package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidActionType = errors.New("invalid action type")
)

// Config represents a config.yml.
type Config struct {
	Version  string         `yaml:"version"`
	Handlers []EventHandler `yaml:"handlers"`
}

// EventHandler is event handler which contains action names.
type EventHandler struct {
	Subscription string `yaml:"subscription"`
	Action       Action `yaml:"action"`
}

// ActionType represents action type.
type ActionType string

const (
	// ActionNone is uncategorized action type.
	ActionNone ActionType = ""
	// ActionHTTP is action type for HTTP action.
	ActionHTTP ActionType = "http"
)

// Validate validates action type.
func (t ActionType) Validate() error {
	switch t {
	case ActionHTTP:
		return nil

	case ActionNone:
		return ErrInvalidActionType

	default:
		return ErrInvalidActionType
	}
}

// Action is action definition.
type Action struct {
	Type              ActionType `yaml:"type"`
	HTTPRequestAction `yaml:",inline"`
}

// HTTPRequestAction is configuration of HTTP action.
type HTTPRequestAction struct {
	Method  string        `yaml:"method"`
	Header  http.Header   `yaml:"header"`
	URL     string        `yaml:"url"`
	Payload *yamlTemplate `yaml:"payload"`
}

type yamlTemplate struct {
	tmpl *template.Template
	res  gjson.Result
}

func (t *yamlTemplate) UnmarshalYAML(node *yaml.Node) error {
	b := &bytes.Buffer{}
	e := yaml.NewEncoder(b)
	defer e.Close()
	if err := e.Encode(node); err != nil {
		return err
	}
	funcs := template.FuncMap{
		"path": func(path string) string { return t.res.Get(path).String() },
	}
	tmpl, err := template.New("tmpl").Delims("${", "}").Funcs(funcs).Parse(b.String())
	if err != nil {
		return err
	}
	t.tmpl = tmpl
	return nil
}

func (t *yamlTemplate) Render(jsonData []byte) ([]byte, error) {
	t.res = gjson.ParseBytes(jsonData)
	b := &bytes.Buffer{}
	if err := t.tmpl.Execute(b, nil); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (t *yamlTemplate) RenderMap(jsonData []byte) (map[string]interface{}, error) {
	y, err := t.Render(jsonData)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(y, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (t *yamlTemplate) RenderJSON(jsonData []byte) ([]byte, error) {
	m, err := t.RenderMap(jsonData)
	if err != nil {
		return nil, err
	}
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// Parse parses config file and returns config data.
func Parse(data []byte) (Config, error) {
	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	if err := config.Validate(); err != nil {
		return Config{}, fmt.Errorf("validation error: %w", err)
	}

	return config, nil
}

// Validate validates config.
func (c *Config) Validate() error {
	if c.Version == "" {
		return errors.New("version is required")
	}
	if len(c.Handlers) == 0 {
		return errors.New("handler should be defined one or more")
	}
	for _, a := range c.Handlers {
		if a.Subscription == "" {
			return errors.New("subscription should be defined")
		}
		if err := a.Action.Type.Validate(); err != nil {
			return fmt.Errorf("action type validation failed: %w", err)
		}
	}
	return nil
}
