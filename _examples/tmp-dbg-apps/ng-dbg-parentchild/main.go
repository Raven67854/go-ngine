package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floorID, boxID, pyrID int
	tmpScene              *ng.Scene
	tmpNode               *ng.SceneNode
)

func main() {
	apputil.Main(setupScene, onAppThread, onWinThread)
}

func onAppThread() {
	apputil.HandleCamCtlKeys()
	if tmpScene = apputil.SceneCam.Scene(); tmpScene != nil {
		if tmpNode = tmpScene.Node(boxID); tmpNode != nil {
			tmpNode.Transform.Pos.Y = (2.125 + (0.5 * math.Sin(ng.Loop.Tick.Now*8)))
			tmpNode.Transform.Rot.Y += 0.005
			tmpNode.Transform.SetScale(1.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*2)))
			tmpScene.ApplyNodeTransforms(tmpNode.ID)
		}
	}
}

func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["pulse"]].GetTex2D(1).SetMixWeight(0.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*4)))
}

func setupScene() {
	var (
		meshFloorID, meshBoxID, meshPyrID int
		err                               error
		bufRest                           *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
		"cat":     "tex/cat.png",
		"gopher":  "tex/gopher.png",
		"crate":   "tex/crate.jpeg",
	})
	fxPulseID := ng.Core.Libs.Effects.AddNew()
	apputil.LibIDs.Fx["pulse"] = fxPulseID
	fxPulse := &ng.Core.Libs.Effects[fxPulseID]
	fxPulse.EnableTex2D(0).Tex_SetImageID(apputil.LibIDs.Img2D["crate"])
	fxPulse.EnableTex2D(1).Tex_SetImageID(apputil.LibIDs.Img2D["gopher"]).SetMixWeight(0.5)
	fxPulse.UpdateRoutine()

	dogMat := &ng.Core.Libs.Materials[apputil.LibIDs.Mat["dog"]]
	dogMat.FaceEffects.ByTag["top"] = apputil.LibIDs.Fx["cat"]
	dogMat.FaceEffects.ByTag["front"] = fxPulseID
	dogMat.FaceEffects.ByTag["back"] = fxPulseID

	//	meshes / models
	if bufRest, err = ng.Core.Mesh.Buffers.AddNew("buf_rest", 200); err != nil {
		panic(err)
	}
	if meshFloorID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.Core.Mesh.Desc.Plane); err != nil {
		panic(err)
	}
	if meshBoxID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_box", ng.Core.Mesh.Desc.Cube); err != nil {
		panic(err)
	}
	if meshPyrID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_pyr", ng.Core.Mesh.Desc.Pyramid); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloorID)
	bufRest.Add(meshBoxID)
	bufRest.Add(meshPyrID)

	scene := apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshPyrID)
	floor := apputil.AddNode(scene, 0, meshFloorID, apputil.LibIDs.Mat["cobbles"], -1)
	floorID = floor.ID
	floor.Transform.SetScale(100)

	box := apputil.AddNode(scene, 0, meshBoxID, apputil.LibIDs.Mat["dog"], -1)
	boxID = box.ID
	box.Transform.Pos.Y = 1.25

	pyr := apputil.AddNode(scene, boxID, meshPyrID, apputil.LibIDs.Mat["cat"], -1)
	pyrID = pyr.ID
	pyr.Transform.Pos.Y = 2.125

	scene.ApplyNodeTransforms(0)
	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-2.5, 2, -7)
	camCtl.EndUpdate()
}
