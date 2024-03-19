package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type control interface {
	setOwner(p *page)
	event(string)
	getTag() string
	getText() string
	getAttributes() *map[string]string
}

type htmlControl struct {
	//	control
	text       string
	tag        string
	attributes map[string]string
	owner      *page
	callbacks  map[string]func()
}

var id_counter int = 0

func (h *htmlControl) setOwner(p *page) {
	if h.owner != nil {
		panic("Control is already owned by a page.")
	}
	if h.attributes == nil {
		h.attributes = map[string]string{}
	}
	_, i := h.attributes["id"]
	if !i {
		id_counter++
		h.attributes["id"] = fmt.Sprintf("id%d", id_counter)
	}

	var cc control = h
	p.ids[h.attributes["id"]] = &cc
	h.owner = p
}

func (h *htmlControl) event(ev string) {
	if h.callbacks == nil {
		return
	}
	h.callbacks[ev]()
}

func (h htmlControl) getTag() string {
	return h.tag
}

func (h htmlControl) getText() string {
	return h.text
}

func (h htmlControl) getAttributes() *map[string]string {
	return &h.attributes
}

func (h *htmlControl) setAttr(name string, val string) {
	if h.attributes == nil {
		h.attributes = make(map[string]string)
	}
	h.attributes[name] = val
}

func (h *htmlControl) mapEvent(ev string, event func()) {
	if h.callbacks == nil {
		h.callbacks = make(map[string]func())
	}
	h.callbacks[ev] = event
}

type composite interface {
	append(n *control)
	insertAt(n *control, pos int)
}

type compositeHtml struct {
	//	composite
	htmlControl
	controls []*control
}

func (c *compositeHtml) append(n *control) {
	c.controls = append(c.controls, n)
	(*n).setOwner(c.owner)
	var cc composite = c
	c.owner.addControl(cc, nil, *n)
}

func (c *compositeHtml) insertAt(n *control, pos int) {
	if pos > len(c.controls) {
		c.append(n)
		return
	}
	c.controls = append(c.controls[:pos+1], c.controls[pos:]...)
	c.controls[pos] = n

	(*n).setOwner(c.owner)
}

type page struct {
	compositeHtml
	ids    map[string]*control
	buffer string
}

func (p *page) addControl(parent composite, where *htmlControl, who control) {
	var parentCtrl control = parent.(control)

	attrs, _ := json.Marshal(*who.getAttributes())
	p.buffer += fmt.Sprintf("addElement('%s',null,'%s','%s','%s');", (*parentCtrl.getAttributes())["id"],
		who.getTag(), attrs, who.getText())
}

func (p page) render() string {
	ret := `
	<!DOCTYPE html>
<html lang="en">

<head>
  <title>HTML5 Example</title>
  <script> 

  function addElement(parent, before, tag, attrs, text ) {
	var p = document.getElementById(parent)
	var x = document.createElement(tag)
	for ( y in Object.keys(attrs)) {
		x[y] = attrs[y]
	}
	x.innerText = text
	if (before == null) {
	   p.append(x)	
	} else {
		var b2 = document.getElementById("before")
		p.insertBefore(x,b2)
	}
  }

  function onClick(button) { 
	  var state = new XMLHttpRequest(); 
	  state.onload = function () { 
		eval (
		 state.responseText
		 ) 
	  } 
	  state.open("GET", button.id+"/onclick", true); 
	  state.send(); 
  } 

  function getData() { 
	var state = new XMLHttpRequest(); 
	state.onload = function () { 
	  eval (
	   state.responseText
	   ) 
	} 
//	state.open("GET", button.id+"/onclick", true); 
	state.open("GET", "/_refresh", true); 
	state.send(); 
} 

</script> </head>

<body onLoad="getData()" id="body0">
</body>
</html>
	`
	return ret
}

func newPage() *page {
	ret := &page{}
	ret.ids = make(map[string]*control)
	ret.setAttr("id", "body0")
	ret.owner = ret
	return ret
}

type Button struct {
	htmlControl
}

func NewButton(text string, onclick func()) *Button {
	ret := &Button{htmlControl{text: text, tag: "Button"}}
	ret.mapEvent("onclick", onclick)
	return ret
}

type app struct {
	page
}

func newApp() *app {
	ret := app{page: *newPage()}
	ret.owner = &ret.page
	return &ret
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
	} else if path == "_refresh" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(a.buffer))
		print("Refresh called with: " + a.buffer + "\n")
		a.buffer = ""
	} else {
		for name, control := range a.page.ids {
			if strings.HasPrefix(path, name+"/") {
				event := path[len(name)+1:]
				(*control).event(event)
			}
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(a.buffer))
		//		a.buffer = ""
	}
}

func main() {

	app1 := newApp()

	/*
		x := &htmlControl{}
		c := x.(*control)
		c.setOwner(&app1.page)
	*/

	button1 := NewButton("Button1", func() { print("Hello world 1\n") })
	button2 := NewButton("Button2", func() {
		b := NewButton("Button 3", nil)
		var cc control = b
		app1.append(&cc)
	})

	var cc control = button1
	app1.append(&cc)
	cc = button2
	app1.append(cc)
	app1.run()
}
