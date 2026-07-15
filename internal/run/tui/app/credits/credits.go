/*
Copyright 2026 Julien Breux

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package credits

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "credits"
)

type Particle struct {
	x, y   float64
	vx, vy float64
	color  tcell.Color
	char   rune
}

type CreditsPage struct {
	*tview.Box
	app          *tview.Application
	onClose      func()
	particles    []*Particle
	contributors []string
	authorColors map[string]tcell.Color
	scrollY      float64
	lastUpdate   time.Time
	stop         chan struct{}
}

var (
	Contributors = []string{
		"@JulienBreux",
		"@rdout2",
	}
)

func New(app *tview.Application, onClose func()) *CreditsPage {
	credits := []string{
		"Run CLI",
		"",
		"Thanks!",
		"",
	}
	credits = append(credits, Contributors...)
	credits = append(credits, "", "Thank you for using!")

	authorColors := make(map[string]tcell.Color)
	for _, a := range Contributors {
		authorColors[a] = tcell.NewRGBColor(int32(rand.Intn(256)), int32(rand.Intn(256)), int32(rand.Intn(256)))
	}

	c := &CreditsPage{
		Box:          tview.NewBox(),
		app:          app,
		onClose:      onClose,
		contributors: credits,
		authorColors: authorColors,
		stop:         make(chan struct{}),
	}
	c.lastUpdate = time.Now()
	return c
}

// StartAnimation starts the animation loop
func (c *CreditsPage) StartAnimation() {
	c.lastUpdate = time.Now()
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.app.QueueUpdateDraw(func() {
					c.update()
				})
			case <-c.stop:
				return
			}
		}
	}()
}

// StopAnimation stops the animation loop
func (c *CreditsPage) StopAnimation() {
	close(c.stop)
}

func (c *CreditsPage) update() {
	now := time.Now()
	dt := now.Sub(c.lastUpdate).Seconds()
	c.lastUpdate = now

	// Update Author Colors (Shimmer effect)
	for _, a := range Contributors {
		if rand.Intn(10) == 0 { // 10% chance to shift color
			c.authorColors[a] = tcell.NewRGBColor(int32(rand.Intn(256)), int32(rand.Intn(256)), int32(rand.Intn(256)))
		}
	}

	// Update Particles
	_, _, w, h := c.GetInnerRect()
	if w <= 0 || h <= 0 {
		return
	}

	// Move particles
	validParticles := c.particles[:0]
	for _, p := range c.particles {
		p.x += p.vx * dt
		p.y += p.vy * dt
		if p.y < float64(h) && p.x >= 0 && p.x < float64(w) {
			validParticles = append(validParticles, p)
		}
	}
	c.particles = validParticles

	// Spawn new particles
	if len(c.particles) < 100 {
		count := 1
		if dt > 0.1 {
			count = 5 // Spawn more if frame drop or test delay
		}
		for i := 0; i < count; i++ {
			c.particles = append(c.particles, &Particle{
				x:     float64(rand.Intn(w)),
				y:     0,
				vx:    (rand.Float64() - 0.5) * 20, // -10 to 10 chars/sec
				vy:    rand.Float64()*20 + 10,      // 10 to 30 chars/sec
				color: tcell.NewRGBColor(int32(rand.Intn(256)), int32(rand.Intn(256)), int32(rand.Intn(256))),
				char:  []rune("★●◼▲")[rand.Intn(4)],
			})
		}
	}

	// Scroll Credits
	c.scrollY -= 3.0 * dt // Speed (lines/sec)
	totalHeight := len(c.contributors) * 2
	if c.scrollY < float64(-totalHeight-h/2) {
		c.scrollY = float64(h) // Reset to bottom
	}
}

func (c *CreditsPage) Draw(screen tcell.Screen) {
	c.DrawForSubclass(screen, c)
	x, y, w, h := c.GetInnerRect()
	if w <= 0 || h <= 0 {
		return
	}

	// Draw Particles
	for _, p := range c.particles {
		screen.SetContent(x+int(p.x), y+int(p.y), p.char, nil, tcell.StyleDefault.Foreground(p.color))
	}

	// Draw Credits
	center := y + h/2
	for i, line := range c.contributors {
		lineY := float64(center) + c.scrollY + float64(i*2)
		if lineY >= float64(y) && lineY < float64(y+h) {
			color := tcell.ColorWhite
			if c.authorColors[line] != 0 {
				color = c.authorColors[line]
			}
			tview.Print(screen, line, x, int(lineY), w, tview.AlignCenter, color)
		}
	}
}

func (c *CreditsPage) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
			c.StopAnimation()
			c.onClose()
		}
	}
}
