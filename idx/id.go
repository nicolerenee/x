package idx

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jaevor/go-nanoid"
)

const (
	PrefixPartLength = 7
	IDPartLength     = 21
	Parts            = 2
	TotalLength      = PrefixPartLength + IDPartLength + Parts - 1
)

type PrefixedID string

func parts(str string) (string, string) {
	p := strings.SplitN(string(str), "-", Parts-1)

	if len(p) != Parts && len(p[0]) != PrefixPartLength {
		return "", ""
	}

	return p[0], p[1]
}

func (p PrefixedID) Prefix() string {
	prefix, _ := parts(string(p))
	return prefix
}

func (p PrefixedID) id() string {
	_, id := parts(string(p))
	return id
}

func MustNewID(prefix string) PrefixedID {
	id, err := NewID(prefix)
	if err != nil {
		panic(err)
	}

	return id
}

func NewID(prefix string) (PrefixedID, error) {
	prefix = strings.ToLower(prefix)
	if len(prefix) != PrefixPartLength {
		fmt.Println("Problem in NewID: prefix is: " + prefix)
		return "", errors.New("invalid prefix")
	}

	id, err := newIDValue()
	if err != nil {
		return "", err
	}

	return PrefixedID(fmt.Sprintf("%s-%s", prefix, id)), nil
}

func newIDValue() (string, error) {
	id, err := nanoid.Standard(IDPartLength)
	if err != nil {
		return "", err
	}

	return id(), nil
}

// Value implements sql.Valuer so that PrefixedIDs can be written to databases
// transparently. PrefixedIDs map to strings.
func (p PrefixedID) Value() (driver.Value, error) {
	if p.Prefix() == "" {
		return "", errors.New("no prefix set")
	}
	if p.id() == "" {
		return "", errors.New("no id set")
	}

	return string(p), nil
}

func Parse(str string) (PrefixedID, error) {
	prefix, id := parts(str)

	if prefix == "" || id == "" {
		return "", errors.New("invalid id, missing prefix")
	}

	if len(prefix) != 7 {
		return "", errors.New("invalid id, prefix is incorrect length")
	}

	return PrefixedID(str), nil
}

// Scan implements sql.Scanner so PrefixedIDs can be read from databases
// transparently. The value returned is not checked to ensure it's a
// properly formatted PrefixedID.
func (p *PrefixedID) Scan(v any) error {
	if v == nil {
		return fmt.Errorf("expected a value")
	}

	switch src := v.(type) {
	case string:
		*p = PrefixedID(src)
	case []byte:
		*p = PrefixedID(string(src))
	case PrefixedID:
		*p = src
	default:
		return fmt.Errorf("unexpected type, %T", src)
	}

	return nil
}

// MarshalGQL provides GraphQL marshaling so that PrefixedIDs can be returned
// in GraphQL results transparently. Only types that map to a string are supported.
func (p PrefixedID) MarshalGQL(w io.Writer) {
	// io.WriteString(w, strconv.Quote(p.String())) //nolint:errcheck
	// graphql ID is a scalar which must be quoted
	_, _ = io.WriteString(w, strconv.Quote(string(p)))
}

// UnmarshalGQL provides GraphQL unmarshaling so that PrefixedIDs can be parsed
// in GraphQL requests transparently. Only input types that map to a string are supported.
func (p *PrefixedID) UnmarshalGQL(v interface{}) error {
	return p.Scan(v)
}

// Checks to ensure NamespaceID meets the Scanner interface for ent
// var _ field.ValueScanner = NamespaceID{}
// var _ field.ValueScanner = uuid.UUID{}
// var _ field.TextValueScanner = NamespaceID{}
// var _ encoding.TextMarshaler = NamespaceID{}
// var _ encoding.TextUnmarshaler = NamespaceID{}
// var _ driver.Valuer = (*NamespaceID)(nil)
// var _ sql.Scanner = NamespaceID{}
