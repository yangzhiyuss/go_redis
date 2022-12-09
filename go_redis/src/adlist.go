package src

const (
	// 从表头向表尾进行迭代
	AL_START_HEAD = 0
	// 从表尾到表头进行迭代
	AL_START_TAIL = 1
)

// 双端链表节点
type ListNode struct {
	pre   *ListNode
	next  *ListNode
	value interface{}
}

// 双端链表迭代器
type listIter struct {
	Next      *ListNode
	Direction int
}

// 双端链表结构
type List struct {
	head  *ListNode
	tail  *ListNode
	dup   func(ptr interface{}) interface{}
	match func(ptr interface{}, key interface{}) int
	len   uint64
}

// 基本元素的获取
func (list *List) Length() uint64 {
	return list.len
}

func (list *List) First() *ListNode {
	return list.head
}

func (list *List) Last() *ListNode {
	return list.tail
}

func (listNode *ListNode) PreNode() *ListNode {
	return listNode.pre
}

func (listNode *ListNode) NextNode() *ListNode {
	return listNode.next
}

func (listNode *ListNode) NodeValue() interface{} {
	return listNode.value
}

func (list *List) SetDupMethod(m func(ptr interface{}) interface{}) {
	list.dup = m
}

func (list *List) SetMatchMethod(m func(ptr interface{}, key interface{}) int) {
	list.match = m
}

func (list *List) GetDupMethod() func(ptr interface{}) interface{} {
	return list.dup
}

func (list *List) GetMatchMethod() func(ptr interface{}, key interface{}) int {
	return list.match
}

// 给链表表头增加新节点
func (list *List) AddNodeHead(value interface{}) {
	node := new(ListNode)
	node.value = value

	if list.len == 0 {
		list.head = node
		list.tail = node
		node.next = nil
		node.pre = nil
	} else {
		node.pre = nil
		node.next = list.head
		list.head.pre = node
		list.head = node
	}
	list.len++
}

// 给链表表尾增加新节点
func (list *List) AddNodeTail(value interface{}) {
	node := new(ListNode)
	node.value = value
	if list.len == 0 {
		list.head = node
		list.tail = node
		node.next = nil
		node.pre = nil
	} else {
		node.next = nil
		node.pre = list.tail
		list.tail.next = node
		list.tail = node
	}
	list.len ++
}

//在链表中插入数据
//0:向前插入节点
//1:向后插入节点
func (list *List) InsertNode(old_node *ListNode, value interface{}, after int) {
	node := new(ListNode)
	node.value = value
	if after == 1 {
		node.pre = old_node
		node.next = old_node.next
		if list.tail == old_node {
			list.tail = node
		}
	} else {
		node.next = old_node
		node.pre = old_node.pre
		if list.head == old_node {
			list.head = node
		}
	}
	if node.pre != nil {
		node.pre.next = node
	}
	if node.next != nil {
		node.next.pre = node
	}
	list.len++
}

func (list *List) DelNode(node *ListNode)  {
	if node.pre != nil {
		node.pre = node.next
	} else {
		list.head = node.next
	}

	if node.next != nil {
		node.next = node.pre
	} else {
		list.tail = node.pre
	}

	//情况node的值
	node.next = nil
	node.pre = nil
	node.value = nil

	list.len--
}

