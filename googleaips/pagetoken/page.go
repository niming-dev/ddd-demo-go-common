package pagetoken

import (
	"bytes"
	"encoding/base64"

	"github.com/niming-dev/ddd-demo/apis/golang/pagination/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type (
	listRequestBase interface {
		GetPageSize() int32
		GetPageToken() string
	}

	listRequest interface {
		listRequestBase
		GetParent() string
		GetFilter() string
		GetOrderBy() string
	}

	listRequestWithoutParent interface {
		listRequestBase
		GetFilter() string
		GetOrderBy() string
	}

	listRequestWithoutParentFilter interface {
		listRequestBase
		GetOrderBy() string
	}
)

// PageToken 页令牌
type PageToken struct {
	*pagination.Pagination
}

// InjectListRequest 从request中注入数据
func (p PageToken) InjectListRequest(req listRequestBase) {
	p.Pagination.Limit = req.GetPageSize()
	if raw, ok := req.(listRequestWithoutParentFilter); ok {
		p.Pagination.OrderBy = raw.GetOrderBy()
	}
	if raw, ok := req.(listRequestWithoutParent); ok {
		p.Pagination.Filter = raw.GetFilter()
	}
	if raw, ok := req.(listRequest); ok {
		p.Pagination.Parent = raw.GetParent()
	}
}

// String 转为字符串
func (p PageToken) String() string {
	bs, _ := proto.Marshal(p.Pagination)
	return base64.RawURLEncoding.EncodeToString(bs)
}

// Equal 判断两个pageToken是否相等
func (p PageToken) Equal(pt PageToken) bool {
	p1, _ := proto.Marshal(p.Pagination)
	p2, _ := proto.Marshal(pt)
	return bytes.Equal(p1, p2)
}

// CompareReq .
func (p PageToken) CompareReq(req listRequestBase) bool {
	if raw, ok := req.(listRequestWithoutParentFilter); ok {
		if p.GetOrderBy() != raw.GetOrderBy() {
			return false
		}
	}
	if raw, ok := req.(listRequestWithoutParent); ok {
		if p.GetFilter() != raw.GetFilter() {
			return false
		}
	}
	if raw, ok := req.(listRequest); ok {
		if p.GetParent() != raw.GetParent() {
			return false
		}
	}
	return true
}

// Next 返回下一页的token
func (p PageToken) Next() PageToken {
	return PageToken{
		Pagination: &pagination.Pagination{
			Offset:  p.GetOffset() + p.GetLimit(),
			Limit:   p.GetLimit(),
			OrderBy: p.GetOrderBy(),
			Filter:  p.GetFilter(),
			Parent:  p.GetParent(),
		},
	}
}

// New 新建一个PageToken实例
func New(offset, limit int32, opts ...Option) PageToken {
	pt := PageToken{
		&pagination.Pagination{
			Offset: offset,
			Limit:  limit,
		},
	}

	return checkAndReset(pt, evaluateOpt(opts))
}

// NewFromReq 根据请求内容创建PageToken
func NewFromReq(req listRequestBase, opts ...Option) (PageToken, error) {
	o := evaluateOpt(opts)
	limit := int32(o.defaultLimit)
	if p := req.GetPageSize(); p > 0 {
		limit = p
	} else if p > int32(o.maxLimit) {
		limit = int32(o.maxLimit)
	}

	// 根据参数生成pageToken
	var pageToken PageToken
	if req.GetPageToken() != "" {
		pt, err := ParsePageToken(req.GetPageToken(), opts...)
		if err != nil {
			return PageToken{}, status.Error(codes.InvalidArgument, "Invalid argument [PageToken]")
		}
		pageToken = pt
		pageToken.Pagination.Limit = limit
	} else {
		pageToken = New(0, req.GetPageSize())
		pageToken.InjectListRequest(req)
		pageToken.Pagination.Limit = limit
		checkAndReset(pageToken, evaluateOpt(opts))
	}

	if !pageToken.CompareReq(req) {
		return PageToken{}, status.Error(codes.InvalidArgument, "Invalid argument [PageToken]")
	}

	return pageToken, nil
}

// ParsePageToken 解析PageToken
func ParsePageToken(token string, opts ...Option) (PageToken, error) {
	rawStr, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return PageToken{}, err
	}

	pt := PageToken{Pagination: new(pagination.Pagination)}
	err = proto.Unmarshal(rawStr, pt.Pagination)
	if err != nil {
		return PageToken{}, err
	}

	return checkAndReset(pt, evaluateOpt(opts)), nil
}

func checkAndReset(pt PageToken, o *options) PageToken {
	if pt.Pagination == nil {
		return pt
	}
	if o.defaultLimit > 0 && int(pt.Pagination.Limit) == 0 {
		pt.Pagination.Limit = int32(o.defaultLimit)
	}
	if o.maxLimit > 0 && int(pt.Pagination.Limit) > o.maxLimit {
		pt.Pagination.Limit = int32(o.maxLimit)
	}

	if o.defaultOrderBy != "" && pt.Pagination.OrderBy == "" {
		pt.Pagination.OrderBy = o.defaultOrderBy
	}

	return pt
}
