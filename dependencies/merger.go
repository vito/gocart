package dependencies

import (
	"github.com/xoebus/gocart/dependency"
)

func Merge(cartridgeList, lockList []dependency.Dependency) []dependency.Dependency {
	merged := []dependency.Dependency{}
	versions := map[string]int{}

	for i, c := range cartridgeList {
		merged = append(merged, c)
		versions[c.Path] = i
	}

	for _, l := range lockList {
		existing, found := versions[l.Path]
		if found {
			merged[existing].Version = l.Version
		}
	}

	return merged
}
