package app

type Vec2Df32 struct {
	x float32
	y float32
}
type Vec2Di struct {
	x int
	y int
}

type Layer interface {
	Attach()
	Detach()
	OnEvent()
	OnUpdate()
	OnRender()
	SetTransform(origin, size Vec2Df32)
}
