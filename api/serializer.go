package api

import (
	slz "github.com/rinosukmandityo/maknews/api/serializer"
	js "github.com/rinosukmandityo/maknews/api/serializer/json"
	ms "github.com/rinosukmandityo/maknews/api/serializer/msgpack"
)

var (
	ContentTypeJson    = "application/json"
	ContentTypeMsgPack = "application/x-msgpack"
)

func GetSerializer(contentType string) slz.UserSerializer {
	if contentType == ContentTypeMsgPack {
		return &ms.News{}
	}
	return &js.News{}
}
