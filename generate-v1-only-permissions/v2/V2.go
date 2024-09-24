package v2

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/project-kessel/ksl-schema-language/pkg/intermediate"
	"github.com/project-kessel/ksl-schema-language/pkg/ksl"
)

const (
	V1OnlyPermissionsFile  = "rbac_v1_permissions.json"
	V2PermissionsExtension = "add_v1_based_permission"
	ExtensionV1PermName    = "v1_perm"
)

func ExtractMigratedPermissions(kslSrcPath string) (map[string]bool, error) {
	migratedPerms := map[string]bool{}

	err := filepath.WalkDir(kslSrcPath, func(path string, d fs.DirEntry, err error) error {
		if err != err {
			return err
		}

		if d.IsDir() { //Skip directories
			return nil
		}

		if d.Name() == V1OnlyPermissionsFile { //Skip existing permissions file
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		defer f.Close()

		var ns *intermediate.Namespace
		switch filepath.Ext(path) {
		case ".ksl": //KSL source
			ns, err = ksl.Compile(f)
		case ".json": //Compiled or generated KSIL
			ns, err = intermediate.Load(f)
		}

		if err != nil {
			return err
		}

		for _, e := range ns.ExtensionReferences {
			extractAddedPermissionsFromExtension(migratedPerms, e)
		}

		for _, t := range ns.Types {
			for _, e := range t.Extensions {
				extractAddedPermissionsFromExtension(migratedPerms, e)
			}
			for _, r := range t.Relations {
				for _, e := range r.Extensions {
					extractAddedPermissionsFromExtension(migratedPerms, e)
				}
			}
		}
		return nil
	})

	return migratedPerms, err
}

func extractAddedPermissionsFromExtension(perms map[string]bool, e *intermediate.ExtensionReference) {
	if e.Name != V2PermissionsExtension {
		return
	}

	perm := e.Params[ExtensionV1PermName]

	perms[perm] = true
}

func WriteV1OnlyPermissionsFile(kslSrcPath string, perms []string) error {
	refs := []*intermediate.ExtensionReference{}

	for _, perm := range perms {
		refs = append(refs, &intermediate.ExtensionReference{Namespace: "rbac", Name: "add_v1only_permission", Params: map[string]string{"perm": perm}})
	}

	ns := &intermediate.Namespace{Name: "rbac_v1_permissions", Imports: []string{"rbac"}, ExtensionReferences: refs}

	f, err := os.OpenFile(filepath.Join(kslSrcPath, V1OnlyPermissionsFile), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return intermediate.Store(ns, f)
}
