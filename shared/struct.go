package shared



// Request and Response Structures
type StructureReq struct {
	OperationType string
	Mat1         [][]float64
	Mat2         [][]float64
}

type StructureResponse struct {
	Res [][]float64
	Worker string
}