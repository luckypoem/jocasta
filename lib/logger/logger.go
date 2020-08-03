package logger

import (
	"log"
	"sync/atomic"
)

// Logger log interface
type Logger interface {
	Debugf(format string, arg ...interface{})
	Infof(format string, arg ...interface{})
	Warnf(format string, arg ...interface{})
	Errorf(format string, arg ...interface{})
	DPanicf(format string, arg ...interface{})
	Fatalf(format string, arg ...interface{})
}

// Std 标准输出,	os.Stdout, os.Stderr, os.Stdin, ioutil.Discard
type Std struct {
	*log.Logger
	has uint32
}

var _ Logger = (*Std)(nil)

// Option option
type Option func(std *Std)

// WithEnable enable log or not
func WithEnable(enable bool) Option {
	return func(std *Std) {
		std.Mode(enable)
	}
}

// New new std logger with option
func New(l *log.Logger, opts ...Option) *Std {
	s := &Std{l, 0}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Mode enable log or not
func (sf *Std) Mode(enable bool) {
	if enable {
		atomic.StoreUint32(&sf.has, 1)
	} else {
		atomic.StoreUint32(&sf.has, 0)
	}
}

// Debugf implement Logger interface.
func (sf Std) Debugf(format string, args ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.Logger.Printf("[D]: "+format, args...)
	}
}

// Infof implement Logger interface.
func (sf Std) Infof(format string, args ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.Logger.Printf("[I]: "+format, args...)
	}
}

// Errorf implement Logger interface.
func (sf Std) Errorf(format string, args ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.Logger.Printf("[E]: "+format, args...)
	}
}

// Warnf implement Logger interface.
func (sf Std) Warnf(format string, args ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.Logger.Printf("[W]: "+format, args...)
	}
}

// DPanicf implement Logger interface.
func (sf Std) DPanicf(format string, args ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.Logger.Printf("[P]: "+format, args...)
	}
}

// Fatalf implement Logger interface.
func (sf Std) Fatalf(format string, args ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.Logger.Printf("[F]: "+format, args...)
	}
}
