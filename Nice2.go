package main

import (
	"sync"
	//E "Nice2"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var identifier int

type Position struct {
	Entity
	*mu
	X int
	Y int
}

type mu struct {
	*sync.Mutex
}

func (s *mu) LockMu() {
	s.Lock()
}

func (s *mu) UnlockMu() {
	s.Unlock()
}

type Image struct {
	Entity
	*mu
	image *ebiten.Image
	op    ebiten.DrawImageOptions
}

type Entity struct {
	id int
}

func (s Entity) Getid() int {
	return s.id
}

type validComp interface {
	Getid() int
	LockMu()
	UnlockMu()
}

func NewComp[T validComp]() *Component[T] {
	new := &Component[T]{index: make(map[int]int)}
	return new
}

type Component[T validComp] struct {
	index    map[int]int
	theArray []*T
}

func (s *Component[T]) Add(object *T) {
	s.index[(*object).Getid()] = len(s.theArray)
	s.theArray = append(s.theArray, object)
}

func (s *Component[T]) Remove(object *T) {
	(*object).LockMu()
	//index of object to be removed
	index := s.index[(*object).Getid()]
	//delete id and index of said object from map
	delete(s.index, (*object).Getid())
	(*s.theArray[len(s.theArray)-1]).LockMu()
	//set value of deleted index to the last object in array, thereby deleting it
	s.theArray[index] = s.theArray[len(s.theArray)-1]
	//get id of moved index
	movedId := (*s.theArray[index]).Getid()
	//set new index of moved object correctly in map
	s.index[movedId] = index
	(*s.theArray[index]).UnlockMu()
	//delete the last (now duplicated) object from the array
	s.theArray = s.theArray[:len(s.theArray)-1]
	(*object).UnlockMu()
}

func (s *Component[T]) IterateWrite(f func(int, *T)) {

	println("locked")
	for i, s := range s.theArray {

		(*s).LockMu()
		f(i, s)
		(*s).UnlockMu()
	}
	println("Unlocked")
}

func (s *Component[T]) IterateRead(f func(int, *T)) {

	for i, s := range s.theArray {
		z := *s
		y := &z
		f(i, y)
	}
}

func (s *Component[T]) GetWrite(id int, f func(*T)) {

	println("locked")
	index := s.index[id]
	comp := s.theArray[index]
	(*comp).LockMu()
	f(comp)
	(*comp).UnlockMu()
	println("Unlocked")
}

func (s *Component[T]) GetRead(id int) *T {

	z := *s.theArray[id]
	y := &z
	return y
}

/*
	func remove(s []int, i int) []int {
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	}
*/

func NewEntity() Entity {
	new := Entity{}
	new.id = identifier
	identifier++
	return new
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

type Components struct {
	Position *Component[Position]
	Image    *Component[Image]
}

var Comps *Components = &Components{}

func init() {

	Comps.Position = NewComp[Position]()
	Comps.Image = NewComp[Image]()

	var err error
	image1, _, err := ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	Ent1 := NewEntity()
	Comps.Position.Add(&Position{Ent1, &mu{&sync.Mutex{}}, 200, 200})
	Comps.Image.Add(&Image{Entity: Ent1, image: image1})

	Ent2 := NewEntity()
	Comps.Position.Add(&Position{Ent2, &mu{&sync.Mutex{}}, 100, 100})
	Comps.Image.Add(&Image{Entity: Ent2, image: image1})

	//Ent1.Add(&Position{}, PosMap)

}
