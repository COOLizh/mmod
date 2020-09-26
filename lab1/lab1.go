package lab1

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Lab1 : includes all needful information and methods for solving Lab1
type Lab1 struct {
	squareMethodBRV       []float64
	congruentialMethodBRV []float64
	brvCount              int
	bitDepth              int
}

func (l *Lab1) init(n int) {
	l.brvCount = n
	l.bitDepth = 8
}

func (l *Lab1) midSquareMethod() {
	var lastNum int64 = 34554334 //l.randInt(10000000, 99999999)
	l.squareMethodBRV = make([]float64, l.brvCount)
	for i := 0; i < l.brvCount; i++ {
		lastNum *= lastNum
		lastNumBytes := []byte(strconv.FormatInt(int64(lastNum), 10))
		zeroCount := 2*l.bitDepth - len(lastNumBytes)
		tmp := make([]byte, 2*l.bitDepth)
		for j := 0; j < 2*l.bitDepth; j++ {
			if j <= zeroCount-1 {
				tmp[j] = 48
			} else {
				tmp[j] = lastNumBytes[j-zeroCount]
			}
		}
		newNumBytes := make([]byte, l.bitDepth)
		for i, j := l.bitDepth/2, 0; i < len(tmp)-l.bitDepth/2; i++ {
			newNumBytes[j] = tmp[i]
			j++
		}
		zeroCount = 0
		for j := 0; j < len(newNumBytes); j++ {
			if newNumBytes[j] != 48 {
				zeroCount = j
				break
			}
		}
		newNumString := string(newNumBytes[zeroCount:])
		newNum, err := strconv.ParseInt(newNumString, 0, 64)
		if err != nil {
			panic(err)
		}
		l.squareMethodBRV[i] = float64(newNum) / math.Pow(10, float64(l.bitDepth))
		lastNum = newNum
	}
}

func (l *Lab1) multiplicativeCongruentialMethod() {
	m, k, prevA := 12312312, 23423423, 34534534 //l.randInt(10000000, 99999999), l.randInt(10000000, 99999999), l.randInt(10000000, 99999999)
	l.congruentialMethodBRV = make([]float64, l.brvCount)
	for i := 0; i < l.brvCount; i++ {
		newA := (k * prevA) % m
		l.congruentialMethodBRV[i] = float64(newA) / float64(m)
		prevA = newA
	}
}

func (l Lab1) testUniformity(n int, method string) {
	l.init(n)
	var mathExpectation float64
	var dispersion float64
	if method == "midSquareMethod" {
		l.midSquareMethod()
		p := l.createP(l.squareMethodBRV, l.brvCount)
		if n == 100 {
			fmt.Println("\n-----Testing uniformity for mid square method (n = 100)-----")
			l.createHistigram(p, "./lab1/hist1.png", "Uniformity test of mid square method (n = 100)")
		} else {
			fmt.Println("\n-----Testing uniformity for mid square method (n = 100000)-----")
			l.createHistigram(p, "./lab1/hist2.png", "Uniformity test of mid square method (n = 100000)")
		}
		mathExpectation = l.calculateMathExpectation(l.squareMethodBRV)
		dispersion = l.calculateDispersion(l.squareMethodBRV, mathExpectation)
	} else {
		l.multiplicativeCongruentialMethod()
		p := l.createP(l.congruentialMethodBRV, l.brvCount)
		if n == 100 {
			fmt.Println("\n-----Testing uniformity for multiplicative congruential method (n = 100)-----")
			l.createHistigram(p, "./lab1/hist3.png", "Uniformity test of for multiplicative congruential method (n = 100)")
		} else {
			fmt.Println("\n-----Testing uniformity for multiplicative congruential method (n = 100000)-----")
			l.createHistigram(p, "./lab1/hist4.png", "Uniformity test of  multiplicative congruential method(n = 100000)")
		}
		mathExpectation = l.calculateMathExpectation(l.congruentialMethodBRV)
		dispersion = l.calculateDispersion(l.congruentialMethodBRV, mathExpectation)
	}

	fmt.Print("Math expection: ")
	fmt.Println(mathExpectation)
	fmt.Print("Dispersion: ")
	fmt.Println(dispersion)
}

func (l Lab1) createHistigram(p []float64, fileName, histogramName string) {
	var values plotter.Values
	values = append(values, p...)
	pl, err := plot.New()
	if err != nil {
		panic(err)
	}
	pl.Title.Text = histogramName
	hist, err := plotter.NewBarChart(values, 10)
	if err != nil {
		panic(err)
	}
	pl.Add(hist)
	if err := pl.Save(5*vg.Inch, 5*vg.Inch, fileName); err != nil {
		panic(err)
	}
}

func (l Lab1) testIndependency(n int, method string) {
	l.init(n)
	var shift int64
	var x, y []float64
	if method == "midSquareMethod" {
		l.midSquareMethod()
		shift = l.randInt(2, int64(len(l.squareMethodBRV)))
		y, x = l.squareMethodBRV[shift:], l.squareMethodBRV[:int64(len(l.squareMethodBRV))-shift]
		if n == 100 {
			fmt.Println("\n-----Testing independency for mid square method (n = 100)-----")
		} else {
			fmt.Println("\n-----Testing independency for mid square method (n = 100000)-----")
		}
	} else {
		l.multiplicativeCongruentialMethod()
		shift = l.randInt(2, int64(len(l.congruentialMethodBRV)))
		y, x = l.congruentialMethodBRV[shift:], l.congruentialMethodBRV[:int64(len(l.congruentialMethodBRV))-shift]
		if n == 100 {
			fmt.Println("\n-----Testing independency for multiplicative congruential method(n = 100)-----")
		} else {
			fmt.Println("\n-----Testing independency for multiplicative congruential method (n = 100000)-----")
		}
	}
	var mathExpectionXY float64
	for i := 0; i < len(x); i++ {
		mathExpectionXY += x[i] * y[i]
	}
	mathExpectionXY /= float64(len(x))
	mathExpectionX, mathExpectionY := l.calculateMathExpectation(x), l.calculateMathExpectation(y)
	dispersionX, dispersionY := l.calculateDispersion(x, mathExpectionX), l.calculateDispersion(y, mathExpectionY)
	var correlationCoeff float64
	correlationCoeff = (mathExpectionXY - (mathExpectionX * mathExpectionY)) / (math.Sqrt(dispersionX * dispersionY))
	fmt.Print("Correaltion coefficient: ")
	fmt.Println(correlationCoeff)
}

func (l Lab1) createP(slice []float64, brvCount int) []float64 {
	n, p := make([]float64, 10), make([]float64, 10)
	for _, brv := range slice {
		switch {
		case brv <= 0.1:
			n[0]++
		case brv > 0.1 && brv <= 0.2:
			n[1]++
		case brv > 0.2 && brv <= 0.3:
			n[2]++
		case brv > 0.3 && brv <= 0.4:
			n[3]++
		case brv > 0.4 && brv <= 0.5:
			n[4]++
		case brv > 0.5 && brv <= 0.6:
			n[5]++
		case brv > 0.6 && brv <= 0.7:
			n[6]++
		case brv > 0.7 && brv <= 0.8:
			n[7]++
		case brv > 0.8 && brv <= 0.9:
			n[8]++
		case brv > 0.9 && brv <= 1:
			n[9]++
		}
	}
	for i := 0; i < 10; i++ {
		p[i] = float64(n[i]) / float64(brvCount)
	}
	return p
}

func (l Lab1) randInt(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Int63n(max-min)
}

func (l Lab1) calculateMathExpectation(slice []float64) float64 {
	var mathExpectation float64
	for _, item := range slice {
		mathExpectation += item
	}
	return mathExpectation / float64(len(slice))
}

func (l Lab1) calculateDispersion(slice []float64, mathExpectation float64) float64 {
	var dispersion float64
	for _, item := range slice {
		dispersion += math.Pow(item, 2) - math.Pow(mathExpectation, 2)
	}
	return dispersion / float64(len(slice))
}

// Solve : solving lab1
func (l Lab1) Solve() {
	n1, n2 := 100, 100000
	l.testUniformity(n1, "midSquareMethod")
	l.testUniformity(n2, "midSquareMethod")
	l.testUniformity(n1, "multiplicativeCongruentialMethod")
	l.testUniformity(n2, "multiplicativeCongruentialMethod")

	l.testIndependency(n1, "midSquareMethod")
	l.testIndependency(n2, "midSquareMethod")
	l.testIndependency(n1, "multiplicativeCongruentialMethod")
	l.testIndependency(n2, "multiplicativeCongruentialMethod")
}
