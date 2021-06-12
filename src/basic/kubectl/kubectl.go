package kubectl

import "fmt"

// https://github.com/kubernauts/practical-kubernetes-problems/blob/master/images/k8s-resources-map.png
// https://github.com/kubernetes/kubernetes/blob/cea1d4e20b4a7886d8ff65f34c6d4f95efcb4742/staging/src/k8s.io/cli-runtime/pkg/resource/visitor.go

// Kubernetes 抽象了多种的 Resource：Pod、ReplicaSet、ConfigMap、Volumes、Namespace、Roles 构成了 Kubernetes 的数据模型。
// kubectl 是 Kubernetes 的一个客户端命令，操作人员用这个命令来操作 Kubernetes。
// kubectl 联系到 Kubernetes 的 API Server，API Server 联系每个节点上的 kubelet，从而控制每个节点。
// kubectl 处理用户提交的内容（命令行参数、YAML 文件等），将其组织成一个数据结构体，发送给 API Server。

type VisitorFunc func(*Info, error) error

// Visitor 的接口，其中需要 Visit(VisitorFunc) error 的方法。

type Visitor interface {
	Visit(VisitorFunc) error
}

// kubectl 主要是用来处理 Info 结构体

type Info struct {
	Namespace   string
	Name        string
	OtherThings string
}

// 实现 Visitor 接口中的 Visit() 方法，直接调用传进来的方法。

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

// 多种不同的 Visitor。

// NameVisitor 的结构体有一个 Visitor 接口成员，这里意味着多态。
// 在实现 Visit() 方法时调用了自己结构体内的 Visitor 的 Visitor() 方法（装饰器模式），用另一个 Visitor 装饰了自己。

type NameVisitor struct {
	visitor Visitor
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	})
}

type OtherThingsVisitor struct {
	visitor Visitor
}

func (v OtherThingsVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	})
}

type LogVisitor struct {
	visitor Visitor
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor() before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor() after call function")
		return err
	})
}

func Example31() {
	info := Info{}
	var v Visitor = &info
	v = LogVisitor{v}
	v = NameVisitor{v}
	v = OtherThingsVisitor{v}

	loadFile := func(info *Info, err error) error {
		info.Name = "Hao Chen"
		info.Namespace = "MegaEase"
		info.OtherThings = "We are running as remote team."
		return nil
	}
	v.Visit(loadFile)
}
