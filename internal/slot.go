package internal

type ValueSlot struct {
	Ts  int `json:"ts"`
	Val int `json:"val"`
}

type FloatSlot struct {
	Ts  int     `json:"ts"`
	Val float32 `json:"val"`
}

type DistinctSlot struct {
	Ts  int    `json:"ts"`
	Val string `json:"val"`
}

func getSlotSize(window int) int {

	if window > 0 && window <= 10 {
		return 1
	}
	if window > 10 && window <= 60 {
		return 2
	}
	if window > 60 && window <= 1800 {
		return 60
	}
	if window > 1800 && window <= 21600 {
		return 600
	}
	if window > 21600 && window <= 86400 {
		return 3600
	}

	return 3600
}
