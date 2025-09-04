package version

import (
	"sort"
	"strconv"
	"strings"
)

// ParseVer converts version string (e.g., "1.10") to float64
func ParseVer(v string) float64 {
	parts := strings.Split(v, ".")
	if len(parts) < 2 {
		f, _ := strconv.ParseFloat(v, 64)
		return f
	}
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	return float64(major) + float64(minor)/100
}

// MatchesVersion checks if version matches constraint (e.g., ">=1.10")
func MatchesVersion(constraint, ver string) bool {
	if constraint == "" {
		return true
	}

	op := ""
	val := constraint
	if strings.HasPrefix(constraint, ">=") {
		op = ">="
		val = strings.TrimPrefix(constraint, ">=")
	} else if strings.HasPrefix(constraint, "<=") {
		op = "<="
		val = strings.TrimPrefix(constraint, "<=")
	} else if strings.HasPrefix(constraint, ">") {
		op = ">"
		val = strings.TrimPrefix(constraint, ">")
	} else if strings.HasPrefix(constraint, "<") {
		op = "<"
		val = strings.TrimPrefix(constraint, "<")
	} else {
		op = "="
		val = constraint
	}

	v1 := ParseVer(ver)
	v2 := ParseVer(val)

	switch op {
	case ">=":
		return v1 >= v2
	case "<=":
		return v1 <= v2
	case ">":
		return v1 > v2
	case "<":
		return v1 < v2
	case "=":
		return v1 == v2
	}
	return false
}

// SortVersions sorts version strings in ascending order
func SortVersions(ver []string) {
	sort.Slice(ver, func(i, j int) bool {
		return ParseVer(ver[i]) < ParseVer(ver[j])
	})
}
