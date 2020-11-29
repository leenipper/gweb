package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/life4/gweb/web"
)

type Game struct {
	Width  int
	Height int
	Window web.Window
	Canvas web.Canvas
	Body   web.HTMLElement

	state    *State
	platform Platform
	ball     Ball
	block    TextBlock
	bricks   Bricks
	credits  int // number of lost ball credits
}

func (game *Game) Init() {
	game.state = &State{Stop: SubState{}}
	game.credits = InitialCredits
	context := game.Canvas.Context2D()

	// draw background
	context.SetFillStyle(BGColor)
	context.BeginPath()
	context.Rectangle(0, 0, game.Width, game.Height).Filled().Draw()
	context.Fill()
	context.ClosePath()

	// make handlers
	rect := Rectangle{
		x:      game.Width / 2,
		y:      game.Height - 60,
		width:  PlatformWidth,
		height: PlatformHeight,
	}
	platformCicrle := CircleFromRectangle(rect)
	game.platform = Platform{
		rect:         &rect,
		circle:       &platformCicrle,
		context:      context,
		element:      game.Canvas,
		mouseX:       game.Width / 2,
		windowWidth:  game.Width,
		windowHeight: game.Height,
	}
	game.block = TextBlock{context: context, updated: time.Now()}
	game.ball = Ball{
		context:      context,
		windowWidth:  game.Width,
		windowHeight: game.Height,
		platform:     &game.platform,
	}
	game.ball.Reset()
	game.bricks = Bricks{
		context:      context,
		windowWidth:  game.Width,
		windowHeight: game.Height,
		ready:        false,
		text:         &game.block,
	}
	go game.bricks.Draw()
	game.block.DrawCredits(game.credits)
}

func (game *Game) handler() {
	if game.state.Stop.Requested {
		game.state.Stop.Complete()
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(5)
	go func() {
		// update FPS
		game.block.handle(game.credits)
		wg.Done()
	}()
	go func() {
		// update platform position
		game.platform.handleFrame()
		wg.Done()
	}()
	go func() {
		// check if the ball should bounce from a brick
		game.bricks.Handle(&game.ball)
		wg.Done()
	}()
	go func() {
		// check if the ball should bounce from border or platform
		game.ball.handle()
		wg.Done()
	}()
	go func() {
		// check if ball got out of playground
		if game.ball.y >= (game.Height + BallSize) {
			// When out of credits, game is over.
			if game.credits == 0 {
				go game.fail()
			} else {
				game.credits -= 1
				game.ball.Reset()
			}
		}
		//
		if game.bricks.Count() == 0 {
			go game.win()
		}
		wg.Done()
	}()
	wg.Wait()

	game.Window.RequestAnimationFrame(game.handler, false)
}

func (game *Game) Register() {
	game.state = &State{Stop: SubState{}}
	// register mouse movement handler
	game.Body.EventTarget().Listen(web.EventTypeMouseMove, game.platform.handleMouse)
	// register frame updaters
	game.Window.RequestAnimationFrame(game.handler, false)
}

func (game *Game) Stop() {
	if game.state.Stop.Completed {
		return
	}
	game.state.Stop.Request()
	game.state.Stop.Wait()
}

func (game *Game) fail() {
	game.Stop()
	game.drawText("Game Over", FailColor)
}

func (game *Game) win() {
	game.Stop()
	game.drawText("You Win", WinColor)
}

func (game *Game) drawText(text, color string) {
	height := TextHeight * 2
	width := TextWidth * 2
	context := game.Canvas.Context2D()
	context.Text().SetFont(fmt.Sprintf("bold %dpx Roboto", height))
	context.SetFillStyle(color)
	context.Text().Fill(text, (game.Width-width)/2, (game.Height-height)/2, width)
}
