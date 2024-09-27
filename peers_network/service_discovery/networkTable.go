package servicediscovery

import (
	"sync"

	"github.com/sebastianopriscan/GNCFD/utils/guid"
)

var Mu sync.Mutex

var GuidNode map[guid.Guid]int = make(map[guid.Guid]int)
var NodeGuid map[int]guid.Guid = make(map[int]guid.Guid)

var LastGiven int = -1

var NetworkTable map[int][]int = map[int][]int{
	0:  {1, 2},
	1:  {0, 3, 4},
	2:  {0, 3},
	3:  {5, 6, 7, 8},
	4:  {1, 8},
	5:  {3, 9, 10},
	6:  {3, 11},
	7:  {3},
	8:  {3, 4},
	9:  {10, 5},
	10: {9, 5},
	11: {6},
}
