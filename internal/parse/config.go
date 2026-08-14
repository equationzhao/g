package parse

import (
	"fmt"
	"os"

	"github.com/Equationzhao/g/internal/request"
	"gopkg.in/yaml.v3"
)

// fileYAML is the structured config surface (same fields as Request).
type fileYAML struct {
	Format *string `yaml:"format"`
	Long   *bool   `yaml:"long"`
}

// ApplyConfigFile reads PATH when ConfigSrc is ConfigPath. Missing file is exit 2.
func ApplyConfigFile(r request.Request) (request.Request, error) {
	if r.Config != request.ConfigPath {
		return r, nil
	}
	b, err := os.ReadFile(r.ConfigPath)
	if err != nil {
		return r, &Error{Msg: fmt.Sprintf("config %s: %v", r.ConfigPath, err), Code: 2}
	}
	var y fileYAML
	if err := yaml.Unmarshal(b, &y); err != nil {
		return r, &Error{Msg: fmt.Sprintf("config %s: %v", r.ConfigPath, err), Code: 2}
	}
	if y.Format != nil && !r.FormatSet {
		f, err := parseFormat(*y.Format)
		if err != nil {
			return r, err
		}
		r.Format = f
		r.FormatSet = true
	}
	if y.Long != nil && !r.Long {
		r.Long = *y.Long
	}
	return r, nil
}
