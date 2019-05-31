package main

func main() {
	rootCmd.AddCommand(buildCACmd)
	rootCmd.AddCommand(buildKeyServerCmd)
	rootCmd.AddCommand(buildKeyCmd)
	rootCmd.AddCommand(serverConfigCmd)
	rootCmd.AddCommand(clientConfigCmd)

	rootCmd.Execute()
}
