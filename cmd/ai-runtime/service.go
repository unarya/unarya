package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"sync"

	"github.com/unarya/unarya/lib/proto/pb/aipb"
	onnx "github.com/yalue/onnxruntime_go"
)

// RuntimeServer implements the AI gRPC inference interface
type RuntimeServer struct {
	aipb.UnimplementedAIInferenceServer

	mu      sync.RWMutex
	session *onnx.Session[float32]
	model   string

	// Store input/output info and tensors
	inputNames   []string
	outputNames  []string
	inputShape   []int64
	outputShape  []int64
	inputTensor  *onnx.Tensor[float32]
	outputTensor *onnx.Tensor[float32]
}

// LoadModel loads an ONNX model into GPU-backed session
func (r *RuntimeServer) LoadModel(modelPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Close old session if exists
	if r.session != nil {
		if err := r.session.Destroy(); err != nil {
			fmt.Printf("[Runtime] Warning: failed to destroy old session: %v\n", err)
		}
		r.session = nil
	}

	// Destroy old tensors if exist
	if r.inputTensor != nil {
		r.inputTensor.Destroy()
		r.inputTensor = nil
	}
	if r.outputTensor != nil {
		r.outputTensor.Destroy()
		r.outputTensor = nil
	}

	// Initialize ONNX Runtime environment
	if err := onnx.InitializeEnvironment(); err != nil {
		return fmt.Errorf("failed to initialize ONNX environment: %w", err)
	}

	// Create CUDA provider options
	cudaOptions, err := onnx.NewCUDAProviderOptions()
	if err != nil {
		return fmt.Errorf("failed to create CUDA options: %w", err)
	}
	cudaOptions.Update(map[string]string{
		"device_id": "0",
	})

	// Create session options
	sessionOptions, err := onnx.NewSessionOptions()
	if err != nil {
		return fmt.Errorf("failed to create session options: %w", err)
	}
	defer sessionOptions.Destroy()

	// Enable CUDA execution provider
	if err := sessionOptions.AppendExecutionProviderCUDA(cudaOptions); err != nil {
		return fmt.Errorf("failed to enable CUDA: %w", err)
	}

	// Set default shapes and names (adjust for your model)
	r.inputNames = []string{"input"}
	r.outputNames = []string{"output"}
	r.inputShape = []int64{1, 3, 224, 224}
	r.outputShape = []int64{1, 1000}

	// Create input and output tensors (these will be reused)
	inputTensor, err := onnx.NewEmptyTensor[float32](onnx.NewShape(r.inputShape...))
	if err != nil {
		return fmt.Errorf("failed to create input tensor: %w", err)
	}

	outputTensor, err := onnx.NewEmptyTensor[float32](onnx.NewShape(r.outputShape...))
	if err != nil {
		inputTensor.Destroy()
		return fmt.Errorf("failed to create output tensor: %w", err)
	}

	// Create session with pre-configured inputs/outputs
	session, err := onnx.NewSession[float32](
		modelPath,
		r.inputNames,
		r.outputNames,
		[]*onnx.Tensor[float32]{inputTensor},
		[]*onnx.Tensor[float32]{outputTensor},
	)
	if err != nil {
		inputTensor.Destroy()
		outputTensor.Destroy()
		return fmt.Errorf("failed to load model %s: %w", modelPath, err)
	}

	r.session = session
	r.inputTensor = inputTensor
	r.outputTensor = outputTensor
	r.model = modelPath

	fmt.Printf("[Runtime] ✅ Model loaded: %s\n", modelPath)
	fmt.Printf("[Runtime] Inputs: %v\n", r.inputNames)
	fmt.Printf("[Runtime] Outputs: %v\n", r.outputNames)
	return nil
}

// Predict handles single inference
func (r *RuntimeServer) Predict(ctx context.Context, req *aipb.PredictRequest) (*aipb.PredictResponse, error) {
	// Read session under read-lock
	r.mu.RLock()
	sess := r.session
	inputTensor := r.inputTensor
	outputTensor := r.outputTensor
	inputShape := append([]int64(nil), r.inputShape...)
	r.mu.RUnlock()

	if sess == nil || inputTensor == nil || outputTensor == nil {
		return nil, fmt.Errorf("model not loaded")
	}

	// Convert bytes → []float32
	inputData := bytesToFloat32(req.Input)
	if inputData == nil {
		return nil, fmt.Errorf("invalid input bytes")
	}

	// Validate input size
	expectedSize := int64(1)
	for _, dim := range inputShape {
		expectedSize *= dim
	}
	if len(inputData) != int(expectedSize) {
		return nil, fmt.Errorf("input size mismatch: got %d, expected %d", len(inputData), expectedSize)
	}

	// Copy input data into the tensor's underlying data slice
	inputSlice := inputTensor.GetData()
	copy(inputSlice, inputData)

	// Run inference
	if err := sess.Run(); err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	// Get output data from tensor
	outputData := outputTensor.GetData()

	return &aipb.PredictResponse{
		Output: float32ToBytes(outputData),
	}, nil
}

// PredictBatch handles batch inference by calling Predict per item.
// (Simple implementation; consider batching at tensor level for performance)
func (r *RuntimeServer) PredictBatch(ctx context.Context, req *aipb.PredictBatchRequest) (*aipb.PredictBatchResponse, error) {
	var outs [][]byte
	for _, in := range req.Inputs {
		resp, err := r.Predict(ctx, &aipb.PredictRequest{Input: in, Model: req.Model})
		if err != nil {
			return nil, err
		}
		outs = append(outs, resp.Output)
	}
	return &aipb.PredictBatchResponse{Outputs: outs}, nil
}

// ReloadModel reloads the model file on demand
func (r *RuntimeServer) ReloadModel(ctx context.Context, req *aipb.ReloadModelRequest) (*aipb.ReloadModelResponse, error) {
	if err := r.LoadModel(req.ModelPath); err != nil {
		return &aipb.ReloadModelResponse{Ok: false, Message: err.Error()}, nil
	}
	return &aipb.ReloadModelResponse{Ok: true, Message: "ok"}, nil
}

// Status returns current runtime status
func (r *RuntimeServer) Status(ctx context.Context, req *aipb.StatusRequest) (*aipb.StatusResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	state := "no_model"
	if r.session != nil {
		state = "loaded"
	}
	return &aipb.StatusResponse{Model: r.model, Version: "gpu-runtime", State: state}, nil
}

// Close releases ONNX resources
func (r *RuntimeServer) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.session != nil {
		if err := r.session.Destroy(); err != nil {
			return err
		}
		r.session = nil
	}

	if r.inputTensor != nil {
		r.inputTensor.Destroy()
		r.inputTensor = nil
	}

	if r.outputTensor != nil {
		r.outputTensor.Destroy()
		r.outputTensor = nil
	}

	return nil
}

// Helper functions for byte/float32 conversion
func bytesToFloat32(b []byte) []float32 {
	if len(b)%4 != 0 {
		return nil
	}
	result := make([]float32, len(b)/4)
	for i := 0; i < len(result); i++ {
		bits := binary.LittleEndian.Uint32(b[i*4 : (i+1)*4])
		result[i] = math.Float32frombits(bits)
	}
	return result
}

func float32ToBytes(f []float32) []byte {
	result := make([]byte, len(f)*4)
	for i, v := range f {
		bits := math.Float32bits(v)
		binary.LittleEndian.PutUint32(result[i*4:(i+1)*4], bits)
	}
	return result
}
