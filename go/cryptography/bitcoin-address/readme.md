# Cryptography - how to get a Bitcoin Address from public key?

Recently topics like Bitcoin, blockchain, cryptocurrency are widely popular not
only in developers circles but generally. As I am very interested and excited
about cryptography and its practical application I decided to write an article
with practical examples about one of the common operations in Bitcoin lifecycle.

Rather than using imaginary pseudo-language I used my favourite Back-End language:
Golang or shorter go.

When initially setting up, a Bitcoin client software creates new public/private
ECDSA key pair which is used to sign and later verify validity of transaction
claim. Same key is used to generate a Bitcoin Address for the owner. So how do
we get a Bitcoin Address from key pair?

From Bitcoin official wiki:

> A Bitcoin address is only a hash, so the sender can't provide a full public key
> in scriptPubKey. When redeeming coins that have been sent to a Bitcoin address,
> the recipient provides both the signature and the public key. The script
> verifies that the provided public key does hash to the hash in scriptPubKey,
> and then it also checks the signature against the public key.
(https://en.bitcoin.it/wiki/Transaction#Pay-to-PubkeyHash)

To calculate a bitcoin address we need public part of ECDSA (shorter acronym
for Elliptic Curve Digital Signature Algorithm) key pair. It consist of two
very long (32 bytes each) numbers between 0 and _n_ where _n_ is the curve order.
Those numbers are coordinates that we get by calculating elliptic curve from
private key. In order to display them, those 64 bytes of data is converted to
hexadecimal format. It takes 130 hex characters at 4 bits per character to
display the full key.

Let's assume we are building an api that will accept public key from client
and our task is to convert that public key into bitcoin format address. You would
receive from client two strings that represent X and Y coordinates in hexadecimal
string:

```go

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

```

If you run above code written in golang it will output string representation
for both coordinates and their size in bytes as well. To get full public key we
concatenate those two coordinates into one string.

Output:
```sh
2017/09/03 07:08:45 === X (132 bytes)  ===
2017/09/03 07:08:45 0144760C30109592B556D2AB2FA13C4383F3A3897E21015EE67650ADCD74289B648485CA966C1748480551246713BB5AB4CEE068CFDC3D9E31BD3502C486448543AA
2017/09/03 07:08:45 === End of X ===
2017/09/03 07:08:45 === Y (132 bytes)  ===
2017/09/03 07:08:45 FD874316135A987D7DABE7CE45A55ADAA296D528B023A133679B431D945DBDBDB0A411E47E584FF1356B6D37A01CDFB1C280178C26C38FA7976252030BA331E635AA
2017/09/03 07:08:45 === End of Y ===
2017/09/03 07:08:45 === XY (264 bytes) ===
2017/09/03 07:08:45 0144760C30109592B556D2AB2FA13C4383F3A3897E21015EE67650ADCD74289B648485CA966C1748480551246713BB5AB4CEE068CFDC3D9E31BD3502C486448543AAFD874316135A987D7DABE7CE45A55ADAA296D528B023A133679B431D945DBDBDB0A411E47E584FF1356B6D37A01CDFB1C280178C26C38FA7976252030BA331E635AA
2017/09/03 07:08:45 === End of XY ===
```

According to Bitcoin standard we append 1 byte 0x04 or **04** to the left side
of our string which now looks like this:

```go

    	pub := "04" + xy

	pubsize := strconv.Itoa(binary.Size([]byte(pub)))
	log.Println("=== pub key (" + pubsize + " bytes) ===")
	log.Println(pub)
	log.Println("=== End of pub key  ===")

```

Output:
```sh
2017/09/03 07:08:45 === pub key (266 bytes) ===
2017/09/03 07:08:45 040144760C30109592B556D2AB2FA13C4383F3A3897E21015EE67650ADCD74289B648485CA966C1748480551246713BB5AB4CEE068CFDC3D9E31BD3502C486448543AAFD874316135A987D7DABE7CE45A55ADAA296D528B023A133679B431D945DBDBDB0A411E47E584FF1356B6D37A01CDFB1C280178C26C38FA7976252030BA331E635AA
2017/09/03 07:08:45 === End of pub key  ===
```

The prefix 04 is used to distinguish uncompressed public keys from compressed
public keys that begin with a 02 or 03.

Now that we have full string in a format 04 + X coordinate  + Y coordinate we
can compress it or  reduce its length by converting it to byte array and hashing
the result with sha256 cryptographic function (designed by NSA):

```go
	op1c := sha256.Sum256([]byte(pub))
	op1 := op1c[:]

	op1size := strconv.Itoa(binary.Size(op1))
	log.Println("=== op1 SHA256 (" + op1size + " bytes) ===")
	log.Println(op1)
	log.Println("=== End of op1 SHA256 ===")

```

Output:
```sh
2017/09/03 07:08:45 === op1 SHA256 (32 bytes) ===
2017/09/03 07:08:45 [51 253 48 121 109 252 192 231 197 225 240 28 83 213 10 101 86 137 250 207 132 196 252 160 137 203 10 70 174 9 203 84]
2017/09/03 07:08:45 === End of op1 SHA256 ===
```

During this process we will use sha256 function three times in order to improve
security and privacy. Result of this operation is a 32 bytes in a string which
we then hash again but this time with _ripemd160_ function:

```go
	op2c := ripemd160.New()
	op2c.Write(op1)
	op2 := op2c.Sum(nil)

	op2size := strconv.Itoa(binary.Size(op2))
	log.Println("=== op2 ripemd160 (" + op2size + " bytes) ===")
	log.Println(op2)
	log.Println("=== End of op2 ripemd160 ===")

```

Output:
```sh
2017/09/03 07:08:45 === op2 ripemd160 (20 bytes) ===
2017/09/03 07:08:45 [243 145 90 252 135 153 181 107 75 211 242 24 56 225 246 132 242 8 74 34]
2017/09/03 07:08:45 === End of op2 ripemd160 ===
```

Now we add to the left side something called address prefix, also referred to as
version id. It is simply agreed prefix to help fast and easily validate
if certain string is really Bitcoin address. This also tells to which network
address belongs. Think of it as a country calling code in a phone number.
Most used is the _0x00_ which translates in decimal value of _0_.

```go
	op3 := append([]byte{0x00}, op2...)

	op3size := strconv.Itoa(binary.Size(op3))
	log.Println("=== op3 apend ver.  byte (" + op3size + " bytes) ===")
	log.Println(op3)
	log.Println("=== End of op3 append ver.  byte ===")

```

Output:
```sh
2017/09/03 07:08:45 === op3 apend ver.  byte (21 bytes) ===
2017/09/03 07:08:45 [0 243 145 90 252 135 153 181 107 75 211 242 24 56 225 246 132 242 8 74 34]
2017/09/03 07:08:45 === End of op3 append ver.  byte ===
```

This will be base of our future Bitcoin Address, but for now we will hash the
result with _sha256_ two times:


```go
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

```

Result is a hash in a 32 bytes long byte array.

Output:
```sh
2017/09/03 07:08:45 === op4 SHA256 (32 bytes) ===
2017/09/03 07:08:45 [192 224 168 49 151 40 227 32 236 85 55 233 79 255 225 44 141 155 65 149 206 38 116 213 201 132 197 172 43 128 136 47]
2017/09/03 07:08:45 === End of op4 SHA256 ===
2017/09/03 07:08:45 === op5 SHA256 (32 bytes) ===
2017/09/03 07:08:45 [196 242 220 247 252 144 46 153 108 48 51 182 195 25 59 155 116 184 118 28 168 143 43 163 104 165 48 106 39 94 126 237]
2017/09/03 07:08:45 === End of op5 SHA256 ===
```

We will need only first four bytes which is referred to as address checksum:


```go
	op6 := op5c[0:4]

	op6size := strconv.Itoa(binary.Size(op6))
	log.Println("=== op6 address checksum (" + op6size + " bytes) ===")
	log.Println(op6)
	log.Println("=== End of op6 address checksum ===")

```

Output:
```sh
2017/09/03 07:08:45 === op6 address checksum (4 bytes) ===
2017/09/03 07:08:45 [196 242 220 247]
2017/09/03 07:08:45 === End of op6 address checksum ===
```

Remember before we mentioned that we have a base of our Bitcoin Address? We will
take it now and we will add four bytes from above to the end of byte array of base:


```go
	op7 := append(op3, op6...)

	op7size := strconv.Itoa(binary.Size(op7))
	log.Println("=== op7 binary address (" + op7size + " bytes) ===")
	log.Println(op7)
	log.Println("=== End of op7 binary address ===")

```

Output:
```sh
2017/09/03 07:08:45 === op7 binary address (25 bytes) ===
2017/09/03 07:08:45 [0 243 145 90 252 135 153 181 107 75 211 242 24 56 225 246 132 242 8 74 34 196 242 220 247]
2017/09/03 07:08:45 === End of op7 binary address ===
```

And this is our new 25-byte binary Bitcoin Address. We will at the end convert
it into a _base58_ string using _Base58Check_ encoding. Base58Check encoding is
used for encoding byte arrays of address into human-typable strings. The
original Bitcoin client source code explains the reasoning behind base58 encoding:

```h
// Why base-58 instead of standard base-64 encoding?
// - Don't want 0OIl characters that look the same in some fonts and
//      could be used to create visually identical looking account numbers.
// - A string with non-alphanumeric characters is not as easily accepted as an account number.
// - E-mail usually won't line-break if there's no punctuation to break at.
// - Doubleclicking selects the whole number as one word if it's all alphanumeric.
```

And this is how we do it:


```go

	op8 := base58.Encode(op7)

	op8size := strconv.Itoa(binary.Size([]byte(op8)))
	log.Println("=== op8 base 58 address (" + op8size + " bytes) ===")
	log.Println(op8)
	log.Println("=== End of op8 base 58 address ===")

```

Output:
```sh
2017/09/03 07:08:45 === op8 base 58 address (34 bytes) ===
2017/09/03 07:08:45 1PCsL5zmQBS1nKBESPW9VozNBDsUUE4YHL
2017/09/03 07:08:45 === End of op8 base 58 address ===
```

Congratulations! Now we should have our new shiny address: **1PCsL5zmQBS1nKBESPW9VozNBDsUUE4YHL**

This whole process is summarised in schema that can be found on a Bitcoin official
wiki documentation:

![Bitcoin Address schema](https://en.bitcoin.it/w/images/en/9/9b/PubKeyToAddr.png "Bitcoin Address schema")

If you followed along with code you should get exactly same result as performing
above operations on same public key should always lead to the same output.

[My Twitter account](https://twitter.com/vedran_s)

[vsrc.pro](http://www.vsrc.pro/)
