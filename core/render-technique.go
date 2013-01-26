package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

var (
	techs     map[string]renderTechnique
	tmpEffect *FxEffect
	tmpMat    *FxMaterial
)

type techniqueCtor func(string) renderTechnique

func initTechniques() {
	techs = map[string]renderTechnique{}
	for techName, techMaker := range map[string]techniqueCtor{"rt_unlit": newTechnique_Unlit} {
		techs[techName] = techMaker(techName)
	}
}

type renderTechnique interface {
	initMeshBuffer(*MeshBuffer) []*ugl.VertexAttribPointer
	name() string
	onPreRender()
	onRenderMesh()
	onRenderMeshModel()
	onRenderNode()
}

type baseTechnique struct {
	prog *ugl.Program
}

func (me *baseTechnique) initMeshBuffer(meshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	atts = append(atts, ugl.NewVertexAttribPointer("aPos", me.prog.AttrLocs["aPos"], 3, 8*4, gl.Pointer(nil)))
	return
}

func (me *baseTechnique) name() string {
	return me.prog.Name
}

func (me *baseTechnique) onPreRender() {
}

func (me *baseTechnique) onRenderMesh() {
}

func (me *baseTechnique) onRenderMeshModel() {
}

func (me *baseTechnique) onRenderNode() {
}

func (me *baseTechnique) setProg(name string, unifs []string, attrs []string) {
	prog := glProgMan.Programs[name]
	prog.SetUnifLocations("uMatModelProj")
	if len(unifs) > 0 {
		prog.SetUnifLocations(unifs...)
	}
	prog.SetAttrLocations("aPos")
	if len(attrs) > 0 {
		prog.SetAttrLocations(attrs...)
	}
	me.prog = prog
}

type techniqueUnlit struct {
	baseTechnique
}

func newTechnique_Unlit(progName string) renderTechnique {
	me := &techniqueUnlit{}
	me.baseTechnique.setProg(progName, []string{"uDiffuse"}, []string{"aTexCoords"})
	return me
}

func (me *techniqueUnlit) initMeshBuffer(meshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	atts = me.baseTechnique.initMeshBuffer(meshBuffer)
	atts = append(atts, ugl.NewVertexAttribPointer("aTexCoords", me.prog.AttrLocs["aTexCoords"], 2, 8*4, gl.Offset(nil, 3*4)))
	return
}

func (me *techniqueUnlit) onRenderNode() {
	me.baseTechnique.onRenderNode()
	if tmpMat = curNode.EffectiveMaterial(); tmpMat != curMat {
		if curMat = tmpMat; curMat != nil {
			tmpEffect = Core.Libs.Effects[curMat.DefaultEffectID]
			gl.ActiveTexture(gl.TEXTURE0)
			Core.Libs.Images.I2D[tmpEffect.Diffuse.Texture.Image2ID].glTex.Bind()
			gl.Uniform1i(curProg.UnifLocs["uDiffuse"], 0)
		}
	}
}
