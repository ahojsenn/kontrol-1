# Rules

R1: AR = Ausgangsrechnungen
R1#S1Partner: Leistung wurde von Partner erbracht
- Partner bekommt 70% seiner Nettoposition
- Kommitment bekommt 25% der Partnernettoposition
- Vertrieb bekommt 5% der Partnernettoposition

R1#S2#Extern: Leistung wurde von Partner erbracht
- Kommitment bekommt 95% der Externnettoposition
- Vertrieb bekommt 5% der Partnernettoposition

# Run all tests

go test ./...

go install bitbucket.org/rwirdemann/kontrol

# Query accounts
```
curl -s http://localhost:8991/kontrol/accounts | python -m json.tool
curl -s http://localhost:8991/kontrol/accounts/AN/bookings | python -m json.tool
```

## Build for different Linux
```
env GOOS=linux GOARCH=amd64 go build bitbucket.org/rwirdemann/kontrol
scp kontrol root@94.130.79.196:~
```
