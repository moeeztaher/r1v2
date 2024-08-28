package Apis

//import "time"

/*type ApiStatus struct {
    AefIds []string `json:"aefIds"`
}*/

/*type Resource struct {
    ResourceName   string   `json:"resourceName"`
    CommType       string   `json:"commType"`
    Uri            string   `json:"uri"`
    CustOpName     string   `json:"custOpName"`
    CustOperations []CustOp `json:"custOperations"`
    Operations     []string `json:"operations"`
    Description    string   `json:"description"`
}*/

/*type CustOp struct {
    CommType    string   `json:"commType"`
    CustOpName  string   `json:"custOpName"`
    Operations  []string `json:"operations"`
    Description string   `json:"description"`
}*/

/*type Version struct {
    ApiVersion     string     `json:"apiVersion"`
    Expiry         time.Time  `json:"expiry"`
    Resources      []Resource `json:"resources"`
    CustOperations []CustOp   `json:"custOperations"`
}*/

/*type InterfaceDescription struct {
    Ipv4Addr       string   `json:"ipv4Addr"`
    Ipv6Addr       string   `json:"ipv6Addr"`
    Fqdn           string   `json:"fqdn"`
    Port           int      `json:"port"`
    ApiPrefix      string   `json:"apiPrefix"`
    SecurityMethods []string `json:"securityMethods"`
}*/

/*type CivicAddr struct {
    Country   string `json:"country"`
    A1        string `json:"A1"`
}*/

/*type GeoArea struct {
    Shape string `json:"shape"`
    Point struct {
        Lon float64 `json:"lon"`
        Lat float64 `json:"lat"`
    } `json:"point"`
}*/

/*type AefLocation struct {
    CivicAddr CivicAddr `json:"civicAddr"`
    GeoArea   GeoArea   `json:"geoArea"`
    DcId      string    `json:"dcId"`
}*/

/*type ServiceKpis struct {
    MaxReqRate int    `json:"maxReqRate"`
    MaxRestime int    `json:"maxRestime"`
    Availability int  `json:"availability"`
    AvalComp  string `json:"avalComp"`
}*/

/*type IpRange struct {
    Start string `json:"start"`
    End   string `json:"end"`
}*/

/*type UeIpRange struct {
    UeIpv4AddrRanges []IpRange `json:"ueIpv4AddrRanges"`
    UeIpv6AddrRanges []IpRange `json:"ueIpv6AddrRanges"`
}*/

/*type AefProfile struct {
    AefId               string                 `json:"aefId"`
    Versions            []Version              `json:"versions"`
    Protocol            string                 `json:"protocol"`
    DataFormat          string                 `json:"dataFormat"`
    SecurityMethods     []string               `json:"securityMethods"`
    DomainName          string                 `json:"domainName"`
    InterfaceDescriptions []InterfaceDescription `json:"interfaceDescriptions"`
    AefLocation         AefLocation            `json:"aefLocation"`
    ServiceKpis         ServiceKpis            `json:"serviceKpis"`
    UeIpRange           UeIpRange              `json:"ueIpRange"`
}*/

/*type PublishServiceRequest struct {
    ApiName            string       `json:"apiName"`
    ApiId              string       `json:"apiId"`
    ApiStatus          ApiStatus    `json:"apiStatus"`
    AefProfiles        []AefProfile `json:"aefProfiles"`
    Description        string       `json:"description"`
    SupportedFeatures  string       `json:"supportedFeatures"`
    ShareableInfo      struct {
        IsShareable    bool     `json:"isShareable"`
        CapifProvDoms  []string `json:"capifProvDoms"`
    } `json:"shareableInfo"`
    ServiceAPICategory string `json:"serviceAPICategory"`
    ApiSuppFeats       string `json:"apiSuppFeats"`
    PubApiPath         struct {
        CcfIds []string `json:"ccfIds"`
    } `json:"pubApiPath"`
    CcfId              string `json:"ccfId"`
    apiProvName        string `apiProvName: "string"`
}/*

/*type ProblemDetails struct {
    Type            string `json:"type"`
    Title           string `json:"title"`
    Status          int    `json:"status"`
    Detail          string `json:"detail"`
    Instance        string `json:"instance"`
    Cause           string `json:"cause"`
    InvalidParams   []struct {
        Param  string `json:"param"`
        Reason string `json:"reason"`
    } `json:"invalidParams"`
    SupportedFeatures string `json:"supportedFeatures"`
}*/

/*type ServiceInfo struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}*/

/*type YamlInfo struct {
    Info struct {
        Title   string `yaml:"title"`
        Version string `yaml:"version"`
    } `yaml:"info"`
}*/

type PublishServiceAPI struct {
	APIName   string `json:"apiName" bson:"apiName" validate:"required"`
	APIID     string `json:"apiId" bson:"apiId"`
	APIStatus struct {
		AefIds []string `json:"aefIds" bson:"aefIds"`
	} `json:"apiStatus" bson:"apiStatus"`
	AefProfiles []struct {
		AefId    string `json:"aefId" bson:"aefId"`
		Versions []struct {
			APIVersion string `json:"apiVersion" bson:"apiVersion"`
			Expiry     string `json:"expiry" bson:"expiry"`
			Resources  []struct {
				ResourceName   string `json:"resourceName" bson:"resourceName"`
				CommType       string `json:"commType" bson:"commType"`
				Uri            string `json:"uri" bson:"uri"`
				CustOpName     string `json:"custOpName" bson:"custOpName"`
				CustOperations []struct {
					CommType    string   `json:"commType" bson:"commType"`
					CustOpName  string   `json:"custOpName" bson:"custOpName"`
					Operations  []string `json:"operations" bson:"operations"`
					Description string   `json:"description" bson:"description"`
				} `json:"custOperations" bson:"custOperations"`
				Operations  []string `json:"operations" bson:"operations"`
				Description string   `json:"description" bson:"description"`
			} `json:"resources" bson:"resources"`
			CustOperations []struct {
				CommType    string   `json:"commType" bson:"commType"`
				CustOpName  string   `json:"custOpName" bson:"custOpName"`
				Operations  []string `json:"operations" bson:"operations"`
				Description string   `json:"description" bson:"description"`
			} `json:"custOperations" bson:"custOperations"`
		} `json:"versions" bson:"versions"`
		Protocol              string   `json:"protocol" bson:"protocol"`
		DataFormat            string   `json:"dataFormat" bson:"dataFormat"`
		SecurityMethods       []string `json:"securityMethods" bson:"securityMethods"`
		DomainName            string   `json:"domainName" bson:"domainName"`
		InterfaceDescriptions []struct {
			Ipv4Addr        string   `json:"ipv4Addr" bson:"ipv4Addr"`
			Ipv6Addr        string   `json:"ipv6Addr" bson:"ipv6Addr"`
			Fqdn            string   `json:"fqdn" bson:"fqdn"`
			Port            int      `json:"port" bson:"port"`
			ApiPrefix       string   `json:"apiPrefix" bson:"apiPrefix"`
			SecurityMethods []string `json:"securityMethods" bson:"securityMethods"`
		} `json:"interfaceDescriptions" bson:"interfaceDescriptions"`
		AefLocation struct {
			CivicAddr struct {
				Country    string `json:"country" bson:"country"`
				A1         string `json:"A1" bson:"A1"`
				A2         string `json:"A2" bson:"A2"`
				A3         string `json:"A3" bson:"A3"`
				A4         string `json:"A4" bson:"A4"`
				A5         string `json:"A5" bson:"A5"`
				A6         string `json:"A6" bson:"A6"`
				PRD        string `json:"PRD" bson:"PRD"`
				POD        string `json:"POD" bson:"POD"`
				STS        string `json:"STS" bson:"STS"`
				HNO        string `json:"HNO" bson:"HNO"`
				HNS        string `json:"HNS" bson:"HNS"`
				LMK        string `json:"LMK" bson:"LMK"`
				LOC        string `json:"LOC" bson:"LOC"`
				NAM        string `json:"NAM" bson:"NAM"`
				PC         string `json:"PC" bson:"PC"`
				BLD        string `json:"BLD" bson:"BLD"`
				UNIT       string `json:"UNIT" bson:"UNIT"`
				FLR        string `json:"FLR" bson:"FLR"`
				ROOM       string `json:"ROOM" bson:"ROOM"`
				PLC        string `json:"PLC" bson:"PLC"`
				PCN        string `json:"PCN" bson:"PCN"`
				POBOX      string `json:"POBOX" bson:"POBOX"`
				ADDCODE    string `json:"ADDCODE" bson:"ADDCODE"`
				SEAT       string `json:"SEAT" bson:"SEAT"`
				RD         string `json:"RD" bson:"RD"`
				RDSEC      string `json:"RDSEC" bson:"RDSEC"`
				RDBR       string `json:"RDBR" bson:"RDBR"`
				RDSUBBR    string `json:"RDSUBBR" bson:"RDSUBBR"`
				PRM        string `json:"PRM" bson:"PRM"`
				POM        string `json:"POM" bson:"POM"`
				UsageRules string `json:"usageRules" bson:"usageRules"`
				Method     string `json:"method" bson:"method"`
				ProvidedBy string `json:"providedBy" bson:"providedBy"`
			} `json:"civicAddr" bson:"civicAddr"`
			GeoArea struct {
				Shape string `json:"shape" bson:"shape"`
				Point struct {
					Lon float64 `json:"lon" bson:"lon"`
					Lat float64 `json:"lat" bson:"lat"`
				} `json:"point" bson:"point"`
			} `json:"geoArea" bson:"geoArea"`
			DcId string `json:"dcId" bson:"dcId"`
		} `json:"aefLocation" bson:"aefLocation"`
		ServiceKpis struct {
			MaxReqRate   int    `json:"maxReqRate" bson:"maxReqRate"`
			MaxRestime   int    `json:"maxRestime" bson:"maxRestime"`
			Availability int    `json:"availability" bson:"availability"`
			AvalComp     string `json:"avalComp" bson:"avalComp"`
			AvalGraComp  string `json:"avalGraComp" bson:"avalGraComp"`
			AvalMem      string `json:"avalMem" bson:"avalMem"`
			AvalStor     string `json:"avalStor" bson:"avalStor"`
			ConBand      int    `json:"conBand" bson:"conBand"`
		} `json:"serviceKpis" bson:"serviceKpis"`
		UeIpRange struct {
			UeIpv4AddrRanges []struct {
				Start string `json:"start" bson:"start"`
				End   string `json:"end" bson:"end"`
			} `json:"ueIpv4AddrRanges" bson:"ueIpv4AddrRanges"`
			UeIpv6AddrRanges []struct {
				Start string `json:"start" bson:"start"`
				End   string `json:"end" bson:"end"`
			} `json:"ueIpv6AddrRanges" bson:"ueIpv6AddrRanges"`
		} `json:"ueIpRange" bson:"ueIpRange"`
	} `json:"aefProfiles" bson:"aefProfiles"`
	Description       string `json:"description" bson:"description"`
	SupportedFeatures string `json:"supportedFeatures" bson:"supportedFeatures"`
	ShareableInfo     struct {
		IsShareable   bool     `json:"isShareable" bson:"isShareable"`
		CapifProvDoms []string `json:"capifProvDoms" bson:"capifProvDoms"`
	} `json:"shareableInfo" bson:"shareableInfo"`
	ServiceAPICategory string `json:"serviceAPICategory" bson:"serviceAPICategory"`
	ApiSuppFeats       string `json:"apiSuppFeats" bson:"apiSuppFeats"`
	PubApiPath         struct {
		CcfIds []string `json:"ccfIds" bson:"ccfIds"`
	} `json:"pubApiPath" bson:"pubApiPath"`
	CcfId       string `json:"ccfId" bson:"ccfId"`
	ApiProvName string `json:"apiProvName" bson:"apiProvName"`
}

type GetServiceAPI struct {
	APIName        string        `json:"apiName"`
	APIID          string        `json:"apiId"`
	APIStatus      APIStatus     `json:"apiStatus"`
	AEFProfiles    []AEFProfile  `json:"aefProfiles"`
	Description    string        `json:"description"`
	SupportedFeats string        `json:"supportedFeatures"`
	ShareableInfo  ShareableInfo `json:"shareableInfo"`
	ServiceAPICat  string        `json:"serviceAPICategory"`
	APISuppFeats   string        `json:"apiSuppFeats"`
	PubAPIPath     PubAPIPath    `json:"pubApiPath"`
	CCFID          string        `json:"ccfId"`
	APIProvName    string        `json:"apiProvName"`
}

type APIStatus struct {
	AEFIDs []string `json:"aefIds"`
}

type AEFProfile struct {
	AEFID                 string                 `json:"aefId"`
	Versions              []Version              `json:"versions"`
	Protocol              string                 `json:"protocol"`
	DataFormat            string                 `json:"dataFormat"`
	SecurityMethods       []string               `json:"securityMethods"`
	DomainName            string                 `json:"domainName"`
	InterfaceDescriptions []InterfaceDescription `json:"interfaceDescriptions"`
	AEFLocation           AEFLocation            `json:"aefLocation"`
	ServiceKPIs           ServiceKPIs            `json:"serviceKpis"`
	UEIPRange             UEIPRange              `json:"ueIpRange"`
}

type Version struct {
	APIVersion     string          `json:"apiVersion"`
	Expiry         string          `json:"expiry"`
	Resources      []Resource      `json:"resources"`
	CustOperations []CustOperation `json:"custOperations"`
}

type Resource struct {
	ResourceName   string          `json:"resourceName"`
	CommType       string          `json:"commType"`
	URI            string          `json:"uri"`
	CustOpName     string          `json:"custOpName"`
	CustOperations []CustOperation `json:"custOperations"`
	Operations     []string        `json:"operations"`
	Description    string          `json:"description"`
}

type CustOperation struct {
	CommType    string   `json:"commType"`
	CustOpName  string   `json:"custOpName"`
	Operations  []string `json:"operations"`
	Description string   `json:"description"`
}

type InterfaceDescription struct {
	IPv4Addr        string   `json:"ipv4Addr"`
	IPv6Addr        string   `json:"ipv6Addr"`
	FQDN            string   `json:"fqdn"`
	Port            int      `json:"port"`
	APIPrefix       string   `json:"apiPrefix"`
	SecurityMethods []string `json:"securityMethods"`
}

type AEFLocation struct {
	CivicAddr CivicAddr `json:"civicAddr"`
	GeoArea   GeoArea   `json:"geoArea"`
	DCID      string    `json:"dcId"`
}

type CivicAddr struct {
	Country string `json:"country"`
	A1      string `json:"A1"`
}

type GeoArea struct {
	Shape string `json:"shape"`
	Point Point  `json:"point"`
}

type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type ServiceKPIs struct {
	MaxReqRate   int    `json:"maxReqRate"`
	MaxRestime   int    `json:"maxRestime"`
	Availability int    `json:"availability"`
	AvalComp     string `json:"avalComp"`
}

type UEIPRange struct {
	UEIPv4AddrRanges []IPRange `json:"ueIpv4AddrRanges"`
	UEIPv6AddrRanges []IPRange `json:"ueIpv6AddrRanges"`
}

type IPRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ShareableInfo struct {
	IsShareable   bool     `json:"isShareable"`
	CAPIFProvDoms []string `json:"capifProvDoms"`
}

type PubAPIPath struct {
	CCFIDs []string `json:"ccfIds"`
}

type AefProfiles struct {
	AefId string `json:"aefId"`
}

type ApiData struct {
	ApiName     string    `json:"apiName"`
	ApiId       string    `json:"apiId"`
	ApiStatus   APIStatus `json:"apiStatus"`
	AefProfiles []struct {
		AefId    string `json:"aefId"`
		Versions []struct {
			APIVersion string `json:"apiVersion"`
			Expiry     string `json:"expiry"`
			Resources  []struct {
				ResourceName   string `json:"resourceName"`
				CommType       string `json:"commType"`
				Uri            string `json:"uri"`
				CustOpName     string `json:"custOpName"`
				CustOperations []struct {
					CommType    string   `json:"commType"`
					CustOpName  string   `json:"custOpName"`
					Operations  []string `json:"operations"`
					Description string   `json:"description"`
				} `json:"custOperations"`
				Operations  []string `json:"operations"`
				Description string   `json:"description"`
			} `json:"resources"`
			CustOperations []struct {
				CommType    string   `json:"commType"`
				CustOpName  string   `json:"custOpName"`
				Operations  []string `json:"operations"`
				Description string   `json:"description"`
			} `json:"custOperations"`
		} `json:"versions"`
		Protocol              string   `json:"protocol"`
		DataFormat            string   `json:"dataFormat"`
		SecurityMethods       []string `json:"securityMethods"`
		DomainName            string   `json:"domainName"`
		InterfaceDescriptions []struct {
			Ipv4Addr        string   `json:"ipv4Addr"`
			Ipv6Addr        string   `json:"ipv6Addr"`
			Fqdn            string   `json:"fqdn"`
			Port            int      `json:"port"`
			ApiPrefix       string   `json:"apiPrefix"`
			SecurityMethods []string `json:"securityMethods"`
		} `json:"interfaceDescriptions"`
		AefLocation struct {
			CivicAddr struct {
				Country    string `json:"country"`
				A1         string `json:"A1"`
				A2         string `json:"A2"`
				A3         string `json:"A3"`
				A4         string `json:"A4"`
				A5         string `json:"A5"`
				A6         string `json:"A6"`
				PRD        string `json:"PRD"`
				POD        string `json:"POD"`
				STS        string `json:"STS"`
				HNO        string `json:"HNO"`
				HNS        string `json:"HNS"`
				LMK        string `json:"LMK"`
				LOC        string `json:"LOC"`
				NAM        string `json:"NAM"`
				PC         string `json:"PC"`
				BLD        string `json:"BLD"`
				UNIT       string `json:"UNIT"`
				FLR        string `json:"FLR"`
				ROOM       string `json:"ROOM"`
				PLC        string `json:"PLC"`
				PCN        string `json:"PCN"`
				POBOX      string `json:"POBOX"`
				ADDCODE    string `json:"ADDCODE"`
				SEAT       string `json:"SEAT"`
				RD         string `json:"RD"`
				RDSEC      string `json:"RDSEC"`
				RDBR       string `json:"RDBR"`
				RDSUBBR    string `json:"RDSUBBR"`
				PRM        string `json:"PRM"`
				POM        string `json:"POM"`
				UsageRules string `json:"usageRules"`
				Method     string `json:"method"`
				ProvidedBy string `json:"providedBy"`
			} `json:"civicAddr"`
			GeoArea struct {
				Shape string `json:"shape"`
				Point struct {
					Lon float64 `json:"lon"`
					Lat float64 `json:"lat"`
				} `json:"point"`
			} `json:"geoArea"`
			DcId string `json:"dcId"`
		} `json:"aefLocation"`
		ServiceKpis struct {
			MaxReqRate   int    `json:"maxReqRate"`
			MaxRestime   int    `json:"maxRestime"`
			Availability int    `json:"availability"`
			AvalComp     string `json:"avalComp"`
			AvalGraComp  string `json:"avalGraComp"`
			AvalMem      string `json:"avalMem"`
			AvalStor     string `json:"avalStor"`
			ConBand      int    `json:"conBand"`
		} `json:"serviceKpis"`
		UeIpRange struct {
			UeIpv4AddrRanges []struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"ueIpv4AddrRanges"`
			UeIpv6AddrRanges []struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"ueIpv6AddrRanges"`
		} `json:"ueIpRange"`
	} `json:"aefProfiles"`
}

type PatchRequest struct {
	APIStatus       *APIStatus     `json:"apiStatus,omitempty"`
	AEFProfiles     []AEFProfile   `json:"aefProfiles,omitempty"`
	Description     *string        `json:"description,omitempty"`
	ShareableInfo   *ShareableInfo `json:"shareableInfo,omitempty"`
	ServiceCategory *string        `json:"serviceAPICategory,omitempty"`
	APISuppFeats    *string        `json:"apiSuppFeats,omitempty"`
	PubAPIPath      *PubAPIPath    `json:"pubApiPath,omitempty"`
	CCFId           *string        `json:"ccfId,omitempty"`
}

type Rapp struct {
	ApfId              string   `bson:"apf_id"`
	IsAuthorized       bool     `bson:"is_authorized"`
	AuthorizedServices []string `bson:"authorized_services"`
}
