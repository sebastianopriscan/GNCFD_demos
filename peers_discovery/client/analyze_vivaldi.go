package main

import (
	"fmt"
	"log"
	"math"

	"github.com/sebastianopriscan/GNCFD/core"
	"github.com/sebastianopriscan/GNCFD/core/impl/vivaldi"
	"github.com/sebastianopriscan/GNCFD/utils/guid"
)

var loops int = 0

func euclideanNorm(first []float64, second []float64) float64 {
	sum := 0.
	for i := 0; i < len(first); i++ {
		sum += math.Pow(first[i]-second[i], 2.)
	}

	return math.Sqrt(sum)
}

func analyze_vivaldi_core(core core.GNCFDCoreInteractionGate) {

	/*
		metadata, err := core.GetStateUpdates()
		if err != nil {
			log.Printf("error retrieving core updates, details: %s", err)
			return
		}

		viv_meta, ok := metadata.(*vivaldi.VivaldiMetadata[float64])
		if !ok {
			log.Println("core not compatible")
			return
		}
	*/

	viv_core, ok := core.(*vivaldi.VivaldiCore[float64])
	if !ok {
		log.Println("core not compatible")
		return
	}
	viv_meta, err := viv_core.DumpCore()
	if err != nil {
		log.Printf("error dumping core, details: %v\n", viv_meta)
		return
	}

	var myGuid guid.Guid
	var myCoors []float64
	for gd, coor := range viv_meta.Data {
		if gd == viv_meta.Communicator {
			myGuid = gd
			myCoors = coor.Coords
			break
		}
	}

	message := ""
	for gd, coor := range viv_meta.Data {
		mssg := ""
		mssg += fmt.Sprintf("Loop %d\n\t GUID: %v\n\tFailed: %v\n\tCoors:\n", loops, gd, coor.IsFailed)

		for _, coor := range coor.Coords {
			mssg += fmt.Sprintf("\t\t%v\n", coor)
		}

		message += mssg
	}

	message += fmt.Sprintf("\nDistance from me, %v:\n", myGuid)
	for gd, coor := range viv_meta.Data {
		if gd != myGuid {
			dist := euclideanNorm(myCoors, coor.Coords)
			message += fmt.Sprintf("\t%v:%v\n", gd, dist)
		}
	}

	log.Println(message)
	loops = loops + 1
}
