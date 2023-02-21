package pkg

import (
	"fmt"
	"runtime"
	"strings"
)

type Platform struct {
	OS   string
	Arch string
}

func (p *Platform) Set(v string) error {
	o, a, ok := strings.Cut(v, "/")
	if !ok {
		return fmt.Errorf("invalid platform: %s", v)
	}
	p.OS = o
	p.Arch = a
	return nil
}

func (p *Platform) Type() string {
	return "os/arch"
}

func (p *Platform) String() string {
	return p.OS + "/" + p.Arch
}

func DefaultPlatform() Platform {
	return Platform{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
