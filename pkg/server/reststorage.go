package server

import (
	"context"

	kobjv1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
)

type KobjRESTStorage struct {
	rest.TableConvertor
	MemObjects map[string]*kobjv1.Kobj
}

var typeAssertKobjStorage *KobjRESTStorage
var _ rest.StandardStorage = typeAssertKobjStorage
var _ rest.Scoper = typeAssertKobjStorage
var _ rest.KindProvider = typeAssertKobjStorage
var _ rest.TableConvertor = typeAssertKobjStorage

// var _ rest.Exporter = typeAssertKobjStorage
// var _ rest.Connecter = typeAssertKobjStorage
// var _ rest.ResourceStreamer = typeAssertKobjStorage
// var _ rest.Redirector = typeAssertKobjStorage

// NewKobjStorage returns a struct that implements methods needed for Kubernetes to satisfy API requests for the `Kobj` resource
func NewKobjStorage() *KobjRESTStorage {
	return &KobjRESTStorage{
		TableConvertor: printerstorage.TableConvertor{
			TableGenerator: printers.NewTableGenerator().
				With(printersinternal.AddHandlers).
				With(KobjTableHandlers),
		},
	}
}

// Kind satisfies the KindProvider interface
func (*KobjRESTStorage) Kind() string {
	return "Kobj"
}

// NamespaceScoped satisfies the Scoper interface
func (*KobjRESTStorage) NamespaceScoped() bool {
	return true
}

// New satisfies the Storage interface
func (*KobjRESTStorage) New() runtime.Object {
	return &kobjv1.Kobj{}
}

// NewList satisfies part of the Lister interface
func (s *KobjRESTStorage) NewList() runtime.Object {
	return &kobjv1.KobjList{}
}

// List satisfies part of the Lister interface
func (s *KobjRESTStorage) List(
	ctx context.Context,
	options *metainternalversion.ListOptions,
) (runtime.Object, error) {
	list := &kobjv1.KobjList{
		Items: make([]kobjv1.Kobj, len(s.MemObjects)),
	}
	i := 0
	for _, ko := range s.MemObjects {
		ko.DeepCopyInto(&list.Items[i])
		i++
	}
	return list, nil
}

// Get satisfies the Getter interface
func (s *KobjRESTStorage) Get(
	ctx context.Context,
	name string,
	opts *metav1.GetOptions,
) (runtime.Object, error) {
	ko := s.MemObjects[name]
	if ko == nil {
		return nil, k8serrors.NewNotFound(KobjGroupResource, name)
	}
	return ko.DeepCopy(), nil
}

func (s *KobjRESTStorage) Watch(
	ctx context.Context,
	options *metainternalversion.ListOptions,
) (watch.Interface, error) {
	return nil, nil
}

func (s *KobjRESTStorage) Create(
	ctx context.Context,
	obj runtime.Object,
	createValidation rest.ValidateObjectFunc,
	options *metav1.CreateOptions,
) (runtime.Object, error) {
	return nil, nil
}

func (s *KobjRESTStorage) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *metav1.UpdateOptions,
) (runtime.Object, bool, error) {
	return nil, false, nil
}

func (s *KobjRESTStorage) Delete(
	ctx context.Context,
	name string,
	deleteValidation rest.ValidateObjectFunc,
	options *metav1.DeleteOptions,
) (runtime.Object, bool, error) {
	return nil, false, nil
}

func (s *KobjRESTStorage) DeleteCollection(
	ctx context.Context,
	deleteValidation rest.ValidateObjectFunc,
	options *metav1.DeleteOptions,
	listOptions *metainternalversion.ListOptions,
) (runtime.Object, error) {
	return nil, nil
}
