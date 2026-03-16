package external

import (
	"encoding/json"
	"os"

	"logCheckLinter/domain"
)

func LoadConfig(path string) (*domain.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg domain.Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
