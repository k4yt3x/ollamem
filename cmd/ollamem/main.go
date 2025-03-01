package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/discover"
	"github.com/ollama/ollama/format"
	"github.com/ollama/ollama/llm"
	"github.com/ollama/ollama/server"
)

type Flags struct {
	contextLength int
	modelName     string
	modelPath     string
	forceCPU      bool
	forceGPU      bool
	verbose       bool
}

//nolint:forbidigo // fmt is used for printing
func main() {
	flags := Flags{}
	flag.IntVar(&flags.contextLength, "c", 2048, "context length")
	flag.StringVar(&flags.modelName, "m", "", "name of the Ollama model file")
	flag.StringVar(&flags.modelPath, "f", "", "path to the GGUF model file")
	flag.BoolVar(&flags.forceCPU, "cpu", false, "force CPU mode")
	flag.BoolVar(&flags.forceGPU, "gpu", false, "force GPU mode")
	flag.BoolVar(&flags.verbose, "v", false, "display all estimates")
	flag.Parse()

	// Disable logger output
	log.SetOutput(io.Discard)

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

	// Check if the both CPU and GPU are forced
	if flags.forceCPU && flags.forceGPU {
		fmt.Fprintln(os.Stderr, "Please provide only one of -cpu or -gpu")
		return
	}

	// Estimate the memory required for the model
	// Here the GPUs needs to be set to the CPU to get an accurate estimate
	gpus := discover.GetCPUInfo()
	projectors := []string{}
	opts := api.DefaultOptions()
	opts.Runner.NumCtx = flags.contextLength
	estimate := llm.EstimateGPULayers(gpus, ggmlFile, projectors, opts)
	totalSize := estimate.TotalSize

	// Discover the available GPUs
	switch {
	case flags.forceCPU:
		gpus = discover.GetCPUInfo()
	case flags.forceGPU:
		gpus = discover.GetGPUInfo()
	default:
		gpus = discover.GetGPUInfo()
		if len(gpus) == 0 {
			gpus = discover.GetCPUInfo()
		}
	}

	// Estimate the maximum context length
	maxCtxLength := 0
	left, right := 0, math.MaxInt32

	for left <= right {
		mid := left + (right-left)/2
		opts.Runner.NumCtx = mid
		estimate = llm.EstimateGPULayers(gpus, ggmlFile, projectors, opts)

		if estimate.TotalSize < gpus[0].FreeMemory {
			left = mid + 1
			maxCtxLength = mid
		} else {
			right = mid - 1
		}
	}

	// Print the total memory estimates
	fmt.Printf(
		"Estimated required memory: %d bytes (%s)\n",
		totalSize,
		format.HumanBytes2(totalSize),
	)
	fmt.Printf("Estimated maximum context length with free memory: %d\n", maxCtxLength)

	// Dump all estimates if in verbose mode
	if flags.verbose {
		fmt.Println("\nFull Memory Estimate:")
		fmt.Println(estimate.LogValue().String())
		spew.Dump(estimate)
	}
}
