package v1

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ExtractRBACPermissions(rbacJSONPath string, skipApps map[string]bool) ([]string, error) {
	return extractPermissionsFromFolder(rbacJSONPath, skipApps)
}

func extractPermissionsFromFolder(path string, skipApps map[string]bool) ([]string, error) {
	permissions := []string{}

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		perms, err := extractPermissionsFromFile(path, skipApps)
		if err != nil {
			return err
		}

		permissions = append(permissions, perms...)
		return nil
	})

	return permissions, err
}

func extractPermissionsFromFile(path string, skipApps map[string]bool) ([]string, error) {
	perms := []string{}

	fileName := filepath.Base(path)
	appName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	if skipApps[appName] {
		return perms, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file := permFile{}
	json.Unmarshal(data, &file)

	// extract map keys to an array so we can sort them and access values in the map by that order
	resources := make([]string, 0, len(file))
	for resource := range file {
		resources = append(resources, resource)
	}
	sort.Strings(resources)

	for _, resourceName := range resources {
		permissions := file[resourceName]
		for _, permission := range permissions {
			perms = append(perms, fmt.Sprintf("%s:%s:%s", appName, resourceName, permission.Verb))
		}
	}

	return perms, nil
}

type perm struct {
	Verb string `json:"verb"`
}
type permFile map[string][]perm
