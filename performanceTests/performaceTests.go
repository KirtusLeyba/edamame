package main

import (
	ednet "edamame/core/networks"
	"fmt"
	"math"
	"time"
)

func sep() {
	fmt.Printf("==============================\n\n")
}

func main() {
	trials := 5

	// fmt.Printf("Testing the edamame networks!\n")
	// sep()
	//
	// fmt.Printf("Simple serial method:\n")
	// for i := range trials{
	// 	numNodes := 100*int(math.Pow(2, float64(i)))
	// 	n := ednet.NewRandomSpatialNet(numNodes, 0.1)
	// 	start := time.Now()
	// 	for range 100 {
	// 		n.SpringUpdate(0.1, 0.1, 1.0, 1.0, 0.001)
	// 	}
	// 	elapsed := time.Since(start)
	// 	fmt.Printf("Nodes: %v, Elapsed time: %v seconds\n", numNodes, elapsed.Seconds())
	// }
	// sep()
	//
	fmt.Printf("Using go routines:\n")
	for i := range trials {
		numNodes := 100 * int(math.Pow(2, float64(i)))
		n := ednet.NewRandomSpatialNet(numNodes, 0.1)
		start := time.Now()
		for range 100 {
			n.SpringUpdateParallel(0.1, 0.1, 1.0, 1.0, 0.001)
		}
		elapsed := time.Since(start)
		fmt.Printf("Nodes: %v, Elapsed time: %v seconds\n", numNodes, elapsed.Seconds())
	}
	sep()

	// fmt.Printf("Using serial spatial hashing\n")
	// for i := range trials {
	// 	numNodes := 100 * int(math.Pow(2, float64(i)))
	// 	n := ednet.NewRandomSpatialNet(numNodes, 0.1)
	// 	start := time.Now()
	// 	for range 100 {
	// 		n.SpringUpdateHashing(0.1, 0.1, 1.0, 1.0, 0.001)
	// 	}
	// 	elapsed := time.Since(start)
	// 	fmt.Printf("Nodes: %v, Elapsed time: %v seconds\n", numNodes, elapsed.Seconds())
	// }
	// sep()

	fmt.Printf("Using parallel spatial hashing\n")
	for i := range trials {
		numNodes := 100 * int(math.Pow(2, float64(i)))
		n := ednet.NewRandomSpatialNet(numNodes, 0.1)
		var binSize float32 = 1000.0
		n.ResetSpatialHashing(binSize)
		start := time.Now()
		for range 100 {
			n.SpringUpdateHashingParallel(0.1, 0.1, 1.0, 1.0, 0.001, binSize)
		}
		elapsed := time.Since(start)
		fmt.Printf("Nodes: %v, Elapsed time: %v seconds\n", numNodes, elapsed.Seconds())
	}
	sep()

}
