package main

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
