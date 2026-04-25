package trafficlight

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type TrafficLight struct {
	RotX, RotY float32
	RedOn      bool
	OrangeOn   bool
	GreenOn    bool
	Ambient    float32
}

// Global variables for traffic light data and state
var (
	trafficLightData []float32
	rotX             float32
	rotY             float32
	State            TrafficLight // Exposed state
	GlProgram        uint32
)

const (
	width  = 250 // nolint: mnd
	height = 500 // nolint: mnd
)

// New creates a new TrafficLight with default values.
func New() TrafficLight {
	return TrafficLight{
		Ambient: 0.3, // nolint: mnd
	}
}

const (
	vertexShaderSource = `
#version 410
layout (location = 0) in vec3 vp;
layout (location = 1) in vec3 normal;

out vec3 Normal;
out vec3 FragPos;

uniform mat4 model;    // Alleen rotatie/positie
uniform mat4 transform; // De volledige Projection * View * Model

void main() {
    // FragPos in wereld-coördinaten voor lichtberekening
    FragPos = vec3(model * vec4(vp, 1.0));
    // Normalen correct meedraaien
    Normal = mat3(transpose(inverse(model))) * normal;

    gl_Position = transform * vec4(vp, 1.0);
}
` + "\x00"

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
uniform float ambientStrength;

void main() {
    // 1. Ambient
    // float ambientStrength = 0.3; // nolint: mnd
    vec3 ambient = ambientStrength * lightColor;

    // 2. Diffuse (De boosdoener: we maken hier de vec3 'diffuse')
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.2); // nolint: mnd
    vec3 diffuse = diff * lightColor;

    // 3. Specular (De glans)
    float specularStrength = 1.5; // nolint: mnd
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

// InitializeTrafficLight sets up OpenGL and GLFW, generates traffic
// light data, and creates the VAO.
// It returns the OpenGL program ID.
func InitializeTrafficLight() uint32 {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // nolint: mnd
	glfw.WindowHint(glfw.ContextVersionMinor, 1) // nolint: mnd
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Window creation is handled by the caller.

	GlProgram = initOpenGL()

	// Generate the data for the traffic light
	trafficLightData = generateTrafficLightData()

	// Create the VAO
	vao := makeVao(trafficLightData)
	_ = vao // VAO is used implicitly in draw functions, but we don't need to return it here.

	// Enable depth testing for 3D rendering
	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE) // Draw both front and back faces of triangles

	return GlProgram
}

// generateTrafficLightData generates the vertex data for the traffic light.
func generateTrafficLightData() []float32 {
	var data []float32

	// HULPFUNCTIE: Voegt een punt toe met positie en normaal
	addVertex := func(x, y, z, nx, ny, nz float32) {
		data = append(data, x, y, z, nx, ny, nz)
	}

	// 1. ACHTERBORD (Een simpele rechthoek/box)
	addVertex(-0.4, 1.0, 0.0, 0.0, 0.0, 1.0)  // nolint: mnd
	addVertex(0.4, 1.0, 0.0, 0.0, 0.0, 1.0)   // nolint: mnd
	addVertex(-0.4, -1.0, 0.0, 0.0, 0.0, 1.0) // nolint: mnd

	addVertex(0.4, 1.0, 0.0, 0.0, 0.0, 1.0)   // nolint: mnd
	addVertex(0.4, -1.0, 0.0, 0.0, 0.0, 1.0)  // nolint: mnd
	addVertex(-0.4, -1.0, 0.0, 0.0, 0.0, 1.0) // nolint: mnd

	// WITTE BIES (Een iets grotere rechthoek achter het zwarte bord)
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

			nx1, ny1 := float32(math.Cos(angle1)), float32(math.Sin(angle1))
			nx2, ny2 := float32(math.Cos(angle2)), float32(math.Sin(angle2))

			addVertex(x1, y1+yOffset, 0.0, nx1, ny1, 0.0)   // nolint: mnd
			addVertex(x2, y2+yOffset, 0.0, nx2, ny2, 0.0)   // nolint: mnd
			addVertex(x1, y1+yOffset, depth, nx1, ny1, 0.0) // nolint: mnd

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

			addVertex(0, yOffset, 0.01, 0, 0, 1) // nolint: mnd
			addVertex(float32(math.Cos(angle1))*radius, yOffset+float32(math.Sin(angle1))*radius,
				0.01, 0, 0, 1) // nolint: lll,mnd
			addVertex(float32(math.Cos(angle2))*radius, yOffset+float32(math.Sin(angle2))*radius,
				0.01, 0, 0, 1) // nolint: lll,mnd
		}
	}

	return data
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW) // nolint: mnd

	var vao uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0) // nolint: mnd
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0)) // nolint: mnd,staticcheck
	gl.EnableVertexAttribArray(0)                                       // nolint: mnd

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4)) // nolint: mnd,staticcheck
	gl.EnableVertexAttribArray(1)                                         // nolint: mnd

	return vao
}

// Draw renders the traffic light. It uses the current State and applies rotations.
func Draw(window *glfw.Window) {
	if GlProgram == 0 {
		log.Println("Error: OpenGL program not initialized. Call InitializeTrafficLight first.")

		return
	}

	gl.ClearColor(0.2, 0.2, 0.3, 0.9) // nolint: mnd
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(GlProgram)

	// Matrices
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), // nolint: mnd
		float32(width)/float32(height), 0.1, 10.0) // nolint: mnd
	view := mgl32.LookAtV(mgl32.Vec3{0, 0, 5},
		mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}) // nolint: mnd
	model := mgl32.HomogRotate3DX(mgl32.DegToRad(rotX)).Mul4(
		mgl32.HomogRotate3DY(mgl32.DegToRad(rotY)),
	)
	mvp := projection.Mul4(view).Mul4(model)

	gl.UniformMatrix4fv(gl.GetUniformLocation(GlProgram, gl.Str("model\x00")), 1, false, &model[0])
	gl.UniformMatrix4fv(gl.GetUniformLocation(GlProgram, gl.Str("transform\x00")), 1, false, &mvp[0])

	// LIGHT
	gl.Uniform3f(gl.GetUniformLocation(GlProgram, gl.Str("lightPos\x00")),
		0.0, 10.0, 2.0) // nolint: mnd
	gl.Uniform3f(gl.GetUniformLocation(GlProgram, gl.Str("lightColor\x00")),
		1.0, 1.0, 0.9) // nolint: mnd
	gl.Uniform3f(gl.GetUniformLocation(GlProgram, gl.Str("viewPos\x00")), 0, 0, 5) // nolint: mnd

	vao := makeVao(trafficLightData)
	gl.BindVertexArray(vao)

	// HET ZWARTE BORD (Eerste 6 vertices)
	gl.Uniform3f(gl.GetUniformLocation(GlProgram, gl.Str("objectColor\x00")),
		0.15, 0.15, 0.15) // nolint: mnd
	gl.Uniform1f(gl.GetUniformLocation(GlProgram, gl.Str("emission\x00")), 0.0) // nolint: mnd
	gl.DrawArrays(gl.TRIANGLES, 0, 6)                                           // nolint: mnd

	// DE WITTE BIES (Volgende 6 vertices)
	gl.Uniform3f(gl.GetUniformLocation(GlProgram, gl.Str("objectColor\x00")),
		0.9, 0.9, 0.9) // nolint: mnd
	gl.DrawArrays(gl.TRIANGLES, 6, 6) // nolint: mnd

	// DE KAPPEN (De rest)
	gl.Uniform3f(gl.GetUniformLocation(GlProgram, gl.Str("objectColor\x00")),
		0.1, 0.1, 0.1) // nolint: mnd
	gl.DrawArrays(gl.TRIANGLES, 12, 288) // nolint: mnd

	// Lampen definitie met kleur
	drawLamp(GlProgram, 300, mgl32.Vec3{1, 0, 0}, State.RedOn)      // nolint: mnd
	drawLamp(GlProgram, 348, mgl32.Vec3{1, 0.5, 0}, State.OrangeOn) // nolint: mnd
	drawLamp(GlProgram, 396, mgl32.Vec3{0, 1, 0.2}, State.GreenOn)  // nolint: mnd

	// ambientStrength doorgeven aan Shader
	gl.Uniform1f(gl.GetUniformLocation(GlProgram, gl.Str("ambientStrength\x00")), State.Ambient)

	glfw.PollEvents()
	window.SwapBuffers()
}

func drawLamp(program uint32, startIndex int32, color mgl32.Vec3, isOn bool) {
	if isOn {
		gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("emission\x00")), 1.0)
		gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("objectColor\x00")),
			color.X(), color.Y(), color.Z())
	} else {
		gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("emission\x00")), 0.0)
		gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("objectColor\x00")),
			color.X()*0.1, color.Y()*0.1, color.Z()*0.1) // nolint
	}

	gl.DrawArrays(gl.TRIANGLES, startIndex, 48) // nolint endindex of lamps
}

// initOpenGL initializes OpenGL and returns an initialized program.
// This function should be called once during application startup.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

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

// SetRotations updates the rotation angles for the traffic light.
func SetRotations(x, y float32) {
	rotX = x
	rotY = y
}

// ToggleRed toggles the state of the red light.
func ToggleRed() {
	State.RedOn = !State.RedOn
	fmt.Println("Rood getoggled:", State.RedOn)
}

// ToggleOrange toggles the state of the orange light.
func ToggleOrange() {
	State.OrangeOn = !State.OrangeOn
}

// ToggleGreen toggles the state of the green light.
func ToggleGreen() {
	State.GreenOn = !State.GreenOn
}

// SetAmbient sets the ambient light strength.
func SetAmbient(ambient float32) {
	State.Ambient = ambient
}
