package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

var clientset *kubernetes.Clientset

func buildCommand() *cobra.Command {
	rootCmd := &cobra.Command{Use: "netiedge"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "Path to the kubeconfig file")

	// create clientset
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		kubeconfigPath, err := cmd.Flags().GetString("kubeconfig")
		if err != nil {
			fmt.Println("Error getting kubeconfig flag: ", err)
			return err
		}

		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		if kubeconfigPath != "" {
			loadingRules.ExplicitPath = kubeconfigPath
		}

		clientConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			fmt.Println("Error creating client config: ", err)
			return err
		}

		// Create the clientset
		clientset, err = kubernetes.NewForConfig(clientConfig)
		if err != nil {
			fmt.Println("Error creating clientset: ", err)
			return err
		}

		return nil
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Perform checks on the Kubernetes cluster",
		Run:   check,
	}

	rootCmd.AddCommand(checkCmd)
	return rootCmd
}

func check(cmd *cobra.Command, args []string) {
	checkBasicNodes()
}

func checkBasicNodes() {
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: "core"})
	if err != nil {
		fmt.Printf("Error listing nodes: %v\n", err)
		return
	}
	if len(nodes.Items) >= 3 {
		fmt.Println("Basic check passed: At least 3 nodes with label 'core' found.")
	} else {
		fmt.Printf("Basic check failed: Found %d nodes with label 'core'.\n", len(nodes.Items))
	}
}

func InitAndExecute() {
	cmd := buildCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
