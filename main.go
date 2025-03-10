package main

import (
	"fmt"
	"time"
)

//const (
//	RedColorRGB   = 0xFF0000
//	GreenColorRGB = 0x00FF00
//)
//
//func (c car) getColorRgb() int {
//	if c.Model == "Opala" {
//		return RedColorRGB
//	}
//	return GreenColorRGB
//}

func contado(n int) {
	for i := range n {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	//   x := 5
	//   x = inc(x)

	ch := make(chan string)

	go func() {
		ch <- "Hello world ! Anonymous function and Channel"
	}()

	msg := <-ch
	fmt.Println(msg)

	fmt.Println("Hello, world!")
	println("Hello, world!")
	go contado(10)
	go contado(5)

	var y int
	y = 12 % 5

	println(y)

	var s []string

	// Crescendo de 1 até 30
	for i := 1; i <= 30; i++ {
		s = append(s, "*")
		fmt.Println(s)
	}

	// Diminuindo de 30 até 0
	for i := 30; i > 0; i-- {
		s = s[:i-1] // Reduz o slice
		fmt.Println(s)
	}

	//DIRETA
	//var car1 car
	//car1.Model = "Opala"
	//car1.getHorsePower()
	//
	////
	//var car2 car = struct {
	//	Make   string
	//	Model  string
	//	Height int
	//	Width  int
	//}{Make: "pokemon", Model: "Marea", Height: 1800, Width: 3000}

	//car2.getRgb()
	//car2.getHorsePower()

	//carFrunFrun(car1, car2)
	//verifyCarTypeAndColor(car1, car2)
}

//func inc(x int)(int, error) {
//  x++
//
//  fmt.Println(x)
//  if(x == 6){
//    return nil
//  }
//  return x
//}

//func carFrunFrun(carFrun ...car) {
//	for _, c := range carFrun {
//		switch c.Model {
//		case "Opala":
//			fmt.Println("Frun Frun")
//		case "Marea":
//			fmt.Println("Don`t Frun Frun")
//		default:
//			fmt.Printf("Unknown model: %s\n", c.Model)
//		}
//	}
//
//}

//
//func verifyCarTypeAndColor(cars ...car) {
//	for _, c := range cars {
//		if c.Model == "Opala" && c.getHorsePower() == 1800 && c.getRgb() == 0xFF0000 { // Red color in RGB
//			fmt.Printf("The car %s is a Sport car - Frun Frun!\n", c.Model)
//		} else {
//			fmt.Printf("The car %s is not a Sport car. It is Green.\n", c.Model)
//		}
//	}
//}

// Exemplo de loop "for" básico com uma variável de contagem
func simpleForLoop() {
	for i := 0; i < 5; i++ {
		fmt.Println("Contagem:", i)
	}
}

// Exemplo de loop "for" simulando um "while"
func whileLoop() {
	x := 0
	for x < 10 {
		fmt.Println("Valor de x:", x)
		x++
	}
}

// Exemplo de loop "for" com range percorrendo um slice
func rangeLoop() {
	cars := []string{"Opala", "Marea", "Fusca"}
	for index, car := range cars {
		fmt.Printf("Carro %d: %s\n", index+1, car)
	}
}

// Exemplo de loop "for" percorrendo um map
func forMap() {
	carColors := map[string]string{"Opala": "Red", "Marea": "Green", "Fusca": "Blue"}
	for model, color := range carColors {
		fmt.Printf("O modelo %s tem a cor %s\n", model, color)
	}
}

// Exemplo de loop "for" infinito com uma condição de parada
func forLoopInfinite() {
	z := 0
	for {
		fmt.Println("Valor de z:", z)
		z++
		if z == 5 {
			fmt.Println("Parando o loop!")
			break
		}
	}
}
