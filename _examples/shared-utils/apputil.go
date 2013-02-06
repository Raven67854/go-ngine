package exampleutils

import (
	"fmt"
	"runtime"

	ng "github.com/go3d/go-ngine/core"
)

var (
	//	The RenderCanvas the example scene is initially being rendered to. This is an off-screen "render-to-texture" RenderCanvas.
	SceneCanvas *ng.RenderCanvas

	//	The primary scene-rendering camera, rendering to SceneCanvas.
	SceneCam *ng.Camera

	//	Unlike the off-screen (render-to-texture) SceneCanvas above, this RenderCanvas epresents the actual screen/window.
	PostFxCanvas *ng.RenderCanvas

	//	Takes the image rendered to SceneCanvas, may post-process it or not, and blits it to PostFxCanvas.
	PostFxCam *ng.Camera

	//	Do not set this field directly, only use PauseResume() to toggle it and effect the associated render-state changes.
	Paused bool

	retro, refreshWinTitle bool
	curKeyHint             = 0
	sec                    = 0
)

//	Refreshes the window title every second, showing the next one entry in KeyHints every 3 seconds.
func OnSec() {
	if sec++; sec == 3 {
		sec = 0
		if curKeyHint++; (curKeyHint > MaxKeyHint) || (curKeyHint >= (len(KeyHints))) {
			curKeyHint = 0
		}
	}
	ng.UserIO.SetWinTitle(WindowTitle())
}

//	Called by each example-app's func main(). Initializes go:ngine, sets SceneCam/SceneCanvas/PostFxCam/PostFxCanvas etc., calls the specified assetLoader function, then enters The Loop.
func Main(setupExampleScene, onAppThread, onWinThread func()) {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())

	width, height, fullscreen := 1280, 720, false
	// width, height, fullscreen := 1920, 1080, true
	opt := ng.NewEngineOptions(AssetRootDirPath(), width, height, 0, fullscreen)

	// While the default for this (force GL core profile on Macs only) is reasonable for "real-world" apps at present, for
	// the example apps we force core profile to implicitly ensure all of go:ngine's GL code is fully core-profile compliant
	opt.Initialization.GlContext.CoreProfile = true
	opt.Loop.ForceThreads.App, opt.Loop.ForceThreads.Prep = true, true

	if err := ng.Init(opt, fmt.Sprintf("Loading example app... (%v CPU cores)", runtime.GOMAXPROCS(0))); err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		defer ng.Dispose()
		ng.Loop.On.EverySec, ng.Loop.On.AppThread, ng.Loop.On.WinThread = OnSec, onAppThread, onWinThread

		PostFxCanvas = ng.Core.Rendering.Canvases.Final()
		PostFxCam = PostFxCanvas.AddNewCameraQuad()

		SceneCanvas = ng.Core.Rendering.Canvases.AddNew(true, 1, 1)
		SceneCam = SceneCanvas.AddNewCamera3D()
		SceneCam.Rendering.States.ClearColor.Set(0.5, 0.6, 0.85, 1)

		setupExampleScene()
		ng.Core.SyncUpdates()
		ng.Loop.Loop()
		PrintPostLoopSummary()
	}
}

//	Returns the window title to be set by OnSec().
func WindowTitle() string {
	cw, ch := SceneCanvas.CurrentAbsoluteSize()
	camPos, camDir := SceneCam.Controller.Pos, SceneCam.Controller.Dir()
	return fmt.Sprintf("%v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", ng.Stats.FpsLastSec, cw, ch, KeyHints[curKeyHint], camPos.String(), camDir.String())
}