package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	name string
	x, y int
}

type Graph struct {
	nodes     map[string]Node
	adjMatrix map[string]map[string]int
}

func NewGraph() *Graph {
	return &Graph{
		nodes:     make(map[string]Node),
		adjMatrix: make(map[string]map[string]int),
	}
}

func (g *Graph) AddNode(name string, x, y int) {
	g.nodes[name] = Node{name, x, y}
}

func (g *Graph) AddEdge(u, v string, capacity int) {
	if g.adjMatrix[u] == nil {
		g.adjMatrix[u] = make(map[string]int)
	}
	if g.adjMatrix[v] == nil {
		g.adjMatrix[v] = make(map[string]int)
	}
	g.adjMatrix[u][v] = capacity
	g.adjMatrix[v][u] = capacity
}

func (g *Graph) BFS(source, sink string, parent map[string]string) bool {
	visited := make(map[string]bool)
	queue := []string{source}
	visited[source] = true
	parent[source] = ""

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for v, capacity := range g.adjMatrix[u] {
			if !visited[v] && capacity > 0 {
				queue = append(queue, v)
				parent[v] = u
				visited[v] = true
				if v == sink {
					return true
				}
			}
		}
	}
	return false
}

func (g *Graph) FordFulkerson(source, sink string) (int, [][]string) {
	parent := make(map[string]string)
	maxFlow := 0
	var allPaths [][]string

	for g.BFS(source, sink, parent) {
		pathFlow := int(^uint(0) >> 1)
		var path []string

		for v := sink; v != source; v = parent[v] {
			u := parent[v]
			pathFlow = min(pathFlow, g.adjMatrix[u][v])
			path = append([]string{v}, path...)
		}
		path = append([]string{source}, path...)

		for v := sink; v != source; v = parent[v] {
			u := parent[v]
			g.adjMatrix[u][v] -= pathFlow
			if g.adjMatrix[v] == nil {
				g.adjMatrix[v] = make(map[string]int)
			}
			g.adjMatrix[v][u] += pathFlow
		}

		maxFlow += pathFlow
		allPaths = append(allPaths, path)
	}
	return maxFlow, allPaths
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SelectPaths(allPaths [][]string, numAnts int) [][]string {
	// Sort paths by length
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	selectedPaths := make([][]string, 0)
	usedNodes := make(map[string]bool)

	for _, path := range allPaths {
		valid := true
		for _, node := range path[1 : len(path)-1] {
			if usedNodes[node] {
				valid = false
				break
			}
		}
		if valid {
			selectedPaths = append(selectedPaths, path)
			for _, node := range path[1 : len(path)-1] {
				usedNodes[node] = true
			}
		}
		if len(selectedPaths) >= numAnts {
			break
		}
	}

	return selectedPaths
}

func SimulateAnts(numAnts int, paths [][]string) []string {
	var lastArr []string

	antPositions := make([]int, numAnts)
	antPaths := make([][]string, numAnts)
	occupied := make(map[string]bool)
	complete := make([]bool, numAnts)

	for i := range antPaths {
		antPaths[i] = make([]string, 0)
	}

	// Distribute ants across paths
	for i := 0; i < numAnts; i++ {
		antPaths[i] = paths[i%len(paths)]
	}

	for {
		moved := false
		stepOutput := []string{}
		for i := 0; i < numAnts; i++ {
			if complete[i] {
				continue
			}
			currentPath := antPaths[i]
			if antPositions[i] < len(currentPath)-1 {
				nextNode := currentPath[antPositions[i]+1]
				if !occupied[nextNode] {
					if antPositions[i] > 0 {
						occupied[currentPath[antPositions[i]]] = false
					}
					antPositions[i]++
					stepOutput = append(stepOutput, fmt.Sprintf("L%d-%s", i+1, nextNode))
					occupied[nextNode] = true
					moved = true
					if antPositions[i] == len(currentPath)-1 {
						complete[i] = true
						occupied[nextNode] = false
					}
				}
			}
		}
		if len(stepOutput) > 0 {
			lastArr = append(lastArr, strings.Join(stepOutput, " "))
		}
		if !moved {
			break
		}
	}
	return lastArr
}

func main() {
	args := os.Args[1:]
	var inputFile string
	var err error
	inputFile, _, err = InputControl(args)
	if err != nil {
		return
	}
	g := NewGraph()

	var startNode, endNode string
	content := ReadAllLines(inputFile)
	lines := strings.Split(content, "\n")
	lines = StrArrCleaner(lines)

	numAnts, _ := strconv.Atoi(lines[0])
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				startNode = lines[i+1]
				i++
			} else if line == "##end" {
				endNode = lines[i+1]
				i++
			}
			continue
		}

		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			g.AddEdge(parts[0], parts[1], 1)
		} else {
			parts := strings.Fields(line)
			name := parts[0]
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			g.AddNode(name, x, y)
		}
	}

	startNodeParts := strings.Fields(startNode)
	endNodeParts := strings.Fields(endNode)
	maxFlow, allPaths := g.FordFulkerson(startNodeParts[0], endNodeParts[0])
	if maxFlow == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	selectedPaths := SelectPaths(allPaths, numAnts)

	newArr := SimulateAnts(numAnts, selectedPaths)
	if len(newArr) == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	} else {
		fmt.Printf("The maximum possible flow is %d\n", maxFlow)
		fmt.Println("Selected paths:")
		for _, path := range selectedPaths {
			fmt.Println(path)
		}
		fmt.Println("Ants' movement:")
		for _, str := range newArr {
			fmt.Println(str)
		}
		fmt.Printf("Ants arrived to end in %d turns.\n", len(newArr))
	}
}

func InputControl(args []string) (string, string, error) {
	var inputFile string
	var outputFile string
	var err error
	outputFile = "exit.txt"
	if len(args) == 1 {
		inputFile = args[0]
	} else if len(args) == 2 {
		inputFile = args[0]
		outputFile = args[1]
	} else {
		fmt.Println("Invalid input", err)
		return "", "", err
	}
	return inputFile, outputFile, nil
}

func StrArrCleaner(Arr []string) []string {
	var deleteArr []int
	for index, item := range Arr {
		if len(item) == 0 || item == " " || item == "\n" {
			deleteArr = append(deleteArr, index)
		}
	}
	Arr = Except(deleteArr, Arr)
	return Arr
}

func Except(lines []int, Arr []string) []string {
	leng := len(lines)
	for i := leng - 1; i >= 0; i-- {
		Arr = RemoveElementStr(Arr, lines[i])
	}
	return Arr
}

func RemoveElementStr(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func ReadAllLines(fileName string) string {
	var file *os.File
	var line string
	var lineArr []string
	var err error
	file, err = os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			if len(line) > 0 {
				lineArr = append(lineArr, line)
			}
			break
		}
		lineArr = append(lineArr, strings.TrimRight(line, "\n"))
	}
	line = strings.Join(lineArr, "\n")
	return line
}
