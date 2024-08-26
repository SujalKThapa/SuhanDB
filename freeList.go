package main

import "encoding/binary"

const metaPage = 0

type freeList struct {
	maxPage       pgnum
	releasedPages []pgnum
}

func newFreeList() *freeList {
	return &freeList{
		maxPage:       metaPage,
		releasedPages: []pgnum{},
	}
}

func (fr *freeList) getNextPage() pgnum {
	if len(fr.releasedPages) != 0 {
		pageID := fr.releasedPages[len(fr.releasedPages)-1]
		fr.releasedPages = fr.releasedPages[:len(fr.releasedPages)-1]
		return pageID
	}
	fr.maxPage += 1
	return fr.maxPage
}

func (fr *freeList) releasePage(page pgnum) {
	fr.releasedPages = append(fr.releasedPages, page)
}

func (fr *freeList) serialize(buf []byte) []byte {
	pos := 0
	binary.LittleEndian.PutUint16(buf[pos:], uint16(fr.maxPage))
	pos += 2
	binary.LittleEndian.PutUint16(buf[pos:], uint16(len(fr.releasedPages)))
	pos += 2

	for _, page := range fr.releasedPages {
		binary.LittleEndian.PutUint64(buf[pos:], uint64(page))
		pos += pageNumSize
	}
	return buf
}

func (fr *freeList) deserialize(buf []byte) {
	pos := 0
	fr.maxPage = pgnum(binary.LittleEndian.Uint16(buf[pos:]))
	pos += 2

	releasedPagesCount := int(binary.LittleEndian.Uint16(buf[pos:]))
	pos += 2

	for i := 0; i < releasedPagesCount; i++ {
		fr.releasedPages = append(fr.releasedPages, pgnum(binary.LittleEndian.Uint64(buf[pos:])))
		pos += pageNumSize
	}

}
