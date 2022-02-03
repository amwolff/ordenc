//go:build go1.18
// +build go1.18

package ordenc

import (
	"testing"

	"github.com/zeebo/assert"
	"github.com/zeebo/mwc"
)

func FuzzEncryptDecrypt(f *testing.F) {
	rng := mwc.Rand()

	k, err := NewRandomKey(&rng)
	assert.NoError(f, err)

	for _, s := range [...][]byte{ // Add paths to the seed corpus.
		[]byte("sample.jpg"),
		[]byte("photos/2022/January/dog.jpg"),
		[]byte("photos/2022/February/dog2.jpg"),
		[]byte("photos/2022/February/dog3.jpg"),
		[]byte("photos/2022/February/dog4.jpg"),
	} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		e := Encrypt(k, input, nil)

		d, ok := Decrypt(k, e, nil)

		assert.That(t, ok)
		assert.Equal(t, input, d)
	})
}

func FuzzEncryptOrderedPairs(f *testing.F) {
	rng := mwc.Rand()

	k, err := NewRandomKey(&rng)
	assert.NoError(f, err)

	for _, s := range [...][2][]byte{ // Add pairs to the seed corpus.
		{
			[]byte("photos/2022/January/dog.jpg"),
			[]byte("sample.jpg"),
		},
		{
			[]byte("photos/2022/February/dog2.jpg"),
			[]byte("photos/2022/January/dog.jpg"),
		},
		{
			[]byte("photos/2022/February/dog2.jpg"),
			[]byte("photos/2022/February/dog3.jpg"),
		},
		{
			[]byte("photos/2022/February/dog3.jpg"),
			[]byte("photos/2022/February/dog4.jpg"),
		},
	} {
		f.Add(s[0], s[1])
	}

	f.Fuzz(func(t *testing.T, p1, p2 []byte) {
		e1 := Encrypt(k, p1, nil)
		e2 := Encrypt(k, p2, nil)

		// t.Logf("p1=%s / p2=%s", string(p1), string(p2))
		// t.Logf("e1=%s / e2=%s", string(e1), string(e2))

		if string(p1) < string(p2) {
			assert.That(t, string(e1) < string(e2))
		} else if string(p1) > string(p2) {
			assert.That(t, string(e1) > string(e2))
		} else {
			assert.That(t, string(e1) == string(e2))
		}
	})
}
