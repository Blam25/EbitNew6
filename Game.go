package main

import (
	"sync"

	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Update() error {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	for _, s := range Systems {
		wg.Add(1)
		go s(wg)
	}
	wg.Wait()
	//go TestSys()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	var wg *sync.WaitGroup = &sync.WaitGroup{}

	for _, s := range DrawSystems {
		wg.Add(1)
		go s(screen, wg)
	}
	wg.Wait()

	//go TestSys2(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

var Systems []func(wg *sync.WaitGroup)
var DrawSystems []func(screen *ebiten.Image, wg *sync.WaitGroup)

func init() {
	//Systems = append(Systems, TestSys)
	//Systems = append(Systems, TestSys)

	Systems = append(Systems, TestSysMovement)
	Systems = append(Systems, setRect)

	DrawSystems = append(DrawSystems, TestSys2)
}

func TestSys(wg *sync.WaitGroup) {
	defer wg.Done()

	hmm := 1
	what := false

	Comps.Position.IterateWrite(func(i int, s *Position) {
		s.X = s.X + hmm
		print(s.X)
		if s.X > 500 {
			what = true
		}
	})

	if what {
		Comps.Position.IterateWrite(func(i int, s *Position) {
			print("hejsan")
		})
	}

	Comps.Position.IterateWrite(func(i int, s *Position) {
		if Comps.Image.GetRead(i) != nil {

		}
	})
}

func TestSys2(screen *ebiten.Image, wg *sync.WaitGroup) {
	defer wg.Done()

	Comps.Position.IterateRead(func(i int, s *Position) {
		if a := Comps.Image.GetRead(i); a != nil {
			a.op.GeoM.Reset()
			a.op.GeoM.Translate(float64(s.X), float64(s.Y))
			screen.DrawImage(a.image, &a.op)
		}
	})
}

func TestSysMovement(wg *sync.WaitGroup) {
	defer wg.Done()

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		Comps.Position.IterateWrite(func(s int, pos *Position) {
			if Comps.Player.GetRead(pos.Getid()) == nil {
				pos.Y = pos.Y + 5
			}
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		Comps.Position.IterateWrite(func(s int, pos *Position) {
			pos.Y = pos.Y - 5
		})
		Comps.Player.IterateRead(func(s int, player *Player) {
			Comps.Position.GetWrite(player.Getid(), func(pos *Position) {
				pos.Y = pos.Y + 5
			})
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		Comps.Position.IterateWrite(func(s int, pos *Position) {
			pos.X = pos.X - 5
		})
		Comps.Player.IterateRead(func(s int, player *Player) {
			Comps.Position.GetWrite(player.Getid(), func(pos *Position) {
				pos.X = pos.X + 5
			})
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		Comps.Position.IterateWrite(func(s int, pos *Position) {
			pos.X = pos.X + 5
		})
		Comps.Player.IterateRead(func(s int, player *Player) {
			Comps.Position.GetWrite(player.Getid(), func(pos *Position) {
				pos.X = pos.X - 5
			})
		})
	}
}

func setRect(wg *sync.WaitGroup) {
	defer wg.Done()

	Comps.Rect.IterateWrite(func(i int, s *Rect) {
		x := int(Comps.Position.GetRead(s.Getid()).X)
		y := int(Comps.Position.GetRead(s.Getid()).Y)
		s.Rect = image.Rect(
			x,
			y,
			x+s.Width,
			y+s.Height)
		s.Top = image.Rect(
			x+2,
			y-s.Height,
			x+s.Width-2,
			y-s.Height+10)
		s.Bottom = image.Rect(
			x+2,
			y,
			x+s.Width-2,
			y-10)
		s.Right = image.Rect(
			x+s.Width-10,
			y+2,
			x+s.Width,
			y+s.Height-2)
		s.Left = image.Rect(
			x,
			y+2,
			x+10,
			y+s.Height-2)
	})

}
