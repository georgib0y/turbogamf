package main

import (
	"fmt"
)

func main() {

}

type Player interface{}

type Game interface {
	Players() Player
}
