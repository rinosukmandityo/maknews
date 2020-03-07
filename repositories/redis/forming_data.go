package redis

func FormingData(data map[string]string) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range data {
		res[k] = v
	}
	return res
}
