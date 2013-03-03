package core

import (
	"sync"

	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
)

var (
	thrApp struct {
		sync.Mutex
	}
	thrPrep struct {
		sync.Mutex
		// nodePreBatch nodeBatchPreps
	}
	thrRend struct {
		curCam                *Camera
		curView               *RenderView
		curEffect, nextEffect *FxEffect
		curTech, nextTech     RenderTechnique
		curProg               *ugl.Program
		quadTex               gl.Uint
	}
)

func init() {
}

func (_ *NgCore) copyAppToPrep() {
	for cid := 0; cid < len(Core.Render.Canvases); cid++ {
		if Core.Render.Canvases[cid].renderThisFrame() {
			Core.Render.Canvases[cid].copyAppToPrep()
		}
	}
}

func (_ *NgCore) copyPrepToRend() {
	for cid := 0; cid < len(Core.Render.Canvases); cid++ {
		if Core.Render.Canvases[cid].renderThisFrame() {
			Core.Render.Canvases[cid].copyPrepToRend()
		}
	}
}

func (me *RenderCanvas) copyAppToPrep() {
	for view := 0; view < len(me.Views); view++ {
		me.Views[view].copyAppToPrep()
	}
}

func (me *RenderCanvas) copyPrepToRend() {
	for view := 0; view < len(me.Views); view++ {
		me.Views[view].copyPrepToRend()
	}
}

func (me *RenderView) copyAppToPrep() {
	me.Technique.copyAppToPrep()
}

func (me *RenderView) copyPrepToRend() {
	me.Technique.copyPrepToRend()
}

func (me *RenderTechniqueScene) copyAppToPrep() {
	me.Camera.copyAppToPrep()
}

func (me *RenderTechniqueScene) copyPrepToRend() {
	me.Camera.copyPrepToRend()
}

func (me *Camera) copyAppToPrep() {
	me.thrPrep.matProj = me.thrApp.matProj
	me.Controller.thrPrep.mat = me.Controller.thrApp.mat
	me.thrPrep.matPos.Translation(&me.Controller.Pos)
	if scene := me.scene(); scene != nil && !scene.thrPrep.copyDone {
		scene.thrPrep.copyDone = true
		scene.RootNode.copyAppToPrep()
	}
}

func (me *Camera) copyPrepToRend() {
	if scene := me.scene(); scene != nil && !scene.thrRend.copyDone {
		scene.thrRend.copyDone = true
		scene.RootNode.copyPrepToRend()
		scene.thrPrep.done = false
	}
}

func (me *Node) copyAppToPrep() {
	me.thrPrep.matModelView = me.Transform.matModelView
	for _, subNode := range me.ChildNodes.M {
		subNode.copyAppToPrep()
	}
}

func (me *Node) copyPrepToRend() {
	var (
		cam *Camera
		mat *unum.Mat4
		cr  bool
	)
	for cam, cr = range me.thrPrep.camRender {
		me.thrRend.camRender[cam] = cr
	}
	for cam, mat = range me.thrPrep.camProjMats {
		me.thrRend.camProjMats[cam].Load(mat)
	}
	for _, subNode := range me.ChildNodes.M {
		subNode.copyPrepToRend()
	}
}
