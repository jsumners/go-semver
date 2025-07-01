package semver

// Compare evaluates the ordinality between two versions.
// Results:
//   - `a > b => 1`
//   - `a < b => -1`
//   - `a == b => 0`
func Compare(a *Version, b *Version) int {
	if a.major > b.major {
		return 1
	}
	if a.major < b.major {
		return -1
	}

	// Major versions are equal, what about minor versions?
	if a.minor > b.minor {
		return 1
	}
	if a.minor < b.minor {
		return -1
	}

	// Major and minor versions are equal, what about patch versions?
	if a.patch > b.patch {
		return 1
	}
	if a.patch < b.patch {
		return -1
	}

	// todo: compare pre-release tags

	// Major, minor, and patch versions are all equal so versions are equal.
	return 0
}
