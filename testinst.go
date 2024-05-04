// compare to a c++ example: http://www.tutorialspoint.com/cplusplus/cpp_interfaces.htm
package main

import (
	"fmt"
	"math"
)

// interface
type Shape interface {
	Area() float64
	GetWidth() float64
	GetHeight() float64
	SetWidth(float64)
	SetHeight(float64)
}

// reusable part, only implement SetWidth and SetHeight method of the interface
// {

type WidthHeight struct {
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

// }

type Rectangle struct {
	WidthHeight
}

func (this *Rectangle) Area() float64 {
	return this.GetWidth() * this.GetHeight() / 2
}

// override
func (this *Rectangle) GetHeight() float64 {
	fmt.Println("in Rectangle.GetHeight")
	// in case you still needs the WidthHeight's GetHeight method
	return this.WidthHeight.GetHeight()
}

type hasArea interface {
	Area() float64
}

func printArea(a hasArea) {
	fmt.Println("Print area: ", a.Area())
}

func setArea(a *hasArea, val float64) {
	r, ok := (*a).(*Rectangle)
	if ok {
		mult := math.Sqrt(val / r.Area())
		r.height *= mult
		r.width *= mult
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

func main2() {
	var r Rectangle
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

	setArea(&x, 1000)
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

}
