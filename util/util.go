package util

func Round(f float64) float64 {
	if f == 0 {
		return 0
	}
	d := 0.005
	if f < 0 {
		d = -0.005
	}
	return float64(int64((f+d)*100)) / 100
}
