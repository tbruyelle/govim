package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/myitcv/govim"
)

type Driver struct {
	*govim.Govim
}

func (d *Driver) Do(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case ErrDriver:
				err = r
			default:
				panic(r)
			}
		}
	}()
	return f()
}

func (d *Driver) DoFunction(f govim.VimFunction) govim.VimFunction {
	return func(args ...json.RawMessage) (interface{}, error) {
		var i interface{}
		err := d.Do(func() error {
			var err error
			i, err = f(args...)
			return err
		})
		if err != nil {
			return nil, err
		}
		return i, nil
	}
}

func (d *Driver) DoRangeFunction(f govim.VimRangeFunction) govim.VimRangeFunction {
	return func(first, last int, args ...json.RawMessage) (interface{}, error) {
		var i interface{}
		err := d.Do(func() error {
			var err error
			i, err = f(first, last, args...)
			return err
		})
		if err != nil {
			return nil, err
		}
		return i, nil
	}
}

func (d *Driver) Errorf(format string, args ...interface{}) {
	panic(ErrDriver{Underlying: fmt.Errorf(format, args...)})
}

func (d *Driver) ChannelExpr(expr string) json.RawMessage {
	i, err := d.Govim.ChannelExpr(expr)
	if err != nil {
		d.Errorf("ChannelExpr(%q) failed: %v", expr, err)
	}
	return i
}

func (d *Driver) ChannelEx(expr string) {
	if err := d.Govim.ChannelEx(expr); err != nil {
		d.Errorf("ChannelEx(%q) failed: %v", expr, err)
	}
}

func (d *Driver) ParseString(j json.RawMessage) string {
	var v string
	if err := json.Unmarshal(j, &v); err != nil {
		d.Errorf("failed to parse string from %q: %v", j, err)
	}
	return v
}

func (d *Driver) ParseInt(j json.RawMessage) int {
	var v int
	if err := json.Unmarshal(j, &v); err != nil {
		d.Errorf("failed to parse int from %q: %v", j, err)
	}
	return v
}

func (d *Driver) ChannelExprf(format string, args ...interface{}) json.RawMessage {
	return d.ChannelExpr(fmt.Sprintf(format, args...))
}

func (d *Driver) ChannelExf(format string, args ...interface{}) {
	d.ChannelEx(fmt.Sprintf(format, args...))
}

func (d *Driver) DefineFunction(name string, args []string, f govim.VimFunction) {
	if err := d.Govim.DefineFunction(name, args, d.DoFunction(f)); err != nil {
		d.Errorf("failed to DefineFunction %q: %v", name, err)
	}
}

func (d *Driver) DefineRangeFunction(name string, args []string, f govim.VimRangeFunction) {
	if err := d.Govim.DefineRangeFunction(name, args, d.DoRangeFunction(f)); err != nil {
		d.Errorf("failed to DefineRangeFunction %q: %v", name, err)
	}
}

func (d *Driver) DefineCommand(name string, f govim.VimCommandFunction, attrs ...govim.CommAttr) {
	if err := d.Govim.DefineCommand(name, f, attrs...); err != nil {
		d.Errorf("failed to DefineCommand %q: %v", name, err)
	}
}

type ErrDriver struct {
	Underlying error
}

func (e ErrDriver) Error() string {
	return fmt.Sprintf("driver error: %v", e.Underlying)
}