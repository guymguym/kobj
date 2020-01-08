package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	Register(&Kobj{}, &KobjList{})
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kobj TODO
type Kobj struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Key               string `json:"key"`
	Value             string `json:"value"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobjList TODO
type KobjList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kobj `json:"items"`
}

/*
func (obj *Kobj) DeepCopyObject() runtime.Object {
	out := *obj
	out.TypeMeta = obj.TypeMeta
	// obj.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return &out
}

func (obj *KobjList) DeepCopyObject() runtime.Object {
	out := *obj
	out.TypeMeta = obj.TypeMeta
	obj.ListMeta.DeepCopyInto(&out.ListMeta)
	if obj.Items != nil {
		out.Items = make([]Kobj, len(obj.Items))
		for i := range obj.Items {
			out.Items[i] = *obj.Items[i].DeepCopyObject().(*Kobj)
		}
	}
	return &out
}
*/
