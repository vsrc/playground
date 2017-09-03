package main

import (
	"crypto/sha256"
	"encoding/binary"
	"log"
	"strconv"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func main() {

	log.Println("Hello World!")

	x := "0144760C30109592B556D2AB2FA13C4383F3A3897E21015EE67650ADCD74289B648485CA966C1748480551246713BB5AB4CEE068CFDC3D9E31BD3502C486448543AA"

	xsize := strconv.Itoa(binary.Size([]byte(x)))
	log.Println("=== X (" + xsize + " bytes)  ===")
	log.Println(x)
	log.Println("=== End of X ===")

	y := "FD874316135A987D7DABE7CE45A55ADAA296D528B023A133679B431D945DBDBDB0A411E47E584FF1356B6D37A01CDFB1C280178C26C38FA7976252030BA331E635AA"

	ysize := strconv.Itoa(binary.Size([]byte(y)))
	log.Println("=== Y (" + ysize + " bytes)  ===")
	log.Println(y)
	log.Println("=== End of Y ===")

	xy := x + y

	xysize := strconv.Itoa(binary.Size([]byte(xy)))
	log.Println("=== XY (" + xysize + " bytes) ===")
	log.Println(xy)
	log.Println("=== End of XY ===")

	pub := "04" + xy

	pubsize := strconv.Itoa(binary.Size([]byte(pub)))
	log.Println("=== pub key (" + pubsize + " bytes) ===")
	log.Println(pub)
	log.Println("=== End of pub key  ===")

	op1c := sha256.Sum256([]byte(pub))
	op1 := op1c[:]

	op1size := strconv.Itoa(binary.Size(op1))
	log.Println("=== op1 SHA256 (" + op1size + " bytes) ===")
	log.Println(op1)
	log.Println("=== End of op1 SHA256 ===")

	op2c := ripemd160.New()
	op2c.Write(op1)
	op2 := op2c.Sum(nil)

	op2size := strconv.Itoa(binary.Size(op2))
	log.Println("=== op2 ripemd160 (" + op2size + " bytes) ===")
	log.Println(op2)
	log.Println("=== End of op2 ripemd160 ===")

	op3 := append([]byte{0x00}, op2...)

	op3size := strconv.Itoa(binary.Size(op3))
	log.Println("=== op3 apend ver.  byte (" + op3size + " bytes) ===")
	log.Println(op3)
	log.Println("=== End of op3 append ver.  byte ===")

	op4c := sha256.Sum256(op3)
	op4 := op4c[:]

	op4size := strconv.Itoa(binary.Size(op4))
	log.Println("=== op4 SHA256 (" + op4size + " bytes) ===")
	log.Println(op4)
	log.Println("=== End of op4 SHA256 ===")

	op5c := sha256.Sum256(op4)
	op5 := op5c[:]

	op5size := strconv.Itoa(binary.Size(op5))
	log.Println("=== op5 SHA256 (" + op5size + " bytes) ===")
	log.Println(op5)
	log.Println("=== End of op5 SHA256 ===")

	op6 := op5c[0:4]

	op6size := strconv.Itoa(binary.Size(op6))
	log.Println("=== op6 address checksum (" + op6size + " bytes) ===")
	log.Println(op6)
	log.Println("=== End of op6 address checksum ===")

	op7 := append(op3, op6...)

	op7size := strconv.Itoa(binary.Size(op7))
	log.Println("=== op7 binary address (" + op7size + " bytes) ===")
	log.Println(op7)
	log.Println("=== End of op7 binary address ===")

	op8 := base58.Encode(op7)

	op8size := strconv.Itoa(binary.Size([]byte(op8)))
	log.Println("=== op8 base 58 address (" + op8size + " bytes) ===")
	log.Println(op8)
	log.Println("=== End of op8 base 58 address ===")

}
