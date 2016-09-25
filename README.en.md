# Go-Sastrawi

[![GoDoc](https://godoc.org/github.com/RadhiFadlillah/go-sastrawi?status.png)](https://godoc.org/github.com/RadhiFadlillah/go-sastrawi)

Go-Sastrawi is a Go package for doing stemming in Indonesian language. It is based from [Sastrawi](https://github.com/sastrawi/sastrawi) for PHP by [Andy Librian](https://github.com/andylibrian).

## Stemming

From [Wikipedia](https://en.wikipedia.org/wiki/Stemming), stemming is the process of reducing inflected (or sometimes derived) words to their word stem, base or root form. For example :

- menahan => tahan
- pewarna => warna

## Usage Examples

This package has two main components. First is [Tokenizer](https://godoc.org/github.com/RadhiFadlillah/go-sastrawi#Tokenizer) that used for splitting sentence into words and remove symbols and URLs from the sentence. The second is [Stemmer](https://godoc.org/github.com/RadhiFadlillah/go-sastrawi#Stemmer) that used for reducing inflected word to its root form.

```go
import (
	"fmt"
	"github.com/RadhiFadlillah/go-sastrawi"
)

func main() {
	// Original sentence
	sentence := "Rakyat memenuhi halaman gedung untuk menyuarakan isi hatinya. Baca berita selengkapnya di http://www.kompas.com."

	// Split sentence into words using tokenizer
	tokenizer := sastrawi.NewTokenizer()
	words := tokenizer.Tokenize(sentence)

	// Reduce inflected words to its root form
	stemmer := sastrawi.NewStemmer(sastrawi.DefaultDictionary)
	for _, word := range words {
		fmt.Printf("%s => %s\n", word, stemmer.Stem(word))
	}
}
```

Beside using the default dictionary, you can also create your own root words dictionary.

```go
import (
	"fmt"
	"github.com/RadhiFadlillah/go-sastrawi"
)

func main() {
	// Create new dictionary
	dictionary := sastrawi.NewDictionary("lapar")
	dictionary.Print("")

	// Add new words to dictionary
	dictionary.Add("ingin", "makan", "gizi", "enak", "lezat")
	dictionary.Print("")

	// Remove some words from dictionary
	dictionary.Remove("enak", "lezat")
	dictionary.Print("")

	// Use your new dictionary for stemming
	sentence := "Aku kelaparan dan menginginkan makanan yang bergizi."
	words := sastrawi.NewTokenizer().Tokenize(sentence)

	stemmer := sastrawi.NewStemmer(dictionary)
	for _, word := range words {
		fmt.Printf("%s => %s\n", word, stemmer.Stem(word))
	}
}
```

## Resource

#### Algorithm

1. Nazief and Adriani Algorith
2. Asian J. 2007. ___Effective Techniques for Indonesian Text Retrieval___. PhD thesis School of Computer Science and Information Technology RMIT University Australia. ([PDF](http://researchbank.rmit.edu.au/eserv/rmit:6312/Asian.pdf) and [Amazon](https://www.amazon.com/Effective-Techniques-Indonesian-Text-Retrieval/dp/3639021649))
3. Arifin, A.Z., I.P.A.K. Mahendra dan H.T. Ciptaningtyas. 2009. ___Enhanced Confix Stripping Stemmer and Ants Algorithm for Classifying News Document in Indonesian Language___, Proceeding of International Conference on Information & Communication Technology and Systems (ICTS). ([PDF](http://personal.its.ac.id/files/pub/2623-agusza-baru%2021%20d%20VIP%20enhanced-confix-stripping-stem.pdf))
4. A. D. Tahitoe, D. Purwitasari. 2010. ___Implementasi Modifikasi Enhanced Confix Stripping Stemmer Untuk Bahasa Indonesia dengan Metode Corpus Based Stemming___, Institut Teknologi Sepuluh Nopember (ITS) – Surabaya, 60111, Indonesia. ([PDF](http://digilib.its.ac.id/public/ITS-Undergraduate-14255-paperpdf.pdf))
5. Additional stemming rules from [Sastrawi's contributors](https://github.com/sastrawi/sastrawi/graphs/contributors).

#### Root Words Dictionary

Stemming process by this package is depends heavily on the root words dictionary. Sastrawi use root words dictionary from [kateglo.com](http://kateglo.com) with some changes.

## License

As [Sastrawi](https://github.com/sastrawi/sastrawi) for PHP, Go-Sastrawi is also distributed using [MIT](http://choosealicense.com/licenses/mit/) license. Root words dictionary is distributed by Kateglo using [CC-BY-NC-SA 3.0](https://github.com/ivanlanin/kateglo#lisensi-isi) license.

## Sastrawi in Other Language

- [Sastrawi](https://github.com/sastrawi/sastrawi) - PHP
- [JSastrawi](https://github.com/jsastrawi/jsastrawi) - Java
- [cSastrawi](https://github.com/mohangk/c_sastrawi) - C
- [PySastrawi](https://github.com/har07/PySastrawi) - Python
