package fizzbuzz

import (
	"context"
	"strconv"
)

func PtrController(ctx context.Context, req *ControllerReq) (*ControllerResp, error) {
	res := make(map[int]string, req.To-req.From)
	for i := req.From; i < req.To; i++ {
		x, err := ptrLogic(ctx, &logicReq{i})
		if err != nil {
			return nil, err
		}

		res[i] = x.value
	}

	return &ControllerResp{res}, nil
}

func ptrLogic(ctx context.Context, req *logicReq) (*logicResp, error) {
	var (
		divisibleBy3 = req.value%3 == 0
		divisibleBy5 = req.value%5 == 0
	)
	switch {
	case divisibleBy3 && divisibleBy5:
		return &logicResp{"fizzbuzz"}, nil
	case divisibleBy3:
		return &logicResp{"fizz"}, nil
	case divisibleBy5:
		return &logicResp{"buzz"}, nil
	default:
		return &logicResp{strconv.FormatInt(int64(req.value), 10)}, nil
	}
}
