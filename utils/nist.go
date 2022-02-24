package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PkgInfoMetrics struct {
}

type CVE struct {
	DataType    string `json:"data_type"`
	DataFormat  string `json:"data_format"`
	DataVersion string `json:"data_version"`
	CVEDataMeta struct {
		ID       string `json:"ID"`
		ASSIGNER string `json:"ASSIGNER"`
	} `json:"CVE_data_meta"`
	Problemtype struct {
		ProblemtypeData []struct {
			Description []struct {
				Lang  string `json:"lang"`
				Value string `json:"value"`
			} `json:"description"`
		} `json:"problemtype_data"`
	} `json:"problemtype"`
	References struct {
		ReferenceData []struct {
			URL       string   `json:"url"`
			Name      string   `json:"name"`
			Refsource string   `json:"refsource"`
			Tags      []string `json:"tags"`
		} `json:"reference_data"`
	} `json:"references"`
	Description struct {
		DescriptionData []struct {
			Lang  string `json:"lang"`
			Value string `json:"value"`
		} `json:"description_data"`
	} `json:"description"`
}

type PkgInfo struct {
	ResultsPerPage int `json:"resultsPerPage"`
	StartIndex     int `json:"startIndex"`
	TotalResults   int `json:"totalResults"`
	Result         struct {
		CVEDataType      string `json:"CVE_data_type"`
		CVEDataFormat    string `json:"CVE_data_format"`
		CVEDataVersion   string `json:"CVE_data_version"`
		CVEDataTimestamp string `json:"CVE_data_timestamp"`
		CVEItems         []struct {
			Cve            CVE `json:"cve"`
			Configurations struct {
				CVEDataVersion string `json:"CVE_data_version"`
				Nodes          []struct {
					Operator string        `json:"operator"`
					Children []interface{} `json:"children"`
					CpeMatch []struct {
						Vulnerable          bool          `json:"vulnerable"`
						Cpe23URI            string        `json:"cpe23Uri"`
						VersionEndExcluding string        `json:"versionEndExcluding"`
						CpeName             []interface{} `json:"cpe_name"`
					} `json:"cpe_match"`
				} `json:"nodes"`
			} `json:"configurations"`
			Impact struct {
				BaseMetricV3 struct {
					CvssV3 struct {
						Version               string  `json:"version"`
						VectorString          string  `json:"vectorString"`
						AttackVector          string  `json:"attackVector"`
						AttackComplexity      string  `json:"attackComplexity"`
						PrivilegesRequired    string  `json:"privilegesRequired"`
						UserInteraction       string  `json:"userInteraction"`
						Scope                 string  `json:"scope"`
						ConfidentialityImpact string  `json:"confidentialityImpact"`
						IntegrityImpact       string  `json:"integrityImpact"`
						AvailabilityImpact    string  `json:"availabilityImpact"`
						BaseScore             float64 `json:"baseScore"`
						BaseSeverity          string  `json:"baseSeverity"`
					} `json:"cvssV3"`
					ExploitabilityScore float64 `json:"exploitabilityScore"`
					ImpactScore         float64 `json:"impactScore"`
				} `json:"baseMetricV3"`
				BaseMetricV2 struct {
					CvssV2 struct {
						Version               string  `json:"version"`
						VectorString          string  `json:"vectorString"`
						AccessVector          string  `json:"accessVector"`
						AccessComplexity      string  `json:"accessComplexity"`
						Authentication        string  `json:"authentication"`
						ConfidentialityImpact string  `json:"confidentialityImpact"`
						IntegrityImpact       string  `json:"integrityImpact"`
						AvailabilityImpact    string  `json:"availabilityImpact"`
						BaseScore             float64 `json:"baseScore"`
					} `json:"cvssV2"`
					Severity                string  `json:"severity"`
					ExploitabilityScore     float64 `json:"exploitabilityScore"`
					ImpactScore             float64 `json:"impactScore"`
					AcInsufInfo             bool    `json:"acInsufInfo"`
					ObtainAllPrivilege      bool    `json:"obtainAllPrivilege"`
					ObtainUserPrivilege     bool    `json:"obtainUserPrivilege"`
					ObtainOtherPrivilege    bool    `json:"obtainOtherPrivilege"`
					UserInteractionRequired bool    `json:"userInteractionRequired"`
				} `json:"baseMetricV2"`
			} `json:"impact"`
			PublishedDate    string `json:"publishedDate"`
			LastModifiedDate string `json:"lastModifiedDate"`
		} `json:"CVE_Items"`
	} `json:"result"`
}

// GetPkgRisk - Fetches vulnerability information about a system package from NIST
// given the system package name and version number. Supports paging.
// Parameters:
//		pkgVersionStr (str) - name and semver of package to fetch. Assumes a format of "PACKAGENAME SEMVER".
//		startIndex (int) - for paging and denotes the index to query at
//		limit (int) - for paging and denotes the number of items to fetch. -1 auto pages to fetch all
func GetPkgRisk(pkgVersionStr string, startIndex int, limit int) (PkgInfo, error) {
	req, _ := http.NewRequest("GET", "https://services.nvd.nist.gov/rest/json/cves/1.0", nil)

	q := req.URL.Query()
	q.Add("keyword", pkgVersionStr)
	q.Add("resultsPerPage", strconv.Itoa(limit))
	q.Add("startIndex", strconv.Itoa(startIndex))
	req.URL.RawQuery = q.Encode()

	resp, _ := http.Get(req.URL.String())

	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return PkgInfo{}, err
	} else {
		var result PkgInfo

		if err := json.Unmarshal(body, &result); err != nil {
			return PkgInfo{}, err
		} else {
			return result, nil
		}
	}
}

func GetPkgInfoCPEMetric(pkgInfo PkgInfo) {

}

func SortPkgInfoCVEItems(pkgInfo []CVE) {

}
