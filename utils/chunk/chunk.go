package chunk

import "strings"

func Chunk(items []string, size int) []string {
	var chunk []string
	var j int
	for i := 0; i < len(items); i += size {
		j += size
		if j > len(items) {
			j = len(items)
		}
		chunk = append(chunk, strings.Join(items[i:j], ","))
	}
	return chunk
}
