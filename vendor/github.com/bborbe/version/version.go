package version

type Version string

func (v Version) String() string {return string (v)}

const (
	LATEST = Version("latest")
)
