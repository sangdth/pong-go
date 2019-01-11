package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 960, 540

type color struct {
	r, g, b byte
}

type pos struct {
	// float32 is big enough to use in most general task,
	// and small enough for performance
	x, y float32
}

// Making the ball
type ball struct {
	pos
	xv float32
	yv float32
	r  int
	color
}

// *ball is a pointer, and an experience dev said that it is better to use it
// I need to research and test more, is it true or not.
func (ball *ball) draw(pixels []byte) {
	for y := -ball.r; y < ball.r; y++ {
		for x := -ball.r; x < ball.r; x++ {
			if x*x+y*y < ball.r*ball.r {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(paddleLeft *paddle, paddleRight *paddle) {
	ball.x += ball.xv
	ball.y += ball.yv

	// handle collision
	if int(ball.y)-ball.r < 0 || int(ball.y)+ball.r > winHeight {
		ball.yv = -ball.yv
	}

	// handle ball goes out of screen
	if ball.x < 0 || int(ball.x) > winWidth {
		ball.x = 300
		ball.y = 300
	}

	if ball.x < paddleLeft.x+float32(paddleLeft.w)/2 ||
		ball.x > paddleRight.x-float32(paddleRight.w)/2 {
		ball.xv = -ball.xv
	}
}

// Making the paddles
type paddle struct {
	pos
	w int
	h int
	color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

// In C they prefix with SDL_, so all func and variables name prepend with SDL
// in Go we import as sdl so it will be sdl.SCANCODE
func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= 5
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += 5
	}
}

// player 2 is the bot, this is A.I. and it never loses
func (paddle *paddle) AIUpdate(ball *ball) {
	paddle.y = ball.y
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func main() {

	// Added after EP06 to address macosx issues
	// this init makes sure all event system was initialized
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Testing SDL2",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth),
		int32(winHeight),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(winWidth),
		int32(winHeight),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{
		pos:   pos{50, 100},
		w:     20,
		h:     100,
		color: color{255, 255, 255},
	}

	player2 := paddle{
		pos:   pos{float32(winWidth) - 50, 100},
		w:     20,
		h:     100,
		color: color{255, 255, 255},
	}

	ball := ball{
		pos:   pos{300, 300},
		xv:    4,
		yv:    2,
		r:     20,
		color: color{255, 255, 255},
	}

	keyState := sdl.GetKeyboardState()

	// Changd after EP 06 to address MacOSX
	// OSX requires that you consume events for windows to open and work properly
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)

		player1.draw(pixels)
		player1.update(keyState)

		player2.draw(pixels)
		player2.AIUpdate(&ball)

		ball.draw(pixels)
		ball.update(&player1, &player2)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}
}
