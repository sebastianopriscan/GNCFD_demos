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
	0: {1, 2},
	1: {0, 3, 4},
	2: {0, 3},
	3: {1, 2, 5, 7},
	4: {1, 7},
	5: {3, 6},
	6: {5},
	7: {3, 4},
}
