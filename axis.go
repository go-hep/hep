package dao

const (
	UnderflowBin = -2
	OverflowBin  = -1
)

type AxisKind int

const (
	FixedBinning    AxisKind = 0
	VariableBinning AxisKind = 1
)

type Axis interface {
	Kind() AxisKind
	LowerEdge() float64
	UpperEdge() float64
	Bins() int
	BinLowerEdge(idx int) float64
	BinUpperEdge(idx int) float64
	BinWidth(idx int) float64
	CoordToIndex(coord float64) int
}

// EOF
