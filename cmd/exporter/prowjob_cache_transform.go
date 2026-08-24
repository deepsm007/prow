/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed under the License is distributed on an "AS IS" BASIS,
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"k8s.io/client-go/tools/cache"

	prowapi "sigs.k8s.io/prow/pkg/apis/prowjobs/v1"
)

// stripProwJobForInformerCache drops large prowjob fields that the exporter does
// not read. The informer otherwise retains full pod specs for every prowjob in
// the cluster, which drives multi-GB RSS on busy app.ci instances.
func stripProwJobForInformerCache(obj interface{}) (interface{}, error) {
	switch t := obj.(type) {
	case *prowapi.ProwJob:
		return slimProwJob(t), nil
	case cache.DeletedFinalStateUnknown:
		if pj, ok := t.Obj.(*prowapi.ProwJob); ok {
			t.Obj = slimProwJob(pj)
		}
		return t, nil
	default:
		return obj, nil
	}
}

func slimProwJob(pj *prowapi.ProwJob) *prowapi.ProwJob {
	if pj == nil {
		return nil
	}
	slim := *pj
	slim.Spec = prowapi.ProwJobSpec{
		Type:      pj.Spec.Type,
		Agent:     pj.Spec.Agent,
		Job:       pj.Spec.Job,
		Refs:      slimRefs(pj.Spec.Refs),
		ExtraRefs: slimExtraRefs(pj.Spec.ExtraRefs),
	}
	slim.Status = prowapi.ProwJobStatus{
		StartTime:      pj.Status.StartTime,
		PendingTime:    pj.Status.PendingTime,
		CompletionTime: pj.Status.CompletionTime,
		State:          pj.Status.State,
	}
	return &slim
}

func slimRefs(refs *prowapi.Refs) *prowapi.Refs {
	if refs == nil {
		return nil
	}
	return &prowapi.Refs{
		Org:     refs.Org,
		Repo:    refs.Repo,
		BaseRef: refs.BaseRef,
	}
}

func slimExtraRefs(refs []prowapi.Refs) []prowapi.Refs {
	if len(refs) == 0 {
		return nil
	}
	out := make([]prowapi.Refs, len(refs))
	for i := range refs {
		if slim := slimRefs(&refs[i]); slim != nil {
			out[i] = *slim
		}
	}
	return out
}
