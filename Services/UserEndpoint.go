package Services

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"gokitdemo/util"
	"strconv"
)

type UserRequest struct {
	Uid    int    `json:"uid"`
	Method string `json:"method"`
}

type UserResponse struct {
	Result string `json:"result"`
}

// 限流
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				//return nil, errors.NewError(429, "too many request")
				//return nil,errors.New("to many request")
				return nil, util.NewMyError(429, "too many requests")
			}
			return next(ctx, request)
		}

	}
}
func GetUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := "no string"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(util.ServicePort)
		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("userid为%d的用户删除成功“", r.Uid)
			}
		}

		return UserResponse{Result: result}, nil
	}
}
