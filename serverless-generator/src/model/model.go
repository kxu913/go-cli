package model

type MetaData struct {
	Name    string
	Version string
	Prefix  string
}

func MetaDataToMap(metaData *MetaData) map[string]any {

	return map[string]any{
		"name": metaData.Name + "-" + metaData.Version,
		"labels": map[string]any{
			"app":     metaData.Name,
			"version": metaData.Version,
		},
	}

}

type Container struct {
	Image        string
	Port         int32
	Environments []map[string]any
	ForceUpdate  bool
	RunAsRoot    bool
}

func ContainerToMap(container *Container) map[string]any {
	c := map[string]any{
		"image": container.Image,
		"ports": []map[string]any{
			{
				"name":          "http",
				"protocol":      "TCP",
				"containerPort": container.Port,
			},
		},
		"env": container.Environments,
	}
	if container.ForceUpdate {
		c["imagePullPolicy"] = "Always"
	} else {
		c["imagePullPolicy"] = "IfNotPresent"
	}
	if container.RunAsRoot {
		c["securityContext"] = map[string]any{
			"allowPrivilegeEscalation": false,
			"runAsUser":                0,
		}
	}
	return c
}

func ToDeploymentSpec(metaData *MetaData, replicas int, container *Container) map[string]any {
	c := ContainerToMap(container)
	c["name"] = metaData.Name //set container name avoid random string
	return map[string]any{
		"replicas": replicas,
		"selector": map[string]any{
			"matchLabels": map[string]any{
				"app":     metaData.Name,
				"version": metaData.Version,
			},
		},
		"template": map[string]any{
			"metadata": MetaDataToMap(metaData),
			"spec": map[string]any{
				"containers":         []map[string]any{c},
				"serviceAccountName": metaData.Name + "-acct",
			},
		},
	}
}
