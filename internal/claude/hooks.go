package claude

type HookDetail struct {
	Event   string
	Matcher string
	Type    string
	Command string
	Async   bool
	Timeout int
	Source  string
}

func ExtractHooks(settings *Settings, source string) []HookDetail {
	if settings.Hooks == nil {
		return nil
	}
	var result []HookDetail
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
	return result
}
