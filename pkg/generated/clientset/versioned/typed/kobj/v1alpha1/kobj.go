/*
   Copyright 2020 Guy Margalit.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"
	scheme "github.com/kobj-io/kobj/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// KobjsGetter has a method to return a KobjInterface.
// A group's client should implement this interface.
type KobjsGetter interface {
	Kobjs(namespace string) KobjInterface
}

// KobjInterface has methods to work with Kobj resources.
type KobjInterface interface {
	Create(*v1alpha1.Kobj) (*v1alpha1.Kobj, error)
	Update(*v1alpha1.Kobj) (*v1alpha1.Kobj, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Kobj, error)
	List(opts v1.ListOptions) (*v1alpha1.KobjList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Kobj, err error)
	KobjExpansion
}

// kobjs implements KobjInterface
type kobjs struct {
	client rest.Interface
	ns     string
}

// newKobjs returns a Kobjs
func newKobjs(c *KobjV1alpha1Client, namespace string) *kobjs {
	return &kobjs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the kobj, and returns the corresponding kobj object, and an error if there is any.
func (c *kobjs) Get(name string, options v1.GetOptions) (result *v1alpha1.Kobj, err error) {
	result = &v1alpha1.Kobj{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kobjs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Kobjs that match those selectors.
func (c *kobjs) List(opts v1.ListOptions) (result *v1alpha1.KobjList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.KobjList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kobjs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kobjs.
func (c *kobjs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("kobjs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a kobj and creates it.  Returns the server's representation of the kobj, and an error, if there is any.
func (c *kobjs) Create(kobj *v1alpha1.Kobj) (result *v1alpha1.Kobj, err error) {
	result = &v1alpha1.Kobj{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("kobjs").
		Body(kobj).
		Do().
		Into(result)
	return
}

// Update takes the representation of a kobj and updates it. Returns the server's representation of the kobj, and an error, if there is any.
func (c *kobjs) Update(kobj *v1alpha1.Kobj) (result *v1alpha1.Kobj, err error) {
	result = &v1alpha1.Kobj{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kobjs").
		Name(kobj.Name).
		Body(kobj).
		Do().
		Into(result)
	return
}

// Delete takes name of the kobj and deletes it. Returns an error if one occurs.
func (c *kobjs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kobjs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kobjs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kobjs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched kobj.
func (c *kobjs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Kobj, err error) {
	result = &v1alpha1.Kobj{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("kobjs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
