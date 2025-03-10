import { useMutation, useQueryClient } from '@tanstack/react-query';
import { ConfigMap } from 'kubernetes-types/core/v1';

import axios from '@/portainer/services/axios';
import { EnvironmentId } from '@/react/portainer/environments/types';
import { withInvalidate } from '@/react-tools/react-query';

import { parseKubernetesAxiosError } from '../../axiosError';

import { configMapQueryKeys } from './useK8sConfigMaps';

/**
 * useUpdateK8sConfigMapMutation returns a mutation hook for updating a Kubernetes ConfigMap using the Kubernetes proxy API.
 */
export function useUpdateK8sConfigMapMutation(
  environmentId: EnvironmentId,
  namespace: string
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      data,
      configMapName,
    }: {
      data: ConfigMap;
      configMapName: string;
    }) => updateConfigMap(environmentId, namespace, configMapName, data),
    ...withInvalidate(queryClient, [
      configMapQueryKeys.configMaps(environmentId, namespace),
    ]),
    // handle success notifications in the calling component
  });
}

async function updateConfigMap(
  environmentId: EnvironmentId,
  namespace: string,
  configMap: string,
  data: ConfigMap
) {
  try {
    return await axios.put(
      `/endpoints/${environmentId}/kubernetes/api/v1/namespaces/${namespace}/configmaps/${configMap}`,
      data
    );
  } catch (e) {
    throw parseKubernetesAxiosError(
      e,
      `Unable to update ConfigMap '${configMap}'`
    );
  }
}
