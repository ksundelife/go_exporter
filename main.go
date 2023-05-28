package main

import (
    "log"
    "fmt"
    "time"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    //"github.com/OpenNebula/one/src/oca/go/src/goca"
    "github.com/megamsys/opennebula-go/api"
    "github.com/megamsys/opennebula-go/compute"
    "github.com/megamsys/opennebula-go/template"
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
    fmt.Println("Hello world!")
    cm := make(map[string]string)
    cm[api.ENDPOINT] = "http://192.168.0.118:2633/RPC2"
    cm[api.USERID] = "oneadmin"
    cm[api.PASSWORD] = "oneadmin"

    cl, _ := api.NewClient(cm)
    v := compute.VirtualMachine {
        Name: "testmegam4",
        TemplateName: "megam",
        Cpu: "1",
        Memory: "1024",
        Image: "megam",
        ClusterId: "100" ,
        T: cl,
        ContextMap: map[string]string{"assembly_id": "ASM-007", "assemblies_id": "AMS-007"},
        Vnets: map[string]string{"0":"ipv4-pub"},
    } //memory in terms of MB! duh!

    response, err := v.Create(template.UserTemplates{})
    if err != nil {
        // handle error
    }

    vmid := response.(string)
    fmt.Println("VirtualMachine created successfully", vmid)
}