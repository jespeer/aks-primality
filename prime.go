package main

import (
  "fmt"
  "math"
  "math/big"
)

func main() {
  var n int
  fmt.Print("Enter an integer n: ")
  fmt.Scan(&n)
  fmt.Println("Your number is:", n)
  fmt.Println(aks(n))
}

func aks(n int) string {
  //1.
  if check_perfect_power(n) {
    return "composite"
  }
  //2.
  r := find_smallest_r(n)
  //3.
  for a := 2; a <= int(math.Min(float64(r), float64(n-1))); a++ {
    if n % a == 0 {
      return "composite"
    }
  }
  //4.
  if n <= r {
    return "prime"
  }
  //5.
  if polynomials(r, n) {
    return "composite"
  }
  //6.
  return "prime"
}

//checks if given integer is a perfect power
func check_perfect_power(n int) bool {
    max := math.Sqrt(float64(n)) + 1
    for i := 2; i < int(max); i++ {
      a := n;
      b := 0;
      for a % i == 0 {
        a /= i
        b += 1
        if a == 1 {
          return true
        }
      }
    }
    return false
}

//finds the smallest r such that ord_r(n) > log_2(n)^2
func find_smallest_r(n int) int {
  bound := int(math.Ceil(math.Log2(float64(n)) * math.Log2(float64(n))))

  for r := 2; r < math.MaxInt64; r++ {
    if multiplicativeOrder(n, r) > bound {
      return r
    }
  }
  return -1
}

//returns the multiplicative order
func multiplicativeOrder(a, n int) int {
  if GCD(a, n) != 1 {
    return -1
  }
  result := 1
  k := 1
  for k < n {
    result = (result * a) % n
    if result == 1 {
      return k
    }
    k++
  }
  return -1
}

// if (X+a)^n != X^(n) + a (mod X^(r) - 1, n), returns true for composite
func polynomials(r, n int) bool {
    //getting the bound for the loop
    phi := math.Sqrt(float64(EulersTotient(r)))
    bin := math.Log2((float64(n)))
    a_bound := math.Floor(phi * bin)
    x := []*big.Int{}

    for a := 1; a < int(a_bound); a++ {
      //finds the coefficients of the polynomials
      //if they are not equal mod (x^r -1, n), return true
      x = multiplyPolynomial([]*big.Int{big.NewInt(int64(a)), big.NewInt(1)}, n, r)
      for i := 0; i < len(x); i++ {
        if (x[i].Mod(x[i], big.NewInt(int64(n)))).Int64() != 0 {
          return true
        }
      }
    }
    return false
}

//returns GCD of two integers
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

//returns eulers totient (phi) for an integer
func EulersTotient(n int) int {
  result := 1
  for i := 2; i < n; i++ {
    if GCD(i, n) == 1 {
      result++
    }
  }
  return result
}

//returns slice of integers
//performs fast polynomial multiplication for multiple values of a
func multiplyPolynomial(poly []*big.Int, exponent, r int) ([]*big.Int) {
  //creates a slice with length r
	x := make([]*big.Int, r)
	for i := 0; i < len(x); i++ {
		x[i] = big.NewInt(0)
	}

  //a is the head of the base polynomial
	a := poly[0]
  //set head equal to 1
	x[0] = big.NewInt(1)
	n := exponent

  //iterates downward through the exponent and gets the coefficients for that exponent
  //performs the odd combination
	for exponent > 0 {
		if exponent % 2 == 1 {
			x = combinePolynomial(x, poly, n, r)
		}
    //performs the even combination
		poly = combinePolynomial(poly, poly, n, r)
		exponent /= 2
	}

	x[0].Sub(x[0], a)
  //performs modulo X^r -1, n
	x[n % r].Sub(x[n % r], big.NewInt(1))

	return x
}

func combinePolynomial(p1 []*big.Int, p2 []*big.Int, n, r int) ([]*big.Int) {
	foo := big.NewInt(0)
  //create a zeroed slice with length r
	x := make([]*big.Int, r)
	for i := 0; i < len(x); i++ {
		x[i] = big.NewInt(0)
    fmt.Println(x[i])
	}

  //loop through polynomials
	for i := 0; i < len(p1); i++ {
		for j := 0; j < len(p2); j++ {
      //if in index
			if (i + j) < r {
        //multiply value in p1 with value in p2
				foo.Mul(p1[i], p2[j])
        //add multiplied value (foo) to position i+j in x
				x[i + j].Add(x[i + j], foo)
        //if overreach, adjust indices accordingly
			} else {
				foo.Mul(p1[i], p2[j])
        //add foo to position
				x[(i + j) % r].Add(x[(i + j) % r], foo)
			}
		}
	}
	return x
}