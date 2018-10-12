/*
Copyright 2018 The Jetstack cert-manager contributors.

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

package v1alpha1

import (
	v1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ChallengeLister helps list Challenges.
type ChallengeLister interface {
	// List lists all Challenges in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Challenge, err error)
	// Challenges returns an object that can list and get Challenges.
	Challenges(namespace string) ChallengeNamespaceLister
	ChallengeListerExpansion
}

// challengeLister implements the ChallengeLister interface.
type challengeLister struct {
	indexer cache.Indexer
}

// NewChallengeLister returns a new ChallengeLister.
func NewChallengeLister(indexer cache.Indexer) ChallengeLister {
	return &challengeLister{indexer: indexer}
}

// List lists all Challenges in the indexer.
func (s *challengeLister) List(selector labels.Selector) (ret []*v1alpha1.Challenge, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Challenge))
	})
	return ret, err
}

// Challenges returns an object that can list and get Challenges.
func (s *challengeLister) Challenges(namespace string) ChallengeNamespaceLister {
	return challengeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ChallengeNamespaceLister helps list and get Challenges.
type ChallengeNamespaceLister interface {
	// List lists all Challenges in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Challenge, err error)
	// Get retrieves the Challenge from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Challenge, error)
	ChallengeNamespaceListerExpansion
}

// challengeNamespaceLister implements the ChallengeNamespaceLister
// interface.
type challengeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Challenges in the indexer for a given namespace.
func (s challengeNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Challenge, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Challenge))
	})
	return ret, err
}

// Get retrieves the Challenge from the indexer for a given namespace and name.
func (s challengeNamespaceLister) Get(name string) (*v1alpha1.Challenge, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("challenge"), name)
	}
	return obj.(*v1alpha1.Challenge), nil
}
