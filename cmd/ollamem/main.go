package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/discover"
	"github.com/ollama/ollama/llm"
	"github.com/ollama/ollama/server"
)

type Flags struct {
	contextLength int
	modelName     string
	modelPath     string
	verbose       bool
}

//nolint:forbidigo // fmt is used for printing
func main() {
	flags := Flags{}
	flag.IntVar(&flags.contextLength, "c", 2048, "context length for model")
	flag.StringVar(&flags.modelName, "m", "", "name of the Ollama model file")
	flag.StringVar(&flags.modelPath, "f", "", "path to the GGUF model file")
	flag.BoolVar(&flags.verbose, "v", false, "display all estimates")
	flag.Parse()

	var modelPath string
	switch {
	case flags.modelName != "":
		model, err := server.GetModel(flags.modelName)
		if err != nil {
			panic(err)
		}
		modelPath = model.ModelPath
	case flags.modelPath != "":
		modelPath = flags.modelPath
	default:
		fmt.Fprintln(os.Stderr, "Please provide either a Ollama model name or GGUF model path")
		return
	}

	// Load the model
	ggmlFile, err := llm.LoadModel(modelPath, 0)
	if err != nil {
		panic(err)
	}

	// Estimate the memory required for the model
	gpus := []discover.GpuInfo{
		{
			Library: "cpu",
		},
	}
	projectors := []string{}
	opts := api.DefaultOptions()
	opts.Runner.NumCtx = flags.contextLength
	estimate := llm.EstimateGPULayers(gpus, ggmlFile, projectors, opts)

	// Print the total memory estimate
	fmt.Printf("Estimated required memory: %d bytes\n", estimate.TotalSize)
	fmt.Printf("Estimated required memory: %f KiB\n", float64(estimate.TotalSize)/1024)
	fmt.Printf("Estimated required memory: %f MiB\n", float64(estimate.TotalSize)/math.Pow(1024, 2))
	fmt.Printf("Estimated required memory: %f GiB\n", float64(estimate.TotalSize)/math.Pow(1024, 3))

	// Dump all estimates if in verbose mode
	if flags.verbose {
		fmt.Println("\nFull Memory Estimate:")
		spew.Dump(estimate)
	}
}
