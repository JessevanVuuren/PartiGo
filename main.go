package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	HEIGHT     = 800
	WIDTH      = 800
	WORLD_SIZE = 300
	SPEED      = 1
	STEPS      = 8
)

type Draw struct {
	source rl.Rectangle
	origin rl.Vector2
	color  rl.Color
	size   rl.Vector2
}

type Particle struct {
	old    rl.Vector2
	pos    rl.Vector2
	acc    rl.Vector2
	radius float32
	draw   Draw
}

func (p *Particle) setDraw(t rl.Texture2D) {
	scale := (p.radius * 2) / float32(t.Width)

	tW := float32(t.Width)
	tH := float32(t.Height)

	p.draw.size = rl.Vector2{X: tW * scale, Y: tH * scale}

	p.draw.source = rl.Rectangle{X: 0.0, Y: 0.0, Width: tW, Height: tH}
	p.draw.origin = rl.Vector2{X: tW / 2 * scale, Y: tH / 2 * scale}
	p.draw.color = rl.GetColor(0xFFFFFFFF)
}

func (p *Particle) getPos() rl.Rectangle {
	return rl.Rectangle{
		X:      p.pos.X,
		Y:      p.old.Y,
		Width:  p.draw.size.X,
		Height: p.draw.size.Y,
	}
}

func render_particle(p []Particle, texture rl.Texture2D) {
	for i := 0; i < len(p); i++ {

		speed := rl.Vector2Length(rl.Vector2Subtract(p[i].old, p[i].pos))
		color := rl.Color{uint8(255 * (speed / 2)), uint8(255 * (speed / 2)), uint8(255 * (speed / 2)), 255}
		rl.DrawTexturePro(texture, p[i].draw.source, p[i].getPos(), p[i].draw.origin, 0, color)

	}
}

func apply_force(particles []Particle) {
	for i := 0; i < len(particles); i++ {
		particles[i].acc = rl.Vector2Add(particles[i].acc, rl.Vector2{X: 0, Y: 10})
	}
}

func word_box(particles []Particle) {
	for i := 0; i < len(particles); i++ {
		dist := rl.Vector2Subtract(rl.Vector2{X: WIDTH / 2, Y: HEIGHT / 2}, particles[i].pos)
		length := rl.Vector2Length(dist)

		if length > WORLD_SIZE-particles[i].radius {
			nX := dist.X / length
			nY := dist.Y / length

			particles[i].pos.X = (WIDTH / 2) - nX*(WORLD_SIZE-particles[i].radius)
			particles[i].pos.Y = (WIDTH / 2) - nY*(WORLD_SIZE-particles[i].radius)
		}
	}
}

func collision(particles []Particle) {
	for i := 0; i < len(particles); i++ {
		for j := 0; j < len(particles); j++ {
			if i != j {
				vec := rl.Vector2Subtract(particles[i].pos, particles[j].pos)
				length := rl.Vector2Length(vec)
				if length < particles[i].radius+particles[j].radius {
					norm := rl.Vector2Scale(vec, 1/length)
					delta := particles[i].radius + particles[j].radius - length

					adjustment := rl.Vector2Scale(norm, delta/2)
					particles[i].pos = rl.Vector2Add(particles[i].pos, adjustment)
					particles[j].pos = rl.Vector2Subtract(particles[j].pos, adjustment)
				}
			}
		}
	}
}

func attraction_force(particles []Particle) {
	mousePos := rl.GetMousePosition()

	for i := 0; i < len(particles); i++ {
		delta := rl.Vector2Subtract(mousePos, particles[i].pos)
		length := rl.Vector2Length(delta)
		if length < 100 {
			norm := rl.Vector2Scale(delta, 1/length)
			particles[i].acc = rl.Vector2Add(particles[i].acc, rl.Vector2Scale(norm, 240))
		}
	}
}

func integration(particles []Particle, dt float32) {
	for i := 0; i < len(particles); i++ {
		var displacement = rl.Vector2Subtract(particles[i].pos, particles[i].old)

		particles[i].old = particles[i].pos

		particles[i].pos.X = particles[i].pos.X + displacement.X + particles[i].acc.X*(dt*dt)
		particles[i].pos.Y = particles[i].pos.Y + displacement.Y + particles[i].acc.Y*(dt*dt)

		particles[i].acc = rl.Vector2Zero()
	}
}

func spawn_particle(tex rl.Texture2D, x, y float32) Particle {
	p := Particle{
		old:    rl.Vector2{X: x, Y: y},
		pos:    rl.Vector2{X: x, Y: y},
		acc:    rl.Vector2Zero(),
		radius: 10,
	}

	p.setDraw(tex)
	return p
}

func spawn_grid(particles []Particle, tex rl.Texture2D) []Particle {
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			X := WIDTH/2 - float32(25*5)
			Y := HEIGHT/2 - float32(25*5)
			particles = append(particles, spawn_particle(tex, X+float32(x*25), Y+float32(y*25)))

		}
	}
	return particles
}

// left mouse = spawn a lot
// right mouse = attraction force
// s = spawn 100 in grid

func main() {
	rl.InitWindow(WIDTH, HEIGHT, "PartiGo")
	rl.SetTargetFPS(60)

	var tex = rl.LoadTextureFromImage(rl.LoadImage("./ball.png"))

	// shader := rl.LoadShader("shaders/particle.vs", "shaders/particle.fs")

	particles := []Particle{}
	particles = spawn_grid(particles, tex)

	for !rl.WindowShouldClose() {

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			particles = append(particles, spawn_particle(tex, 200, 200))
		}

		if rl.IsKeyPressed(rl.KeyS) {
			particles = spawn_grid(particles, tex)
		}

		for i := 0; i < STEPS; i++ {
			if rl.IsMouseButtonDown(rl.MouseButtonRight) {
				attraction_force(particles)
			}
			apply_force(particles)
			integration(particles, rl.GetFrameTime()*SPEED)
			collision(particles)
			word_box(particles)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.GetColor(0x181818FF))
		rl.DrawCircle(WIDTH/2, HEIGHT/2, WORLD_SIZE, rl.GetColor(0x101010FF))
		// rl.BeginShaderMode(shader)
		render_particle(particles, tex)
		// rl.EndShaderMode()
		rl.DrawText("PartiGos: "+strconv.Itoa(len(particles)), 5, 5, 20, rl.GetColor(0xFFFFFFFF))
		rl.DrawFPS(5, 30)
		rl.EndDrawing()
	}
}
