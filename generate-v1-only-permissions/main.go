package main

import (
	"flag"
	"strings"

	v1 "github.com/RedHatInsights/rbac-config-actions/generatepermissions/v1"
	v2 "github.com/RedHatInsights/rbac-config-actions/generatepermissions/v2"
)

func main() {
	kslSrc := flag.String("ksl-src", "", "The path to the directory containing .ksl source files and .json precompiled files for the current environment.")
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

	migratedPerms, err := v2.ExtractMigratedPermissions(*kslSrc)
	if err != nil {
		panic(err)
	}

	allPerms, err := v1.ExtractRBACPermissions(*rbacPermissions)
	if err != nil {
		panic(err)
	}

	v1OnlyPerms := make([]string, 0, len(allPerms))
	for _, perm := range allPerms {
		v2Perm := v1PermToV2Perm(perm)
		if !migratedPerms[v2Perm] {
			v1OnlyPerms = append(v1OnlyPerms, v2Perm)
		}
	}

	v2.WriteV1OnlyPermissionsFile(*kslSrc, v1OnlyPerms)
}

func v1PermToV2Perm(v1Perm string) string {
	s := strings.ReplaceAll(v1Perm, ":", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "*", "all")

	return s
}
