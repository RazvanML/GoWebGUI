// compare to a c++ example: http://www.tutorialspoint.com/cplusplus/cpp_interfaces.htm
package main

import (
	"fmt"
	"math"
)

// interface

type hasArea interface {
	Area() float64
}

type Shape interface {
	hasArea
	GetWidth() float64
	GetHeight() float64
	SetWidth(float64)
	SetHeight(float64)
}

// reusable part, only implement SetWidth and SetHeight method of the interface
// {

type WidthHeight struct {
	hasArea
	width  float64
	height float64
}

func (this *WidthHeight) SetWidth(w float64) {
	this.width = w
}
func (this *WidthHeight) SetHeight(h float64) {
	this.height = h
}
func (this *WidthHeight) GetWidth() float64 {
	return this.width
}
func (this *WidthHeight) GetHeight() float64 {
	fmt.Println("in WidthHeight.GetHeight")
	return this.height
}

func (this *WidthHeight) Area() float64 {
	return this.GetWidth() * this.GetHeight()
}

// }

type Rectangle struct {
	WidthHeight
}

// override
func (this *Rectangle) GetHeight() float64 {
	fmt.Println("in Rectangle.GetHeight")
	// in case you still needs the WidthHeight's GetHeight method
	return this.WidthHeight.GetHeight()
}

type Square struct {
	Rectangle
}

func (this *Square) SetHeight(h float64) {
	this.WidthHeight.height = h
	this.WidthHeight.width = h
}

func (this *Square) SetWidth(h float64) {
	this.WidthHeight.height = h
	this.WidthHeight.width = h
}

func printArea(a hasArea) {
	fmt.Println("Print area: ", a.Area())
}

func setArea2(a hasArea, val float64) {
	r, ok := a.(Shape)
	if ok {
		m := math.Sqrt(val / a.Area())
		r.SetHeight(m * r.GetHeight())
		r.SetWidth(m * r.GetWidth())
	} else {
		panic("I don't know how to extend area for this shape.")
	}
}

type Defective struct {
	Shape
}

func (this *Defective) GetHeight() float64 {
	fmt.Println("in defective.GetHeight")
	return 7.0
}

type Stringer interface {
	String() string
	setString(string)
}

type MyType struct {
	value string
}

func (m MyType) String() string { return m.value }

func (m MyType) setString(s string) {
	m.value = s
}

func setStringer(s Stringer, val string) {
	s.setString(val)
}

type IFace interface {
	SetSomeField(newValue string)
	GetSomeField() string
}

type Implementation struct {
	someField string
}

func (i Implementation) GetSomeField() string {
	return i.someField
}

func (i *Implementation) SetSomeField(newValue string) {
	i.someField = newValue
}

func Create() *Implementation {
	return &Implementation{someField: "Hello"}
}

func setField(i IFace, val string) {
	i.SetSomeField(val)
}

func main1() {
	var r Square
	r.SetHeight(5)
	var i Shape = &r
	i.SetWidth(4)
	i.SetHeight(6)

	fmt.Println(i)
	fmt.Println("width: ", i.GetWidth())
	fmt.Println("height: ", i.GetHeight())
	fmt.Println("area: ", i.Area())

	x := hasArea(&r)
	fmt.Println("cast: ", x.Area())
	printArea(x)

	//	setArea(&x, 1000)
	setArea2(&r, 100)
	fmt.Println("area: ", i.Area())

	//	setArea(i, 100)

	var d Defective
	d.GetHeight()

	wh1 := WidthHeight{width: 10, height: 200}
	rec1 := Rectangle{wh1}
	rec2 := Rectangle{WidthHeight{width: 10, height: 200}}

	m := make(map[int]*Rectangle)
	m[1] = &rec1
	m[2] = &rec2

	m[1].height = 7.9

	fmt.Println("rec1 height: ", rec1.height)

	m1 := MyType{value: "something"}

	setStringer(m1, "something else1")
	fmt.Println("Value: ", m1.String())

	var s1 Stringer
	s1 = m1
	s1.setString("something else2")
	fmt.Println("Value: ", m1.String())

	i1 := *Create()
	fmt.Println("Field:", i1.someField)
	setField(&i1, "test")
	fmt.Println("Field:", i1.someField)
}
