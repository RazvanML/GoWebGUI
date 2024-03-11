package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type control interface {
	render() string
	setOwner(p *page)
	event(string)
}

type htmlControl struct {
	control
	owner     *page
	id        string
	callbacks map[string]func()
}

func (h *htmlControl) event(ev string) {
	if h.callbacks == nil {
		return
	}
	h.callbacks[ev]()
}

func (h *htmlControl) mapEvent(ev string, event func()) {
	if h.callbacks == nil {
		h.callbacks = make(map[string]func())
	}
	h.callbacks[ev] = event
}

var id_counter int = 0

func (h *htmlControl) setOwner(p *page) {
	if h.owner != nil {
		panic("Control is already owned by a page.")
	}
	if h.id == "" {
		id_counter++
		h.id = fmt.Sprintf("id%d", id_counter)
	}

	h.owner = p
	p.ids[h.id] = h
}

type composite struct {
	htmlControl
	controls []control
}

func (c *composite) append(n control) {
	c.controls = append(c.controls, n)
	n.setOwner(c.owner)
}

func (c *composite) render() string {
	ret := ""
	for _, c1 := range c.controls {
		ret = ret + c1.render()
	}
	return ret
}

type page struct {
	composite
	ids map[string]control
}

func (p *page) render() string {
	ret := fmt.Sprintf(`
	<!DOCTYPE html>
<html lang="en">

<head>
  <title>HTML5 Boilerplate</title>
  <script> 
  function onClick(button) { 
	  var state = new XMLHttpRequest(); 
	  state.onload = function () { 
		  document.getElementById("container") 
			  .innerHTML = state.getAllResponseHeaders(); 
	  } 
	  state.open("GET", button.id+"/onclick", true); 
	  state.send(); 
  } 
</script> </head>

<body>
  %s
</body>

</html>
	`, p.composite.render())
	return ret
}

func newPage() page {
	ret := page{}
	ret.owner = &ret
	ret.ids = make(map[string]control)
	return ret
}

type Button struct {
	htmlControl
	text string
}

func NewButton(text string, onclick func()) Button {
	ret := Button{text: text}
	ret.mapEvent("onclick", onclick)
	return ret
}

func (c *Button) render() string {
	return fmt.Sprintf("<Button onClick='onClick(this)' id='%s'>%s</Button>", c.id, c.text)
}

type app struct {
	page
}

func newApp() app {
	return app{newPage()}
}

func (a *app) run() {
	http.HandleFunc("/", a.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func (a *app) handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if len(path) == 0 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(a.render()))
	} else {
		for name, control := range a.page.ids {
			if strings.HasPrefix(path, name+"/") {
				event := path[len(name)+1:]
				control.event(event)
			}

		}
	}
}

func main() {

	app1 := newApp()
	button1 := NewButton("Button1", func() { print("Hello world 1\n") })
	button2 := NewButton("Button2", func() {
		b := NewButton("Button 3", nil)
		app1.append(&b)
	})

	app1.append(&button1)
	app1.append(&button2)
	app1.run()

}
