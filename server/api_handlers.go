package server

import (
	"encoding/json"
	"github.com/savsgio/atreugo"
)

type postUserModel struct {
	Organization string `json:"organization"`
	Username     string `json:"username"`
}

func APIPostUser(ctx *atreugo.RequestCtx) error {

	posted := postUserModel{}
	he(json.Unmarshal(ctx.PostBody(), &posted))

	return nil

}
