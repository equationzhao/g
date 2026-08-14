package print

import "fmt"

func HumanSize(n int64, si bool) string {
	if n < 0 {
		n = 0
	}
	base := int64(1024)
	suf := [...]string{"B", "K", "M", "G", "T", "P"}
	if si {
		base = 1000
		suf = [...]string{"B", "k", "M", "G", "T", "P"}
	}
	if n < base {
		return fmt.Sprintf("%d", n)
	}
	v := float64(n)
	i := 0
	for v >= float64(base) && i < len(suf)-1 {
		v /= float64(base)
		i++
	}
	if v >= 10 {
		return fmt.Sprintf("%.0f%s", v, suf[i])
	}
	return fmt.Sprintf("%.1f%s", v, suf[i])
}
