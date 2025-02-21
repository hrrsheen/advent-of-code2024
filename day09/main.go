package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Node struct {
	id        int
	length    int
	freeSpace int
	next      *Node
	prev      *Node
}

type FileSystem struct {
	start *Node
	end   *Node
	size  int
}

func NewNode(id int, length int, freeSpace int) *Node {
	node := Node{id: id, length: length, freeSpace: freeSpace, next: nil, prev: nil}

	return &node
}

func (list *FileSystem) Append(newNode *Node) {
	if list.end != nil {
		list.end.next = newNode
	}

	newNode.prev = list.end
	list.end = newNode
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
			file := NewNode(id, value, 0)
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

func MoveFile(left *Node, right *Node, fileSys *FileSystem) {
	if left.next == right {
		right.freeSpace += left.freeSpace
		left.freeSpace = 0
		return
	}

	// Create a new node to represent the right-most file moving to the next free space.
	interspersed := NewNode(right.id, right.length, left.freeSpace-right.length)

	// Insert a copy of the right node into position after the left node.
	interspersed.prev = left
	interspersed.next = left.next
	left.next.prev = interspersed
	left.next = interspersed

	left.freeSpace = 0

	// Sever the connection to the right node.
	if right == fileSys.end {
		fileSys.end = right.prev
		fileSys.size -= (int(right.length + right.freeSpace))
	}
	right.prev.next = right.next
	if right.next != nil {
		right.next.prev = right.prev
	}
	// right.prev.freeSpace += right.length + right.freeSpace

	if right.length < 0 || fileSys.end.freeSpace < 0 {
		panic("File size can't be negative.")
	}
}

func SeekFit(openFile *Node, toFit *Node) (*Node, bool) {
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

func SeekFreeSpace(file *Node, freeSpace int) (*Node, int) {
	seekDistance := 0
	i := 0
	for file != nil && file.freeSpace < freeSpace {
		if i > 0 {
			seekDistance += file.length + file.freeSpace
		}
		file = file.next
		i++
	}

	if i > 0 {
		seekDistance += file.length
	}

	return file, seekDistance
}

func Defrag(fileSystem FileSystem) {
	leftmostSpace, seekDistance := SeekFreeSpace(fileSystem.start, 1)
	endFile := fileSystem.end
	cachedIDs := make(map[int]bool)

	leftPos := seekDistance
	rightPos := fileSystem.size - (int(endFile.length + endFile.freeSpace))
	for leftPos < rightPos {
		openFile, found := SeekFit(leftmostSpace, endFile)
		if found && !cachedIDs[endFile.id] {
			MoveFile(openFile, endFile, &fileSystem)
			leftmostSpace, seekDistance = SeekFreeSpace(leftmostSpace, 1)
			leftPos += seekDistance
		}

		cachedIDs[endFile.id] = true
		rightPos -= int(endFile.prev.length + endFile.prev.freeSpace)

		if found && endFile.prev != openFile {
			endFile.prev.freeSpace += endFile.length + endFile.freeSpace
		}
		endFile = endFile.prev
	}
	fmt.Printf("defragged\n")
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

	WriteListToFile(fileSystem, fmt.Sprintf("%s_defrag.txt", inputName))
	fmt.Printf("Checksum: %d\n", CalculateChecksum(fileSystem))
}
