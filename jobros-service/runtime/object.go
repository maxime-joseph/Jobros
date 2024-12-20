package runtime

// GroupVersionKind uniquely identifies a resource type.
type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

// Object is a common interface that all resources should implement.
type Object interface {
	GetGroupVersionKind() GroupVersionKind
	DeepCopy() Object
}
