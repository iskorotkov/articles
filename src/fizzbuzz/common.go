package fizzbuzz

type ControllerReq struct {
	From, To int
}

type ControllerResp struct {
	Values map[int]string
}

type logicReq struct {
	value int
}

type logicResp struct {
	value string
}
