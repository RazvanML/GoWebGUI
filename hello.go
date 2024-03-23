package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	getCallbacks() map[string]func(string)
}

type htmlControl struct {
	control
	text       string
	tag        string
	attributes map[string]string
	owner      *page
	callbacks  map[string]func(string)
	syncProps  []string
}

var id_counter int = 0

func (h *htmlControl) getCallbacks() map[string]func(string) {
	return h.callbacks
}

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
	h.callbacks[ev]("")
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

func (h *htmlControl) mapEvent(ev string, event func(string)) {
	if h.callbacks == nil {
		h.callbacks = make(map[string]func(string))
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
	} else {
		c.controls = append(c.controls[:pos+1], c.controls[pos:]...)
		c.controls[pos] = n
		(*n).setOwner(c.owner)
		var cc composite = c
		c.owner.addControl(cc, c.controls[pos+1], *n)
	}

}

type page struct {
	compositeHtml
	ids    map[string]*control
	buffer string
}

func (p *page) addControl(parent composite, where *control, who control) {
	var parentCtrl control = parent.(control)

	attrs, _ := json.Marshal(who.getAttributes())
	ev := who.getCallbacks()
	keys := make([]string, 0, len(ev))
	for k := range ev {
		keys = append(keys, k)
	}
	events, _ := json.Marshal(keys)

	wherestr := "null"
	if where != nil {
		wherestr = "'" + (*(*where).getAttributes())["id"] + "'"
	}

	p.buffer += fmt.Sprintf("addElement('%s',%s,'%s',%s,'%s',%s);\n",
		(*parentCtrl.getAttributes())["id"], wherestr,
		who.getTag(), attrs, who.getText(), events)
}

func (p page) render() string {
	ret := `
	<!DOCTYPE html>
<html lang="en">

<head>
  <title>HTML5 Example</title>
  <script> 


  function stringify_object(object, depth=0, max_depth=2) {
    // change max_depth to see more levels, for a touch event, 2 is good
    if (depth > max_depth)
        return 'Object';

    const obj = {};
    for (let key in object) {
        let value = object[key];
        if (value instanceof Node)
            // specify which properties you want to see from the node
            value = {id: value.id};
        else if (value instanceof Window)
            value = 'Window';
        else if (value instanceof Object)
            value = stringify_object(value, depth+1, max_depth);

        obj[key] = value;
    }

    return depth? obj: JSON.stringify(obj);
}

  controls = {}

  function addElement(parent, before, tag, attrs, text, events ) {
	var p = document.getElementById(parent)
	var x = document.createElement(tag)
	for ( y in attrs) {
		x.setAttribute(y,attrs[y])
	}
	x.innerText = text

	for (e in events) {
		x.addEventListener(events[e], (ev)=>conveyEvent(x,events[e],ev))
	}
	if (before == null) {
	   p.append(x)	
	} else {
		var b2 = document.getElementById(before)
		p.insertBefore(x,b2)
	}
	controls[x.id] = x
  }


  function conveyEvent(control, eventStr, event) {
	var state = new XMLHttpRequest(); 
	state.onload = function () { 
	  eval (
	   state.responseText
	   ) 
	} 
	state.open("POST", "convey/"+control.id+"/"+eventStr, true); 
	state.send(stringify_object(event)); 
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

func NewButton(text string, onclick func(string)) *Button {
	ret := &Button{htmlControl{text: text, tag: "Button"}}
	if onclick != nil {
		ret.mapEvent("click", onclick)
	}
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
	} else if strings.HasPrefix(path, "convey/") {
		arr := strings.Split(path, "/")
		ctrl := *a.page.ids[arr[1]]
		reqBody, _ := io.ReadAll(r.Body)
		ctrl.getCallbacks()[arr[2]](string(reqBody))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(a.buffer))
		print("Refresh called with: " + a.buffer + "\n")
		a.buffer = ""
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "Cannot find: "+path)
	}
}

func main() {

	app1 := newApp()
	button1 := NewButton("Button1", func(ss string) { print("Hello world 1\n" + ss) })
	button2 := NewButton("Button2", func(ss string) {
		b := NewButton("Button 3", nil)
		var cc control = b
		app1.insertAt(&cc, 1)
	})

	var cc control = button1
	app1.append(&cc)
	cc = button2
	app1.append(&cc)
	app1.run()
}
