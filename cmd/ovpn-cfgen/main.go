package main

func main() {
	rootCmd.AddCommand(buildCACmd)
	rootCmd.AddCommand(buildKeyServerCmd)

	rootCmd.Execute()
}
