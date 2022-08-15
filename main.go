package main
import (
    "os"
    // "fmt"
    "net/http"
    "k8s.io/component-base/cli"
    "github.com/spf13/cobra"
    flag "github.com/spf13/pflag"
    "math/rand"
    "time"
)

var (
    healthy = true
    reliability *float32 = flag.Float32("reliability", 1.0, "target reliability for the health-probe results")
    rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func pingHandler(w http.ResponseWriter, req *http.Request) {
    if healthy && (rng.Float32() < *reliability) {
        w.WriteHeader(200)
    } else {
        w.WriteHeader(500)
    }
}

// TODO(janrous): add feature for marking the server as (un)healthy
// TODO(janrous): add command-line flags for simulating unhealthy

var CmdProbeServer = &cobra.Command{
    Use: "probeserver",
    Short: "Starts server with simple health-probe implementations",
    Args: cobra.MaximumNArgs(0),
    Run: probeServer,
}

func probeServer(cmd *cobra.Command, args []string) {
    http.HandleFunc("/ping", pingHandler)
    // http.HandleFunc("/shutdown, shutdownHandler)
    http.ListenAndServe(":8090", nil)
}
func main() {
    rootCmd := &cobra.Command{
        Use: "app",
    }
    rootCmd.AddCommand(CmdProbeServer)

    code := cli.Run(rootCmd)
    os.Exit(code)
}
