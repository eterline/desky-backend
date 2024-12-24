package utils

func RemoveSliceIndex[Type any](s []Type, index int) []Type {
	return append(s[:index], s[index+1:]...)
}
