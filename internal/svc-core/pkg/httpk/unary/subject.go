package unaryHttpk

import (
	f "github.com/valyala/fasthttp"
)

const SubjectKey = "subject"

type Subject struct {
	Id       string `json:"id"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
}

var AnonymousSubject = Subject{
	Id:       "ANON",
	FullName: "Anonymous User",
	Role:     "USER",
}

func SetSubject(ctx *f.RequestCtx, id, fullName, role string) {
	ctx.SetUserValue(SubjectKey, Subject{
		Id:       id,
		FullName: fullName,
		Role:     role,
	})
}

func GetSubject(ctx *f.RequestCtx) Subject {
	val, ok := ctx.UserValue(SubjectKey).(Subject)
	if !ok {
		return AnonymousSubject
	}
	return val
}

func ExtractSubject(next f.RequestHandler) f.RequestHandler {
	return func(ctx *f.RequestCtx) {
		subject := parseSubjectFromHeader(ctx)
		ctx.SetUserValue(SubjectKey, subject)
		next(ctx)
	}
}

func parseSubjectFromHeader(ctx *f.RequestCtx) Subject {
	id := string(ctx.Request.Header.Peek("X-User-Id"))
	if id == "" {
		return AnonymousSubject
	}

	return Subject{
		Id:       id,
		FullName: string(ctx.Request.Header.Peek("X-User-Name")),
		Role:     string(ctx.Request.Header.Peek("X-User-Role")),
	}
}
