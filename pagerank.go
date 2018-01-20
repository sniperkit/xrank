package main

import (
	"fmt"
	"math"
)

/*
 * ## グラフの基礎知識
 *
 * - ノード(node): ネットワークの頂点
 * - リンク(link): ネットワークの頂点と頂点を結ぶ線
 * - 有向グラフ(directed graph): リンクのつながりに方向があるグラフ
 *     - ex) Twitterのフォロー関係
 * - 無向グラフ(undirected graph): リンクのつながりに方向がないグラフ
 *     - ex) Facebookの友人関係
 * - グラフ構造の値: 各ノードとリンクは値を持つことができる
 *     - ex) 地点Aの人口
 *     - ex) 地点Aから地点Bまでの距離
 *
 *
 * ex) 3つのノードを持つ有向グラフ
 *
 *  <------------,
 * 1<----,       |
 * └─ -> 2 <--> 3
 *
 * - nodes: {1, 2, 3}
 * - links: { 1: {2}, 2: {1, 3}, 3: {1, 2}}
 * - adjacencyList(matrix):
 *  {1: {1: 0,   2: 1,   3: 0},
 *   2: {1: 0.5, 2: 0,   3: 0.5}
 *   3: {1: 0.5, 2: 0.5, 3: 0}}
 *
 * |  |1  |2  |3  |
 * |:-|--:|--:|--:|
 * |1 |  0|  1|  0|
 * |2 |1/2|  0|1/2|
 * |3 |1/2|1/2|  0|
 *
 */

type Links map[int][]int
type Nodes map[int]float64
type AdjacencyList map[int]map[int]float64

func main() {
	// rankの初期値
	const v = float64(1.0)

	// 有向グラフのリンク
	//   ex) map[int][]int {fromNodeId: {toNodeIds}}
	links := Links{
		1: {2, 3},
		2: {3},
		3: {1, 4, 5},
		4: {1, 2, 3, 4, 5, 6},
		5: {6},
		6: {5},
	}
	printLinks(links)

	// 有向グラフのノードとそのRank
	//	 ex) map[int]float64 {nodeId: rank}
	nodes := toNodes(links)

	// 各ノードからノードへの遷移確率をもつ隣接行列G
	G := meyersAdjacencyList(links, nodes)

	// 50 step以内にグラフのrankが収束するはず
	for step := 1; step < 50; step++ {
		fmt.Printf("\n===== step %v =====\n", step)
		nodes = updateRank(nodes, G)
		printRank(nodes)
	}
}

func Round(f float64) float64 {
	const place = 15
	shift := math.Pow(10, float64(place))
	return math.Floor(f*shift+.5) / shift
}

// リンクからノードに変換
func toNodes(links Links) Nodes {
	const v = float64(1.0)
	nodes := Nodes{}

	for fromId, toIds := range links {
		nodes[fromId] = v
		for _, id := range toIds {
			nodes[id] = v
		}
	}
	return nodes
}

// リンク、ノードから隣接行列に変換する
func toMatrix(links Links, nodes Nodes) AdjacencyList {
	matrix := AdjacencyList{}

	// nodeIdの隣接行列を0で埋める
	for i, _ := range nodes {
		tmpMap := map[int]float64{}
		for j, _ := range nodes {
			tmpMap[j] = 0
		}
		matrix[i] = tmpMap
	}
	return matrix
}

/*
 * 遷移確率の隣接行列S
 *   リンクがないノードjへの遷移確率は0とする
 *   ai = {0: リンクあり, 1: リンクなし}
 *
 *   Niはノードiのリンク数
 *   nはノード数
 *
 *   あるノードiからあるノードjへの遷移確率pij
 *   pij = (1/Ni) + ai(1/n)
 *
 * |  |1  |2  |3  |4  |5  |6  |
 * |:-|--:|--:|--:|--:|--:|--:|
 * |1 |  0|1/2|1/2|  0|  0|  0|
 * |2 |  0|  0|  1|  0|  0|  0|
 * |3 |1/3|  0|  0|1/3|1/3|  0|
 * |4 |1/6|1/6|1/6|1/6|1/6|1/6|
 * |5 |  0|  0|  0|  0|  0|  1|
 * |6 |  0|  0|  0|  0|  1|  0|
 */
func probabilityAdjacencyList(links Links, nodes Nodes) AdjacencyList {
	matrix := toMatrix(links, nodes)

	for row, link := range links {
		ni := float64(len(link))
		for _, col := range link {
			matrix[row][col] = 1.0 / ni
		}
	}
	return matrix
}

/*
 * Meyer's Random Surfer Model
 *   damping factorを考慮した遷移確率
 *   d(damping factor)は、リンクがなくても一定確率で別ノードに遷移する確率
 *
 *   nはノード数
 *
 *   あるノードiからノードjへの遷移確率pij
 *   pij = d * Sij + (1 - d) * (1 / n)
 *
 * 以下は d = 0.9の場合の隣接行列
 * |  |1    |2    |3    |4    |5    |6    |
 * |:-|----:|----:|----:|----:|----:|----:|
 * |1 | 1/60|28/60|28/60| 1/60| 1/60| 1/60|
 * |2 | 1/60| 1/60|55/60| 1/60| 1/60| 1/60|
 * |3 |19/60| 1/60| 1/60|19/60|19/60| 1/60|
 * |4 |10/60|10/60|10/60|10/60|10/60|10/60|
 * |5 | 1/60| 1/60| 1/60| 1/60| 1/60|55/60|
 * |6 | 1/60| 1/60| 1/60| 1/60|55/60| 1/60|
 */
func meyersAdjacencyList(links Links, nodes Nodes) AdjacencyList {
	const d = float64(0.85)

	S := probabilityAdjacencyList(links, nodes)
	G := toMatrix(links, nodes)
	n := float64(len(G))

	for i, cols := range S {
		for j, _ := range cols {
			s := S[i][j]
			G[i][j] = d*(s) + (1.0-d)*(1.0/n)
		}
	}
	return G
}

func copyNodeKey(nodes Nodes) Nodes {
	newnodes := Nodes{}
	for k, _ := range nodes {
		newnodes[k] = float64(0)
	}
	return newnodes
}

// ノードからノードへ遷移確率pの分のrankをそれぞれ配分する
func updateRank(nodes Nodes, G AdjacencyList) Nodes {
	nextNodes := copyNodeKey(nodes)

	for i, rank := range nodes {
		for j, p := range G[i] {
			nextNodes[j] += p * rank
		}
	}

	return nextNodes
}

func printLinks(links Links) {
	for k, links := range links {
		fmt.Println("key: ", k, "links: ", links)
	}
}

func printRank(nodes Nodes) {
	for id, rank := range nodes {
		fmt.Println("id: ", id, ",rank: ", rank)
	}
}
