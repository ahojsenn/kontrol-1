package processing

import (
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type AusgangsrechnungTestSuite struct {
	suite.Suite
	repository        accountSystem.AccountSystem
	accountBank       *account.Account
	accountHannes     *account.Account
	accountRalf       *account.Account
	accountBen        *account.Account
	accountKommitment *account.Account
}

func (suite *AusgangsrechnungTestSuite) SetupTest() {
	suite.repository = accountSystem.NewDefaultAccountSystem()
	suite.accountBank = suite.repository.GetCollectiveAccount()
	suite.accountRalf, _ = suite.repository.Get(  valueMagnets.StakeholderRepository{}.Get("RW").Id )
	suite.accountHannes, _ = suite.repository.Get(  valueMagnets.StakeholderRepository{}.Get("JM").Id )
	suite.accountBen, _ = suite.repository.Get(  valueMagnets.StakeholderRepository{}.Get("BW").Id )
	suite.accountKommitment, _ = suite.repository.Get(valueMagnets.StakeholderRepository{}.Get("K").Id)
}

func TestAusgangsRechnungTestSuite(t *testing.T) {
	suite.Run(t, new(AusgangsrechnungTestSuite))
}

/*
func (suite *AusgangsrechnungTestSuite) TestPartnerNettoAnteil() {

	// given: a booking
	net := make(map[valueMagnets.Stakeholder]float64)
	net[valueMagnets.StakeholderRepository{}.Get("RW")] = 10800.0
	net[valueMagnets.StakeholderRepository{}.Get("JM")] = 3675.0
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13,"AR", "", "", "JM", net, 17225.25, "Rechnung 1234", 1, 2017, its2018)

	// when: the position is processed
	Process(suite.repository, *p)

	// then ralf 1 booking: his own net share
	accountRalf, _ := suite.repository.Get(valueMagnets.StakeholderRepository{}.Get("RW").Id)
	bookingsRalf := accountRalf.Bookings
	suite.Equal(1, len(bookingsRalf))
	bRalf, _ := findBookingByText(bookingsRalf, "Rechnung 1234#NetShare#RW")
	suite.InDelta(10800.0*PartnerShare, bRalf.Amount, 0.01)
	suite.Equal(1, bRalf.Month)
	suite.Equal(2017, bRalf.Year)
	suite.Equal(booking.CC_Nettoanteil, bRalf.Type)

	// and hannes got 3 bookings: his own net share and 2 provisions ostCenter booking
	accountHannes, _ := suite.repository.Get(valueMagnets.StakeholderRepository{}.Get("JM").Id)
	bookingsHannes := accountHannes.Bookings
	suite.Equal(3, len(bookingsHannes))

	// net share
	b, _ := findBookingByText(bookingsHannes, "Rechnung 1234#NetShare#JM")
	suite.Equal(3675.0*PartnerShare, b.Amount)
	suite.Equal(1, b.Month)
	suite.Equal(2017, b.Year)

	// provision from ralf
	provisionRalf, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#RW")
	suite.Equal(10800.0*PartnerProvision, provisionRalf.Amount)
	suite.Equal(booking.CC_Vertriebsprovision, provisionRalf.Type)

	// // provision from hannes
	provisionHannes, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#JM")
	suite.Equal(3675.0*PartnerProvision, provisionHannes.Amount)
	suite.Equal(booking.CC_Vertriebsprovision, provisionHannes.Type)

	// kommitment got 25% from ralfs net booking
	accountKommitment, _ := suite.repository.Get(valueMagnets.StakeholderKM.Id)
	bookingsKommitment := accountKommitment.Bookings
	suite.Equal(2, len(bookingsKommitment))
	kommitmentRalf, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#RW")
	suite.Equal(10800.0*KommmitmentShare, kommitmentRalf.Amount)
	suite.Equal(booking.CC_Kommitmentanteil, kommitmentRalf.Type)

	// and kommitment got 25% from hannes net booking
	kommitmentHannes, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#JM")
	suite.Equal(3675.0*KommmitmentShare, kommitmentHannes.Amount)
	suite.Equal(booking.CC_Kommitmentanteil, kommitmentHannes.Type)
}
*/


//
// Tests für Offene Posten
//
func (suite *AusgangsrechnungTestSuite) TestOffeneRechnung() {

	// given: a booking with empty timestamp in position "BankCreated"
	net := make(map[valueMagnets.Stakeholder]float64)
	net[valueMagnets.StakeholderRepository{}.Get("RW")] = 10800.0
	net[valueMagnets.StakeholderRepository{}.Get("JM")] = 3675.0
	p := booking.NewBooking(13, "AR", "", "", "JM", "Project-X", net, 17225.25, "Rechnung 1234", 1, 2017, time.Time{})

	// when: the position is processed
	Process(suite.repository, *p)

	// then the booking is on SKR03_1400
	account1400, _ := suite.repository.Get(accountSystem.SKR03_1400.Id)
	bookings1400 := account1400.Bookings
	suite.Equal(1, len(bookings1400))

	// the booking is booked to partners
	/*
	accountHannes, _ := suite.repository.Get(valueMagnets.StakeholderRepository{}.Get("JM").Id)
	bookingsHannes := accountHannes.Bookings
	suite.Equal(3, len(bookingsHannes))
	*/
}



//
// Tests für CC_Vertriebsprovision
//

// - Kommitment bekommt den 95% der Nettoposition
// - Dealbringer ist Partner => Partner bekommt je 5% der Nettoposition(en)
/*
func (suite *AusgangsrechnungTestSuite) TestDealbringerIstPartner() {

	// Eine Buchung mit 2 Nettopositionen
	net := make(map[valueMagnets.Stakeholder]float64)
	net[valueMagnets.StakeholderRepository{}.Get("RW")] = 10800.0
	net[valueMagnets.StakeholderRepository{}.Get("JM")] = 3675.0
	dealbringer := "JM"
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.Ausgangsrechnung(13, dealbringer, net, 17225.25, "Rechnung 1234", 1, 2017, its2018)

	Process(suite.repository, *p)

	// Ralfs CC_Nettoanteil
	provisionRalf, _ := findBookingByText(suite.accountRalf.Bookings, "Rechnung 1234#NetShare#RW")
	suite.assertBooking(10800.00*PartnerShare, booking.CC_Nettoanteil, provisionRalf)

	// Hannes bekommt Provision für Hannes CC_Nettoanteil
	provisionHannes, _ := findBookingByText(suite.accountHannes.Bookings, "Rechnung 1234#Provision#JM")
	suite.assertBooking(3675.0*PartnerProvision, booking.CC_Vertriebsprovision, provisionHannes)
}

// - Kommitment bekommt den 95% der Nettoposition
// - Dealbringer ist Angestellter => Angestellter bekommt 5% der Nettoposition,
//   Kostenstelle Dealbringer
func (suite *AusgangsrechnungTestSuite) TestDealbringerIstAngestellter() {

	// Given a booking where dealbringes is an employee
	net := make(map[valueMagnets.Stakeholder]float64)
	net[valueMagnets.StakeholderRepository{}.Get("RW")] = 10800.0
	dealbringer := "BW"
	its2017 := time.Date(2017, 1, 23, 0, 0, 0, 0, time.UTC)

	// when booked
	b := booking.Ausgangsrechnung(13, dealbringer, net, 17225.25, "Rechnung 1234", 1, 2017, its2017)

	Process(suite.repository, *b)

	// Provision ist auf Ben-Account gebucht
	provision, err := findBookingByText(suite.accountBen.Bookings, "Rechnung 1234#Provision#RW")
	suite.Nil(err)
	suite.NotEqual(provision, nil)
	suite.assertBooking(10800.0*PartnerProvision, booking.CC_Vertriebsprovision, provision)
	suite.Equal("BW", provision.CostCenter)
}
*/

func (suite *AusgangsrechnungTestSuite) assertBooking(amount float64, _type string, b *booking.Booking) {
	suite.Equal(amount, b.Amount)
	suite.Equal(_type, b.Type)
}

func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}