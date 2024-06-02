package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type control interface {
	setOwner(p *page)
	event(string)
	getTag() string
	getText() string
	getAttributes() *map[string]string
	getCallbacks() map[string]func(string)
	append(n control)
	insertAt(n control, pos int)
	enable()
	disable()
	isEnabled() bool
}

type htmlControl struct {
	control
	text       string
	tag        string
	attributes map[string]string
	owner      *page
	callbacks  map[string]func(string)
	syncProps  []string
	controls   []*control // child controls
}

var id_counter int = 0

func (h *htmlControl) enable() {
	h.setAttr("disabled")
}

func (h *htmlControl) disable() {
	h.setAttr("disabled", "1")
}

func (h *htmlControl) isEnabled() bool {
	disabled := true
	_, ok := (*h.getAttributes())["disabled"]
	if !ok {
		disabled = false
	}
	return !disabled
}

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

func (h htmlControl) event(ev string) {
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

func (h *htmlControl) setAttr(name ...string) {
	if h.attributes == nil {
		h.attributes = make(map[string]string)
	}
	if len(name) == 2 {
		h.attributes[name[0]] = name[1]
	} else {
		delete(h.attributes, name[0])
	}
	if h.owner != nil {
		if len(name) == 2 {
			h.owner.buffer += fmt.Sprintf("setElementAttr('%s','%s','%s')", h.attributes["id"], name[0], name[1])
		} else {
			h.owner.buffer += fmt.Sprintf("removeElementAttr('%s','%s')", h.attributes["id"], name[0])
		}
	}
}

func (h *htmlControl) mapEvent(ev string, event func(string)) {
	if h.callbacks == nil {
		h.callbacks = make(map[string]func(string))
	}
	h.callbacks[ev] = event
}

func (c htmlControl) append(n control) {
	c.controls = append(c.controls, &n)
	n.setOwner(c.owner)
	c.owner.addControl(&c, nil, n)
}

func (c *htmlControl) insertAt(n control, pos int) {
	if pos > len(c.controls) {
		c.append(n)
		return
	} else {
		c.controls = append(c.controls[:pos+1], c.controls[pos:]...)
		c.controls[pos] = &n
		n.setOwner(c.owner)
		c.owner.addControl(c, *c.controls[pos+1], n)
	}

}

type page struct {
	htmlControl
	id     string
	ids    map[string]*control
	buffer string
}

func (p *page) addControl(parent control, where control, who control) {
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
		wherestr = "'" + (*where.getAttributes())["id"] + "'"
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

  function setElementAttr(id, attr, val) {
	var e = document.getElementById(id)
	e.setAttribute(attr,val)
  }

  function removeElementAttr(id, attr) {
	var e = document.getElementById(id)
	e.removeAttribute(attr)
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

func newPage(id string) *page {
	ret := &page{}
	ret.id = id
	ret.ids = make(map[string]*control)
	ret.setAttr("id", "body0")
	ret.owner = ret
	return ret
}

type Button struct {
	htmlControl
}

func NewButton(text string, onClick func(string)) Button {
	ret := Button{htmlControl{text: text, tag: "Button"}}
	if onClick != nil {
		ret.mapEvent("click", onClick)
	}
	return ret
}

type TextInput struct {
	htmlControl
	val string
}

func NewTextInput(text string, onChange func(string)) TextInput {
	ret := TextInput{htmlControl: htmlControl{text: "", tag: "input"}, val: text}
	ret.setAttr("type", "input")
	(*ret.getAttributes())["type"] = "input"
	if onChange != nil {
		ret.mapEvent("change", onChange)
	}
	return ret
}

type app struct {

	// pattern and page. This is the initial state of the pages.
	pages map[string]func(string) *page

	// pages after they are created
	userPages map[string]*page
}

func newApp() *app {

	ret := app{pages: make(map[string]func(string) *page), userPages: make(map[string]*page)}
	return &ret
}

func (a *app) run() {
	http.HandleFunc("/", a.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (a *app) handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	pg, exists := a.pages[path]
	if exists {
		id := uuid.New().String()
		page := pg(id)
		page.owner = page
		a.userPages[id] = page
		// todo - cleanup user pages of the expired sessions
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session", Value: id, Expires: expiration}
		http.SetCookie(w, &cookie)

		w.Write([]byte(page.render()))
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		w.Write([]byte("Invalid access"))
		return
	}
	page := a.userPages[cookie.Value]

	if path == "_refresh" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(page.buffer))
		print("Refresh called with: " + page.buffer + "\n")
		page.buffer = ""
	} else if strings.HasPrefix(path, "convey/") {
		arr := strings.Split(path, "/")
		ctrl := *page.ids[arr[1]]
		reqBody, _ := io.ReadAll(r.Body)
		ctrl.getCallbacks()[arr[2]](string(reqBody))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(page.buffer))
		print("Refresh called with: " + page.buffer + "\n")
		page.buffer = ""
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "Cannot find: "+path)
	}
}

func main() {

	app := newApp()

	app.pages["page1"] = func(id string) *page {
		p := newPage(id)
		button1 := NewButton("Button1", func(ss string) { print("Hello world 1\n" + ss) })
		button2 := NewButton("Button2", func(ss string) {
			b := NewButton("Button 3", nil)
			p.insertAt(&b, 1)
		})
		p.append(&button1)
		p.append(&button2)

		text := NewTextInput("text value", nil)
		p.append(&text)

		return p
	}

	app.pages["page2"] = func(id string) *page {
		p := newPage(id)
		button1 := NewButton("Button1", func(ss string) { print("Hello world 1\n" + ss) })
		button1_d := NewButton("Trigger button1", func(ss string) {
			if button1.isEnabled() {
				button1.disable()
			} else {
				button1.enable()
			}
		})
		p.append(&button1)
		p.append(&button1_d)
		return p
	}

	app.run()
}
