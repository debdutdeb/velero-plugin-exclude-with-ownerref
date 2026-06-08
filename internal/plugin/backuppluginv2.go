/*
Copyright the Velero contributors.

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

package plugin

import (
	"github.com/debdutdeb/velero-plugin-exclude-with-ownerref/internal/config"
	"github.com/debdutdeb/velero-plugin-exclude-with-ownerref/internal/kubernetes"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
)

var cfg *config.Config

func init() {
	// var err error
	// cfg, err = kubernetes.GetOurConfig()
	// if err != nil {
	// 	logrus.WithError(err).Fatal("Error getting our config")
	// }
}

// BackupPluginV2 is a v2 backup item action plugin for Velero.
type BackupPluginV2 struct {
	log logrus.FieldLogger
}

// NewBackupPluginV2 instantiates a v2 BackupPlugin.
func NewBackupPluginV2(log logrus.FieldLogger) *BackupPluginV2 {
	return &BackupPluginV2{log: log}
}

// Name is required to implement the interface, but the Velero pod does not delegate this
// method -- it's used to tell velero what name it was registered under. The plugin implementation
// must define it, but it will never actually be called.
func (p *BackupPluginV2) Name() string {
	return "exclude-with-ownerref"
}

// AppliesTo returns information about which resources this action should be invoked for.
// The IncludedResources and ExcludedResources slices can include both resources
// and resources with group names. These work: "ingresses", "ingresses.extensions".
// A BackupPlugin's Execute function will only be invoked on items that match the returned
// selector. A zero-valued ResourceSelector matches all resources.
func (p *BackupPluginV2) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{}, nil // for all
}

// Execute allows the ItemAction to perform arbitrary logic with the item being backed up,
// in this case, setting a custom annotation on the item being backed up.
func (p *BackupPluginV2) Execute(item runtime.Unstructured, backup *v1.Backup) (runtime.Unstructured, []velero.ResourceIdentifier, string, []velero.ResourceIdentifier, error) {
	p.log.Info("Exclude with OwnerRef backup plugin")

	var err error
	if cfg == nil {
		cfg, err = kubernetes.GetOurConfig()
		if err != nil {
			return nil, nil, "", nil, err
		}
	}

	metadata, err := meta.Accessor(item)
	if err != nil {
		return nil, nil, "", nil, err
	}

	gvk, err := meta.TypeAccessor(item)
	if err != nil {
		return nil, nil, "", nil, err
	}

	ownerReferences := metadata.GetOwnerReferences()
	p.log.Infof("Resource name: %v, namespace: %v, kind: %v, ownerReferences: %v", metadata.GetName(), metadata.GetNamespace(), gvk.GetKind(), ownerReferences)

	if cfg.ExcludeResourcesWithOwnerReferences && len(ownerReferences) > 0 {
		p.log.Infof("Skipping resource: %v, namespace: %v, kind: %v", metadata.GetName(), metadata.GetNamespace(), gvk.GetKind())
		return nil, nil, "", nil, nil
	}

	for _, regex := range cfg.ExcludeNamesRegexes {
		if regex.MatchString(metadata.GetName()) {
			p.log.Infof("Skipping resource: %v, namespace: %v, kind: %v", metadata.GetName(), metadata.GetNamespace(), gvk.GetKind())
			return nil, nil, "", nil, nil
		}
	}

	for _, regex := range cfg.ExcludeRegexesForKinds {
		if regex.Kind == gvk.GetKind() {
			for _, regex := range regex.Regexes {
				if regex.MatchString(metadata.GetName()) {
					p.log.Infof("Skipping resource: %v, namespace: %v, kind: %v", metadata.GetName(), metadata.GetNamespace(), gvk.GetKind())
					return nil, nil, "", nil, nil
				}
			}
		}
	}

	return item, nil, "", nil, nil
}

func (p *BackupPluginV2) Progress(operationID string, backup *v1.Backup) (velero.OperationProgress, error) {
	return velero.OperationProgress{}, nil
}

func (p *BackupPluginV2) Cancel(operationID string, backup *v1.Backup) error {
	return nil
}
