package extedit

// Diff represents the output of an extedit session.
type Diff struct {
	content     Content
	Differences []int // Lines containing between session input and session output.
}

// Content returns session output as a single string.
func (d Diff) Content() string {
	return d.content.String()
}

// Lines returns session output split by newline.
func (d Diff) Lines() []string {
	return d.content.c
}

// Line of session output at a given index.
func (d Diff) Line(i int) string {
	return d.content.c[i]
}

// NewDiff constructs a new Diff.
func NewDiff(a, b Content) Diff {
	d := Diff{content: b}

	for i, line := range b.c {
		if len(a.c) <= i || line != a.c[i] {
			d.Differences = append(d.Differences, i)
		}
	}

	return d
}
