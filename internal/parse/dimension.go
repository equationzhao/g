package parse

// Dimension is one axis of the 12-cell interaction matrix in docs/rewrite-architecture.md §17.2.
type Dimension uint8

const (
	DimMeta Dimension = iota
	DimFormat
	DimLong
	DimVisibility
	DimWalk
	DimFilter
	DimSort
	DimSize
	DimTime
	DimPresent
	DimGit
	DimDeref
	dimCount
)

// Budget is the hard cap on primary flags. Changing this value is a Key Decision, not a drive-by.
const Budget = 40

func (d Dimension) String() string {
	if int(d) < len(dimNames) {
		return dimNames[d]
	}
	return "?"
}

func parseDimension(s string) (Dimension, bool) {
	for i, n := range dimNames {
		if n == s {
			return Dimension(i), true
		}
	}
	return 0, false
}

var dimNames = [...]string{
	DimMeta:       "Meta",
	DimFormat:     "Format",
	DimLong:       "Long",
	DimVisibility: "Visibility",
	DimWalk:       "Walk",
	DimFilter:     "Filter",
	DimSort:       "Sort",
	DimSize:       "Size",
	DimTime:       "Time",
	DimPresent:    "Present",
	DimGit:        "Git",
	DimDeref:      "Deref",
}

// Relation is one of the five interaction classes. Every dimension pair and every
// same-dimension flag pair must resolve to exactly one.
type Relation uint8

const (
	RelOrthogonal Relation = iota
	RelSameDim
	RelSuppresses
	RelError
	RelIgnored
)

func (r Relation) String() string {
	if int(r) < len(relNames) {
		return relNames[r]
	}
	return "?"
}

func parseRelation(s string) (Relation, bool) {
	for i, n := range relNames {
		if n == s {
			return Relation(i), true
		}
	}
	return 0, false
}

var relNames = [...]string{
	RelOrthogonal: "orthogonal",
	RelSameDim:    "same-dim",
	RelSuppresses: "suppresses",
	RelError:      "error",
	RelIgnored:    "ignored",
}

// ConfigExempt explains why a primary flag has no YAML key.
type ConfigExempt uint8

const (
	ConfigRequired ConfigExempt = iota
	ConfigMeta                  // --help / --version / --config
	ConfigCLIOnly               // GNU muscle, not worth a YAML key
)

func (e ConfigExempt) String() string {
	switch e {
	case ConfigRequired:
		return ""
	case ConfigMeta:
		return "meta"
	case ConfigCLIOnly:
		return "cli-only"
	default:
		return "?"
	}
}
