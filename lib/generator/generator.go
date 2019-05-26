package generator

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

type configType uint

const (
	configTypeNone configType = iota
	configTypeString
	configTypeEmbed
)

type configValue struct {
	Name string

	Type configType

	String []string
	Embed  []byte
}

type Config struct {
	values []configValue
	keys   map[string]struct{}
	mu     sync.Mutex
}

func New() *Config {
	return &Config{
		values: []configValue{},
		keys:   make(map[string]struct{}),
	}
}

func (cfg *Config) pushValue(value *configValue, isUnique bool) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.keys[value.Name]; isUnique && ok {
		return errors.New("key was already defined")
	}

	cfg.values = append(cfg.values, *value)
	cfg.keys[value.Name] = struct{}{}

	return nil
}

func (cfg *Config) Remove(name string) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.keys[name]; !ok {
		return errors.New("key doesn't exist")
	}

	for i := range cfg.values {
		if cfg.values[i].Name == name {
			cfg.values = append(cfg.values[:i], cfg.values[i+1:]...)
			delete(cfg.keys, name)
			return nil
		}
	}

	panic("unreachable")
}

func (cfg *Config) MustEnable(name string) {
	panicIfErr(cfg.Enable(name))
}

func (cfg *Config) MustEmbed(name string, value []byte) {
	panicIfErr(cfg.Embed(name, value))
}

func (cfg *Config) MustAdd(name string, value ...interface{}) {
	panicIfErr(cfg.Add(name, value...))
}

func (cfg *Config) MustSet(name string, value ...interface{}) {
	panicIfErr(cfg.Set(name, value...))
}

func (cfg *Config) Embed(name string, value []byte) error {
	value = bytes.TrimSpace(value)

	if len(value) == 0 {
		return errors.New("missing embedded value")
	}

	return cfg.pushValue(&configValue{
		Name:  name,
		Type:  configTypeEmbed,
		Embed: value,
	}, true)
}

func (cfg *Config) Enable(name string) error {
	return cfg.pushValue(&configValue{
		Name: name,
	}, true)
}

func (cfg *Config) Add(name string, value ...interface{}) error {
	if len(value) == 0 {
		return errors.New("at least one value is required")
	}

	strValues := make([]string, 0, len(value))
	for i := range value {
		strValues = append(strValues, fmt.Sprintf("%v", value[i]))
	}

	return cfg.pushValue(&configValue{
		Name:   name,
		Type:   configTypeString,
		String: strValues,
	}, false)
}

func (cfg *Config) Set(name string, value ...interface{}) error {
	_ = cfg.Remove(name)
	return cfg.Add(name, value...)
}

func (cfg *Config) Compile() ([]byte, error) {
	return compile(cfg.values)
}

func panicIfErr(err error) {
	if err != nil {
		panic(fmt.Errorf("unexpected internal error: %v", err.Error()))
	}
}
