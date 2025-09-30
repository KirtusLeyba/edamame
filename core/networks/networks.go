package networks

import (
	"errors"
	"math/rand"
	"strconv"
)

type Node struct {
	Name string
}

type Network struct {
	Node_count  uint
	Node_list   []Node
	Adjacencies map[Node][]Node
	Directed    bool
}

func Create_empty_network(directed bool) Network {
	net := Network{0, make([]Node, 0), make(map[Node][]Node), directed}
	return net
}

func Add_node(net *Network, name string) error {
	_, ok := net.Adjacencies[Node{name}]
	if ok {
		return errors.New("attempted to add node named " + name + " twice.")
	}

	new_node := Node{name}
	net.Node_list = append(net.Node_list, new_node)
	net.Adjacencies[new_node] = make([]Node, 0)
	net.Node_count++
	return nil
}

func Has_node(net *Network, name string) bool {
	_, ok := net.Adjacencies[Node{name}]
	return ok
}

func Has_edge(net *Network, node_a_name, node_b_name string) bool {
	if !Has_node(net, node_a_name) {
		return false
	}
	if !Has_node(net, node_b_name) {
		return false
	}
	for _, node := range net.Adjacencies[Node{node_a_name}] {
		if node.Name == node_b_name {
			return true
		}
	}
	return false
}

func Add_edge(net *Network, node_a_name, node_b_name string) error {
	if !Has_node(net, node_a_name) {
		return errors.New("Cannot add edge with nonexistant node " + node_a_name)
	}
	if !Has_node(net, node_b_name) {
		return errors.New("Cannot add edge with nonexistant node " + node_b_name)
	}

	if net.Directed {
		if !Has_edge(net, node_a_name, node_b_name) {
			net.Adjacencies[Node{node_a_name}] = append(net.Adjacencies[Node{node_a_name}], Node{node_b_name})
		} else {
			return errors.New("Edge from " + node_a_name + " and " + node_b_name + " already exists.")
		}
	} else {
		if !Has_edge(net, node_a_name, node_b_name) {
			net.Adjacencies[Node{node_a_name}] = append(net.Adjacencies[Node{node_a_name}], Node{node_b_name})
		} else {
			return errors.New("Edge between " + node_a_name + " and " + node_b_name + " already exists.")
		}
		if !Has_edge(net, node_b_name, node_a_name) {
			net.Adjacencies[Node{node_b_name}] = append(net.Adjacencies[Node{node_b_name}], Node{node_a_name})
		} else {
			return errors.New("Edge between " + node_b_name + " and " + node_a_name + " already exists.")
		}
	}
	return nil
}

func Create_rand_network(seed int64, num_nodes int, edge_chance float32, directed bool) Network {
	local_rng := rand.New(rand.NewSource(seed))

	g := Create_empty_network(directed)

	for i := range num_nodes {
		Add_node(&g, strconv.Itoa(i))
	}

	if directed {
		for i := range num_nodes {
			for j := range num_nodes {
				if i == j {
					continue
				}
				if local_rng.Float32() < edge_chance {
					Add_edge(&g, strconv.Itoa(i), strconv.Itoa(j))
				}
			}
		}
	} else {
		for i := range num_nodes {
			for j := i + 1; j < num_nodes; j++ {
				if local_rng.Float32() < edge_chance {
					Add_edge(&g, strconv.Itoa(i), strconv.Itoa(j))
				}
			}
		}
	}
	return g
}
