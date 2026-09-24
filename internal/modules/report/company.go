package report

import reportpkg "icmongolang/pkg/report"

// CompanyInfo returns the company details used on all report documents.
func CompanyInfo() reportpkg.CompanyInfo {
	return reportpkg.CompanyInfo{
		Name:    "ICMON Auto Repair",
		Address: "123/4 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
		Phone:   "02-123-4567",
		TaxID:   "0123456789012",
	}
}
