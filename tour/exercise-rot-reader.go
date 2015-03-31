package main
import (
    "io"
    "strings"
    "os"
)


type rot13Reader struct {
    r io.Reader
}


func rot13(b byte) byte {
    var dr byte
    switch {
        case b < 65:
            dr = b
        case b <= 109:
            dr = b + 13
        case b <= 122:
            dr = b - 13
        default:
            dr = b
    }
    return dr
}

func (rt13Reader rot13Reader) Read(b []byte) (int, error) {
    n, err := rt13Reader.r.Read(b)
    if err != nil {
        return n, err
    }
    for i, c := range b {
        b[i] = rot13(c)
    }
    return n, nil
}

func main() {
    s := strings.NewReader("Lbh penpxrq gur pbqr!")
    r := rot13Reader{s}
    io.Copy(os.Stdout, &r)
}
/*

Lbh penpxrq gur pbqr
You cracked the code

76  98  104_112 101  110 112 120 114 113
89  111 117_99  114  97  99  107 101 100

+13 +13 +13_-13 +13  -13 -13 -13 -13 -13


ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm

26 - (p+13)
*/