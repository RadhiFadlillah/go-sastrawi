# Go-Sastrawi

[![GoDoc](https://godoc.org/github.com/RadhiFadlillah/go-sastrawi?status.png)](https://godoc.org/github.com/RadhiFadlillah/go-sastrawi)


Go-Sastrawi is a Go package for doing stemming in Indonesian language. It is based from [Sastrawi](https://github.com/sastrawi/sastrawi) for PHP by [Andy Librian](https://github.com/andylibrian). For more information in English, see [readme](https://github.com/RadhiFadlillah/go-sastrawi/blob/master/README.en.md).

Go-Sastrawi adalah package Go untuk melakukan _stemming_ pada bahasa Indonesia. Dikembangkan dari [Sastrawi](https://github.com/sastrawi/sastrawi) untuk PHP yang dibuat oleh [Andy Librian](https://github.com/andylibrian).

## Stemming

Dari [Wikipedia](https://en.wikipedia.org/wiki/Stemming), _stemming_ adalah proses untuk mengubah kata berimbuhan menjadi kata dasar. Contohnya :

- menahan => tahan
- pewarna => warna

## Contoh Penggunaan

Package ini memiliki dua komponen utama. Yang pertama adalah [_Tokenizer_](https://godoc.org/github.com/RadhiFadlillah/go-sastrawi#Tokenizer) untuk memecah kalimat menjadi kata-kata sekaligus menghapus simbol-simbol dan URL dalam kalimat tersebut. Yang kedua adalah [_Stemmer_](https://godoc.org/github.com/RadhiFadlillah/go-sastrawi#Stemmer) untuk mengubah kata berimbuhan menjadi kata dasar.

```go
import (
	"fmt"
	"github.com/RadhiFadlillah/go-sastrawi"
)

func main() {
	// Kalimat asal
	sentence := "Rakyat memenuhi halaman gedung untuk menyuarakan isi hatinya. Baca berita selengkapnya di http://www.kompas.com."

	// Pecah kalimat menjadi kata-kata menggunakan tokenizer
	tokenizer := sastrawi.NewTokenizer()
	words := tokenizer.Tokenize(sentence)

	// Ubah kata berimbuhan menjadi kata dasar
	stemmer := sastrawi.NewStemmer(sastrawi.DefaultDictionary)
	for _, word := range words {
		fmt.Printf("%s => %s\n", word, stemmer.Stem(word))
	}
}
```

Selain menggunakan kamus kata dasar default, user juga dapat membuat kamus kata dasar sendiri.

```go
import (
	"fmt"
	"github.com/RadhiFadlillah/go-sastrawi"
)

func main() {
	// Buat kamus baru
	dictionary := sastrawi.NewDictionary("lapar")
	dictionary.Print("")

	// Tambah kata dasar ke kamus
	dictionary.Add("ingin", "makan", "gizi", "enak", "lezat")
	dictionary.Print("")

	// Hapus kata dasar dari kamus
	dictionary.Remove("enak", "lezat")
	dictionary.Print("")

	// Gunakan kamus yang telah dibuat untuk stemming
	sentence := "Aku kelaparan dan menginginkan makanan yang bergizi."
	words := sastrawi.NewTokenizer().Tokenize(sentence)

	stemmer := sastrawi.NewStemmer(dictionary)
	for _, word := range words {
		fmt.Printf("%s => %s\n", word, stemmer.Stem(word))
	}
}
```

## Pustaka

#### Algoritma

1. Algoritma Nazief dan Adriani
2. Asian J. 2007. ___Effective Techniques for Indonesian Text Retrieval___. PhD thesis School of Computer Science and Information Technology RMIT University Australia. ([PDF](http://researchbank.rmit.edu.au/eserv/rmit:6312/Asian.pdf) dan [Amazon](https://www.amazon.com/Effective-Techniques-Indonesian-Text-Retrieval/dp/3639021649))
3. Arifin, A.Z., I.P.A.K. Mahendra dan H.T. Ciptaningtyas. 2009. ___Enhanced Confix Stripping Stemmer and Ants Algorithm for Classifying News Document in Indonesian Language___, Proceeding of International Conference on Information & Communication Technology and Systems (ICTS). ([PDF](http://personal.its.ac.id/files/pub/2623-agusza-baru%2021%20d%20VIP%20enhanced-confix-stripping-stem.pdf))
4. A. D. Tahitoe, D. Purwitasari. 2010. ___Implementasi Modifikasi Enhanced Confix Stripping Stemmer Untuk Bahasa Indonesia dengan Metode Corpus Based Stemming___, Institut Teknologi Sepuluh Nopember (ITS) – Surabaya, 60111, Indonesia. ([PDF](http://digilib.its.ac.id/public/ITS-Undergraduate-14255-paperpdf.pdf))
5. Tambahan aturan _stemming_ dari [kontributor Sastrawi](https://github.com/sastrawi/sastrawi/graphs/contributors).

#### Kamus Kata Dasar

Proses stemming oleh Sastrawi sangat bergantung pada kamus kata dasar. Sastrawi menggunakan kamus kata dasar dari [kateglo.com](http://kateglo.com) dengan sedikit perubahan.

## Lisensi

Sebagaimana [Sastrawi](https://github.com/sastrawi/sastrawi) untuk PHP, Go-Sastrawi juga disebarkan dengan lisensi [MIT](http://choosealicense.com/licenses/mit/). Untuk lisensi kamus kata dasar dari Kateglo adalah [CC-BY-NC-SA 3.0](https://github.com/ivanlanin/kateglo#lisensi-isi).

## Di Bahasa Pemrograman Lain

- [Sastrawi](https://github.com/sastrawi/sastrawi) - PHP
- [JSastrawi](https://github.com/jsastrawi/jsastrawi) - Java
- [cSastrawi](https://github.com/mohangk/c_sastrawi) - C
- [PySastrawi](https://github.com/har07/PySastrawi) - Python
