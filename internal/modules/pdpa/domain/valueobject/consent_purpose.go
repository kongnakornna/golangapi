package valueobject

type ConsentPurpose int8

const (
	ConsentPurposeNecessary ConsentPurpose = 1
	ConsentPurposeMarketing ConsentPurpose = 5
	ConsentPurposeAnalytics ConsentPurpose = 6
)

func (c ConsentPurpose) IsValid() bool {
	switch c {
	case ConsentPurposeNecessary, ConsentPurposeAnalytics, ConsentPurposeMarketing:
		return true
	}
	return false
}

func (c ConsentPurpose) String() string {
	switch c {
	case ConsentPurposeNecessary:
		return "NECESSARY"
	case ConsentPurposeAnalytics:
		return "ANALYTICS"
	case ConsentPurposeMarketing:
		return "MARKETING"
	}
	return "UNKNOWN"
}