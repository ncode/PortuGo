// Package cp1252 converts the single-byte text repertoire used by source files
// and character-code built-ins.
package cp1252

import "unicode/utf8"

// DecodeByte returns a Windows-1252 character, or RuneError for an undefined byte.
func DecodeByte(c byte) rune {
	if c < 0x80 || c >= 0xa0 {
		return rune(c)
	}
	switch c {
	case 0x80:
		return '€'
	case 0x82:
		return '‚'
	case 0x83:
		return 'ƒ'
	case 0x84:
		return '„'
	case 0x85:
		return '…'
	case 0x86:
		return '†'
	case 0x87:
		return '‡'
	case 0x88:
		return 'ˆ'
	case 0x89:
		return '‰'
	case 0x8a:
		return 'Š'
	case 0x8b:
		return '‹'
	case 0x8c:
		return 'Œ'
	case 0x8e:
		return 'Ž'
	case 0x91:
		return '‘'
	case 0x92:
		return '’'
	case 0x93:
		return '“'
	case 0x94:
		return '”'
	case 0x95:
		return '•'
	case 0x96:
		return '–'
	case 0x97:
		return '—'
	case 0x98:
		return '˜'
	case 0x99:
		return '™'
	case 0x9a:
		return 'š'
	case 0x9b:
		return '›'
	case 0x9c:
		return 'œ'
	case 0x9e:
		return 'ž'
	case 0x9f:
		return 'Ÿ'
	default:
		return utf8.RuneError
	}
}

// EncodeRune returns the Windows-1252 byte for r when it is representable.
func EncodeRune(r rune) (byte, bool) {
	if r >= 0 && r < 0x80 || r >= 0xa0 && r <= 0xff {
		return byte(r), true
	}
	if r != utf8.RuneError {
		for b := byte(0x80); b < 0xa0; b++ {
			if DecodeByte(b) == r {
				return b, true
			}
		}
	}
	return 0, false
}
