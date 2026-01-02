package service

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncodeBase62(n uint64) string {
	if n == 0 {
		return "0"
	}

	buf := make([]byte, 0, 8)
	for n > 0 {
		r := n % 62
		buf = append(buf, base62Chars[r])
		n = n / 62
	}

	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}
