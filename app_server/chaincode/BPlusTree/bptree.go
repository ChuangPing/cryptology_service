package bptree

import (
	"sync"
)

type ValueData struct {
	Authority string // 权限
	EncData   string // 关键字密文
}

//	BPItem：叶节点数据结构
type BPItem struct {
	Key int64
	Val interface{}
}

//type StringKV struct {
//	Key string
//	Val interface{}
//}

//	BPNode：树节点数据结构：叶节点与索引节点统一为BPNode
type BPNode struct {
	MaxKey int64
	Nodes  []*BPNode // 索引节点
	Items  []BPItem  //	叶节点
	Next   *BPNode
}

// BPTree：B+树结构
type BPTree struct {
	mutex sync.RWMutex // 互斥锁，支持并发读，互斥写
	root  *BPNode      // 根节点
	width int          // 树的阶
	halfw int          // 节点上线
}

//	NewBPTree：初始化B+树
func NewBPTree(width int) *BPTree {
	if width < 3 {
		width = 3
	}

	var bt = &BPTree{}
	bt.root = NewLeafNode(width) // 初始只有叶节点
	bt.width = width
	bt.halfw = (bt.width + 1) / 2
	return bt
}

//申请width+1是因为插入时可能暂时出现节点key大于申请width的情况,待后期再分裂处理  --初始化叶子节点
func NewLeafNode(width int) *BPNode {
	var node = &BPNode{}
	node.Items = make([]BPItem, width+1)
	node.Items = node.Items[0:0]
	return node
}

//申请width+1是因为插入时可能暂时出现节点key大于申请width的情况,待后期再分裂处理 -- 初始化索引节点
func NewIndexNode(width int) *BPNode {
	var node = &BPNode{}
	node.Nodes = make([]*BPNode, width+1)
	node.Nodes = node.Nodes[0:0]
	return node
}

func (node *BPNode) findItem(key int64) int {
	num := len(node.Items)
	for i := 0; i < num; i++ {
		if node.Items[i].Key > key {
			return -1
		} else if node.Items[i].Key == key {
			return i
		}
	}
	return -1
}

//	插入函数
func (t *BPTree) Set(key int64, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.setValue(nil, t.root, key, value)
}

func (t *BPTree) setValue(parent *BPNode, node *BPNode, key int64, value interface{}) {
	//	通过索引节点递归找到需要插入的位置

	for i := 0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey || i == len(node.Nodes)-1 {
			t.setValue(node, node.Nodes[i], key, value)
			break
		}
	}

	//叶子结点，添加数据  -- 直接添加不考虑是否需要分裂
	if len(node.Nodes) < 1 {
		node.setValue(key, value)

	}

	//结点分裂   -- 添加完成后在考虑是否需要分裂
	childNode := t.splitNode(node)

	//	有节点分裂
	if childNode != nil {
		//若父结点不存在，则创建一个父节点  -- 第一此初始化时还没有索引节点，因此parent = nil
		if parent == nil {
			parent = NewIndexNode(t.width)
			parent.addChild(node)
			//	将向上分裂的节点赋值给根节点
			t.root = parent
		}
		//添加结点到父亲结点
		parent.addChild(childNode)
	} else {
		//	当插入的节点未达到分裂要求且key是最大的，此时更新树的索引节点maxKey
		if node.MaxKey > t.root.MaxKey {
			t.root.MaxKey = node.MaxKey
		}
		if parent != nil {
			for i := 0; i < len(parent.Nodes); i++ {
				if parent.Nodes[i].MaxKey > parent.MaxKey {
					parent.MaxKey = parent.Nodes[i].MaxKey
				}
			}
		}

	}
}

func (node *BPNode) setValue(key int64, value interface{}) {
	item := BPItem{key, value}
	num := len(node.Items)

	//先确定三种特殊情况  -- 还没有叶节点直接插入  key比第一个叶节点小插入在左边  key比最后一个叶节点大插入到右边
	if num < 1 {
		node.Items = append(node.Items, item)
		// 更新节点的最大索引
		node.MaxKey = item.Key
		return
	} else if key < node.Items[0].Key {
		// 放在最左边
		node.Items = append([]BPItem{item}, node.Items...)
		return
	} else if key > node.Items[num-1].Key {
		//	放最右边
		node.Items = append(node.Items, item)
		node.MaxKey = item.Key
		return
	}
	// 排除上面的特殊情况，现在所有的一般情况 node.MaxKey > item.Key,因此不用在设置node.MaxKey
	for i := 0; i < num; i++ {
		if node.Items[i].Key > key {
			// 末尾多开辟一个空间
			node.Items = append(node.Items, BPItem{})
			// 将node.Items[i:] 赋值到 node.Items[i+1:] -- 相当于整体向右移一位
			copy(node.Items[i+1:], node.Items[i:])
			//将item插入到正确位置
			node.Items[i] = item
			return
		} else if node.Items[i].Key == key {
			//	插入关键字一样，可以修改记录
			node.Items[i] = item
			return
		}
	}
}

func (t *BPTree) splitNode(node *BPNode) *BPNode {
	if len(node.Nodes) > t.width {
		//创建新结点 -- 索引节点分裂
		halfw := t.width
		//	初始化索引节点
		node2 := NewIndexNode(t.width)
		//	将分裂的部分添加到新的索引节点
		node2.Nodes = append(node2.Nodes, node.Nodes[halfw:len(node.Nodes)]...)
		//	初始化索引节点的maxKey
		node2.MaxKey = node2.Nodes[len(node2.Nodes)-1].MaxKey

		//修改原结点数据
		node.Nodes = node.Nodes[0:halfw]
		node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey

		return node2
	} else if len(node.Items) > t.width {
		//创建新结点 -- 叶子节点分裂
		halfw := t.width
		//	初始化叶节点
		node2 := NewLeafNode(t.width)
		node2.Items = append(node2.Items, node.Items[halfw:len(node.Items)]...)
		node2.MaxKey = node2.Items[len(node2.Items)-1].Key

		//修改原结点数据
		node.Next = node2
		node.Items = node.Items[0:halfw]
		node.MaxKey = node.Items[len(node.Items)-1].Key

		return node2
	}

	return nil
}

func (node *BPNode) addChild(child *BPNode) {
	num := len(node.Nodes)
	if num < 1 {
		//	还没有父索引节点
		node.Nodes = append(node.Nodes, child)
		node.MaxKey = child.MaxKey
		return
	} else if child.MaxKey < node.Nodes[0].MaxKey {
		//	放最左边
		node.Nodes = append([]*BPNode{child}, node.Nodes...)
		return
	} else if child.MaxKey > node.Nodes[num-1].MaxKey {
		//	放最右边
		node.Nodes = append(node.Nodes, child)
		node.MaxKey = child.MaxKey
		return
	}

	//	分裂节点需要插入到父节点的中间某个位置
	for i := 0; i < num; i++ {
		if node.Nodes[i].MaxKey > child.MaxKey {
			// 将child插入到 i 位置
			node.Nodes = append(node.Nodes, nil)
			copy(node.Nodes[i+1:], node.Nodes[i:])
			node.Nodes[i] = child
			return
		}
	}
}

//	B+树查询的方法： 根据key获取value

func (t *BPTree) Get(key int64) (interface{}, int) {
	var count, j int
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := t.root
	for i := 0; i < len(node.Nodes); i++ {
		j++
		if key <= node.Nodes[i].MaxKey {
			// 统计索引节点查找次数
			count = j
			// 下一层节点
			node = node.Nodes[i]
			j = i
			i = -1
		} else {
			continue
		}
	}

	//没有到达叶子结点说明没有查找到
	if len(node.Nodes) > 0 {
		return nil, -1
	}

	//定位到叶节点进行关键字比对  -- 叶节点不使用二分查找
	//for i := 0; i < len(node.Items); i++ {
	//	count++
	//	if node.Items[i].Key == key {
	//		return node.Items[i].Val, count
	//	}
	//}

	//叶节点使用二分查找
	return node.binarySearch(key, count)
}

func (node *BPNode) binarySearch(key int64, count int) (interface{}, int) {
	var low, mid int
	var high = len(node.Items) - 1
	//mid := 0
	for {
		count++
		if low > high {
			break
		}
		mid = (low + high) / 2
		if node.Items[mid].Key == key {
			return node.Items[mid].Val, count
		} else if node.Items[mid].Key > key {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return nil, -1
}

//	获取叶节点全部数据
func (t *BPTree) GetData() map[int64]interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.getData(t.root)
}

func (t *BPTree) getData(node *BPNode) map[int64]interface{} {
	data := make(map[int64]interface{})
	for {
		if len(node.Nodes) > 0 {
			for i := 0; i < len(node.Nodes); i++ {
				// 递归遍历叶节点，并记录每个叶节点的最大值
				data[node.Nodes[i].MaxKey] = t.getData(node.Nodes[i])
			}
			break
		} else {
			for i := 0; i < len(node.Items); i++ {
				//	获取全部叶节点信息
				data[node.Items[i].Key] = node.Items[i].Val
			}
			break
		}
	}
	return data
}

//	删除节点方法
func (t *BPTree) Remove(key int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.deleteItem(nil, t.root, key)
}

func (t *BPTree) deleteItem(parent *BPNode, node *BPNode, key int64) {
	//	通过索引节点递归的找到删除位置
	for i := 0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey {
			t.deleteItem(node, node.Nodes[i], key)
			break
		}
	}

	if len(node.Nodes) < 1 {
		//待删除元素为叶子节点
		node.deleteItem(key)
		if len(node.Items) < t.halfw {
			//删除记录后若结点的子项<m/2，则从兄弟结点移动记录，或者合并结点
			t.itemMoveOrMerge(parent, node)
		}
	} else {
		//待删除元素为非叶节点
		//若结点的子项<m/2，则从兄弟结点移动记录，或者合并结点
		node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey
		if len(node.Nodes) < t.halfw {
			t.childMoveOrMerge(parent, node)
		}
	}
}

//	根据key在叶节点进行删除
func (node *BPNode) deleteItem(key int64) bool {
	num := len(node.Items)
	for i := 0; i < num; i++ {
		if node.Items[i].Key > key {
			return false
		} else if node.Items[i].Key == key {

			//	删除node.Items[i] 位置元素
			copy(node.Items[i:], node.Items[i+1:])
			//	去除末尾重复元素
			node.Items = node.Items[0 : len(node.Items)-1]
			//	跟新剩余叶节点的maxkey
			node.MaxKey = node.Items[len(node.Items)-1].Key
			return true
		}
	}
	return false
}

func (t *BPTree) itemMoveOrMerge(parent *BPNode, node *BPNode) {
	//获取兄弟结点
	var node1 *BPNode = nil
	var node2 *BPNode = nil
	for i := 0; i < len(parent.Nodes); i++ {
		if parent.Nodes[i] == node {
			if i < len(parent.Nodes)-1 {
				// 父节的右兄弟  --
				node2 = parent.Nodes[i+1]
			} else if i > 0 {
				// 父节的右兄弟左兄弟
				node1 = parent.Nodes[i-1]
			}
			break
		}
	}

	//将左侧结点的记录移动到删除结点
	if node1 != nil && len(node1.Items) > t.halfw {
		//获取被删除节点的父节点的左兄弟的一个孩子节点 -- 需要将他添加到删除节点的位置
		item := node1.Items[len(node1.Items)-1]

		//修改被需要合并元素的item（删除一个） 重新记录最大的Maxkey
		//此时被删除的父节点的左兄弟剩下的节点 -- 不包括item节点 因为 [0 : len(node1.Items)-1] 不包括最后一一个元素
		node1.Items = node1.Items[0 : len(node1.Items)-1]
		node1.MaxKey = node1.Items[len(node1.Items)-1].Key
		//将从左边移动过来的节点添加到node.Items的左边
		node.Items = append([]BPItem{item}, node.Items...)
		return
	}

	//将右侧结点的记录移动到删除结点
	if node2 != nil && len(node2.Items) > t.halfw {
		// 右边的第一个元素需要移动到被删除元素位置
		item := node2.Items[0]

		//修改需要合并元素的node  -- item 和 maxKey
		node2.Items = node1.Items[1:]
		node.Items = append(node.Items, item)
		node.MaxKey = node.Items[len(node.Items)-1].Key
		return
	}

	//与左侧结点进行合并 -- 左侧节点不够不能合并到删除元素位置
	if node1 != nil && len(node1.Items)+len(node.Items) <= t.width {
		// 将被删除的节点合并到左侧节点的右边
		node1.Items = append(node1.Items, node.Items...)
		// 移动链接指针
		node1.Next = node.Next
		//叶子节点本身就是从小到大，因此maxKey取最后一个即可
		node1.MaxKey = node1.Items[len(node1.Items)-1].Key
		//	此时node节点已经合并到左边节点，删除node节点
		parent.deleteChild(node)
		return
	}

	//与右侧结点进行合并
	if node2 != nil && len(node2.Items)+len(node.Items) <= t.width {
		//将被删除节点合并到右侧节点的左侧
		node.Items = append(node.Items, node2.Items...)
		node.Next = node2.Next
		node.MaxKey = node.Items[len(node.Items)-1].Key
		parent.deleteChild(node2)
		return
	}
}

// 索引节点合并或者移动
func (t *BPTree) childMoveOrMerge(parent *BPNode, node *BPNode) {
	if parent == nil {
		return
	}

	//获取兄弟结点
	var node1 *BPNode = nil
	var node2 *BPNode = nil
	for i := 0; i < len(parent.Nodes); i++ {
		if parent.Nodes[i] == node {
			if i < len(parent.Nodes)-1 {
				node2 = parent.Nodes[i+1]
			} else if i > 0 {
				node1 = parent.Nodes[i-1]
			}
			break
		}
	}

	//将左侧结点的子结点移动到删除结点
	if node1 != nil && len(node1.Nodes) > t.halfw {
		item := node1.Nodes[len(node1.Nodes)-1]
		node1.Nodes = node1.Nodes[0 : len(node1.Nodes)-1]
		node.Nodes = append([]*BPNode{item}, node.Nodes...)
		return
	}

	//将右侧结点的子结点移动到删除结点
	if node2 != nil && len(node2.Nodes) > t.halfw {
		item := node2.Nodes[0]
		node2.Nodes = node1.Nodes[1:]
		node.Nodes = append(node.Nodes, item)
		return
	}

	//与左侧结点进行合并
	if node1 != nil && len(node1.Nodes)+len(node.Nodes) <= t.width {
		node1.Nodes = append(node1.Nodes, node.Nodes...)
		parent.deleteChild(node)
		return
	}

	//与右侧结点进行合并
	if node2 != nil && len(node2.Nodes)+len(node.Nodes) <= t.width {
		node.Nodes = append(node.Nodes, node2.Nodes...)
		parent.deleteChild(node2)
		return
	}
}

func (node *BPNode) deleteChild(child *BPNode) bool {
	num := len(node.Nodes)
	for i := 0; i < num; i++ {
		if node.Nodes[i] == child {
			copy(node.Nodes[i:], node.Nodes[i+1:])
			node.Nodes = node.Nodes[0 : len(node.Nodes)-1]
			node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey
			return true
		}
	}
	return false
}
