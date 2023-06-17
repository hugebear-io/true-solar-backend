package helper

func ChunkBy(items map[string]map[string]interface{}, chunkSize int) (chunks [][]map[string]map[string]interface{}) {
	slice := make([]map[string]map[string]interface{}, 0)

	for k, v := range items {
		slice = append(slice, map[string]map[string]interface{}{k: v})
	}

	for chunkSize < len(slice) {
		slice, chunks = slice[chunkSize:], append(chunks, slice[0:chunkSize:chunkSize])
	}

	return append(chunks, slice)
}
