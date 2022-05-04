package service

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

type MmapTest struct {
}

func (m *MmapTest) Test() {
	pagesize := os.Getpagesize()
	file, err := os.OpenFile("test.mmap", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return
	}

	state, err := file.Stat()
	if err != nil {
		return
	}

	if state.Size() == 0 {
		var n int
		n, err = file.WriteAt(bytes.Repeat([]byte{'0'}, pagesize), 0)
		if err != nil {
			return
		}
		fmt.Println("写的n : ", n)
		state, err = file.Stat()
		if err != nil {
			return
		}
	}

	fmt.Printf("pagesize : %d\nfilesize: %d\n", pagesize, state.Size())
	data, err := unix.Mmap(int(file.Fd()), 0, int(state.Size()), unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return
	}

	err = file.Close()
	if err != nil {
		return
	}

	for i := 0; i < 8; i++ {
		data[i] = 'l'
	}

	err = unix.Munmap(data)
}
