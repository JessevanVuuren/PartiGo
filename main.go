package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	HEIGHT     = 800
	WIDTH      = 800
	WORLD_SIZE = 300
	SPEED      = 1
	STEPS      = 8
)

type Particle struct {
	old    rl.Vector2
	pos    rl.Vector2
	acc    rl.Vector2
	radius float32
	color  rl.Color
}

func render_particle(particles []Particle) {
	for i := 0; i < len(particles); i++ {
		rl.DrawCircle(int32(particles[i].pos.X), int32(particles[i].pos.Y), particles[i].radius, particles[i].color)
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

func vec3_divide_f(v rl.Vector2, l float32) rl.Vector2 {
	return rl.Vector2{X: v.X / l, Y: v.Y / 2}
}

func collision(particles []Particle) {
	for i := 0; i < len(particles); i++ {
		for j := 0; j < len(particles); j++ {
			if i != j {
				vec := rl.Vector2Subtract(particles[i].pos, particles[j].pos)
				length := rl.Vector2Length(vec)
				if length < particles[i].radius+particles[j].radius {
					norm := vec3_divide_f(vec, length)
					delta := particles[i].radius + particles[j].radius - length
					norm = rl.Vector2Scale(norm, delta)
					// particles[i].pos = rl.Vector2Add(particles[i].pos, norm)
					// particles[j].pos = rl.Vector2Subtract(particles[j].pos, norm)

				}
			}
		}

		// mfloat_t axis[VEC3_SIZE];
		// vec3_subtract(axis, a->current, b->current);
		// mfloat_t dist = vec3_length(axis);
		// if (dist < a->radius + b->radius) {
		// 	mfloat_t norm[VEC3_SIZE];
		// 	vec3_divide_f(norm, axis, dist);
		// 	mfloat_t delta = a->radius + b->radius - dist;
		// 	vec3_multiply_f(norm, norm, 0.5 * delta);
		// 	vec3_add(a->current, a->current, norm);
		// 	vec3_subtract(b->current, b->current, norm);
		// }
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

func main() {
	rl.InitWindow(WIDTH, HEIGHT, "PartiGo")
	rl.SetTargetFPS(60)

	particles := []Particle{}

	for !rl.WindowShouldClose() {

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			particles = append(particles, Particle{
				old:    rl.Vector2{X: 200, Y: 200},
				pos:    rl.Vector2{X: 200, Y: 200},
				acc:    rl.Vector2Zero(),
				radius: 10,
				color:  rl.GetColor(0xffffffff),
			})
		}

		for i := 0; i < STEPS; i++ {
			apply_force(particles)
			integration(particles, rl.GetFrameTime()*SPEED)
			collision(particles)
			word_box(particles)
		}

		rl.BeginDrawing()

		rl.DrawCircle(WIDTH/2, HEIGHT/2, WORLD_SIZE, rl.GetColor(0x101010FF))

		render_particle(particles)
		rl.ClearBackground(rl.GetColor(0x181818FF))

		rl.EndDrawing()
	}
}
