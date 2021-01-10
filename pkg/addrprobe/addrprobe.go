package addrprobe

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Log   zerolog.Logger
	Quit  chan bool
	Delay time.Duration
}

func New(cfg Config) *Service {
	if cfg.Delay == 0 {
		cfg.Delay = 5 * time.Second
	}

	return &Service{cfg: cfg}
}

type Service struct {
	cfg Config

	wg sync.WaitGroup
}

func (s *Service) Run(addresses []string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	probesOk := true

	res := make(chan bool)
	for _, a := range addresses {
		s.wg.Add(1)
		go func(addr string) {
			defer s.wg.Done()
			res <- s.Probe(ctx, addr)
		}(a)
	}

	go func() {
		for {
			select {
			case <-s.cfg.Quit:
				cancel()
			case v := <-res:
				if !v {
					probesOk = false
					cancel()
				}
			}
		}
	}()

	s.wg.Wait()
	close(res)

	return probesOk
}

func (s *Service) Probe(ctx context.Context, address string) bool {
	log := s.cfg.Log.With().Str("address", address).Logger()
	log.Info().Msg("Probing")

	ch := make(chan bool)
	go s.dial(ctx, address, ch, log)
	res := <-ch

	log.Info().Bool("success", res).Msg("Probing ended")

	return res
}

func (s *Service) dial(ctx context.Context, address string, chres chan bool, log zerolog.Logger) {
	var d net.Dialer

DialLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error().Err(ctx.Err()).Msg("Timed out before connection was made")
			chres <- false
			break DialLoop
		default:
		}

		cctx, cancel := context.WithTimeout(ctx, time.Millisecond*250)
		conn, err := d.DialContext(cctx, "tcp", address)
		cancel()

		if err == nil {
			conn.Close()
			chres <- true
			break DialLoop
		}

		log.Warn().Err(err).Msg("Error dialing")
		time.Sleep(s.cfg.Delay)
	}
}
