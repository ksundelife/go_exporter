package main

import (
    "log"
    "fmt"
    "os"
    "strconv"
    "time"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/OpenNebula/one/src/oca/go/src/goca"
)

func recordMetrics() {
    go func() {
        for {
            opsProcessed.Inc()
            time.Sleep(2 * time.Second)
        }
    }()
}

var (
    opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
        Name: "myapp_processed_ops_total",
        Help: "The total number of processed events",
    })
)

func main() {
    recordMetrics()

    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":8082", nil))
        // Initialize connection with OpenNebula
    con := map[string]string{
        "user": "astra",
        "password": "12345678",
        "endpoint": "192.168.100.14",
    }

    client := goca.NewDefaultClient(
        goca.NewConfig(con["user"], con["password"], con["endpoint"]),
    )

    controller := goca.NewController(client)

    // Read VM ID from arguments
    id, _ := strconv.Atoi(os.Args[1])

    vmctrl := controller.VM(id)

    // Fetch informations of the created VM
    vm, err := vmctrl.Info(false)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Shutting down %+v\n", vm.Name)

    // Poweroff the VM
    err = vmctrl.Poweroff()
    if err != nil {
        log.Fatal(err)
    }

    // Get short informations of the VMs
    vms, err := controller.VMs().Info()
    if err != nil {
        log.Fatal(err)
    }
    for i := 0; i < len(vms.VMs); i++ {
        // This Info method, per VM instance, give us detailed informations on the instance
        // Check xsd files to see the difference
        vm, err := controller.VM(vms.VMs[i].ID).Info(false)
        if err != nil {
            log.Fatal(err)
        }
        //Do some others stuffs on vm
        fmt.Printf("%+v\n", vm)
    }
}