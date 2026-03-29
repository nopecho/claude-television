package tui

import (
	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
)

func defaultCfg() *config.Config {
	return &config.Config{
		Channels: config.ChannelsConfig{
			Groups: map[string][]string{},
		},
	}
}

func makeChannel(id, name string, opts ...func(*channel.Channel)) channel.Channel {
	ch := channel.Channel{
		ID:     id,
		Name:   name,
		Path:   "/tmp/" + name,
		Status: channel.StatusHealthy,
		Data:   &channel.ChannelData{},
	}
	for _, o := range opts {
		o(&ch)
	}
	return ch
}

func withPinned(pinned bool) func(*channel.Channel) {
	return func(ch *channel.Channel) { ch.Pinned = pinned }
}

func withGroup(group string) func(*channel.Channel) {
	return func(ch *channel.Channel) { ch.Group = group }
}

func withGlobal() func(*channel.Channel) {
	return func(ch *channel.Channel) { ch.IsGlobal = true }
}
