package specutil

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func ConvertSwaggerToOpenAPI(file string) error {
	cmd := exec.Command("primecodegen",
		"--log-level", "trace",
		"openapi-convert",
		"--format-in", "swagger20",
		"--format-out", "openapi30",
		"--input", file,
	)
	cmd.Stderr = os.Stderr
	log.Trace().Str("cmd", cmd.String()).Msg("calling primecodegen to convert openapi specification")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute primecodegen: %w", err)
	}

	err = os.WriteFile(file, output, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write updated OpenAPI spec: %w", err)
	}

	return nil
}

func MergeAndPatchOpenAPI(files []string, inputPatches []string, patches []string, outputFile string) error {
	args := []string{
		"--log-level", "trace",
		"openapi-patch",
		"-o", outputFile,
	}
	for _, f := range files {
		args = append(args, "-i", f)
	}
	for _, p := range inputPatches {
		args = append(args, "--input-patch", p)
	}
	for _, p := range patches {
		args = append(args, "--patch", p)
	}
	cmd := exec.Command("primecodegen", args...)
	cmd.Stderr = os.Stderr
	log.Trace().Str("cmd", cmd.String()).Msg("calling primecodegen to patch and merge openapi specifications")
	return cmd.Run()
}
