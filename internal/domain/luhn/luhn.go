package luhn

import "strconv"

func Valid(number string) bool {
	n, err := strconv.Atoi(number)
	if err != nil {
		return false
	}
	return (n%10+checksum(n/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
