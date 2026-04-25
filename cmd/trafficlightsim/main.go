package main

import (
	"runtime"

	"github.com/MelleKoning/trafficlightsim/pkg/trafficlight"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 250 // nolint: mnd
	height = 500 // nolint: mnd
)

var (
	rotX float32
	rotY float32
)

func main() {
	window := initGlfw()

	defer glfw.Terminate()

	trafficlight.InitializeTrafficLight()

	window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, _ int,
		action glfw.Action, _ glfw.ModifierKey,
	) {
		if action == glfw.Press {
			switch key { // nolint:exhaustive
			case glfw.Key1:
				trafficlight.ToggleRed()
			case glfw.Key2:
				trafficlight.ToggleOrange()
			case glfw.Key3:
				trafficlight.ToggleGreen()
			case glfw.KeyPageUp:
				trafficlight.SetAmbient(trafficlight.State.Ambient + 0.1) // nolint: mnd
			case glfw.KeyPageDown:
				trafficlight.SetAmbient(trafficlight.State.Ambient - 0.1) // nolint: mnd
			}
		}
	})

	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyUp) == glfw.Press {
			rotX -= 2.0 // nolint: mnd
		}

		if window.GetKey(glfw.KeyDown) == glfw.Press {
			rotX += 2.0 // nolint: mnd
		}

		if window.GetKey(glfw.KeyLeft) == glfw.Press {
			rotY -= 2.0 // nolint: mnd
		}

		if window.GetKey(glfw.KeyRight) == glfw.Press {
			rotY += 2.0 // nolint: mnd
		}

		trafficlight.SetRotations(rotX, rotY)
		trafficlight.Draw(window)
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // nolint: mnd
	glfw.WindowHint(glfw.ContextVersionMinor, 1) // nolint: mnd
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Verkeerslicht", nil, nil) // nolint: mnd
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	runtime.LockOSThread()

	return window
}
