package main
import (
    "os"
    "fmt"
    "net/http"
    "k8s.io/component-base/cli"
    "github.com/spf13/cobra"
    "math/rand"
    "time"
)

var (
    // TODO(janrous): we can emulate outages every A-B minutes that last for C-D minutes.
    rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)
// TODO(janrous): add feature for marking the server as (un)healthy
// TODO(janrous): add command-line flags for simulating unhealthy

var CmdProbeServer = &cobra.Command{
    Use: "probeserver",
    Short: "Starts server with simple health-probe implementations",
    Args: cobra.MaximumNArgs(0),
    Run: probeServer,
}

type ServerState struct {
    Reliability float32
    Healthy bool
}


func probeServer(cmd *cobra.Command, args []string) {
    rel, err := cmd.Flags().GetFloat32("reliability")
    if err != nil {
        panic(err)
    }
    state := ServerState{
        Reliability: rel,
        Healthy: true,
    }
    healthFn := func(w http.ResponseWriter, req *http.Request) {
        resp := 500
        if state.Healthy && (rng.Float32() < state.Reliability) {
            resp = 200
        }
        w.WriteHeader(resp)
        fmt.Printf("health-check on %s (by %s) responded %d\n", req.URL.String(), req.UserAgent(), resp)
    }
    healthEps, err := cmd.Flags().GetStringArray("healthEndpoint")
    for _, endpoint := range healthEps {
        fmt.Printf("Serving health endpoint on /%s\n", endpoint)
        http.HandleFunc(fmt.Sprintf("/%s", endpoint), healthFn)
    }
    http.HandleFunc("/startoutage", func(w http.ResponseWriter, req *http.Request) {
        state.Healthy = false
        w.WriteHeader(200)
    })
    http.HandleFunc("/endoutage", func(w http.ResponseWriter, req *http.Request) {
        state.Healthy = true
        w.WriteHeader(200)
    })
    // http.HandleFunc("/shutdown, shutdownHandler)
    http.ListenAndServe(":8090", nil)
}
func main() {
    rootCmd := &cobra.Command{
        Use: "app",
    }
    rootCmd.AddCommand(CmdProbeServer)
    CmdProbeServer.PersistentFlags().Float32("reliability", 0.95, "target reliability for the health-probe results")
    CmdProbeServer.PersistentFlags().StringArray("healthEndpoint", []string{"ping"}, "list of endpoints to server health on")

    code := cli.Run(rootCmd)
    os.Exit(code)
}
