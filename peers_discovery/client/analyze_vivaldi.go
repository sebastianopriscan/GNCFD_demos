package main

import (
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
		log.Printf("Loop %d\n\t GUID: %v\n\tCoors:\n", loops, gd)

		for _, coor := range coor.Coords {
			log.Printf("\t\t%v\n", coor)
		}
	}

	loops = loops + 1
}
