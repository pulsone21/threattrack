package entities

type Entity interface {
	ScanTo(ScanFunc) error
}
