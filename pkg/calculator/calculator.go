package calculator

import (
	"errors"
	"fmt"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"strconv"
)

var buttons = []string{
	"C", "()", "%", "/",
	"7", "8", "9", "*",
	"4", "5", "6", "-",
	"1", "2", "3", "+",
	"+=", "0", ".", "=",
}

type operation func(a float64, b float64) (string, float64, error)

func add(a float64, b float64) (string, float64, error) {
	return "+", a+b, nil
}

func subtract(a float64, b float64) (string, float64, error) {
	return "-", a-b, nil
}

func multiply(a float64, b float64) (string, float64, error) {
	return "*", a * b, nil
}

func divide(a float64, b float64) (string, float64, error) {
	if b == 0 {
		return "/", 0, errors.New("Cannot divide by zerro")
	}
	return "/", a / b, nil
}

type Calculator struct {
	app.Compo
	previous *float64
	symbol *string
	current *float64
	operation *operation
	err error
}

func (c *Calculator) numberButton(number uint64) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		if c.current == nil {
			v := float64(number)
			c.current = &v
		} else {
			v := *c.current * 10.0 + float64(number)
			c.current = &v
		}
		c.err = nil
		c.Update()
	}
}

func (c *Calculator) clear(ctx app.Context, e app.Event) {
	c.previous = nil
	c.current = nil
	c.operation = nil
	c.symbol = nil
	c.err = nil
	c.Update()
}

func (c *Calculator) operationButton(op operation) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		var v float64
		var s string
		if c.operation != nil && &c.previous != nil && c.current != nil {
			s, v, c.err = (*c.operation)(*c.previous, *c.current)
		} else if c.current != nil {
			v = *c.current
			s = ""
		} else if c.previous != nil {
			v = *c.previous
			s = ""
		} else {
			v = 0
			s = ""
		}
		if c.err == nil {
			c.previous = &v
			c.current = nil
			c.operation = &op
			if s == "" {
				c.symbol = nil
			} else {
				c.symbol = &s
			}
		}
		c.Update()
	}
}

func (c *Calculator) HandleButton(button string) app.EventHandler {
	if '0' <= button[0] && button[0] <= '9' {
		number, _ := strconv.ParseUint(button, 10, 64)
		return c.numberButton(number)
 	}
 	switch button {
	case "C":
		return c.clear
	case "+":
		return c.operationButton(add)
	case "-":
		return c.operationButton(subtract)
	case "/":
		return c.operationButton(divide)
	case "*":
		return c.operationButton(multiply)
	case "=":
		return c.operationButton(nil)
	default:
		return func(ctx app.Context, e app.Event) {

		}
	}

}

func (c *Calculator) renderNumber() app.UI {
	if c.current != nil {
		return app.Input().Class("display").Value(strconv.FormatFloat(*c.current, 'f', -1, 64))
	} else if c.previous != nil {
		return app.Input().Class("display").Value(strconv.FormatFloat(*c.previous, 'f', -1, 64))
	} else {
		return app.Input().Class("display").Value(fmt.Sprintf("0"))
	}
}

func (c *Calculator) Render() app.UI {
	return app.Div().Class("wrapper").Body(
		c.renderNumber(),
		app.Range(buttons).Slice(func(i int) app.UI {
			return app.Button().Text(buttons[i]).OnClick(c.HandleButton(buttons[i]))
		}),
		app.Div().Class("error").Text(c.err),
	)
}

