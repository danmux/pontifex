package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	cCount = 54
	ja     = byte(53)
	jb     = byte(54)
	maxKey = 26
)

var verbose = false

func main() {

	var keyName, inName, outName string
	flag.StringVar(&keyName, "k", "shared.key", "the shared key to use")
	flag.StringVar(&inName, "i", "in.txt", "the input")
	flag.StringVar(&outName, "o", "out.txt", "the output")
	op := flag.String("op", "enc", "[enc, dec, key, dkey] key - takes a dec of cards input and makes a key, dkey, makes a key from a default deck")
	flag.BoolVar(&verbose, "v", false, "verbose prints the deck at each stage")

	orgHelp := flag.Usage
	flag.Usage = func() {
		orgHelp()
		setupMaps(true)
	}

	flag.Parse()

	setupMaps(false)

	in, err := ioutil.ReadFile(inName)
	if err != nil {
		log.Fatal(err)
	}

	switch *op {
	case "enc":
		key := readKey(keyName)
		in = cleanInput(in)
		o := encrypt(key, in)
		err = writeChars(o, outName, true)
	case "dec":
		key := readKey(keyName)
		in = cleanInput(in)
		o := decrypt(key, in)
		err = writeChars(o, outName, false)
	case "key":
		cards, err := ioutil.ReadFile(inName)
		if err != nil {
			break
		}
		key := cardsToKey(cards)
		err = writeKey(key, keyName)
	case "dkey":
		err = defaultKey(keyName)
	default:
		log.Fatal(*op, "unsupported")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func decrypt(k key, in []byte) []byte {
	sz := len(in)
	ks := k.stream(sz)
	o := make([]byte, sz)
	for i, v := range in {
		if v < ks[i] {
			v += maxKey
		}
		o[i] = v - ks[i]
	}
	return o
}

func encrypt(k key, in []byte) []byte {
	sz := len(in)
	ks := k.stream(sz)
	o := make([]byte, sz)
	for i, v := range in {
		o[i] = ks[i] + v
		if o[i] > maxKey {
			o[i] = o[i] - maxKey
		}
	}
	return o
}

// defaultKey just generates a default key and writes it in base64
func defaultKey(fName string) error {
	k := cardsToKey([]byte(`
						1H 2H 3H 4H 5H 6H 7H 8H 9H 10H JH QH KH *A
						1S 2S 3S 4S 5S 6S 7S 8S 9S 10S JS QS KS
						1D 2D 3D 4D 5D 6D 7D 8D 9D 10D JD QD KD *B
						1C 2C 3C 4C 5C 6C 7C 8C 9C 10C JC QC KC
						`))
	return writeKey(k, fName)
}

// readKey reads a base64 encoded key
func readKey(fName string) key {
	keyB64, err := ioutil.ReadFile(fName)
	if err != nil {
		log.Fatal(err)
	}

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(keyB64)))
	if _, err := base64.StdEncoding.Decode(dst, keyB64); err != nil {
		log.Fatal(err)
	}

	return key(dst)
}

// writeKey writes a key to base 64
func writeKey(k key, fName string) error {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(k)))
	base64.StdEncoding.Encode(dst, k)
	return ioutil.WriteFile(fName, dst, 0644)
}

// writeChars writes the bytes out mapped to chars from 'A' onwards
// if block is true then writes 4 blocks of 5 chars per line - like old school encrypted stuff
func writeChars(d []byte, fName string, block bool) error {
	f, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	for i, v := range d {
		if err := w.WriteByte('A' + v - 1); err != nil {
			return err
		}
		if !block {
			continue
		}
		// block version
		if i%5 == 4 {
			if err := w.WriteByte(' '); err != nil {
				return err
			}
		}
		if i%20 == 19 {
			if err := w.WriteByte('\n'); err != nil {
				return err
			}
		}
	}
	if err := w.WriteByte('\n'); err != nil {
		return err
	}
	return w.Flush()
}

// cleanInput takes the input bytes, removes whitespace
// maps upper case chars to values where 'A' = 1
func cleanInput(in []byte) []byte {
	in = bytes.Join(bytes.Fields(in), nil)
	for i, v := range in {
		in[i] = v - 'A' + 1
	}
	pads := 4 - (len(in)-1)%5
	for i := 0; i < pads; i++ {
		in = append(in, 9+byte(i)) // pads with IJKL
	}

	return in
}

// stream generates a key stream of count size
func (k *key) stream(count int) []byte {
	r := make([]byte, count)
	tot := 0
	for tot < count {
		cardVal := k.oneRound()
		if cardVal == 255 {
			continue
		}
		r[tot] = cardVal
		tot++
	}

	return r
}

// oneRound generates one key round
func (k *key) oneRound() byte {

	sz := len(*k)

	if verbose {
		fmt.Println("\nstart")
		k.log()
	}

	// move A joker down 1
	k.move(ja, 1)

	if verbose {
		fmt.Println("\nmoved ja")
		k.log()
	}

	// move B joker down 2
	k.move(jb, 2)

	if verbose {
		fmt.Println("\nmoved jb")
		k.log()
	}

	// find the top most joker
	top := k.findVal(ja)
	bot := k.findVal(jb)
	if top >= bot {
		b := bot
		bot = top
		top = b
	}

	if verbose {
		fmt.Printf("top:%d, bot:%d\n", top, bot)
	}

	// tripple cut into new key
	newKey := make([]byte, sz)

	p := 0
	for i := bot + 1; i < sz; i++ {
		newKey[p] = (*k)[i]
		p++
	}
	for i := top; i <= bot; i++ {
		newKey[p] = (*k)[i]
		p++
	}
	for i := 0; i < top; i++ {
		newKey[p] = (*k)[i]
		p++
	}

	if verbose {
		fmt.Println("\ntripple cut")
		key(newKey).log()
	}

	// get the value of the bottom of the deck for the count cut
	// handling jokers
	count := enforceMaxCardVal(newKey[sz-1])

	if verbose {
		fmt.Println("\ncount cut:", count)
	}

	// the count cut copies from the newKey above back into the original key
	p = 0
	// from the count onwards becomes the top
	for i := int(count); i < sz-1; i++ {
		(*k)[p] = newKey[i]
		p++
	}
	// and the top becomes the bottom - up to the last card which stays the same
	for i := 0; i < int(count); i++ {
		(*k)[p] = newKey[i]
		p++
	}
	// put the last card back in place
	(*k)[sz-1] = newKey[sz-1]

	if verbose {
		fmt.Println("\nafter count cut")
		k.log()
	}

	// get the value of the top card to count down
	tc := (*k)[0]
	if verbose {
		fmt.Println("\ntop card", valToCard[tc], tc)
	}
	// handle jokers
	tc = enforceMaxCardVal(tc)

	// get the value of this card
	val := (*k)[tc]

	// if we land on a joker - ignore this round
	if val > 53 {
		return 255
	}

	// only generate up to our max modulus value
	if val > maxKey {
		val -= maxKey
	}
	return val
}

// enforceMaxCardVal simply makes sure both jokers are counted as 53
func enforceMaxCardVal(in byte) byte {
	if in > 52 {
		return 53 // both jokers value 53
	}
	return in
}

// cardsToKey converts the string based deck format to a key
func cardsToKey(cards []byte) key {
	r := key{}
	cs := bytes.Fields(cards)
	for _, c := range cs {
		r = append(r, cardToVal[string(c)])
	}
	return r
}

// key models the deck based key
type key []byte

// findVal returns the index of the card with value val
func (k key) findVal(val byte) int {
	for i, v := range k {
		if val == v {
			return i
		}
	}
	return -1
}

// move moves one card forward in the key deck
func (k key) move(val byte, off int) {
	sz := len(k)
	pos := k.findVal(val)

	// if this will wrap around then becomes an insert
	if pos+off >= sz {
		// slide the deck down one between the old and new positions
		insert := (pos + off) % sz
		for i := pos; i > insert; i-- {
			k[i] = k[i-1]
		}
		// and set the new position to val
		k[insert+1] = val
		return
	}

	// otherwise move everything up one and set the new position to val
	for i := 0; i < off; i++ {
		k[pos+i] = k[pos+i+1]
	}
	k[pos+off] = val
	return
}

// log prints the deck out
func (k key) log() {
	for i, v := range k {
		if i%13 == 0 {
			println()
		}
		print(valToCard[v], " ")
	}
	println()
}

// some maps to go from and to card string to their representative value
var (
	cardToVal = map[string]byte{}
	valToCard = map[byte]string{}
)

// setupMaps sets up the maps from card strings to values and visa versa
func setupMaps(show bool) {
	if show {
		fmt.Fprintln(os.Stderr, "    Card Values:")
	}
	val := byte(1)
	// bridge clubs, diamonds, hearts, and spades
	for _, suit := range []string{"C", "D", "H", "S"} {
		for _, card := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"} {
			c := card + suit
			if show {
				fmt.Fprintf(os.Stderr, "    %s\t%d\n", c, int(val))
			}
			cardToVal[c] = val
			valToCard[val] = c
			val++
		}
	}
	j := "*A"
	cardToVal[j] = ja
	valToCard[ja] = j

	j = "*B"
	cardToVal[j] = jb
	valToCard[jb] = j
}
