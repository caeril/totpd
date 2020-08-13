package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image/png"
	"time"

	"github.com/savsgio/atreugo"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/caeril/totpd/data"
)

func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func addUser(org string, uid string) string {

	user := data.User{Organization: org, Username: uid}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      user.Organization,
		AccountName: user.Username,
	})
	if err != nil {
		panic(err)
	}

	user.Id = hex.EncodeToString(NewSHA256([]byte(key.URL())))
	user.Secret = key.Secret()
	user.URL = key.URL()

	data.PutUser(user)

	return user.Id

}

func InitHandlers() {

}

func _Post_Validation() error {

	return nil
}

func _Get_Index(ctx *atreugo.RequestCtx) error {

	page := get_markup("Index.html")

	page.vars.Set("title", "Home")

	users := []data.User{}

	for _, uid := range data.ListUsers() {
		user := data.GetUser(uid)
		users = append(users, user)
	}

	page.vars.Set("users", users)

	return render_markup(ctx, page)

}

func _Post_User(ctx *atreugo.RequestCtx) error {

	sUsername := string(ctx.Request.PostArgs().Peek("uid"))
	sOrganization := string(ctx.Request.PostArgs().Peek("org"))

	uid := addUser(sOrganization, sUsername)

	ctx.Redirect("/users/"+uid, 302)
	return nil
}

func _Get_User(ctx *atreugo.RequestCtx) error {

	sValid := string(ctx.QueryArgs().Peek("valid"))
	sInvalid := string(ctx.QueryArgs().Peek("invalid"))

	valid := sValid == "true"
	invalid := sInvalid == "true"

	iId := ctx.UserValue("id")
	sId := fmt.Sprintf("%s", iId)
	user := data.GetUser(sId)

	if len(user.Id) == 0 {
		return ctx.HTTPResponse("404 Nuh uh", 404)
	}

	page := get_markup("User.html")

	page.vars.Set("title", "User")

	page.vars.Set("user", user)

	futureCodes := []string{}

	t := time.Now()

	for i := 0; i < 6; i++ {
		code, _ := totp.GenerateCode(user.Secret, t)
		futureCodes = append(futureCodes, code)
		t = t.Add(time.Second * 30)
	}

	page.vars.Set("codes", futureCodes)
	page.vars.Set("invalid", invalid)
	page.vars.Set("valid", valid)

	return render_markup(ctx, page)

}

func _Get_QR(ctx *atreugo.RequestCtx) error {

	iId := ctx.UserValue("id")
	sId := fmt.Sprintf("%s", iId)
	user := data.GetUser(sId)

	if len(user.Id) == 0 {
		return ctx.HTTPResponse("404 Nuh uh", 404)
	}

	key, err := otp.NewKeyFromURL(user.URL)

	if err != nil {
		panic(err)
	}

	w := bytes.Buffer{}
	image, _ := key.Image(512, 512)
	png.Encode(&w, image)
	ctx.SetContentType("image/png")

	return ctx.HTTPResponseBytes(w.Bytes(), 200)

}

func _Post_Validate(ctx *atreugo.RequestCtx) error {

	sUserId := string(ctx.Request.PostArgs().Peek("uid"))
	sCode := string(ctx.Request.PostArgs().Peek("code"))

	user := data.GetUser(sUserId)

	result := totp.Validate(sCode, user.Secret)

	if result {
		ctx.Redirect("/users/"+user.Id+"?valid=true", 302)
	} else {
		ctx.Redirect("/users/"+user.Id+"?invalid=true", 302)
	}

	return nil
}
