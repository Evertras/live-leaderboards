package repo

const (
	keyPrimary = "pk"
	keyRange   = "sk"
)

type primaryKey struct {
	PK string `dynamodbav:"pk" json:"-"`
	SK string `dynamodbav:"sk" json:"-"`
}
