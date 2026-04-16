package unique

// UniqueStrings 字符数组去重
func UniqueStrings(strs []string) []string {
	m := make(map[string]struct{})
	res := make([]string, 0, len(strs))
	for _, v := range strs {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			res = append(res, v)
		}
	}
	return res
}
