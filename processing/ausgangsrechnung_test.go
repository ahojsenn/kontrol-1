package processing

import (
	"testing"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/booking"
	"bitbucket.org/rwirdemann/kontrol/owner"
	"github.com/stretchr/testify/suite"
)

type AusgangsrechnungTestSuite struct {
	suite.Suite
	repository        account.Repository
	accountBank       *account.Account
	accountHannes     *account.Account
	accountKommitment *account.Account
}

func (suite *AusgangsrechnungTestSuite) SetupTest() {
	suite.repository = account.NewDefaultRepository()
	suite.accountBank = suite.repository.BankAccount()
	suite.accountHannes, _ = suite.repository.Get(owner.StakeholderJM.Id)
	suite.accountKommitment, _ = suite.repository.Get(owner.StakeholderKM.Id)
}

func TestAusgangsRechnungTestSuite(t *testing.T) {
	suite.Run(t, new(AusgangsrechnungTestSuite))
}

func (suite *AusgangsrechnungTestSuite) TestPartnerNettoAnteil() {

	// given: a booking
	net := make(map[owner.Stakeholder]float64)
	net[owner.StakeholderRW] = 10800.0
	net[owner.StakeholderJM] = 3675.0
	p := booking.NewBooking("AR", "JM", net, 17225.25, "Rechnung 1234", 1, 2017)

	// when: the position is processed
	Process(suite.repository, *p)

	// then ralf 1 booking: his own net share
	accountRalf, _ := suite.repository.Get(owner.StakeholderRW.Id)
	bookingsRalf := accountRalf.Bookings
	suite.Equal(1, len(bookingsRalf))
	bRalf, _ := findBookingByText(bookingsRalf, "Rechnung 1234#NetShare#RW")
	suite.InDelta(10800.0*owner.PartnerShare, bRalf.Amount, 0.01)
	suite.Equal(1, bRalf.Month)
	suite.Equal(2017, bRalf.Year)
	suite.Equal(booking.Nettoanteil, bRalf.Type)

	// and hannes got 3 bookings: his own net share and 2 provisions
	accountHannes, _ := suite.repository.Get(owner.StakeholderJM.Id)
	bookingsHannes := accountHannes.Bookings
	suite.Equal(3, len(bookingsHannes))

	// net share
	b, _ := findBookingByText(bookingsHannes, "Rechnung 1234#NetShare#JM")
	suite.Equal(3675.0*owner.PartnerShare, b.Amount)
	suite.Equal(1, b.Month)
	suite.Equal(2017, b.Year)

	// provision from ralf
	provisionRalf, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#RW")
	suite.Equal(10800.0*owner.PartnerProvision, provisionRalf.Amount)
	suite.Equal(booking.Vertriebsprovision, provisionRalf.Type)

	// // provision from hannes
	provisionHannes, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#JM")
	suite.Equal(3675.0*owner.PartnerProvision, provisionHannes.Amount)
	suite.Equal(booking.Vertriebsprovision, provisionHannes.Type)

	// kommitment got 25% from ralfs net booking
	accountKommitment, _ := suite.repository.Get(owner.StakeholderKM.Id)
	bookingsKommitment := accountKommitment.Bookings
	suite.Equal(2, len(bookingsKommitment))
	kommitmentRalf, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#RW")
	suite.Equal(10800.0*owner.KommmitmentShare, kommitmentRalf.Amount)
	suite.Equal(booking.Kommitmentanteil, kommitmentRalf.Type)

	// and kommitment got 25% from hannes net booking
	kommitmentHannes, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#JM")
	suite.Equal(3675.0*owner.KommmitmentShare, kommitmentHannes.Amount)
	suite.Equal(booking.Kommitmentanteil, kommitmentHannes.Type)
}
