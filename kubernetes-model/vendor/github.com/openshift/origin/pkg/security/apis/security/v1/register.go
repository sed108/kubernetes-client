/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package v1

import (
	"github.com/openshift/api/security/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName       = "security.openshift.io"
	LegacyGroupName = ""
)

// SchemeGroupVersion is group version used to register these objects
var (
	SchemeGroupVersion       = schema.GroupVersion{Group: GroupName, Version: "v1"}
	LegacySchemeGroupVersion = schema.GroupVersion{Group: LegacyGroupName, Version: "v1"}

	LegacySchemeBuilder    = runtime.NewSchemeBuilder(v1.LegacySchemeBuilder.AddToScheme, addConversionFuncs, addDefaultingFuncs, RegisterDefaults, RegisterConversions)
	AddToSchemeInCoreGroup = LegacySchemeBuilder.AddToScheme

	SchemeBuilder = runtime.NewSchemeBuilder(v1.SchemeBuilder.AddToScheme, addConversionFuncs, addDefaultingFuncs)
	AddToScheme   = SchemeBuilder.AddToScheme

	localSchemeBuilder = &SchemeBuilder
)

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
