package stringutil

import "testing"

func TestLength(t *testing.T) {
    cases := []struct {
        in string
        want int
    }{
        {"Teste", 5},
        {"123456", 6},
        {"", 0},
    }
    for _, c := range cases {
        got := Length(c.in)
        if got != c.want {
            t.Errorf("Length(%q) == %d, but want %d", c.in, got, c.want)
        }
    }
}