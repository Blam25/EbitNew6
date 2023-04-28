package EbitNew6

import (
	//E "Nice2"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var identifier int

type Position struct {
	X int
	Y int
}

type Entity struct {
	Id int
}

func NewComp[T any]() *Component[T] {
	new := &Component[T]{theMap: make(map[int]*T)}
	return new

}

type Component[T any] struct {
	theMap map[int]*T
}

/*
func Add[T any](object T, objectArray map[int]T) {
	//objectArray = append(objectArray, object)
	objectArray[s.Id] = object
}*/

func NewEntity() *Entity {
	new := &Entity{}
	new.Id = identifier
	identifier++
	return new
}

var PosMap map[int]*Position
var PosArr []*Position

func init() {
	PosMap = make(map[int]*Position)

	//Ent1 := NewEntity()
	//Ent1.Add(&Position{}, PosMap)

}

func main() {
	ebiten.SetWindowSize(640, 480)
	//ebiten.SetVsyncEnabled(false)
	//ebiten.SetTPS(ebiten.SyncWithFPS)
	//ebiten.SetFPSMode(ebiten.)
	ebiten.SetWindowTitle("Render an image")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
