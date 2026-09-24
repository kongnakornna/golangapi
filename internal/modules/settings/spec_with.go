package settings

// FixedFilter is an unconditional WHERE clause (used by endpoint variants).
type FixedFilter struct {
	Column string
	Value  string
}

// With returns a copy of the spec with extra fixed filters appended.
func (s ListSpec) With(fixed ...FixedFilter) ListSpec {
	c := s
	c.FixedFilters = append(append([]FixedFilter{}, s.FixedFilters...), fixed...)
	return c
}
