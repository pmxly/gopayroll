package taskspool

func Add(args ...int64) (int64, error) {
	/*sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil*/
	num := fibonacci(45)
	return num, nil

}

func fibonacci(num int64) int64 {
	if num < 2 {
		return 1
	}
	return fibonacci(num-1) + fibonacci(num-2)
}