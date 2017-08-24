package gomoku

type GomokuPoint struct{
    X int `json:"x"`
    Y int `json:"y"`
}

type GomokuBoard [][]int
