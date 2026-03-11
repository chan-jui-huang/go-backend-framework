package config

import (
	"reflect"
	"sync"

	booterpkg "github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/spf13/viper"
)

type Loader struct {
	booterConfig *booterpkg.Config
	viper        *viper.Viper
}

func NewLoader(v *viper.Viper, booterConfig *booterpkg.Config) *Loader {
	return &Loader{
		booterConfig: booterConfig,
		viper:        v,
	}
}

func (l *Loader) Unmarshal(key string, cfg any) {
	if err := l.viper.UnmarshalKey(key, cfg); err != nil {
		panic(err)
	}
}

func (l *Loader) UnmarshalMany(configs map[string]any) {
	for key, cfg := range configs {
		l.Unmarshal(key, cfg)
	}
}

func (l *Loader) BooterConfig() *booterpkg.Config {
	return l.booterConfig
}

func (l *Loader) Viper() *viper.Viper {
	return l.viper
}

type registry struct {
	viper   *viper.Viper
	configs map[string]any
	mu      sync.RWMutex
}

var Registry = &registry{configs: map[string]any{}}

func NewRegistry(v *viper.Viper) *registry {
	return &registry{
		viper:   v,
		configs: map[string]any{},
	}
}

func (r *registry) Set(key string, cfg any) {
	if reflect.ValueOf(cfg).Kind() != reflect.Pointer {
		panic("config is not the pointer")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.configs[key] = cfg
}

func (r *registry) SetMany(configs map[string]any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for key, cfg := range configs {
		if reflect.ValueOf(cfg).Kind() != reflect.Pointer {
			panic("config is not the pointer")
		}
		r.configs[key] = cfg
	}
}

func (r *registry) Register(key string, cfg any) {
	if err := r.viper.UnmarshalKey(key, cfg); err != nil {
		panic(err)
	}
	r.Set(key, cfg)
}

func (r *registry) RegisterMany(configs map[string]any) {
	for key, cfg := range configs {
		r.Register(key, cfg)
	}
}

func (r *registry) Get(key string) any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return reflect.ValueOf(r.configs[key]).Elem().Interface()
}

func (r *registry) Unset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.configs, key)
}

func (r *registry) UnsetMany(keys ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, key := range keys {
		delete(r.configs, key)
	}
}

func (r *registry) SetViper(v *viper.Viper) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.viper = v
}

func (r *registry) GetViper() viper.Viper {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return *r.viper
}
