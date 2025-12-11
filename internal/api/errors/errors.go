package errors

import "errors"

var (
    ErrInvalidUsername = errors.New("invalid username format").Error()
    ErrInvalidEmail    = errors.New("invalid email format").Error()
    ErrInvalidPassword = errors.New("invalid password").Error()
    ErrDuplicateUser   = errors.New("username or email already exists").Error()
    ErrFileNotFound    = errors.New("cannot read file from request").Error()
	ErrBadRequest = errors.New("bad Request").Error()
	ErrConnection = errors.New("connection Error").Error()
	Errfilesize = errors.New("give the minimum file size").Error()
    ErrInserterr = errors.New("cannot insert the value in the database").Error()
    ErrDelete = errors.New("cannot delete the value from the database").Error()
    Errredis = errors.New("cannot fetch the data").Error()
    Errcookie = errors.New("cannot get the cookie").Error()
    Errfetch = errors.New("cannot fetch from the postgres").Error()
    Errredisfetcherr = errors.New("cannot fetch the value from the redis").Error()
    ErrImage = errors.New("cannot get the image ").Error()
    Errminio = errors.New("cannot put the object into bucket")
    ErrPresignedUrl = errors.New("cannot generate the presigned Url").Error()
)
