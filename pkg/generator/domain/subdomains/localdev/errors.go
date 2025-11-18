package localdev

import "errors"

var (
	// ErrComposeConfigRequired Compose配置必需
	ErrComposeConfigRequired = errors.New("compose configuration is required")

	// ErrInvalidVolumeType 无效的卷类型
	ErrInvalidVolumeType = errors.New("invalid volume type, must be 'bind' or 'volume'")

	// ErrVolumeSourceRequired 卷源路径必需
	ErrVolumeSourceRequired = errors.New("volume source is required")

	// ErrVolumeTargetRequired 卷目标路径必需
	ErrVolumeTargetRequired = errors.New("volume target is required")

	// ErrInvalidK8sVolumeType 无效的K8s卷类型
	ErrInvalidK8sVolumeType = errors.New("invalid k8s volume type, must be one of: configMap, persistentVolumeClaim, emptyDir, hostPath")
)
