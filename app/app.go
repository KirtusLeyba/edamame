package app

import (
	"edamame/core/networks"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Execute() {

    g := networks.Create_rand_network(4, 1000, 0.0025, false)
    g_transform := Create_network_transform(&g, 1717)

    var width int32 = 4096
    var height int32 = 2160

    rl.InitWindow(width, height, "test")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)

    var iter int = 0

    for !rl.WindowShouldClose(){

        if(iter < 1000){
            Spring_update(&g_transform, &g, 0.005, 0.01, 0.01)
        } else if (iter == 1000) {
            rl.TakeScreenshot("./test.png")
        }
        cx, cy := Get_COM(&g_transform)

        rl.BeginDrawing()
        rl.ClearBackground(rl.Black)
        for i := range g.Node_count {
            for j := range g.Adjacencies[g.Node_list[i]] {

                other_node_idx, _ := strconv.Atoi(g.Adjacencies[g.Node_list[i]][j].Name)

                rl.DrawLine(int32(g_transform.Node_xs[i] - cx) + width/2, int32(g_transform.Node_ys[i] - cy) + height/2,
                            int32(g_transform.Node_xs[other_node_idx] - cx) + width/2, int32(g_transform.Node_ys[other_node_idx] - cy) + height/2, rl.Gray)
            }

        }
        for i := range g.Node_count {
            rl.DrawCircle(int32(g_transform.Node_xs[i] - cx) + width/2, int32(g_transform.Node_ys[i] - cy) + height/2, 8.0, rl.Black)
            rl.DrawCircle(int32(g_transform.Node_xs[i] - cx) + width/2, int32(g_transform.Node_ys[i] - cy) + height/2, 6.0, rl.Blue)
        }

        rl.EndDrawing()

        iter++
    }

}
