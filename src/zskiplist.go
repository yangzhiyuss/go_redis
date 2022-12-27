package src

import (
	"math/rand"
	"time"
)

const (
	ZSKIPLIST_MAXLEVEL = 32
	ZSKIPLIST_P        = 0.25
)

type zSkipListLevel struct {
	//向前指针
	forward *zSkipListNode
	//跨度
	span uint
}

// 跳跃表节点
type zSkipListNode struct {
	//分值
	score float64
	//
	value interface{}
	//后退指针
	backward *zSkipListNode
	//层
	level []zSkipListLevel
}

// 跳跃表结构
type zSkipList struct {
	header *zSkipListNode
	tail   *zSkipListNode
	length uint64
	//表中最大节点的层数
	level int
}

func zslCreateNode(level int, score float64) *zSkipListNode {
	node := new(zSkipListNode)
	node.score = score
	node.level = make([]zSkipListLevel, level)
	return node
}

func ZSLCreate() *zSkipList {
	zsl := new(zSkipList)
	zsl.length = 0
	zsl.level = 1
	zsl.tail = nil

	zsl.header = zslCreateNode(ZSKIPLIST_MAXLEVEL, 0)
	for i := 0; i < ZSKIPLIST_MAXLEVEL; i++ {
		zsl.header.level[i].forward = nil
		zsl.header.level[i].span = 0
	}
	zsl.header.backward = nil
	return zsl
}

func randomLevel() int {
	level := 1
	for float64(rand.Int() & 0xFFFF) < (ZSKIPLIST_P * 0xFFFF) {
		level += 1
		rand.Seed(time.Now().UnixNano())
	}

	if level < ZSKIPLIST_MAXLEVEL {
		return level
	} else {
		return ZSKIPLIST_MAXLEVEL
	}
}

func (zsl *zSkipList) Insert(score float64) *zSkipListNode {
	updata := make([]*zSkipListNode, ZSKIPLIST_MAXLEVEL)
	rank := make([]uint, ZSKIPLIST_MAXLEVEL)
	x := zsl.header
	//找到对应的节点
	for i := zsl.level - 1; i >= 0; i++ {
		if i == zsl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i + 1]
		}

		for x.level[i].forward != nil && x.level[i].forward.score < score {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		updata[i] = x
	}

	level := randomLevel()
	//层数超过已存的最大层
	if level > zsl.level {
		for i := zsl.level; i < level; i++ {
			rank[i] = 0
			updata[i] = zsl.header
			updata[i].level[i].span = uint(zsl.length)
		}
	}

	x = zslCreateNode(level, score)
	for i := 0; i < level; i++ {
		x.level[i].forward = updata[i].level[i].forward
		updata[i].level[i].forward = x

		x.level[i].span = updata[i].level[i].span - (rank[0] - rank[i])
		updata[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	for i := level; i < zsl.level; i++ {
		updata[i].level[i].span ++
	}

	if updata[0] == zsl.header {
		x.backward = nil
	} else {
		x.backward = updata[0]
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		zsl.tail = x
	}
	zsl.length++

	return x
}

func (zsl *zSkipList) Delete(score float64) int {
	update := make([]*zSkipListNode, ZSKIPLIST_MAXLEVEL)
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for  x.level[i].forward != nil && x.level[i].forward.score > score {
			x = x.level[i].forward
		}
		update[i] = x
	}

	x = x.level[0].forward

	if x != nil && score == x.score {
		zsl.deleteNode(x, update)
		return 1
	} else {
		return 0
	}
}

func (zsl *zSkipList) deleteNode(x *zSkipListNode, update []*zSkipListNode)  {
	for i := 0; i < zsl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
		x.level[i].forward = nil
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		zsl.tail = x.backward
	}
	x.backward = nil

	for zsl.level > 1 && zsl.header.level[zsl.level - 1].forward == nil {
		zsl.level--
	}

	zsl.length--
}