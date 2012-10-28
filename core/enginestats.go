package core

import (
	glfw "github.com/go-gl/glfw"
)

type tEngineStats struct {
	FpsLastSec int
	Frame, FrameRenderBoth, FrameRenderCpu, FrameRenderGpu, FrameCoreCode, FrameUserCode, Gc *TEngineStatsTiming

	fpsCounter int
	fpsAll float64
}

	func newEngineStats () *tEngineStats {
		var stats = &tEngineStats {}
		stats.reset()
		return stats
	}

	func (me *tEngineStats) AverageFps () float64 {
		return me.fpsAll / glfw.Time()
	}

	func (me *tEngineStats) reset () {
		var nt = func () *TEngineStatsTiming { return &TEngineStatsTiming {} }
		me.FpsLastSec, me.fpsCounter, me.fpsAll = 0, 0, 0
		me.Frame, me.FrameRenderBoth, me.FrameRenderCpu, me.FrameRenderGpu, me.FrameCoreCode, me.FrameUserCode, me.Gc = nt(), nt(), nt(), nt(), nt(), nt(), nt()
	}

type TEngineStatsTiming struct {
	max, measuredCounter, measureStartTime, thisTime, totalAccum float64
	comb1, comb2 *TEngineStatsTiming
}

	func (me *TEngineStatsTiming) Average () float64 {
		return me.totalAccum / me.measuredCounter
	}

	func (me *TEngineStatsTiming) combine () {
		me.max = me.comb1.max + me.comb2.max
		me.measuredCounter = (me.comb1.measuredCounter + me.comb2.measuredCounter) * 0.5
		me.totalAccum = me.comb1.totalAccum + me.comb2.totalAccum
	}

	func (me *TEngineStatsTiming) begin () {
		me.measureStartTime = glfw.Time()
	}

	func (me *TEngineStatsTiming) end () {
		if me.thisTime = glfw.Time() - me.measureStartTime; me.thisTime > me.max { me.max = me.thisTime }
		me.measuredCounter++
		me.totalAccum += me.thisTime
	}

	func (me *TEngineStatsTiming) Max () float64 {
		return me.max
	}