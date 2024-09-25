package v2

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/project-kessel/ksl-schema-language/pkg/intermediate"
)

const (
	V1OnlyPermissionsFile  = "rbac_v1_permissions.json"
	V2PermissionsExtension = "add_v1_based_permission"
	ExtensionV1PermName    = "v1_perm"
)

func GetMigratedApps(kslPath string) (map[string]bool, error) {
	apps := map[string]bool{}

	listFile, err := os.Open(filepath.Join(kslPath, "migrated_apps.lst"))
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

func WriteV1OnlyPermissionsFile(kslPath string, perms []string) error {
	refs := []*intermediate.ExtensionReference{}

	for _, perm := range perms {
		refs = append(refs, &intermediate.ExtensionReference{Namespace: "rbac", Name: "add_v1only_permission", Params: map[string]string{"perm": perm}})
	}

	ns := &intermediate.Namespace{Name: "rbac_v1_permissions", Imports: []string{"rbac"}, ExtensionReferences: refs}

	f, err := os.OpenFile(filepath.Join(kslPath, "src", V1OnlyPermissionsFile), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return intermediate.Store(ns, f)
}
