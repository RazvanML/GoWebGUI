GoWebGUI allows you to build web applications in a Model-View-Controller fashion, in the same style of a classical desktop GUI application.
A minimal application looks like:

```go
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
```

The main advantages of this approach are:
- releases the developer of the HTTP way of thinking and focus on the business logic
- Javascript controls can be separately developed and tested
- Simplicity, one needs to know Go only
