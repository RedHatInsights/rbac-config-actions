package main

import (
	"flag"
	"strings"

	v1 "github.com/RedHatInsights/rbac-config-actions/generatepermissions/v1"
	v2 "github.com/RedHatInsights/rbac-config-actions/generatepermissions/v2"
)

func main() {
	kslSrc := flag.String("ksl", "", "The path to the ksl project directory (where the migrated_apps.lst file is)")
	rbacPermissions := flag.String("rbac-permissions-json", "", "The path to the directory containing RBAC permissions .json files for the current environment.")

	flag.Parse()

	if kslSrc != nil && *kslSrc == "" {
		kslSrc = nil
	}
	if rbacPermissions != nil && *rbacPermissions == "" {
		rbacPermissions = nil
	}

	if kslSrc == nil || rbacPermissions == nil {
		flag.Usage()
		return
	}

	migratedApps, err := v2.GetMigratedApps(*kslSrc)
	if err != nil {
		panic(err)
	}

	hostsonlyApps, err := v2.GetHostOnlyApps(*kslSrc)
	if err != nil {
		panic(err)
	}

	perms, err := v1.ExtractRBACPermissions(*rbacPermissions, migratedApps)
	if err != nil {
		panic(err)
	}

	v2.WriteV1OnlyPermissionsFile(*kslSrc, hostsonlyApps, perms)
}

func v1PermToV2Perm(v1Perm string) string {
	s := strings.ReplaceAll(v1Perm, ":", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "*", "all")

	return s
}
