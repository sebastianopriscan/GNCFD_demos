package main

import (
	"fmt"
	"log"

	"github.com/sebastianopriscan/GNCFD/core"
	"github.com/sebastianopriscan/GNCFD/core/impl/vivaldi"
)

var loops int = 0

func analyze_vivaldi_core(core core.GNCFDCore) {
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

	for gd, coor := range viv_meta.Data {
		mssg := ""
		mssg += fmt.Sprintf("Loop %d\n\t GUID: %v\n\tFailed: %v\n\tCoors:\n", loops, gd, coor.IsFailed)

		for _, coor := range coor.Coords {
			mssg += fmt.Sprintf("\t\t%v\n", coor)
		}

		log.Println(mssg)
	}

	loops = loops + 1
}
