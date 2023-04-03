// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"fmt"
	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
)

var (
	client          *bssopenapi20171214.Client
	MetricsName     = "alibaba_billing_package"
	GaugeVecTotal   *prometheus.GaugeVec
	GaugeVecPercent *prometheus.GaugeVec
	LABELS          = []string{"region", "remark", "instance_id", "status", "name"}
	handler         = promhttp.Handler()
)

func init() {
	v := viper.New()
	v.AutomaticEnv()
	v.MustBindEnv(
		"ACCESS_KEY",
		"SECRET_KEY",
		"REGION",
		"ENDPOINT",
	)
	v.SetDefault("REGION", "us-east-1")
	var err error
	client, err = bssopenapi20171214.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(v.GetString("ACCESS_KEY")),
		AccessKeySecret: tea.String(v.GetString("SECRET_KEY")),
		RegionId:        tea.String(v.GetString("REGION")),
	})
	if err != nil {
		panic(err)
	}

	GaugeVecTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_total", MetricsName),
		}, LABELS)
	GaugeVecPercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_percent", MetricsName),
		}, LABELS)
	prometheus.MustRegister(GaugeVecPercent, GaugeVecTotal)
}

func Data() *bssopenapi20171214.QueryResourcePackageInstancesResponseBodyData {
	queryDPUtilizationDetailRequest := &bssopenapi20171214.QueryResourcePackageInstancesRequest{}
	result, err := client.QueryResourcePackageInstances(queryDPUtilizationDetailRequest)
	if err != nil {
		panic(err)
	}
	return result.Body.Data
}

func parseValue(amount, unit string) float64 {
	i, _ := strconv.ParseFloat(amount, 32)
	switch unit {
	case "Byte":
		return i
	case "KB":
		return i * 1024
	case "MB":
		return i * 1024 * 1024
	case "GB":
		return i * 1024 * 1024 * 1024
	case "TB":
		return i * 1024 * 1024 * 1024 * 1024
	default:
		return i
	}
}

func expvarHandler(w http.ResponseWriter, r *http.Request) {
	data := Data()
	instances := data.Instances.Instance
	for _, instance := range instances {
		namePrefix := strings.ToLower(*instance.PackageType)
		labels := map[string]string{
			"region":      *instance.Region,
			"remark":      *instance.Remark,
			"instance_id": *instance.InstanceId,
			"status":      *instance.Status,
			"name":        namePrefix,
		}

		remain := parseValue(*instance.RemainingAmount, *instance.RemainingAmountUnit)
		total := parseValue(*instance.TotalAmount, *instance.TotalAmountUnit)
		GaugeVecTotal.With(labels).Set(remain)
		GaugeVecPercent.With(labels).Set(remain / total)

	}
	handler.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/metrics", expvarHandler)
	http.ListenAndServe(":8080", nil)
}
