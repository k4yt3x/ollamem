# Ollamem

Accurately estimate the memory required to run GGUF models using Ollama's original memory estimation functions.

Since this project is using Ollama's original `llm.EstimateGPULayers` function, it should get the exact estimation results as Ollama.

This project should probably be a part of Ollama. This tool is inefficient standalone as it requires much of Ollama's code to function. The built binary is around 20 MiB in size. But hey, it works.

![Image](https://github.com/user-attachments/assets/052fca4d-b43b-444f-aa8b-589cac0aeb8f)

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
Estimated required memory: 43964024832 bytes (40.9 GiB)

# Example: Estimate memory required to run the L3.3-MS-Nevoria-70b model with a context length of 40960
$ ollamem -m hf.co/bartowski/L3.3-MS-Nevoria-70b-GGUF:Q4_K_M -c 40960
Estimated required memory: 62596040704 bytes (58.3 GiB)
```

You can also specify the path to a GGUF model file:

```bash
$ ollamem -f <model_file>

# Example
$ ollamem -f ~/.ollama/models/blobs/sha256-65bd93d5deec246dab52b752ba7f37182a1aedb69bd1be364979abe087905363
Estimated required memory: 41003425792 bytes (38.2 GiB)
```

Use the `-v` flag to dump all fields of the memory estimate:

```bash
$ ollamem -m hf.co/bartowski/L3.3-MS-Nevoria-70b-GGUF:Q4_K_M -v
2025/03/01 04:31:39 INFO looking for compatible GPUs
Estimated required memory: 43964024832 bytes (40.9 GiB)

Full Memory Estimate:
[library=cuda layers=[requested=-1 model=81 offload=81 split=] memory=[available=[44.9 GiB] gpu_overhead=0 B required=[full=40.9 GiB partial=40.9 GiB kv=640.0 MiB allocations=[40.9 GiB]] weights=[total=38.9 GiB repeating=38.1 GiB nonrepeating=822.0 MiB] graph=[full=324.0 MiB partial=1.1 GiB]]]
(llm.MemoryEstimate) {
 Layers: (int) 81,
 Graph: (uint64) 339740672,
 VRAMSize: (uint64) 43964024832,
 TotalSize: (uint64) 43964024832,
 TensorSplit: (string) "",
 GPUSizes: ([]uint64) (len=1 cap=1) {
  (uint64) 43964024832
 },
 inferenceLibrary: (string) (len=4) "cuda",
 layersRequested: (int) -1,
 layersModel: (int) 81,
 availableList: ([]string) (len=1 cap=1) {
  (string) (len=8) "44.9 GiB"
 },
 kv: (uint64) 671088640,
 allocationsList: ([]string) (len=1 cap=1) {
  (string) (len=8) "40.9 GiB"
 },
 memoryWeights: (uint64) 41730703360,
 memoryLayerOutput: (uint64) 861913088,
 graphFullOffload: (uint64) 339740672,
 graphPartialOffload: (uint64) 1158103040,
 projectorWeights: (uint64) 0,
 projectorGraph: (uint64) 0
}
```
