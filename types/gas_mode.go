package types

type GasMode int8

const (
	GasModeAuto GasMode = iota + 1
	GasModeLegacy
	GasModeEip1559
)

var (
	GasMode_name = map[int32]string{
		1: "auto",
		2: "legacy",
		3: "eip1559",
	}

	GasMode_value = map[string]int32{
		"auto":    1,
		"legacy":  2,
		"eip1559": 3,
	}
)
