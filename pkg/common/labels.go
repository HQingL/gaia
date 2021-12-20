/*
Copyright 2021 The Clusternet Authors.

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

package common

// label key
const (
	ObjectCreatedByLabel    = "clusternet.io/created-by"
	ObjectControlledByLabel = "apps.clusternet.io/owned-by"

	// the source info where this object belongs to or controlled by
	ConfigGroupLabel     = "apps.clusternet.io/config.group"
	ConfigVersionLabel   = "apps.clusternet.io/config.version"
	ConfigKindLabel      = "apps.clusternet.io/config.kind"
	ConfigNameLabel      = "apps.clusternet.io/config.name"
	ConfigNamespaceLabel = "apps.clusternet.io/config.namespace"
	ConfigUIDLabel       = "apps.clusternet.io/config.uid"

	ConfigSubscriptionUIDLabel       = "apps.clusternet.io/subs.uid"
	ConfigSubscriptionNameLabel      = "apps.clusternet.io/subs.name"
	ConfigSubscriptionNamespaceLabel = "apps.clusternet.io/subs.namespace"
)

// label value
const (
	RBACDefaults    = "rbac-defaults"

	GaiaControllerManager   = "gaia-controller-manager"
)