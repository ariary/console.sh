package stringslice

//cartesianProductPlusPlus: Perform cartesian product between a slice of string slice and a string slice. Beware: complexity -> quadratic
func CartesianProduct(list1 [][]string, list2 []string) (product [][]string) {
	var lenProduct int
	if len(list1) == 1 {
		lenProduct = len(list1) * (len(list2) - 1)
	} else {
		lenProduct = len(list1) * (len(list2))
	}
	product = make([][]string, lenProduct)
	productIndex := 0
	for i := 0; i < len(list1); i++ { //for each item of first list
		for j := 0; j < len(list2); j++ { //couple it with other
			product[productIndex] = append(product[productIndex], list1[i]...)
			product[productIndex] = append(product[productIndex], list2[j])
			productIndex++
		}
	}
	return product
}

//Return true if  specific string is within a string slice
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//ChunksString: split a string in chunk of fixed size
func ChunksString(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}
