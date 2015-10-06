package system

import (
	"fmt"
	"net"
	"sort"
	"time"
)

type DNS struct {
	host        string
	resolveable bool
	addrs       []string
	Timeout     int64
	loaded      bool
	err         error
}

func NewDNS(host string, system *System) DNS {
	return DNS{host: host}
}

func (d *DNS) Host() string {
	return d.host
}

func (d *DNS) setup() error {
	if d.loaded {
		return d.err
	}
	d.loaded = true

	timeout := d.Timeout
	if timeout == 0 {
		timeout = 500
	}

	addrs, err := lookupHost(d.host, timeout)
	if err != nil {
		d.resolveable = false
		d.addrs = []string{}
		d.err = err
		return d.err
	}
	sort.Strings(addrs)
	d.resolveable = true
	d.addrs = addrs
	return nil
}

func (d *DNS) Addrs() ([]string, error) {
	err := d.setup()

	return d.addrs, err
}

func (d *DNS) Resolveable() (interface{}, error) {
	err := d.setup()

	return d.resolveable, err
}

// Stub out
func (d *DNS) Exists() (interface{}, error) {
	return false, nil
}

func lookupHost(host string, timeout int64) ([]string, error) {
	c1 := make(chan []string, 1)
	e1 := make(chan error, 1)
	go func() {
		addrs, err := net.LookupHost(host)
		if err != nil {
			e1 <- err
		}
		c1 <- addrs
	}()
	select {
	case res := <-c1:
		return res, nil
	case err := <-e1:
		return nil, err
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return nil, fmt.Errorf("DNS lookup timed out (%d milliseconds)", timeout)
	}
}
