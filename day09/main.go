package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type FileNode struct {
	address   int
	id        int
	length    int
	freeSpace int
	next      *FileNode
	prev      *FileNode
}

type FileSystem struct {
	start *FileNode
	end   *FileNode
	size  int
}

func NewFileNode(address int, id int, length int, freeSpace int) *FileNode {
	file := FileNode{address: address, id: id, length: length, freeSpace: freeSpace, next: nil, prev: nil}

	return &file
}

func (list *FileSystem) Append(newFileNode *FileNode) {
	if list.end != nil {
		list.end.next = newFileNode
	}

	newFileNode.prev = list.end
	list.end = newFileNode
	if list.start == nil {
		list.start = list.end
	}
}

func RuneToDigit(ch rune) int {
	return int(ch - '0')
}

func PopulateFileSystem(r *bufio.Reader) FileSystem {
	fileSystem := FileSystem{}

	var id int = 0
	toggle := true
	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if ch == '\n' {
			break
		}

		value := RuneToDigit(ch)
		if toggle {
			file := NewFileNode(fileSystem.size, id, value, 0)
			fileSystem.Append(file)
		} else {
			fileSystem.end.freeSpace = value
			id++
		}
		fileSystem.size += int(value)

		toggle = !toggle
	}

	return fileSystem
}

func MoveFile(left *FileNode, right *FileNode, fileSys *FileSystem) {
	if left.next == right {
		right.freeSpace += left.freeSpace
		left.freeSpace = 0
		return
	}

	// Create a new node to represent the right-most file moving to the next free space.
	newAddress := left.address + left.length
	movedFile := NewFileNode(newAddress, right.id, right.length, left.freeSpace-right.length)

	// Insert a copy of the right node into position after the left node.
	movedFile.prev = left
	movedFile.next = left.next
	left.next.prev = movedFile
	left.next = movedFile

	left.freeSpace = 0

	// Sever the connections to the right node.
	right.prev.next = right.next
	if right == fileSys.end {
		fileSys.end = right.prev
		fileSys.size -= (right.length + right.freeSpace)
	}
	if right.next != nil {
		right.next.prev = right.prev
	}
	// Update the left file with the newly added free space.
	if right.prev != left {
		right.prev.freeSpace += right.length + right.freeSpace
	}

	if right.length < 0 || fileSys.end.freeSpace < 0 {
		panic("File size can't be negative.")
	}
}

func SeekFit(openFile *FileNode, toFit *FileNode) (*FileNode, bool) {
	found := false
	for openFile != toFit {
		if openFile.freeSpace >= toFit.length {
			found = true
			break
		}
		openFile = openFile.next
	}

	return openFile, found
}

func SeekFreeSpace(file *FileNode, freeSpace int) *FileNode {
	for file != nil && file.freeSpace < freeSpace {
		file = file.next
	}

	return file
}

func Defrag(fileSystem FileSystem) {
	firstSpace := SeekFreeSpace(fileSystem.start, 1) // The left-most file with any available space.
	currentFile := fileSystem.end                    // The file that we're currently attempting to move.

	cachedIDs := make(map[int]bool)
	for firstSpace.address+firstSpace.length < currentFile.address {
		// Ignore files that we've already moved.
		if cachedIDs[currentFile.id] {
			currentFile = currentFile.prev
			continue
		}
		cachedIDs[currentFile.id] = true

		// Look for the first file with enough space for the current file.
		openFile, found := SeekFit(firstSpace, currentFile)
		if found {
			MoveFile(openFile, currentFile, &fileSystem)
			firstSpace = SeekFreeSpace(firstSpace, 1)
		}

		currentFile = currentFile.prev
	}
}

func Sum(start int, end int) int {
	return (end - start + 1) * (start + end) / 2
}

func CalculateChecksum(fileSys FileSystem) int {
	position := 0
	checkSum := 0
	for f := fileSys.start; f != nil; f = f.next {
		checkSum += Sum(position, position+int(f.length)-1) * int(f.id)
		position += int(f.length + f.freeSpace)
	}

	return checkSum
}

func WriteListToFile(fileSystem FileSystem, filename string) {
	file, err := os.Create(filename)
	Check(err)

	for f := fileSystem.start; f != nil; f = f.next {
		file.WriteString(fmt.Sprintf("ID: %4d, size: %d, free space %d\n", f.id, f.length, f.freeSpace))
	}
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	const inputName = "./input"

	file, err := os.Open(fmt.Sprintf("%s.txt", inputName))
	Check(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	fileSystem := PopulateFileSystem(reader)
	Defrag(fileSystem)

	// WriteListToFile(fileSystem, fmt.Sprintf("%s_defrag.txt", inputName))
	fmt.Printf("Checksum: %d\n", CalculateChecksum(fileSystem))
}
