package dns

type Zone struct {
	Id     string
	Name   string
	Status string
}

type Record struct {
	Id     string
	Type   string
	Name   string
	Value  string
	TTL    int
	ZoneId string
}

type Provider interface {
	RequestZone(zoneName string) (Zone, error)
	RequestRecord(zone Zone, recordName string, recordType string) (Record, error)
	UpdateZoneRecord(zone Zone, record Record) (Record, error)
}

type ZoneNotFoundError struct {
	ZoneName string
}

func (e *ZoneNotFoundError) Error() string {
	return "The zone " + e.ZoneName + " was not found"
}

type RecordNotFoundError struct {
	Zone       Zone
	RecordName string
	RecordType string
}

func (e *RecordNotFoundError) Error() string {
	return "The record " + e.RecordName + " with type " + e.RecordType + " was not found in the zone " + e.Zone.Name
}

type UpdateNotPossibleError struct {
	Zone   Zone
	Record Record
}

func (e *UpdateNotPossibleError) Error() string {
	return "The update of " + e.Record.Value + " in zone " + e.Zone.Name + " is not possible"
}
