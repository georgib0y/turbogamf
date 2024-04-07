package main

import ()

func main() {

}

type Player interface{}

type Game interface {
	Players() Player
}
