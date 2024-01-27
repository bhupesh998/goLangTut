package main 

import "fmt"

const SIZE = 2

type Node struct{

	Val string
	Left *Node
	Right *Node

}

type Queue struct{
	Head *Node
	Tail *Node
	Length int
}

type Cache struct{
	Queue Queue
	Hash Hash
}

type Hash map[string]*Node

func NewCache() Cache{

	return Cache{Queue : NewQueue(),Hash : Hash{}}

}

func NewQueue() Queue{

	head:= &Node{}
	tail :=&Node{}

	head.Right = tail
	tail.Left = head

	return Queue{Head: head, Tail: tail}
}

func (c *Cache) Check(word string){
	node := &Node{}

	if val , ok := c.Hash[word]; ok{
		node = c.Remove(val)
	}else{
		node = &Node{Val: word}
	}
	c.Add(node)
	c.Hash[word]= node

}


func (c *Cache) Remove (n *Node) *Node {
	fmt.Printf("Remove %s\n", n.Val)
	left := n.Left 
	right := n.Right

	right.Left = left
	left.Right = right

	c.Queue.Length -= 1
	delete(c.Hash, n.Val)

	return n
}

func (c *Cache) Add (n *Node) {
	fmt.Printf("Add %s\n", n.Val)
	temp := c.Queue.Head.Right

	c.Queue.Head.Right = n
	n.Left = c.Queue.Head
	n.Right = temp
	temp.Left = n

	c.Queue.Length++

	if c.Queue.Length > SIZE{
		c.Remove(c.Queue.Tail.Left)
	}

}

func (c *Cache) Display (){
	c.Queue.Display()

} 

func(q *Queue) Display (){
	node := q.Head.Right
	fmt.Printf("%d- [" , q.Length)

	for i:=0;i<q.Length;i++ {
		fmt.Printf("{%s}", node.Val)
		if i<q.Length -1 {
			fmt.Printf("<--->")
		}
		node = node.Right
	}

	fmt.Println("]")
}


func main(){
	fmt.Println("GO Lang Cache Start")
	cache := NewCache()

	for _, word := range []string{"A", "B", "C","D"}{
				cache.Check(word)
				cache.Display()
		}

	fmt.Println("final cache display")
	cache.Display()
}