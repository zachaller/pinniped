// Copyright 2020 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "go.pinniped.dev/generated/1.17/apis/concierge/login/v1alpha1"
	scheme "go.pinniped.dev/generated/1.17/client/concierge/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TokenCredentialRequestsGetter has a method to return a TokenCredentialRequestInterface.
// A group's client should implement this interface.
type TokenCredentialRequestsGetter interface {
	TokenCredentialRequests() TokenCredentialRequestInterface
}

// TokenCredentialRequestInterface has methods to work with TokenCredentialRequest resources.
type TokenCredentialRequestInterface interface {
	Create(*v1alpha1.TokenCredentialRequest) (*v1alpha1.TokenCredentialRequest, error)
	Update(*v1alpha1.TokenCredentialRequest) (*v1alpha1.TokenCredentialRequest, error)
	UpdateStatus(*v1alpha1.TokenCredentialRequest) (*v1alpha1.TokenCredentialRequest, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.TokenCredentialRequest, error)
	List(opts v1.ListOptions) (*v1alpha1.TokenCredentialRequestList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TokenCredentialRequest, err error)
	TokenCredentialRequestExpansion
}

// tokenCredentialRequests implements TokenCredentialRequestInterface
type tokenCredentialRequests struct {
	client rest.Interface
}

// newTokenCredentialRequests returns a TokenCredentialRequests
func newTokenCredentialRequests(c *LoginV1alpha1Client) *tokenCredentialRequests {
	return &tokenCredentialRequests{
		client: c.RESTClient(),
	}
}

// Get takes name of the tokenCredentialRequest, and returns the corresponding tokenCredentialRequest object, and an error if there is any.
func (c *tokenCredentialRequests) Get(name string, options v1.GetOptions) (result *v1alpha1.TokenCredentialRequest, err error) {
	result = &v1alpha1.TokenCredentialRequest{}
	err = c.client.Get().
		Resource("tokencredentialrequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TokenCredentialRequests that match those selectors.
func (c *tokenCredentialRequests) List(opts v1.ListOptions) (result *v1alpha1.TokenCredentialRequestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.TokenCredentialRequestList{}
	err = c.client.Get().
		Resource("tokencredentialrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested tokenCredentialRequests.
func (c *tokenCredentialRequests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("tokencredentialrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a tokenCredentialRequest and creates it.  Returns the server's representation of the tokenCredentialRequest, and an error, if there is any.
func (c *tokenCredentialRequests) Create(tokenCredentialRequest *v1alpha1.TokenCredentialRequest) (result *v1alpha1.TokenCredentialRequest, err error) {
	result = &v1alpha1.TokenCredentialRequest{}
	err = c.client.Post().
		Resource("tokencredentialrequests").
		Body(tokenCredentialRequest).
		Do().
		Into(result)
	return
}

// Update takes the representation of a tokenCredentialRequest and updates it. Returns the server's representation of the tokenCredentialRequest, and an error, if there is any.
func (c *tokenCredentialRequests) Update(tokenCredentialRequest *v1alpha1.TokenCredentialRequest) (result *v1alpha1.TokenCredentialRequest, err error) {
	result = &v1alpha1.TokenCredentialRequest{}
	err = c.client.Put().
		Resource("tokencredentialrequests").
		Name(tokenCredentialRequest.Name).
		Body(tokenCredentialRequest).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *tokenCredentialRequests) UpdateStatus(tokenCredentialRequest *v1alpha1.TokenCredentialRequest) (result *v1alpha1.TokenCredentialRequest, err error) {
	result = &v1alpha1.TokenCredentialRequest{}
	err = c.client.Put().
		Resource("tokencredentialrequests").
		Name(tokenCredentialRequest.Name).
		SubResource("status").
		Body(tokenCredentialRequest).
		Do().
		Into(result)
	return
}

// Delete takes name of the tokenCredentialRequest and deletes it. Returns an error if one occurs.
func (c *tokenCredentialRequests) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("tokencredentialrequests").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *tokenCredentialRequests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("tokencredentialrequests").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched tokenCredentialRequest.
func (c *tokenCredentialRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TokenCredentialRequest, err error) {
	result = &v1alpha1.TokenCredentialRequest{}
	err = c.client.Patch(pt).
		Resource("tokencredentialrequests").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
