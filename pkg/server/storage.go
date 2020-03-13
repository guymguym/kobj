package server

import (
	"context"
	"strings"

	kobjv1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/klog"
)

type KobjMemStorage struct {
	Index map[string]*kobjv1.Kobj
}

func NewKobjMemRESTStorage() *rest.Storage {
	storage := NewKobjMemStorage()
	strategy := &KobjStrategy{}
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &kobjv1.Kobj{} },
		NewListFunc:              func() runtime.Object { return &kobjv1.KobjList{} },
		DefaultQualifiedResource: KobjGroupResource,
		PredicateFunc:            KobjMatcher,
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
		TableConvertor: printerstorage.TableConvertor{
			TableGenerator: printers.NewTableGenerator().
				With(printersinternal.AddHandlers).
				With(KobjTableHandlers),
		},
		Storage: genericregistry.DryRunnableStorage{
			Storage: storage,
			Codec:   codecs.LegacyCodec(kobjv1.SchemeGroupVersion),
			// Codec:   codecs.CodecForVersions(kobjv1.SchemeGroupVersion),
		},
	}
	util.Assert(store.CompleteWithOptions(&generic.StoreOptions{
		RESTOptions: generic.RESTOptions{ResourcePrefix: KobjResource},
		AttrFunc:    KobjGetAttrs,
	}))
	return store
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
	klog.Infof("KobjMemStorage.Create: %s => %v", key, obj)
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
	klog.Infof("KobjMemStorage.Delete: %s", key)
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
	klog.Infof("KobjMemStorage.GuaranteedUpdate: %s => %v", key, ptrToType)
	panic("KobjMemStorage.GuaranteedUpdate: TODO implement me")
}

func (s *KobjMemStorage) Get(
	ctx context.Context,
	key string,
	resourceVersion string,
	objPtr runtime.Object,
	ignoreNotFound bool,
) error {
	klog.Infof("KobjMemStorage.Get: %s => %v", key, objPtr)
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

func (s *KobjMemStorage) GetToList(
	ctx context.Context,
	key string,
	resourceVersion string,
	p storage.SelectionPredicate,
	listObj runtime.Object,
) error {
	klog.Infof("KobjMemStorage.GetToList: %s => %v", key, listObj)
	list := listObj.(*kobjv1.KobjList)
	ko := s.Index[key]
	items := make([]kobjv1.Kobj, 1)
	ko.DeepCopyInto(&items[0])
	list.Items = items
	return nil
}

func (s *KobjMemStorage) List(
	ctx context.Context,
	key string,
	resourceVersion string,
	p storage.SelectionPredicate,
	listObj runtime.Object,
) error {
	klog.Infof("KobjMemStorage.List: %s => %v", key, listObj)
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
	klog.Infof("KobjMemStorage.Watch: %s", key)
	panic("KobjMemStorage.Watch: TODO implement me")
}

func (s *KobjMemStorage) WatchList(
	ctx context.Context,
	key string,
	resourceVersion string,
	p storage.SelectionPredicate,
) (watch.Interface, error) {
	klog.Infof("KobjMemStorage.WatchList: %s", key)
	panic("KobjMemStorage.WatchList: TODO implement me")
}

func (s *KobjMemStorage) Count(key string) (int64, error) {
	klog.Infof("KobjMemStorage.Count: %s", key)
	var count int64
	for k := range s.Index {
		if strings.HasPrefix(k, key) {
			count++
		}
	}
	return count, nil
}

func (s *KobjMemStorage) Versioner() storage.Versioner {
	klog.Infof("KobjMemStorage.Versioner: TODO")
	panic("KobjMemStorage.Versioner: TODO implement me")
}
