package server

import (
	"context"
	"fmt"

	kobjv1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

type KobjStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (*KobjStrategy) ObjectKinds(obj runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	return []schema.GroupVersionKind{kobjv1.SchemeGroupVersion.WithKind("Kobj")}, false, nil
}

// NamespaceScoped returns true because all objects need to be within a namespace.
func (*KobjStrategy) NamespaceScoped() bool {
	return true
}

// Canonicalize normalizes the object after validation.
func (*KobjStrategy) Canonicalize(obj runtime.Object) {
}

func (*KobjStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// AllowCreateOnUpdate is false for kobj; this means a POST is needed to create one.
func (*KobjStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForCreate clears the status of an object before creation.
func (*KobjStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	// ko := obj.(*kobjv1.Kobj)
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (*KobjStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	// ko := obj.(*kobjv1.Kobj)
	// oldko := old.(*kobjv1.Kobj)
}

// Validate validates a new object.
func (*KobjStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	// ko := obj.(*kobjv1.Kobj)
	return nil
}

// ValidateUpdate is the default update validation for an end user.
func (*KobjStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	// ko := obj.(*kobjv1.Kobj)
	// oldko := old.(*kobjv1.Kobj)
	return nil
}

// KobjGetAttrs returns labels and fields of a given object for filtering purposes.
func KobjGetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	ko, ok := obj.(*kobjv1.Kobj)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a kobj.")
	}
	return labels.Set(ko.ObjectMeta.Labels), KobjToSelectableFields(ko), nil
}

// KobjToSelectableFields returns a field set that represents the object for matching purposes.
func KobjToSelectableFields(ko *kobjv1.Kobj) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&ko.ObjectMeta, true)
	return objectMetaFieldsSet
}

// KobjMatcher is the filter used by the generic storage backend to watch events
// from storage to clients of the apiserver only interested in specific labels/fields.
func KobjMatcher(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: KobjGetAttrs,
	}
}
