package src

func Min(x uint64, y uint64) uint64 {
	if x < y {
		return x
	} else {
		return y
	}
}

func Max(x uint64, y uint64) uint64 {
	if x < y {
		return y
	} else {
		return x
	}
}
