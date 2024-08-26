package main

import (
	"errors"
	"fmt"
	"os"
)

type pgnum uint64

type page struct {
	num  pgnum
	data []byte
}

type dal struct {
	file     *os.File
	pageSize int
	*meta
	*freeList
}

func newDal(path string) (*dal, error) {
	dal := &dal{
		meta:     newEmptyMeta(),
		pageSize: os.Getpagesize(),
	}

	// exist
	if _, err := os.Stat(path); err == nil {
		dal.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			_ = dal.close()
			return nil, err
		}

		meta, err := dal.readMeta()
		if err != nil {
			return nil, err
		}
		dal.meta = meta

		freeList, err := dal.readFreeList()
		if err != nil {
			return nil, err
		}
		dal.freeList = freeList
		// doesn't exist
	} else if errors.Is(err, os.ErrNotExist) {
		// init freelist
		dal.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			_ = dal.close()
			return nil, err
		}

		dal.freeList = newFreeList()
		dal.freeListPage = dal.getNextPage()
		_, err := dal.writeFreeList()
		if err != nil {
			return nil, err
		}

		// write meta page
		_, err = dal.writeMeta(dal.meta) // other error
	} else {
		return nil, err
	}
	return dal, nil
}

func (d *dal) readFreeList() (*freeList, error) {
	p, err := d.readPage(d.freeListPage)
	if err != nil {
		return nil, err
	}

	freelist := newFreeList()
	freelist.deserialize(p.data)
	return freelist, nil
}

func (d *dal) writeFreeList() (*page, error) {
	p := d.allocateEmptyPage()
	p.num = d.freeListPage
	d.freeList.serialize(p.data)

	err := d.writePage(p)
	if err != nil {
		return nil, err
	}
	d.freeListPage = p.num
	return p, nil
}

func (d *dal) close() error {
	if d.file != nil {
		err := d.file.Close()
		if err != nil {
			return fmt.Errorf("could not close file: %s", err)
		}
		d.file = nil
	}
	return nil
}

func (d *dal) allocateEmptyPage() *page {
	return &page{
		data: make([]byte, d.pageSize),
	}
}

func (d *dal) readPage(pageNum pgnum) (*page, error) {
	p := d.allocateEmptyPage()
	offset := int(pageNum) * d.pageSize
	//Reads d file from offset position and puts it in p.data
	_, err := d.file.ReadAt(p.data, int64(offset))
	if err != nil {
		return nil, err
	}
	return p, err
}

func (d *dal) writePage(p *page) error {
	offset := int64(p.num) * int64(d.pageSize)
	_, err := d.file.WriteAt(p.data, offset)
	return err
}

func (d *dal) writeMeta(meta *meta) (*page, error) {
	p := d.allocateEmptyPage()
	p.num = metaPageNum
	meta.serialize(p.data)

	err := d.writePage(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *dal) readMeta() (*meta, error) {
	p, err := d.readPage(metaPageNum)
	if err != nil {
		return nil, err
	}

	meta := newEmptyMeta()
	meta.deserialize(p.data)
	return meta, nil
}
