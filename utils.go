package shell

func prependString(x []string, y string) []string {
	x = append(x, "")
	copy(x[1:], x)
	x[0] = y
	return x
}
