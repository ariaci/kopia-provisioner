package identity

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var providerRegistry = map[string]func(context ProviderPipelineContext, value, next string) (Password, error){
	"env":    newEnvPasswordProvider,
	"file":   newFilePasswordProvider,
	"inline": newInlinePasswordProvider,
}

type ProviderPipelineContext struct {
	BaseDir string
}

type ProviderStage struct {
	Context       ProviderPipelineContext
	NextStageType string
	NextStage     Password
}

type InlinePasswordProvider struct {
	ProviderStage
	Value string
}

type EnvPasswordProvider struct {
	ProviderStage
	Variable string
}

type FilePasswordProvider struct {
	ProviderStage
	Context  ProviderPipelineContext
	FilePath string
}

func newProviderPipelineContext(filePath string) ProviderPipelineContext {
	var err error
	switch {
	case filePath == "":
		filePath, err = os.Executable()
		if err == nil {
			filePath = filepath.Dir(filePath)
		}
	case filepath.IsAbs(filePath):
		filePath = filepath.Dir(filePath)
	default:
		var wd string
		wd, err = os.Getwd()
		if err == nil {
			filePath = filepath.Join(wd, filepath.Dir(filePath))
		}
	}

	if err != nil {
		panic(err)
	}

	return ProviderPipelineContext{
		BaseDir: filePath,
	}
}

func (p *ProviderStage) EnsureNextStage(resolveFn func() (string, error)) (Password, error) {
	if p.NextStage != nil {
		return p.NextStage, nil
	}

	input, err := resolveFn()
	if err != nil {
		return nil, err
	}

	next, err := newPassword(p.Context, p.NextStageType, input)
	if err != nil {
		return nil, err
	}

	p.NextStage = next
	return next, nil
}

func newInlinePasswordProvider(context ProviderPipelineContext, value, nextType string) (Password, error) {
	return InlinePasswordProvider{
		Value: value,
		ProviderStage: ProviderStage{
			Context:       context,
			NextStageType: nextType,
		},
	}, nil
}

func (p InlinePasswordProvider) KopiaArguments() ([]string, error) {
	next, err := p.ProviderStage.EnsureNextStage(func() (string, error) {
		return p.Value, nil
	})
	if err != nil {
		return nil, err
	}

	return next.KopiaArguments()
}

func newEnvPasswordProvider(context ProviderPipelineContext, variable, nextType string) (Password, error) {
	return EnvPasswordProvider{
		Variable: strings.TrimSpace(variable),
		ProviderStage: ProviderStage{
			Context:       context,
			NextStageType: nextType,
		},
	}, nil
}

func (p EnvPasswordProvider) KopiaArguments() ([]string, error) {
	next, err := p.ProviderStage.EnsureNextStage(func() (string, error) {
		if p.Variable == "" {
			return "", fmt.Errorf("environment variable name cannot be empty")
		}
		return os.Getenv(p.Variable), nil
	})
	if err != nil {
		return nil, err
	}

	return next.KopiaArguments()
}

func newFilePasswordProvider(context ProviderPipelineContext, filePath, nextType string) (Password, error) {
	return FilePasswordProvider{
		Context:  context,
		FilePath: strings.TrimSpace(filePath),
		ProviderStage: ProviderStage{
			Context:       context,
			NextStageType: nextType,
		},
	}, nil
}

func (p FilePasswordProvider) resolveFilePath() string {
	if p.FilePath == "" {
		return ""
	}
	if filepath.IsAbs(p.FilePath) {
		return p.FilePath
	}

	return filepath.Join(p.Context.BaseDir, p.FilePath)
}

func (p FilePasswordProvider) KopiaArguments() ([]string, error) {
	next, err := p.ProviderStage.EnsureNextStage(func() (string, error) {
		filePath := p.resolveFilePath()
		if filePath == "" {
			return "", fmt.Errorf("file password provider requires a non-empty file path")
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}

		value := strings.TrimSpace(string(data))
		return value, nil
	})
	if err != nil {
		return nil, err
	}

	return next.KopiaArguments()
}
