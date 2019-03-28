package openapi

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// sign with sha 256
func Sign(content, key string) string {
	h := sha256.New()
	h.Write([]byte(content + key))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Verify(result, content, key string) bool {
	signed := Sign(content, key)
	return strings.EqualFold(signed, result)
}

type KvPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Pairs []KvPair

func (p Pairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Pairs) Len() int {
	return len(p)
}
func (p Pairs) Less(i, j int) bool {
	return p[i].Key < p[j].Key
}

// 防止存放攻击, 客户端需要显式的传过来时间戳
func BuildParams(params Pairs) string {
	sort.Sort(params)
	var result string
	for _, v := range params {
		r := v.Key + "=" + v.Value + "&"
		result += r
	}
	return result
}
