package navigator

type navigatorOptions struct {
	elevationMin, elevationMax, elevationAmplitude float64
	minOffline, maxOffline                         int
	skipOffline                                    bool
}

type Option func(*navigatorOptions)

func defaultOptions() *navigatorOptions {
	return &navigatorOptions{
		elevationMin:       1,
		elevationMax:       500,
		elevationAmplitude: 4,
		minOffline:         1,
		maxOffline:         10,
	}
}

func SkipOfflineMode() Option {
	return func(opt *navigatorOptions) {
		opt.skipOffline = true
	}
}

func WithElevation(min, max float64, amplitude float64) Option {
	return func(opt *navigatorOptions) {
		opt.elevationMin, opt.elevationMax = min, max
		opt.elevationAmplitude = amplitude
	}
}

func WithOffline(min, max int) Option {
	return func(opt *navigatorOptions) {
		opt.minOffline, opt.maxOffline = min, max
	}
}
