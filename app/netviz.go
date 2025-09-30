package app

import (
	"edamame/core/networks"
	"math"
	"math/rand"
)

type Network_transform struct {
	Node_xs  []float32
	Node_ys  []float32
	Node_vxs []float32
	Node_vys []float32
}

func Create_network_transform(g *networks.Network, seed int64) Network_transform {

	local_rng := rand.New(rand.NewSource(seed))
	transform := Network_transform{make([]float32, 0), make([]float32, 0),
		make([]float32, 0), make([]float32, 0)}

	for _ = range g.Node_count {
		transform.Node_xs = append(transform.Node_xs, 300.0*local_rng.Float32())
		transform.Node_ys = append(transform.Node_ys, 300.0*local_rng.Float32())
		transform.Node_vxs = append(transform.Node_vxs, 0.0)
		transform.Node_vys = append(transform.Node_vys, 0.0)
	}

	return transform
}

func Get_COM(nt *Network_transform) (float32, float32) {
	var cx float32
	var cy float32
	for i, _ := range nt.Node_xs {
		cx += nt.Node_xs[i]
		cy += nt.Node_ys[i]
	}
	cx /= float32(len(nt.Node_xs))
	cy /= float32(len(nt.Node_ys))
	return cx, cy
}

func Spring_update(nt *Network_transform,
	g *networks.Network,
	attr, repl, ts float32) {
	fx := make([]float32, 0)
	fy := make([]float32, 0)
	for i := range g.Node_count {

		fx = append(fx, 0.0)
		fy = append(fy, 0.0)

		for j := range g.Node_count {

			node_a := g.Node_list[i]
			node_b := g.Node_list[j]
			a_x := nt.Node_xs[i]
			a_y := nt.Node_ys[i]
			b_x := nt.Node_xs[j]
			b_y := nt.Node_ys[j]

			r := math.Sqrt(float64((b_x-a_x)*(b_x-a_x) + (b_y-a_y)*(b_y-a_y)))

			var f float32

			if networks.Has_edge(g, node_a.Name, node_b.Name) {
				f = float32(r) * attr
			} else {
				if r < 0.01 {
					r = 0.01
				}
				f = float32(-1.0/(r*r)) * repl
			}

			theta := math.Atan2(float64(b_y-a_y), float64(b_x-a_x))
			fx[i] += f * float32(math.Cos(theta))
			fy[i] += f * float32(math.Sin(theta))
		}

		nt.Node_vxs[i] += fx[i]
		nt.Node_vys[i] += fy[i]

	}

	for i := range g.Node_count {
		nt.Node_xs[i] += nt.Node_vxs[i] * ts
		nt.Node_ys[i] += nt.Node_vys[i] * ts
	}

}
