package kubernetes

import (
	"fmt"

	"github.com/debdutdeb/velero-plugin-exclude-with-ownerref/internal/config"
	"github.com/vmware-tanzu/velero/pkg/plugin/framework/common"
	"go.yaml.in/yaml/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClient() (*kubernetes.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	clientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("getting client config: %w", err)
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	return client, nil
}

func GetOurConfig() (*config.Config, error) {
	config := &config.Config{
		ExcludeResourcesWithOwnerReferences: false,
		ExcludeNamesRegexes:                 []*config.RegexWrapper{},
		ExcludeRegexesForKinds:              []config.ExcludeRegexForKind{},
	}

	cl, err := GetClient()
	if err != nil {
		return nil, fmt.Errorf("getting kubernetes client: %w", err)
	}

	configMapClient := cl.CoreV1().ConfigMaps("velero")

	pluginConfig, err := common.GetPluginConfig(common.PluginKindBackupItemActionV2, "debdutdeb.com/velero-plugin-exclude-with-ownerref", configMapClient)
	if err != nil {
		return nil, fmt.Errorf("getting plugin config: %w", err)
	}

	err = yaml.Unmarshal([]byte(pluginConfig.Data["config"]), config)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling plugin config: %w", err)
	}

	return config, nil
}
