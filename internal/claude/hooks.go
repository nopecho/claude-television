package claude

type HookDetail struct {
	Event   string `json:"event"`
	Matcher string `json:"matcher"`
	Type    string `json:"type"`
	Command string `json:"command"`
	Async   bool   `json:"async"`
	Timeout int    `json:"timeout"`
	Source  string `json:"source"`
}

func ExtractHooks(settings *Settings, source string) ([]HookDetail, error) {
	if settings == nil || settings.Hooks == nil {
		return []HookDetail{}, nil
	}
	result := []HookDetail{}
	for event, rules := range settings.Hooks {
		for _, rule := range rules {
			for _, action := range rule.Hooks {
				result = append(result, HookDetail{
					Event: event, Matcher: rule.Matcher,
					Type: action.Type, Command: action.Command,
					Async: action.Async, Timeout: action.Timeout,
					Source: source,
				})
			}
		}
	}
	return result, nil
}
