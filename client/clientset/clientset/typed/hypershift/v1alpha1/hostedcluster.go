/*


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
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha1 "github.com/openshift/hypershift/api/hypershift/v1alpha1"
	hypershiftv1alpha1 "github.com/openshift/hypershift/client/applyconfiguration/hypershift/v1alpha1"
	scheme "github.com/openshift/hypershift/client/clientset/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HostedClustersGetter has a method to return a HostedClusterInterface.
// A group's client should implement this interface.
type HostedClustersGetter interface {
	HostedClusters(namespace string) HostedClusterInterface
}

// HostedClusterInterface has methods to work with HostedCluster resources.
type HostedClusterInterface interface {
	Create(ctx context.Context, hostedCluster *v1alpha1.HostedCluster, opts v1.CreateOptions) (*v1alpha1.HostedCluster, error)
	Update(ctx context.Context, hostedCluster *v1alpha1.HostedCluster, opts v1.UpdateOptions) (*v1alpha1.HostedCluster, error)
	UpdateStatus(ctx context.Context, hostedCluster *v1alpha1.HostedCluster, opts v1.UpdateOptions) (*v1alpha1.HostedCluster, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.HostedCluster, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.HostedClusterList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.HostedCluster, err error)
	Apply(ctx context.Context, hostedCluster *hypershiftv1alpha1.HostedClusterApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.HostedCluster, err error)
	ApplyStatus(ctx context.Context, hostedCluster *hypershiftv1alpha1.HostedClusterApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.HostedCluster, err error)
	HostedClusterExpansion
}

// hostedClusters implements HostedClusterInterface
type hostedClusters struct {
	client rest.Interface
	ns     string
}

// newHostedClusters returns a HostedClusters
func newHostedClusters(c *HypershiftV1alpha1Client, namespace string) *hostedClusters {
	return &hostedClusters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the hostedCluster, and returns the corresponding hostedCluster object, and an error if there is any.
func (c *hostedClusters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.HostedCluster, err error) {
	result = &v1alpha1.HostedCluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hostedclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HostedClusters that match those selectors.
func (c *hostedClusters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.HostedClusterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.HostedClusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hostedclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested hostedClusters.
func (c *hostedClusters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("hostedclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a hostedCluster and creates it.  Returns the server's representation of the hostedCluster, and an error, if there is any.
func (c *hostedClusters) Create(ctx context.Context, hostedCluster *v1alpha1.HostedCluster, opts v1.CreateOptions) (result *v1alpha1.HostedCluster, err error) {
	result = &v1alpha1.HostedCluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("hostedclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hostedCluster).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a hostedCluster and updates it. Returns the server's representation of the hostedCluster, and an error, if there is any.
func (c *hostedClusters) Update(ctx context.Context, hostedCluster *v1alpha1.HostedCluster, opts v1.UpdateOptions) (result *v1alpha1.HostedCluster, err error) {
	result = &v1alpha1.HostedCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hostedclusters").
		Name(hostedCluster.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hostedCluster).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *hostedClusters) UpdateStatus(ctx context.Context, hostedCluster *v1alpha1.HostedCluster, opts v1.UpdateOptions) (result *v1alpha1.HostedCluster, err error) {
	result = &v1alpha1.HostedCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hostedclusters").
		Name(hostedCluster.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hostedCluster).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the hostedCluster and deletes it. Returns an error if one occurs.
func (c *hostedClusters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hostedclusters").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *hostedClusters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hostedclusters").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched hostedCluster.
func (c *hostedClusters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.HostedCluster, err error) {
	result = &v1alpha1.HostedCluster{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("hostedclusters").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied hostedCluster.
func (c *hostedClusters) Apply(ctx context.Context, hostedCluster *hypershiftv1alpha1.HostedClusterApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.HostedCluster, err error) {
	if hostedCluster == nil {
		return nil, fmt.Errorf("hostedCluster provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(hostedCluster)
	if err != nil {
		return nil, err
	}
	name := hostedCluster.Name
	if name == nil {
		return nil, fmt.Errorf("hostedCluster.Name must be provided to Apply")
	}
	result = &v1alpha1.HostedCluster{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("hostedclusters").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *hostedClusters) ApplyStatus(ctx context.Context, hostedCluster *hypershiftv1alpha1.HostedClusterApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.HostedCluster, err error) {
	if hostedCluster == nil {
		return nil, fmt.Errorf("hostedCluster provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(hostedCluster)
	if err != nil {
		return nil, err
	}

	name := hostedCluster.Name
	if name == nil {
		return nil, fmt.Errorf("hostedCluster.Name must be provided to Apply")
	}

	result = &v1alpha1.HostedCluster{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("hostedclusters").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
