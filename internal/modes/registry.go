package modes

var registry = make(map[string]*Mode)
var orderedIDs []string

// Register adds a mode to the registry.
func Register(m *Mode) {
	registry[m.ID] = m
	orderedIDs = append(orderedIDs, m.ID)
}

// Get retrieves a mode by ID.
func Get(id string) (*Mode, bool) {
	m, ok := registry[id]
	return m, ok
}

// All returns all registered modes in registration order.
func All() []*Mode {
	modes := make([]*Mode, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		modes = append(modes, registry[id])
	}
	return modes
}

// ByCategory returns modes filtered by category.
func ByCategory(cat ModeCategory) []*Mode {
	var modes []*Mode
	for _, id := range orderedIDs {
		m := registry[id]
		if m.Category == cat {
			modes = append(modes, m)
		}
	}
	return modes
}
