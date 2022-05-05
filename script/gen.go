package script

func GetCard(s string) []int {
	r := len(s)
	l := 0
	v := 0
	m := map[int]int{}
	out := []int{}
	for l < r {
		switch s[l] {
		case '2', '3', '4', '5', '6', '7', '8', '9':
			v = int(s[l] - '1')
			l++
		case '1':
			l += 2
			v = 9
		case 'J':
			if s[l+1] == 'o' {
				l += 5
				v = 53
			} else {
				l++
				v = 10
			}
		case 'Q':
			l++
			v = 11
		case 'K':
			l++
			v = 12
		case 'A':
			l++
			v = 0
		case 'V':
			l += 9
			v = 52
		}
		m[v] += 1
		switch v {
		case 52, 53:
			out = append(out, (m[v]-1)*54+v)
		default:
			out = append(out, ((m[v]-1)%4)*13+v+m[v]/4*54)
		}
	}
	return out
}
