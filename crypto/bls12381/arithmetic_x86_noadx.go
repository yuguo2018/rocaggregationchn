//go:build amd64 && blsasm
// +build amd64,blsasm

package bls12381

// enableADX is true if the ADX/BMI2 instruction set was requested for the BLS
// implementation. The system may still fall back to plain ASM if the necessary
// instructions are unavailable on the CPU.
const enableADX = false
