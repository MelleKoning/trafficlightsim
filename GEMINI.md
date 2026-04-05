# GEMINI.md - Traffic Light Sim

## Project Overview
`trafficlightsim` is a Go-based 3D simulation of a traffic light. It renders a detailed 3D model of a traffic light (including a backboard, white border, and lamp hoods) and uses OpenGL for real-time rendering.

### Main Technologies
- **Language:** Go 1.26+
- **Graphics API:** OpenGL 4.1 Core Profile
- **Windowing/Input:** [GLFW 3.2](https://github.com/go-gl/glfw)
- **Math Library:** [MathGL (mgl32)](https://github.com/go-gl/mathgl) for 3D transformations.

## Project Structure
- `cmd/trafficlightsim/main.go`: The primary entry point. It contains:
  - Shader source code (GLSL 410).
  - Procedural geometry generation for the traffic light.
  - The main rendering loop and input handling.
- `vendor/`: Local copies of dependencies (standard for some Go projects).
- `.golangci.yml`: Configuration for `golangci-lint`.
- `.pre-commit-config.yaml`: Hooks for code quality and consistency.

## Building and Running

### Prerequisites
Ensure you have the necessary system libraries for OpenGL and GLFW development (e.g., on Linux: `libgl1-mesa-dev`, `libglfw3-dev`, `libx11-dev`, `libxcursor-dev`, `libxinerama-dev`, `libxrandr-dev`, `libxi-dev`).

### Commands
- **Run the simulation:**
  ```bash
  go run cmd/trafficlightsim/main.go
  ```
- **Build the executable:**
  ```bash
  go build -o trafficlightsim ./cmd/trafficlightsim
  ```
- **Linting:**
  ```bash
  golangci-lint run
  ```

## Development Conventions

### Controls
While the simulation is running, use the **Arrow Keys** to interact with the 3D model:
- **Up/Down:** Rotate around the X-axis.
- **Left/Right:** Rotate around the Y-axis.

### Coding Standards
- **Style:** Adhere to standard Go formatting (`go fmt`).
- **Quality:** Before committing, ensure that `golangci-lint` passes and all `pre-commit` hooks are satisfied.
- **Shaders:** Shaders are currently embedded as string constants in `main.go`. Any modifications should maintain OpenGL 4.1 compatibility.
- **Geometry:** The traffic light geometry is procedurally generated in `generateTrafficLightData()`. Modifications to the model's structure should happen there.
