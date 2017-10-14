package examples

func ExampleSum() {
	Sum(1, 5)
	//Output:
	// 1 + 5 = 6
}

func ExampleSum_negative() {
	Sum(-3, -9)
	//Output:
	// -3 + -9 = -12
}

func ExampleSum_zeroes() {
	Sum(0, 0)
	//Output:
	// 0 + 0 = 0
}
