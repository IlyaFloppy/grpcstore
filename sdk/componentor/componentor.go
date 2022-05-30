package componentor

import (
	"context"

	"github.com/rs/zerolog"
)

type Component interface {
	Name() string
	Run(ctx context.Context) error
	ReadyCh() <-chan struct{}
}

func Run(ctx context.Context, logger zerolog.Logger, components []Component) (exitCode int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	baseCtx := context.Background()
	contexts := make([]context.Context, len(components))
	cancels := make([]context.CancelFunc, len(components))
	errorsCh := make([]chan error, len(components))
	for i := 0; i < len(components); i++ {
		contexts[i], cancels[i] = context.WithCancel(baseCtx)
		errorsCh[i] = make(chan error, 1)
	}

	for i, c := range components {
		go func(i int, c Component) {
			err := c.Run(contexts[i])
			if err != nil {
				cancel()
			}
			errorsCh[i] <- err
			close(errorsCh[i])
		}(i, c)

		select {
		case <-ctx.Done():
			break // NOTE: break from select; components[i+1:] will still run.
		case <-c.ReadyCh():
			logger.Info().
				Str("component", components[i].Name()).
				Int("idx", i).
				Int("len", len(components)).
				Msg("component is ready")
		}
	}

	<-ctx.Done()

	var errors []error
	for i := len(components) - 1; i >= 0; i-- {
		cancels[i]()
		err := <-errorsCh[i]

		logger := logger.With().
			Str("component", components[i].Name()).
			Int("idx", i).
			Int("len", len(components)).
			Logger()

		if err != nil {
			errors = append(errors, err)
			logger.Err(err).Msg("component finished with error")
		} else {
			logger.Info().Msg("component finished successfully")
		}
	}

	if len(errors) > 0 {
		logger.Error().Interface("errors", errors).Msg("grpcstore finished with errors")
		return 1
	}

	logger.Info().Msg("grpcstore finished")
	return 0
}
