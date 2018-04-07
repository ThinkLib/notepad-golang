package main

import (
	"github.com/valyala/fasthttp"
	"crypto/md5"
	"fmt"
	"strings"
	"html/template"
	"log"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	if cmd.debug {
		ctx.Request.Header.VisitAll(func(key, value []byte) {
			log.Printf("%s: %s", key, value)
		})
	}
	log.Println(string(getIp(ctx)) + " " + string(ctx.Method()) + " " + string(ctx.UserAgent()) + " " + ctx.URI().String())
	switch {
	case string(ctx.Path()) == "/":
		has := md5.Sum(getIp(ctx))
		ctx.SetStatusCode(302)
		ctx.Response.Header.Set("Location", "/"+fmt.Sprintf("%x", has))

	case strings.HasPrefix(string(ctx.Path()), "/static/"):
		fasthttp.FSHandler(cmd.staticPath+"static", 1)(ctx)

	case strings.HasPrefix(string(ctx.Path()), "/ajax/update_contents/"):
		contents := ctx.FormValue("contents")
		enc_str := ctx.FormValue("enc_str")
		uris := strings.Split(string(ctx.RequestURI()), "/")
		if len(uris) < 4 {
			ctx.Error("not found", fasthttp.StatusNotFound)
			return
		}
		if !checkUri(uris[3]) {
			ctx.Error("uri error", fasthttp.StatusForbidden)
			return
		}
		var notepad Notepad
		notepad.Id = uris[3]
		if err := db.One("Id", notepad.Id, &notepad); err != nil {
			ctx.Error("not init", fasthttp.StatusInternalServerError)
			return
		}
		notepad.Contents = string(contents)

		if len(notepad.Password) > 0 {
			if notepad.Password == string(ctx.Request.Header.Cookie("password_"+notepad.Id)) {
				if len(enc_str) > 0 {
					notepad.Password = string(ctx.FormValue("enc_str"))
				} else {
					notepad.Password = ""
				}
				if err := db.Update(&notepad); err != nil {
					ctx.Error("error", fasthttp.StatusInternalServerError)
					return
				}
				ctx.WriteString("update success")
			} else {
				ctx.Error("password error", fasthttp.StatusForbidden)
			}
		} else {
			if len(enc_str) > 0 {
				notepad.Password = string(ctx.FormValue("enc_str"))
			} else {
				notepad.Password = ""
			}
			if err := db.Update(&notepad); err != nil {
				ctx.Error("error", fasthttp.StatusInternalServerError)
				return
			}
			ctx.WriteString("update success")
		}

	case strings.HasPrefix(string(ctx.Path()), "/ajax/get_contents/"):
		uris := strings.Split(string(ctx.RequestURI()), "/")
		if len(uris) < 4 {
			ctx.Error("not found", fasthttp.StatusNotFound)
			return
		}
		if !checkUri(uris[3]) {
			ctx.Error("uri error", fasthttp.StatusForbidden)
			return
		}
		var notepad Notepad
		notepad.Id = uris[3]
		if err := db.One("Id", notepad.Id, &notepad); err != nil {
			ctx.Error("not found", fasthttp.StatusForbidden)
			return
		}
		if len(notepad.Password) > 0 {
			if notepad.Password == string(ctx.Request.Header.Cookie("password_"+notepad.Id)) {
				ctx.WriteString(notepad.Contents)
			} else {
				ctx.Error("password error", fasthttp.StatusForbidden)
			}
		} else {
			ctx.WriteString(notepad.Contents)
		}

	default:
		uris := strings.Split(string(ctx.RequestURI()), "/")
		if len(uris) < 2 {
			ctx.Error("not found", fasthttp.StatusNotFound)
			return
		}
		if !checkUri(uris[1]) {
			ctx.Error("uri error", fasthttp.StatusForbidden)
			return
		}
		var notepad Notepad
		notepad.Id = uris[1]
		if err := db.One("Id", notepad.Id, &notepad); err != nil {
			notepad.Contents = ""
			notepad.Password = ""
			db.Save(&notepad)
		}
		if len(notepad.Password) == 0 {
			ctx.Response.Header.SetContentType("text/html; charset=utf-8")
			t, _ := template.ParseFiles("static/index.html")
			t.Execute(ctx, notepad)
		} else {
			password := ctx.Request.Header.Cookie("password_" + notepad.Id)
			if string(password) == notepad.Password {
				ctx.Response.Header.SetContentType("text/html; charset=utf-8")
				t, _ := template.ParseFiles("static/index.html")
				t.Execute(ctx, notepad)
			} else {
				ctx.Response.Header.SetContentType("text/html; charset=utf-8")
				t, _ := template.ParseFiles("static/password.html")
				t.Execute(ctx, notepad)
			}
		}
	}
}

func checkUri(uri string) bool {
	return true
}

func getIp(ctx *fasthttp.RequestCtx) []byte {
	var ip []byte
	if cmd.ipHttpHeader == "NOT" {
		ip = ctx.RemoteIP()
	} else {
		ip = ctx.Request.Header.Peek(cmd.ipHttpHeader)
	}
	return ip
}
