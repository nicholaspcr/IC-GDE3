package mo

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// MultiExecutions returns the pareto front of the total of 30 executions of the same problem
func MultiExecutions(params Params, prob ProblemFn, variant VariantFn) {
	basePath := os.Getenv("HOME")
	paretoPath := ".go-de/mode/paretoFront"
	checkFilePath(basePath, paretoPath)

	startTimer := time.Now()                 //	timer start
	rand.Seed(time.Now().UTC().UnixNano())   // Rand Seed
	population := generatePopulation(params) // random generated population
	var wg sync.WaitGroup                    // number of working go routines
	elemChan := make(chan Elements)
	for i := 0; i < params.EXECS; i++ {
		f, err := os.Create(basePath +
			"/" +
			paretoPath +
			"/exec-" +
			strconv.Itoa(i+1) +
			".csv")
		checkError(err)
		wg.Add(1)
		// worker
		go func() {
			defer wg.Done()
			elemChan <- DE(params, prob, variant, population.Copy(), f)
		}()
	}
	// closer
	go func() {
		wg.Wait()
		close(elemChan)
	}()

	var pareto Elements // DE pareto front
	for i := 0; i < params.EXECS; i++ {
		v, ok := <-elemChan
		if !ok {
			fmt.Println("one of the goroutine workers didn't work")
		}
		pareto = append(pareto, v...)
	}
	result := filterDominated(pareto) // non dominated set

	multiExecPath := ".go-de/mode/multiExecutions"
	checkFilePath(basePath, multiExecPath)
	// todo: add the use of the variant name here
	f, err := os.Create(basePath + "/" + multiExecPath + "/" + variant.Name + ".csv")
	checkError(err)
	defer f.Close()
	writeHeader(result, f)
	writeGeneration(result, f)
	fmt.Println("Done writing file!")

	timeSpent := time.Since(startTimer)
	fmt.Println(timeSpent)
}

// DE -> runs a simple multiObjective DE in the ZDT1 case
func DE(
	p Params,
	evaluate ProblemFn,
	variant VariantFn,
	population Elements,
	f *os.File,
) Elements {
	defer f.Close()

	for i := range population {
		err := evaluate(&population[i], p.M)
		checkError(err)
	}
	writeHeader(population, f)
	writeGeneration(population, f)

	for ; p.GEN > 0; p.GEN-- {
		trial := population.Copy() // trial population slice
		for i, t := range trial {
			v, err := variant.fn(population, p)
			checkError(err)
			// CROSS OVER
			currInd := rand.Int() % p.DIM
			for j := 0; j < p.DIM; j++ {
				changeProb := rand.Float64()
				if changeProb < p.CR || currInd == p.DIM-1 {
					t.X[currInd] = v.X[currInd]
				}
				if t.X[currInd] < p.FLOOR {
					t.X[currInd] = p.FLOOR
				}
				if t.X[currInd] > p.CEIL {
					t.X[currInd] = p.CEIL
				}
				currInd = (currInd + 1) % p.DIM
			}
			evalErr := evaluate(&t, p.M)
			checkError(evalErr)
			if t.dominates(population[i]) {
				population[i] = t.Copy()
			} else if !population[i].dominates(t) {
				population = append(population, t.Copy())
			}
		}

		population = reduceByCrowdDistance(population, p.NP)
		writeGeneration(population, f)
	}
	return population
}
