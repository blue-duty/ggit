package common

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// RemoveStringSlice 通过字符串切片2，去掉字符串切片1中的元素
func RemoveStringSlice(a []string, b []string) []string {
	// 将b装进map
	bMap := make(map[string]bool)
	for _, v := range b {
		bMap[v] = true
	}
	// 遍历a，如果a中的元素在b中，就删除
	for i := 0; i < len(a); i++ {
		if _, ok := bMap[a[i]]; ok {
			a = append(a[:i], a[i+1:]...)
			i--
		}
	}
	return a
}
