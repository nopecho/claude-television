package channel

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Registry struct {
	Channels  []Channel `json:"channels"`
	UpdatedAt time.Time `json:"updated_at"`
}

func LoadRegistry(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Registry{}, nil
		}
		return nil, fmt.Errorf("read registry: %w", err)
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("parse registry: %w", err)
	}
	return &reg, nil
}

func SaveRegistry(reg *Registry, path string) error {
	reg.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal registry: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func (r *Registry) FindByID(id string) *Channel {
	for i := range r.Channels {
		if r.Channels[i].ID == id {
			return &r.Channels[i]
		}
	}
	return nil
}

func (r *Registry) Add(ch Channel) {
	r.Channels = append(r.Channels, ch)
}

func (r *Registry) Remove(id string) {
	filtered := make([]Channel, 0, len(r.Channels))
	for _, ch := range r.Channels {
		if ch.ID != id {
			filtered = append(filtered, ch)
		}
	}
	r.Channels = filtered
}
