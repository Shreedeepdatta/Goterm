package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

const (
	wall        = 1
	nothing     = 0
	playerconst = 10
	max_samples = 100
)

var flag bool

type position struct {
	x, y int
}
type player struct {
	pos   position
	level *level
}

func (p *player) update() {
	if flag == true {
		p.pos.x += 1
		flag = false
	} else {
		p.pos.x -= 1
		flag = true
	}

}

type stats struct {
	start  time.Time
	frames int
	fps    float64
}

func newstats() *stats {
	return &stats{
		start: time.Now(),
	}
}
func (s *stats) update() {
	s.frames++
	if s.frames == max_samples {
		s.fps = float64(s.frames) / float64(time.Since(s.start).Milliseconds())
		s.frames = 0
		s.start = time.Now()
	}
}

type level struct {
	height int
	width  int
	data   [][]int
}

func (l *level) set(pos position, v int) {
	l.data[pos.y][pos.x] = v
}

type game struct {
	isRunning bool
	level     *level
	stats     *stats
	player    *player
	drawbuf   *bytes.Buffer
}

func newLevel(width, height int) *level {

	data := make([][]int, height)
	for h := 0; h < height; h++ {
		for w := 0; w < height; w++ {
			data[h] = make([]int, width)
		}
	}
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if w == width-1 {
				data[h][w] = wall
			}
			if h == height-1 {
				data[h][w] = wall
			}
			if w == 0 {
				data[h][w] = wall
			}
			if h == 0 {
				data[h][w] = wall
			}
		}
	}
	return &level{
		width:  width,
		height: height,
		data:   data,
	}
}

func newGame(width, height int) *game {
	level := newLevel(width, height)
	return &game{
		level:   level,
		drawbuf: new(bytes.Buffer),
		stats:   newstats(),
		player: &player{
			level: level,
			pos:   position{x: 2, y: 5},
		},
	}
}
func (g *game) start() {
	g.isRunning = true
	g.loop()
}
func (g *game) loop() {
	for g.isRunning {
		g.update()
		g.render()
		g.stats.update()

		time.Sleep(time.Second * 1)
	}
}
func (g *game) update() {
	g.level.set(g.player.pos, nothing)
	g.player.update()
	g.level.set(g.player.pos, playerconst)
}
func (g *game) renderPlayer() {
	//g.level.data[g.player.pos.y][g.player.pos.x] = playerconst

}
func (g *game) renderLevel() {
	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			if g.level.data[h][w] == nothing {
				g.drawbuf.WriteString(" ")
			}
			if g.level.data[h][w] == wall {
				g.drawbuf.WriteString("|")
			}
			if g.level.data[h][w] == playerconst {
				g.drawbuf.WriteString("0")
			}
		}
		g.drawbuf.WriteString("\n")
	}
	//g.drawbuf.WriteString("happy")

}
func (g *game) render() {
	g.drawbuf.Reset()
	fmt.Fprint(os.Stdout, "\033[H\033[2J")
	g.renderLevel()
	g.renderstats()
	fmt.Fprint(os.Stdout, g.drawbuf.String())

}
func (g *game) renderstats() {
	g.drawbuf.WriteString("--Stats\n")
	g.drawbuf.WriteString(fmt.Sprintf("FPS: %.5f", g.stats.fps))
}
func main() {

	width := 80
	height := 20
	//fmt.Fprint(os.Stdout, "\033[2J")
	g := newGame(width, height)
	g.start()
}
