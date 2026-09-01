package testdata

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
)

type ValueOnly struct {
	Name    string
	Age     int
	Enabled bool
	Hash    [32]byte
}

type WithPrivateValue struct {
	Name string
	hash [32]byte
}

func NewWithPrivateValue() WithPrivateValue {
	v := WithPrivateValue{Name: "private"}
	v.hash[0] = 7
	v.hash[31] = 9
	return v
}

func (v WithPrivateValue) Hash() [32]byte {
	return v.hash
}

type Child struct {
	Name string
}

type ExportedRefs struct {
	Child *Child
	Names []string
	Meta  map[string]int
}

type Cat struct {
	Name    string
	Sibling *Cat
}

type SharedChild struct {
	Left  *Child
	Right *Child
}

type WithPrivateRef struct {
	Name  string
	child *Child
}

func NewWithPrivateRef() WithPrivateRef {
	return WithPrivateRef{Name: "private-ref", child: &Child{Name: "hidden-child"}}
}

func (v WithPrivateRef) Child() *Child {
	return v.child
}

type RuntimeState struct {
	Mutex     sync.Mutex
	RWMutex   sync.RWMutex
	Once      sync.Once
	WaitGroup sync.WaitGroup
	Cond      *sync.Cond
	Map       sync.Map
	Atomic    atomic.Value
	Chan      chan string
	Func      func() string
	Ctx       context.Context
	File      *os.File
	Conn      net.Conn
	DB        *sql.DB
	Tx        *sql.Tx
	Client    *http.Client
	Name      string
}

type InterfaceFields struct {
	Pointer any
	Slice   any
	Map     any
}

type NestedCollections struct {
	Groups map[string][]*Child
	Items  []map[string]*Child
}

type NilRefs struct {
	Ptr   *Child
	Slice []string
	Map   map[string]int
	Any   any
}

type LargeValue struct {
	Name  string
	Count int
	A     [64]byte
	B     [64]byte
	C     [64]byte
	D     [64]byte
	E     [64]byte
	F     [64]byte
	G     [64]byte
	H     [64]byte
}

type ManyPointers struct {
	One   *Child
	Two   *Child
	Three *Child
	Four  *Child
	Five  *Child
	Six   *Child
	Seven *Child
	Eight *Child
}

type Tree struct {
	Name     string
	Children []*Tree
}

type DomainObject struct {
	Customer ExportedRefs
	Shared   SharedChild
	Nested   NestedCollections
	Private  WithPrivateRef
	Runtime  RuntimeState
}
