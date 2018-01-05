package main

import (
	"encoding/base64"
	"fmt"

	"github.com/ugorji/go/codec"
)

var (
	msgpackC codec.MsgpackHandle
)

func processMsgpack(data []byte) (t int, o interface{}, err error) {
	fmt.Println("Base64:", base64.StdEncoding.EncodeToString(data))
	var tmp Wrapper
	codec.NewDecoderBytes(data, &msgpackC).Decode(&tmp)
	t = tmp.Type
	switch findPacket(t) {
	case "error":
		o := new(struct {
			Type    int
			Content struct {
				Content int
			}
		})
		err = codec.NewDecoderBytes(data, &msgpackC).Decode(&o)
		return t, o.Content.Content, err
	case "userReceive":
		o := new(struct {
			Type    int
			Content struct {
				Content UserReceive
			}
		})
		err = codec.NewDecoderBytes(data, &msgpackC).Decode(&o)
		return t, o.Content.Content, err
	case "loginSuccess":
		o := new(struct {
			Type    int
			Content struct {
				Content LoginSuccess
			}
		})
		err = codec.NewDecoderBytes(data, &msgpackC).Decode(&o)
		return t, o.Content.Content, err
	}
	return
}
