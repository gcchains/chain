package chainmetrics

import (
	"net"

	"github.com/gcchains/chain/commons/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	// gauge items
	blockNumberCounter = prometheus.NewGauge(prometheus.GaugeOpts{Name: "gcchain_block_number",
		Help: "current blockNumber."})

	txsNumberCounter = prometheus.NewGauge(prometheus.GaugeOpts{Name: "gcchain_txs_number",
		Help: "current txsNumber."})

	insertionElapsedTime = prometheus.NewGauge(prometheus.GaugeOpts{Name: "gcchain_insertion_elapsed_time",
		Help: "current insertion elapsed time."})
)

// configuration items
var (
	chainId        = ""
	gatewayAddress = ""
)

func InitMetrics(port, gatewayURL string) {
	ips := getIP()
	log.Debug("InitMetrics", "ips", ips)
	if len(ips) > 0 {
		chainId = ips[0] + ":" + port
	}
	gatewayAddress = gatewayURL
	log.Debug("InitMetrics", "chainId", chainId, "gatewayAddress", gatewayAddress)
}

func NeedMetrics() bool {
	return gatewayAddress != ""
}

func ReportBlockNumberGauge(exportedJob string, blockNumber float64) {
	blockNumberCounter.Set(blockNumber)
	reportGauge(gatewayAddress, exportedJob, chainId, blockNumberCounter)
}

func ReportTxsNumberGauge(exportedJob string, txsNumber float64) {
	txsNumberCounter.Set(txsNumber)
	reportGauge(gatewayAddress, exportedJob, chainId, txsNumberCounter)
}

func ReportInsertionElapsedTime(exportedJob string, elapsed float64) {
	insertionElapsedTime.Set(elapsed)
	reportGauge(gatewayAddress, exportedJob, chainId, insertionElapsedTime)
}

func reportGauge(monitorURL, exportedJob, host string, gauge prometheus.Gauge) {
	if err := push.New(monitorURL, exportedJob).
		Collector(gauge).
		Grouping("host", host).
		Push(); err != nil {
		log.Error("Could not push blockNumber to Pushgateway.", "error", err)
	}
}

func getIP() []string {
	ips := []string{}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Error(err.Error())
	}
	for _, address := range addresses {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				ips = append(ips, ip.IP.String())
			}
		}
	}
	return ips
}
