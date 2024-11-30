package stats

func periodise[V any](series []V, getTimestamp func(V) string) [][]V {
	var bins [][]V
	for i, e := range series {
		if (i == 0) {
			bins = append(bins, *new ([]V))
			bins[len(bins) - 1] = append(bins[len(bins) - 1], e)
			continue
		}
		thisDate := getTimestamp(e)
		previousDate := getTimestamp(series[i - 1])
		if (thisDate[:7] != previousDate[:7]) {
			bins = append(bins, *new ([]V))
		}
		bins[len(bins) - 1] = append(bins[len(bins) - 1], e)
	}
	return bins
}
