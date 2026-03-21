package claude

import "encoding/json"

type MCPServer struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
	URL     string            `json:"url"`
	Source  string            `json:"source"`
}

func ExtractMCPServers(settings *Settings, source string) ([]MCPServer, error) {
	if settings == nil || settings.Raw == nil {
		return []MCPServer{}, nil
	}
	raw, ok := settings.Raw["mcpServers"]
	if !ok {
		return []MCPServer{}, nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var servers map[string]struct {
		Type    string            `json:"type"`
		Command string            `json:"command"`
		Args    []string          `json:"args"`
		Env     map[string]string `json:"env"`
		URL     string            `json:"url"`
	}
	if err := json.Unmarshal(data, &servers); err != nil {
		return nil, err
	}
	result := []MCPServer{}
	for name, s := range servers {
		result = append(result, MCPServer{
			Name: name, Type: s.Type, Command: s.Command,
			Args: s.Args, Env: s.Env, URL: s.URL, Source: source,
		})
	}
	return result, nil
}
