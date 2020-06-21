package sastrawi

import (
	"regexp"
)

var (
	// Regex for tokenizer
	rxURL       = regexp.MustCompile(`(?i)(www\.|https?|s?ftp)\S+`)
	rxEmail     = regexp.MustCompile(`(?i)\S+@\S+`)
	rxTwitter   = regexp.MustCompile(`(?i)(@|#)\S+`)
	rxEscapeStr = regexp.MustCompile(`(?i)&.*;`)
	rxSymbol    = regexp.MustCompile(`(?i)[^a-z\s]`)

	// Regex for stemmer
	rxPrefixFirst = regexp.MustCompile(`^(be.+lah|be.+an|me.+i|di.+i|pe.+i|ter.+i)$`)
	rxParticle    = regexp.MustCompile(`-*(lah|kah|tah|pun)$`)
	rxPossesive   = regexp.MustCompile(`-*(ku|mu|nya)$`)
	rxSuffix      = regexp.MustCompile(`-*(is|isme|isasi|i|kan|an)$`)

	rxPrefixMe1  = regexp.MustCompile(`^me([lrwy][aiueo].*)$`) // me{l|r|w|y}V => me-{l|r|w|y}V
	rxPrefixMe2  = regexp.MustCompile(`^mem([bfv].*)$`)        // mem{b|f|v} => mem-{b|f|v}
	rxPrefixMe3  = regexp.MustCompile(`^mem(pe.*)$`)           // mempe => mem-pe
	rxPrefixMe4  = regexp.MustCompile(`^mem(r?[aiueo].*)$`)    // mem{rV|V} => mem-{rV|V} OR me-p{rV|V}
	rxPrefixMe5  = regexp.MustCompile(`^men([cdjstz].*)$`)     // men{c|d|j|s|t|z} => men-{c|d|j|s|t|z}
	rxPrefixMe6  = regexp.MustCompile(`^men([aiueo].*)$`)      // menV => me-nV OR me-tV
	rxPrefixMe7  = regexp.MustCompile(`^meng([ghqk].*)$`)      // meng{g|h|q|k} => meng-{g|h|q|k}
	rxPrefixMe8  = regexp.MustCompile(`^meng(([aiueo])(.*))$`) // mengV => meng-V OR meng-kV OR me-ngV OR mengV- where V = 'e'
	rxPrefixMe9  = regexp.MustCompile(`^meny(([aiueo])(.*))$`) // menyV => meny-sV OR me-nyV to stem menyala
	rxPrefixMe10 = regexp.MustCompile(`^mem(p[^e].*)$`)        // mempV => mem-pA where A != 'e'

	rxPrefixPe1  = regexp.MustCompile(`^pe([wy][aiueo].*)$`)              // pe{w|y}V => pe-{w|y}V
	rxPrefixPe2  = regexp.MustCompile(`^per([aiueo].*)$`)                 // perV => per-V OR pe-rV
	rxPrefixPe3  = regexp.MustCompile(`^per([^aiueor][a-z][^e].*)$`)      // perCAP => per-CAP where C != 'r' and P != 'er'
	rxPrefixPe4  = regexp.MustCompile(`^per([^aiueor][a-z]er[aiueo].*)$`) // perCAerV => per-CAerV where C != 'r'
	rxPrefixPe5  = regexp.MustCompile(`^pem([bfv].*)$`)                   // pem{b|f|v} => pem-{b|f|v}
	rxPrefixPe6  = regexp.MustCompile(`^pem(r?[aiueo].*)$`)               // pem{rV|V} => pe-m{rV|V} OR pe-p{rV|V}
	rxPrefixPe7  = regexp.MustCompile(`^pen([cdjstz].*)$`)                // pen{c|d|j|s|t|z} => pen-{c|d|j|s|t|z}
	rxPrefixPe8  = regexp.MustCompile(`^pen([aiueo].*)$`)                 // penV => pe-nV OR pe-tV
	rxPrefixPe9  = regexp.MustCompile(`^peng([^aiueo].*)$`)               // pengC => peng-C
	rxPrefixPe10 = regexp.MustCompile(`^peng(([aiueo])(.*))$`)            // pengV => peng-V OR peng-kV OR pengV- where V = 'e'
	rxPrefixPe11 = regexp.MustCompile(`^peny([aiueo].*)$`)                // penyV => peny-sV OR pe-nyV
	rxPrefixPe12 = regexp.MustCompile(`^pe(l[aiueo].*)$`)                 // pelV => pe-lV OR pel-V for pelajar
	rxPrefixPe13 = regexp.MustCompile(`^pe[^aiueorwylmn](er[aiueo].*)$`)  // peCerV => per-erV where C != {r|w|y|l|m|n}
	rxPrefixPe14 = regexp.MustCompile(`^pe([^aiueorwylmn][^e].*)$`)       // peCP => pe-CP where C != {r|w|y|l|m|n} and P != 'er'
	rxPrefixPe15 = regexp.MustCompile(`^pe([^aiueorwylmn]er[^aiueo].*)$`) // peC1erC2 => pe-C1erC2 where C1 != {r|w|y|l|m|n}

	rxPrefixBe1 = regexp.MustCompile(`^ber([aiueo].*)$`)                 // berV => ber-V || be-rV
	rxPrefixBe2 = regexp.MustCompile(`^ber([^aiueor][a-z][^e].*)$`)      // berCAP => ber-CAP where C != 'r' and P != 'er'
	rxPrefixBe3 = regexp.MustCompile(`^ber([^aiueor][a-z]er[aiueo].*)$`) // berCAerV => ber-CAerV where C != 'r'
	rxPrefixBe4 = regexp.MustCompile(`^bel(ajar)$`)                      // belajar => bel-ajar
	rxPrefixBe5 = regexp.MustCompile(`^be([^aiueorl]er[^aiueo].*)$`)     // beC1erC2 => be-C1erC2 where C1 != {'r'|'l'}

	rxPrefixTe1 = regexp.MustCompile(`^ter([aiueo].*)$`)             // terV => ter-V OR te-rV
	rxPrefixTe2 = regexp.MustCompile(`^ter([^aiueor]er[aiueo].*)$`)  // terCerV => ter-CerV where C != 'r'
	rxPrefixTe3 = regexp.MustCompile(`^ter([^aiueor][^e].*)$`)       // terCP => ter-CP where C != 'r' and P != 'er'
	rxPrefixTe4 = regexp.MustCompile(`^te([^aiueor]er[^aiueo].*)$`)  // teC1erC2 => te-C1erC2 where C1 != 'r'
	rxPrefixTe5 = regexp.MustCompile(`^ter([^aiueor]er[^aiueo].*)$`) // terC1erC2 => ter-C1erC2 where C1 != 'r'

	rxInfix1 = regexp.MustCompile(`^(([^aiueo])e[rlm])([aiueo].*)$`) // Ce{r|l|m}V => Ce{r|l|m}V OR CV
	rxInfix2 = regexp.MustCompile(`^(([^aiueo])in)([aiueo].*)$`)     // CinV => CinV OR CV
)
