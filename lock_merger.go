package gocart

func MergeDependencies(cartridgeList, lockList []Dependency) []Dependency {
	merged := []Dependency{}
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
