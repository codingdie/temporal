package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
	"go.temporal.io/server/common/number"
)

type (
	TwoStepCircuitBreaker interface {
		Name() string
		State() gobreaker.State
		Counts() gobreaker.Counts
		Allow() (done func(success bool), err error)
	}

	// TwoStepCircuitBreakerWithDynamicSettings is a wrapper of gobreaker.TwoStepCircuitBreaker
	// that calls the settingsFn everytime the Allow function is called and replaces the circuit
	// breaker if there is a change in the settings object. Note that in this case, the previous
	// state of the circuit breaker is lost.
	TwoStepCircuitBreakerWithDynamicSettings struct {
		cb *gobreaker.TwoStepCircuitBreaker

		settingsFn   func() map[string]any
		baseSettings baseSettings

		name          string
		readyToTrip   func(counts gobreaker.Counts) bool
		onStateChange func(name string, from gobreaker.State, to gobreaker.State)
		isSuccessful  func(err error) bool
	}

	baseSettings struct {
		MaxRequests uint32
		Interval    time.Duration
		Timeout     time.Duration
	}
)

var _ TwoStepCircuitBreaker = (*TwoStepCircuitBreakerWithDynamicSettings)(nil)

const (
	maxRequestsKey = "MaxRequests"
	intervalKey    = "intervalKey"
	timeoutKey     = "timeout"

	// Zero values indicate to use the default values from gobreaker.Settings.
	defaultMaxRequests = uint32(0)
	defaultInterval    = 0 * time.Second
	defaultTimeout     = 0 * time.Second
)

func NewTwoStepCircuitBreakerWithDynamicSettings(
	settingsFn func() map[string]any,
) *TwoStepCircuitBreakerWithDynamicSettings {
	return &TwoStepCircuitBreakerWithDynamicSettings{
		settingsFn: settingsFn,
	}
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) WithName(
	name string,
) *TwoStepCircuitBreakerWithDynamicSettings {
	if c == nil {
		return nil
	}
	ret := *c
	ret.cb = nil
	ret.name = name
	return &ret
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) WithReadyToTrip(
	readyToTrip func(counts gobreaker.Counts) bool,
) *TwoStepCircuitBreakerWithDynamicSettings {
	if c == nil {
		return nil
	}
	ret := *c
	ret.cb = nil
	ret.readyToTrip = readyToTrip
	return &ret
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) WithOnStateChange(
	onStateChange func(name string, from gobreaker.State, to gobreaker.State),
) *TwoStepCircuitBreakerWithDynamicSettings {
	if c == nil {
		return nil
	}
	ret := *c
	ret.cb = nil
	ret.onStateChange = onStateChange
	return &ret
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) WithIsSuccessful(
	isSuccessful func(err error) bool,
) *TwoStepCircuitBreakerWithDynamicSettings {
	if c == nil {
		return nil
	}
	ret := *c
	ret.cb = nil
	ret.isSuccessful = isSuccessful
	return &ret
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) Name() string {
	if c.cb == nil {
		return ""
	}
	return c.cb.Name()
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) State() gobreaker.State {
	if c.cb == nil {
		return 0
	}
	return c.cb.State()
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) Counts() gobreaker.Counts {
	if c.cb == nil {
		return gobreaker.Counts{}
	}
	return c.cb.Counts()
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) Allow() (done func(success bool), err error) {
	if err := c.checkAndUpdateSettings(); err != nil {
		return nil, err
	}
	return c.cb.Allow()
}

func (c *TwoStepCircuitBreakerWithDynamicSettings) checkAndUpdateSettings() error {
	settingsMap := c.settingsFn()
	bs := baseSettings{
		MaxRequests: defaultMaxRequests,
		Interval:    defaultInterval,
		Timeout:     defaultTimeout,
	}

	if maxRequests, ok := settingsMap[maxRequestsKey]; ok {
		bs.MaxRequests = uint32(
			number.NewNumber(maxRequests).GetUintOrDefault(uint(defaultMaxRequests)),
		)
	}
	if interval, ok := settingsMap[intervalKey]; ok {
		bs.Interval = time.Duration(
			number.NewNumber(interval).GetIntOrDefault(int(defaultInterval.Seconds())),
		) * time.Second
	}
	if timeout, ok := settingsMap[timeoutKey]; ok {
		bs.Timeout = time.Duration(
			number.NewNumber(timeout).GetIntOrDefault(int(defaultTimeout.Seconds())),
		) * time.Second
	}

	if c.cb != nil && bs == c.baseSettings {
		return nil
	}

	c.baseSettings = bs
	c.cb = gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
		Name:          c.name,
		MaxRequests:   bs.MaxRequests,
		Interval:      bs.Interval,
		Timeout:       bs.Timeout,
		ReadyToTrip:   c.readyToTrip,
		OnStateChange: c.onStateChange,
		IsSuccessful:  c.isSuccessful,
	})
	return nil
}
