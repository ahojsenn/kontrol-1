package owner

const (
	PartnerShare             = 0.7
	KommmitmentShare         = 0.25
	KommmitmentExternShare   = 0.95
	KommmitmentEmployeeShare = 0.95
	PartnerProvision         = 0.05
)

var StakeholderRW = Stakeholder{Id: "RW", Name: "Ralf Wirdemann", Type: StakeholderTypePartner}
var StakeholderAN = Stakeholder{Id: "AN", Name: "Anke Nehrenberg", Type: StakeholderTypePartner}
var StakeholderJM = Stakeholder{Id: "JM", Name: "Johannes Mainusch", Type: StakeholderTypePartner}
var StakeholderBW = Stakeholder{Id: "BW", Name: "Ben Wiedenmann", Type: StakeholderTypeEmployee}
var StakeholderKM = Stakeholder{Id: "K", Name: "Kommitment", Type: StakeholderTypeCompany}
var StakeholderEX = Stakeholder{Id: "EX", Name: "Extern", Type: StakeholderTypeExtern}

var AllStakeholder = []Stakeholder{StakeholderRW, StakeholderAN, StakeholderJM, StakeholderBW, StakeholderEX, StakeholderKM}