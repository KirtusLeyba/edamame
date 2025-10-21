import networkx as nx
import pandas as pd


n = 20000
p = 0.00001


G = nx.erdos_renyi_graph(n, p)

print("Saving random graph with {} nodes and {} edges...".format(len(G.nodes()), len(G.edges())))

with open("./sample_data/random_nodes.csv", "w") as fp:
    fp.write("name,radius,r,g,b,a\n")
    for node in G.nodes():
        fp.write("{},5.0,40,94,150,255\n".format(str(node)))

with open("./sample_data/random_edges.csv", "w") as fp:
    fp.write("nodeA,nodeB,width\n")
    for edge in G.edges():
        fp.write("{},{},2.5\n".format(edge[0], edge[1]))
