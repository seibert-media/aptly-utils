package key

type KeySlice []Key

func (v KeySlice) Len() int           { return len(v) }
func (v KeySlice) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (p KeySlice) Less(i, j int) bool { return p[i] < p[j] }
