package deployments

import (
	"context"
	"sync"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	dockerclient "github.com/portainer/portainer/api/docker/client"
	k "github.com/portainer/portainer/api/kubernetes"

	"github.com/pkg/errors"
)

type BaseStackDeployer interface {
	DeploySwarmStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry, prune, pullImage bool) error
	DeployComposeStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry, forcePullImage, forceRecreate bool) error
	DeployKubernetesStack(stack *portainer.Stack, endpoint *portainer.Endpoint, user *portainer.User) error
}

type StackDeployer interface {
	BaseStackDeployer
	RemoteStackDeployer
}

type stackDeployer struct {
	lock                *sync.Mutex
	swarmStackManager   portainer.SwarmStackManager
	composeStackManager portainer.ComposeStackManager
	kubernetesDeployer  portainer.KubernetesDeployer
	ClientFactory       *dockerclient.ClientFactory
	dataStore           dataservices.DataStore
}

// NewStackDeployer inits a stackDeployer struct with a SwarmStackManager, a ComposeStackManager and a KubernetesDeployer
func NewStackDeployer(swarmStackManager portainer.SwarmStackManager, composeStackManager portainer.ComposeStackManager,
	kubernetesDeployer portainer.KubernetesDeployer, clientFactory *dockerclient.ClientFactory, dataStore dataservices.DataStore) *stackDeployer {
	return &stackDeployer{
		lock:                &sync.Mutex{},
		swarmStackManager:   swarmStackManager,
		composeStackManager: composeStackManager,
		kubernetesDeployer:  kubernetesDeployer,
		ClientFactory:       clientFactory,
		dataStore:           dataStore,
	}
}
func (d *stackDeployer) DeploySwarmStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry, prune, pullImage bool) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.swarmStackManager.Login(registries, endpoint)
	defer d.swarmStackManager.Logout(endpoint)

	return d.swarmStackManager.Deploy(stack, prune, pullImage, endpoint)
}

func (d *stackDeployer) DeployComposeStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry, forcePullImage, forceRecreate bool) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.swarmStackManager.Login(registries, endpoint)
	defer d.swarmStackManager.Logout(endpoint)

	// --force-recreate doesn't pull updated images
	if forcePullImage {
		if err := d.composeStackManager.Pull(context.TODO(), stack, endpoint); err != nil {
			return err
		}
	}

	err := d.composeStackManager.Up(context.TODO(), stack, endpoint, portainer.ComposeUpOptions{
		ForceRecreate: forceRecreate,
	})
	if err != nil {
		d.composeStackManager.Down(context.TODO(), stack, endpoint)
	}
	return err
}

func (d *stackDeployer) DeployKubernetesStack(stack *portainer.Stack, endpoint *portainer.Endpoint, user *portainer.User) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	appLabels := k.KubeAppLabels{
		StackID:   int(stack.ID),
		StackName: stack.Name,
		Owner:     user.Username,
	}

	if stack.GitConfig == nil {
		appLabels.Kind = "content"
	} else {
		appLabels.Kind = "git"
	}

	k8sDeploymentConfig, err := CreateKubernetesStackDeploymentConfig(stack, d.kubernetesDeployer, appLabels, user, endpoint)
	if err != nil {
		return errors.Wrap(err, "failed to create temp kub deployment files")
	}

	if err := k8sDeploymentConfig.Deploy(); err != nil {
		return errors.Wrap(err, "failed to deploy kubernetes application")
	}

	return nil
}
