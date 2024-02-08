package updator

import (
	"github.com/ThCompiler/go.beget.api/pkg/beget/api/dns"
	"github.com/ThCompiler/go.beget.api/pkg/beget/api/dns/build"
	"github.com/ThCompiler/go.beget.api/pkg/beget/api/dns/result"
	generalresult "github.com/ThCompiler/go.beget.api/pkg/beget/api/result"
	"github.com/ThCompiler/go.beget.api/pkg/beget/core"
	"github.com/ThCompiler/go.beget.api/pkg/beget/core/request"
	"github.com/pkg/errors"
	"update_hostname/internal/logger"
)

const (
	basicPriority = 10
)

type Updater struct {
	getRequest *request.BegetRequest[result.GetData]
	client     core.Client
	log        logger.Interface
	domain     string
}

func NewUpdater(client core.Client, log logger.Interface, domain string) *Updater {
	prepareRequest, err := core.PrepareRequest[result.GetData](client, dns.CallGetData(domain))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "try create request for domain %s", domain))
	}

	return &Updater{
		getRequest: prepareRequest,
		client:     client,
		log:        log,
		domain:     domain,
	}
}

func (updater *Updater) Update() {
	ip, err := getIp()
	if err != nil {
		updater.log.Error(err)
		return
	}

	updater.log.Info("gotten public ip %s", ip)

	data := updater.getDomainInfo()
	if data == nil {
		return
	}

	if data.TypeRecords() != result.Basic || len(data.BasicRecords().A) == 0 || data.BasicRecords().A[0].Address != ip {
		updater.setDomainIp(ip)
	} else {
		updater.log.Info("no changes for ip %s for domain %s", ip, updater.domain)
	}
}

func (updater *Updater) setDomainIp(ip string) {
	req, err := core.PrepareRequest[generalresult.BoolResult](
		updater.client,
		dns.CallChangeRecords(updater.domain,
			dns.SetRecords(
				build.NewBasicRecordsCreator().AddARecords(
					build.NewARecordCreator().AddRecord(basicPriority, ip),
				).Create(),
			),
		),
	)
	if err != nil {
		updater.log.Error(errors.Wrapf(err, "error of creating request to change domain %s ip on %s", updater.domain, ip))
		return
	}

	resp, err := req.Do()
	if err != nil {
		updater.log.Error(errors.Wrapf(err, "error of execution request to change domain %s ip on %s", updater.domain, ip))
		return
	}

	answer, err := resp.GetResult()
	if err != nil {
		updater.log.Error(errors.Wrapf(resp.MustGetError(), "error of request of setting domain %s ip %s", updater.domain, ip))
		return
	}

	result, err := answer.GetResult()
	if err != nil {
		updater.log.Error(errors.Wrapf(answer.MustGetError(), "error of method of setting domain %s ip %s", updater.domain, ip))
		return
	}

	if *result {
		updater.log.Info("ip of domain %s successfully changed to ip %s", updater.domain, ip)
		return
	}

	updater.log.Warn("ip of domain %s not changed to ip %s", updater.domain, ip)
}

func (updater *Updater) getDomainInfo() *result.GetData {
	resp, err := updater.getRequest.Do()
	if err != nil {
		updater.log.Error(errors.Wrapf(err, "try get info for domain %s", updater.domain))
		return nil
	}

	answer, err := resp.GetResult()
	if err != nil {
		updater.log.Error(errors.Wrapf(resp.MustGetError(), "response was with error for %s", updater.domain))
		return nil
	}

	data, err := answer.GetResult()
	if err != nil {
		updater.log.Error(errors.Wrapf(answer.MustGetError(), "method return error for %s", updater.domain))
		return nil
	}

	return data
}
