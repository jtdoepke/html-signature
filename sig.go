package htmlsig

import (
	"io"

	"github.com/ekzhu/minhash-lsh"
	"golang.org/x/net/html"
)

type lzwEntryLink interface{} // lzwEntry | nil

type lzwEntry struct {
	previous lzwEntryLink
	token    string
}

type HTMLSignature struct {
	SignatureSize int

	StructureFingerprint []byte
	previousLZWEntry     lzwEntryLink
	tagLZWDict           map[lzwEntry]byte

	TextFingerprint []uint64
	minhash         *minhashlsh.Minhash
}

func (sig *HTMLSignature) PushTag(tag []byte) {
	newLZWEntry := lzwEntry{
		previous: sig.previousLZWEntry,
		token:    string(tag),
	}
	if _, ok := sig.tagLZWDict[newLZWEntry]; ok {
		sig.previousLZWEntry = newLZWEntry
	} else {
		if len(sig.tagLZWDict) < sig.SignatureSize {
			sig.tagLZWDict[newLZWEntry] = byte(len(sig.tagLZWDict)) + 1
		}
		if newLZWEntry.previous == nil {
			sig.StructureFingerprint = append(sig.StructureFingerprint, 0)
		} else {
			sig.StructureFingerprint = append(sig.StructureFingerprint, sig.tagLZWDict[newLZWEntry.previous.(lzwEntry)])
		}
		sig.previousLZWEntry = nil
	}
}

func (sig *HTMLSignature) PushText(text []byte) {
	sig.minhash.Push(text)
}

func (sig *HTMLSignature) Close() {
	sig.previousLZWEntry = nil
	sig.tagLZWDict = nil

	sig.TextFingerprint = sig.minhash.Signature()
	sig.minhash = nil
}

func NewHTMLSignature(signatureSize int) HTMLSignature {
	return HTMLSignature{
		StructureFingerprint: make([]byte, 0),
		SignatureSize:        signatureSize,
		tagLZWDict:           make(map[lzwEntry]byte, signatureSize),
		minhash:              minhashlsh.NewMinhash(0, signatureSize),
	}
}

type Token struct {
	Type html.TokenType
	Text []byte
}

func (signature *HTMLSignature) FromReader(html_reader io.Reader) {
	tokenizer := html.NewTokenizer(html_reader)

	c := make(chan Token)
	go func() {
		defer close(c)
		for tokenType := tokenizer.Next(); tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
			switch tokenType {
			case html.StartTagToken, html.EndTagToken:
				t, _ := tokenizer.TagName()
				c <- Token{
					Type: tokenType,
					Text: t,
				}
			case html.TextToken:
				t := tokenizer.Text()
				c <- Token{
					Type: tokenType,
					Text: t,
				}
			}
		}
	}()

	for token := range c {
		switch token.Type {
		case html.StartTagToken, html.EndTagToken:
			signature.PushTag(token.Text)
		case html.TextToken:
			signature.PushText(token.Text)
		}
	}
}
