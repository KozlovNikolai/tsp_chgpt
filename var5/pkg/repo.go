package pkg

type Repo struct {
	Store         map[int]*Node
	CurrentNodeID int
	NextID        int
}

type Node struct {
	ID       int
	Level    int
	Path     []int
	Cost     int
	ParentID int
	Matrix   [][]int
}

func NewRepo(node *Node) *Repo {
	return &Repo{
		CurrentNodeID: node.ID,
		NextID:        1,
		Store: map[int]*Node{
			0: node,
		},
	}
}

func (r *Repo) CreateLeaves() {
	mx := CloneMx(r.Store[r.CurrentNodeID].Matrix)
	nextID := r.NextID
	for j := 1; j < len(mx[0]); j++ {
		if mx[0][j] != INF {
			reducedMatrix := RemoveCellFromMatrixByIdx(mx, 0, j)

			path := make([]int, len(r.Store[r.CurrentNodeID].Path))
			copy(path, r.Store[r.CurrentNodeID].Path)
			path = append(path, nextID)
			cost := CalculateCost(reducedMatrix)
			r.CreateLeaf(reducedMatrix, path, nextID, cost)
			nextID++
		}
	}
	r.NextID = nextID
}

func (r *Repo) CreateLeaf(mx [][]int, path []int, nextID, cost int) {
	r.Store[nextID] = &Node{
		ID:       nextID,
		Level:    r.Store[r.CurrentNodeID].Level + 1,
		Path:     path,
		Cost:     cost,
		ParentID: r.CurrentNodeID,
		Matrix:   mx,
	}
}
