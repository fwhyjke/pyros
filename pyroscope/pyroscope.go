package pyroscope

import (
	"context"
	"runtime"

	"git.ecom.tech/ecom/dev/pap/backend/go/log"
	"github.com/grafana/pyroscope-go"
)

const PyroscoprAddr = ":4040"

type Profiler struct {
	cfg pyroscope.Config
	log *log.Logger
	p   *pyroscope.Profiler
}

func New(opts ...Option) *Profiler {
	o := &Options{pyroscopeServerAddress: PyroscoprAddr}
	for _, optionFn := range opts {
		optionFn(o)
	}

	cfg := pyroscope.Config{
		ApplicationName: o.serviceName,
		ServerAddress:   o.pyroscopeServerAddress,
		Logger:          pyroscope.StandardLogger,

		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	}

	return &Profiler{
		cfg: cfg,
		log: log.NewNop(),
	}
}

func (p *Profiler) Start(ctx context.Context) error {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	profiler, err := pyroscope.Start(p.cfg)
	if err != nil {
		p.log.ErrorContext(ctx, "failed to listen pyroscope", log.Error(err))
		return err
	}
	p.log.Info("pyroscope started")
	p.p = profiler

	return nil
}

func (p *Profiler) Stop(ctx context.Context) error {
	if p.p != nil {
		p.log.Info("pyroscope shutting down...")
		return p.p.Stop()
	}
	return nil
}

type Option func(*Options)

type Options struct {
	pyroscopeServerAddress string
	serviceName            string
}

func WithPyroscopeServerAddress(addr string) Option {
	return func(h *Options) {
		h.pyroscopeServerAddress = addr
	}
}

func WithServiceName(name string) Option {
	return func(h *Options) {
		h.serviceName = name
	}
}

func (p *Profiler) WithLogger(l *log.Logger) {
	if l != nil {
		p.log = l
	}
}
