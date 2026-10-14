# chbs

Generate xkcd-style passphrases from the shell, with the same options as the web
version at [correcthorsebatterystaple.net](https://correcthorsebatterystaple.net).

```console
$ chbs
Manner-Conversation-Doubt-Replace-1
```

The wordlists are embedded in the binary and every draw comes from `crypto/rand`,
so it runs offline, starts instantly, and depends on nothing outside the standard
library.

## Install

```console
go install github.com/kenaqshal/chbs@latest
```

Needs Go 1.27.1 or newer. From a checkout, `go build` is enough.

## Usage

| Flag | Default | Meaning |
| --- | --- | --- |
| `-n` | `4` | minimum number of words |
| `-l` | `20` | minimum phrase length, separators included |
| `-s` | `-` | separator; with more than one character, one is drawn per gap |
| `-upper` | `true` | capitalize the first letter of every word |
| `-digits` | `true` | append one random digit (0-9) |
| `-dict` | none | extra wordlists on top of the primary one, comma separated: `jargon,science` |

`-n` and `-l` are floors rather than exact counts: words are drawn until both are
met, so a long `-l` produces a longer phrase. Standard Go flag syntax applies, so
boolean flags switch off with `-upper=false`.

## Examples

```console
$ chbs
Manner-Conversation-Doubt-Replace-1

$ chbs -n 6 -dict jargon,science
Practices-Lowest-Unbalanced-Sew-Agency-Colleague-1

$ chbs -s . -n 4
Behavior.Piece.Ear.Decide.1

$ chbs -upper=false -digits=false -s _ -n 4
reward_want_self_neck

$ chbs -n 2 -l 40
Pronunciation-Fair-Host-Happen-Brighten-Bush-2
```

## Entropy

Every word is drawn uniformly and independently, so it carries `log2` of the
number of entries it was drawn from, and the trailing digit adds 3.3 bits. At the
defaults that is four or more words from the primary list, so 48 bits or more.

| Lists in play | Entries | Bits per word |
| --- | --- | --- |
| `wordlist.txt` (always on) | 2,281 | 11.2 |
| `jargon.txt` | 8,803 | 13.1 |
| `science-terms.txt` | 576 | 9.2 |
| `wordlist.txt` + `jargon.txt` + `science-terms.txt` | 11,660 | 13.2 |

The last row sits below `log2(11,660) = 13.5` because the lists overlap: taken
together they hold 9,752 distinct words. Note also that `-dict` only adds entropy
against an attacker who does not know which lists you chose.

## Data

The wordlists in `data/` come from
[correcthorsebatterystaple.net](https://correcthorsebatterystaple.net), are MIT
licensed, and are copyright © John Van Der Loo. They are parsed as comma- or
newline-separated fields and embedded into the binary at build time.

## License

MIT, see [LICENSE](LICENSE).
