package gomoku
import(
    // "log"
    // "math/rand"
)

const (
    GOMOKU_DEFAULT_SIZE = 15
    GOMOKU_ROLE_BLACK   = 2
    GOMOKU_ROLE_WHITE   = 1
    GOMOKU_ROLE_EMPTY   = 0
)

type free_tree struct{
    root []*free_node
}

type free_node struct{
    nexts   []*free_node
    score   int
    point   free_point
}

type free_point struct{
    x       int
    y       int
}

func free_score() {

}



func Gomoku_free(board GomokuBoard, role int, point *GomokuPoint)error {


    return nil
}

