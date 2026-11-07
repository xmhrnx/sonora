package onnx

import (
	"errors"
	"path/filepath"

	"github.com/mehran9675/sonora/internal/engines/onnx/sherpa"
)

// recognizerConfig finds a package's files.
func recognizerConfig(dir string, m model) (sherpa.RecognizerConfig, error) {
	cfg := sherpa.RecognizerConfig{Family: m.Family}
	find := func(dst *string, patterns ...string) error {
		for _, p := range patterns {
			if found, _ := filepath.Glob(filepath.Join(dir, p)); len(found) > 0 {
				*dst = found[0]
				return nil
			}
		}
		return errors.New("the downloaded model is incomplete; remove the engine and install it again")
	}
	var errs []error
	switch m.Family {
	case sherpa.Whisper:
		errs = append(errs, find(&cfg.Encoder, "*-encoder.int8.onnx"), find(&cfg.Decoder, "*-decoder.int8.onnx"),
			find(&cfg.Tokens, "*-tokens.txt"))
	case sherpa.Transducer:
		errs = append(errs, find(&cfg.Encoder, "encoder*.onnx"), find(&cfg.Decoder, "decoder*.onnx"),
			find(&cfg.Joiner, "joiner*.onnx"), find(&cfg.Tokens, "tokens.txt"))
	case sherpa.SenseVoice:
		errs = append(errs, find(&cfg.Model, "model*.onnx"), find(&cfg.Tokens, "tokens.txt"))
	case sherpa.Moonshine:
		errs = append(errs, find(&cfg.Preprocessor, "preprocess*.onnx"), find(&cfg.Encoder, "encode*.onnx"),
			find(&cfg.Decoder, "uncached_decode*.onnx"), find(&cfg.CachedDecoder, "cached_decode*.onnx"),
			find(&cfg.Tokens, "tokens.txt"))
	}
	return cfg, errors.Join(errs...)
}
