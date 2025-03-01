# Ollamem

Accurately estimate the memory required to run GGUF models using Ollama's original memory estimation functions.

Since this project is using Ollama's original `llm.EstimateGPULayers` function, it should get the exact estimation results as Ollama.

This project should probably be a part of Ollama. This tool is inefficient standalone as it requires much of Ollama's code to function. The built binary is around 20 MiB in size. But hey, it works.

![Image](https://github.com/user-attachments/assets/8b226ed2-cbdf-4dd1-9548-bd8a7a4a04d7)

## Build

Install the following dependencies:

- [Go 1.24+](https://golang.org/dl/)
- [Git](https://git-scm.com/downloads)

Use the following commands to build Ollamem:

```bash
git clone https://github.com/k4yt3x/ollamem.git
cd ollamem
go build -ldflags="-s -w" -trimpath -o bin/ollamem ./cmd/ollamem
```

## Usage

To estimate the memory required to run an installed Ollama model, use the following command:

```bash
$ ollamem -m <model_name> -c <context_length>

# Example: Estimate memory required to run the L3.3-MS-Nevoria-70b model with a context length of 2048
$ ollamem -m hf.co/bartowski/L3.3-MS-Nevoria-70b-GGUF:Q4_K_M
Estimated required memory: 48312662016 bytes (45.0 GiB)
Estimated maximum context length with free memory: 11402

# Example: Estimate memory required to run the L3.3-MS-Nevoria-70b model with a context length of 40960
$ ollamem -m hf.co/bartowski/L3.3-MS-Nevoria-70b-GGUF:Q4_K_M -c 40960
Estimated required memory: 48312662016 bytes (45.0 GiB)
Estimated maximum context length with free memory: 11402
```

You can also specify the path to a GGUF model file:

```bash
$ ollamem -f <model_file>

# Example
$ ollamem -f ./example.gguf
Estimated required memory: 47814449152 bytes (44.5 GiB)
Estimated maximum context length with free memory: 16746
```

Use the `-v` flag to dump all fields of the memory estimate:

```bash
$ ollamem -m hf.co/bartowski/L3.3-MS-Nevoria-70b-GGUF:Q4_K_M -v
Estimated required memory: 48312662016 bytes (45.0 GiB)
Estimated maximum context length with free memory: 11402

Full Memory Estimate:
[library=cuda layers=[requested=-1 model=81 offload=81 split=] memory=[available=[45.0 GiB] gpu_overhead=0 B required=[full=45.0 GiB partial=45.0 GiB kv=3.5 GiB allocations=[45.0 GiB]] weights=[total=41.7 GiB repeating=40.9 GiB nonrepeating=822.0 MiB] graph=[full=1.5 GiB partial=1.5 GiB]]]
(llm.MemoryEstimate) {
 Layers: (int) 81,
 Graph: (uint64) 1584945152,
 VRAMSize: (uint64) 48312662016,
 TotalSize: (uint64) 48312662016,
 TensorSplit: (string) "",
 GPUSizes: ([]uint64) (len=1 cap=1) {
  (uint64) 48312662016
 },
 inferenceLibrary: (string) (len=4) "cuda",
 layersRequested: (int) -1,
 layersModel: (int) 81,
 availableList: ([]string) (len=1 cap=1) {
  (string) (len=8) "45.0 GiB"
 },
 kv: (uint64) 3736207360,
 allocationsList: ([]string) (len=1 cap=1) {
  (string) (len=8) "45.0 GiB"
 },
 memoryWeights: (uint64) 44795822080,
 memoryLayerOutput: (uint64) 861913088,
 graphFullOffload: (uint64) 1584945152,
 graphPartialOffload: (uint64) 1635842048,
 projectorWeights: (uint64) 0,
 projectorGraph: (uint64) 0
}
```
