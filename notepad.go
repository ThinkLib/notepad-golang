package main

type Notepad struct {
	Id       string `storm:"id"`
	Contents string
	Password string
}
