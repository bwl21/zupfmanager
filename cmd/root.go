package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"   // Standardwert für lokale Builds
	GitCommit = "dirty" // Standardwert für uncommitted changes
)

// rootCmd represents the base command when called without any subcommands
// This is the root command for the zupfmanager application.
// It provides the entry point for all other commands.
var rootCmd = &cobra.Command{
	Use:   "zupfmanager",
	Short: "Tool um Notenmappen mit Zupfnoter zusammenzustellen",
	Long: `
Zupfmanager ist ein Werkzeug zum Verwalten und Erstellen von Musikprojekten, 
insbesondere für Zupfinstrumente.

Es hilft Ihnen, Ihre Lieder zu organisieren, Projektkonfigurationen zu verwalten und 
Notenblätter mit dem Zupfnoter-Renderer zu generieren.

Zupfmanager unterstützt:
- Erstellen, Auflisten, Aktualisieren und Löschen von Musikprojekten
- Hinzufügen, Entfernen und Verwalten von Lieder innerhalb von Projekten
- Erstellen von Notenblättern
- Importieren von ABC-Dateien
- Interaktion mit einer Terminalbenutzeroberfläche`,
	RunE: func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetBool("version")
		if version {
			fmt.Printf("zupfmanager %s\nCommit: %s\n", Version, GitCommit)
			os.Exit(0)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init is called after all imported packages have been initialized.
func init() {
	// Add the version flag
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")

	// Enable shell completion
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `Erzeugt ein Completion-Skript für verschiedene Shells.

Zum Laden der Completion:

Bash:

  $ source <(zupfmanager completion bash)

  # Um die Completion für jede Sitzung zu laden, einmalig ausführen:
  $ zupfmanager completion bash > /usr/local/etc/bash_completion.d/zupfmanager

Zsh:

  $ source <(zupfmanager completion zsh)

  # Wenn die Shell-Completion noch nicht in Ihrer Umgebung aktiviert ist,
  # müssen Sie sie aktivieren. Sie können Folgendes einmalig ausführen:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # Um die Completion für jede Sitzung zu laden, einmalig ausführen:
  $ zupfmanager completion zsh > "${fpath[1]}/_zupfmanager"

  # Sie müssen eine neue Shell starten, damit diese Einrichtung wirksam wird.

fish:

  $ zupfmanager completion fish | source

  # Um die Completion für jede Sitzung zu laden, einmalig ausführen:
  $ zupfmanager completion fish > ~/.config/fish/completions/zupfmanager.fish
PowerShell:

  > zupfmanager completion powershell | Out-String | Invoke-Expression

  # Um die Completion für jede neue Sitzung zu laden, ausführen:
  > zupfmanager completion powershell > zupfmanager.ps1
  # und diese Datei von Ihrem PowerShell-Profil aus sourcen.
		`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Generate the completion script based on the shell type
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				completionDir := filepath.Join(os.Getenv("HOME"), ".zsh", "completion")
				err := os.MkdirAll(completionDir, 0755)
				if err != nil {
					fmt.Println("Error creating zsh completion directory:", err)
					os.Exit(1)
				}
				completionFile := filepath.Join(completionDir, "_zupfmanager")
				if err := rootCmd.GenZshCompletionFile(completionFile); err != nil {
					fmt.Println("Error generating zsh completion file:", err)
					os.Exit(1)
				}
				fmt.Println("zsh completion file generated at", completionFile)

			case "fish":
				rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			}
			os.Exit(0)
		},
	}
	rootCmd.AddCommand(completionCmd)
}
