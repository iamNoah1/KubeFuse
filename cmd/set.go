/*
Copyright © 2025 Noah Ispas <noahispas.public@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"kubefuse/internal/app"
	"strings"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <kind/name> <path=value> [path=value]...",
	Short: "Patch fields on a resource with optional TTL and audit",
	Args:  cobra.MinimumNArgs(2),
	Long:  `Patch fields on a resource with optional TTL and audit`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		namespace, _ := cmd.Flags().GetString("namespace")

		if strings.Contains(toComplete, "/") {
			parts := strings.SplitN(toComplete, "/", 2)
			kind := parts[0]
			namePrefix := parts[1]

			names, err := app.ListResourceNames(cmd.Context(), kind, namespace)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			completions := make([]string, 0, len(names))
			for _, name := range names {
				if namePrefix != "" && !strings.HasPrefix(name, namePrefix) {
					continue
				}
				completions = append(completions, kind+"/"+name)
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		}

		kinds, err := app.ListResourceKinds()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		completions := make([]string, 0, len(kinds))
		for _, kind := range kinds {
			if toComplete != "" && !strings.HasPrefix(kind, toComplete) {
				continue
			}
			completions = append(completions, kind+"/")
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		targetRaw := args[0]
		patchesRaw := args[1:]
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}

		reason, err := cmd.Flags().GetString("reason")
		if err != nil {
			return err
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}

		ttl, err := cmd.Flags().GetString("ttl")
		if err != nil {
			return err
		}

		dto := app.SetDTO{
			TargetRaw:     targetRaw,
			PatchesRaw:    patchesRaw,
			NamespaceFlag: namespace,
			Reason:        reason,
			TTL:           ttl,
			DryRun:        dryRun,
		}

		return app.SetHandler(dto)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	setCmd.Flags().StringP("namespace", "n", "", "K8s namespace")
	setCmd.Flags().StringP("ttl", "t", "10m", "Time to life before the change gets rolled back")
	setCmd.Flags().StringP("reason", "r", "test", "Reason for the patch, gets added to k8s resource annotation")
	setCmd.Flags().Bool("dry-run", false, "Preview the patch without applying it")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
