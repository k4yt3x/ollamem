package main

import (
	"flag"
	"fmt"
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
	flag.IntVar(&flags.contextLength, "c", 2048, "context length for model")
	flag.StringVar(&flags.modelName, "m", "", "name of the Ollama model file")
	flag.StringVar(&flags.modelPath, "f", "", "path to the GGUF model file")
	flag.BoolVar(&flags.forceCPU, "cpu", false, "force CPU mode")
	flag.BoolVar(&flags.forceGPU, "gpu", false, "force GPU mode")
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

	// Check if the both CPU and GPU are forced
	if flags.forceCPU && flags.forceGPU {
		fmt.Fprintln(os.Stderr, "Please provide only one of -cpu or -gpu")
		return
	}

	// Discover the available GPUs
	var gpus discover.GpuInfoList
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

	// Estimate the memory required for the model
	projectors := []string{}
	opts := api.DefaultOptions()
	opts.Runner.NumCtx = flags.contextLength
	estimate := llm.EstimateGPULayers(gpus, ggmlFile, projectors, opts)

	// Print the total memory estimate
	fmt.Printf(
		"Estimated required memory: %d bytes (%s)\n",
		estimate.TotalSize,
		format.HumanBytes2(estimate.TotalSize),
	)

	// Dump all estimates if in verbose mode
	if flags.verbose {
		fmt.Println("\nFull Memory Estimate:")
		fmt.Println(estimate.LogValue().String())
		spew.Dump(estimate)
	}
}
