package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 1920, 1080

type color struct {
	r, g, b byte
}

type position struct {
	// float32 is big enough to use in most general task,
	// and small enough for performance
	x, y float32
}

// Making the ball
type ball struct {
	color
	position
	r  int
	xv float32
	yv float32
}

func (ball *ball) draw(pixels []byte) {

}

// Making the paddles
type paddle struct {
	color
	position
	h int
	w int
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, color{255, 255, 255}, pixels)
		}
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
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

	for y := 0; y < winHeight; y++ {
		for x := 0; x < winWidth; x++ {
			setPixel(x, y, color{byte(x % 255), byte(y % 255), 0}, pixels)
		}
	}

	tex.Update(nil, pixels, winWidth*4)
	renderer.Copy(tex, nil, nil)
	renderer.Present()

	// Changd after EP 06 to address MacOSX
	// OSX requires that you consume events for windows to open and work properly
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
	}
}
