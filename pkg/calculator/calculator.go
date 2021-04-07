package calculator

import (
	"errors"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"math/big"
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

var numbers = []*big.Float{
	big.NewFloat(0),
	big.NewFloat(1),
	big.NewFloat(2),
	big.NewFloat(3),
	big.NewFloat(4),
	big.NewFloat(5),
	big.NewFloat(6),
	big.NewFloat(7),
	big.NewFloat(8),
	big.NewFloat(9),
	big.NewFloat(10),
}

func  calcTength() *big.Float {
	v := &big.Float{}
	v.Quo(numbers[1], numbers[10])
	return v
}

var tenth = calcTength()
var minusOne = big.NewFloat(-1)

type operation func(a *big.Float, b *big.Float) (string, *big.Float, error)

func add(a *big.Float, b *big.Float) (string, *big.Float, error) {
	c := &big.Float{}
	c.Add(a, b)
	return "+", c, nil
}

func subtract(a *big.Float, b *big.Float) (string, *big.Float, error) {
	c := &big.Float{}
	return "-", c.Sub(a, b), nil
}

func multiply(a *big.Float, b *big.Float) (string, *big.Float, error) {
	c := &big.Float{}
	return "×", c.Mul(a, b), nil
}

func divide(a *big.Float, b *big.Float) (string, *big.Float, error) {
	if b.Cmp(numbers[0]) == 0 {
		return "÷", nil, errors.New("Cannot divide by zerro")
	}
	c := &big.Float{}
	return "÷", c.Quo(a, b), nil
}

type Calculator struct {
	app.Compo
	previous   *big.Float
	symbol     *string
	current    *big.Float
	multiplier *big.Float
	operation  *operation
	err        error
}

func (c *Calculator) numberButton(number uint64) app.EventHandler {
	n := numbers[number]
	return func(ctx app.Context, e app.Event) {
		if c.multiplier == nil {
			if c.current == nil {
				c.current = n
			} else {
				x, v := &big.Float{}, &big.Float{}
				x.Mul(c.current, numbers[10])
				v.Add(x, n)
				c.current = v
			}
		} else {
			if c.current == nil {
				v := &big.Float{}
				v.Mul(n, c.multiplier)
				c.current = v
			} else {
				v, x := &big.Float{}, &big.Float{}
				x.Mul(n, c.multiplier)
				v.Add(c.current, x)
				c.current = v
			}
			nf := &big.Float{}
			nf.Quo(c.multiplier, numbers[10])
			c.multiplier = nf
		}
		c.err = nil
		c.Update()
	}
}

func (c *Calculator) clear(ctx app.Context, e app.Event) {
	c.previous = nil
	c.current = nil
	c.operation = nil
	c.multiplier = nil
	c.symbol = nil
	c.err = nil
	c.Update()
}

func (c *Calculator) operationButton(op operation) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		var v *big.Float
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
			c.previous = v
			c.current = nil
			c.operation = &op
			if s == "" {
				c.symbol = nil
			} else {
				c.symbol = &s
			}
			c.multiplier = nil
		}
		c.Update()
	}
}

func (c *Calculator) decimalPoint(ctx app.Context, e app.Event) {
	if c.current == nil {
		c.current = numbers[0]
	}
	if c.multiplier == nil {
		c.multiplier = tenth
	}
	c.Update()
}

func (c *Calculator) toggleSign(ctx app.Context, e app.Event) {
	if c.current == nil {
		c.current = numbers[0]
	}
	v := &big.Float{}
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
		displayValue = c.current.Text('g', 16)
		if c.multiplier != nil && !strings.Contains(displayValue, ".") {
			displayValue = displayValue + "."
		}
	} else if c.previous != nil {
		displayValue = c.previous.Text('g', 16)
	} else {
		if c.multiplier == nil {
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

