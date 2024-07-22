package pointer

func True() *bool {
	res := true
	return &res
}

func False() *bool {
	res := false
	return &res
}

func Float32(v float32) *float32 {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

func Int(v int) *int {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Interface(v interface{}) *interface{} {
	return &v
}

func String(v string) *string {
	return &v
}
