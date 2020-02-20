package server

import (
	"context"
	"strings"

	kobjv1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/storage"
)

type KobjMemStorage struct {
	Index map[string]*kobjv1.Kobj
}

func NewKobjMemStorage() *KobjMemStorage {
	return &KobjMemStorage{
		Index: make(map[string]*kobjv1.Kobj),
	}
}

var _ storage.Interface = &KobjMemStorage{}

func (s *KobjMemStorage) Create(
	ctx context.Context,
	key string,
	obj runtime.Object,
	out runtime.Object,
	ttl uint64,
) error {
	ko := obj.(*kobjv1.Kobj)
	old := s.Index[key]
	s.Index[key] = ko
	if old != nil && out != nil {
		old.DeepCopyInto(out.(*kobjv1.Kobj))
	}
	return nil
}

func (s *KobjMemStorage) Delete(
	ctx context.Context,
	key string,
	out runtime.Object,
	preconditions *storage.Preconditions,
	validateDeletion storage.ValidateObjectFunc,
) error {
	old := s.Index[key]
	delete(s.Index, key)
	if old != nil && out != nil {
		old.DeepCopyInto(out.(*kobjv1.Kobj))
	}
	return nil
}

func (s *KobjMemStorage) GuaranteedUpdate(
	ctx context.Context,
	key string,
	ptrToType runtime.Object,
	ignoreNotFound bool,
	precondtions *storage.Preconditions,
	tryUpdate storage.UpdateFunc,
	suggestion ...runtime.Object,
) error {
	panic("implement me")
}

func (s *KobjMemStorage) Get(ctx context.Context, key string, resourceVersion string, objPtr runtime.Object, ignoreNotFound bool) error {
	ko := s.Index[key]
	if ko == nil {
		if ignoreNotFound {
			return nil
		}
		return errors.NewNotFound(KobjGroupResource, key)
	}
	ko.DeepCopyInto(objPtr.(*kobjv1.Kobj))
	return nil
}

func (s *KobjMemStorage) GetToList(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate, listObj runtime.Object) error {
	list := listObj.(*kobjv1.KobjList)
	ko := s.Index[key]
	items := make([]kobjv1.Kobj, 1)
	ko.DeepCopyInto(&items[0])
	list.Items = items
	return nil
}

func (s *KobjMemStorage) List(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate, listObj runtime.Object) error {
	list := listObj.(*kobjv1.KobjList)
	items := make([]kobjv1.Kobj, len(s.Index))
	i := 0
	for _, ko := range s.Index {
		ko.DeepCopyInto(&items[i])
		i++
	}
	list.Items = items
	return nil
}

func (s *KobjMemStorage) Watch(
	ctx context.Context,
	key string,
	resourceVersion string,
	p storage.SelectionPredicate,
) (watch.Interface, error) {
	panic("implement me")
}

func (s *KobjMemStorage) WatchList(
	ctx context.Context,
	key string,
	resourceVersion string,
	p storage.SelectionPredicate,
) (watch.Interface, error) {
	panic("implement me")
}

func (s *KobjMemStorage) Count(key string) (int64, error) {
	var count int64
	for k := range s.Index {
		if strings.HasPrefix(k, key) {
			count++
		}
	}
	return count, nil
}

func (s *KobjMemStorage) Versioner() storage.Versioner {
	panic("implement me")
}
