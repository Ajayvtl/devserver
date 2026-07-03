package tasks

// Progress describes a single progress update for a task.
type Progress struct {
	Percent int
	Detail  string
}

// normalize clamps the percentage and keeps the detail string intact.
func (p Progress) normalize() Progress {
	if p.Percent < 0 {
		p.Percent = 0
	}
	if p.Percent > 100 {
		p.Percent = 100
	}
	return p
}
