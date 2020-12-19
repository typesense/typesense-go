package pointer

func True() *bool {
	res := true
	return &res
}

func False() *bool {
	res := false
	return &res
}

func Int(v int) *int {
	return &v
}

func Interface(v interface{}) *interface{} {
	return &v
}

func String(v string) *string {
	return &v
}
