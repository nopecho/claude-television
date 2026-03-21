package cmd

import (
	"os"
	"sync"

	"github.com/nopecho/claude-television/internal/channel"
)

func loadAllChannels(channels []channel.Channel, claudeHome string, cache *channel.Cache) {
	var wg sync.WaitGroup
	for i := range channels {
		wg.Add(1)
		go func(ch *channel.Channel) {
			defer wg.Done()

			if cache != nil {
				expected := channel.ExpectedFiles(ch)
				if entry, valid := cache.LoadIfValid(ch.ID, expected); valid {
					ch.Data = &entry.Data
					return
				}
			}

			data, mtimes, err := channel.LoadChannelData(ch, claudeHome)
			if err != nil {
				ch.Status = channel.StatusError
				return
			}
			ch.Data = data
			ch.Status = determineStatus(ch, data)

			if cache != nil {
				cache.Save(&channel.CacheEntry{
					ChannelID:  ch.ID,
					Data:       *data,
					FileMtimes: mtimes,
				})
			}
		}(&channels[i])
	}
	wg.Wait()
}

func determineStatus(ch *channel.Channel, data *channel.ChannelData) channel.ChannelStatus {
	if _, err := os.Stat(ch.Path); err != nil {
		return channel.StatusError
	}
	if data.Settings == nil && data.ClaudeMD == nil {
		return channel.StatusWarning
	}
	return channel.StatusHealthy
}
