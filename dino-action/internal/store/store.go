package store

import (
	"context"
	"os"
	"strings"

	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/config"
	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/delay"
	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/domain"
	"gitlab.com/stackvista/demo/kubecon2024/poi/pkg/otel"
	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/yaml.v2"
)

type Store interface {
	Get(ctx context.Context, name string) ([]domain.Action, error)
}

type theStore struct {
	cfg config.Config
}

var _ Store = &theStore{} // Compile time check to ensure that theStore implements Store

func NewStore(cfg config.Config) (Store, error) {
	return &theStore{cfg: cfg}, nil
}

func (s *theStore) loadData(ctx context.Context) (map[string][]domain.Action, error) {
	delay.PretendHeavyOperation()

	f, err := os.Open(s.cfg.StoreContents)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var actions []domain.DinoActions
	if err := yaml.NewDecoder(f).Decode(&actions); err != nil {
		return nil, err
	}

	return s.convert(ctx, actions), nil
}

func (s *theStore) convert(_ context.Context, actions []domain.DinoActions) map[string][]domain.Action {
	delay.PretendHeavyOperation()

	poisMap := make(map[string][]domain.Action)
	for _, poi := range actions {
		poisMap[poi.DinoName] = poi.Actions
	}

	return poisMap
}

func (s *theStore) Get(ctx context.Context, name string) ([]domain.Action, error) {
	ctx, span := otel.Tracer.Start(ctx, "Store.Get")
	defer span.End()

	if span.IsRecording() {
		span.SetAttributes(attribute.String("name", name))
	}

	delay.PretendHeavyOperation()

	data, err := s.loadData(ctx)
	if err != nil {
		return nil, err
	}

	for c, d := range data {
		if strings.EqualFold(c, name) {
			return d, nil
		}
	}

	err = &domain.DinoNotFound{Name: name}

	span.RecordError(err)

	return nil, err
}
