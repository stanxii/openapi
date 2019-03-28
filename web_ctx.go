package openapi

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const defaultExpireTime = 60000

// to check if the request is valid  from the signing key
func CheckValid(req *http.Request, keeper SecretKeeper) (bool, error) {
	if req == nil {
		return false, errors.New("illegal request")
	}
	// time in millis
	timeStr := getParamFromRequest(req, "time")
	signResult := getParamFromRequest(req, "sign")

	rt, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return false, errors.New("error parameter")
	}
	now := time.Now().UnixNano() / 1e6
	duration := math.Abs(float64(rt - now))
	if duration > defaultExpireTime {
		return false, errors.New("error timestamp")
	}

	pairs := getPairs(req)
	content := BuildParams(pairs)
	secret, err := keeper.GetSecret()
	if err != nil {
		return false, err
	}
	result := Verify(signResult, content, secret)
	if result {
		return result, nil
	}
	return result, errors.New("error verifying")
}

func getParamFromRequest(req *http.Request, param string) string {
	if req == nil {
		return ""
	}
	return req.URL.Query().Get(param)
}

func getPairs(req *http.Request) Pairs {
	pairs := make([]KvPair, 0, 10)
	headers := req.Header
	headerPairs := getPairsFromMap(headers)
	// add all params
	paramsMap := req.URL.Query()
	paramPairs := getPairsFromMap(paramsMap)
	pairs = append(pairs, headerPairs...)
	pairs = append(pairs, paramPairs...)
	return pairs
}

// get params and headers except the param sign
func getPairsFromMap(m map[string][]string) Pairs {
	pairs := make([]KvPair, 0, 10)
	for k, v := range m {
		if len(k) < 1 {
			continue
		}
		var val string
		for _, e := range v {
			val += e
		}
		if strings.EqualFold(k, "sign") {
			continue
		}
		p := KvPair{
			Key:   k,
			Value: val,
		}
		pairs = append(pairs, p)
	}
	return pairs
}
