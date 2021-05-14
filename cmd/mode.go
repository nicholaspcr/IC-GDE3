package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	mo "github.com/nicholaspcr/gde3/pkg/mode"
	"github.com/nicholaspcr/gde3/pkg/problems"
	"github.com/nicholaspcr/gde3/pkg/problems/models"
	"github.com/nicholaspcr/gde3/pkg/variants"
	"github.com/spf13/cobra"
)

// local flags
var mConst int
var functionName string
var variantName string
var disablePlot bool

// modeCmd represents the mode command
var modeCmd = &cobra.Command{
	Use:   "multi",
	Short: "Multi-objective implementation of DE",
	Long:  `An implementation that allows the processing of multiple objective functions, these are a bit more complex and time consuming overall.`,

	Run: func(cmd *cobra.Command, args []string) {
		problem := problems.GetProblemByName(functionName)
		variant := variants.GetVariantByName(variantName)
		if problem.Name == "" {
			fmt.Println("Invalid problem")
			return
		}
		if variant.Name == "" {
			fmt.Println("Invalid variant.")
			return
		}
		params := models.Params{
			NP:    np,
			M:     mConst,
			DIM:   dim,
			GEN:   gen,
			EXECS: execs,
			FLOOR: floor,
			CEIL:  ceil,
			CR:    crConst,
			F:     fConst,
			P:     pConst,
		}
		if cpuprofile != "" {
			cpuF, err := os.Create(cpuprofile)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			defer cpuF.Close() // error handling omitted for example
			if err := pprof.StartCPUProfile(cpuF); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()

		}

		mo.MultiExecutions(params, problem, variant, disablePlot)

		if memprofile != "" {
			memF, err := os.Create(memprofile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			defer memF.Close() // error handling omitted for example
			runtime.GC()       // get up-to-date statistics
			if err := pprof.WriteHeapProfile(memF); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(modeCmd)
	modeCmd.Flags().IntVar(&mConst,
		"M",
		3,
		"M -> DE constant")

	modeCmd.Flags().StringVar(&functionName,
		"fn",
		"DTLZ1",
		"name of the problem to be used.")

	modeCmd.Flags().StringVar(&variantName,
		"vr",
		"rand1",
		"name fo the variant to be used")

	modeCmd.Flags().BoolVar(&disablePlot,
		"disable-plot",
		false,
		"to write in files the result of the gde3 to be able to plot it with the python scripts")

}
