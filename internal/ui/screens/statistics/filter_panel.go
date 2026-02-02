package statistics

import (
	"github.com/gurselcakar/arithmego/internal/analytics"
)

// FilterPanelModel manages filter state for statistics views.
type FilterPanelModel struct {
	// Current selected indices
	categoryIdx   int
	difficultyIdx int
	timePeriodIdx int

	// Available options
	categories   []string
	difficulties []string
	timePeriods  []analytics.TimePeriod
}

// NewFilterPanel creates a new filter panel model.
func NewFilterPanel() FilterPanelModel {
	return FilterPanelModel{
		categories:   analytics.AllCategories(),
		difficulties: analytics.AllDifficulties(),
		timePeriods:  analytics.AllTimePeriods(),
	}
}

// CycleCategory cycles to the next category.
func (m *FilterPanelModel) CycleCategory() {
	m.categoryIdx = (m.categoryIdx + 1) % len(m.categories)
}

// CycleDifficulty cycles to the next difficulty.
func (m *FilterPanelModel) CycleDifficulty() {
	m.difficultyIdx = (m.difficultyIdx + 1) % len(m.difficulties)
}

// CycleTimePeriod cycles to the next time period.
func (m *FilterPanelModel) CycleTimePeriod() {
	m.timePeriodIdx = (m.timePeriodIdx + 1) % len(m.timePeriods)
}

// GetCategoryDisplay returns the display name for current category.
func (m FilterPanelModel) GetCategoryDisplay() string {
	if m.categoryIdx >= 0 && m.categoryIdx < len(m.categories) {
		return analytics.CategoryDisplayName(m.categories[m.categoryIdx])
	}
	return "All Categories"
}

// GetDifficultyDisplay returns the display name for current difficulty.
func (m FilterPanelModel) GetDifficultyDisplay() string {
	if m.difficultyIdx >= 0 && m.difficultyIdx < len(m.difficulties) {
		return analytics.DifficultyDisplayName(m.difficulties[m.difficultyIdx])
	}
	return "All Difficulties"
}

// GetTimePeriodDisplay returns the display name for current time period.
func (m FilterPanelModel) GetTimePeriodDisplay() string {
	if m.timePeriodIdx >= 0 && m.timePeriodIdx < len(m.timePeriods) {
		return m.timePeriods[m.timePeriodIdx].String()
	}
	return "All Time"
}

// GetFilters returns the current filter configuration.
func (m FilterPanelModel) GetFilters() analytics.AggregateFilter {
	var category, difficulty string
	var timePeriod analytics.TimePeriod

	if m.categoryIdx >= 0 && m.categoryIdx < len(m.categories) {
		category = m.categories[m.categoryIdx]
	}
	if m.difficultyIdx >= 0 && m.difficultyIdx < len(m.difficulties) {
		difficulty = m.difficulties[m.difficultyIdx]
	}
	if m.timePeriodIdx >= 0 && m.timePeriodIdx < len(m.timePeriods) {
		timePeriod = m.timePeriods[m.timePeriodIdx]
	}

	return analytics.AggregateFilter{
		Category:   category,
		Difficulty: difficulty,
		TimePeriod: timePeriod,
	}
}
