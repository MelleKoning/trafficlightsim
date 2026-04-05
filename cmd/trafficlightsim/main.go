package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func generateTrafficLightData() []float32 {
	var data []float32

	// HULPFUNCTIE: Voegt een punt toe met positie en normaal
	addVertex := func(x, y, z, nx, ny, nz float32) {
		data = append(data, x, y, z, nx, ny, nz)
	}

	// 1. ACHTERBORD (Een simpele rechthoek/box)
	// We tekenen de voorkant (Z = 0)
	// Twee driehoeken vormen één rechthoek
	// Punten: x, y, z,  nx, ny, nz
	addVertex(-0.4, 1.0, 0.0, 0.0, 0.0, 1.0)  // nolint: mnd
	addVertex(0.4, 1.0, 0.0, 0.0, 0.0, 1.0)   // nolint: mnd
	addVertex(-0.4, -1.0, 0.0, 0.0, 0.0, 1.0) // nolint: mnd

	addVertex(0.4, 1.0, 0.0, 0.0, 0.0, 1.0)   // nolint: mnd
	addVertex(0.4, -1.0, 0.0, 0.0, 0.0, 1.0)  // nolint: mnd
	addVertex(-0.4, -1.0, 0.0, 0.0, 0.0, 1.0) // nolint: mnd

	// WITTE BIES (Een iets grotere rechthoek achter het zwarte bord)
	// We tekenen deze op Z = -0.01 (net achter het zwarte bord)
	addVertex(-0.48, 1.05, -0.01, 0.0, 0.0, 1.0)  // nolint: mnd
	addVertex(0.48, 1.05, -0.01, 0.0, 0.0, 1.0)   // nolint: mnd
	addVertex(-0.48, -1.05, -0.01, 0.0, 0.0, 1.0) // nolint: mnd

	addVertex(0.48, 1.05, -0.01, 0.0, 0.0, 1.0)   // nolint: mnd
	addVertex(0.48, -1.05, -0.01, 0.0, 0.0, 1.0)  // nolint: mnd
	addVertex(-0.48, -1.05, -0.01, 0.0, 0.0, 1.0) // nolint: mnd
	// 2. DE DRIE KAPPEN (Bogen boven de lampen)
	yPositions := []float32{0.6, 0.0, -0.6}
	for _, yOffset := range yPositions {
		const segments = 16 // nolint: mnd

		const radius = 0.2 // nolint: mnd

		const depth = 0.4 // nolint: mnd

		for i := 0; i < segments; i++ {
			angle1 := float64(i) * math.Pi / segments
			angle2 := float64(i+1) * math.Pi / segments

			x1, y1 := float32(math.Cos(angle1))*radius, float32(math.Sin(angle1))*radius
			x2, y2 := float32(math.Cos(angle2))*radius, float32(math.Sin(angle2))*radius

			// Elke segment van de kap is een rechthoek (2 driehoeken)
			// We berekenen de normalen op basis van de hoek voor gladde ronding
			nx1, ny1 := float32(math.Cos(angle1)), float32(math.Sin(angle1))
			nx2, ny2 := float32(math.Cos(angle2)), float32(math.Sin(angle2))

			// Driehoek 1
			addVertex(x1, y1+yOffset, 0.0, nx1, ny1, 0.0)   // nolint: mnd
			addVertex(x2, y2+yOffset, 0.0, nx2, ny2, 0.0)   // nolint: mnd
			addVertex(x1, y1+yOffset, depth, nx1, ny1, 0.0) // nolint: mnd
			// Driehoek 2
			addVertex(x2, y2+yOffset, 0.0, nx2, ny2, 0.0)   // nolint: mnd
			addVertex(x2, y2+yOffset, depth, nx2, ny2, 0.0) // nolint: mnd
			addVertex(x1, y1+yOffset, depth, nx1, ny1, 0.0) // nolint: mnd
		}
	}

	// 3. DE LAMPEN (Rood, Oranje, Groen)
	for _, yOffset := range yPositions {
		const segments = 16 // nolint: mnd

		const radius = 0.18 // Iets kleiner dan de kappen // nolint: mnd

		for i := 0; i < segments; i++ {
			angle1 := float64(i) * 2 * math.Pi / segments
			angle2 := float64(i+1) * 2 * math.Pi / segments

			// Teken een driehoek van het midden van de lamp naar de rand
			// Z staat op 0.01 om net VOOR het bord te liggen
			addVertex(0, yOffset, 0.01, 0, 0, 1)                                                                 // nolint: mnd
			addVertex(float32(math.Cos(angle1))*radius, yOffset+float32(math.Sin(angle1))*radius, 0.01, 0, 0, 1) // nolint: lll,mnd
			addVertex(float32(math.Cos(angle2))*radius, yOffset+float32(math.Sin(angle2))*radius, 0.01, 0, 0, 1) // nolint: lll,mnd
		}
	}

	return data
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// Why the number 4? The reason is we are using float32 slices, and a 32-bit float has 4 bytes,
	// so we are saying the size of the buffer, in bytes, is 4 times the number of points.
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW) // nolint: mnd

	var vao uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0) // nolint: mnd
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil) // nolint: mnd

	// Positie attribuut (Locatie 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0)) // nolint: mnd,staticcheck
	gl.EnableVertexAttribArray(0)                                       // nolint: mnd

	// Normaal attribuut (Locatie 1) - Nieuw!
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4)) // nolint: mnd,staticcheck
	gl.EnableVertexAttribArray(1)                                         // nolint: mnd

	return vao
}

const (
	width  = 500 // nolint: mnd
	height = 500 // nolint: mnd

	vertexShaderSource = `
#version 410
layout (location = 0) in vec3 vp;
layout (location = 1) in vec3 normal;

out vec3 Normal;
out vec3 FragPos;

uniform mat4 model;

void main() {
    FragPos = vec3(model * vec4(vp, 1.0));
    Normal = mat3(transpose(inverse(model))) * normal;
    gl_Position = model * vec4(vp, 1.0);
}` + "\x00"

	// the vec4(1, 1, 1, 1) is the colour, red, green, blue, alpha
	fragmentShaderSource = `
#version 410
in vec3 Normal;
in vec3 FragPos;

out vec4 frag_colour;

uniform vec3 lightPos;
uniform vec3 lightColor;
uniform vec3 objectColor;
uniform vec3 viewPos;

uniform float emission; // 0.0 voor plastic, 1.0 voor brandende lamp

void main() {
    // 1. Ambient
    float ambientStrength = 0.9; // nolint: mnd
    vec3 ambient = ambientStrength * lightColor;

    // 2. Diffuse (De boosdoener: we maken hier de vec3 'diffuse')
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0); // nolint: mnd
    vec3 diffuse = diff * lightColor;

    // 3. Specular (De glans)
    float specularStrength = 0.5; // nolint: mnd
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32); // nolint: mnd
    vec3 specular = specularStrength * spec * lightColor;

    // Resultaat
    vec3 result = (ambient + diffuse + specular) * objectColor;

    // Voeg extra licht toe als het een lamp is
    if (emission > 0.5) { // nolint: mnd
        result = objectColor * 1.5; // De lamp "gloeit" // nolint: mnd
    }


    frag_colour = vec4(result, 1.0); // nolint: mnd
}
	` + "\x00"
)

var (
	trafficLightData []float32
	rotX             float32
	rotY             float32
)

func main() {
	window := initGlfw()

	defer glfw.Terminate()

	program := initOpenGL()

	// STAP 1: Genereer de data
	trafficLightData = generateTrafficLightData()

	// STAP 2: Maak de VAO en geef de data mee
	vao := makeVao(trafficLightData)

	// Om diepte goed te zien (3D)
	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE) // Teken zowel voor- als achterkant van de driehoeken

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
		// STAP 3: Teken de data
		draw(vao, window, program)
	}
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.ClearColor(0.2, 0.2, 0.3, 0.9) // nolint: mnd
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// Matrices
	// nolint: lll,mnd
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 10.0)
	view := mgl32.LookAtV(mgl32.Vec3{0, 0, 5}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}) // nolint: mnd
	// nolint: lll,mnd
	model := mgl32.HomogRotate3DX(mgl32.DegToRad(rotX)).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotY)))

	// STUUR DE MATRICES CORRECT
	// We hebben in de shader 'uniform mat4 model' nodig voor positie
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &model[0]) // nolint: lll,mnd

	// We sturen de gecombineerde MVP naar een nieuwe naam (of we passen de shader aan)
	// Voor nu: we gebruiken de shader die gl_Position = model * vec4(vp, 1.0) doet.
	// Dat betekent dat 'model' daar eigenlijk MVP moet zijn.
	mvp := projection.Mul4(view).Mul4(model)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &mvp[0]) // nolint: lll,mnd

	// LICHT (De zon staat rechtsboven de camera)
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("lightPos\x00")), 2.0, 2.0, 5.0)   // nolint: mnd
	gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("emission\x00")), 0.0)             // Standaard uit // nolint: mnd
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("lightColor\x00")), 1.0, 1.0, 1.0) // nolint: mnd
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("viewPos\x00")), 0, 0, 5)          // nolint: mnd

	gl.BindVertexArray(vao)

	// --- DEEL 1: HET ZWARTE BORD (Eerste 6 vertices) ---
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("objectColor\x00")), 0.05, 0.05, 0.05) // nolint: mnd
	gl.DrawArrays(gl.TRIANGLES, 0, 6)                                                         // nolint: mnd

	// --- DEEL 2: DE WITTE BIES (Volgende 6 vertices) ---
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("objectColor\x00")), 0.9, 0.9, 0.9) // BIJNA WIT // nolint: mnd
	gl.DrawArrays(gl.TRIANGLES, 6, 6)                                                      // nolint: mnd

	// --- DEEL 3: DE KAPPEN (De rest) ---
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("objectColor\x00")), 0.1, 0.1, 0.1) // nolint: mnd

	numVertices := int32(len(trafficLightData) / 6) // nolint: gosec,mnd
	gl.DrawArrays(gl.TRIANGLES, 12, numVertices-12) // nolint: mnd

	// 4. DE LAMPEN (Hier komt de kleur!)
	gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("emission\x00")), 1.0) // Zet "aan" // nolint: mnd

	// Rood
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("objectColor\x00")), 1.0, 0.0, 0.0) // nolint: mnd
	gl.DrawArrays(gl.TRIANGLES, 300, 48)                                                   // nolint: mnd
	// Oranje
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("objectColor\x00")), 1.0, 0.5, 0.0) // nolint: mnd
	gl.DrawArrays(gl.TRIANGLES, 348, 48)                                                   // nolint: mnd
	// Groen
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("objectColor\x00")), 0.0, 1.0, 0.2) // nolint: mnd
	gl.DrawArrays(gl.TRIANGLES, 396, 48)                                                   // nolint: mnd

	glfw.PollEvents()
	window.SwapBuffers()
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

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	// compile shaders
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32

	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32

		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
