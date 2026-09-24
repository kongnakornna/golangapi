# 📘 คู่มือสถาปัตยกรรมซอฟต์แวร์และวิศวกรรมซอฟต์แวร์ฉบับสมบูรณ์
# (The Complete Software Architecture & Engineering Master Guide)
## เวอร์ชัน Full Code 2.0 (เมษายน 2026)

> **เอกสารต้นแบบเชิงลึกพร้อมโค้ดเต็ม** สำหรับนักพัฒนาซอฟต์แวร์ทุกระดับ
> ครอบคลุมตั้งแต่พื้นฐานจนถึงโปรเจกต์จริง `icmongolang` พร้อมโค้ดที่รันได้

---

# สารบัญ (Table of Contents)

## ภาคที่ 1: พื้นฐาน
1. บทนำสู่สถาปัตยกรรมซอฟต์แวร์
2. Data Structures และ Algorithms (Full Code)
3. Big O Notation (Full Code)
4. Functional Programming (Full Code)
5. OOP (Full Code)
6. System Design
7. Software Architecture
8. CI/CD
9. Cloud Platform AWS
10. Domain-Driven Design
11. Clean Architecture

## ภาคที่ 2: การจัดการข้อมูล
12. CRUD (Full Code)
13. ORM (Full Code)
14. TypeORM (Full Code)
15. Schema (Full SQL)
16. Entity (Full Code)
17. DTO และ Validation (Full Code)
18. Transaction (Full Code)
19. Cache (Full Code)

## ภาคที่ 3: icmongolang Project
20. ภาพรวมโปรเจกต์
21. Full Code: Domain Layer
22. Full Code: Application Layer
23. Full Code: Infrastructure Layer
24. Full Code: Interface Layer
25. Full Code: Bootstrap / Entry Points
26. Full Code: WebSocket Server
27. Full Code: Worker Pool
28. Full Code: MQTT Client
29. Full Code: IoT Module (Complete)
30. Full Code: Auth Module (Complete)

## ภาคที่ 4: DevOps และ EDA
31. Event-Driven Design
32. Code Review Workflow
33. Monitoring Tools
34. POC & ROI

## ภาคที่ 5: RCA และ Checklists
35. RCA Diagrams
36. Best Practices & Checklists

## ภาคผนวก
A-F: Templates, Conventions, References

---

# ภาคที่ 1: พื้นฐาน

---

## 1. บทนำสู่สถาปัตยกรรมซอฟต์แวร์

### 1.1 หลักการสำคัญ

```
┌───────────────────────────────────────────┐
│  Interface Layer   (HTTP, WS, CLI)        │
├───────────────────────────────────────────┤
│  Infrastructure    (DB, Kafka, ES, Redis) │
├───────────────────────────────────────────┤
│  Application       (Use Cases)            │
├───────────────────────────────────────────┤
│  Domain            (Entities, VOs) ❤️     │
└───────────────────────────────────────────┘
```

### 1.2 Dependency Rule

| Layer | รู้จัก | ห้ามรู้จัก |
|-------|-------|-----------|
| Domain | stdlib, uuid | ทุก layer อื่น |
| Application | Domain | Infrastructure, Interface |
| Infrastructure | Domain, Application | Interface |
| Interface | ทุก layer | – |

---

## 2. Data Structures และ Algorithms (Full Code)

### 2.1 Linked List

```go
package main

import "fmt"

type Node struct {
    Value int
    Next  *Node
}

type LinkedList struct {
    Head *Node
    Size int
}

func (l *LinkedList) Append(value int) {
    newNode := &Node{Value: value}
    if l.Head == nil {
        l.Head = newNode
    } else {
        current := l.Head
        for current.Next != nil {
            current = current.Next
        }
        current.Next = newNode
    }
    l.Size++
}

func (l *LinkedList) Prepend(value int) {
    newNode := &Node{Value: value, Next: l.Head}
    l.Head = newNode
    l.Size++
}

func (l *LinkedList) Delete(value int) bool {
    if l.Head == nil {
        return false
    }
    if l.Head.Value == value {
        l.Head = l.Head.Next
        l.Size--
        return true
    }
    current := l.Head
    for current.Next != nil {
        if current.Next.Value == value {
            current.Next = current.Next.Next
            l.Size--
            return true
        }
        current = current.Next
    }
    return false
}

func (l *LinkedList) Print() {
    current := l.Head
    for current != nil {
        fmt.Printf("%d -> ", current.Value)
        current = current.Next
    }
    fmt.Println("nil")
}

func main() {
    list := &LinkedList{}
    list.Append(1)
    list.Append(2)
    list.Append(3)
    list.Prepend(0)
    list.Print() // 0 -> 1 -> 2 -> 3 -> nil
    list.Delete(2)
    list.Print() // 0 -> 1 -> 3 -> nil
}
```

### 2.2 Stack และ Queue

```go
package main

import "fmt"

// Stack (LIFO)
type Stack struct {
    items []int
}

func (s *Stack) Push(item int) {
    s.items = append(s.items, item)
}

func (s *Stack) Pop() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    last := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return last, true
}

func (s *Stack) Peek() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    return s.items[len(s.items)-1], true
}

// Queue (FIFO)
type Queue struct {
    items []int
}

func (q *Queue) Enqueue(item int) {
    q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (int, bool) {
    if len(q.items) == 0 {
        return 0, false
    }
    first := q.items[0]
    q.items = q.items[1:]
    return first, true
}

func main() {
    s := &Stack{}
    s.Push(1)
    s.Push(2)
    s.Push(3)
    v, _ := s.Pop()
    fmt.Println("Popped from stack:", v) // 3

    q := &Queue{}
    q.Enqueue(1)
    q.Enqueue(2)
    q.Enqueue(3)
    v, _ = q.Dequeue()
    fmt.Println("Dequeued from queue:", v) // 1
}
```

### 2.3 Binary Search Tree

```go
package main

import "fmt"

type BSTNode struct {
    Value int
    Left  *BSTNode
    Right *BSTNode
}

type BST struct {
    Root *BSTNode
}

func (b *BST) Insert(value int) {
    b.Root = insertNode(b.Root, value)
}

func insertNode(node *BSTNode, value int) *BSTNode {
    if node == nil {
        return &BSTNode{Value: value}
    }
    if value < node.Value {
        node.Left = insertNode(node.Left, value)
    } else if value > node.Value {
        node.Right = insertNode(node.Right, value)
    }
    return node
}

func (b *BST) Search(value int) bool {
    return searchNode(b.Root, value)
}

func searchNode(node *BSTNode, value int) bool {
    if node == nil {
        return false
    }
    if value == node.Value {
        return true
    }
    if value < node.Value {
        return searchNode(node.Left, value)
    }
    return searchNode(node.Right, value)
}

func (b *BST) InOrder() {
    inOrder(b.Root)
    fmt.Println()
}

func inOrder(node *BSTNode) {
    if node == nil {
        return
    }
    inOrder(node.Left)
    fmt.Printf("%d ", node.Value)
    inOrder(node.Right)
}

func main() {
    bst := &BST{}
    bst.Insert(50)
    bst.Insert(30)
    bst.Insert(70)
    bst.Insert(20)
    bst.Insert(40)
    bst.Insert(60)
    bst.Insert(80)
    bst.InOrder() // 20 30 40 50 60 70 80
    fmt.Println("Search 40:", bst.Search(40))
    fmt.Println("Search 100:", bst.Search(100))
}
```

### 2.4 Hash Table

```go
package main

import "fmt"

type HashEntry struct {
    Key   string
    Value interface{}
}

type HashTable struct {
    buckets [][]HashEntry
    size    int
}

func NewHashTable(size int) *HashTable {
    return &HashTable{
        buckets: make([][]HashEntry, size),
        size:    size,
    }
}

func (h *HashTable) hash(key string) int {
    hash := 0
    for _, c := range key {
        hash = (hash*31 + int(c)) % h.size
    }
    return hash
}

func (h *HashTable) Put(key string, value interface{}) {
    idx := h.hash(key)
    for i, entry := range h.buckets[idx] {
        if entry.Key == key {
            h.buckets[idx][i].Value = value
            return
        }
    }
    h.buckets[idx] = append(h.buckets[idx], HashEntry{Key: key, Value: value})
}

func (h *HashTable) Get(key string) (interface{}, bool) {
    idx := h.hash(key)
    for _, entry := range h.buckets[idx] {
        if entry.Key == key {
            return entry.Value, true
        }
    }
    return nil, false
}

func (h *HashTable) Delete(key string) {
    idx := h.hash(key)
    for i, entry := range h.buckets[idx] {
        if entry.Key == key {
            h.buckets[idx] = append(h.buckets[idx][:i], h.buckets[idx][i+1:]...)
            return
        }
    }
}

func main() {
    ht := NewHashTable(10)
    ht.Put("name", "John")
    ht.Put("age", 30)
    ht.Put("city", "Bangkok")

    if v, ok := ht.Get("name"); ok {
        fmt.Println("name:", v)
    }
    if v, ok := ht.Get("age"); ok {
        fmt.Println("age:", v)
    }
    ht.Delete("city")
    if _, ok := ht.Get("city"); !ok {
        fmt.Println("city deleted")
    }
}
```

### 2.5 Suffix Automaton

```go
package main

import "fmt"

type State struct {
    Len  int
    Link int
    Next map[byte]int
}

type SuffixAutomaton struct {
    States []State
    Last   int
}

func NewSuffixAutomaton() *SuffixAutomaton {
    return &SuffixAutomaton{
        States: []State{{Len: 0, Link: -1, Next: make(map[byte]int)}},
        Last:   0,
    }
}

func (sa *SuffixAutomaton) Extend(c byte) {
    cur := len(sa.States)
    sa.States = append(sa.States, State{
        Len:  sa.States[sa.Last].Len + 1,
        Link: 0,
        Next: make(map[byte]int),
    })

    p := sa.Last
    for p != -1 && sa.States[p].Next[c] == 0 {
        if _, exists := sa.States[p].Next[c]; !exists {
            sa.States[p].Next[c] = cur
        }
        p = sa.States[p].Link
    }

    if p == -1 {
        sa.States[cur].Link = 0
    } else {
        q := sa.States[p].Next[c]
        if sa.States[p].Len+1 == sa.States[q].Len {
            sa.States[cur].Link = q
        } else {
            clone := len(sa.States)
            sa.States = append(sa.States, State{
                Len:  sa.States[p].Len + 1,
                Link: sa.States[q].Link,
                Next: make(map[byte]int),
            })
            for k, v := range sa.States[q].Next {
                sa.States[clone].Next[k] = v
            }
            for p != -1 && sa.States[p].Next[c] == q {
                sa.States[p].Next[c] = clone
                p = sa.States[p].Link
            }
            sa.States[q].Link = clone
            sa.States[cur].Link = clone
        }
    }
    sa.Last = cur
}

func (sa *SuffixAutomaton) Build(s string) {
    for i := 0; i < len(s); i++ {
        sa.Extend(s[i])
    }
}

func (sa *SuffixAutomaton) Contains(pattern string) bool {
    state := 0
    for i := 0; i < len(pattern); i++ {
        next, ok := sa.States[state].Next[pattern[i]]
        if !ok {
            return false
        }
        state = next
    }
    return true
}

func main() {
    sa := NewSuffixAutomaton()
    sa.Build("banana")
    fmt.Println("Contains 'ana':", sa.Contains("ana"))
    fmt.Println("Contains 'ban':", sa.Contains("ban"))
    fmt.Println("Contains 'xyz':", sa.Contains("xyz"))
}
```

### 2.6 Palindromic Automaton

```go
package main

import "fmt"

type PALNode struct {
    Next [26]int
    Len  int
    Link int
    Cnt  int
}

type PalindromicAutomaton struct {
    Tree []PALNode
    S    []byte
    Last int
}

func NewPalindromicAutomaton() *PalindromicAutomaton {
    pa := &PalindromicAutomaton{}
    // Node 0: root with len -1
    n0 := PALNode{Len: -1, Link: 0, Cnt: 0}
    // Node 1: root with len 0
    n1 := PALNode{Len: 0, Link: 0, Cnt: 0}
    pa.Tree = []PALNode{n0, n1}
    pa.Last = 1
    return pa
}

func (pa *PalindromicAutomaton) AddChar(c byte) {
    idx := int(c - 'a')
    pa.S = append(pa.S, c)
    pos := len(pa.S) - 1
    cur := pa.Last

    for {
        curlen := pa.Tree[cur].Len
        if pos-1-curlen >= 0 && pa.S[pos-1-curlen] == c {
            break
        }
        cur = pa.Tree[cur].Link
    }

    if pa.Tree[cur].Next[idx] != 0 {
        pa.Last = pa.Tree[cur].Next[idx]
        pa.Tree[pa.Last].Cnt++
        return
    }

    now := len(pa.Tree)
    newNode := PALNode{Len: pa.Tree[cur].Len + 2, Cnt: 0}
    pa.Tree = append(pa.Tree, newNode)
    pa.Tree[cur].Next[idx] = now

    if pa.Tree[now].Len == 1 {
        pa.Tree[now].Link = 1
    } else {
        link := pa.Tree[cur].Link
        for {
            linklen := pa.Tree[link].Len
            if pos-1-linklen >= 0 && pa.S[pos-1-linklen] == c {
                break
            }
            link = pa.Tree[link].Link
        }
        pa.Tree[now].Link = pa.Tree[link].Next[idx]
    }
    pa.Last = now
    pa.Tree[pa.Last].Cnt++
}

func main() {
    pa := NewPalindromicAutomaton()
    for i := 0; i < len("ababa"); i++ {
        pa.AddChar("ababa"[i])
    }
    total := 0
    for i := 2; i < len(pa.Tree); i++ {
        total += pa.Tree[i].Cnt
    }
    fmt.Println("Total palindromic substrings:", total)
}
```

### 2.7 HLD + Lazy Segment Tree

```go
package main

import "fmt"

type HLD struct {
    adj        [][]int
    parent     []int
    depth      []int
    heavy      []int
    head       []int
    pos        []int
    curPos     int
    n          int
}

func NewHLD(n int) *HLD {
    h := &HLD{
        adj:    make([][]int, n),
        parent: make([]int, n),
        depth:  make([]int, n),
        heavy:  make([]int, n),
        head:   make([]int, n),
        pos:    make([]int, n),
        n:      n,
    }
    for i := 0; i < n; i++ {
        h.parent[i] = -1
        h.heavy[i] = -1
    }
    return h
}

func (h *HLD) AddEdge(u, v int) {
    h.adj[u] = append(h.adj[u], v)
    h.adj[v] = append(h.adj[v], u)
}

func (h *HLD) dfs(u int) int {
    size := 1
    maxCSize := 0
    for _, v := range h.adj[u] {
        if v != h.parent[u] {
            h.parent[v] = u
            h.depth[v] = h.depth[u] + 1
            cSize := h.dfs(v)
            size += cSize
            if cSize > maxCSize {
                maxCSize = cSize
                h.heavy[u] = v
            }
        }
    }
    return size
}

func (h *HLD) decompose(u, hd int) {
    h.head[u] = hd
    h.pos[u] = h.curPos
    h.curPos++
    if h.heavy[u] != -1 {
        h.decompose(h.heavy[u], hd)
    }
    for _, v := range h.adj[u] {
        if v != h.parent[u] && v != h.heavy[u] {
            h.decompose(v, v)
        }
    }
}

func (h *HLD) Init() {
    h.dfs(0)
    h.decompose(0, 0)
}

func (h *HLD) LCA(u, v int) int {
    for h.head[u] != h.head[v] {
        if h.depth[h.head[u]] > h.depth[h.head[v]] {
            u = h.parent[h.head[u]]
        } else {
            v = h.parent[h.head[v]]
        }
    }
    if h.depth[u] > h.depth[v] {
        return v
    }
    return u
}

func main() {
    h := NewHLD(7)
    h.AddEdge(0, 1)
    h.AddEdge(0, 2)
    h.AddEdge(1, 3)
    h.AddEdge(1, 4)
    h.AddEdge(2, 5)
    h.AddEdge(2, 6)
    h.Init()
    fmt.Println("LCA(3,4):", h.LCA(3, 4))
    fmt.Println("LCA(3,5):", h.LCA(3, 5))
    fmt.Println("LCA(5,6):", h.LCA(5, 6))
}
```

---

## 3. Big O Notation (Full Code)

```go
package main

import "fmt"

// O(1) - Constant Time
func getFirst(arr []int) int {
    return arr[0]
}

// O(n) - Linear Time
func findMax(arr []int) int {
    max := arr[0]
    for i := 1; i < len(arr); i++ {
        if arr[i] > max {
            max = arr[i]
        }
    }
    return max
}

// O(n²) - Quadratic Time
func bubbleSort(arr []int) []int {
    n := len(arr)
    for i := 0; i < n; i++ {
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
            }
        }
    }
    return arr
}

// O(log n) - Logarithmic Time
func binarySearch(arr []int, target int) int {
    left, right := 0, len(arr)-1
    for left <= right {
        mid := (left + right) / 2
        if arr[mid] == target {
            return mid
        }
        if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1
}

// O(n log n) - Linearithmic Time
func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])
    right := mergeSort(arr[mid:])
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}

// O(2ⁿ) - Exponential (Fibonacci naive)
func fibNaive(n int) int {
    if n <= 1 {
        return n
    }
    return fibNaive(n-1) + fibNaive(n-2)
}

// O(n) - Dynamic Programming
func fibDP(n int) int {
    if n <= 1 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}

func main() {
    arr := []int{5, 2, 8, 1, 9, 3}
    fmt.Println("First:", getFirst(arr))
    fmt.Println("Max:", findMax(arr))
    fmt.Println("Sorted:", bubbleSort([]int{5, 2, 8, 1, 9, 3}))
    fmt.Println("Binary Search 8:", binarySearch([]int{1, 2, 3, 5, 8, 9}, 8))
    fmt.Println("Merge Sort:", mergeSort([]int{5, 2, 8, 1, 9, 3}))
    fmt.Println("Fib(10) naive:", fibNaive(10))
    fmt.Println("Fib(10) DP:", fibDP(10))
}
```

---

## 4. Functional Programming (Full Code)

```go
package main

import "fmt"

// Pure Function
func add(a, b int) int {
    return a + b
}

// Higher-Order Functions
func Map(arr []int, fn func(int) int) []int {
    result := make([]int, len(arr))
    for i, v := range arr {
        result[i] = fn(v)
    }
    return result
}

func Filter(arr []int, fn func(int) bool) []int {
    result := []int{}
    for _, v := range arr {
        if fn(v) {
            result = append(result, v)
        }
    }
    return result
}

func Reduce(arr []int, fn func(int, int) int, initial int) int {
    acc := initial
    for _, v := range arr {
        acc = fn(acc, v)
    }
    return acc
}

// Function Composition
func Compose(f, g func(int) int) func(int) int {
    return func(x int) int {
        return f(g(x))
    }
}

// Currying
func Adder(a int) func(int) int {
    return func(b int) int {
        return a + b
    }
}

// Immutability - returns new slice
func AppendImmutable(arr []int, item int) []int {
    newArr := make([]int, len(arr)+1)
    copy(newArr, arr)
    newArr[len(arr)] = item
    return newArr
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}

    doubled := Map(numbers, func(n int) int { return n * 2 })
    fmt.Println("Doubled:", doubled) // [2 4 6 8 10]

    evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
    fmt.Println("Evens:", evens) // [2 4]

    sum := Reduce(numbers, func(a, b int) int { return a + b }, 0)
    fmt.Println("Sum:", sum) // 15

    addOne := func(x int) int { return x + 1 }
    double := func(x int) int { return x * 2 }
    addOneThenDouble := Compose(double, addOne)
    fmt.Println("Compose(5):", addOneThenDouble(5)) // 12

    add5 := Adder(5)
    fmt.Println("Add5(3):", add5(3)) // 8

    original := []int{1, 2, 3}
    modified := AppendImmutable(original, 4)
    fmt.Println("Original:", original) // [1 2 3]
    fmt.Println("Modified:", modified) // [1 2 3 4]
}
```

---

## 5. OOP (Full Code)

```go
package main

import (
    "errors"
    "fmt"
    "math"
)

// Encapsulation
type BankAccount struct {
    accountNumber string
    balance       float64
}

func NewBankAccount(accountNumber string, initialBalance float64) *BankAccount {
    return &BankAccount{
        accountNumber: accountNumber,
        balance:       initialBalance,
    }
}

func (b *BankAccount) Deposit(amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    b.balance += amount
    return nil
}

func (b *BankAccount) Withdraw(amount float64) error {
    if amount > b.balance {
        return errors.New("insufficient funds")
    }
    b.balance -= amount
    return nil
}

func (b *BankAccount) GetBalance() float64 {
    return b.balance
}

// Abstraction
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Inheritance (via embedding + interface)
type Circle struct {
    radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.radius
}

type Rectangle struct {
    width, height float64
}

func (r Rectangle) Area() float64 {
    return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.width + r.height)
}

// Polymorphism
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// SOLID: Single Responsibility
type User struct {
    ID    int
    Name  string
    Email string
}

type UserRepository interface {
    Save(u *User) error
    FindByID(id int) (*User, error)
}

type UserValidator struct{}

func (v *UserValidator) Validate(u *User) error {
    if u.Email == "" {
        return errors.New("email required")
    }
    return nil
}

type InMemoryUserRepo struct {
    users map[int]*User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
    return &InMemoryUserRepo{users: make(map[int]*User)}
}

func (r *InMemoryUserRepo) Save(u *User) error {
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepo) FindByID(id int) (*User, error) {
    u, ok := r.users[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return u, nil
}

// SOLID: Open/Closed
type Discount interface {
    Apply(price float64) float64
}

type PercentageDiscount struct {
    Percent float64
}

func (d PercentageDiscount) Apply(price float64) float64 {
    return price * (1 - d.Percent/100)
}

type FixedDiscount struct {
    Amount float64
}

func (d FixedDiscount) Apply(price float64) float64 {
    return price - d.Amount
}

func main() {
    // Encapsulation
    acc := NewBankAccount("12345", 1000)
    _ = acc.Deposit(500)
    _ = acc.Withdraw(200)
    fmt.Printf("Balance: %.2f\n", acc.GetBalance())

    // Polymorphism
    shapes := []Shape{Circle{5}, Rectangle{4, 6}}
    for _, s := range shapes {
        PrintShapeInfo(s)
    }

    // User repo
    repo := NewInMemoryUserRepo()
    user := &User{ID: 1, Name: "John", Email: "john@example.com"}
    validator := &UserValidator{}
    if err := validator.Validate(user); err == nil {
        _ = repo.Save(user)
    }

    // Open/Closed
    discounts := []Discount{
        PercentageDiscount{10},
        FixedDiscount{50},
    }
    for _, d := range discounts {
        fmt.Printf("Discounted price: %.2f\n", d.Apply(1000))
    }
}
```

---

## 6-11. System Design, Architecture, DDD, Clean Architecture

(ดูรายละเอียดใน ภาคที่ 3 ซึ่งเป็นโค้ดจริงจากโปรเจกต์ icmongolang)

---

# ภาคที่ 2: การจัดการข้อมูล

---

## 12. CRUD (Full Code)

### 12.1 Raw SQL CRUD

```sql
-- Create
INSERT INTO users (name, email, age) VALUES ('John', 'john@example.com', 30);

-- Read
SELECT * FROM users WHERE id = 1;
SELECT * FROM users WHERE age > 18 ORDER BY name;
SELECT * FROM users LIMIT 10 OFFSET 20;

-- Update
UPDATE users SET name = 'Jane' WHERE id = 1;

-- Delete
DELETE FROM users WHERE id = 1;

-- Soft Delete
UPDATE users SET deleted_at = NOW() WHERE id = 1;
SELECT * FROM users WHERE deleted_at IS NULL;
```

### 12.2 CRUD ใน Go with GORM

```go
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:100;not null"`
    Email string `gorm:"uniqueIndex;size:255;not null"`
    Age   int
}

func main() {
    dsn := "host=localhost user=postgres password=password dbname=test port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // AutoMigrate
    db.AutoMigrate(&User{})

    // CREATE
    user := User{Name: "John", Email: "john@example.com", Age: 30}
    db.Create(&user)
    fmt.Println("Created:", user.ID)

    // READ
    var found User
    db.First(&found, user.ID)
    fmt.Println("Found:", found.Name)

    var users []User
    db.Where("age > ?", 18).Find(&users)
    fmt.Println("Users over 18:", len(users))

    // UPDATE
    db.Model(&user).Update("name", "Jane")
    db.Model(&user).Updates(User{Name: "Johnny", Age: 31})

    // DELETE
    db.Delete(&user, user.ID)

    // Soft Delete (ต้องมี gorm.DeletedAt)
    // type UserWithSoftDelete struct {
    //     ID uint
    //     gorm.DeletedAt
    // }
    // db.Delete(&user) // soft delete
    // db.Unscoped().Delete(&user) // hard delete
}
```

---

## 13. ORM (Full Code)

```go
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// Entity Relationships
type User struct {
    ID      uint     `gorm:"primaryKey"`
    Name    string   `gorm:"size:100;not null"`
    Email   string   `gorm:"uniqueIndex;size:255"`
    Posts   []Post   `gorm:"foreignKey:UserID"`
    Profile *Profile `gorm:"foreignKey:UserID"`
}

type Profile struct {
    ID     uint   `gorm:"primaryKey"`
    UserID uint   `gorm:"uniqueIndex"`
    Bio    string `gorm:"type:text"`
    Avatar string
}

type Post struct {
    ID       uint   `gorm:"primaryKey"`
    UserID   uint   `gorm:"index"`
    Title    string `gorm:"size:200;not null"`
    Content  string `gorm:"type:text"`
    Tags     []Tag  `gorm:"many2many:post_tags"`
}

type Tag struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"uniqueIndex;size:50"`
    Posts []Post `gorm:"many2many:post_tags"`
}

func main() {
    dsn := "host=localhost user=postgres password=password dbname=test port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    db.AutoMigrate(&User{}, &Profile{}, &Post{}, &Tag{})

    // CREATE with relations
    user := User{
        Name:  "John",
        Email: "john@example.com",
        Profile: &Profile{
            Bio:    "Software Engineer",
            Avatar: "avatar.jpg",
        },
        Posts: []Post{
            {
                Title:   "Hello GORM",
                Content: "This is my first post",
                Tags: []Tag{
                    {Name: "golang"},
                    {Name: "orm"},
                },
            },
        },
    }
    db.Create(&user)

    // READ with Preload
    var found User
    db.Preload("Posts").Preload("Profile").Preload("Posts.Tags").First(&found, user.ID)
    fmt.Printf("User: %s, Posts: %d\n", found.Name, len(found.Posts))
    fmt.Printf("Profile Bio: %s\n", found.Profile.Bio)
    fmt.Printf("First Post Tags: %d\n", len(found.Posts[0].Tags))

    // UPDATE with relations
    db.Model(&found).Association("Posts").Append(&Post{Title: "Second Post"})

    // DELETE with cascade
    db.Select("Posts", "Profile").Delete(&found)
}
```

---

## 14. TypeORM (Full Code)

### 14.1 Entity Definition

```typescript
// src/entity/User.ts
import {
    Entity,
    PrimaryGeneratedColumn,
    Column,
    CreateDateColumn,
    UpdateDateColumn,
    DeleteDateColumn,
    OneToMany,
    OneToOne,
    Index,
} from "typeorm";
import { Post } from "./Post";
import { Profile } from "./Profile";

@Entity("users")
@Index(["email"], { unique: true })
export class User {
    @PrimaryGeneratedColumn("uuid")
    id: string;

    @Column({ length: 100 })
    @Index()
    firstName: string;

    @Column({ length: 100 })
    lastName: string;

    @Column({ unique: true })
    email: string;

    @Column({ select: false })
    password: string;

    @Column({ default: true })
    isActive: boolean;

    @OneToOne(() => Profile, (profile) => profile.user, { cascade: true })
    profile: Profile;

    @OneToMany(() => Post, (post) => post.author)
    posts: Post[];

    @CreateDateColumn()
    createdAt: Date;

    @UpdateDateColumn()
    updatedAt: Date;

    @DeleteDateColumn()
    deletedAt: Date | null;
}
```

### 14.2 Post Entity with Many-to-Many

```typescript
// src/entity/Post.ts
import {
    Entity,
    PrimaryGeneratedColumn,
    Column,
    ManyToOne,
    ManyToMany,
    JoinTable,
    CreateDateColumn,
    Index,
} from "typeorm";
import { User } from "./User";
import { Category } from "./Category";

@Entity("posts")
export class Post {
    @PrimaryGeneratedColumn()
    id: number;

    @Column({ length: 200 })
    @Index()
    title: string;

    @Column("text")
    content: string;

    @Column({ default: 0 })
    views: number;

    @ManyToOne(() => User, (user) => user.posts)
    @JoinTable({ name: "post_author" })
    author: User;

    @Column()
    authorId: string;

    @ManyToMany(() => Category, (category) => category.posts)
    @JoinTable({
        name: "post_categories",
        joinColumn: { name: "post_id", referencedColumnName: "id" },
        inverseJoinColumn: { name: "category_id", referencedColumnName: "id" },
    })
    categories: Category[];

    @CreateDateColumn()
    createdAt: Date;
}
```

### 14.3 Data Source Configuration

```typescript
// src/data-source.ts
import "reflect-metadata";
import { DataSource } from "typeorm";
import { User } from "./entity/User";
import { Post } from "./entity/Post";
import { Profile } from "./entity/Profile";
import { Category } from "./entity/Category";

export const AppDataSource = new DataSource({
    type: "postgres",
    host: process.env.DB_HOST || "localhost",
    port: parseInt(process.env.DB_PORT || "5432"),
    username: process.env.DB_USERNAME || "postgres",
    password: process.env.DB_PASSWORD || "password",
    database: process.env.DB_NAME || "myapp",
    synchronize: process.env.NODE_ENV === "development",
    logging: process.env.NODE_ENV === "development",
    entities: [User, Post, Profile, Category],
    migrations: ["src/migration/*.ts"],
    subscribers: [],
    poolSize: 10,
    extra: {
        max: 20,
        min: 5,
        idleTimeoutMillis: 30000,
    },
});
```

### 14.4 Repository Usage

```typescript
// src/services/UserService.ts
import { Repository, DataSource } from "typeorm";
import { User } from "../entity/User";

export class UserService {
    private userRepository: Repository<User>;

    constructor(dataSource: DataSource) {
        this.userRepository = dataSource.getRepository(User);
    }

    // CREATE
    async createUser(data: Partial<User>): Promise<User> {
        const user = this.userRepository.create(data);
        return await this.userRepository.save(user);
    }

    // READ
    async getUserById(id: string): Promise<User | null> {
        return await this.userRepository.findOne({
            where: { id },
            relations: ["posts", "profile"],
        });
    }

    async getAllUsers(page = 1, limit = 10): Promise<[User[], number]> {
        return await this.userRepository.findAndCount({
            skip: (page - 1) * limit,
            take: limit,
            order: { createdAt: "DESC" },
        });
    }

    // UPDATE
    async updateUser(id: string, data: Partial<User>): Promise<User | null> {
        await this.userRepository.update(id, data);
        return this.getUserById(id);
    }

    // DELETE
    async deleteUser(id: string): Promise<boolean> {
        const result = await this.userRepository.delete(id);
        return (result.affected ?? 0) > 0;
    }

    // SOFT DELETE
    async softDeleteUser(id: string): Promise<boolean> {
        const result = await this.userRepository.softDelete(id);
        return (result.affected ?? 0) > 0;
    }

    // QueryBuilder - Complex
    async searchUsers(query: string): Promise<User[]> {
        return await this.userRepository
            .createQueryBuilder("user")
            .leftJoinAndSelect("user.posts", "post")
            .where("user.firstName ILIKE :q", { q: `%${query}%` })
            .orWhere("user.lastName ILIKE :q", { q: `%${query}%` })
            .orderBy("user.createdAt", "DESC")
            .take(20)
            .getMany();
    }

    // Transaction
    async transferPost(postId: number, fromId: string, toId: string): Promise<void> {
        await this.userRepository.manager.transaction(async (manager) => {
            const fromUser = await manager.findOne(User, { where: { id: fromId } });
            const toUser = await manager.findOne(User, { where: { id: toId } });
            if (!fromUser || !toUser) throw new Error("User not found");
            // Update posts
            await manager.update("posts", { authorId: fromId }, { authorId: toId });
        });
    }
}
```

---

## 15. Schema (Full SQL)

```sql
-- =====================================================
-- Complete E-Commerce Schema
-- =====================================================

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role);

-- Profiles (1-to-1)
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    avatar_url VARCHAR(500),
    date_of_birth DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Categories
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_parent ON categories(parent_id);

-- Products
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category_id UUID REFERENCES categories(id),
    sku VARCHAR(50) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_name_fts ON products USING gin(to_tsvector('english', name));

-- Orders
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total DECIMAL(12, 2) NOT NULL CHECK (total >= 0),
    shipping_address JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);

-- Order Items
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10, 2) NOT NULL CHECK (unit_price >= 0),
    subtotal DECIMAL(12, 2) GENERATED ALWAYS AS (quantity * unit_price) STORED
);

CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_product ON order_items(product_id);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- View
CREATE VIEW user_orders_summary AS
SELECT
    u.id AS user_id,
    u.email,
    COUNT(DISTINCT o.id) AS order_count,
    COALESCE(SUM(o.total), 0) AS total_spent
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.email;
```

---

## 16. Entity (Full Code)

```go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
    valueobject "icmongolang/internal/modules/pdpa/domain/value_object"
)

// ConsentLog – Aggregate Root
type ConsentLog struct {
    ID           uuid.UUID                `json:"id"`
    UserID       uuid.UUID                `json:"user_id"`
    Purpose      valueobject.ConsentPurpose `json:"purpose"`
    Status       valueobject.ConsentStatus `json:"status"`
    PolicyID     uuid.UUID                `json:"policy_id"`
    IPAddress    string                   `json:"ip_address"`
    UserAgent    string                   `json:"user_agent"`
    GrantedAt    time.Time                `json:"granted_at"`
    RevokedAt    *time.Time               `json:"revoked_at,omitempty"`
    ExpiresAt    *time.Time               `json:"expires_at,omitempty"`
    CreatedAt    time.Time                `json:"created_at"`
    UpdatedAt    time.Time                `json:"updated_at"`
}

// Constructor – บังคับ invariants ตอนสร้าง
func NewConsentLog(
    userID uuid.UUID,
    purpose valueobject.ConsentPurpose,
    policyID uuid.UUID,
    ipAddress, userAgent string,
) (*ConsentLog, error) {
    if userID == uuid.Nil {
        return nil, domainerrors.ErrInvalidUserID
    }
    if !purpose.IsValid() {
        return nil, domainerrors.ErrInvalidPurpose
    }
    now := time.Now()
    return &ConsentLog{
        ID:        uuid.New(),
        UserID:    userID,
        Purpose:   purpose,
        Status:    valueobject.ConsentStatusActive,
        PolicyID:  policyID,
        IPAddress: ipAddress,
        UserAgent: userAgent,
        GrantedAt: now,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// Behavior: Revoke
func (c *ConsentLog) Revoke() error {
    if c.Status == valueobject.ConsentStatusRevoked {
        return domainerrors.ErrConsentAlreadyRevoked
    }
    now := time.Now()
    c.Status = valueobject.ConsentStatusRevoked
    c.RevokedAt = &now
    c.UpdatedAt = now
    return nil
}

// Behavior: IsActive
func (c *ConsentLog) IsActive() bool {
    return c.Status == valueobject.ConsentStatusActive &&
        (c.ExpiresAt == nil || c.ExpiresAt.After(time.Now()))
}
```

---

## 17. DTO และ Validation (Full Code)

### 17.1 Go DTO with validator

```go
package dto

import (
    "time"

    "github.com/google/uuid"
)

// CreateUserRequest
type CreateUserRequest struct {
    Email     string `json:"email" validate:"required,email,max=255"`
    Password  string `json:"password" validate:"required,min=8,max=100"`
    FirstName string `json:"first_name" validate:"required,min=1,max=100"`
    LastName  string `json:"last_name" validate:"required,min=1,max=100"`
    Phone     string `json:"phone" validate:"omitempty,e164"`
    Age       int    `json:"age" validate:"omitempty,gte=18,lte=120"`
}

// UpdateUserRequest
type UpdateUserRequest struct {
    FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
    LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
    Phone     *string `json:"phone,omitempty" validate:"omitempty,e164"`
}

// UserResponse
type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Phone     string    `json:"phone,omitempty"`
    Role      string    `json:"role"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
}

// LoginRequest
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// TokenResponse
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
}
```

### 17.2 Validator Setup

```go
package validator

import (
    "errors"
    "strings"

    "github.com/go-playground/validator/v10"
)

type Validator struct {
    validate *validator.Validate
}

func New() *Validator {
    v := validator.New()
    _ = v.RegisterValidation("no_space", noSpace)
    return &Validator{validate: v}
}

func (v *Validator) Validate(i interface{}) error {
    if err := v.validate.Struct(i); err != nil {
        var errs validator.ValidationErrors
        if errors.As(err, &errs) {
            messages := make([]string, 0, len(errs))
            for _, e := range errs {
                messages = append(messages, formatError(e))
            }
            return errors.New(strings.Join(messages, "; "))
        }
        return err
    }
    return nil
}

func noSpace(fl validator.FieldLevel) bool {
    return !strings.Contains(fl.Field().String(), " ")
}

func formatError(e validator.FieldError) string {
    switch e.Tag() {
    case "required":
        return e.Field() + " is required"
    case "email":
        return e.Field() + " must be a valid email"
    case "min":
        return e.Field() + " must be at least " + e.Param() + " characters"
    case "max":
        return e.Field() + " must be at most " + e.Param() + " characters"
    case "gte":
        return e.Field() + " must be >= " + e.Param()
    case "lte":
        return e.Field() + " must be <= " + e.Param()
    default:
        return e.Field() + " is invalid"
    }
}
```

### 17.3 Usage in Handler

```go
package handler

import (
    "encoding/json"
    "net/http"

    "icmongolang/internal/delivery/rest/dto"
    "icmongolang/internal/usecase"
    validator "icmongolang/internal/pkg/validator"
)

type UserHandler struct {
    userUC    usecase.UserUsecase
    validator *validator.Validator
}

func NewUserHandler(uc usecase.UserUsecase, v *validator.Validator) *UserHandler {
    return &UserHandler{userUC: uc, validator: v}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid JSON")
        return
    }
    if err := h.validator.Validate(req); err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    user, err := h.userUC.Create(r.Context(), &req)
    if err != nil {
        respondError(w, http.StatusConflict, err.Error())
        return
    }

    respondJSON(w, http.StatusCreated, user)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
    respondJSON(w, status, map[string]string{"error": msg})
}
```

### 17.4 TypeScript DTO with class-validator

```typescript
// src/dto/create-user.dto.ts
import {
    IsEmail,
    IsString,
    MinLength,
    MaxLength,
    IsOptional,
    IsInt,
    Min,
    Max,
    IsPhoneNumber,
} from "class-validator";

export class CreateUserDto {
    @IsEmail()
    @MaxLength(255)
    email: string;

    @IsString()
    @MinLength(8)
    @MaxLength(100)
    password: string;

    @IsString()
    @MinLength(1)
    @MaxLength(100)
    firstName: string;

    @IsString()
    @MinLength(1)
    @MaxLength(100)
    lastName: string;

    @IsOptional()
    @IsPhoneNumber()
    phone?: string;

    @IsOptional()
    @IsInt()
    @Min(18)
    @Max(120)
    age?: number;
}

// src/middleware/validation.middleware.ts
import { plainToInstance } from "class-transformer";
import { validate, ValidationError } from "class-validator";
import { Request, Response, NextFunction } from "express";

export function validateDto<T extends object>(dtoClass: new () => T) {
    return async (req: Request, res: Response, next: NextFunction) => {
        const dto = plainToInstance(dtoClass, req.body);
        const errors: ValidationError[] = await validate(dto);

        if (errors.length > 0) {
            const messages = errors.flatMap((err) =>
                Object.values(err.constraints || {})
            );
            return res.status(400).json({ errors: messages });
        }

        req.body = dto;
        next();
    };
}

// Usage
import express from "express";
const app = express();
app.post("/users", validateDto(CreateUserDto), (req, res) => {
    // req.body เป็น CreateUserDto ที่ validate แล้ว
    res.json({ ok: true });
});
```

---

## 18. Transaction (Full Code)

### 18.1 GORM Transaction

```go
package repository

import (
    "context"

    "gorm.io/gorm"
)

type TransactionManager struct {
    db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManager {
    return &TransactionManager{db: db}
}

// Do – executes fn in a transaction
func (tm *TransactionManager) Do(ctx context.Context, fn func(tx *gorm.DB) error) error {
    return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        return fn(tx)
    })
}

// DoWithRetry – retry on deadlock
func (tm *TransactionManager) DoWithRetry(ctx context.Context, maxRetries int, fn func(tx *gorm.DB) error) error {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        lastErr = tm.Do(ctx, fn)
        if lastErr == nil {
            return nil
        }
        // check if retryable error
    }
    return lastErr
}
```

### 18.2 Usage Example – Bank Transfer

```go
package usecase

import (
    "context"
    "errors"

    "gorm.io/gorm"

    "icmongolang/internal/models"
    "icmongolang/internal/repository"
)

type TransferUsecase struct {
    txManager *repository.TransactionManager
    db        *gorm.DB
}

func NewTransferUsecase(tm *repository.TransactionManager, db *gorm.DB) *TransferUsecase {
    return &TransferUsecase{txManager: tm, db: db}
}

type TransferInput struct {
    FromAccountID uint
    ToAccountID   uint
    Amount        float64
}

func (u *TransferUsecase) Execute(ctx context.Context, input TransferInput) error {
    if input.Amount <= 0 {
        return errors.New("amount must be positive")
    }
    if input.FromAccountID == input.ToAccountID {
        return errors.New("cannot transfer to same account")
    }

    return u.txManager.Do(ctx, func(tx *gorm.DB) error {
        var from models.Account
        if err := tx.Set("gorm:query_option", "FOR UPDATE").
            First(&from, input.FromAccountID).Error; err != nil {
            return err
        }

        if from.Balance < input.Amount {
            return errors.New("insufficient funds")
        }

        var to models.Account
        if err := tx.Set("gorm:query_option", "FOR UPDATE").
            First(&to, input.ToAccountID).Error; err != nil {
            return err
        }

        // Update balances
        if err := tx.Model(&from).
            Update("balance", gorm.Expr("balance - ?", input.Amount)).Error; err != nil {
            return err
        }
        if err := tx.Model(&to).
            Update("balance", gorm.Expr("balance + ?", input.Amount)).Error; err != nil {
            return err
        }

        // Log transaction
        log := models.TransferLog{
            FromAccountID: input.FromAccountID,
            ToAccountID:   input.ToAccountID,
            Amount:        input.Amount,
        }
        return tx.Create(&log).Error
    })
}
```

### 18.3 TypeORM Transaction

```typescript
// src/services/transfer.service.ts
import { DataSource, EntityManager } from "typeorm";

export class TransferService {
    constructor(private dataSource: DataSource) {}

    async transfer(fromId: number, toId: number, amount: number): Promise<void> {
        await this.dataSource.transaction(async (manager: EntityManager) => {
            // Lock rows
            const from = await manager
                .createQueryBuilder("account")
                .setLock("pessimistic_write")
                .where("account.id = :id", { id: fromId })
                .getOne();

            const to = await manager
                .createQueryBuilder("account")
                .setLock("pessimistic_write")
                .where("account.id = :id", { id: toId })
                .getOne();

            if (!from || !to) throw new Error("Account not found");
            if (from.balance < amount) throw new Error("Insufficient funds");

            await manager.update("accounts", fromId, {
                balance: () => `balance - ${amount}`,
            });
            await manager.update("accounts", toId, {
                balance: () => `balance + ${amount}`,
            });

            await manager.insert("transfer_logs", {
                fromAccountId: fromId,
                toAccountId: toId,
                amount,
            });
        });
    }
}
```

---

## 19. Cache (Full Code)

### 19.1 Redis Cache Implementation

```go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
)

type Cache interface {
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Get(ctx context.Context, key string, dest interface{}) error
    Delete(ctx context.Context, keys ...string) error
    Exists(ctx context.Context, key string) (bool, error)
}

type redisCache struct {
    client *redis.Client
    prefix string
}

func NewRedisCache(client *redis.Client, prefix string) Cache {
    return &redisCache{client: client, prefix: prefix}
}

func (c *redisCache) key(k string) string {
    return fmt.Sprintf("%s:%s", c.prefix, k)
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, c.key(key), data, ttl).Err()
}

func (c *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := c.client.Get(ctx, c.key(key)).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}

func (c *redisCache) Delete(ctx context.Context, keys ...string) error {
    fullKeys := make([]string, len(keys))
    for i, k := range keys {
        fullKeys[i] = c.key(k)
    }
    return c.client.Del(ctx, fullKeys...).Err()
}

func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
    n, err := c.client.Exists(ctx, c.key(key)).Result()
    if err != nil {
        return false, err
    }
    return n > 0, nil
}
```

### 19.2 Cache-Aside Pattern

```go
package service

import (
    "context"
    "errors"
    "time"

    "icmongolang/internal/cache"
    "icmongolang/internal/repository"
)

type UserService struct {
    repo  repository.UserRepository
    cache cache.Cache
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 1. Try cache
    var user User
    if err := s.cache.Get(ctx, cacheKey, &user); err == nil {
        return &user, nil
    }

    // 2. Cache miss → DB
    dbUser, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. Save to cache
    _ = s.cache.Set(ctx, cacheKey, dbUser, 10*time.Minute)

    return dbUser, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, data *User) error {
    if err := s.repo.Update(ctx, id, data); err != nil {
        return err
    }
    // Invalidate cache
    return s.cache.Delete(ctx, fmt.Sprintf("user:%d", id))
}
```

### 19.3 Write-Through Pattern

```go
func (s *UserService) CreateUserWriteThrough(ctx context.Context, user *User) error {
    // 1. Write to DB first
    if err := s.repo.Create(ctx, user); err != nil {
        return err
    }
    // 2. Write to cache
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    return s.cache.Set(ctx, cacheKey, user, 10*time.Minute)
}
```

### 19.4 Cache Invalidation หลายแบบ

```go
// Pattern 1: Time-based TTL
cache.Set(ctx, "key", value, 5*time.Minute)

// Pattern 2: Explicit delete on update
cache.Delete(ctx, "user:1", "users:list")

// Pattern 3: Version-based
version := time.Now().Unix()
key := fmt.Sprintf("user:%d:v%d", userID, version)

// Pattern 4: Pattern-based deletion (ใช้ SCAN)
func (c *redisCache) DeletePattern(ctx context.Context, pattern string) error {
    iter := c.client.Scan(ctx, 0, c.key(pattern), 100).Iterator()
    for iter.Next(ctx) {
        if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
            return err
        }
    }
    return iter.Err()
}
```

---

# ภาคที่ 3: icmongolang Project — Full Code

---

## 20. ภาพรวมโปรเจกต์

```
icmongolang/
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── websocket/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── config/
│   ├── config.go
│   ├── config.dev.yml
│   └── config.prod.yml
├── internal/
│   ├── modules/
│   │   ├── auth/
│   │   ├── user/
│   │   ├── iot/
│   │   ├── websocket/
│   │   └── ...
│   └── shared/
│       ├── db/
│       ├── logger/
│       └── httpErrors/
├── pkg/
│   ├── jwt/
│   ├── redis/
│   ├── mqtt/
│   ├── influxdb/
│   └── helpers/
├── migrations/
├── docker-compose.dev.yml
├── .air.toml
└── go.mod
```

---

## 21. Full Code: Domain Layer

### 21.1 Entity

```go
// internal/modules/iot/domain/entity/device.go
package entity

import (
    "time"

    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/iot/domain/errors"
    valueobject "icmongolang/internal/modules/iot/domain/value_object"
)

type Device struct {
    ID         uuid.UUID                  `json:"id"`
    SerialNo   string                     `json:"serial_no"`
    Name       string                     `json:"name"`
    UserID     uuid.UUID                  `json:"user_id"`
    LocationID *uuid.UUID                 `json:"location_id,omitempty"`
    Status     valueobject.DeviceStatus   `json:"status"`
    Type       valueobject.DeviceType     `json:"type"`
    LastSeen   *time.Time                 `json:"last_seen,omitempty"`
    Metadata   map[string]interface{}     `json:"metadata,omitempty"`
    CreatedAt  time.Time                  `json:"created_at"`
    UpdatedAt  time.Time                  `json:"updated_at"`
}

func NewDevice(userID uuid.UUID, serialNo, name string, deviceType valueobject.DeviceType) (*Device, error) {
    if userID == uuid.Nil {
        return nil, domainerrors.ErrInvalidUserID
    }
    if serialNo == "" {
        return nil, domainerrors.ErrInvalidSerialNo
    }
    if !deviceType.IsValid() {
        return nil, domainerrors.ErrInvalidDeviceType
    }
    now := time.Now()
    return &Device{
        ID:        uuid.New(),
        SerialNo:  serialNo,
        Name:      name,
        UserID:    userID,
        Status:    valueobject.DeviceStatusInactive,
        Type:      deviceType,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

func (d *Device) Activate() error {
    if d.Status == valueobject.DeviceStatusActive {
        return domainerrors.ErrDeviceAlreadyActive
    }
    d.Status = valueobject.DeviceStatusActive
    d.UpdatedAt = time.Now()
    return nil
}

func (d *Device) Deactivate() error {
    d.Status = valueobject.DeviceStatusInactive
    d.UpdatedAt = time.Now()
    return nil
}

func (d *Device) UpdateLastSeen() {
    now := time.Now()
    d.LastSeen = &now
    d.UpdatedAt = now
}

func (d *Device) IsActive() bool {
    return d.Status == valueobject.DeviceStatusActive
}
```

### 21.2 Value Objects

```go
// internal/modules/iot/domain/value_object/status.go
package valueobject

type DeviceStatus string

const (
    DeviceStatusActive      DeviceStatus = "ACTIVE"
    DeviceStatusInactive    DeviceStatus = "INACTIVE"
    DeviceStatusMaintenance DeviceStatus = "MAINTENANCE"
    DeviceStatusError       DeviceStatus = "ERROR"
)

func (s DeviceStatus) IsValid() bool {
    switch s {
    case DeviceStatusActive, DeviceStatusInactive, DeviceStatusMaintenance, DeviceStatusError:
        return true
    }
    return false
}

func (s DeviceStatus) String() string {
    return string(s)
}
```

```go
// internal/modules/iot/domain/value_object/device_type.go
package valueobject

type DeviceType string

const (
    DeviceTypeWaterMeter DeviceType = "WATER_METER"
    DeviceTypeGasMeter   DeviceType = "GAS_METER"
    DeviceTypeSensor     DeviceType = "SENSOR"
    DeviceTypeGPS        DeviceType = "GPS"
)

func (t DeviceType) IsValid() bool {
    switch t {
    case DeviceTypeWaterMeter, DeviceTypeGasMeter, DeviceTypeSensor, DeviceTypeGPS:
        return true
    }
    return false
}
```

### 21.3 Repository Interface

```go
// internal/modules/iot/domain/repository/device_repository.go
package repository

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iot/domain/entity"
)

type DeviceRepository interface {
    Save(ctx context.Context, device *entity.Device) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error)
    FindBySerialNo(ctx context.Context, serialNo string) (*entity.Device, error)
    FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Device, error)
    Update(ctx context.Context, device *entity.Device) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### 21.4 Domain Errors

```go
// internal/modules/iot/domain/errors/errors.go
package domainerrors

import "errors"

var (
    ErrDeviceNotFound       = errors.New("device not found")
    ErrDeviceAlreadyActive  = errors.New("device already active")
    ErrInvalidUserID        = errors.New("invalid user id")
    ErrInvalidSerialNo      = errors.New("invalid serial number")
    ErrInvalidDeviceType    = errors.New("invalid device type")
    ErrDeviceAlreadyExists  = errors.New("device already exists")
    ErrUnauthorized         = errors.New("unauthorized")
)
```

---

## 22. Full Code: Application Layer

### 22.1 Use Case

```go
// internal/modules/iot/application/create_device.go
package application

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iot/domain/entity"
    domainerrors "icmongolang/internal/modules/iot/domain/errors"
    "icmongolang/internal/modules/iot/domain/repository"
    valueobject "icmongolang/internal/modules/iot/domain/value_object"
    "icmongolang/pkg/logger"
)

type CreateDeviceUseCase struct {
    repo   repository.DeviceRepository
    logger logger.Logger
}

func NewCreateDeviceUseCase(
    repo repository.DeviceRepository,
    log logger.Logger,
) *CreateDeviceUseCase {
    return &CreateDeviceUseCase{repo: repo, logger: log}
}

type CreateDeviceInput struct {
    UserID    uuid.UUID
    SerialNo  string
    Name      string
    Type      valueobject.DeviceType
    IPAddress string
}

type CreateDeviceOutput struct {
    ID       uuid.UUID `json:"id"`
    SerialNo string    `json:"serial_no"`
    Name     string    `json:"name"`
    Status   string    `json:"status"`
}

func (uc *CreateDeviceUseCase) Execute(ctx context.Context, input CreateDeviceInput) (*CreateDeviceOutput, error) {
    // Check duplicate
    existing, _ := uc.repo.FindBySerialNo(ctx, input.SerialNo)
    if existing != nil {
        return nil, domainerrors.ErrDeviceAlreadyExists
    }

    device, err := entity.NewDevice(input.UserID, input.SerialNo, input.Name, input.Type)
    if err != nil {
        return nil, err
    }

    if err := uc.repo.Save(ctx, device); err != nil {
        uc.logger.Error("failed to save device", "error", err, "serial", input.SerialNo)
        return nil, err
    }

    uc.logger.Info("device created", "id", device.ID, "user_id", input.UserID)

    return &CreateDeviceOutput{
        ID:       device.ID,
        SerialNo: device.SerialNo,
        Name:     device.Name,
        Status:   device.Status.String(),
    }, nil
}
```

### 22.2 Get Device Use Case

```go
// internal/modules/iot/application/get_device.go
package application

import (
    "context"

    "github.com/google/uuid"
    "icmongolang/internal/modules/iot/domain/repository"
)

type GetDeviceUseCase struct {
    repo repository.DeviceRepository
}

func NewGetDeviceUseCase(repo repository.DeviceRepository) *GetDeviceUseCase {
    return &GetDeviceUseCase{repo: repo}
}

type GetDeviceInput struct {
    ID     uuid.UUID
    UserID uuid.UUID
}

type GetDeviceOutput struct {
    ID       uuid.UUID `json:"id"`
    SerialNo string    `json:"serial_no"`
    Name     string    `json:"name"`
    Status   string    `json:"status"`
    Type     string    `json:"type"`
}

func (uc *GetDeviceUseCase) Execute(ctx context.Context, input GetDeviceInput) (*GetDeviceOutput, error) {
    device, err := uc.repo.FindByID(ctx, input.ID)
    if err != nil {
        return nil, err
    }
    // Authorization check
    if device.UserID != input.UserID {
        return nil, domainerrors.ErrUnauthorized
    }
    return &GetDeviceOutput{
        ID:       device.ID,
        SerialNo: device.SerialNo,
        Name:     device.Name,
        Status:   device.Status.String(),
        Type:     device.Type.String(),
    }, nil
}
```

---

## 23. Full Code: Infrastructure Layer

### 23.1 GORM Model

```go
// internal/modules/iot/infrastructure/persistence/postgres/models.go
package postgres

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
)

type DeviceModel struct {
    ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    SerialNo   string         `gorm:"type:varchar(50);uniqueIndex;not null"`
    Name       string         `gorm:"type:varchar(200);not null"`
    UserID     uuid.UUID      `gorm:"type:uuid;index;not null"`
    LocationID *uuid.UUID     `gorm:"type:uuid;index"`
    Status     string         `gorm:"type:varchar(20);not null;index"`
    Type       string         `gorm:"type:varchar(30);not null"`
    LastSeen   *time.Time
    Metadata   datatypes.JSON `gorm:"type:jsonb"`
    CreatedAt  time.Time      `gorm:"default:now()"`
    UpdatedAt  time.Time      `gorm:"default:now()"`
}

func (DeviceModel) TableName() string {
    return "iot_devices"
}
```

### 23.2 Repository Implementation

```go
// internal/modules/iot/infrastructure/persistence/postgres/device_repo_impl.go
package postgres

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "gorm.io/gorm"

    "icmongolang/internal/modules/iot/domain/entity"
    domainerrors "icmongolang/internal/modules/iot/domain/errors"
    valueobject "icmongolang/internal/modules/iot/domain/value_object"
)

type deviceRepoImpl struct {
    db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *deviceRepoImpl {
    return &deviceRepoImpl{db: db}
}

func (r *deviceRepoImpl) Save(ctx context.Context, device *entity.Device) error {
    m := r.toModel(device)
    return r.db.WithContext(ctx).Create(m).Error
}

func (r *deviceRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
    var m DeviceModel
    err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrDeviceNotFound
    }
    if err != nil {
        return nil, err
    }
    return r.toEntity(&m), nil
}

func (r *deviceRepoImpl) FindBySerialNo(ctx context.Context, serialNo string) (*entity.Device, error) {
    var m DeviceModel
    err := r.db.WithContext(ctx).First(&m, "serial_no = ?", serialNo).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrDeviceNotFound
    }
    if err != nil {
        return nil, err
    }
    return r.toEntity(&m), nil
}

func (r *deviceRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Device, error) {
    var models []DeviceModel
    if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
        return nil, err
    }
    devices := make([]*entity.Device, len(models))
    for i := range models {
        devices[i] = r.toEntity(&models[i])
    }
    return devices, nil
}

func (r *deviceRepoImpl) Update(ctx context.Context, device *entity.Device) error {
    m := r.toModel(device)
    return r.db.WithContext(ctx).Save(m).Error
}

func (r *deviceRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&DeviceModel{}, "id = ?", id).Error
}

func (r *deviceRepoImpl) toModel(e *entity.Device) *DeviceModel {
    return &DeviceModel{
        ID:         e.ID,
        SerialNo:   e.SerialNo,
        Name:       e.Name,
        UserID:     e.UserID,
        LocationID: e.LocationID,
        Status:     string(e.Status),
        Type:       string(e.Type),
        LastSeen:   e.LastSeen,
        CreatedAt:  e.CreatedAt,
        UpdatedAt:  e.UpdatedAt,
    }
}

func (r *deviceRepoImpl) toEntity(m *DeviceModel) *entity.Device {
    return &entity.Device{
        ID:         m.ID,
        SerialNo:   m.SerialNo,
        Name:       m.Name,
        UserID:     m.UserID,
        LocationID: m.LocationID,
        Status:     valueobject.DeviceStatus(m.Status),
        Type:       valueobject.DeviceType(m.Type),
        LastSeen:   m.LastSeen,
        CreatedAt:  m.CreatedAt,
        UpdatedAt:  m.UpdatedAt,
    }
}
```

### 23.3 Redis Cache

```go
// internal/modules/iot/infrastructure/persistence/redis/device_cache.go
package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"

    "icmongolang/internal/modules/iot/domain/entity"
)

type DeviceCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewDeviceCache(client *redis.Client) *DeviceCache {
    return &DeviceCache{client: client, ttl: 10 * time.Minute}
}

func (c *DeviceCache) key(id uuid.UUID) string {
    return fmt.Sprintf("iot:device:%s", id)
}

func (c *DeviceCache) Set(ctx context.Context, device *entity.Device) error {
    data, err := json.Marshal(device)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, c.key(device.ID), data, c.ttl).Err()
}

func (c *DeviceCache) Get(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
    data, err := c.client.Get(ctx, c.key(id)).Bytes()
    if err != nil {
        return nil, err
    }
    var device entity.Device
    if err := json.Unmarshal(data, &device); err != nil {
        return nil, err
    }
    return &device, nil
}

func (c *DeviceCache) Delete(ctx context.Context, id uuid.UUID) error {
    return c.client.Del(ctx, c.key(id)).Err()
}
```

---

## 24. Full Code: Interface Layer

### 24.1 HTTP Handler

```go
// internal/modules/iot/interfaces/http/device_handler.go
package http

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"

    "icmongolang/internal/modules/iot/application"
    "icmongolang/internal/modules/iot/interfaces/http/dto"
    "icmongolang/pkg/httputil"
    "icmongolang/pkg/validator"
)

type DeviceHandler struct {
    createUC  *application.CreateDeviceUseCase
    getUC     *application.GetDeviceUseCase
    validator *validator.Validator
}

func NewDeviceHandler(
    createUC *application.CreateDeviceUseCase,
    getUC *application.GetDeviceUseCase,
    v *validator.Validator,
) *DeviceHandler {
    return &DeviceHandler{
        createUC:  createUC,
        getUC:     getUC,
        validator: v,
    }
}

func (h *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(uuid.UUID)
    if !ok {
        httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
        return
    }

    var req dto.CreateDeviceRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
        return
    }
    if err := h.validator.Validate(req); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    output, err := h.createUC.Execute(r.Context(), application.CreateDeviceInput{
        UserID:    userID,
        SerialNo:  req.SerialNo,
        Name:      req.Name,
        Type:      valueobject.DeviceType(req.Type),
        IPAddress: r.RemoteAddr,
    })
    if err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }
    httputil.JSON(w, http.StatusCreated, output)
}

func (h *DeviceHandler) Get(w http.ResponseWriter, r *http.Request) {
    userID, _ := r.Context().Value("user_id").(uuid.UUID)
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
        return
    }
    output, err := h.getUC.Execute(r.Context(), application.GetDeviceInput{
        ID:     id,
        UserID: userID,
    })
    if err != nil {
        httputil.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
        return
    }
    httputil.JSON(w, http.StatusOK, output)
}
```

### 24.2 Routes

```go
// internal/modules/iot/interfaces/http/routes.go
package http

import (
    "github.com/go-chi/chi/v5"

    "icmongolang/internal/shared/middleware"
)

type Handlers struct {
    Device *DeviceHandler
}

func RegisterRoutes(r chi.Router, h *Handlers, authMiddleware func(http.Handler) http.Handler) {
    r.Route("/api/v1/iot/devices", func(r chi.Router) {
        r.Use(authMiddleware)
        r.Post("/", h.Device.Create)
        r.Get("/{id}", h.Device.Get)
    })
}
```

### 24.3 DTO

```go
// internal/modules/iot/interfaces/http/dto/device_dto.go
package dto

type CreateDeviceRequest struct {
    SerialNo string `json:"serial_no" validate:"required,min=1,max=50"`
    Name     string `json:"name" validate:"required,min=1,max=200"`
    Type     string `json:"type" validate:"required,oneof=WATER_METER GAS_METER SENSOR GPS"`
}
```

---

## 25. Full Code: Bootstrap / Entry Points

### 25.1 main.go

```go
// cmd/api/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/joho/godotenv"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    iotHTTP "icmongolang/internal/modules/iot/interfaces/http"
    iotApp "icmongolang/internal/modules/iot/application"
    iotPG "icmongolang/internal/modules/iot/infrastructure/persistence/postgres"
    appMiddleware "icmongolang/internal/shared/middleware"
    "icmongolang/pkg/jwt"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/validator"
)

func main() {
    _ = godotenv.Load()

    // 1. Logger
    log := logger.New()

    // 2. Database
    db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
    if err != nil {
        log.Fatal("db connection failed", "error", err)
    }

    // Auto-migrate (dev only)
    if err := db.AutoMigrate(&iotPG.DeviceModel{}); err != nil {
        log.Fatal("migration failed", "error", err)
    }

    // 3. Redis
    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

    // 4. JWT
    jwtMaker, err := jwt.NewRSAMaker(
        []byte(os.Getenv("JWT_PRIVATE_KEY")),
        []byte(os.Getenv("JWT_PUBLIC_KEY")),
    )
    if err != nil {
        log.Fatal("jwt init failed", "error", err)
    }

    // 5. Validator
    v := validator.New()

    // 6. Repositories
    deviceRepo := iotPG.NewDeviceRepository(db)

    // 7. Use Cases
    createDeviceUC := iotApp.NewCreateDeviceUseCase(deviceRepo, log)
    getDeviceUC := iotApp.NewGetDeviceUseCase(deviceRepo)

    // 8. Handlers
    deviceHandler := iotHTTP.NewDeviceHandler(createDeviceUC, getDeviceUC, v)

    // 9. Router
    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(60 * time.Second))

    // Auth middleware
    authMW := appMiddleware.Auth(jwtMaker)

    // Register routes
    iotHTTP.RegisterRoutes(r, &iotHTTP.Handlers{
        Device: deviceHandler,
    }, authMW)

    // Health check
    r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("OK"))
    })

    // 10. Server with graceful shutdown
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    go func() {
        log.Info("server starting", "addr", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("server error", "error", err)
        }
    }()

    // Wait for signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Info("shutting down...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("forced shutdown", "error", err)
    }
    log.Info("server exited")
}
```

### 25.2 Auth Middleware

```go
// internal/shared/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/google/uuid"

    "icmongolang/pkg/jwt"
    "icmongolang/pkg/httputil"
)

func Auth(jwtMaker *jwt.RSAMaker) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            h := r.Header.Get("Authorization")
            if !strings.HasPrefix(h, "Bearer ") {
                httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
                return
            }
            tokenStr := strings.TrimPrefix(h, "Bearer ")
            payload, err := jwtMaker.VerifyToken(tokenStr)
            if err != nil {
                httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
                return
            }
            ctx := context.WithValue(r.Context(), "user_id", payload.UserID)
            ctx = context.WithValue(ctx, "email", payload.Email)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

## 26. Full Code: WebSocket Server

### 26.1 Hub

```go
// pkg/websocket/hub.go
package websocket

import (
    "encoding/json"
    "log"
    "sync"
)

type Message struct {
    Type    string          `json:"type"`
    UserID  string          `json:"user_id,omitempty"`
    Payload json.RawMessage `json:"payload"`
}

type Client struct {
    Hub    *Hub
    Conn   *Conn
    Send   chan []byte
    UserID string
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case c := <-h.register:
            h.mu.Lock()
            h.clients[c] = true
            h.mu.Unlock()
            log.Printf("client registered: %s", c.UserID)

        case c := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[c]; ok {
                delete(h.clients, c)
                close(c.Send)
            }
            h.mu.Unlock()
            log.Printf("client unregistered: %s", c.UserID)

        case msg := <-h.broadcast:
            h.mu.RLock()
            for c := range h.clients {
                select {
                case c.Send <- msg:
                default:
                    close(c.Send)
                    delete(h.clients, c)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) Broadcast(msg Message) {
    data, err := json.Marshal(msg)
    if err != nil {
        return
    }
    h.broadcast <- data
}

func (h *Hub) SendToUser(userID string, msg Message) {
    data, err := json.Marshal(msg)
    if err != nil {
        return
    }
    h.mu.RLock()
    defer h.mu.RUnlock()
    for c := range h.clients {
        if c.UserID == userID {
            select {
            case c.Send <- data:
            default:
            }
        }
    }
}

func (h *Hub) Register(c *Client) {
    h.register <- c
}

func (h *Hub) Unregister(c *Client) {
    h.unregister <- c
}
```

### 26.2 Handler

```go
// internal/modules/websocket/interfaces/http/ws_handler.go
package http

import (
    "net/http"

    "github.com/gorilla/websocket"

    "icmongolang/pkg/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // Configure properly in production
    },
}

type WSHandler struct {
    hub *websocket.Hub
}

func NewWSHandler(hub *websocket.Hub) *WSHandler {
    return &WSHandler{hub: hub}
}

func (h *WSHandler) Handle(w http.ResponseWriter, r *http.Request) {
    userID, _ := r.Context().Value("user_id").(string)

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    client := &websocket.Client{
        Hub:    h.hub,
        Conn:   conn,
        Send:   make(chan []byte, 256),
        UserID: userID,
    }

    h.hub.Register(client)

    go client.WritePump()
    go client.ReadPump()
}
```

### 26.3 Client Pumps

```go
// pkg/websocket/client.go
package websocket

import (
    "time"

    "github.com/gorilla/websocket"
)

type Conn = websocket.Conn

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 1024 * 1024
)

func (c *Client) ReadPump() {
    defer func() {
        c.Hub.Unregister(c)
        _ = c.Conn.Close()
    }()
    c.Conn.SetReadLimit(maxMessageSize)
    _ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
    c.Conn.SetPongHandler(func(string) error {
        return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
    })

    for {
        _, _, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err,
                websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                // log
            }
            return
        }
        // Handle incoming message if needed
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        _ = c.Conn.Close()
    }()

    for {
        select {
        case msg, ok := <-c.Send:
            _ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                _ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                return
            }

        case <-ticker.C:
            _ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

---

## 27. Full Code: Worker Pool

```go
// pkg/worker/pool.go
package worker

import (
    "context"
    "log"
    "sync"
    "time"
)

type Job func(ctx context.Context) error

type Pool struct {
    numWorkers int
    jobs       chan Job
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

func NewPool(numWorkers, queueSize int) *Pool {
    ctx, cancel := context.WithCancel(context.Background())
    return &Pool{
        numWorkers: numWorkers,
        jobs:       make(chan Job, queueSize),
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (p *Pool) Start() {
    for i := 0; i < p.numWorkers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

func (p *Pool) worker(id int) {
    defer p.wg.Done()
    for {
        select {
        case <-p.ctx.Done():
            log.Printf("[worker %d] shutting down", id)
            return
        case job, ok := <-p.jobs:
            if !ok {
                return
            }
            func() {
                defer func() {
                    if r := recover(); r != nil {
                        log.Printf("[worker %d] panic: %v", id, r)
                    }
                }()
                ctx, cancel := context.WithTimeout(p.ctx, 30*time.Second)
                defer cancel()
                if err := job(ctx); err != nil {
                    log.Printf("[worker %d] job error: %v", id, err)
                }
            }()
        }
    }
}

func (p *Pool) Submit(job Job) bool {
    select {
    case p.jobs <- job:
        return true
    case <-p.ctx.Done():
        return false
    default:
        log.Println("queue full, dropping job")
        return false
    }
}

func (p *Pool) Shutdown(timeout time.Duration) {
    close(p.jobs)
    done := make(chan struct{})
    go func() {
        p.wg.Wait()
        close(done)
    }()
    select {
    case <-done:
        log.Println("all workers stopped")
    case <-time.After(timeout):
        log.Println("shutdown timeout, forcing cancel")
        p.cancel()
    }
}
```

### 27.1 Usage Example

```go
func main() {
    pool := worker.NewPool(5, 100)
    pool.Start()

    for i := 0; i < 100; i++ {
        i := i
        pool.Submit(func(ctx context.Context) error {
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("Processed job %d\n", i)
            return nil
        })
    }

    time.Sleep(5 * time.Second)
    pool.Shutdown(10 * time.Second)
}
```

---

## 28. Full Code: MQTT Client

```go
// pkg/mqtt/client.go
package mqtt

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
    Broker   string
    ClientID string
    Username string
    Password string
    KeepAlive time.Duration
}

type Client struct {
    client mqtt.Client
    config Config
    mu     sync.RWMutex
    handlers map[string]MessageHandler
}

type MessageHandler func(topic string, payload []byte)

func NewClient(cfg Config) (*Client, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(cfg.Broker)
    opts.SetClientID(cfg.ClientID)
    opts.SetUsername(cfg.Username)
    opts.SetPassword(cfg.Password)
    opts.SetKeepAlive(cfg.KeepAlive)
    opts.SetAutoReconnect(true)
    opts.SetConnectRetry(true)
    opts.SetConnectRetryInterval(5 * time.Second)

    c := mqtt.NewClient(opts)

    token := c.Connect()
    if !token.WaitTimeout(10 * time.Second) {
        return nil, fmt.Errorf("mqtt connect timeout")
    }
    if err := token.Error(); err != nil {
        return nil, err
    }

    return &Client{
        client:   c,
        config:   cfg,
        handlers: make(map[string]MessageHandler),
    }, nil
}

func (c *Client) Publish(ctx context.Context, topic string, payload interface{}) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    token := c.client.Publish(topic, 1, false, data)
    token.Wait()
    return token.Error()
}

func (c *Client) Subscribe(topic string, handler MessageHandler) error {
    c.mu.Lock()
    c.handlers[topic] = handler
    c.mu.Unlock()

    token := c.client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
        c.mu.RLock()
        h, ok := c.handlers[msg.Topic()]
        c.mu.RUnlock()
        if ok {
            h(msg.Topic(), msg.Payload())
        }
    })
    token.Wait()
    return token.Error()
}

func (c *Client) Disconnect() {
    c.client.Disconnect(250)
}

// GetDataFromTopic – ดึงข้อมูลจาก topic
func (c *Client) GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error) {
    resultCh := make(chan []byte, 1)
    errCh := make(chan error, 1)

    token := c.client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
        select {
        case resultCh <- msg.Payload():
        default:
        }
    })
    token.Wait()
    if err := token.Error(); err != nil {
        return nil, err
    }

    defer c.client.Unsubscribe(topic)

    select {
    case data := <-resultCh:
        return data, nil
    case err := <-errCh:
        return nil, err
    case <-time.After(timeout):
        return nil, fmt.Errorf("timeout waiting for topic %s", topic)
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}
```

---

## 29. Full Code: IoT Module (Complete)

### 29.1 IoT Usecase (Complete)

```go
// internal/modules/iot/application/iot_usecase.go
package application

import (
    "context"
    "fmt"

    "github.com/google/uuid"

    "icmongolang/internal/modules/iot/domain/entity"
    domainerrors "icmongolang/internal/modules/iot/domain/errors"
    "icmongolang/internal/modules/iot/domain/repository"
    valueobject "icmongolang/internal/modules/iot/domain/value_object"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/mqtt"
    "icmongolang/pkg/websocket"
)

type IotUsecase struct {
    deviceRepo repository.DeviceRepository
    mqttClient *mqtt.Client
    wsHub      *websocket.Hub
    logger     logger.Logger
}

func NewIotUsecase(
    deviceRepo repository.DeviceRepository,
    mqttClient *mqtt.Client,
    wsHub *websocket.Hub,
    log logger.Logger,
) *IotUsecase {
    return &IotUsecase{
        deviceRepo: deviceRepo,
        mqttClient: mqttClient,
        wsHub:      wsHub,
        logger:     log,
    }
}

// RegisterDevice – ลงทะเบียนอุปกรณ์ใหม่
func (uc *IotUsecase) RegisterDevice(ctx context.Context, userID uuid.UUID, serial, name, deviceType string) (*entity.Device, error) {
    device, err := entity.NewDevice(userID, serial, name, valueobject.DeviceType(deviceType))
    if err != nil {
        return nil, err
    }

    if err := uc.deviceRepo.Save(ctx, device); err != nil {
        return nil, err
    }

    // Broadcast via WebSocket
    uc.wsHub.Broadcast(websocket.Message{
        Type: "device.registered",
        Payload: mustJSON(map[string]interface{}{
            "device_id": device.ID,
            "serial_no": device.SerialNo,
        }),
    })

    return device, nil
}

// SendCommand – ส่งคำสั่งไปยังอุปกรณ์
func (uc *IotUsecase) SendCommand(ctx context.Context, deviceID uuid.UUID, command string) error {
    device, err := uc.deviceRepo.FindByID(ctx, deviceID)
    if err != nil {
        return err
    }
    if !device.IsActive() {
        return domainerrors.ErrDeviceNotActive
    }

    topic := fmt.Sprintf("device/%s/command", device.SerialNo)
    payload := map[string]interface{}{
        "command":   command,
        "timestamp": time.Now().Unix(),
        "device_id": deviceID,
    }

    return uc.mqttClient.Publish(ctx, topic, payload)
}

// HandleGpsData – รับข้อมูล GPS จาก MQTT
func (uc *IotUsecase) HandleGpsData(serial string, payload []byte) error {
    device, err := uc.deviceRepo.FindBySerialNo(context.Background(), serial)
    if err != nil {
        return err
    }

    device.UpdateLastSeen()
    if err := uc.deviceRepo.Update(context.Background(), device); err != nil {
        return err
    }

    // Broadcast real-time location
    uc.wsHub.Broadcast(websocket.Message{
        Type: "device.gps",
        Payload: payload,
    })

    return nil
}

func mustJSON(v interface{}) []byte {
    data, _ := json.Marshal(v)
    return data
}
```

---

## 30. Full Code: Auth Module (Complete)

### 30.1 Auth Usecase

```go
// internal/modules/auth/application/auth_usecase.go
package application

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"

    "icmongolang/internal/modules/auth/domain/entity"
    "icmongolang/internal/modules/auth/domain/repository"
    "icmongolang/pkg/jwt"
)

type AuthUsecase struct {
    userRepo repository.UserRepository
    jwtMaker *jwt.RSAMaker
    refresh  repository.RefreshTokenRepository
}

func NewAuthUsecase(
    userRepo repository.UserRepository,
    jwtMaker *jwt.RSAMaker,
    refresh repository.RefreshTokenRepository,
) *AuthUsecase {
    return &AuthUsecase{
        userRepo: userRepo,
        jwtMaker: jwtMaker,
        refresh:  refresh,
    }
}

type RegisterInput struct {
    Email     string
    Password  string
    FirstName string
    LastName  string
}

func (uc *AuthUsecase) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
    existing, _ := uc.userRepo.FindByEmail(ctx, input.Email)
    if existing != nil {
        return nil, errors.New("email already exists")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &entity.User{
        ID:           uuid.New(),
        Email:        input.Email,
        PasswordHash: string(hash),
        FirstName:    input.FirstName,
        LastName:     input.LastName,
        IsActive:     true,
        IsVerified:   false,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if err := uc.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

type LoginInput struct {
    Email    string
    Password string
    IP       string
    UserAgent string
}

type LoginOutput struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int
}

func (uc *AuthUsecase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
    user, err := uc.userRepo.FindByEmail(ctx, input.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
        return nil, errors.New("invalid credentials")
    }

    if !user.IsActive {
        return nil, errors.New("account disabled")
    }

    accessToken, err := uc.jwtMaker.CreateToken(user.ID, user.Email, 15*time.Minute)
    if err != nil {
        return nil, err
    }

    refreshToken := uuid.New().String()
    if err := uc.refresh.Store(ctx, refreshToken, user.ID, 7*24*time.Hour); err != nil {
        return nil, err
    }

    return &LoginOutput{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    900,
    }, nil
}

func (uc *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (*LoginOutput, error) {
    userID, err := uc.refresh.Validate(ctx, refreshToken)
    if err != nil {
        return nil, errors.New("invalid refresh token")
    }

    user, err := uc.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    accessToken, err := uc.jwtMaker.CreateToken(user.ID, user.Email, 15*time.Minute)
    if err != nil {
        return nil, err
    }

    // Rotate refresh token
    newRefresh := uuid.New().String()
    if err := uc.refresh.Store(ctx, newRefresh, user.ID, 7*24*time.Hour); err != nil {
        return nil, err
    }
    _ = uc.refresh.Delete(ctx, refreshToken)

    return &LoginOutput{
        AccessToken:  accessToken,
        RefreshToken: newRefresh,
        ExpiresIn:    900,
    }, nil
}
```

### 30.2 JWT Maker (Complete)

```go
// pkg/jwt/rsa_maker.go
package jwt

import (
    "crypto/rsa"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type Payload struct {
    UserID uuid.UUID
    Email  string
}

type RSAMaker struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
}

func NewRSAMaker(privatePEM, publicPEM []byte) (*RSAMaker, error) {
    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
    if err != nil {
        return nil, err
    }
    publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
    if err != nil {
        return nil, err
    }
    return &RSAMaker{privateKey: privateKey, publicKey: publicKey}, nil
}

func (m *RSAMaker) CreateToken(userID uuid.UUID, email string, duration time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "sub":   userID.String(),
        "email": email,
        "exp":   time.Now().Add(duration).Unix(),
        "iat":   time.Now().Unix(),
        "iss":   "icmongolang",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(m.privateKey)
}

func (m *RSAMaker) VerifyToken(tokenString string) (*Payload, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return m.publicKey, nil
    })
    if err != nil {
        return nil, err
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    userID, err := uuid.Parse(claims["sub"].(string))
    if err != nil {
        return nil, err
    }
    return &Payload{
        UserID: userID,
        Email:  claims["email"].(string),
    }, nil
}
```

---

# ภาคที่ 4: DevOps และ EDA

---

## 31. Event-Driven Design

### 31.1 Event Bus Implementation

```go
// pkg/eventbus/bus.go
package eventbus

import (
    "context"
    "log"
    "sync"
)

type Event struct {
    Type    string
    Payload interface{}
}

type Handler func(ctx context.Context, event Event) error

type Bus struct {
    mu       sync.RWMutex
    handlers map[string][]Handler
}

func NewBus() *Bus {
    return &Bus{
        handlers: make(map[string][]Handler),
    }
}

func (b *Bus) Subscribe(eventType string, handler Handler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *Bus) Publish(ctx context.Context, event Event) {
    b.mu.RLock()
    handlers := b.handlers[event.Type]
    b.mu.RUnlock()

    for _, h := range handlers {
        go func(h Handler) {
            defer func() {
                if r := recover(); r != nil {
                    log.Printf("handler panic: %v", r)
                }
            }()
            if err := h(ctx, event); err != nil {
                log.Printf("handler error for %s: %v", event.Type, err)
            }
        }(h)
    }
}
```

---

## 32. Code Review Workflow

### 32.1 GitHub Actions CI

```yaml
# .github/workflows/ci.yml
name: CI

on:
  pull_request:
    branches: [main, develop]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          args: --timeout=5m

      - name: Test
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Coverage check
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $coverage%"
          if (( $(echo "$coverage < 70" | bc -l) )); then
            echo "Coverage below 70%"
            exit 1
          fi

      - name: Build
        run: go build ./...

      - name: Security scan
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          severity: 'HIGH,CRITICAL'
          exit-code: '1'
```

---

## 33-34. Monitoring Tools & POC/ROI

(ดูรายละเอียดในเวอร์ชันก่อนหน้า — บทที่ 23-25)

---

# ภาคที่ 5: RCA และ Checklists

---

## 35. RCA Diagrams

### 35.1 Fishbone Diagram

```mermaid
flowchart LR
    subgraph Man["👤 คน"]
        M1[ทักษะไม่พอ] --> M2[ขาดการฝึก]
    end
    subgraph Machine["⚙️ เครื่อง"]
        Mac1[สึกหรอ]
    end
    subgraph Method["📋 วิธี"]
        Met1[ขั้นตอนซับซ้อน]
    end
    subgraph Material["📦 วัตถุดิบ"]
        Mat1[คุณภาพต่ำ]
    end
    subgraph Measurement["📊 การวัด"]
        Meas1[เครื่องมือไม่แม่น]
    end
    subgraph Environment["🌍 สิ่งแวดล้อม"]
        Env1[อุณหภูมิสูง]
    end

    Man --> Problem((ปัญหา))
    Machine --> Problem
    Method --> Problem
    Material --> Problem
    Measurement --> Problem
    Environment --> Problem
```

### 35.2 5 Whys

```mermaid
flowchart TD
    P[("ปัญหา: สินค้ามีตำหนิ")]
    P --> W1[ทำไม? เครื่องไม่ตรง]
    W1 --> W2[ทำไม? ตั้งเร็วเกินไป]
    W2 --> W3[ทำไม? ไม่มีคู่มือ]
    W3 --> W4[ทำไม? ไม่ปรับปรุงเอกสาร]
    W4 --> RC[("Root Cause: ขาดระบบจัดการความรู้")]
```

---

## 36. Checklists

### 36.1 Domain Layer

- [ ] Entity มี constructor (`New{{Entity}}`)
- [ ] ไม่มี setter ตรง – เปลี่ยนผ่าน behavior method
- [ ] Value Object มี `IsValid()`
- [ ] Repository เป็น interface
- [ ] Domain error เป็น sentinel
- [ ] Domain **ไม่ import** gorm, gin, sarama, redis

### 36.2 Application Layer

- [ ] Use case ละ 1 ไฟล์
- [ ] Input/Output DTO แยกชัดเจน
- [ ] Execute() return error
- [ ] ไม่มี SQL/HTTP ใน use case

### 36.3 Infrastructure Layer

- [ ] GORM model มี `TableName()` + prefix
- [ ] Repository impl map model ↔ entity
- [ ] Kafka message ใช้ JSON
- [ ] Redis key มี namespace

### 36.4 Interface Layer

- [ ] Handler ดึง `user_id` จาก context
- [ ] Route ลงทะเบียนครบ
- [ ] Error response เป็น JSON
- [ ] Auth middleware ถูก apply

### 36.5 Build & Test

- [ ] `go build ./...` ผ่าน
- [ ] `go vet ./...` ผ่าน
- [ ] `go test ./...` ผ่าน
- [ ] Coverage > 70%

---

# ภาคผนวก

## A. Dependency Matrix

| Layer | stdlib | uuid | gorm | chi | sarama | redis | jwt |
|-------|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| Domain | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Application | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Infrastructure | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ |
| Interface | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ |

## B. Kafka Topic Convention

```
<module>.<entity>.<action>
iot.device.registered
iot.device.gps
iot.order.created
alarm.triggered
payment.completed
```

## C. Error Code Convention

| Domain Error | HTTP Status |
|--------------|:-:|
| ErrNotFound | 404 |
| ErrAlreadyExists | 409 |
| ErrInvalid* | 400 |
| ErrUnauthorized | 401 |
| ErrForbidden | 403 |
| ErrRateLimit | 429 |

## D. Environment Variables

```env
# Database
DB_DSN=postgres://appuser:secret@localhost:5432/icmongo?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# JWT
JWT_PRIVATE_KEY_PATH=private.pem
JWT_PUBLIC_KEY_PATH=public.pem

# MQTT
MQTT_BROKER=tcp://localhost:1883

# SMTP
SMTP_HOST=localhost
SMTP_PORT=1025

# Kafka
KAFKA_BROKERS=localhost:9092
```

## E. สรุปสูตร

### ROI

```
ROI = (Total Benefit - Total Cost) / Total Cost × 100%
Payback Period = Total Cost / (Total Benefit / 12)
```

## F. CLI Commands

```bash
# Full workflow
go clean -cache
go clean -modcache
go mod tidy
go mod download
go mod verify
go run cmd/api/main.go migrate
swag init -g cmd/api/main.go
mockery --all
go mod vendor
go test ./...
go build ./...
air
```

---

## 📝 ข้อมูลผู้เขียน

**ผู้เขียน:** คงนคร จันทะคุณ
**อีเมล:** kongnakornjantakun@gmail.com
**Line:** kongnakornna
**วันที่:** เมษายน 2026

📞 **โทรศัพท์ / พร้อมเพย์:** 0955088091

---

**END OF DOCUMENT — Full Code Version 2.0**