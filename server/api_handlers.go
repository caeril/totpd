package server

import (
	"encoding/json"
	"fmt"
	"github.com/caeril/totpd/data"
	"github.com/pquerna/otp/totp"
	"github.com/savsgio/atreugo"
)

type postUserModel struct {
	Organization string `json:"organization"`
	Username     string `json:"username"`
}

func APIPostUser(ctx *atreugo.RequestCtx) error {

	posted := postUserModel{}
	he(json.Unmarshal(ctx.PostBody(), &posted))

	uid := addUser(posted.Organization, posted.Username)

	response := struct {
		Message string `json:"message"`
		UID     string `json:"uid"`
	}{"success", uid}

	responseJSON, _ := json.Marshal(response)

	ctx.Response.Header.Add("Content-Type", "application/json")
	return ctx.HTTPResponseBytes(responseJSON, 200)

}

func APIValidateCode(ctx *atreugo.RequestCtx) error {

	iUID := ctx.UserValue("uid")
	sUID := iUID.(string)
	sCode := string(ctx.QueryArgs().Peek("code"))

	user := data.GetUser(sUID)

	if len(user.Id) == 0 {
		return ctx.HTTPResponse("404 Nuh uh", 404)
	}

	result := totp.Validate(sCode, user.Secret)

	response := struct {
		Message string `json:"message"`
	}{}

	if result {
		response.Message = "valid"
	} else {
		response.Message = "invalid"
	}

	responseJSON, _ := json.Marshal(response)

	ctx.Response.Header.Add("Content-Type", "application/json")
	return ctx.HTTPResponseBytes(responseJSON, 200)

}
