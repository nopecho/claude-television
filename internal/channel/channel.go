package channel

import (
	"time"

	"github.com/nopecho/claude-television/internal/claude"
)

type ChannelStatus string

const (
	StatusHealthy ChannelStatus = "healthy"
	StatusWarning ChannelStatus = "warning"
	StatusError   ChannelStatus = "error"
)

type Channel struct {
	ID         string        `json:"id"`
	Path       string        `json:"path"`
	Name       string        `json:"name"`
	Pinned     bool          `json:"pinned"`
	Group      string        `json:"group"`
	Status     ChannelStatus `json:"status"`
	LastSynced time.Time     `json:"last_synced"`
	Data       *ChannelData  `json:"-"`
}

type ChannelData struct {
	Settings      *claude.Settings    `json:"settings"`
	LocalSettings *claude.Settings    `json:"local_settings"`
	ClaudeMD      *claude.ClaudeMD    `json:"claude_md"`
	SubClaudeMDs  []claude.ClaudeMD   `json:"sub_claude_mds"`
	Hooks         []claude.HookDetail `json:"hooks"`
	MCPServers    []claude.MCPServer  `json:"mcp_servers"`
	Plugins       []claude.Plugin     `json:"plugins"`
	LocalSkills   []claude.Skill      `json:"local_skills"`
	GitInfo       *GitInfo            `json:"git_info"`
	MemoryFiles   []claude.MemoryFile `json:"memory_files"`
}

type GitInfo struct {
	Branch        string `json:"branch"`
	LastCommit    string `json:"last_commit"`
	LastCommitMsg string `json:"last_commit_msg"`
	LastCommitAt  string `json:"last_commit_at"`
	DirtyFiles    int    `json:"dirty_files"`
}
