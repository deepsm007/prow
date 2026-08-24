/*
Copyright 2026 The Kubernetes Authors.

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

package main

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	prowapi "sigs.k8s.io/prow/pkg/apis/prowjobs/v1"
)

func TestStripProwJobForInformerCache(t *testing.T) {
	heavy := prowJobWithHeavySpec()
	got, err := stripProwJobForInformerCache(heavy)
	if err != nil {
		t.Fatalf("stripProwJobForInformerCache() error = %v", err)
	}
	slim, ok := got.(*prowapi.ProwJob)
	if !ok {
		t.Fatalf("got %T, want *prowapi.ProwJob", got)
	}
	if slim.Spec.PodSpec != nil {
		t.Fatal("expected pod spec to be stripped from cache object")
	}
	if slim.Spec.Refs == nil || slim.Spec.Refs.Org != "openshift" {
		t.Fatalf("expected refs org to be preserved, got %#v", slim.Spec.Refs)
	}
	if len(slim.Spec.Refs.Pulls) != 0 {
		t.Fatal("expected pulls to be stripped from refs")
	}
	if slim.Status.State != prowapi.PendingState {
		t.Fatalf("expected status state %q, got %q", prowapi.PendingState, slim.Status.State)
	}
}

func TestStripProwJobForInformerCacheDeletedFinalState(t *testing.T) {
	heavy := prowJobWithHeavySpec()
	tombstone := cache.DeletedFinalStateUnknown{Obj: heavy}
	got, err := stripProwJobForInformerCache(tombstone)
	if err != nil {
		t.Fatalf("stripProwJobForInformerCache() error = %v", err)
	}
	deleted, ok := got.(cache.DeletedFinalStateUnknown)
	if !ok {
		t.Fatalf("got %T, want cache.DeletedFinalStateUnknown", got)
	}
	slim, ok := deleted.Obj.(*prowapi.ProwJob)
	if !ok {
		t.Fatalf("deleted obj %T, want *prowapi.ProwJob", deleted.Obj)
	}
	if slim.Spec.PodSpec != nil {
		t.Fatal("expected pod spec to be stripped from tombstone object")
	}
}

func prowJobWithHeavySpec() *prowapi.ProwJob {
	return &prowapi.ProwJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pj",
			Namespace: "ci",
			Labels: map[string]string{
				"created-by-prow": "true",
			},
			Annotations: map[string]string{
				"foo": "bar",
			},
		},
		Spec: prowapi.ProwJobSpec{
			Type:  prowapi.PresubmitJob,
			Agent: prowapi.KubernetesAgent,
			Job:   "integration",
			Refs: &prowapi.Refs{
				Org:     "openshift",
				Repo:    "release",
				BaseRef: "main",
				Pulls: []prowapi.Pull{{
					Number: 1,
					Title:  "example",
					SHA:    "deadbeef",
				}},
			},
			PodSpec: &corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "test",
					Image: "busybox",
					Env: []corev1.EnvVar{{
						Name:  "SCRIPT",
						Value: "very long script payload",
					}},
				}},
			},
		},
		Status: prowapi.ProwJobStatus{
			State: prowapi.PendingState,
		},
	}
}
