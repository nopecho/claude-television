package tui

type TabID int

const (
	TabGlobal TabID = iota
	TabProjects
	TabSkills
	TabHooks
)

var tabNames = []string{"Global", "Projects", "Skills", "Hooks"}

func renderTabs(active TabID) string {
	var tabs string
	for i, name := range tabNames {
		if TabID(i) == active {
			tabs += activeTabStyle.Render("[" + name + "]")
		} else {
			tabs += inactiveTabStyle.Render(" " + name + " ")
		}
	}
	return tabs
}
