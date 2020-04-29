package main

import (
	"math/rand"
	"time"
	"fmt"

	//"github.com/Masterminds/glide/tree"
)

func init(){  //优先于main函数执行
	rand.Seed(time.Now().UnixNano())
	fmt.Println("时间随机数已经设定")
}
//判断a<b
type LessFunc func(a,b interface{})bool
//判断是否重叠
type OverlapFunc func(a,b interface{})bool
type Key interface {}
type Item interface {}

type Node struct {
	key  Key  //输入的值
	item Item
	priority int //优先级
	left *Node
	right *Node //左右孩子节点
}
//新建一个节点
func NewNode (key Key,item Item,priority int)*Node{
	res:=new(Node)
	res.key=key
	res.item=item
	res.priority=priority
	return res
}
type Tree struct{
	less LessFunc //函数指针，小于
	overlap  OverlapFunc
	count int  //数量
	root *Node//根节点
}
//新建一棵树
func NewTree(lessfn LessFunc)*Tree{
	t:=new(Tree)
	t.less=lessfn
	return t
}
//33333
//新建重叠树，每个节点都是一个队列，提供函数处理
func NewOverlapTree(lessfn LessFunc,overlap OverlapFunc  )*Tree{
	t:=new(Tree)
	t.less=lessfn
	t.overlap=overlap
	return t
}
//重置0
func (t*Tree)Reset(){
	t.root=nil
	t.count=0
}
//返回长度
func (t*Tree)Len ()int{
	return t.count
}
//抓取数据
func  (t*Tree)Get(key Key)Item{
	return  t.get(t.root,key)
}
func  (t*Tree)get(node*Node,key Key)Item{
	if node==nil{
		return nil
	}
	if t.less(key,node.key){
		return t.get(node.left,key) //左边递归
	}
	if t.less(node.key,key){
		return t.get(node.right,key)//右边递归
	}
	return node.item //相等


}
//判断是否存在
func (t*Tree) Exists(key Key)bool{
	return t.exists(t.root,key)
}

func (t*Tree) exists(node*Node,key Key)bool{
	if node==nil{
		return false
	}
	if t.less(key,node.key){
		return t.exists(node.left,key) //左边递归
	}
	if t.less(node.key,key){
		return  t.exists(node.right,key) //右边递归
	}
	return true

}
func (t*Tree)Insert(key Key,item Item){
	priority:=rand.Int()//提取优先级
	t.root=t.insert(t.root,key,item,priority)//插入数据

}

func (t*Tree)insert(node*Node,key Key,item Item,priority int )*Node{
	if node==nil{
		t.count++
		return NewNode(key,item,priority)
	}
	if t.less(key,node.key){
		node.left=t.insert(node.left,key,item,priority)
		if node.left.priority<node.priority{
			//左旋
			return t.leftRotate(node)
		}
		return node
	}
	if t.less(node.key,key){
		node.right=t.insert(node.right,key,item,priority)
		if node.right.priority<node.priority{
			//右旋
			return t.rightRotate(node)
		}
		return node
	}

	node.item=item//更新数据
	return node

}
//左旋
func (t*Tree)leftRotate(node*Node)*Node{
	res:=node.left
	x:=res.right
	res.right=node
	node.left=x
	return res
}
//右旋
func (t*Tree)rightRotate(node*Node)*Node{
	res:=node.right
	x:=res.left
	res.left=node
	node.right=x
	return res
}
//切割
func (t*Tree)Split(key Key)(*Node,*Node){
	inserted:=t.insert(t.root,key,nil,-1)
	return inserted.left,inserted.right
}
//归并，
func (t *Tree)Merge(left,right*Node)*Node{
	if left==nil{
		return right
	}
	if right==nil{
		return left
	}
	if left.priority<right.priority{
		res:=left
		x:=left.right
		res.right=t.Merge(x,right)//归并
	}
	res:=left
	x:=right.left
	res.left=t.Merge(x,left)//返回结果
	return res



}

func (t *Tree)Delete(key Key){
	if t.Exists(key)==false{
		return
	}
	t.root=t.delete(t.root,key)//删除
}
func (t *Tree)delete(node *Node,key Key)*Node{
	if node==nil{
		return nil
	}

	if t.less(key,node.key){
		res:=node
		x:=node.left
		res.left=t.delete(x,key)
		return res
	}
	if t.less(node.key,key){
		res:=node
		x:=node.right
		res.right=t.delete(x,key)
		return res
	}


	t.count-- //删除当前节点
	return t.Merge(node.left,node.right)


}
//求高度
func (t *Tree)Height(key Key)int{
	return t.height(t.root,key)
}

func (t *Tree)height(node *Node,key Key)int{
	if node==nil{
		return 0
	}
	if t.less(key,node.key){
		depth:=t.height(node.left,key)
		return depth+1
	}
	if t.less(node.key,key){
		depth:=t.height(node.right,key)
		return depth+1
	}
	return 0
}
//循环所有节点
func (t*Tree)IterAscend()<-chan Item{
	c:=make(chan Item)
	go func() {
		iterateInorder(t.root,c)//前序遍历
		close(c)
	}()
	return c
}
func iterateInorder(h*Node,c chan <-Item){
	if h==nil{
		return
	}
	iterateInorder(h.left,c)
	c<-h.item
	iterateInorder(h.right,c)
}
//循环所有节点
func (t*Tree)IterKeyAscend()<-chan Key{
	c:=make(chan Key)
	go func() {
		iteratekeyInorder(t.root,c)//前序遍历
		close(c)
	}()
	return c
}
func iteratekeyInorder(h*Node,c chan <-Key){
	if h==nil{
		return
	}
	iteratekeyInorder(h.left,c)
	c<-h.key
	iteratekeyInorder(h.right,c)

	//c<-h.key前序
	//iteratekeyInorder(h.left,c)
	//iteratekeyInorder(h.right,c)


	//iteratekeyInorder(h.left,c)
	//iteratekeyInorder(h.right,c)
	//c<-h.key后序
}

func (t*Tree)Min()Item{
	return min(t.root)
}
func min(h*Node)Item{
	if h==nil{
		return nil
	}
	if h.left==nil{
		return h.item
	}
	return min(h.left)
}

func (t*Tree)Max()Item{
	return max(t.root)
}
func max(h*Node)Item{
	if h==nil{
		return nil
	}
	if h.right==nil{
		return h.item
	}
	return min(h.right)
}

func (t*Tree)IterateOverlap(key Key)<-chan Item{
	c:=make(chan Item)
	go func (){
		if t.overlap!=nil{
			t.iterateOverlap(t.root,key,c)
		}
		close(c)
	}()
	return c
}
func (t*Tree)iterateOverlap(h*Node,  key Key,c chan<- Item){
	if h==nil{
		return
	}
	lessThanLower:=t.overlap(h.key,key)
	greaterThanLower:=t.overlap(h.key,key)
	if !lessThanLower{
		t.iterateOverlap(h.left,key,c)
	}
	if!lessThanLower&& !greaterThanLower{
		c<-h.item
	}

	if !greaterThanLower{
		t.iterateOverlap(h.right,key,c)
	}

}
//协助比较大小
func StringLess(p,q interface{})bool{
	return p.(string)<q.(string )
}
func IntLess(p,q interface{})bool{
	return p.(int)<q.(int )
}




func main2(){
	//测试平衡
	tree:=NewTree(IntLess)
	for i:=0;i<10000;i++{
		tree.Insert(i,i)
	}
	for i:=0;i<10000;i+=1000{
		fmt.Println(tree.Height(i))
	}


}
func main1(){

	tree:=NewTree(StringLess)
	tree.Insert("xyz","123")
	tree.Insert("mnb","123")
	tree.Insert("abc","123")
	tree.Insert("bcd","123")
	fmt.Println(tree.Get("bcd"))


}
func main24(){
	tree:=NewTree(IntLess)
	n:=100
	for i:=0;i<100;i++{
		tree.Insert(n-i,true)
	}
	//插入数据，在不断取出数据，满足堆的性质
	c:=tree.IterKeyAscend()
	for j,item:=1,<-c;item!=nil;j,item=j+1,<-c{
		fmt.Println(item.(int))
	}
}
func main(){
	tree1:=NewTree(IntLess)
	n:=100
	for i:=0;i<100;i++{
		tree1.Insert(n-i,true)
	}


	//tree2:=NewTree(IntLess)
	n=300
	for i:=0;i<100;i++{
		tree1.Insert(n-i,true)
	}
	//tree1.root=tree1.Merge(tree1.root,tree2.root)
	//插入数据，在不断取出数据，满足堆的性质
	c:=tree1.IterKeyAscend()
	for j,item:=1,<-c;item!=nil;j,item=j+1,<-c{
		fmt.Println(item.(int))
	}

	
}
