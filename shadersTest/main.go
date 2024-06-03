package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialize the window
	rl.InitWindow(800, 600, "raylib-go shader example")
	rl.SetTargetFPS(60)

	// Load a texture
	texture := rl.LoadTexture("raylib_logo.png") // Ensure you have a texture file "raylib_logo.png"

	shader := rl.LoadShader("", "grayscale.fs")

	for !rl.WindowShouldClose() {
		// Update logic here

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Begin using the shader
		rl.BeginShaderMode(shader)

		// Draw the texture with the shader applied
		rl.DrawTexture(texture, 200, 150, rl.White)

		// End using the shader
		rl.EndShaderMode()

		rl.DrawText("Grayscale Shader Example", 10, 10, 20, rl.DarkGray)
		rl.EndDrawing()
	}

	// Unload resources
	rl.UnloadShader(shader)
	rl.UnloadTexture(texture)
	rl.CloseWindow()
}
