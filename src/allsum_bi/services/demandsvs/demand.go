package demandsvs

import "allsum_bi/models"

func ChangeStatus(demanduuid string, demandstatus int, reportstatus int) (err error) {
	demand, err := models.GetDemandByUuid(demanduuid)
	if err != nil {
		return
	}
	report, err := models.GetReport(demand.Reportid)
	if err != nil {
		return
	}
	demand.Status = demandstatus
	report.Status = reportstatus
	err = models.UpdateDemand(demand, "status")
	if err != nil {
		return
	}
	err = models.UpdateReport(report, "status")
	return
}
