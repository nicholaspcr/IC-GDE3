package mo

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
)

func generatePopulation(p Params) []Elem {
	ret := make([]Elem, p.NP)
	constant := p.CEIL - p.CEIL // range between floor and ceiling
	for i := 0; i < p.NP; i++ {
		ret[i].X = make([]float64, p.DIM)

		for j := 0; j < p.DIM; j++ {
			ret[i].X[j] = rand.Float64()*constant + p.CEIL // value varies within [ceil,upper]
		}

		// for ZDT4
		// ret[i].X[0] = rand.Float64()
	}
	return ret
}

// generates random indices in the int slice, r -> it's a pointer
func generateIndices(startInd, NP int, r []int) error {
	if len(r) > NP {
		return errors.New("insufficient elements in population to generate random indices")
	}
	for i := startInd; i < len(r); i++ {
		for done := false; !done; {
			r[i] = rand.Int() % NP
			done = true
			for j := 0; j < i; j++ {
				done = done && r[j] != r[i]
			}
		}
	}
	return nil
}

func checkFilePath(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			log.Fatalf("error creating file in path: %v", filePath)
		}
	}
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func writeHeader(pop []Elem, f *os.File) {
	for i := range pop {
		fmt.Fprintf(f, "pop[%d]\t", i)
	}
	fmt.Fprintf(f, "\n")
}

func writeGeneration(pop []Elem, f *os.File) {
	qtdObjs := len(pop[0].objs)
	for i := 0; i < qtdObjs; i++ {
		for _, p := range pop {
			fmt.Fprintf(f, "%10.3f\t", p.objs[i])
		}
		fmt.Fprintf(f, "\n")
	}
}
