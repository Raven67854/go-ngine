package main

import (
	"math"

	ng "github.com/go3d/go-ngine/core"
	nga "github.com/go3d/go-ngine/assets"
	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	floor, pyr, box *ng.Node
)

func main () {
	ngsamples.SamplesMainFunc(LoadSampleScene_02_PyrCube)
}

func onLoop () {
	ngsamples.CheckToggleKeys()
	ngsamples.CheckCamCtlKeys()

	//	animate mesh nodes
	pyr.NodeTransform.Rot.X -= 0.0005
	pyr.NodeTransform.Rot.Y -= 0.0005
	pyr.NodeTransform.Pos.Set(-13.75, 2 * math.Sin(ng.Loop.TickNow), 2)
	pyr.NodeTransform.OnPosRotChanged()

	box.NodeTransform.Rot.Y += 0.0004
	box.NodeTransform.Rot.Z += 0.0006
	box.NodeTransform.Pos.Set(-8.125, 2 * math.Cos(ng.Loop.TickNow), -2)
	box.NodeTransform.OnPosRotChanged()
}

func LoadSampleScene_02_PyrCube () {
	var err error
	var meshFloor, meshPyr, meshCube *ng.Mesh
	var bufFloor, bufRest *ng.MeshBuffer

	ng.Loop.OnLoop = onLoop
	ngsamples.Cam.Options.BackfaceCulling = false

	//	textures / materials
	ng.Core.Textures["tex_cobbles"] = ng.Core.Textures.LoadAsync(ng.TextureProviders.RemoteFile, "http://dl.dropbox.com/u/136375/go-ngine/assets/tex/cobbles.png")
	ng.Core.Textures["tex_crate"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/crate.jpeg")
	ng.Core.Textures["tex_mosaic"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/mosaic.jpeg")
	nga.Materials["mat_cobbles"] = nga.Materials.New("tex_cobbles")

	nga.Materials["mat_crate"] = nga.Materials.New("tex_crate")

	nga.Materials["mat_mosaic"] = nga.Materials.New("tex_mosaic")


	//	meshes / models
	if bufFloor, err = ng.Core.MeshBuffers.Add("buf_floor", ng.Core.MeshBuffers.NewParams(6, 6)); err != nil { panic(err) }
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(36 + 12, 36 + 12)); err != nil { panic(err) }
	if meshFloor, err = ng.Core.Meshes.Load("mesh_plane", ng.MeshProviders.PrefabPlane); err != nil { panic(err) }
	if meshPyr, err = ng.Core.Meshes.Load("mesh_pyramid", ng.MeshProviders.PrefabPyramid); err != nil { panic(err) }
	if meshCube, err = ng.Core.Meshes.Load("mesh_cube", ng.MeshProviders.PrefabCube); err != nil { panic(err) }
	ng.Core.Meshes.AddRange(meshFloor, meshPyr, meshCube)
	bufFloor.Add(meshFloor); bufRest.Add(meshCube); bufRest.Add(meshPyr)
	meshPyr.Models.Default().SetMatName("mat_mosaic")
	meshCube.Models.Default().SetMatName("mat_crate")

	//	scene
	var scene = ng.NewScene()
	ng.Core.Scenes[""] = scene
	scene.RootNode.SubNodes.MakeN("node_floor", "mesh_plane", "", "node_pyr", "mesh_pyramid", "", "node_box", "mesh_cube", "")
	floor, pyr, box = scene.RootNode.SubNodes.M["node_floor"], scene.RootNode.SubNodes.M["node_pyr"], scene.RootNode.SubNodes.M["node_box"]

	floor.SetMatName("mat_cobbles")
	floor.NodeTransform.SetPosXYZ(0.1, 0, -8)
	floor.NodeTransform.SetScalingN(1000)

	ngsamples.CamCtl.BeginUpdate(); ngsamples.CamCtl.Pos.Y = 1.6; ngsamples.CamCtl.EndUpdate()

	//	upload everything
	ng.Core.SyncUpdates()
}