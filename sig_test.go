package htmlsig

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureFromReader(t *testing.T) {
	doc := `<html>
				<body>
				<!-- Example web document -->
					<p><b>This text is bold</b></p>
					<p><strong>This text is strong</strong></p>
					<p><big>This text is big</big></p>
					<p><em>This text is emphasized</em></p>
					<p><i>This text is italic</i></p>
					<p><small>This text is small</small></p>
					<p>This is<sub> subscript</sub> and
					<sup>superscript</sup></p>
				</body>
			</html>`
	reader := strings.NewReader(doc)
	signature := NewHTMLSignature(25)
	signature.FromReader(reader)
	signature.Close()
	assert.Equal(
		t,
		[]byte{0, 0, 0, 0, 4, 3, 0, 3, 0, 9, 3, 0, 8, 0, 8, 0, 8, 0, 0, 19, 2},
		signature.StructureFingerprint)
	assert.Equal(
		t,
		[]uint64{
			0xc3bfe37667d7507, 0x40e505931a078f28, 0xdeffcf9a7d3f5c2,
			0x1c9192791bd6e870, 0x1dcff8a00209aa3e, 0x3b490dea0dd81f36,
			0x16a397439478b1b6, 0xe48b7d51f34c2e0, 0x94cbc618e63d2ff,
			0x10b59c0e0d1a7afc, 0x3e9fab5170e14bb2, 0x283303b15d5a363c,
			0xac7a0d885bc4442, 0x34fe877694d9e7c, 0x1c1077a6c81b1464,
			0x4d9a5a2fe5e0d88, 0x121ea154885460c, 0x4ffdc594c66bf64,
			0x18be82ee6ee96e3d, 0x46b357876098dc, 0x392fb13aece78f0d,
			0x1e148002b3eddf4, 0x93da3adbb03c26e, 0xe248152ce826a18,
			0x36e38b5de308eef},
		signature.TextFingerprint)
}

func benchmarkSignatureFromReader(signatureSize int, b *testing.B) {
	doc := `<html>
		<body>
		<!-- Example web document -->
			<p><b>This text is bold</b></p>
			<p><strong>This text is strong</strong></p>
			<p><big>This text is big</big></p>
			<p><em>This text is emphasized</em></p>
			<p><i>This text is italic</i></p>
			<p><small>This text is small</small></p>
			<p>This is<sub> subscript</sub> and
			<sup>superscript</sup></p>
		</body>
	</html>`
	for n := 0; n < b.N; n++ {
		reader := strings.NewReader(doc)
		signature := NewHTMLSignature(signatureSize)
		signature.FromReader(reader)
		signature.Close()
	}
}

func BenchmarkSignatureFromReader1(b *testing.B)  { benchmarkSignatureFromReader(1, b) }
func BenchmarkSignatureFromReader2(b *testing.B)  { benchmarkSignatureFromReader(2, b) }
func BenchmarkSignatureFromReader4(b *testing.B)  { benchmarkSignatureFromReader(4, b) }
func BenchmarkSignatureFromReader8(b *testing.B)  { benchmarkSignatureFromReader(8, b) }
func BenchmarkSignatureFromReader16(b *testing.B) { benchmarkSignatureFromReader(16, b) }
func BenchmarkSignatureFromReader24(b *testing.B) { benchmarkSignatureFromReader(24, b) }
