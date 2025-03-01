# Ollamem

Accurately estimate the memory required to run GGUF models using Ollama's original memory estimation functions.

Since this project is using Ollama's original `llm.EstimateGPULayers` function, it should get the exact estimation results as Ollama.

This project should probably be a part of Ollama. This tool is inefficient standalone as it requires much of Ollama's code to function. The built binary is around 20 MiB in size. But hey, it works.

![Image](https://github.com/user-attachments/assets/00a99945-4ffe-437f-b4b0-f6bc5d817158)

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

# Example: Estimate memory required to run the DeepSeek-R1-Distill-Llama-70B-Uncensored-i1 model with a context length of 2048
$ ollamem -m hf.co/mradermacher/DeepSeek-R1-Distill-Llama-70B-Uncensored-i1-GGUF:Q4_K_M
Estimated required memory: 45059398208 bytes
Estimated required memory: 44003318.562500 KiB
Estimated required memory: 42971.990784 MiB
Estimated required memory: 41.964835 GiB

# Example: Estimate memory required to run the DeepSeek-R1-Distill-Llama-70B-Uncensored-i1 model with a context length of 40960
$ ollamem -m hf.co/mradermacher/DeepSeek-R1-Distill-Llama-70B-Uncensored-i1-GGUF:Q4_K_M -c 40960
Estimated required memory: 57810082368 bytes
Estimated required memory: 56455158.562500 KiB
Estimated required memory: 55131.990784 MiB
Estimated required memory: 53.839835 GiB
```

You can also specify the path to a GGUF model file:

```bash
$ ollamem -f <model_file>

# Example
$ ollamem -f ~/.ollama/models/blobs/sha256-65bd93d5deec246dab52b752ba7f37182a1aedb69bd1be364979abe087905363
Estimated required memory: 39712391168 bytes
Estimated required memory: 38781632.000000 KiB
Estimated required memory: 37872.687500 MiB
Estimated required memory: 36.985046 GiB
```

Use the `-v` flag to dump all fields of the memory estimate:

```bash
$ ollamem -m hf.co/mradermacher/DeepSeek-R1-Distill-Llama-70B-Uncensored-i1-GGUF:Q4_K_M -v
Estimated required memory: 45059398208 bytes
Estimated required memory: 44003318.562500 KiB
Estimated required memory: 42971.990784 MiB
Estimated required memory: 41.964835 GiB

Full Memory Estimate:
(llm.MemoryEstimate) {
 Layers: (int) 0,
 Graph: (uint64) 0,
 VRAMSize: (uint64) 0,
 TotalSize: (uint64) 45059398208,
 TensorSplit: (string) "",
 GPUSizes: ([]uint64) {
 },
 inferenceLibrary: (string) (len=3) "cpu",
 layersRequested: (int) -1,
 layersModel: (int) 81,
 availableList: ([]string) (len=1 cap=1) {
  (string) (len=3) "0 B"
 },
 kv: (uint64) 671088640,
 allocationsList: ([]string) (len=1 cap=1) {
  (string) (len=3) "0 B"
 },
 memoryWeights: (uint64) 41730703360,
 memoryLayerOutput: (uint64) 861919808,
 graphFullOffload: (uint64) 339740672,
 graphPartialOffload: (uint64) 1158111808,
 projectorWeights: (uint64) 0,
 projectorGraph: (uint64) 0
}
```
