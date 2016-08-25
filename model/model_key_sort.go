package model

type KeySlice []Key

func (k KeySlice) Len() int           { return len(k) }
func (k KeySlice) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k KeySlice) Less(i, j int) bool { return k[i] < k[j] }
