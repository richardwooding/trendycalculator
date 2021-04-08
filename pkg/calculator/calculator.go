package calculator

import (
	"errors"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/ericlagergren/decimal"
	"strconv"
	"strings"
)

var buttons = []string{
	"C", "()", "%", "÷",
	"7", "8", "9", "×",
	"4", "5", "6", "-",
	"1", "2", "3", "+",
	"±", "0", ".", "=",
}

var numbers = []*decimal.Big{
	decimal.New(0, 0),
	decimal.New(1, 0),
	decimal.New(2, 0),
	decimal.New(3, 0),
	decimal.New(4, 0),
	decimal.New(5, 0),
	decimal.New(6, 0),
	decimal.New(7, 0),
	decimal.New(8, 0),
	decimal.New(9, 0),
	decimal.New(10, 0),
}

var minusOne = decimal.New(-1, 0);

type operation func(a *decimal.Big, b *decimal.Big) (string, *decimal.Big, error)

func add(a *decimal.Big, b *decimal.Big) (string, *decimal.Big, error) {
	c := &decimal.Big{}
	c.Add(a, b)
	return "+", c, nil
}

func subtract(a *decimal.Big, b *decimal.Big) (string, *decimal.Big, error) {
	c := &decimal.Big{}
	return "-", c.Sub(a, b), nil
}

func multiply(a *decimal.Big, b *decimal.Big) (string, *decimal.Big, error) {
	c := &decimal.Big{}
	return "×", c.Mul(a, b), nil
}

func divide(a *decimal.Big, b *decimal.Big) (string, *decimal.Big, error) {
	if b.Cmp(numbers[0]) == 0 {
		return "÷", nil, errors.New("Cannot divide by zerro")
	}
	c := &decimal.Big{}
	return "÷", c.Quo(a, b), nil
}

type Calculator struct {
	app.Compo
	previous  *decimal.Big
	symbol    *string
	current   *decimal.Big
	scale     int
	operation *operation
	err       error
}

func (c *Calculator) numberButton(number uint64) app.EventHandler {
	n := numbers[number]
	return func(ctx app.Context, e app.Event) {
		if c.scale == 0 {
			if c.current == nil {
				c.current = n
			} else {
				x, v := &decimal.Big{}, &decimal.Big{}
				x.Mul(c.current, numbers[10])
				v.Add(x, n)
				c.current = v
			}
		} else {
			if c.scale <= 99 {
				if c.current == nil {
					v := &decimal.Big{}
					v.Mul(n, decimal.New(1, c.scale))
					c.current = v
				} else {
					v, x := &decimal.Big{}, &decimal.Big{}
					x.Mul(n, decimal.New(1, c.scale))
					v.Add(c.current, x)
					c.current = v
				}
				c.scale++
			} else {
				c.err = errors.New("limited to 99 decimal points")
			}
		}
		c.err = nil
		c.Update()
	}
}

func (c *Calculator) clear(ctx app.Context, e app.Event) {
	c.previous = nil
	c.current = nil
	c.operation = nil
	c.scale = 0
	c.symbol = nil
	c.err = nil
	c.Update()
}

func (c *Calculator) operationButton(op operation) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		var v *decimal.Big
		var s string
		if c.operation != nil && &c.previous != nil && c.current != nil {
			s, v, c.err = (*c.operation)(c.previous, c.current)
		} else if c.current != nil {
			v = c.current
			s = ""
		} else if c.previous != nil {
			v = c.previous
			s = ""
		} else {
			v = numbers[0]
			s = ""
		}
		if c.err == nil {
			c.previous = v.Reduce()
			c.current = nil
			c.operation = &op
			if s == "" {
				c.symbol = nil
			} else {
				c.symbol = &s
			}
			c.scale = 0
		}
		c.Update()
	}
}

func (c *Calculator) decimalPoint(ctx app.Context, e app.Event) {
	if c.current == nil {
		c.current = numbers[0]
	}
	if c.scale == 0 {
		c.scale = 1
	}
	c.Update()
}

func (c *Calculator) toggleSign(ctx app.Context, e app.Event) {
	if c.current == nil {
		c.current = numbers[0]
	}
	v := &decimal.Big{}
	v.Mul(c.current, minusOne)
	c.current = v
	c.Update()
}

func (c *Calculator) HandleButton(button string) app.EventHandler {
	if '0' <= button[0] && button[0] <= '9' {
		number, _ := strconv.ParseUint(button, 10, 64)
		return c.numberButton(number)
 	}
 	switch button {
	case ".":
		return c.decimalPoint
	case "C":
		return c.clear
	case "+":
		return c.operationButton(add)
	case "-":
		return c.operationButton(subtract)
	case "÷":
		return c.operationButton(divide)
	case "×":
		return c.operationButton(multiply)
	case "=":
		return c.operationButton(nil)
	case "±":
		return c.toggleSign
	default:
		return func(ctx app.Context, e app.Event) {

		}
	}
}

func (c *Calculator) renderNumber() app.UI {
	var displayValue string
	if c.current != nil {
		displayValue = c.current.String()
		if c.scale > 0 && !strings.Contains(displayValue, ".") {
			displayValue = displayValue + "."
		}
	} else if c.previous != nil {
		displayValue = c.previous.String()
	} else {
		if c.scale == 0 {
			displayValue = "0"
		} else {
			displayValue = "0."
		}
	}
	return app.Input().Class("display").Value(displayValue)
}

func (c *Calculator) Render() app.UI {
	return app.Div().Class("wrapper").Body(
		c.renderNumber(),
		app.Range(buttons).Slice(func(i int) app.UI {
			return app.Button().Class("key").Text(buttons[i]).OnClick(c.HandleButton(buttons[i]))
		}),
		app.Div().Class("error").Text(c.err),
	)
}

