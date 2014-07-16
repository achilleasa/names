// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package names

import (
	"fmt"
	"strings"
)

// A Tag tags things that are taggable.
type Tag interface {
	// Kind returns the kind of the tag.
	// This method is for legacy compatibility, callers should
	// use equality or type assertions to verify the Kind, or type
	// of a Tag.
	Kind() string

	// Id returns an identifier for this Tag.
	// The contents and format of the identifier are specific
	// to the implementer of the Tag.
	Id() string

	fmt.Stringer // all Tags should be able to print themselves
}

// TagKind returns one of the *TagKind constants for the given tag, or
// an error if none matches.
func TagKind(tag string) (string, error) {
	i := strings.Index(tag, "-")
	if i <= 0 || !validKinds(tag[:i]) {
		return "", fmt.Errorf("%q is not a valid tag", tag)
	}
	return tag[:i], nil
}

func validKinds(kind string) bool {
	switch kind {
	case UnitTagKind, MachineTagKind, ServiceTagKind, EnvironTagKind, UserTagKind, RelationTagKind, NetworkTagKind, ActionTagKind:
		return true
	}
	return false
}

func splitTag(tag string) (string, string, error) {
	kind, err := TagKind(tag)
	if err != nil {
		return "", "", err
	}
	return kind, tag[len(kind)+1:], nil
}

// ParseTag parses a string representation into a Tag.
func ParseTag(tag string) (Tag, error) {
	kind, id, err := splitTag(tag)
	if err != nil {
		return nil, invalidTagError(tag, "")
	}
	switch kind {
	case UnitTagKind:
		id = unitTagSuffixToId(id)
		if !IsUnit(id) {
			return nil, invalidTagError(tag, kind)
		}
		return NewUnitTag(id), nil
	case MachineTagKind:
		id = machineTagSuffixToId(id)
		if !IsMachine(id) {
			return nil, invalidTagError(tag, kind)
		}
		return NewMachineTag(id), nil
	case ServiceTagKind:
		if !IsService(id) {
			return nil, invalidTagError(tag, kind)
		}
		return NewServiceTag(id), nil
	case UserTagKind:
		if !IsUser(id) {
			return nil, invalidTagError(tag, kind)
		}
		return NewUserTag(id), nil
	case EnvironTagKind:
		if !IsEnvironment(id) {
			return nil, invalidTagError(tag, kind)
		}
		return NewEnvironTag(id), nil
	case RelationTagKind:
		id = relationTagSuffixToKey(id)
		if !IsRelation(id) {
			return nil, invalidTagError(tag, kind)
		}
		return NewRelationTag(id), nil
	case NetworkTagKind:
		if !IsNetwork(id) {
			return nil, invalidTagError(tag, kind)
		}
		return NewNetworkTag(id), nil
	case ActionTagKind:
		if !IsAction(id) {
			return nil, invalidTagError(tag, kind)
		}
		return NewActionTag(id), nil
	default:
		return nil, invalidTagError(tag, "")
	}
}

func invalidTagError(tag, kind string) error {
	if kind != "" {
		return fmt.Errorf("%q is not a valid %s tag", tag, kind)
	}
	return fmt.Errorf("%q is not a valid tag", tag)
}
