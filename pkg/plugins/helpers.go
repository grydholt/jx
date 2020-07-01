package plugins

import (
	"fmt"
	"os"
	"strings"

	jenkinsv1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-helpers/pkg/extensions"
	"github.com/jenkins-x/jx-helpers/pkg/homedir"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetJXPlugin returns the path to the locally installed jx plugin
func GetJXPlugin(name string, version string) (string, error) {
	pluginBinDir, err := PluginBinDir()
	if err != nil {
		return "", err
	}
	plugin := CreateJXPlugin(name, version)
	return extensions.EnsurePluginInstalled(plugin, pluginBinDir)
}

// CreateJXPlugin creates the jx plugin
func CreateJXPlugin(name, version string) jenkinsv1.Plugin {
	binaries := extensions.CreateBinaries(func(p extensions.Platform) string {
		return fmt.Sprintf("https://storage.googleapis.com/cloudbees-jx-plugins/plugin/%s/%s/jx-%s-%s-%s.%s", name, version, name, strings.ToLower(p.Goos), strings.ToLower(p.Goarch), p.Extension())
	})

	plugin := jenkinsv1.Plugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: jenkinsv1.PluginSpec{
			SubCommand:  name,
			Binaries:    binaries,
			Description: name + "  binary",
			Name:        "jx-" + name,
			Version:     version,
		},
	}
	return plugin
}

// PluginBinDir returns where the binary plugins are installed
func PluginBinDir() (string, error) {
	pluginBinDir, err := homedir.PluginBinDir(os.Getenv("JX_CLI_HOME"), ".jx-cli")
	if err != nil {
		return "", errors.Wrapf(err, "failed to find plugin home dir")
	}
	return pluginBinDir, nil
}
