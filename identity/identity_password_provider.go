package identity

import "os"

var providerRegistry = map[string]func(value, next string) (Password, error){
	"env":    newEnvPasswordProvider,
	"inline": newInlinePasswordProvider,
}

type ProviderStage struct {
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

func (p *ProviderStage) EnsureNextStage(resolveFn func() (string, error)) (Password, error) {
	if p.NextStage != nil {
		return p.NextStage, nil
	}

	input, err := resolveFn()
	if err != nil {
		return nil, err
	}

	next, err := newPassword(p.NextStageType, input)
	if err != nil {
		return nil, err
	}

	p.NextStage = next
	return next, nil
}

func newInlinePasswordProvider(value, nextType string) (Password, error) {
	return InlinePasswordProvider{
		Value: value,
		ProviderStage: ProviderStage{
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

func newEnvPasswordProvider(variable, nextType string) (Password, error) {
	return EnvPasswordProvider{
		Variable: variable,
		ProviderStage: ProviderStage{
			NextStageType: nextType,
		},
	}, nil
}

func (p EnvPasswordProvider) KopiaArguments() ([]string, error) {
	next, err := p.ProviderStage.EnsureNextStage(func() (string, error) {
		return os.Getenv(p.Variable), nil
	})
	if err != nil {
		return nil, err
	}

	return next.KopiaArguments()
}
