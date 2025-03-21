package v2

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	v1 "github.com/RedHatInsights/rbac-config-actions/generatepermissions/v1"
	"github.com/project-kessel/ksl-schema-language/pkg/intermediate"
)

const (
	V1OnlyPermissionsFile  = "rbac_v1_permissions.json"
	V2PermissionsExtension = "add_v1_based_permission"
	ExtensionV1PermName    = "v1_perm"
)

func GetMigratedApps(kslPath string) (map[string]bool, error) {
	return readAppList(filepath.Join(kslPath, "migrated_apps.lst"))
}

func GetHostOnlyApps(kslPath string) (map[string]bool, error) {
	return readAppList(filepath.Join(kslPath, "hostsonly_apps.lst"))
}

func WriteV1OnlyPermissionsFile(kslPath string, hostsonly_apps map[string]bool, perms []v1.Permission) error {
	refs := []*intermediate.ExtensionReference{}

	for _, perm := range perms {
		v2_perm := to_v2_perm(perm)

		if hostsonly_apps[perm.App] {
			if !perm.IsWildcard() {
				/*
					refs = append(refs, &intermediate.ExtensionReference{Namespace: "rbac", Name: "add_v1_based_permission", Params: map[string]string{"app": escape_v1_string(perm.App), "resource": escape_v1_string(perm.ResourceType), "verb": escape_v1_string(perm.Verb), "v2_perm": v2_perm}})
					refs = append(refs, &intermediate.ExtensionReference{Namespace: "hbi", Name: "expose_host_permission", Params: map[string]string{"v2_perm": v2_perm, "host_perm": to_host_perm(perm)}})
				*/

				v2_assigned_perm := v2_perm + "_assigned"
				refs = append(refs, &intermediate.ExtensionReference{Namespace: "rbac", Name: "add_v1_based_permission", Params: map[string]string{"app": escape_v1_string(perm.App), "resource": escape_v1_string(perm.ResourceType), "verb": escape_v1_string(perm.Verb), "v2_perm": v2_assigned_perm}})
				refs = append(refs, &intermediate.ExtensionReference{Namespace: "rbac", Name: "add_contingent_permission", Params: map[string]string{"first": "inventory_host_view", "second": v2_assigned_perm, "contingent": v2_perm}})
				refs = append(refs, &intermediate.ExtensionReference{Namespace: "hbi", Name: "expose_host_permission", Params: map[string]string{"v2_perm": v2_perm, "host_perm": to_host_perm(perm)}})
			}
		} else {
			refs = append(refs, &intermediate.ExtensionReference{Namespace: "rbac", Name: "add_v1only_permission", Params: map[string]string{"perm": v2_perm}})
		}
	}

	ns := &intermediate.Namespace{Name: "rbac_v1_permissions", Imports: []string{"rbac", "hbi"}, ExtensionReferences: refs}

	f, err := os.OpenFile(filepath.Join(kslPath, "src", V1OnlyPermissionsFile), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return intermediate.Store(ns, f)
}

func to_v2_perm(perm v1.Permission) string {
	return fmt.Sprintf("%s_%s_%s", escape_v1_string(perm.App), escape_v1_string(perm.Verb), escape_v1_string(perm.ResourceType))
}

func to_host_perm(perm v1.Permission) string {
	return fmt.Sprintf("%s_%s", escape_v1_string(perm.Verb), escape_v1_string(perm.ResourceType))
}

func escape_v1_string(val string) string {
	val = strings.ReplaceAll(val, "-", "_")
	val = strings.ReplaceAll(val, ".", "_")
	val = strings.ReplaceAll(val, "*", "all")

	return val
}

func readAppList(filePath string) (map[string]bool, error) {
	apps := map[string]bool{}

	listFile, err := os.Open(filePath)
	if err != nil {
		return apps, err
	}
	defer listFile.Close()

	scanner := bufio.NewScanner(listFile)
	for scanner.Scan() {
		apps[scanner.Text()] = true
	}

	return apps, scanner.Err()
}
