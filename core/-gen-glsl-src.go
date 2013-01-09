package core

func init() {
	var rss glShaderSources
	rss.init()
	rss.vertex["postfx"] = "in vec3 aPos;\n\nvoid main () {\n\tgl_Position = vec4(aPos, 1.0);\n}\n"
	rss.vertex["rt_unlit_colored"] = "in vec3 aPos;\n\nuniform mat4 uMatCam;\nuniform mat4 uMatModelView;\nuniform mat4 uMatProj;\n\nsmooth out vec3 vPos;\n\nvoid main () {\n\tvPos = (aPos * 0.33) + 0.66;\n\tgl_Position = uMatProj * uMatCam * uMatModelView * vec4(aPos, 1.0);\n}\n"
	rss.vertex["rt_unlit_textured"] = "in vec3 aPos;\nin vec2 aTexCoords;\n\nuniform mat4 uMatCam;\nuniform mat4 uMatModelView;\nuniform mat4 uMatProj;\n\nout vec2 vTexCoords;\n\nvoid main () {\n\tgl_Position = uMatProj * uMatCam * uMatModelView * vec4(aPos, 1.0);\n\tvTexCoords = aTexCoords;\n}\n"
	rss.fragment["postfx"] = "out vec4 oColor;\n\nvoid main () {\n\toColor = vec4(0.66, 0.99, 0.33, 1.0);\n}\n"
	rss.fragment["rt_unlit_colored"] = "smooth in vec3 vPos;\n\nout vec4 oColor;\n\nvoid main () {\n\toColor = vec4(vPos, 1.0);\n}\n"
	rss.fragment["rt_unlit_textured"] = "in vec2 vTexCoords;\n\nuniform sampler2D uTex0;\n\nout vec4 oColor;\n\nvoid main () {\n\toColor = texture(uTex0, vTexCoords);\n}\n"
	glShaderMan.sources = rss
	glShaderMan.names = []string{"postfx", "rt_unlit_colored", "rt_unlit_textured"}
}